package v1

import (
	"github.com/go-chi/chi/v5"
	authmodule "github.com/vadxq/go-rest-starter/internal/core/auth"
	authhandler "github.com/vadxq/go-rest-starter/internal/core/auth/handler"
	userhandler "github.com/vadxq/go-rest-starter/internal/core/user/handler"
)

// RouterConfig 路由配置
type RouterConfig struct {
	UserHandler *userhandler.UserHandler
	AuthHandler *authhandler.AuthHandler
	JWTSecret   string
}

// SetupPublicRoutes 设置公共路由（不需要认证）
func SetupPublicRoutes(r chi.Router, config RouterConfig) {
	authmodule.RegisterPublicRoutes(r, config.AuthHandler)
}
