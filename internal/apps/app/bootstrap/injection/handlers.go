package injection

import (
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	authhandler "github.com/vadxq/go-rest-starter/internal/core/auth/handler"
	healthhandler "github.com/vadxq/go-rest-starter/internal/core/health/handler"
	userhandler "github.com/vadxq/go-rest-starter/internal/core/user/handler"
	"github.com/vadxq/go-rest-starter/pkg/logger"
)

// Handlers 包含所有HTTP处理器
type Handlers struct {
	UserHandler   *userhandler.UserHandler
	AuthHandler   *authhandler.AuthHandler
	HealthHandler *healthhandler.HealthHandler
}

// InitHandlers 初始化所有HTTP处理器
func InitHandlers(
	services *Services,
	log logger.Logger,
	validator *validator.Validate,
	db *gorm.DB,
	redis *redis.Client,
) *Handlers {
	// 初始化用户处理器
	userHandler := userhandler.NewUserHandler(
		services.UserService,
		log,
		validator,
	)

	// 初始化认证处理器
	authHandler := authhandler.NewAuthHandler(
		services.AuthService,
		log,
		validator,
	)

	// 初始化健康检查处理器
	healthHandler := healthhandler.NewHealthHandler(
		db,
		redis,
		log,
	)

	return &Handlers{
		UserHandler:   userHandler,
		AuthHandler:   authHandler,
		HealthHandler: healthHandler,
	}
}
