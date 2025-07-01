package repositories

import "gorm.io/gorm"

// BaseRepository 基础数据访问层接口
type BaseRepository[T any] interface {
	Create(entity *T) error
	GetByID(id uint) (*T, error)
	GetAll() ([]T, error)
	Update(entity *T) error
	Delete(id uint) error
	GetDB() *gorm.DB
}

// baseRepository 基础数据访问层实现
type baseRepository[T any] struct {
	db *gorm.DB
}

// NewBaseRepository 创建基础数据访问层实例
func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &baseRepository[T]{db: db}
}

// Create 创建实体
func (r *baseRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// GetByID 根据ID获取实体
func (r *baseRepository[T]) GetByID(id uint) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// GetAll 获取所有实体
func (r *baseRepository[T]) GetAll() ([]T, error) {
	var entities []T
	err := r.db.Find(&entities).Error
	return entities, err
}

// Update 更新实体
func (r *baseRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Delete 删除实体
func (r *baseRepository[T]) Delete(id uint) error {
	var entity T
	return r.db.Delete(&entity, id).Error
}

// GetDB 获取数据库连接
func (r *baseRepository[T]) GetDB() *gorm.DB {
	return r.db
}
