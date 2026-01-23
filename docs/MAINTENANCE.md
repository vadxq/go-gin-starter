# 维护指南

## 配置管理
- 本地使用 `configs/config.yaml`
- 生产使用 `configs/config.production.yaml` 或环境变量
- Redis 可通过 `app.redis.enabled` 禁用

## 数据库迁移
- 新增迁移文件：`migrations/app/*.sql`
- 修改模型或仓库后建议同步迁移

## 依赖升级
- 使用 `go get -u` 升级依赖
- 升级后运行 `go test ./...`
- 重点关注数据库驱动与缓存库的变更说明

## 日志与监控
- 日志统一使用 `pkg/logger`
- 健康检查：`/health`、`/health/detailed`、`/ready`、`/live`

## 安全与合规
- JWT 密钥必须通过环境变量提供
- 确认 CORS、安全头策略符合生产要求
- 定期更新依赖以修复漏洞

## 版本与发布
- 提交遵循 `feat:`/`fix:` 风格
- 发布前确认：测试通过、配置更新、Swagger 同步
