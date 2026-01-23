# 部署指南

## 本地容器
```bash
./scripts/build.sh
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
- 生产环境优先通过环境变量覆盖配置
- Redis 可选，禁用时使用 `APP_REDIS_ENABLED=false`
