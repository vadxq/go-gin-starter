# 部署指南

## 本地构建
```bash
./scripts/build.sh
```

## Docker 镜像
```bash
./scripts/docker-build.sh
```

## Docker Compose
```bash
docker compose -f deploy/docker/docker-compose.yaml up -d
```

## Kubernetes
参考 `deploy/k8s/deployment.yaml`，需要根据环境配置：
- 数据库连接信息
- Redis 连接信息
- JWT 密钥

## 配置建议
- 默认读取 `configs/config.yaml`，可通过 `CONFIG_PATH` 指定
- 生产环境优先通过环境变量覆盖配置（前缀 `APP_`，如 `APP_DB_HOST`）
- Redis 可选，禁用时使用 `APP_REDIS_ENABLED=false`
