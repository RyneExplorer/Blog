# 开发指南

## 开发前准备

运行环境要求：

- Go `1.25.0` 或与 `go.mod` 兼容的版本
- MySQL `8.0+`
- Redis `6.0+`

初始化项目：

```bash
go mod tidy
cp configs/config.yaml.example configs/config.yaml
make run
```

敏感配置建议通过环境变量注入：

- `MYSQL_PASSWORD`
- `REDIS_PASSWORD`
- `JWT_SECRET`

## 常用命令

```bash
make run            # 启动服务
make build          # 编译程序
make test           # 运行测试
make test-coverage  # 生成覆盖率报告
make fmt            # 格式化代码
make vet            # 静态检查
make deps           # 整理并下载依赖
make clean          # 清理构建产物
```

## 开发约定

项目遵循固定分层：

- Controller 只处理请求、参数绑定和响应
- Service 负责业务逻辑、权限判断、状态流转
- Repository 负责数据库访问

不要跳层调用：

- Controller 不直接操作数据库
- Service 不依赖 Gin 上下文
- Repository 不承载业务规则

## 新增功能的推荐流程

以新增一个业务模块为例，建议按顺序完成：

1. 在 `internal/model/entity` 中新增实体
2. 在 `internal/model/dto/request` 中新增请求结构
3. 在 `internal/model/dto/response` 中新增响应结构
4. 在 `internal/repository` 中新增接口与实现
5. 在 `internal/service` 中新增接口与实现
6. 在 `internal/api/<module>` 中新增 controller 与路由
7. 在 `internal/app/app.go` 的 `initDependencies()` 中接入依赖
8. 在 `internal/api/router.go` 中注册路由
9. 在 `AutoMigrate(...)` 中加入新实体

这个项目使用手工依赖注入，因此第 7 步和第 8 步很容易漏掉。

## DTO 与参数校验

请求结构统一放在 `internal/model/dto/request/`，使用 Gin 的 `binding` 标签校验字段，例如：

```go
type LoginRequest struct {
    Username  string `json:"username" binding:"required"`
    Password  string `json:"password" binding:"required"`
    CaptchaID string `json:"captcha_id" binding:"required"`
    Captcha   string `json:"captcha" binding:"required,len=4"`
}
```

建议：

- 参数校验尽量放在 request DTO 中完成
- Controller 只做绑定与错误返回
- 更复杂的业务校验放到 service 中

## 错误处理

业务错误统一使用 `pkg/errors` 中的错误码和错误对象。调用方不能只依赖 HTTP 状态码，因为很多业务失败仍返回 HTTP `200`。

推荐做法：

- 参数错误走绑定校验
- 业务错误返回统一 `BizError`
- 日志中记录关键上下文，不暴露敏感信息

## 日志与调试

日志由 `pkg/logger` 统一管理，常见调试方式：

- 查看程序日志输出
- 打印结构化字段，例如 `zap.Any(...)`
- 在开发环境下观察 GORM SQL 输出

如果涉及上传功能，还需要确认 `uploads/` 目录可写，且静态资源映射正常。

## 测试建议

当前仓库中测试文件较少，新增功能时建议至少补充 service 或 repository 层测试。测试文件命名使用 `*_test.go`，例如：

- `user_service_test.go`
- `article_repository_test.go`

提交前至少执行：

```bash
make fmt
make vet
make test
```

## 协作注意事项

- 文档统一使用 UTF-8 编码
- 不要把敏感配置写入仓库
- Redis 虽然不是启动强依赖，但缺失会影响验证码、注册和重置密码流程
- `Logout` 当前没有服务端黑名单机制，如需更强安全性应补充 token 失效策略
