package models

import (
	"testing"
)

func TestUser_SetDefaultRole(t *testing.T) {
	user := &User{}
	user.SetDefaultRole()
	if user.Role != RoleUser {
		t.Errorf("默认角色应该是 %s，实际是 %s", RoleUser, user.Role)
	}
	user.Role = RoleAdmin
	user.SetDefaultRole()
	if user.Role != RoleAdmin {
		t.Errorf("已有角色不应该被覆盖，应该是 %s，实际是 %s", RoleAdmin, user.Role)
	}
}

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{"admin_role", RoleAdmin, true},
		{"user_role", RoleUser, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			result := user.IsAdmin()
			if result != tt.expected {
				t.Errorf("IsAdmin() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUser_IsUser(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{"admin_role", RoleAdmin, true},
		{"user_role", RoleUser, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			result := user.IsUser()
			if result != tt.expected {
				t.Errorf("IsUser() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUser_HasRole(t *testing.T) {
	tests := []struct {
		name         string
		userRole     string
		requiredRole string
		expected     bool
	}{
		{"admin_has_admin", RoleAdmin, RoleAdmin, true},
		{"admin_has_user", RoleAdmin, RoleUser, true},
		{"user_has_user", RoleUser, RoleUser, true},
		{"user_has_admin", RoleUser, RoleAdmin, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.userRole}
			result := user.HasRole(tt.requiredRole)
			if result != tt.expected {
				t.Errorf("HasRole(%s) = %v, want %v", tt.requiredRole, result, tt.expected)
			}
		})
	}
}

func TestUser_HasAnyRole(t *testing.T) {
	user := &User{Role: RoleUser}
	if !user.HasAnyRole(RoleUser, RoleAdmin) {
		t.Error("HasAnyRole() 应该返回 true，当用户有其中一个角色时")
	}
	if user.HasAnyRole(RoleAdmin) == false {
		// 只要有一个匹配就行
	} else if user.HasAnyRole("invalid") {
		t.Error("HasAnyRole() 应该返回 false，当用户没有任何一个角色时")
	}
}

func TestUser_HasAllRoles(t *testing.T) {
	user := &User{Role: RoleAdmin}
	if !user.HasAllRoles(RoleUser) {
		t.Error("HasAllRoles() 应该返回 true，当用户有所有角色时")
	}
	user.Role = RoleUser
	if user.HasAllRoles(RoleUser, RoleAdmin) {
		t.Error("HasAllRoles() 应该返回 false，当用户没有所有角色时")
	}
}

func TestUser_GetRoleLevel(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected int
		hasError bool
	}{
		{"admin_level", RoleAdmin, RoleLevelAdmin, false},
		{"user_level", RoleUser, RoleLevelUser, false},
		{"invalid_role", "invalid", -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			level, err := user.GetRoleLevel()
			if tt.hasError {
				if err == nil {
					t.Error("GetRoleLevel() 应该返回错误，当角色无效时")
				}
			} else {
				if err != nil {
					t.Errorf("GetRoleLevel() 不应该返回错误: %v", err)
				}
				if level != tt.expected {
					t.Errorf("GetRoleLevel() = %d, want %d", level, tt.expected)
				}
			}
		})
	}
}

func TestUser_GetRoleDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected string
	}{
		{"admin_display", RoleAdmin, "管理员"},
		{"user_display", RoleUser, "用户"},
		{"invalid_role", "invalid", "invalid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			result := user.GetRoleDisplayName()
			if result != tt.expected {
				t.Errorf("GetRoleDisplayName() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestUser_GetRoleColor(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected string
	}{
		{"admin_color", RoleAdmin, "#dc3545"},
		{"user_color", RoleUser, "#007bff"},
		{"invalid_role", "invalid", "#6c757d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			result := user.GetRoleColor()
			if result != tt.expected {
				t.Errorf("GetRoleColor() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestUser_ValidateRole(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{"valid_admin", RoleAdmin, true},
		{"valid_user", RoleUser, true},
		{"invalid_role", "invalid", false},
		{"empty_role", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			result := user.ValidateRole()
			if result != tt.expected {
				t.Errorf("ValidateRole() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRoleConstants(t *testing.T) {
	if RoleAdmin != "admin" {
		t.Errorf("RoleAdmin 应该是 'admin'，实际是 '%s'", RoleAdmin)
	}
	if RoleUser != "user" {
		t.Errorf("RoleUser 应该是 'user'，实际是 '%s'", RoleUser)
	}
}

func TestRoleLevelConstants(t *testing.T) {
	if RoleLevelAdmin != 2 {
		t.Errorf("RoleLevelAdmin 应该是 2，实际是 %d", RoleLevelAdmin)
	}
	if RoleLevelUser != 1 {
		t.Errorf("RoleLevelUser 应该是 1，实际是 %d", RoleLevelUser)
	}
}

func TestRoleLevelMap(t *testing.T) {
	if RoleLevelMap[RoleAdmin] != RoleLevelAdmin {
		t.Errorf("RoleLevelMap[RoleAdmin] 应该是 %d，实际是 %d", RoleLevelAdmin, RoleLevelMap[RoleAdmin])
	}
	if RoleLevelMap[RoleUser] != RoleLevelUser {
		t.Errorf("RoleLevelMap[RoleUser] 应该是 %d，实际是 %d", RoleLevelUser, RoleLevelMap[RoleUser])
	}
}
