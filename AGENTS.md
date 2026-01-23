# 仓库指南

## 项目结构与模块组织
服务入口为 `cmd/app/main.go`。应用装配位于 `internal/apps/app/bootstrap/`，路由装配位于 `internal/apps/app/router/`。领域模块位于 `internal/core/`（如 `auth`、`user`、`health`），每个模块包含 `handler/service/repository/dto/model/routes` 等子结构。基础设施位于 `internal/platform/`（如 `config`、`db`），传输层工具位于 `internal/transport/`（如 `httpx`、`middleware`）。可复用组件在 `pkg/`（cache、jwt、logger、queue、transaction、utils）。生成的 API 文档在 `api/app/`，配置在 `configs/`，SQL 迁移在 `migrations/`，部署资源在 `deploy/`，脚本在 `scripts/`。

## 构建、测试与开发命令
- `./scripts/dev.sh` 生成 Swagger 并使用 `air` 热重载（`.air.toml`）。
- `go run cmd/app/main.go` 直接运行服务。
- `./scripts/swagger.sh` 重新生成 Swagger 到 `api/app/`。
- `./scripts/build.sh` 构建多平台二进制到 `build/`。
- `go test ./...` 运行全部测试；可加 `-race` 或 `-coverprofile=coverage.out`。
- `docker compose -f deploy/docker/docker-compose.yaml up -d` 启动 PostgreSQL 与 Redis。

## 编码风格与命名
使用 `gofmt`（tab 缩进）。包名小写，文件名使用 snake_case（如 `user_service.go`、`auth_handler.go`），导出标识符 PascalCase，非导出 camelCase。新增类型优先放入对应领域模块下，并保持 handler/service/repository 分层清晰。

## 测试规范
测试与代码同目录，文件名 `*_test.go`，函数名 `TestXxx`。项目使用 `testify`（参考 `internal/core/user/service/user_service_test.go`）。修改服务逻辑时可运行：`go test ./internal/core/user/service/`。

## 提交与 PR 规范
提交信息遵循历史前缀风格（如 `feat: ...`、`fix: ...`）。PR 需包含简要说明、测试命令，以及是否更新了配置/迁移/Swagger（如 `api/app/` 变更请注明）。

## 配置与安全
本地开发请复制 `configs/config.example.yaml` 为 `configs/config.yaml`，并用 `APP_` 环境变量覆盖。不要提交密钥；生产使用 `configs/config.production.yaml` 或环境变量覆盖。若修改路由或请求/响应结构，请运行 `./scripts/swagger.sh` 更新文档。

## 回复

- 使用中文交流和回复
- 文档使用中文 Markdown 格式
