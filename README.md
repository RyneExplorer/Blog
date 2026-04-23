# Blog API

基于 Go 构建的博客后端项目，采用 `Gin + GORM + MySQL + Redis` 的分层架构。当前实现已覆盖博客系统的核心闭环，不是简单的登录示例，而是一套可继续扩展的 API 后端骨架。

## 项目概览

当前已实现的主要能力：

- 用户注册、登录、刷新 token、登出
- 图形验证码、邮箱验证码、找回密码
- 文章草稿、发布、详情、分页、浏览量、点赞、收藏
- 评论树、楼中楼回复、评论点赞
- 分类列表
- 管理员文章审核

技术栈：

- Web 框架：`gin`
- ORM：`gorm`
- 数据库：MySQL
- 缓存：Redis
- 配置管理：`viper`
- 日志：`zap` + `lumberjack`
- 认证：JWT
- 邮件发送：`gomail`

## 目录结构

```text
blog/
├── cmd/server/main.go              # 程序入口
├── internal/
│   ├── app/                        # 应用初始化与依赖注入
│   ├── api/                        # 控制器与路由
│   ├── service/                    # 业务逻辑
│   ├── repository/                 # 数据访问
│   ├── model/entity/               # GORM 实体
│   ├── model/dto/request/          # 请求 DTO
│   ├── model/dto/response/         # 响应 DTO
│   └── middleware/                 # 日志、恢复、CORS、鉴权
├── pkg/
│   ├── config/                     # 配置加载
│   ├── database/                   # MySQL / Redis 初始化
│   ├── logger/                     # 日志封装
│   ├── jwt/                        # JWT 工具
│   ├── response/                   # 统一响应结构
│   └── errors/                     # 业务错误定义
├── configs/                        # 配置文件
├── docs/                           # 补充文档
├── uploads/                        # 上传资源
├── Makefile
└── README.md
```

## 架构说明

项目遵循标准三层结构：

`Controller -> Service -> Repository`

- Controller 负责 HTTP 请求处理、参数绑定与返回响应
- Service 负责业务编排、权限判断、状态流转
- Repository 负责 GORM 查询、事务和数据持久化

启动逻辑集中在 `internal/app/app.go`，初始化顺序为：

1. `initConfig()`
2. `initLogger()`
3. `initDatabase()`
4. `initDependencies()`
5. `initRouter()`
6. `initServer()`

新增模块时，最容易漏改的地方是：

- `internal/app/app.go` 中的依赖注入
- `internal/api/router.go` 中的路由注册
- `AutoMigrate(...)` 中的实体迁移

## 快速开始

### 环境要求

- Go `1.25.0` 或与 `go.mod` 兼容的版本
- MySQL `8.0+`
- Redis `6.0+`

### 安装与运行

```bash
go mod tidy
cp configs/config.yaml.example configs/config.yaml
go run cmd/server/main.go
```

或使用 Makefile：

```bash
make run
```

默认会从以下位置查找配置文件：

- `./configs/config.yaml`
- 当前目录下的 `config.yaml`

敏感配置支持通过环境变量覆盖：

- `MYSQL_PASSWORD`
- `REDIS_PASSWORD`
- `JWT_SECRET`

## 常用命令

```bash
make build          # 编译到 bin/blog-api.exe
make run            # 启动服务
make test           # 运行测试
make test-coverage  # 生成覆盖率报告
make fmt            # 格式化代码
make vet            # 静态检查
make deps           # 整理并下载依赖
make clean          # 清理构建产物
```

## 核心接口概览

健康检查：

- `GET /api/health`

认证模块：

- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/refresh`
- `POST /api/auth/logout`
- `POST /api/auth/email/code`
- `POST /api/auth/password/reset`
- `GET /api/auth/captcha`

用户模块：

- `GET /api/user/profile`
- `PUT /api/user/profile`
- `POST /api/user/avatar`
- `POST /api/user/password`

文章模块：

- `POST /api/articles`
- `POST /api/articles/cover_image`
- `GET /api/articles`
- `GET /api/articles/favorites`
- `GET /api/articles/mine`
- `GET /api/articles/mine/:id`
- `GET /api/articles/:id/comments`
- `GET /api/articles/:id`
- `PUT /api/articles/:id`
- `POST /api/articles/:id/publish`
- `DELETE /api/articles/:id`
- `POST /api/articles/:id/view`
- `POST /api/articles/:id/like`
- `POST /api/articles/:id/unlike`
- `POST /api/articles/:id/favorite`
- `POST /api/articles/:id/unfavorite`

评论模块：

- `POST /api/comments`
- `POST /api/comments/replies`
- `POST /api/comments/:id/like`
- `POST /api/comments/:id/unlike`
- `DELETE /api/comments/:id`

管理审核模块：

- `GET /api/super/articles`
- `GET /api/super/articles/:id`
- `GET /api/super/userlist`
- `POST /api/super/articles/:id/approve`
- `POST /api/super/articles/:id/reject`
- `POST /api/super/articles/:id/ban`
- `PUT /api/super/articles/:id/category`

## 响应与权限约定

统一响应定义在 `pkg/response/response.go`，格式类似：

```json
{
  "code": 0,
  "msg": "成功",
  "data": {}
}
```

需要注意：

- 项目中的业务错误很多情况下仍返回 HTTP `200`
- 前端或调用方需要同时判断 `code` 和 `msg`
- JWT 鉴权从 `Authorization: Bearer <token>` 读取用户信息
- 当前角色约定为 `role=0` 管理员，`role=1` 普通用户

## 重要业务说明

- Redis 初始化失败不会阻止服务启动，但会影响邮箱验证码、注册和找回密码等功能。
- `Logout` 当前没有做服务端 token 黑名单，更接近客户端删除 token。
- 文章详情接口 `GET /api/articles/:id` 当前是公开路由，但业务返回仍受文章状态与权限判断影响。
- 评论删除前需要先处理其子评论，否则会被拒绝。

## 开发建议

- 不要跳层调用：Controller 不直接查库，Repository 不承载业务规则。
- 新增功能优先补齐 `entity`、`dto`、`repository`、`service`、`api` 五层结构。
- 提交前至少执行 `make fmt`、`make vet`、`make test`。
- 若继续维护文档，建议统一保存为 UTF-8 编码，避免终端和编辑器出现乱码。

## 补充文档

- `docs/ARCHITECTURE.md`
- `docs/DEVELOPMENT.md`
- `docs/API.md`
