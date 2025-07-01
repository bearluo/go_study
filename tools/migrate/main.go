package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go-study/db"
	"go-study/db/migrations"
)

func main() {
	// 解析命令行参数
	var (
		action = flag.String("action", "migrate", "迁移操作: migrate, rollback, status")
		help   = flag.Bool("help", false, "显示帮助信息")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// 初始化数据库连接
	if err := db.InitDB(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 创建迁移管理器
	manager := migrations.RegisterAllMigrations(db.DB)

	// 执行相应的操作
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
	default:
		fmt.Printf("未知操作: %s\n", *action)
		showHelp()
		os.Exit(1)
	}
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
	fmt.Println("  -help")
	fmt.Println("        显示帮助信息")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  go run tools/migrate/main.go -action migrate")
	fmt.Println("  go run tools/migrate/main.go -action rollback")
	fmt.Println("  go run tools/migrate/main.go -action status")
}
