# API 文档

## 基础信息

- Base URL：`http://localhost:8080`
- Content-Type：`application/json`
- 鉴权方式：`Authorization: Bearer <token>`

健康检查接口：

```http
GET /api/health
```

## 统一响应格式

项目中的接口统一返回类似结构：

```json
{
  "code": 0,
  "msg": "成功",
  "data": {}
}
```

说明：

- `code = 0` 表示业务成功
- `msg` 为响应消息
- `data` 为返回数据

注意：很多业务错误仍返回 HTTP `200`，客户端必须同时判断 `code`。

## 常见错误码

| code | 说明 |
| --- | --- |
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 1001 | 用户不存在 |
| 1002 | 用户已存在 |
| 1003 | 用户名或密码错误 |
| 1004 | 用户被禁用 |
| 1005 | 无效 token |
| 1006 | token 过期 |
| 1007 | 验证码无效 |

## 认证模块

路由前缀：`/api/auth`

### 注册

```http
POST /api/auth/register
```

请求体：

```json
{
  "username": "testuser",
  "password": "123456",
  "confirm_password": "123456",
  "email": "test@example.com",
  "captcha": "123456"
}
```

### 登录

```http
POST /api/auth/login
```

请求体：

```json
{
  "username": "testuser",
  "password": "123456",
  "captcha_id": "captcha-id",
  "captcha": "abcd"
}
```

### 刷新 token

```http
POST /api/auth/refresh
```

请求体：

```json
{
  "token": "your-token"
}
```

### 发送邮箱验证码

```http
POST /api/auth/email/code
```

请求体：

```json
{
  "email": "test@example.com"
}
```

### 重置密码

```http
POST /api/auth/password/reset
```

请求体：

```json
{
  "email": "test@example.com",
  "captcha": "123456",
  "new_password": "newpass123",
  "confirm_password": "newpass123"
}
```

### 获取图形验证码

```http
GET /api/auth/captcha
```

## 用户模块

路由前缀：`/api/user`

以下接口需要登录：

- `GET /profile`
- `PUT /profile`
- `POST /avatar`
- `POST /password`

### 获取个人资料

```http
GET /api/user/profile
```

### 更新个人资料

```http
PUT /api/user/profile
```

请求体示例：

```json
{
  "nickname": "新的昵称",
  "avatar": "https://example.com/avatar.jpg",
  "bio": "个人简介",
  "email": "test@example.com"
}
```

### 上传头像

```http
POST /api/user/avatar
```

该接口通常使用表单上传文件，返回上传后的资源地址。

### 修改密码

```http
POST /api/user/password
```

请求体：

```json
{
  "old_password": "123456",
  "new_password": "654321"
}
```

## 文章模块

路由前缀：`/api/articles`

### 公开接口

- `GET /api/articles`
- `GET /api/articles/:id`
- `GET /api/articles/:id/comments`
- `POST /api/articles/:id/view`

### 登录后接口

- `POST /api/articles`
- `POST /api/articles/cover_image`
- `GET /api/articles/favorites`
- `GET /api/articles/mine`
- `GET /api/articles/mine/:id`
- `PUT /api/articles/:id`
- `POST /api/articles/:id/publish`
- `DELETE /api/articles/:id`
- `POST /api/articles/:id/like`
- `POST /api/articles/:id/unlike`
- `POST /api/articles/:id/favorite`
- `POST /api/articles/:id/unfavorite`

### 列表查询参数

适用于 `GET /api/articles` 与 `GET /api/articles/mine`：

- `page`：页码，从 1 开始
- `pageSize`：每页条数
- `category_id`：可选分类
- `sort`：`latest` 或 `hottest`

### 创建文章

```http
POST /api/articles
```

请求体：

```json
{
  "title": "文章标题",
  "content": "文章内容",
  "summary": "摘要",
  "cover_image": "/uploads/example.jpg",
  "category_ids": [1, 2]
}
```

### 更新文章

```http
PUT /api/articles/:id
```

请求体与创建文章一致。

## 评论模块

路由前缀：`/api/comments`

以下接口都需要登录：

- `POST /api/comments`
- `POST /api/comments/replies`
- `POST /api/comments/:id/like`
- `POST /api/comments/:id/unlike`
- `DELETE /api/comments/:id`

### 发表评论

```http
POST /api/comments
```

请求体：

```json
{
  "article_id": 1,
  "content": "这是一条评论"
}
```

### 回复评论

```http
POST /api/comments/replies
```

请求体：

```json
{
  "article_id": 1,
  "content": "这是一条回复",
  "parent_id": 10,
  "root_id": 10
}
```

## 分类模块

```http
GET /api/categories
```

当前为公开只读接口。

## 管理审核模块

路由前缀：`/api/super`

这些接口需要登录，且业务层会校验管理员权限：

- `GET /api/super/articles`
- `GET /api/super/articles/:id`
- `GET /api/super/userlist`
- `POST /api/super/articles/:id/approve`
- `POST /api/super/articles/:id/reject`
- `POST /api/super/articles/:id/ban`
- `PUT /api/super/articles/:id/category`

### 审核列表参数

- `page`
- `pageSize`
- `category_id`
- `username`
- `status`

### 驳回请求体

```json
{
  "reason": "内容不符合规范"
}
```

### 更新文章分类请求体

```json
{
  "category_ids": [1, 2]
}
```

## 使用示例

### cURL

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"testuser\",\"password\":\"123456\",\"captcha_id\":\"id\",\"captcha\":\"abcd\"}"
```

```bash
curl http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## 注意事项

- Token 是否过期需要结合业务错误码判断
- 注册与重置密码依赖 Redis 中的邮箱验证码
- `POST /api/user/avatar` 与 `POST /api/articles/cover_image` 涉及上传目录权限
- 文章详情接口当前是公开路由，但业务返回仍受文章状态影响
