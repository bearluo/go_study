package repositories

import "gorm.io/gorm"

// RepositoryManager 数据访问层管理器
type RepositoryManager struct {
	User         UserRepository
	RefreshToken RefreshTokenRepositoryInterface
	// 未来可以添加其他实体的Repository
	// Product ProductRepository
	// Order   OrderRepository
}

// NewRepositoryManager 创建数据访问层管理器
func NewRepositoryManager(db *gorm.DB) *RepositoryManager {
	return &RepositoryManager{
		User:         NewUserRepository(db),
		RefreshToken: NewRefreshTokenRepository(db),
		// 未来可以添加其他实体的Repository
		// Product: NewProductRepository(db),
		// Order:   NewOrderRepository(db),
	}
}
