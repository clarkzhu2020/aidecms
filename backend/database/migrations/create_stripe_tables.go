package migrations

import (
	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"github.com/clarkzhu2020/aidecms/pkg/database"
)

// CreateStripeTables 创建Stripe相关表
func CreateStripeTables(db *database.Database) error {
	return db.DB.AutoMigrate(
		&models.StripePayment{},
		&models.StripeRefund{},
		&models.StripeWebhook{},
	)
}
