package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	v1 "github.com/vadxq/go-rest-starter/internal/apps/app/router/v1"
	authhandler "github.com/vadxq/go-rest-starter/internal/core/auth/handler"
	healthmodule "github.com/vadxq/go-rest-starter/internal/core/health"
	healthhandler "github.com/vadxq/go-rest-starter/internal/core/health/handler"
	userhandler "github.com/vadxq/go-rest-starter/internal/core/user/handler"
	custommiddleware "github.com/vadxq/go-rest-starter/internal/transport/middleware"
	"github.com/vadxq/go-rest-starter/pkg/logger"
)

// 路由组类型定义
type RouterGroup struct {
	Router     chi.Router
	Middleware []func(http.Handler) http.Handler
}

// RouterConfig 路由配置
type RouterConfig struct {
	UserHandler   *userhandler.UserHandler
	AuthHandler   *authhandler.AuthHandler
	HealthHandler *healthhandler.HealthHandler
	JWTSecret     string
	Logger        logger.Logger
}

// Setup 设置所有API路由
func Setup(r chi.Router, config RouterConfig) {
	// 应用全局中间件
	applyGlobalMiddleware(r, config.Logger)

	// API文档路由
	v1.SetupSwaggerRoutes(r)

	// 健康检查和状态监控
	healthmodule.RegisterRoutes(r, config.HealthHandler)

	// API v1
	setupV1Routes(r, config)
}

// applyGlobalMiddleware 应用全局中间件
func applyGlobalMiddleware(r chi.Router, log logger.Logger) {
	// 基础中间件
	r.Use(middleware.RequestID)                    // 请求ID
	r.Use(middleware.RealIP)                       // 真实IP
	r.Use(custommiddleware.RequestContext)         // 请求上下文
	r.Use(custommiddleware.LoggingMiddleware(log)) // 日志
	r.Use(custommiddleware.RecoveryMiddleware)     // 恢复
	r.Use(middleware.Timeout(60 * time.Second))    // 超时
	r.Use(middleware.CleanPath)                    // 清理路径
	r.Use(middleware.StripSlashes)                 // 去除尾部斜杠

	// 安全中间件
	r.Use(custommiddleware.CORSMiddleware) // 跨域
	r.Use(securityHeaders)                 // 安全头

	// 速率限制中间件
	rateLimiter := custommiddleware.NewRateLimitMiddleware(custommiddleware.DefaultRateLimitConfig)
	r.Use(rateLimiter.Handler) // 速率限制
}

// setupV1Routes 设置 API v1 路由
func setupV1Routes(r chi.Router, config RouterConfig) {
	// 定义排除认证的路径
	excludePaths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/refresh",
		"/swagger",
		"/health",
		"/version",
		"/status",
	}

	// 创建JWT认证配置
	jwtConfig := &custommiddleware.JWTConfig{
		Secret:       config.JWTSecret,
		ExcludePaths: excludePaths,
	}

	// API v1 基础路径
	r.Route("/api/v1", func(r chi.Router) {
		v1Config := v1.RouterConfig{
			UserHandler: config.UserHandler,
			AuthHandler: config.AuthHandler,
			JWTSecret:   config.JWTSecret,
		}
		// 公共路由组 - 不需要认证
		v1.SetupPublicRoutes(r, v1Config)
		// 受保护路由组 - 需要认证
		v1.SetupProtectedRoutes(r, v1Config, jwtConfig, config.Logger)
	})
}

// securityHeaders 添加安全相关的HTTP头
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 防止MIME类型嗅探
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// 开启XSS过滤
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// 防止点击劫持
		w.Header().Set("X-Frame-Options", "DENY")
		// HTTP严格传输安全
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		// 引用策略
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// 内容安全策略
		cspValue := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data:; " +
			"style-src 'self' 'unsafe-inline'; " +
			"font-src 'self' data:;"
		w.Header().Set("Content-Security-Policy", cspValue)

		next.ServeHTTP(w, r)
	})
}
