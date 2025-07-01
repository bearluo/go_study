package main

import (
	"fmt"
	"log"

	"go-study/db"
	"go-study/db/models"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Printf("警告：无法加载 .env 文件: %v", err)
	}

	// 初始化数据库
	if err := db.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 测试创建用户
	testUser := &models.User{
		Name:     "测试用户",
		Email:    "test@example.com",
		Password: "123456",
		Role:     "user",
	}

	if err := db.DB.Create(testUser).Error; err != nil {
		log.Fatalf("创建用户失败: %v", err)
	}

	fmt.Printf("用户创建成功，ID: %d, 创建时间: %v\n", testUser.ID, testUser.CreatedAt)

	// 测试查询用户
	var retrievedUser models.User
	if err := db.DB.First(&retrievedUser, testUser.ID).Error; err != nil {
		log.Fatalf("查询用户失败: %v", err)
	}

	fmt.Printf("用户查询成功，ID: %d, 创建时间: %v, 更新时间: %v\n",
		retrievedUser.ID, retrievedUser.CreatedAt, retrievedUser.UpdatedAt)

	// 测试更新用户
	retrievedUser.Name = "更新后的用户名"
	if err := db.DB.Save(&retrievedUser).Error; err != nil {
		log.Fatalf("更新用户失败: %v", err)
	}

	fmt.Printf("用户更新成功，更新时间: %v\n", retrievedUser.UpdatedAt)

	// 测试查询所有用户
	var allUsers []models.User
	if err := db.DB.Find(&allUsers).Error; err != nil {
		log.Fatalf("查询所有用户失败: %v", err)
	}

	fmt.Printf("查询到 %d 个用户:\n", len(allUsers))
	for _, user := range allUsers {
		fmt.Printf("  - ID: %d, 姓名: %s, 邮箱: %s, 创建时间: %v\n",
			user.ID, user.Name, user.Email, user.CreatedAt)
	}

	fmt.Println("数据库时间字段测试完成！")
}
