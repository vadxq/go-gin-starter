# 路由系统

```text
┌──────────────────────────────┐
│ 全局中间件                    │
│ RequestID/RealIP/Timeout/... │
└───────────────┬──────────────┘
                │
    ┌───────────┼───────────┐
    │ Swagger   │ 健康检查   │
    │ /swagger  │ /health…   │
    └───────────┴───────────┘
                │
             /api/v1
        ┌────────┴────────┐
        │ 公共路由         │
        │ /auth/login      │
        │ /auth/refresh    │
        └────────┬─────────┘
                 │
           受保护路由 (JWT)
           /account/logout
           /users
```

## 路由说明
- Swagger：`/swagger` 与 `/swagger/*` 提供文档界面与 JSON。
- 健康检查：`/health`、`/health/detailed`、`/ready`、`/live`、`/version`、`/status`。
- API v1：`/api/v1` 下分为公共与受保护路由，受保护路由需 JWT。
