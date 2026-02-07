package migrations

import (
	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"github.com/clarkzhu2020/aidecms/pkg/database"
)

// CreatePaymentTables 创建支付相关表
func CreatePaymentTables(db *database.Database) error {
	return db.DB.AutoMigrate(
		&models.Payment{},
		&models.PaymentRefund{},
		&models.PaymentWebhook{},
	)
}
