# Go项目Makefile

.PHONY: help build run test clean migrate migrate-status migrate-rollback test-db

# 默认目标
help:
	@echo "可用的命令:"
	@echo "  build          - 构建项目"
	@echo "  run            - 运行项目"
	@echo "  test           - 运行测试"
	@echo "  clean          - 清理构建文件"
	@echo "  migrate        - 执行数据库迁移"
	@echo "  migrate-status - 查看迁移状态"
	@echo "  migrate-rollback - 回滚最后一个迁移"
	@echo "  test-db        - 测试数据库连接和时间字段"

# 构建项目
build:
	go build -o bin/app cmd/main.go

# 运行项目
run:
	go run cmd/main.go

# 运行测试
test:
	go test ./...

# 清理构建文件
clean:
	rm -rf bin/
	go clean

# 数据库迁移相关命令
migrate:
	go run tools/migrate/main.go -action migrate

migrate-status:
	go run tools/migrate/main.go -action status

migrate-rollback:
	go run tools/migrate/main.go -action rollback

# 测试数据库连接
test-db:
	go run tools/test_db.go

# 开发相关命令
dev:
	air

# 安装依赖
deps:
	go mod download
	go mod tidy

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run

# 生成文档
docs:
	godoc -http=:6060 