# 架构设计

## 目标与原则
- **领域优先**：业务能力集中在 `internal/core`，保持高内聚。
- **应用装配隔离**：每个应用独立组合根（`internal/apps/<app>`），支持多应用共享模块。
- **基础设施可替换**：数据库/缓存/队列等集中在 `internal/platform`。
- **传输层统一**：统一请求/响应工具和中间件（`internal/transport`）。

## 目录结构概览
```text
internal/
  apps/
    app/
      bootstrap/      # 应用装配：DI、启动、生命周期
      router/         # 应用路由组装
  core/
    user/             # 领域模块示例
      handler/
      service/
      repository/
      dto/
      model/
      routes.go
    auth/
    health/
  platform/
    config/           # 配置加载与校验
    db/               # 数据库与 Redis 连接
  transport/
    httpx/            # 请求/响应封装
    middleware/       # 中间件
```

## 分层职责
- **bootstrap（组合根）**：构建依赖、启动服务、生命周期管理。
- **router（路由装配）**：组装模块路由、中间件链与版本路由。
- **core（领域模块）**：业务逻辑与模型，严格面向接口。
- **platform（基础设施）**：外部依赖与配置，供应用装配层使用。
- **transport（传输层）**：HTTP 请求/响应统一规范与跨切面能力。

## 请求链路
```
cmd/app/main.go
  -> internal/apps/app/bootstrap/app.go
     -> internal/apps/app/router
        -> internal/transport/middleware
           -> internal/core/<domain>/handler
              -> internal/core/<domain>/service
                 -> internal/core/<domain>/repository
                    -> internal/platform/db
```

## 依赖方向
- `core` 不依赖 `apps` 与 `platform`。
- `apps` 依赖 `core`、`platform`、`transport`。
- `transport` 可依赖 `pkg`，但不依赖 `core`。

## 日志与追踪
- 统一使用 `pkg/logger`，并在响应中返回 `trace_id`。
- 请求日志与错误日志统一格式化，避免多套日志体系。

## 响应规范
统一响应结构：
```json
{ "code": 200, "message": "OK", "data": {...}, "trace_id": "..." }
```

错误响应 `data` 字段包含错误类型与字段信息。

## Redis 与缓存策略
- Redis 通过 `app.redis.enabled` 开关控制。
- Redis 不可用时自动降级为 Noop（缓存与队列静默失效）。

## 多应用扩展
新增应用时推荐结构：
```text
cmd/app2/main.go
internal/apps/app2/bootstrap
internal/apps/app2/router
```
领域模块仍复用 `internal/core`。
