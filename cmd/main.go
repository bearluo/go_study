package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	db "go-study/db"
	"go-study/db/repositories"
	"go-study/middleware"
	routers "go-study/routers"
	"go-study/services"
	"go-study/utils"
)

func main() {
	// 加载环境变量
	loadEnv()
	// 初始化 JWT 配置
	initJWTConfig()
	// 初始化数据库连接
	initDataBases()
	// 启动 HTTP 服务器
	startServer()
}

func loadEnv() {
	env := os.Getenv("ENV")
	fmt.Println("环境变量", env)
	if env != "production" { // 生产环境是用docker运行的，会用--env-file参数指定.env文件，不需要手动加载
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}

func initJWTConfig() {
	utils.InitJWTConfig()
	fmt.Println("JWT 配置已初始化")
}

func initDataBases() {
	db.InitDB()
}

func startServer() {
	// 初始化 Echo 实例
	e := echo.New()

	// 设置自定义验证器
	utils.SetupEchoValidator(e)

	// 创建 Repository Manager
	repoManager := repositories.NewRepositoryManager(db.DB)

	// 创建服务层实例
	serviceManager := services.NewServiceManager(repoManager)

	// 创建中间件管理器
	middlewareManager := middleware.NewMiddlewareManager(serviceManager)

	// 设置路由
	routers.SetupRoutes(e, repoManager, serviceManager, middlewareManager)

	// 启动服务器
	e.Logger.Fatal(e.Start(":8080"))
}
