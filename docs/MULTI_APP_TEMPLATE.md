# 多应用模板

本模板用于在同一仓库内新增第二个应用（例如 `app2`），复用 `internal/core` 领域模块与 `internal/platform` 基础设施。

## 目录结构模板
```text
cmd/
  app2/
    main.go
internal/
  apps/
    app2/
      bootstrap/
        app.go
        injection/
          dependencies.go
          repositories.go
          services.go
          handlers.go
      router/
        router.go
        v1/
          public_routes.go
          protected_routes.go
          swagger.go
```

## 新应用创建步骤
1. 复制现有应用装配目录作为模板：
   - `internal/apps/app` → `internal/apps/app2`
2. 新增入口文件：`cmd/app2/main.go`
3. 需要时调整配置路径（建议使用独立环境变量）。
4. 按需裁剪路由与模块注入（只保留该应用需要的模块）。

## 入口模板（`cmd/app2/main.go`）
```go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	app2bootstrap "github.com/vadxq/go-rest-starter/internal/apps/app2/bootstrap"
	"github.com/vadxq/go-rest-starter/pkg/logger"
)

func main() {
	var appLogger logger.Logger = logger.Default()

	application, err := app2bootstrap.New()
	if err != nil {
		appLogger.Error("创建应用失败", "error", err)
		os.Exit(1)
	}
	appLogger = application.Logger()

	serverErrCh := application.StartServer()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case err := <-serverErrCh:
		appLogger.Error("服务器错误", "error", err)
	case sig := <-signalCh:
		appLogger.Info("接收到系统信号，开始优雅关闭", "signal", sig.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		appLogger.Error("应用关闭失败", "error", err)
		os.Exit(1)
	}
}
```

## 装配模板（`internal/apps/app2/bootstrap/app.go`）
> 可以直接复制 `internal/apps/app/bootstrap/app.go`，并按需调整配置路径：
- 默认配置：读取 `CONFIG_PATH`（与现有应用共用）。
- 独立配置：使用 `APP2_CONFIG_PATH`（建议）。

建议示例：
```go
func getConfigPath() string {
	configPath := os.Getenv("APP2_CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	return configPath
}
```

## 路由模板（`internal/apps/app2/router/router.go`）
- 复用 `internal/core` 模块
- 仅装配需要的模块

示例要点：
```go
api.Setup(router, api.RouterConfig{
	UserHandler:   deps.Handlers.UserHandler,
	AuthHandler:   deps.Handlers.AuthHandler,
	HealthHandler: deps.Handlers.HealthHandler,
	JWTSecret:     deps.Config.JWT.Secret,
	Logger:        appLogger,
})
```

## 裁剪模块与依赖
- 不需要的模块可以从注入与路由中移除：
  - `internal/apps/app2/bootstrap/injection/*`
  - `internal/apps/app2/router/*`
- 领域模块始终位于 `internal/core`，可被多个应用复用。

## 推荐约定
- 应用入口与装配尽量保持统一结构，便于团队协作与扩展。
- 不同应用可使用不同端口与配置文件，避免冲突。
