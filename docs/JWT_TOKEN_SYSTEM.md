# JWT 令牌系统设计

## 概述

本系统实现了基于 JWT 的双令牌认证机制，包括：

- **Access Token（访问令牌）**：短期有效（默认15分钟），用于API访问认证
- **Refresh Token（刷新令牌）**：长期有效（默认7天），存储在数据库中，用于刷新 Access Token

## 系统架构

### 1. 令牌类型

#### Access Token
- **类型**：JWT（JSON Web Token）
- **有效期**：15分钟（可配置）
- **存储位置**：客户端（内存或安全存储）
- **用途**：API 访问认证
- **特点**：
  - 包含用户信息（ID、用户名、邮箱、角色）
  - 自动过期，提高安全性
  - 无需数据库查询即可验证

#### Refresh Token
- **类型**：随机字符串
- **有效期**：7天（可配置）
- **存储位置**：数据库（refresh_tokens 表）
- **用途**：刷新 Access Token
- **特点**：
  - 可撤销（支持登出功能）
  - 支持多设备登录
  - 可追踪和管理

### 2. 数据库设计

#### refresh_tokens 表
```sql
CREATE TABLE refresh_tokens (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_token (token)
);
```

### 3. 核心组件

#### JWT 工具类 (`utils/jwt.go`)
- `GenerateTokenPair()`: 生成令牌对
- `ValidateAccessToken()`: 验证 Access Token
- `RefreshAccessTokenWithUserInfo()`: 使用 Refresh Token 刷新 Access Token
- `RevokeRefreshToken()`: 撤销 Refresh Token

#### 认证服务 (`services/auth_service.go`)
- `Login()`: 用户登录，返回令牌对
- `RefreshToken()`: 刷新 Access Token
- `Logout()`: 用户登出，撤销 Refresh Token
- `LogoutAll()`: 撤销用户的所有令牌

#### 认证中间件 (`middleware/auth_middleware.go`)
- `RequireAuth()`: 要求认证的中间件
- `RequireRole()`: 要求特定角色的中间件
- `OptionalAuth()`: 可选的认证中间件

## 使用流程

### 1. 用户登录
```http
POST /auth/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123"
}
```

**响应**：
```json
{
    "code": 0,
    "message": "登录成功",
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "refresh_token": "a1b2c3d4e5f6...",
        "expires_in": 900,
        "user": {
            "id": 1,
            "name": "张三",
            "email": "user@example.com",
            "role": "user"
        }
    }
}
```

### 2. API 访问
```http
GET /api/protected-resource
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 3. 令牌刷新
当 Access Token 过期时，使用 Refresh Token 获取新的令牌对：

```http
POST /auth/refresh
Content-Type: application/json

{
    "refresh_token": "a1b2c3d4e5f6..."
}
```

**响应**：
```json
{
    "code": 0,
    "message": "令牌刷新成功",
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "refresh_token": "f6e5d4c3b2a1...",
        "expires_in": 900,
        "user": {
            "id": 1,
            "name": "张三",
            "email": "user@example.com",
            "role": "user"
        }
    }
}
```

### 4. 用户登出
```http
POST /auth/logout
Content-Type: application/json

{
    "refresh_token": "a1b2c3d4e5f6..."
}
```

## 安全特性

### 1. 令牌安全
- Access Token 短期有效，减少被盗用的风险
- Refresh Token 存储在数据库中，支持撤销
- 每次刷新都会生成新的 Refresh Token，旧令牌自动失效

### 2. 撤销机制
- 支持单个令牌撤销（登出）
- 支持批量令牌撤销（撤销用户所有令牌）
- 支持管理员强制撤销

### 3. 清理机制
- 自动清理过期的 Refresh Token
- 自动清理已撤销的 Refresh Token
- 定期清理无效令牌

## 配置选项

### JWT 配置 (`utils/jwt.go`)
```go
var DefaultJWTConfig = &JWTConfig{
    AccessTokenSecret:    "your-access-token-secret-key-here",
    AccessTokenDuration:  15 * time.Minute,
    RefreshTokenDuration: 7 * 24 * time.Hour,
}
```

### 环境变量配置
建议在生产环境中使用环境变量：

```bash
JWT_ACCESS_TOKEN_SECRET=your-secure-secret-key
JWT_ACCESS_TOKEN_DURATION=15m
JWT_REFRESH_TOKEN_DURATION=168h
```

## 最佳实践

### 1. 客户端实现
- 将 Access Token 存储在内存中（如 JavaScript 变量）
- 将 Refresh Token 存储在安全的存储中（如 HttpOnly Cookie）
- 在 Access Token 过期前主动刷新
- 实现自动重试机制

### 2. 服务器端实现
- 使用强密钥签名 JWT
- 定期轮换密钥
- 实现令牌黑名单机制（可选）
- 监控异常登录行为

### 3. 安全建议
- 使用 HTTPS 传输
- 设置适当的 Cookie 安全属性
- 实现速率限制
- 记录认证日志

## 错误处理

### 常见错误码
- `401 Unauthorized`: 认证失败
- `403 Forbidden`: 权限不足
- `400 Bad Request`: 请求参数错误

### 错误响应格式
```json
{
    "code": 2004,
    "message": "令牌错误",
    "error": "认证令牌无效"
}
```

## 扩展功能

### 1. 多设备支持
- 每个设备可以有不同的 Refresh Token
- 支持设备管理和监控
- 支持远程登出特定设备

### 2. 令牌轮换
- 支持定期强制刷新所有令牌
- 支持安全事件后的令牌轮换
- 支持用户密码修改后的令牌轮换

### 3. 审计日志
- 记录令牌生成和刷新
- 记录登录和登出事件
- 记录异常访问行为 