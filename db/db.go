package db

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	// 获取.env的DSN变量
	dsn := os.Getenv("DSN")
	fmt.Println("dsn", dsn)

	// 调用 Open 方法，传入驱动名和连接字符串
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	// 检查是否有错误
	if err != nil {
		fmt.Println("连接数据库失败：", err)
		return err
	}

	// 获取底层的 sql.DB 对象并设置连接池参数
	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println("获取底层数据库连接失败：", err)
		return err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 打印成功信息
	fmt.Println("连接数据库成功", DB)
	return nil
}
