# CMDB Skeleton

该仓库提供 Kratos 后端与 Ant Design Pro Vue 前端的基础骨架。

## 后端

```bash
cd backend
# 安装依赖
# go mod tidy

# 生成 ent 代码（第一次执行或修改 schema 后）
# go run entgo.io/ent/cmd/ent generate ./internal/data/ent/schema

# 启动
go run ./cmd/cmdb
```

## 认证与 SSO

默认情况下，HTTP 服务会要求提供认证信息（`Authorization`、`X-User` 或 `sso_session` Cookie）。可通过访问 `/auth/sso/login` 触发 SSO 登录流程，登录完成后回调 `/auth/sso/callback` 并设置 `sso_session` Cookie。

可以使用环境变量 `SSO_REDIRECT_URL` 或 `configs/config.yaml` 中的 `sso.redirect_url` 配置跳转到真实的 SSO 入口地址。

## 前端

```bash
cd frontend
npm install
npm run dev
```

首次访问会进入登录页，可点击 “使用 SSO 登录” 跳转到后端 `/auth/sso/login`，或使用 “本地模拟登录” 进入控制台。
