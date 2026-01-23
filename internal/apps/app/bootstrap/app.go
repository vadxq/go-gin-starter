package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/vadxq/go-rest-starter/internal/apps/app/bootstrap/injection"
	api "github.com/vadxq/go-rest-starter/internal/apps/app/router"
	"github.com/vadxq/go-rest-starter/internal/platform/config"
	"github.com/vadxq/go-rest-starter/internal/platform/db"
	httpx "github.com/vadxq/go-rest-starter/internal/transport/httpx"
	"github.com/vadxq/go-rest-starter/pkg/cache"
	"github.com/vadxq/go-rest-starter/pkg/logger"
)

// App 应用结构体
type App struct {
	DB        *gorm.DB
	Redis     *redis.Client
	Router    *chi.Mux
	Cache     cache.Cache
	Validator *validator.Validate
	Deps      *injection.Dependencies
	Server    *http.Server
	Config    *config.AppConfig
	logger    logger.Logger
}

// New 创建新的应用实例
func New() (*App, error) {
	// 配置日志输出
	configPath := getConfigPath()

	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	// 创建日志器
	appLogger, err := logger.NewLogger(&logger.LogConfig{
		Level:   cfg.Log.Level,
		File:    cfg.Log.File,
		Console: cfg.Log.Console,
	})
	if err != nil {
		return nil, fmt.Errorf("创建日志器失败: %w", err)
	}

	appLogger.Info("配置加载完成", "config_path", configPath)
	httpx.SetLogger(appLogger)

	// 创建应用实例
	app := &App{
		Config: cfg,
		logger: appLogger,
	}

	// 初始化应用
	if err := app.initialize(); err != nil {
		return nil, fmt.Errorf("初始化应用失败: %w", err)
	}

	return app, nil
}

// initialize 初始化应用组件
func (app *App) initialize() error {
	app.logger.Info("开始初始化应用...")

	// 初始化数据库连接
	if err := app.initDatabase(); err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}

	// 初始化Redis连接
	if err := app.initRedis(); err != nil {
		return fmt.Errorf("初始化Redis失败: %w", err)
	}

	// 初始化缓存
	if err := app.initCache(); err != nil {
		return fmt.Errorf("初始化缓存失败: %w", err)
	}

	// 初始化验证器
	app.Validator = validator.New()

	// 初始化依赖注入
	if err := app.initDependencies(); err != nil {
		return fmt.Errorf("初始化依赖注入失败: %w", err)
	}

	// 初始化路由
	if err := app.initRouter(); err != nil {
		return fmt.Errorf("初始化路由失败: %w", err)
	}

	app.logger.Info("应用初始化完成")
	return nil
}

// initDatabase 初始化数据库连接
func (app *App) initDatabase() error {
	app.logger.Info("连接数据库...")

	database, err := db.InitDB(&app.Config.Database)
	if err != nil {
		return err
	}

	app.DB = database
	app.logger.Info("数据库连接成功")
	return nil
}

// initRedis 初始化Redis连接
func (app *App) initRedis() error {
	app.logger.Info("连接Redis...")

	if !app.Config.Redis.Enabled {
		app.logger.Info("Redis已禁用，跳过初始化")
		return nil
	}

	redisClient, err := db.InitRedis(&app.Config.Redis)
	if err != nil {
		app.logger.Warn("Redis连接失败，降级为Noop", "error", err)
		return nil
	}

	if redisClient == nil {
		app.logger.Info("Redis未启用")
		return nil
	}

	app.Redis = redisClient
	app.logger.Info("Redis连接成功")
	return nil
}

// initCache 初始化缓存
func (app *App) initCache() error {
	app.logger.Info("初始化缓存...")

	// 缓存服务必须依赖Redis
	if app.Redis == nil {
		app.logger.Warn("Redis不可用，缓存降级为Noop")
		app.Cache = cache.NewNoop()
		return nil
	}

	cacheOpts := cache.Options{
		DefaultExpiration: 10 * time.Minute,
		CleanupInterval:   5 * time.Minute,
		RedisAddress:      fmt.Sprintf("%s:%d", app.Config.Redis.Host, app.Config.Redis.Port),
		RedisPassword:     app.Config.Redis.Password,
		RedisDB:           app.Config.Redis.DB,
	}

	app.logger.Info("使用Redis作为缓存存储")

	cacheInstance, err := cache.NewCache(cacheOpts)
	if err != nil {
		app.logger.Warn("初始化Redis缓存失败，缓存降级为Noop", "error", err)
		app.Cache = cache.NewNoop()
		return nil
	}

	app.Cache = cacheInstance
	app.logger.Info("缓存初始化成功")
	return nil
}

// initDependencies 初始化依赖注入
func (app *App) initDependencies() error {
	app.logger.Info("初始化依赖注入系统...")

	deps := injection.NewDependencies(
		app.DB,
		app.Redis,
		app.Validator,
		app.Config,
		app.Cache,
		app.logger,
	)

	app.Deps = deps
	app.logger.Info("依赖注入系统初始化完成")
	return nil
}

// initRouter 初始化路由
func (app *App) initRouter() error {
	app.logger.Info("配置API路由...")

	router := chi.NewRouter()

	api.Setup(router, api.RouterConfig{
		UserHandler:   app.Deps.Handlers.UserHandler,
		AuthHandler:   app.Deps.Handlers.AuthHandler,
		HealthHandler: app.Deps.Handlers.HealthHandler,
		JWTSecret:     app.Deps.Config.JWT.Secret,
		Logger:        app.logger,
	})

	app.Router = router
	app.logger.Info("API路由配置完成")
	return nil
}

// StartServer 启动HTTP服务器
func (app *App) StartServer() <-chan error {
	errCh := make(chan error, 1)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Server.Port),
		Handler:      app.Router,
		ReadTimeout:  app.Config.Server.ReadTimeout,
		WriteTimeout: app.Config.Server.WriteTimeout,
	}

	app.Server = server

	// 启动服务器
	go func() {
		app.logger.Info("HTTP服务器启动", "port", app.Config.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("HTTP服务器错误: %w", err)
		}
	}()

	return errCh
}

// Shutdown 优雅关闭应用
func (app *App) Shutdown(ctx context.Context) error {
	app.logger.Info("开始优雅关闭应用...")

	// 使用channel收集错误
	errChan := make(chan error, 3)

	// 并发关闭各个组件
	go func() {
		if app.Server != nil {
			app.logger.Info("关闭HTTP服务器...")
			errChan <- app.Server.Shutdown(ctx)
		} else {
			errChan <- nil
		}
	}()

	go func() {
		if app.DB != nil {
			app.logger.Info("关闭数据库连接...")
			if sqlDB, err := app.DB.DB(); err == nil {
				errChan <- sqlDB.Close()
			} else {
				errChan <- err
			}
		} else {
			errChan <- nil
		}
	}()

	go func() {
		if app.Redis != nil {
			app.logger.Info("关闭Redis连接...")
			errChan <- app.Redis.Close()
		} else {
			errChan <- nil
		}
	}()

	// 等待所有关闭操作完成
	var hasError bool
	for i := 0; i < 3; i++ {
		if err := <-errChan; err != nil {
			app.logger.Error("关闭组件失败", "error", err)
			hasError = true
		}
	}

	if hasError {
		app.logger.Warn("应用关闭时出现错误")
	} else {
		app.logger.Info("应用优雅关闭完成")
	}

	return nil
}

// Logger returns the application logger.
func (app *App) Logger() logger.Logger {
	return app.logger
}

// 获取配置文件路径
func getConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	return configPath
}
