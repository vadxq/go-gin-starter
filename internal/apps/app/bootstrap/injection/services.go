package injection

import (
	"os"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	authservice "github.com/vadxq/go-rest-starter/internal/core/auth/service"
	userservice "github.com/vadxq/go-rest-starter/internal/core/user/service"
	"github.com/vadxq/go-rest-starter/internal/platform/config"
	"github.com/vadxq/go-rest-starter/pkg/cache"
	"github.com/vadxq/go-rest-starter/pkg/jwt"
	"github.com/vadxq/go-rest-starter/pkg/logger"
)

// Services 所有服务的集合
// 包含所有业务逻辑层对象，处理核心业务规则
type Services struct {
	// 用户相关业务逻辑
	UserService userservice.UserService

	// 认证相关业务逻辑
	AuthService authservice.AuthService

	// 可以在此添加更多服务...
	// ProductService userservice.ProductService
	// OrderService userservice.OrderService
}

// InitServices 初始化所有服务
// 这是依赖注入的第二层，依赖于仓库层
func InitServices(
	repos *Repositories,
	validate *validator.Validate,
	db *gorm.DB,
	config *config.AppConfig,
	cacheInstance cache.Cache,
	log logger.Logger,
) *Services {
	// 参数验证
	if repos == nil {
		log.Error("仓库依赖不能为空")
		os.Exit(1)
	}
	if validate == nil {
		log.Error("验证器不能为空")
		os.Exit(1)
	}
	if db == nil {
		log.Error("数据库连接不能为空")
		os.Exit(1)
	}
	if config == nil {
		log.Error("配置不能为空")
		os.Exit(1)
	}

	// 创建JWT配置
	jwtConfig := createJWTConfig(config, log)

	// 创建所有服务实例
	userService := userservice.NewUserService(repos.UserRepo, validate, db, cacheInstance)
	authService := authservice.NewAuthService(repos.UserRepo, validate, db, jwtConfig, cacheInstance)

	// 返回服务集合
	return &Services{
		UserService: userService,
		AuthService: authService,
	}
}

// createJWTConfig 从应用配置创建JWT配置
// 这是一个辅助函数，用于创建JWT服务所需的配置
func createJWTConfig(config *config.AppConfig, log logger.Logger) *jwt.Config {
	if config.JWT.Secret == "" {
		log.Warn("JWT密钥为空，这可能导致安全问题")
	}

	return &jwt.Config{
		Secret:          config.JWT.Secret,
		AccessTokenExp:  config.JWT.AccessTokenExp,
		RefreshTokenExp: config.JWT.RefreshTokenExp,
		Issuer:          config.JWT.Issuer,
	}
}
