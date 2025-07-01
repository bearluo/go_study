package migrations

import (
	"go-study/db/models"

	"gorm.io/gorm"
)

// CreateUsersTableMigration 创建用户表迁移
type CreateUsersTableMigration struct{}

// Up 执行迁移
func (m *CreateUsersTableMigration) Up(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}

// Down 回滚迁移
func (m *CreateUsersTableMigration) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(&models.User{})
}

// Version 获取版本号
func (m *CreateUsersTableMigration) Version() string {
	return "2025_07_01_000001"
}

// Name 获取迁移名称
func (m *CreateUsersTableMigration) Name() string {
	return "create_users_table"
}
