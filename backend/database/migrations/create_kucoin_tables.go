package migrations

import (
	"aidecms/internal/app/models"
	"gorm.io/gorm"
)

// CreateKuCoinTables 创建KuCoin相关表
func CreateKuCoinTables(db *gorm.DB) error {
	// 创建KuCoin订单表
	if err := db.AutoMigrate(&models.KuCoinOrder{}); err != nil {
		return err
	}

	// 创建KuCoin成交记录表
	if err := db.AutoMigrate(&models.KuCoinTrade{}); err != nil {
		return err
	}

	// 创建KuCoin账户表
	if err := db.AutoMigrate(&models.KuCoinAccount{}); err != nil {
		return err
	}

	// 创建KuCoin余额快照表
	if err := db.AutoMigrate(&models.KuCoinBalanceSnapshot{}); err != nil {
		return err
	}

	return nil
}
