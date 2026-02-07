# API 生成说明

接口定义位于 `backend/api/cmdb/v1/cmdb.proto`，建议使用 `buf` 或 `protoc` 生成代码。

示例（使用 buf）：

```bash
buf generate
```

示例（使用 protoc）：

```bash
protoc --go_out=. --go-grpc_out=. --go-http_out=. api/cmdb/v1/cmdb.proto
```
