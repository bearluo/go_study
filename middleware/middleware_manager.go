package middleware

import (
	"go-study/services"

	"github.com/labstack/echo/v4"
)

// MiddlewareManager 中间件管理器
type MiddlewareManager struct {
	AuthMiddleware *AuthMiddleware
}

// NewMiddlewareManager 创建中间件管理器
func NewMiddlewareManager(serviceManager *services.ServiceManager) *MiddlewareManager {
	return &MiddlewareManager{
		AuthMiddleware: NewAuthMiddleware(serviceManager.GetAuthService()),
	}
}

// GetAuthMiddleware 获取认证中间件
func (mm *MiddlewareManager) GetAuthMiddleware() *AuthMiddleware {
	return mm.AuthMiddleware
}

// SetupGlobalMiddlewares 设置全局中间件
func (mm *MiddlewareManager) SetupGlobalMiddlewares(e *echo.Echo) {
	// 添加全局中间件
	e.Use(mm.AuthMiddleware.OptionalAuth())
}

// SetupAuthMiddlewares 设置认证相关中间件
func (mm *MiddlewareManager) SetupAuthMiddlewares(e *echo.Echo) {
	// 这里可以添加认证相关的全局中间件
	// 例如：日志中间件、CORS中间件等
}

// RequireAuth 获取要求认证的中间件
func (mm *MiddlewareManager) RequireAuth() echo.MiddlewareFunc {
	return mm.AuthMiddleware.RequireAuth()
}

// RequireAdmin 获取要求管理员权限的中间件
func (mm *MiddlewareManager) RequireAdmin() echo.MiddlewareFunc {
	return mm.AuthMiddleware.RequireAdmin()
}

// RequireUser 获取要求用户权限的中间件
func (mm *MiddlewareManager) RequireUser() echo.MiddlewareFunc {
	return mm.AuthMiddleware.RequireUser()
}

// RequireRole 获取要求特定角色的中间件
func (mm *MiddlewareManager) RequireRole(role string) echo.MiddlewareFunc {
	return mm.AuthMiddleware.RequireRole(role)
}

// OptionalAuth 获取可选认证的中间件
func (mm *MiddlewareManager) OptionalAuth() echo.MiddlewareFunc {
	return mm.AuthMiddleware.OptionalAuth()
}
