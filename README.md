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

默认情况下，HTTP 服务会要求提供认证信息（`Authorization`、`X-User`、`local_session` 或 `sso_session` Cookie）。本地认证通过 `/auth/login` 提交用户名与密码（示例用途），成功后会写入 `local_session` Cookie。

SSO 默认关闭，需在后端配置后启用。使用环境变量 `SSO_ENABLED=true` 开启，并配置 `SSO_SAML2_LOGIN_URL` 指向 IdP 的 SAML2 登录地址。配置示例见 `configs/config.yaml` 的 `sso` 段落。

## 前端

```bash
cd frontend
npm install
npm run dev
```

首次访问会进入登录页，可点击 “使用 SSO 登录” 跳转到后端 `/auth/sso/login`，或使用 “本地模拟登录” 进入控制台。
