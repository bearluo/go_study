package migrations

import (
	"fmt"

	"go-study/db/models"

	"gorm.io/gorm"
)

// Migration 迁移记录结构
type Migration struct {
	ID        uint        `gorm:"primaryKey;autoIncrement"`
	Version   string      `gorm:"unique;not null"` // 迁移版本号
	Name      string      `gorm:"not null"`        // 迁移名称
	CreatedAt models.Time `gorm:"autoCreateTime"`
}

// MigrationInterface 迁移接口
type MigrationInterface interface {
	Up(db *gorm.DB) error   // 执行迁移
	Down(db *gorm.DB) error // 回滚迁移
	Version() string        // 获取版本号
	Name() string           // 获取迁移名称
}

// MigrationManager 迁移管理器
type MigrationManager struct {
	db         *gorm.DB
	migrations []MigrationInterface
}

// NewMigrationManager 创建迁移管理器
func NewMigrationManager(db *gorm.DB) *MigrationManager {
	return &MigrationManager{
		db:         db,
		migrations: make([]MigrationInterface, 0),
	}
}

// RegisterMigration 注册迁移
func (m *MigrationManager) RegisterMigration(migration MigrationInterface) {
	m.migrations = append(m.migrations, migration)
}

// InitMigrationTable 初始化迁移表
func (m *MigrationManager) InitMigrationTable() error {
	return m.db.AutoMigrate(&Migration{})
}

// GetExecutedMigrations 获取已执行的迁移
func (m *MigrationManager) GetExecutedMigrations() ([]Migration, error) {
	var migrations []Migration
	err := m.db.Find(&migrations).Error
	return migrations, err
}

// IsMigrationExecuted 检查迁移是否已执行
func (m *MigrationManager) IsMigrationExecuted(version string) bool {
	var count int64
	m.db.Model(&Migration{}).Where("version = ?", version).Count(&count)
	return count > 0
}

// ExecuteMigration 执行单个迁移
func (m *MigrationManager) ExecuteMigration(migration MigrationInterface) error {
	version := migration.Version()

	// 检查是否已执行
	if m.IsMigrationExecuted(version) {
		fmt.Printf("迁移 %s 已执行，跳过\n", version)
		return nil
	}

	// 开始事务
	tx := m.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 执行迁移
	if err := migration.Up(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("执行迁移 %s 失败: %v", version, err)
	}

	// 记录迁移
	migrationRecord := Migration{
		Version: version,
		Name:    migration.Name(),
	}
	if err := tx.Create(&migrationRecord).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("记录迁移 %s 失败: %v", version, err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交迁移 %s 失败: %v", version, err)
	}

	fmt.Printf("迁移 %s 执行成功\n", version)
	return nil
}

// Migrate 执行所有迁移
func (m *MigrationManager) Migrate() error {
	// 初始化迁移表
	if err := m.InitMigrationTable(); err != nil {
		return fmt.Errorf("初始化迁移表失败: %v", err)
	}

	// 执行所有迁移
	for _, migration := range m.migrations {
		if err := m.ExecuteMigration(migration); err != nil {
			return err
		}
	}

	fmt.Println("所有迁移执行完成")
	return nil
}

// Rollback 回滚最后一个迁移
func (m *MigrationManager) Rollback() error {
	// 获取最后一个已执行的迁移
	var lastMigration Migration
	if err := m.db.Order("id desc").First(&lastMigration).Error; err != nil {
		return fmt.Errorf("没有找到已执行的迁移: %v", err)
	}

	// 查找对应的迁移对象
	var targetMigration MigrationInterface
	for _, migration := range m.migrations {
		if migration.Version() == lastMigration.Version {
			targetMigration = migration
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("未找到版本 %s 的迁移对象", lastMigration.Version)
	}

	// 开始事务
	tx := m.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 执行回滚
	if err := targetMigration.Down(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("回滚迁移 %s 失败: %v", lastMigration.Version, err)
	}

	// 删除迁移记录
	if err := tx.Delete(&lastMigration).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除迁移记录失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交回滚失败: %v", err)
	}

	fmt.Printf("迁移 %s 回滚成功\n", lastMigration.Version)
	return nil
}

// Status 显示迁移状态
func (m *MigrationManager) Status() error {
	executedMigrations, err := m.GetExecutedMigrations()
	if err != nil {
		return fmt.Errorf("获取已执行迁移失败: %v", err)
	}

	fmt.Println("迁移状态:")
	fmt.Println("版本号\t\t状态\t\t迁移名称")
	fmt.Println("----------------------------------------")

	// 创建已执行迁移的映射
	executedMap := make(map[string]bool)
	for _, migration := range executedMigrations {
		executedMap[migration.Version] = true
	}

	// 显示所有迁移的状态
	for _, migration := range m.migrations {
		status := "未执行"
		if executedMap[migration.Version()] {
			status = "已执行"
		}
		fmt.Printf("%s\t%s\t%s\n", migration.Version(), status, migration.Name())
	}

	return nil
}
