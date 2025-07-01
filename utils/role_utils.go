package utils

import (
	"fmt"
	"go-study/db/models"
)

// RoleInfo 角色信息结构
type RoleInfo struct {
	Name  string `json:"name"`  // 角色名称
	Level int    `json:"level"` // 角色级别
	Desc  string `json:"desc"`  // 角色描述
}

// 角色信息映射表
var RoleInfoMap = map[string]RoleInfo{
	models.RoleUser: {
		Name:  models.RoleUser,
		Level: models.RoleLevelUser,
		Desc:  "普通用户，可以发布内容和评论",
	},
	models.RoleAdmin: {
		Name:  models.RoleAdmin,
		Level: models.RoleLevelAdmin,
		Desc:  "管理员，拥有所有权限",
	},
}

// GetRoleLevel 获取角色级别
func GetRoleLevel(role string) (int, error) {
	level, exists := models.RoleLevelMap[role]
	if !exists {
		return -1, fmt.Errorf("未知角色: %s", role)
	}
	return level, nil
}

// GetRoleInfo 获取角色信息
func GetRoleInfo(role string) (*RoleInfo, error) {
	info, exists := RoleInfoMap[role]
	if !exists {
		return nil, fmt.Errorf("未知角色: %s", role)
	}
	return &info, nil
}

// HasRole 检查用户是否具有指定角色
func HasRole(userRole, requiredRole string) bool {
	userLevel, err := GetRoleLevel(userRole)
	if err != nil {
		return false
	}

	requiredLevel, err := GetRoleLevel(requiredRole)
	if err != nil {
		return false
	}

	return userLevel >= requiredLevel
}

// HasAnyRole 检查用户是否具有任意一个指定角色
func HasAnyRole(userRole string, requiredRoles ...string) bool {
	for _, role := range requiredRoles {
		if HasRole(userRole, role) {
			return true
		}
	}
	return false
}

// HasAllRoles 检查用户是否具有所有指定角色
func HasAllRoles(userRole string, requiredRoles ...string) bool {
	for _, role := range requiredRoles {
		if !HasRole(userRole, role) {
			return false
		}
	}
	return true
}

// IsAdmin 检查是否为管理员
func IsAdmin(role string) bool {
	return HasRole(role, models.RoleAdmin)
}

// IsUser 检查是否为普通用户或更高权限
func IsUser(role string) bool {
	return HasRole(role, models.RoleUser)
}

// GetValidRoles 获取所有有效角色列表
func GetValidRoles() []string {
	return []string{
		models.RoleUser,
		models.RoleAdmin,
	}
}

// GetRoleHierarchy 获取角色层级关系
func GetRoleHierarchy() map[string][]string {
	return map[string][]string{
		models.RoleUser: {
			models.RoleUser,
		},
		models.RoleAdmin: {
			models.RoleUser,
			models.RoleAdmin,
		},
	}
}

// GetSubordinateRoles 获取指定角色的下级角色
func GetSubordinateRoles(role string) ([]string, error) {
	hierarchy := GetRoleHierarchy()
	subordinates, exists := hierarchy[role]
	if !exists {
		return nil, fmt.Errorf("未知角色: %s", role)
	}
	return subordinates, nil
}

// ValidateRole 验证角色是否有效
func ValidateRole(role string) bool {
	_, exists := models.RoleLevelMap[role]
	return exists
}

// GetRoleDisplayName 获取角色显示名称
func GetRoleDisplayName(role string) string {
	displayNames := map[string]string{
		models.RoleUser:  "用户",
		models.RoleAdmin: "管理员",
	}

	if displayName, exists := displayNames[role]; exists {
		return displayName
	}
	return role
}

// GetRoleColor 获取角色对应的颜色（用于前端显示）
func GetRoleColor(role string) string {
	colors := map[string]string{
		models.RoleUser:  "#007bff",
		models.RoleAdmin: "#dc3545",
	}

	if color, exists := colors[role]; exists {
		return color
	}
	return "#6c757d" // 默认灰色
}
