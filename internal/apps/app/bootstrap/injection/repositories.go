package injection

import (
	"os"

	"gorm.io/gorm"

	userrepo "github.com/vadxq/go-rest-starter/internal/core/user/repository"
	"github.com/vadxq/go-rest-starter/pkg/logger"
)

// Repositories 所有仓库的集合
// 包含所有数据访问层对象，负责与数据源交互
type Repositories struct {
	// 用户数据访问对象
	UserRepo userrepo.UserRepository

	// 可以在此添加更多仓库...
	// ProductRepo repository.ProductRepository
	// OrderRepo repository.OrderRepository
}

// InitRepositories 初始化所有仓库
// 这是依赖注入的第一层，负责创建所有数据访问对象
func InitRepositories(db *gorm.DB, log logger.Logger) *Repositories {
	// 参数验证
	if db == nil {
		log.Error("数据库连接不能为空")
		os.Exit(1)
	}

	// 创建所有仓库实例
	userRepo := userrepo.NewUserRepository(db)

	// 返回仓库集合
	return &Repositories{
		UserRepo: userRepo,
	}
}
