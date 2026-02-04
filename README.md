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

## 前端

```bash
cd frontend
npm install
npm run dev
```
