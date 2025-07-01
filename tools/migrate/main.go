package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go-study/db"
	"go-study/db/migrations"
)

func main() {
	// 解析命令行参数
	var (
		action = flag.String("action", "migrate", "迁移操作: migrate, rollback, status, create")
		name   = flag.String("name", "", "迁移名称 (用于 create 操作)")
		help   = flag.Bool("help", false, "显示帮助信息")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// 执行相应的操作
	switch *action {
	case "migrate", "rollback", "status":
		// 这些操作需要数据库连接
		if err := db.InitDB(); err != nil {
			log.Fatalf("数据库连接失败: %v", err)
		}

		// 创建迁移管理器
		manager := migrations.RegisterAllMigrations(db.DB)

		switch *action {
		case "migrate":
			if err := manager.Migrate(); err != nil {
				log.Fatalf("迁移失败: %v", err)
			}
			fmt.Println("迁移完成")
		case "rollback":
			if err := manager.Rollback(); err != nil {
				log.Fatalf("回滚失败: %v", err)
			}
			fmt.Println("回滚完成")
		case "status":
			if err := manager.Status(); err != nil {
				log.Fatalf("获取状态失败: %v", err)
			}
		}
	case "create":
		if *name == "" {
			fmt.Println("错误: create 操作需要指定迁移名称")
			fmt.Println("用法: go run tools/migrate/main.go -action create -name <迁移名称>")
			os.Exit(1)
		}
		if err := createMigrationFile(*name); err != nil {
			log.Fatalf("创建迁移文件失败: %v", err)
		}
		fmt.Println("迁移文件创建成功")
	default:
		fmt.Printf("未知操作: %s\n", *action)
		showHelp()
		os.Exit(1)
	}
}

// createMigrationFile 创建迁移文件
func createMigrationFile(name string) error {
	// 生成时间戳
	timestamp := time.Now().Format("20060102_150405")
	version := timestamp

	// 生成文件名
	fileName := fmt.Sprintf("%s_%s.go", timestamp, name)
	filePath := filepath.Join("db", "migrations", fileName)

	// 生成结构体名称
	structName := fmt.Sprintf("%sMigration", toCamelCase(name))

	// 生成迁移文件内容
	content := fmt.Sprintf(`package migrations

import (
	"gorm.io/gorm"
)

// %s %s迁移
type %s struct{}

// Up 执行迁移
func (m *%s) Up(db *gorm.DB) error {
	// TODO: 在这里实现迁移逻辑
	// 示例:
	// return db.AutoMigrate(&models.YourModel{})
	return nil
}

// Down 回滚迁移
func (m *%s) Down(db *gorm.DB) error {
	// TODO: 在这里实现回滚逻辑
	// 示例:
	// return db.Migrator().DropTable(&models.YourModel{})
	return nil
}

// Version 获取版本号
func (m *%s) Version() string {
	return "%s"
}

// Name 获取迁移名称
func (m *%s) Name() string {
	return "%s"
}
`, structName, name, structName, structName, structName, structName, version, structName, name)

	// 写入文件
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	fmt.Printf("创建迁移文件: %s\n", filePath)
	fmt.Printf("结构体名称: %s\n", structName)
	fmt.Printf("版本号: %s\n", version)
	fmt.Println("请记得在 register.go 中注册新的迁移")

	return nil
}

// toCamelCase 将下划线分隔的名称转换为驼峰命名
func toCamelCase(s string) string {
	var result string
	capitalize := true

	for _, char := range s {
		if char == '_' {
			capitalize = true
		} else {
			if capitalize {
				result += string(char - 32) // 转换为大写
				capitalize = false
			} else {
				result += string(char)
			}
		}
	}

	return result
}

func showHelp() {
	fmt.Println("数据库迁移工具")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  go run tools/migrate/main.go [选项]")
	fmt.Println()
	fmt.Println("选项:")
	fmt.Println("  -action string")
	fmt.Println("        迁移操作 (默认: migrate)")
	fmt.Println("        migrate  - 执行所有未执行的迁移")
	fmt.Println("        rollback - 回滚最后一个迁移")
	fmt.Println("        status   - 显示迁移状态")
	fmt.Println("        create   - 创建新的迁移文件")
	fmt.Println("  -name string")
	fmt.Println("        迁移名称 (用于 create 操作)")
	fmt.Println("  -help")
	fmt.Println("        显示帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  go run tools/migrate/main.go -action migrate")
	fmt.Println("  go run tools/migrate/main.go -action rollback")
	fmt.Println("  go run tools/migrate/main.go -action status")
	fmt.Println("  go run tools/migrate/main.go -action create -name create_products_table")
}
