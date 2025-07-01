# 用户注册功能

本文档介绍了用户注册功能的实现和使用方法。

## 功能特性

### 1. 支持的请求格式
- **JSON 格式**：`Content-Type: application/json`
- **表单格式**：`Content-Type: application/x-www-form-urlencoded`

### 2. 数据验证
- 用户名：2-10个字符，只能包含字母和数字
- 邮箱：有效的邮箱格式，最大50个字符
- 密码：6-18个字符，支持字母、数字和特殊字符

### 3. 安全特性
- 密码使用 bcrypt 加密存储
- 用户名和邮箱唯一性检查
- 自动设置默认角色（user）

### 4. 响应格式
- 返回完整的用户信息
- 包含 Access Token 和 Refresh Token
- 提供令牌过期时间

## API 接口

### 新版本注册接口

**POST** `/api/auth/register`

#### 请求参数

**JSON 格式：**
```json
{
  "name": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

**表单格式：**
```
name=testuser&email=test@example.com&password=password123
```

#### 响应示例

**成功响应：**
```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "user": {
      "id": 1,
      "name": "testuser",
      "email": "test@example.com",
      "role": "user",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "a1b2c3d4e5f6...",
    "expires_in": 900
  }
}
```

**错误响应：**
```json
{
  "code": 2002,
  "message": "邮箱已存在",
  "error": "邮箱已存在"
}
```

### 兼容版本注册接口

**POST** `/user/create`

保持向后兼容，功能与旧版本相同。

## 使用示例

### 1. 使用 curl 进行注册

**JSON 格式：**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

**表单格式：**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "name=testuser&email=test@example.com&password=password123"
```

### 2. 使用 JavaScript 进行注册

```javascript
// JSON 格式
async function registerUser() {
  const response = await fetch('/api/auth/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      name: 'testuser',
      email: 'test@example.com',
      password: 'password123'
    })
  });
  
  const result = await response.json();
  console.log(result);
}

// 表单格式
async function registerUserForm() {
  const formData = new FormData();
  formData.append('name', 'testuser');
  formData.append('email', 'test@example.com');
  formData.append('password', 'password123');
  
  const response = await fetch('/api/auth/register', {
    method: 'POST',
    body: formData
  });
  
  const result = await response.json();
  console.log(result);
}
```

### 3. 使用 Python 进行注册

```python
import requests

# JSON 格式
def register_user_json():
    url = "http://localhost:8080/api/auth/register"
    data = {
        "name": "testuser",
        "email": "test@example.com",
        "password": "password123"
    }
    
    response = requests.post(url, json=data)
    return response.json()

# 表单格式
def register_user_form():
    url = "http://localhost:8080/api/auth/register"
    data = {
        "name": "testuser",
        "email": "test@example.com",
        "password": "password123"
    }
    
    response = requests.post(url, data=data)
    return response.json()
```

## 错误处理

### 常见错误码

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 1001 | 参数错误 | 检查请求参数格式 |
| 1002 | 验证错误 | 检查输入数据是否符合要求 |
| 2002 | 用户已存在 | 使用不同的用户名或邮箱 |

### 验证规则

#### 用户名验证
- 不能为空
- 长度：2-10个字符
- 格式：只能包含字母和数字
- 唯一性：不能与现有用户名重复

#### 邮箱验证
- 不能为空
- 格式：有效的邮箱格式
- 长度：最大50个字符
- 唯一性：不能与现有邮箱重复

#### 密码验证
- 不能为空
- 长度：6-18个字符
- 格式：支持字母、数字和特殊字符

## 安全注意事项

1. **密码安全**
   - 密码在传输前应进行适当的安全处理
   - 建议使用 HTTPS 进行传输
   - 密码在数据库中已加密存储

2. **令牌安全**
   - Access Token 有效期较短（15分钟）
   - Refresh Token 用于刷新 Access Token
   - 请妥善保管 Refresh Token

3. **输入验证**
   - 所有输入都经过严格验证
   - 防止 SQL 注入和 XSS 攻击
   - 用户名和邮箱唯一性检查

## 数据库结构

用户表结构：
```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(10) NOT NULL,
    email VARCHAR(20) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(10) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 相关文件

- `handlers/user_handle.go` - 用户注册处理器
- `services/user_service.go` - 用户注册服务
- `db/repositories/user_repository.go` - 用户数据访问层
- `db/models/user.go` - 用户模型
- `utils/validators.go` - 数据验证工具
- `utils/constants.go` - 验证常量定义 