# vue-blog

基于 `Vue 3 + Vite` 的博客前端项目，当前已经按常规 Vue 3 生产级目录结构完成整理，并接入现有后端接口。

## 技术栈

- Vue 3
- Vite
- Vue Router
- Pinia
- Axios

## 当前功能

- 首页公开文章列表
- 登录、注册、找回密码弹层
- 登录成功后右上角显示默认头像
- 个人中心页面
- 个人资料维护
- 我的文章管理
- 我的收藏查看
- 账户安全密码修改

## 目录结构

```text
public/
├── favicon.ico               # 不经过构建处理的静态资源

src/
├── api/                      # 后端接口封装
├── assets/                   # 全局样式与静态资源
├── components/
│   ├── auth/                 # 登录/注册/找回密码相关组件
│   ├── common/               # Toast 等通用组件
│   └── layout/               # Header 等布局组件
├── composables/              # 预留给可复用 Composition API 逻辑
├── layouts/                  # 页面整体布局
├── router/                   # Vue Router 配置
├── stores/                   # Pinia 状态管理
├── utils/                    # 工具函数与格式化逻辑
├── views/                    # 路由页面组件
├── App.vue                   # 根组件
└── main.js                   # 应用入口

index.html                    # HTML 模板
vite.config.js                # Vite 配置
package.json                  # 依赖与脚本
README.md                     # 项目说明
```

## 命名规范

- 页面组件：`PascalCaseView.vue`
- 通用组件：`PascalCase.vue`
- 布局组件：`PascalCase.vue`
- Store 文件：`xxxStore.js`
- 工具模块：`camelCase.js`
- 路由入口：`src/router/index.js`
- 应用入口：`src/main.js`

## 运行方式

安装依赖：

```sh
npm install
```

本地开发：

```sh
npm run dev
```

生产构建：

```sh
npm run build
```

## 开发代理

开发环境通过 `vite.config.js` 代理后端接口：

- `/api` -> `http://127.0.0.1:8082`
- `/uploads` -> `http://127.0.0.1:8082`

如果后端地址变更，可以修改 `vite.config.js`，或通过 `VITE_API_BASE_URL` 指定接口基础地址。

## 路由说明

- `/`：首页
- `/profile`：个人中心

个人中心通过查询参数切换分区：

- `/profile?tab=profile`
- `/profile?tab=articles`
- `/profile?tab=favorites`
- `/profile?tab=security`

## 说明

项目已经移除旧的单文件堆叠式实现，当前运行路径统一为：

- `router`
- `layouts`
- `views`
- `stores`
- `components`

后续新增功能建议继续按该结构扩展，不再把业务逻辑集中写入 `App.vue`。
