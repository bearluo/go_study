# 数据库迁移系统

这是一个基于GORM的数据库迁移系统，支持版本化的数据库结构变更管理。

## 功能特性

- ✅ 版本化迁移管理
- ✅ 自动迁移记录
- ✅ 支持回滚操作
- ✅ 事务安全
- ✅ 迁移状态查询
- ✅ 命令行工具

## 目录结构

```
db/migrations/
├── migration.go              # 迁移管理器核心代码
├── register.go               # 迁移注册文件
├── README.md                 # 说明文档
├── 20240101_001_create_users_table.go    # 用户表迁移
├── 20240101_002_add_roles_table.go       # 角色表迁移
├── 20240101_003_add_role_id_to_users.go  # 添加用户角色字段迁移
└── 20240101_004_insert_default_roles.go  # 插入默认角色数据迁移
```

## 使用方法

### 1. 执行迁移

```bash
# 执行所有未执行的迁移
go run tools/migrate/main.go -action migrate

# 或者简写（migrate是默认操作）
go run tools/migrate/main.go
```

### 2. 查看迁移状态

```bash
go run tools/migrate/main.go -action status
```

### 3. 回滚最后一个迁移

```bash
go run tools/migrate/main.go -action rollback
```

### 4. 查看帮助

```bash
go run tools/migrate/main.go -help
```

## 创建新的迁移

### 1. 创建迁移文件

创建一个新的迁移文件，命名格式：`YYYYMMDD_XXX_description.go`

```go
package migrations

import (
    "gorm.io/gorm"
)

// YourMigrationName 迁移描述
type YourMigrationName struct{}

// Version 获取版本号
func (m *YourMigrationName) Version() string {
    return "20240101_005" // 使用下一个版本号
}

// Name 获取迁移名称
func (m *YourMigrationName) Name() string {
    return "your_migration_name"
}

// Up 执行迁移
func (m *YourMigrationName) Up(db *gorm.DB) error {
    // 在这里编写迁移逻辑
    return db.Exec(`
        CREATE TABLE IF NOT EXISTS your_table (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name VARCHAR(100) NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `).Error
}

// Down 回滚迁移
func (m *YourMigrationName) Down(db *gorm.DB) error {
    // 在这里编写回滚逻辑
    return db.Exec(`DROP TABLE IF EXISTS your_table`).Error
}
```

### 2. 注册迁移

在 `register.go` 文件中注册新的迁移：

```go
func RegisterAllMigrations() *MigrationManager {
    manager := NewMigrationManager(db.DB)

    // 按版本号顺序注册迁移
    manager.RegisterMigration(&CreateUsersTable{})
    manager.RegisterMigration(&AddRolesTable{})
    manager.RegisterMigration(&AddRoleIDToUsers{})
    manager.RegisterMigration(&InsertDefaultRoles{})
    manager.RegisterMigration(&YourMigrationName{}) // 添加新迁移

    return manager
}
```

## 迁移最佳实践

### 1. 版本号命名规范

- 使用时间戳格式：`YYYYMMDD_XXX`
- 例如：`20240101_001`, `20240101_002`
- 确保版本号唯一且按时间顺序递增

### 2. 迁移文件命名

- 使用描述性的文件名
- 格式：`YYYYMMDD_XXX_description.go`
- 例如：`20240101_001_create_users_table.go`

### 3. 迁移内容

- **Up方法**：执行数据库变更
- **Down方法**：回滚数据库变更
- 使用事务确保数据一致性
- 添加适当的错误处理

### 4. 数据迁移

- 对于数据插入，使用 `INSERT OR IGNORE` 避免重复
- 对于数据更新，考虑数据一致性
- 大量数据操作时考虑分批处理

## 示例迁移

### 创建表

```go
func (m *CreateUsersTable) Up(db *gorm.DB) error {
    return db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name VARCHAR(10) NOT NULL,
            email VARCHAR(20) UNIQUE NOT NULL,
            password VARCHAR(18) NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `).Error
}
```

### 添加字段

```go
func (m *AddRoleIDToUsers) Up(db *gorm.DB) error {
    return db.Exec(`ALTER TABLE users ADD COLUMN role_id INTEGER`).Error
}
```

### 插入数据

```go
func (m *InsertDefaultRoles) Up(db *gorm.DB) error {
    roles := []map[string]interface{}{
        {"name": "admin", "description": "系统管理员"},
        {"name": "user", "description": "普通用户"},
    }

    for _, role := range roles {
        if err := db.Exec(`
            INSERT OR IGNORE INTO roles (name, description) 
            VALUES (?, ?)
        `, role["name"], role["description"]).Error; err != nil {
            return err
        }
    }
    return nil
}
```

## 注意事项

1. **版本号唯一性**：确保每个迁移的版本号都是唯一的
2. **顺序性**：迁移按版本号顺序执行，不要修改已提交的迁移
3. **回滚安全**：确保Down方法能够正确回滚Up方法的操作
4. **事务处理**：迁移系统自动使用事务，确保数据一致性
5. **测试**：在生产环境使用前，先在测试环境验证迁移

## 故障排除

### 迁移失败

如果迁移失败，检查：
1. 数据库连接是否正常
2. SQL语法是否正确
3. 是否有权限问题
4. 版本号是否重复

### 回滚失败

如果回滚失败，检查：
1. Down方法的逻辑是否正确
2. 是否有依赖关系问题
3. 数据是否被其他操作修改

## 扩展功能

可以根据需要扩展以下功能：
- 批量回滚多个迁移
- 迁移文件自动生成
- 迁移依赖关系管理
- 迁移执行日志
- 数据库备份和恢复 