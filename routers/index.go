package routers

import (
	"go-study/db/repositories"
	"go-study/middleware"
	"go-study/services"

	"github.com/labstack/echo/v4"
)

// SetupRoutes 统一管理路由
func SetupRoutes(e *echo.Echo, repoManager *repositories.RepositoryManager, serviceManager *services.ServiceManager, middlewareManager *middleware.MiddlewareManager) {

	// 设置全局中间件
	middlewareManager.SetupGlobalMiddlewares(e)

	// 设置各模块路由
	SetupUserRoutes(e, serviceManager, middlewareManager)
	SetupAuthRoutes(e, serviceManager, middlewareManager)
}
