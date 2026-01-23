# 故障排查

## 配置未生效
- 检查 `configs/config.yaml` 是否存在
- 确认环境变量前缀为 `APP_`

## 数据库连接失败
- 检查 `app.database` 配置
- 确认数据库服务可达

## Redis 连接失败
- 如不需要 Redis，可设置 `APP_REDIS_ENABLED=false`
- 启动 Redis：`docker compose -f deploy/docker/docker-compose.yaml up -d`

## Swagger 不一致
- 运行 `./scripts/swagger.sh` 重新生成文档

## 端口占用
- 修改 `app.server.port` 或释放占用端口
