package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// 用户角色常量
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// 角色权限级别
const (
	RoleLevelUser  = 1
	RoleLevelAdmin = 2
)

// 角色映射表
var RoleLevelMap = map[string]int{
	RoleUser:  RoleLevelUser,
	RoleAdmin: RoleLevelAdmin,
}

// Time 自定义时间类型，实现 driver.Valuer 和 sql.Scanner 接口
type Time struct {
	time.Time
}

// Value 实现 driver.Valuer 接口
func (t Time) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Time.Format("2006-01-02 15:04:05"), nil
}

// Scan 实现 sql.Scanner 接口
func (t *Time) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		// 处理 MySQL 返回的 []uint8 类型
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			// 尝试其他时间格式
			parsedTime, err = time.Parse("2006-01-02T15:04:05Z07:00", string(v))
			if err != nil {
				return fmt.Errorf("无法解析时间格式: %v", err)
			}
		}
		t.Time = parsedTime
	case string:
		parsedTime, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			parsedTime, err = time.Parse("2006-01-02T15:04:05Z07:00", v)
			if err != nil {
				return fmt.Errorf("无法解析时间格式: %v", err)
			}
		}
		t.Time = parsedTime
	case time.Time:
		t.Time = v
	default:
		return fmt.Errorf("不支持的时间类型: %T", value)
	}

	return nil
}

// 定义一个User 结构体,用来表示user表
// User 结构体表示用户表
type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`                                                            // 主键，自动递增
	Name      string `gorm:"size:10;not null"`                                                                    // 用户名，最大长度10，不能为空
	Email     string `gorm:"unique;size:20;not null"`                                                             // 邮箱，唯一索引，最大长度20，不能为空
	Password  string `gorm:"size:100;not null"`                                                                   // 密码，不能为空
	Role      string `gorm:"size:10;not null"`                                                                    // 角色，最大长度10，不能为空
	CreatedAt Time   `gorm:"autoCreateTime;type:timestamp;default:CURRENT_TIMESTAMP"`                             // 在创建时，如果该字段值为零值，则使用当前时间填充
	UpdatedAt Time   `gorm:"autoUpdateTime;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 在创建时该字段值为零值或者在更新时，使用当前时间戳秒数填充
}

// SetDefaultRole 设置默认角色
func (u *User) SetDefaultRole() {
	if u.Role == "" {
		u.Role = RoleUser
	}
}

// IsAdmin 检查是否为管理员
func (u *User) IsAdmin() bool {
	return u.HasRole(RoleAdmin)
}

// IsUser 检查是否为普通用户或更高权限
func (u *User) IsUser() bool {
	return u.HasRole(RoleUser)
}

// HasRole 检查是否具有指定角色
func (u *User) HasRole(role string) bool {
	userLevel, exists := RoleLevelMap[u.Role]
	if !exists {
		return false
	}
	requiredLevel, exists := RoleLevelMap[role]
	if !exists {
		return false
	}
	return userLevel >= requiredLevel
}

// HasAnyRole 检查是否具有任意一个指定角色
func (u *User) HasAnyRole(roles ...string) bool {
	for _, role := range roles {
		if u.HasRole(role) {
			return true
		}
	}
	return false
}

// HasAllRoles 检查是否具有所有指定角色
func (u *User) HasAllRoles(roles ...string) bool {
	for _, role := range roles {
		if !u.HasRole(role) {
			return false
		}
	}
	return true
}

// GetRoleLevel 获取角色级别
func (u *User) GetRoleLevel() (int, error) {
	level, exists := RoleLevelMap[u.Role]
	if !exists {
		return -1, fmt.Errorf("未知角色: %s", u.Role)
	}
	return level, nil
}

// GetRoleDisplayName 获取角色显示名称
func (u *User) GetRoleDisplayName() string {
	displayNames := map[string]string{
		RoleUser:  "用户",
		RoleAdmin: "管理员",
	}
	if displayName, exists := displayNames[u.Role]; exists {
		return displayName
	}
	return u.Role
}

// GetRoleColor 获取角色颜色
func (u *User) GetRoleColor() string {
	colors := map[string]string{
		RoleUser:  "#007bff",
		RoleAdmin: "#dc3545",
	}
	if color, exists := colors[u.Role]; exists {
		return color
	}
	return "#6c757d"
}

// ValidateRole 验证角色是否有效
func (u *User) ValidateRole() bool {
	_, exists := RoleLevelMap[u.Role]
	return exists
}

// BeforeCreate GORM 钩子：创建前设置默认角色
func (u *User) BeforeCreate() error {
	u.SetDefaultRole()
	return nil
}
