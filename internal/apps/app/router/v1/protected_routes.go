package v1

import (
	"github.com/go-chi/chi/v5"
	authmodule "github.com/vadxq/go-rest-starter/internal/core/auth"
	usermodule "github.com/vadxq/go-rest-starter/internal/core/user"
	custommiddleware "github.com/vadxq/go-rest-starter/internal/transport/middleware"
	"github.com/vadxq/go-rest-starter/pkg/logger"
)

// SetupProtectedRoutes 设置受保护路由（需要认证）
func SetupProtectedRoutes(r chi.Router, config RouterConfig, jwtConfig *custommiddleware.JWTConfig, log logger.Logger) {
	// 创建需要JWT认证的路由组
	r.Group(func(r chi.Router) {
		r.Use(custommiddleware.JWTAuth(jwtConfig, log))

		// 认证相关路由（需要认证）
		authmodule.RegisterProtectedRoutes(r, config.AuthHandler)

		// 用户资源路由
		usermodule.RegisterRoutes(r, config.UserHandler)
	})
}
