package main

import (
	"fmt"
	"log"

	"github.com/clarkzhu2020/aidecms/database/migrations"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 直接创建数据库连接
	dsn := "root:@tcp(localhost:3306)/aidecms?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 执行迁移
	if err := db.AutoMigrate(&migrations.User{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migration completed successfully")
}
