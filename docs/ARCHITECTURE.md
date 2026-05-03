# 架构设计

## 概览

`blog` 是一个基于 Go 的博客后端，采用清晰的三层结构：

`Controller -> Service -> Repository`

项目入口位于 `cmd/server/main.go`，真正的启动与装配逻辑集中在 `internal/app/app.go`。

## 分层职责

### API 层

目录：`internal/api/`

职责：

- 处理 HTTP 请求与参数绑定
- 调用 service 执行业务逻辑
- 返回统一响应结构

这一层不直接访问数据库，也不承载复杂业务规则。

### Service 层

目录：`internal/service/`

职责：

- 编排核心业务逻辑
- 做权限判断与状态流转
- 组合多个 repository 的操作
- 构建返回给控制器的结果

例如文章发布、评论回复、管理员审核，都应在 service 层完成。

### Repository 层

目录：`internal/repository/`

职责：

- 封装 GORM 查询
- 执行事务和关联更新
- 处理 MySQL 持久化细节

这一层只关心数据访问，不做业务语义判断。

## 核心目录

- `internal/app`：应用初始化、依赖注入、HTTP 服务启动与优雅关闭
- `internal/model/entity`：数据库实体
- `internal/model/dto/request`：请求 DTO
- `internal/model/dto/response`：响应 DTO
- `internal/middleware`：日志、恢复、CORS、JWT 鉴权
- `pkg/config`：配置加载
- `pkg/database`：MySQL / Redis 初始化与关闭
- `pkg/logger`：日志封装
- `pkg/jwt`：JWT 工具
- `pkg/response`：统一响应格式
- `pkg/errors`：业务错误码与错误对象

## 启动流程

`internal/app/app.go` 的初始化顺序为：

1. `initConfig()`
2. `initLogger()`
3. `initDatabase()`
4. `initDependencies()`
5. `initRouter()`
6. `initServer()`

其中：

- `initDatabase()` 会初始化 MySQL，并执行 `AutoMigrate(...)`
- Redis 初始化失败不会阻止服务启动，但会影响验证码、注册、找回密码等功能

## 路由结构

路由注册集中在 `internal/api/router.go`：

- 健康检查：`GET /api/health`
- 认证模块：`/api/auth`
- 用户模块：`/api/user`
- 文章模块：`/api/articles`
- 评论模块：`/api/comments`
- 分类模块：`/api/categories`
- 管理审核模块：`/api/super`

## 权限模型

鉴权使用 JWT，Token 从 `Authorization: Bearer <token>` 读取。

当前角色约定：

- `role = 0`：管理员
- `role = 1`：普通用户

需要注意的是，部分接口的“是否管理员”判断不在路由层，而在 service 层完成。

## 数据与状态约定

文章状态在现有实现中包括：

- `0`：草稿
- `1`：待审核
- `2`：已发布
- `3`：已拒绝
- `4`：已封禁

统一响应格式定义在 `pkg/response/response.go`，业务错误通常仍返回 HTTP `200`，调用方必须结合 JSON 中的 `code` 和 `msg` 判断结果。

## 扩展建议

新增一个完整业务模块时，建议按以下顺序扩展：

1. 在 `entity` 中新增实体
2. 在 `dto/request` 与 `dto/response` 中新增结构
3. 在 `repository` 中新增接口与实现
4. 在 `service` 中新增接口与实现
5. 在 `api/<module>` 中新增 controller 与 routes
6. 在 `app.go` 的 `initDependencies()` 中注入依赖
7. 在 `router.go` 中注册路由
8. 必要时加入 `AutoMigrate(...)`

最容易漏改的地方通常是 `initDependencies()`、`router.go` 和数据库迁移注册。
