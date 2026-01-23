# 缓存组件（pkg/cache）

## 设计决策：仅 Redis

该包仅支持 Redis 作为缓存存储，不提供内存缓存实现，主要考虑如下：

1. **避免内存不可控增长**
   - 进程内缓存容易出现内存膨胀和泄漏风险
   - 容量无法与容器/节点资源配额严格对齐

2. **分布式一致性**
   - 多实例部署下，本地缓存会造成数据不一致
   - Redis 提供集中式缓存，易于一致性管理

3. **监控与运维友好**
   - Redis 有成熟的可视化与告警体系
   - 便于设置内存上限与淘汰策略

4. **可恢复性**
   - Redis 支持持久化与重启恢复
   - 无需额外的缓存预热流程

## 使用方式

```go
// 初始化缓存（依赖 Redis）
cacheOpts := cache.Options{
    RedisAddress:      "localhost:6379",
    RedisPassword:     "password",
    RedisDB:           0,
    DefaultExpiration: 10 * time.Minute,
    CleanupInterval:   5 * time.Minute,
}

cacheInstance, err := cache.NewCache(cacheOpts)
if err != nil {
    // Redis 不可用时可降级为 Noop
    cacheInstance = cache.NewNoop()
}

// Set / Get
_ = cacheInstance.Set(ctx, "key", []byte("value"), 5*time.Minute)
value, err := cacheInstance.Get(ctx, "key")

// SetObject / GetObject
user := &User{ID: 1, Name: "John"}
_ = cacheInstance.SetObject(ctx, "user:1", user, 10*time.Minute)

var cachedUser User
err = cacheInstance.GetObject(ctx, "user:1", &cachedUser)
```

## 缓存策略

包内提供以下策略（见 `pkg/cache/strategies.go`）：

- **Cache-Aside**：缓存未命中时回源加载并写入缓存
- **SingleFlight**：同一 key 并发请求只触发一次回源

## 配置

Redis 配置来自 `app.redis`：

```yaml
app:
  redis:
    enabled: true
    host: localhost
    port: 6379
    password: ""
    db: 0
```

环境变量覆盖使用 `APP_` 前缀（如 `APP_REDIS_HOST`、`APP_REDIS_ENABLED`）。

## 降级行为

- `cache.NewCache` 需要 Redis 地址，否则返回错误
- 应用启动时 Redis 禁用或连接失败，会使用 `cache.NewNoop()`
- `Noop` 实现的 `Get/GetObject` 返回 `cache.ErrNotFound`，写操作为无副作用
