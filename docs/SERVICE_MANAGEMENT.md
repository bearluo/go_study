# 服务管理和中间件管理

本文档介绍项目中的服务管理器和中间件管理器的使用方式。

## 服务管理器 (ServiceManager)

服务管理器用于统一管理所有的服务实例，提供便捷的访问方式。

### 创建服务管理器

```go
// 创建仓库管理器
repoManager := repositories.NewRepositoryManager(db.DB)

// 创建服务管理器
serviceManager := services.NewServiceManager(repoManager)
```

### 获取服务实例

```go
// 获取用户服务
userService := serviceManager.GetUserService()

// 获取认证服务
authService := serviceManager.GetAuthService()

// 获取服务接口（用于依赖注入和测试）
userServiceInterface := serviceManager.GetUserServiceInterface()
authServiceInterface := serviceManager.GetAuthServiceInterface()
```

### 服务管理器结构

```go
type ServiceManager struct {
    UserService *UserService
    AuthService *AuthService
}
```

## 中间件管理器 (MiddlewareManager)

中间件管理器用于统一管理所有的中间件实例，提供便捷的访问方式。

### 创建中间件管理器

```go
// 创建服务管理器
serviceManager := services.NewServiceManager(repoManager)

// 创建中间件管理器
middlewareManager := middleware.NewMiddlewareManager(serviceManager)
```

### 获取中间件实例

```go
// 获取认证中间件
authMiddleware := middlewareManager.GetAuthMiddleware()

// 获取各种认证中间件函数
requireAuth := middlewareManager.RequireAuth()
requireAdmin := middlewareManager.RequireAdmin()
requireUser := middlewareManager.RequireUser()
requireRole := middlewareManager.RequireRole("admin")
optionalAuth := middlewareManager.OptionalAuth()
```

### 设置全局中间件

```go
// 设置全局中间件
middlewareManager.SetupGlobalMiddlewares(e)

// 设置认证相关中间件
middlewareManager.SetupAuthMiddlewares(e)
```

### 中间件管理器结构

```go
type MiddlewareManager struct {
    AuthMiddleware *AuthMiddleware
}
```

## 在路由中使用

### 更新路由设置

```go
func SetupRoutes(e *echo.Echo, repoManager *repositories.RepositoryManager) {
    // 创建服务管理器
    serviceManager := services.NewServiceManager(repoManager)
    
    // 创建中间件管理器
    middlewareManager := middleware.NewMiddlewareManager(serviceManager)
    
    // 设置全局中间件
    middlewareManager.SetupGlobalMiddlewares(e)
    
    // 设置各模块路由
    SetupUserRoutes(e, serviceManager, middlewareManager)
    SetupAuthRoutes(e, serviceManager, middlewareManager)
}
```

### 在具体路由中使用

```go
func SetupUserRoutes(e *echo.Echo, serviceManager *services.ServiceManager, middlewareManager *middleware.MiddlewareManager) {
    // 创建处理器
    userHandler := handles.NewUserHandler(
        serviceManager.GetUserServiceInterface(),
        serviceManager.GetAuthServiceInterface(),
    )

    // 设置路由
    e.POST("/api/auth/register", userHandler.Register)
    
    // 使用中间件保护路由
    protected := e.Group("/api/users")
    protected.Use(middlewareManager.RequireAuth())
    {
        protected.GET("/profile", userHandler.GetProfile)
    }
}
```

## 优势

### 1. 统一管理
- 所有服务实例在一个地方创建和管理
- 避免重复创建服务实例
- 便于依赖注入和测试

### 2. 类型安全
- 提供类型安全的访问方法
- 编译时检查类型错误
- 支持接口和具体类型的访问

### 3. 易于扩展
- 新增服务时只需在管理器中添加
- 不影响现有代码
- 便于维护和重构

### 4. 测试友好
- 可以轻松替换为 mock 对象
- 支持接口注入
- 便于单元测试

## 最佳实践

1. **单一职责**: 每个管理器只负责管理相关的组件
2. **依赖注入**: 通过构造函数注入依赖，而不是在内部创建
3. **接口优先**: 优先使用接口类型，便于测试和扩展
4. **错误处理**: 在管理器中添加适当的错误处理逻辑
5. **文档化**: 为每个管理器提供清晰的文档和使用示例 