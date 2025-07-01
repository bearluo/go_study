package migrations

import (
	"gorm.io/gorm"
)

// UpdateUsersTablePasswordLenMigration 更新用户表密码长度
type UpdateUsersTablePasswordLenMigration struct{}

// Up 执行迁移
func (m *UpdateUsersTablePasswordLenMigration) Up(db *gorm.DB) error {
	return db.Exec("ALTER TABLE users MODIFY password VARCHAR(100) NOT NULL").Error
}

// Down 回滚迁移
func (m *UpdateUsersTablePasswordLenMigration) Down(db *gorm.DB) error {
	return db.Exec("ALTER TABLE users MODIFY password VARCHAR(18) NOT NULL").Error
}

// Version 获取版本号
func (m *UpdateUsersTablePasswordLenMigration) Version() string {
	return "2025_07_01_000003"
}

// Name 获取迁移名称
func (m *UpdateUsersTablePasswordLenMigration) Name() string {
	return "update_users_table_password_len"
}
