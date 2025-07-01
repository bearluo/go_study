# go + echo + grom 项目

my-go-project/
|-- cmd/                      # 应用程序的启动点
|   |-- main.go               # 主入口文件
|-- config/                   # 配置文件
|   |-- database.go           # 数据库配置
|   |-- server.go             # 服务器配置
|   |-- env.example           # 环境变量示例文件
|-- db/                       # 数据库相关文件
|   |-- migrations/           # 数据库迁移文件
|   |   |-- 20230801_create_users_table.go
|   |   |-- 20230802_add_roles_to_users_table.go
|   |-- models/               # 数据库模型
|   |   |-- user.go           # 用户模型
|   |   |-- role.go           # 角色模型
|   |-- repositories/         # 数据访问层 (DAL)
|   |   |-- user.go           # 用户数据访问层
|   |   |-- role.go           # 角色数据访问层
|-- handlers/                 # HTTP 请求处理器
|   |-- user.go               # 用户请求处理器
|   |-- role.go               # 角色请求处理器
|-- middleware/               # 中间件
|   |-- auth.go               # 认证中间件
|   |-- logging.go            # 日志记录中间件
|-- routers/                   # 路由配置
|   |-- user.go               # 用户路由
|   |-- role.go               # 角色路由
|-- services/                 # 业务逻辑层
|   |-- user.go               # 用户服务
|   |-- role.go               # 角色服务
|-- utils/                    # 辅助函数和工具
|   |-- jwt.go                # JWT 工具
|   |-- validators.go         # 验证工具
|-- internal/                 # 内部模块或包
|   |-- ...
|-- static/                   # 静态文件
|-- templates/                # 模板文件
|-- go.mod                    # Go 模块元数据
|-- go.sum                    # Go 模块校验和
|-- .gitignore                # Git 忽略文件列表
|-- .env                      # 环境变量文件
|-- .air.toml                 # air热重启配置
|-- .dockerignore             # Docker 忽略文件列表
|-- Dockerfile                # Docker 镜像构建文件
|-- Makefile                  # 自动化任务管理文件
|-- README.md                 # 项目说明文件
|-- LICENSE                   # 许可证文件