# Blog 后端项目协作说明

## 1. 项目定位

这是一个使用 Go 编写的博客后端接口项目，当前采用 `Gin + GORM + MySQL + Redis` 的常见分层架构。项目已经不只是“用户登录示例”，而是包含了博客业务的核心闭环：

- 用户注册、登录、登出、刷新 token
- 邮箱验证码与找回密码
- 图形验证码
- 文章草稿、发布、列表、详情、点赞、收藏、浏览量
- 评论与楼中楼回复
- 分类列表
- 管理员文章审核

从代码实际实现来看，这个项目已经是一套完整的博客 API 后端骨架，适合继续迭代功能。

## 2. 技术栈

- Web 框架：`github.com/gin-gonic/gin`
- ORM：`gorm.io/gorm` + `gorm.io/driver/mysql`
- 缓存：`github.com/redis/go-redis/v9`
- 配置：`github.com/spf13/viper`
- 日志：`zap` + `lumberjack`
- 认证：JWT
- 密码加密：`bcrypt`
- 邮件发送：`gomail`
- 验证码：`base64Captcha`

`go.mod` 当前声明为 `go 1.25.0`，实际协作时建议确认本地 Go 版本是否匹配。

## 3. 项目结构

### 入口

- `cmd/server/main.go`

职责很单一：

1. 创建 `app.NewApp()`
2. 调用 `Initialize()`
3. 调用 `Run()`

也就是说，这个项目的真正启动逻辑都集中在 `internal/app/app.go`。

### 核心目录

- `internal/app`
  - 应用初始化、依赖注入、HTTP 服务启动、优雅关闭
- `internal/api`
  - 控制器和路由注册
- `internal/service`
  - 业务逻辑层
- `internal/repository`
  - 数据访问层
- `internal/model/entity`
  - GORM 实体
- `internal/model/dto/request`
  - 请求参数结构体
- `internal/model/dto/response`
  - 返回结构体
- `internal/middleware`
  - 日志、恢复、跨域、JWT 鉴权
- `pkg/config`
  - 配置加载
- `pkg/database`
  - MySQL / Redis 初始化与关闭
- `pkg/logger`
  - 日志封装
- `pkg/jwt`
  - JWT 生成与解析
- `pkg/response`
  - 统一响应格式
- `pkg/errors`
  - 业务错误码与错误对象
- `configs`
  - 配置样例
- `scripts`
  - SQL 初始化脚本

## 4. 启动流程

`internal/app/app.go` 的初始化顺序很重要：

1. `initConfig()`
2. `initLogger()`
3. `initDatabase()`
4. `initDependencies()`
5. `initRouter()`
6. `initServer()`

### 关键点

- `initDatabase()` 会初始化 MySQL，并执行 `AutoMigrate`
- 当前自动迁移的实体包括：
  - `User`
  - `Category`
  - `Article`
  - `ArticleCategory`
  - `ReviewLog`
  - `Comment`
  - `Like`
  - `Favorite`
  - `Image`
- Redis 初始化失败不会阻止服务启动，但会影响依赖 Redis 的功能
  - 比如邮箱验证码、注册、重置密码

### 运行方式

常用命令：

```bash
go run cmd/server/main.go
```

或者：

```bash
make run
```

## 5. 配置说明

配置结构定义在 `pkg/config/config.go`，加载逻辑在 `pkg/config/loader.go`。

### 配置文件

默认从以下位置寻找 `config.yaml`：

- `./configs`
- 当前目录 `.`

仓库里目前看到的是：

- `configs/config.yaml.example`

### 配置项

- `app`
  - 应用名、版本、运行模式、端口、AES 密钥
- `database.mysql`
  - MySQL 连接信息和连接池参数
- `database.redis`
  - Redis 连接信息
- `jwt`
  - JWT 密钥和过期时间
- `log`
  - 日志等级、文件、滚动配置
- `cors`
  - 跨域配置
- `email`
  - 邮件服务器配置

### 环境变量覆盖

这里有一个容易忽略的细节：

- Viper 设置了前缀 `blog`
- 但代码又额外直接读取了这几个环境变量：
  - `MYSQL_PASSWORD`
  - `REDIS_PASSWORD`
  - `JWT_SECRET`

所以实际部署时，这三个敏感值可以直接用上面的名字覆盖，不一定要写成 `BLOG_...`。

## 6. 分层约定

这个项目的代码组织比较标准，遵循：

`Controller -> Service -> Repository`

### Controller 层

位于 `internal/api/*`，负责：

- 接收 HTTP 请求
- 参数绑定与校验
- 调用 service
- 返回统一格式响应

### Service 层

位于 `internal/service/*`，负责：

- 编排业务逻辑
- 权限判断
- 状态流转
- 构建 response DTO

### Repository 层

位于 `internal/repository/*`，负责：

- GORM 查询
- 事务操作
- 聚合计数更新
- 多表关联读取

### 一个非常重要的协作原则

新增功能时，尽量不要跳层调用。

比如：

- Controller 不要直接操作数据库
- Service 不要依赖 Gin 上下文
- Repository 不要做业务含义判断

当前项目整体上是遵守这个原则的，后续扩展也应保持一致。

## 7. 统一响应与错误处理

### 统一响应格式

定义在 `pkg/response/response.go`：

```json
{
  "code": 0,
  "msg": "成功",
  "data": {}
}
```

### 特别注意

项目里的业务错误通常仍然返回 HTTP `200`，真正的结果要看 JSON 中的：

- `code`
- `msg`

这意味着：

- 前端不能只看 HTTP 状态码
- 联调时要明确区分“HTTP 成功”和“业务成功”

### 常见错误码

定义在 `pkg/errors/code.go`，例如：

- `0`：成功
- `400`：请求参数错误
- `401`：未授权
- `403`：禁止访问
- `404`：资源不存在
- `409`：资源冲突
- `1001`：用户不存在
- `1002`：用户已存在
- `1003`：用户名或密码错误
- `1004`：用户被禁用
- `1005`：无效 token
- `1006`：token 过期
- `1007`：验证码无效

业务错误类型使用 `pkg/errors/errors.go` 中的 `BizError`。

## 8. 认证与权限

### JWT 鉴权

中间件在 `internal/middleware/auth.go`。

请求头格式：

```http
Authorization: Bearer <token>
```

中间件会把以下字段写入 Gin context：

- `user_id`
- `username`

### 角色约定

从 `internal/model/entity/user.go` 和 `internal/service/review_service.go` 看：

- `role = 0`：管理员
- `role = 1`：普通用户

管理员审核接口只是挂了 JWT 中间件，真正的管理员校验是在 `reviewService.assertAdmin()` 里做的。

也就是说：

- 路由层只校验“是否登录”
- Service 层才校验“是否管理员”

## 9. 主要业务模块

### 9.1 认证模块 `internal/api/auth`

路由前缀：`/api/auth`

主要接口：

- `POST /register`
- `POST /login`
- `POST /refresh`
- `POST /logout`
- `POST /email/code`
- `POST /password/reset`
- `GET /captcha`

实现细节：

- 登录前需要校验图形验证码
- 注册前需要先发送邮箱验证码，并把验证码写入 Redis
- 重置密码同样依赖 Redis 中的邮箱验证码
- `Logout` 当前基本是空实现，没有 token 黑名单机制

这意味着当前系统的登出更接近客户端删除 token，而不是服务端强制失效。

### 9.2 用户模块 `internal/api/user`

路由前缀：`/api/user`

这些接口都需要登录：

- `GET /favorites`
- `GET /profile`
- `PUT /profile`
- `POST /password`
- `GET /list`

### 9.3 文章模块 `internal/api/article`

路由前缀：`/api/articles`

主要接口：

- `POST /api/articles` 创建草稿
- `GET /api/articles` 文章列表
- `GET /api/articles/mine` 我的文章
- `GET /api/articles/:id/comments` 文章评论树
- `GET /api/articles/:id` 文章详情
- `PUT /api/articles/:id` 更新草稿
- `POST /api/articles/:id/publish` 提交发布
- `DELETE /api/articles/:id` 删除文章
- `POST /api/articles/:id/view` 增加浏览量
- `POST /api/articles/:id/like` 点赞
- `POST /api/articles/:id/unlike` 取消点赞
- `POST /api/articles/:id/favorite` 收藏
- `POST /api/articles/:id/unfavorite` 取消收藏

状态流转从代码可见：

- `0`：草稿
- `1`：待审核
- `2`：已发布
- `3`：已拒绝
- `4`：已封禁

一个值得注意的实现细节：

- 文章详情接口 `GET /api/articles/:id` 当前挂了 `middleware.Auth()`
- Service 中允许“已发布文章”或者“作者本人”查看

也就是说，这个接口虽然支持看已发布文章，但代码层面仍要求带 token。若后续要开放游客访问，需要改路由或控制器逻辑。

### 9.4 评论模块 `internal/api/comment`

路由前缀：`/api/comments`

这些接口都要求登录：

- `POST /api/comments`
- `POST /api/comments/replies`
- `POST /api/comments/:id/like`
- `POST /api/comments/:id/unlike`
- `DELETE /api/comments/:id`

实现特征：

- 支持楼中楼
- 一级评论 `parent_id` 为空
- 子评论会校验 `root_id` 是否与评论树一致
- 删除评论前必须先删子评论，否则会被拒绝

### 9.5 分类模块 `internal/api/category`

路由：

- `GET /api/categories`

目前是只读分类列表。

### 9.6 管理审核模块 `internal/api/super`

路由前缀：`/api/super`

需要登录，且必须是管理员：

- `GET /api/super/articles`
- `GET /api/super/articles/:id`
- `POST /api/super/articles/:id/approve`
- `POST /api/super/articles/:id/reject`
- `POST /api/super/articles/:id/ban`
- `PUT /api/super/articles/:id/category`

审核逻辑集中在 `review_service.go` 和 `review_repository.go`。

## 10. 数据模型重点

### User

核心字段：

- `username`
- `password`
- `email`
- `nickname`
- `avatar`
- `bio`
- `role`
- `status`

### Article

核心字段：

- `user_id`
- `title`
- `content`
- `summary`
- `cover_image`
- `reject_reason`
- `status`
- `view_count`
- `like_count`
- `favorite_count`
- `comment_count`

并且和 `Category` 是多对多关系，通过 `article_categories` 中间表关联。

### Comment

核心字段：

- `article_id`
- `user_id`
- `parent_id`
- `root_id`
- `content`
- `status`
- `like_count`
- `reply_count`

### Category

核心字段：

- `name`
- `slug`

## 11. 中间件

全局中间件在 `internal/api/router.go` 中注册：

- `Recovery()`
- `Logger()`
- `CORS()`

健康检查接口：

- `GET /api/health`

注意这里实际代码中的健康检查是 `/api/health`，不是文档里常见的 `/api/v1/health`。

## 12. 新增功能时应该改哪里

如果要新增一个完整业务模块，建议按下面顺序做：

1. 在 `internal/model/entity` 增加实体
2. 在 `internal/model/dto/request` 增加请求 DTO
3. 在 `internal/model/dto/response` 增加响应 DTO
4. 在 `internal/repository` 增加接口与实现
5. 在 `internal/service` 增加接口与实现
6. 在 `internal/api/<module>` 增加 controller 与 routes
7. 在 `internal/app/app.go` 的 `initDependencies()` 中注入依赖
8. 在 `internal/api/router.go` 中注册路由
9. 如需建表，加入 `AutoMigrate(...)`

这个项目的依赖注入目前是手工写在 `app.go` 里的，不是自动 DI 容器，所以新增模块时最容易漏掉的地方就是：

- `initDependencies()`
- `router.go`
- `AutoMigrate()`

## 13. 已知协作注意点

### 文档存在编码问题

仓库里的 `README.md`、`docs/*.md` 内容大意还能看出来，但终端读取时有明显乱码。后续如果要继续维护文档，建议统一保存为 UTF-8 编码并重新整理。

### Redis 不是“完全可选”

虽然 Redis 初始化失败不会阻止服务启动，但以下功能会直接受影响：

- 邮箱验证码发送
- 注册
- 重置密码

### 登出未做服务端失效

当前 `Logout` 没有 token 黑名单机制。如果后续有安全要求，建议补充：

- Redis 黑名单
- token 版本号
- refresh token 机制

### 业务错误返回 200

这是当前项目的明确约定，但如果未来要接第三方平台或统一网关，可能要评估是否改成更标准的 HTTP 状态码方案。

## 14. 一句话总结

这个 blog 项目已经具备“可继续开发的博客后端雏形”：分层清晰、依赖注入集中、鉴权与审核流程完整，后续扩展时最值得遵守的是现有的分层边界、统一响应规范，以及 `app.go` 中的初始化入口约定。
