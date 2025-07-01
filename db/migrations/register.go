package migrations

import "gorm.io/gorm"

// RegisterAllMigrations 注册所有迁移
func RegisterAllMigrations(db *gorm.DB) *MigrationManager {
	manager := NewMigrationManager(db)

	// 注册所有迁移（按版本号顺序）
	manager.RegisterMigration(&CreateUsersTableMigration{})
	manager.RegisterMigration(&CreateRefreshTokenTableMigration{})
	manager.RegisterMigration(&UpdateUsersTablePasswordLenMigration{})

	return manager
}
