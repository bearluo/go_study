package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddlewareManager_Structure(t *testing.T) {
	// 测试中间件管理器的结构
	middlewareManager := &MiddlewareManager{}

	// 验证结构体字段
	assert.NotNil(t, middlewareManager)
}

func TestMiddlewareManager_Methods(t *testing.T) {
	// 测试中间件管理器的方法存在性
	middlewareManager := &MiddlewareManager{
		AuthMiddleware: &AuthMiddleware{},
	}

	// 验证方法可以调用
	authMiddleware := middlewareManager.GetAuthMiddleware()
	assert.NotNil(t, authMiddleware)
}
