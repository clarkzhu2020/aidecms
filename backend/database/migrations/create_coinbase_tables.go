package main

import (
	"log"

	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"github.com/clarkzhu2020/aidecms/pkg/database"
	"gorm.io/gorm"
)

// CreateCoinbaseTables 创建Coinbase相关数据表
func CreateCoinbaseTables(db *gorm.DB) error {
	log.Println("Creating Coinbase tables...")

	// 创建coinbase_payment_links表
	if err := db.AutoMigrate(&models.CoinbasePaymentLink{}); err != nil {
		return err
	}
	log.Println("Table coinbase_payment_links created/updated")

	// 创建coinbase_orders表
	if err := db.AutoMigrate(&models.CoinbaseOrder{}); err != nil {
		return err
	}
	log.Println("Table coinbase_orders created/updated")

	// 创建coinbase_webhooks表
	if err := db.AutoMigrate(&models.CoinbaseWebhook{}); err != nil {
		return err
	}
	log.Println("Table coinbase_webhooks created/updated")

	log.Println("Coinbase tables created successfully")
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

	// 创建Coinbase表
	if err := CreateCoinbaseTables(gormDB); err != nil {
		log.Fatalf("Failed to create Coinbase tables: %v", err)
	}

	log.Println("Migration completed successfully")
}
