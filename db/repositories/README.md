# 数据访问层 (DAL) 架构说明

## 概述

本项目实现了完整的数据访问层 (Data Access Layer, DAL)，采用分层架构设计，提供了清晰的数据访问接口和实现。

## 架构层次

```
Handlers (控制器层)
    ↓
Services (业务逻辑层)
    ↓
Repositories (数据访问层)
    ↓
Models (数据模型层)
    ↓
Database (数据库)
```

## 文件结构

```
db/repositories/
├── base_repository.go      # 基础数据访问层接口和实现
├── user_repository.go      # 用户数据访问层接口和实现
├── repository_manager.go   # 数据访问层管理器
├── example_usage.go        # 使用示例
└── README.md              # 说明文档
```

## 核心组件

### 1. BaseRepository (基础数据访问层)

提供了通用的 CRUD 操作接口：

```go
type BaseRepository[T any] interface {
    Create(entity *T) error
    GetByID(id uint) (*T, error)
    GetAll() ([]T, error)
    Update(entity *T) error
    Delete(id uint) error
    GetDB() *gorm.DB
}
```

### 2. UserRepository (用户数据访问层)

继承基础接口，并添加用户特定的操作：

```go
type UserRepository interface {
    BaseRepository[models.User]
    GetByEmail(email string) (*models.User, error)
    ExistsByEmail(email string) (bool, error)
}
```

### 3. RepositoryManager (数据访问层管理器)

统一管理所有 Repository 实例：

```go
type RepositoryManager struct {
    User UserRepository
    // 未来可以添加其他实体的Repository
}
```

## 使用方法

### 1. 初始化

```go
// 在 main.go 中
repoManager := repositories.NewRepositoryManager(db.DB)
```

### 2. 在服务层中使用

```go
// 创建用户服务
userService := services.NewUserService(repoManager.User)

// 使用服务层方法
user := &models.User{
    Name:     "张三",
    Email:    "zhangsan@example.com",
    Password: "123456",
}
err := userService.Create(user)
```

### 3. 在处理器中使用

```go
// 创建处理器
userHandler := handles.NewUserHandler(userService)

// 处理器会自动使用服务层进行数据操作
```

## API 路由

### 用户相关 API

- `POST /users` - 创建用户
- `GET /users` - 获取所有用户
- `GET /users/:id` - 根据ID获取用户
- `PUT /users/:id` - 更新用户
- `DELETE /users/:id` - 删除用户

### 向后兼容的路由

- `POST /user/create` - 创建用户（保持向后兼容）

## 扩展新实体

要添加新的实体（如 Product），需要：

1. 创建模型文件 `db/models/product.go`
2. 创建 Repository 文件 `db/repositories/product_repository.go`
3. 在 `repository_manager.go` 中添加 Product 字段
4. 创建服务层 `services/product_service.go`
5. 创建处理器 `handlers/product_handler.go`
6. 在路由中添加相应的路由

## 优势

1. **分层清晰**: 每层职责明确，便于维护和测试
2. **依赖注入**: 通过接口实现松耦合
3. **可扩展性**: 易于添加新的实体和功能
4. **可测试性**: 每层都可以独立测试
5. **代码复用**: 基础 Repository 提供通用功能

## 错误处理

所有数据访问操作都包含适当的错误处理：

- 数据库连接错误
- 记录不存在错误
- 唯一约束违反错误
- 其他数据库操作错误

## 事务支持

可以通过 `GetDB()` 方法获取底层的 GORM 实例来支持事务操作：

```go
tx := userRepo.GetDB().Begin()
// 执行事务操作
tx.Commit()
``` 