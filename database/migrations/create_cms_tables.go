package migrations

import (
	"github.com/clarkgo/clarkgo/internal/app/models"
	"github.com/clarkgo/clarkgo/pkg/database"
)

// CreateCMSTables 创建CMS相关表
func CreateCMSTables(db *database.Database) error {
	// 自动迁移所有模型
	return db.DB.AutoMigrate(
		&models.Media{},
		&models.Role{},
		&models.Permission{},
		&models.Category{},
		&models.Tag{},
		&models.Post{},
		&models.RolePermission{},
		&models.UserRole{},
		&models.PostTag{},
	)
}
