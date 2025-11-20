package migrations

import (
	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"gorm.io/gorm"
)

// MigrateUserModel 迁移用户模型
func MigrateUserModel(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}
