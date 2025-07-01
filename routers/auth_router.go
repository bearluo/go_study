package routers

import (
	handles "go-study/handlers"
	"go-study/middleware"
	"go-study/services"

	"github.com/labstack/echo/v4"
)

// SetupAuthRoutes 设置认证路由
func SetupAuthRoutes(e *echo.Echo, serviceManager *services.ServiceManager, middlewareManager *middleware.MiddlewareManager) {
	// 创建处理器
	authHandler := handles.NewAuthHandler(serviceManager.GetAuthService())

	// 认证路由组
	auth := e.Group("/auth")

	// 公开路由（不需要认证）
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)            // 用户登录
	auth.POST("/validate", authHandler.ValidateToken) // 验证令牌

	// 受保护的路由（需要认证）
	protected := e.Group("/api/auth")
	protected.Use(middlewareManager.RequireAuth())
	{
		protected.GET("/profile", authHandler.GetProfile)     // 获取用户信息
		protected.POST("/refresh", authHandler.RefreshToken)  // 刷新令牌
		protected.POST("/logout", authHandler.Logout)         // 登出
		protected.POST("/logout-all", authHandler.LogoutAll)  // 撤销所有令牌
		protected.GET("/validate", authHandler.ValidateToken) // 验证令牌
	}
}
