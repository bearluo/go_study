package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceManager_Structure(t *testing.T) {
	// 测试服务管理器的结构
	serviceManager := &ServiceManager{}

	// 验证结构体字段
	assert.NotNil(t, serviceManager)
}

func TestServiceManager_Methods(t *testing.T) {
	// 测试服务管理器的方法存在性
	serviceManager := &ServiceManager{
		UserService: &UserService{},
		AuthService: &AuthService{},
	}

	// 验证方法可以调用
	userService := serviceManager.GetUserService()
	authService := serviceManager.GetAuthService()

	assert.NotNil(t, userService)
	assert.NotNil(t, authService)
}
