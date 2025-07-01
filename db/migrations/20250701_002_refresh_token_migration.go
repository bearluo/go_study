package migrations

import (
	"go-study/db/models"

	"gorm.io/gorm"
)

// CreateRefreshTokenTableMigration 创建刷新令牌表迁移
type CreateRefreshTokenTableMigration struct{}

// Up 执行迁移
func (m *CreateRefreshTokenTableMigration) Up(db *gorm.DB) error {
	return db.AutoMigrate(&models.RefreshToken{})
}

// Down 回滚迁移
func (m *CreateRefreshTokenTableMigration) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.RefreshToken{})
}

// Version 获取版本号
func (m *CreateRefreshTokenTableMigration) Version() string {
	return "2025_07_01_000002"
}

// Name 获取迁移名称
func (m *CreateRefreshTokenTableMigration) Name() string {
	return "create_refresh_tokens_table"
}
