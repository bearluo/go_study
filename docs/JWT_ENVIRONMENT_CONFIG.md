# JWT 环境变量配置

本文档说明如何在生产环境中配置 JWT 相关的环境变量。

## 环境变量说明

### JWT_SECRET_KEY
- **描述**: JWT 签名密钥
- **类型**: 字符串
- **默认值**: `your-access-token-secret-key-here`
- **生产环境要求**: **必须**设置一个强密钥
- **示例**: `JWT_SECRET_KEY=your-super-secret-jwt-key-here-change-in-production`

### JWT_ACCESS_TOKEN_DURATION
- **描述**: Access Token 有效期（秒）
- **类型**: 整数
- **默认值**: `900` (15分钟)
- **示例**: `JWT_ACCESS_TOKEN_DURATION=900`

### JWT_REFRESH_TOKEN_DURATION
- **描述**: Refresh Token 有效期（秒）
- **类型**: 整数
- **默认值**: `604800` (7天)
- **示例**: `JWT_REFRESH_TOKEN_DURATION=604800`

## 配置方式

### 1. 开发环境 (.env 文件)

在项目根目录创建 `.env` 文件：

```env
# JWT 配置
JWT_SECRET_KEY=your-super-secret-jwt-key-here-change-in-production
JWT_ACCESS_TOKEN_DURATION=900
JWT_REFRESH_TOKEN_DURATION=604800
```

### 2. 生产环境 (Docker)

使用 Docker 运行时，通过环境变量传递：

```bash
docker run -e JWT_SECRET_KEY=your-production-secret-key \
           -e JWT_ACCESS_TOKEN_DURATION=900 \
           -e JWT_REFRESH_TOKEN_DURATION=604800 \
           your-app
```

或者使用 `.env` 文件：

```bash
docker run --env-file .env.production your-app
```

### 3. 生产环境 (Kubernetes)

在 Kubernetes 配置中设置：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-study-app
spec:
  template:
    spec:
      containers:
      - name: app
        image: your-app:latest
        env:
        - name: JWT_SECRET_KEY
          valueFrom:
            secretKeyRef:
              name: jwt-secret
              key: secret-key
        - name: JWT_ACCESS_TOKEN_DURATION
          value: "900"
        - name: JWT_REFRESH_TOKEN_DURATION
          value: "604800"
```

### 4. 生产环境 (系统环境变量)

直接在系统中设置环境变量：

```bash
export JWT_SECRET_KEY=your-production-secret-key
export JWT_ACCESS_TOKEN_DURATION=900
export JWT_REFRESH_TOKEN_DURATION=604800
```

## 安全建议

### 1. JWT 密钥安全
- **长度**: 至少 32 个字符
- **复杂度**: 包含大小写字母、数字和特殊字符
- **随机性**: 使用密码生成器生成
- **保密性**: 不要提交到版本控制系统

### 2. 密钥轮换
- 定期更换 JWT 密钥
- 实现密钥版本管理
- 支持多密钥验证

### 3. 有效期设置
- **Access Token**: 15-30 分钟（短期）
- **Refresh Token**: 7-30 天（长期）
- 根据安全要求调整

## 示例配置

### 开发环境 (.env)
```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=go_study
DB_CHARSET=utf8mb4

# JWT 配置
JWT_SECRET_KEY=dev-secret-key-change-in-production
JWT_ACCESS_TOKEN_DURATION=900
JWT_REFRESH_TOKEN_DURATION=604800

# 应用环境
ENV=development
```

### 生产环境 (.env.production)
```env
# 数据库配置
DB_HOST=production-db-host
DB_PORT=3306
DB_USER=prod_user
DB_PASSWORD=strong-production-password
DB_NAME=go_study_prod
DB_CHARSET=utf8mb4

# JWT 配置
JWT_SECRET_KEY=your-super-secret-production-jwt-key-here
JWT_ACCESS_TOKEN_DURATION=1800
JWT_REFRESH_TOKEN_DURATION=2592000

# 应用环境
ENV=production
```

## 验证配置

启动应用后，检查日志确认 JWT 配置已正确加载：

```
JWT 配置已初始化
```

## 故障排除

### 1. 配置未生效
- 检查环境变量是否正确设置
- 确认应用重启后配置生效
- 查看应用启动日志

### 2. 密钥问题
- 确保密钥长度足够
- 检查密钥是否包含特殊字符
- 验证密钥格式正确

### 3. 时间格式问题
- 确保时间单位为秒
- 检查数值格式正确
- 验证时间范围合理

## 使用示例

### 1. 开发环境测试

创建 `.env` 文件并启动应用：

```bash
# 创建 .env 文件
cat > .env << EOF
JWT_SECRET_KEY=my-dev-secret-key-123
JWT_ACCESS_TOKEN_DURATION=300
JWT_REFRESH_TOKEN_DURATION=86400
EOF

# 启动应用
go run cmd/main.go
```

### 2. 生产环境部署

使用 Docker 部署：

```bash
# 创建生产环境配置文件
cat > .env.production << EOF
JWT_SECRET_KEY=your-super-secret-production-key-here
JWT_ACCESS_TOKEN_DURATION=1800
JWT_REFRESH_TOKEN_DURATION=2592000
EOF

# 构建并运行 Docker 容器
docker build -t go-study-app .
docker run --env-file .env.production -p 8080:8080 go-study-app
```

### 3. 验证配置生效

启动应用后，可以通过以下方式验证配置：

```bash
# 查看启动日志
docker logs <container-id>

# 或者直接运行应用查看输出
go run cmd/main.go
```

应该看到类似输出：
```
环境变量 development
JWT 配置已初始化
数据库连接已初始化
```

### 4. 动态配置测试

可以通过修改环境变量来测试不同配置：

```bash
# 测试不同的 Access Token 有效期
export JWT_ACCESS_TOKEN_DURATION=60  # 1分钟
go run cmd/main.go

# 测试不同的 Refresh Token 有效期
export JWT_REFRESH_TOKEN_DURATION=3600  # 1小时
go run cmd/main.go
``` 