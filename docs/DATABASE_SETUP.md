# 数据库设置指南

## 环境变量配置

### 创建 .env 文件

在项目根目录创建 `.env` 文件，包含以下配置：

```bash
# 数据库配置
# 重要：DSN 必须包含 parseTime=true 参数来解决时间字段扫描问题
DSN=root:password@tcp(localhost:3306)/go_study?charset=utf8mb4&parseTime=true&loc=Local

# 环境配置
ENV=development

# JWT 配置
JWT_SECRET=your-secret-key-here
JWT_DURATION=24h
```

### DSN 参数说明

- `parseTime=true`: **必需参数**，让 MySQL 驱动自动将 DATETIME 和 TIMESTAMP 字段转换为 Go 的 time.Time 类型
- `loc=Local`: 设置时区为本地时区
- `charset=utf8mb4`: 支持完整的 UTF-8 字符集，包括 emoji

## 常见问题解决

### 1. 时间字段扫描错误

**错误信息：**
```
sql: Scan error on column index 4, name "created_at": unsupported Scan, storing driver.Value type []uint8 into type *time.Time
```

**解决方案：**
确保 DSN 连接字符串包含 `parseTime=true` 参数：

```bash
# 错误示例
DSN=root:password@tcp(localhost:3306)/go_study

# 正确示例
DSN=root:password@tcp(localhost:3306)/go_study?parseTime=true&loc=Local
```

### 2. 时区问题

如果遇到时区相关的问题，可以在 DSN 中指定时区：

```bash
# 使用 UTC 时区
DSN=root:password@tcp(localhost:3306)/go_study?parseTime=true&loc=UTC

# 使用本地时区
DSN=root:password@tcp(localhost:3306)/go_study?parseTime=true&loc=Local

# 使用特定时区
DSN=root:password@tcp(localhost:3306)/go_study?parseTime=true&loc=Asia%2FShanghai
```

### 3. 字符集问题

确保使用 `utf8mb4` 字符集以支持完整的 Unicode 字符：

```bash
DSN=root:password@tcp(localhost:3306)/go_study?charset=utf8mb4&parseTime=true&loc=Local
```

## 数据库表结构

### User 表

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(10) NOT NULL,
    email VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(18) NOT NULL,
    role VARCHAR(10) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 连接池配置

在 `db/db.go` 中已经配置了连接池参数：

- 最大空闲连接数：10
- 最大打开连接数：100
- 连接最大生命周期：1小时

## 测试连接

运行应用程序时，如果看到以下输出，说明数据库连接成功：

```
dsn root:password@tcp(localhost:3306)/go_study?charset=utf8mb4&parseTime=true&loc=Local
连接数据库成功
创建表成功
```

## 故障排除

### 1. 连接被拒绝

检查 MySQL 服务是否运行，以及端口是否正确。

### 2. 认证失败

检查用户名和密码是否正确。

### 3. 数据库不存在

确保数据库 `go_study` 已经创建：

```sql
CREATE DATABASE go_study CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 权限问题

确保用户有足够的权限：

```sql
GRANT ALL PRIVILEGES ON go_study.* TO 'root'@'localhost';
FLUSH PRIVILEGES;
``` 