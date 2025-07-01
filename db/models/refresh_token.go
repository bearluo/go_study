package models

import (
	"time"
)

// RefreshToken 结构体表示刷新令牌表
type RefreshToken struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`                                                            // 主键，自动递增
	UserID    uint   `gorm:"not null;index"`                                                                      // 用户ID，不能为空，建立索引
	Token     string `gorm:"size:255;not null;uniqueIndex"`                                                       // 刷新令牌，不能为空，唯一索引
	ExpiresAt Time   `gorm:"not null;type:timestamp"`                                                             // 过期时间，不能为空
	IsRevoked bool   `gorm:"default:false"`                                                                       // 是否已撤销，默认为false
	CreatedAt Time   `gorm:"autoCreateTime;type:timestamp;default:CURRENT_TIMESTAMP"`                             // 创建时间
	UpdatedAt Time   `gorm:"autoUpdateTime;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 更新时间
}

// IsExpired 检查刷新令牌是否过期
func (rt *RefreshToken) IsExpired() bool {
	return rt.ExpiresAt.Time.Before(time.Now())
}

// IsValid 检查刷新令牌是否有效（未过期且未撤销）
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.IsRevoked
}
