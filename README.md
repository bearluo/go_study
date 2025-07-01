# Go + Echo + GORM 项目

一个基于 Go 语言、Echo 框架和 GORM 的现代化 Web 应用程序模板，采用清晰的分层架构设计。

## 📋 目录

- [项目简介](#项目简介)
- [功能特性](#功能特性)
- [技术栈](#技术栈)
- [项目结构](#项目结构)
- [快速开始](#快速开始)
- [开发指南](#开发指南)

## 🚀 项目简介

这是一个使用 Go 语言开发的 Web 应用程序，采用 Echo 作为 Web 框架，GORM 作为 ORM 框架。项目遵循清晰的分层架构，包含用户认证、角色管理等功能模块。

## ✨ 功能特性

- 🔐 **用户认证系统** - JWT 令牌认证
- 👥 **用户管理** - 用户注册、登录、信息管理
- 🎭 **角色管理** - 基于角色的权限控制
- 🗄️ **数据库迁移** - 自动化的数据库结构管理
- 🔒 **中间件支持** - 认证、日志记录等中间件
- 🔄 **热重载** - 开发环境热重载支持

## 🛠️ 技术栈

- **后端框架**: [Echo](https://echo.labstack.com/) - 高性能、可扩展的 Go Web 框架
- **ORM**: [GORM](https://gorm.io/) - Go 语言的 ORM 库
- **数据库**: MySQL/PostgreSQL/SQLite
- **认证**: JWT (JSON Web Tokens)
- **配置管理**: 环境变量 + 配置文件
- **开发工具**: Air (热重载)

## 📁 项目结构

```
go_study/
├── cmd/                          # 应用程序入口点
│   └── main.go                   # 主程序入口
├── config/                       # 配置文件目录
├── db/                           # 数据库相关
│   ├── migrations/               # 数据库迁移文件
│   ├── models/                   # 数据模型
│   └── repositories/             # 数据访问层
├── handlers/                     # HTTP 请求处理器
├── middleware/                   # 中间件
├── routers/                      # 路由配置
├── services/                     # 业务逻辑层
├── utils/                        # 工具函数
├── docs/                         # 项目文档
├── tools/                        # 开发工具
├── static/                       # 静态文件
├── templates/                    # 模板文件
├── tmp/                         # 临时文件
├── go.mod                       # Go 模块文件
├── go.sum                       # Go 依赖校验文件
├── Makefile                     # 构建脚本
└── README.md                    # 项目说明文档
```

## 🚀 快速开始

### 环境要求

- Go 1.19+
- MySQL 8.0+ / PostgreSQL 12+ / SQLite 3
- Git

### 安装步骤

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd go_study
   ```

2. **安装依赖**
   ```bash
   go mod download
   ```

3. **配置环境变量**
   ```bash
   cp config/env.example .env
   # 编辑 .env 文件，配置数据库连接等信息
   ```

4. **运行数据库迁移**
   ```bash
   go run tools/migrate/main.go
   ```

5. **启动应用**
   ```bash
   go run cmd/main.go
   ```

### 使用 Makefile

```bash
# 安装依赖
make deps

# 运行测试
make test

# 构建应用
make build

# 运行应用
make run

# 清理构建文件
make clean
```

## 🛠️ 开发指南

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 编写单元测试，测试覆盖率 > 80%

### 添加新功能

1. 在 `models/` 中定义数据模型
2. 在 `repositories/` 中实现数据访问层
3. 在 `services/` 中实现业务逻辑
4. 在 `handlers/` 中实现 HTTP 处理器
5. 在 `routers/` 中配置路由

### 数据库迁移

```bash
# 创建新的迁移文件
go run tools/migrate/main.go -action create -name migration_name

# 运行迁移
go run tools/migrate/main.go -action migrate

# 回滚迁移
go run tools/migrate/main.go -action rollback

# 查看迁移状态
go run tools/migrate/main.go -action status

# 显示帮助信息
go run tools/migrate/main.go -help
```