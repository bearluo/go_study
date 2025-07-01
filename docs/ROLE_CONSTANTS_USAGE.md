# 角色常量使用指南

本文档介绍了项目中角色常量的定义和使用方法。

## 角色常量定义

### 在 `db/models/user.go` 中定义的角色常量

```go
const (
    RoleAdmin = "admin"     // 管理员角色
    RoleUser  = "user"      // 普通用户角色
)

const (
    RoleLevelUser  = 1 // 普通用户权限级别
    RoleLevelAdmin = 2 // 管理员权限级别
)
```

## 角色权限层级

```
管理员 (admin) - 最高权限
    ↓
用户 (user) - 基础权限
```

## 使用方法

### 1. 在 User 模型中使用角色方法

```go
user := &models.User{
    Name: "张三",
    Role: models.RoleUser,
}

// 检查角色权限
if user.IsAdmin() {
    // 管理员逻辑
}

if user.HasRole(models.RoleUser) {
    // 用户或更高权限逻辑
}

displayName := user.GetRoleDisplayName() // "用户"
roleColor := user.GetRoleColor()         // "#007bff"
```

### 2. 在中间件中使用角色验证

```go
// 要求管理员权限
adminGroup := e.Group("/api/admin")
adminGroup.Use(authMiddleware.RequireAdmin())

// 要求用户或更高权限
userGroup := e.Group("/api/user")
userGroup.Use(authMiddleware.RequireUser())

// 要求特定角色
customGroup := e.Group("/api/custom")
customGroup.Use(authMiddleware.RequireRole(models.RoleUser))
```

### 3. 在工具函数中使用角色验证

```go
import "go-study/utils"

if utils.IsAdmin(userRole) {
    // 管理员逻辑
}

if utils.HasRole(userRole, models.RoleUser) {
    // 用户或更高权限逻辑
}

roleInfo, err := utils.GetRoleInfo(userRole)
if err == nil {
    fmt.Printf("角色: %s, 级别: %d, 描述: %s\n", roleInfo.Name, roleInfo.Level, roleInfo.Desc)
}

displayName := utils.GetRoleDisplayName(userRole)
color := utils.GetRoleColor(userRole)
```

### 4. 在业务逻辑中使用角色判断

```go
func (s *UserService) UpdateUserRole(userID uint, newRole string) error {
    if !utils.ValidateRole(newRole) {
        return errors.New("无效的角色")
    }
    currentUser := getCurrentUser()
    if !currentUser.IsAdmin() {
        return errors.New("权限不足")
    }
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return err
    }
    user.Role = newRole
    return s.userRepo.Update(user)
}
```

## 角色验证最佳实践

### 1. 始终验证角色有效性

```go
if !utils.ValidateRole(role) {
    return errors.New("无效的角色")
}
```

### 2. 使用角色常量而不是字符串

```go
if user.Role == models.RoleAdmin {
    // 管理员逻辑
}
```

### 3. 使用层级权限检查

```go
if user.HasRole(models.RoleUser) {
    // 用户或管理员可以执行
}
```

### 4. 在数据库查询中使用角色过滤

```go
admins, err := userRepo.FindByRole(models.RoleAdmin)
users, err := userRepo.FindByRole(models.RoleUser)
```

## 扩展角色系统

### 添加新角色

1. 在 `db/models/user.go` 中添加新角色常量：

```go
const (
    RoleAdmin = "admin"
    RoleUser  = "user"
    RoleVIP   = "vip"        // 新增 VIP 角色
)
```

2. 添加对应的权限级别：

```go
const (
    RoleLevelUser  = 1
    RoleLevelVIP   = 2        // 新增 VIP 级别
    RoleLevelAdmin = 3        // 调整管理员级别
)
```

3. 更新角色映射表：

```go
var RoleLevelMap = map[string]int{
    RoleUser:  RoleLevelUser,
    RoleVIP:   RoleLevelVIP,
    RoleAdmin: RoleLevelAdmin,
}
```

4. 在工具函数中添加新角色的处理逻辑。

## 注意事项

1. **角色常量是全局的**：一旦定义，整个项目都应该使用这些常量
2. **权限级别要合理**：确保权限级别的数字大小关系正确
3. **向后兼容**：添加新角色时要注意不要破坏现有的权限逻辑
4. **测试覆盖**：确保所有角色相关的逻辑都有相应的测试用例
5. **文档更新**：角色系统变更时要及时更新相关文档

## 相关文件

- `db/models/user.go` - 角色常量定义和 User 模型方法
- `utils/role_utils.go` - 角色相关的工具函数
- `utils/constants.go` - 其他系统常量
- `middleware/auth_middleware.go` - 认证中间件中的角色验证
- `routers/auth_router.go` - 路由中的角色权限控制 