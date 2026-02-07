package main

import (
	"log"

	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"github.com/clarkzhu2020/aidecms/pkg/database"
	"gorm.io/gorm"
)

// CreateMoonPayTables 创建MoonPay相关数据表
func CreateMoonPayTables(db *gorm.DB) error {
	log.Println("Creating MoonPay tables...")

	// 创建moonpay_transactions表
	if err := db.AutoMigrate(&models.MoonPayTransaction{}); err != nil {
		return err
	}
	log.Println("Table moonpay_transactions created/updated")

	// 创建moonpay_webhooks表
	if err := db.AutoMigrate(&models.MoonPayWebhook{}); err != nil {
		return err
	}
	log.Println("Table moonpay_webhooks created/updated")

	log.Println("MoonPay tables created successfully")
	return nil
}

func main() {
	// 初始化数据库连接
	dbConfig := &database.Config{
		Driver:   "sqlite",
		Database: "database/data.db",
		Debug:    true,
	}

	db := database.NewDatabase(dbConfig)
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 获取GORM实例
	gormDB := db.GetDB()

	// 创建MoonPay表
	if err := CreateMoonPayTables(gormDB); err != nil {
		log.Fatalf("Failed to create MoonPay tables: %v", err)
	}

	log.Println("Migration completed successfully")
}
