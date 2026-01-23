# 开发指南

## 环境准备
- Go 1.24+
- PostgreSQL 12+
- Redis 6+（可选）

初始化配置：
```bash
cp configs/config.example.yaml configs/config.yaml
```

可通过环境变量覆盖配置（前缀 `APP_`）。

## 常用命令
- 开发运行：`./scripts/dev.sh`
- 直接运行：`go run cmd/app/main.go`
- 生成 Swagger：`./scripts/swagger.sh`
- 运行测试：`go test ./...`

## 模块开发流程
以新增 `order` 模块为例：
1. 创建目录：`internal/core/order/{handler,service,repository,dto,model}`
2. 实现 handler/service/repository
3. 在 `internal/core/order/routes.go` 注册路由
4. 更新注入：`internal/apps/app/bootstrap/injection/*`
5. 更新路由装配：`internal/apps/app/router`（若需接入 v1 路由）

## 新增 API 端点
- 在模块 `handler` 中新增方法
- 在模块 `routes.go` 中注册
- 补齐 Swagger 注解并运行 `./scripts/swagger.sh`

## 多应用扩展
- 新增 `cmd/<app>/main.go`
- 新增 `internal/apps/<app>/bootstrap` 与 `router`
- 复用 `internal/core` 模块

## 编码规范
- `gofmt`（tab 缩进）
- 文件名 snake_case
- 模块内分层清晰：handler/service/repository
