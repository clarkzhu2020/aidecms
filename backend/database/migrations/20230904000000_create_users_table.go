package migrations

import (
	"github.com/clarkzhu2020/aidecms/pkg/database"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id       uint   `gorm:"primaryKey;autoIncrement;not null;index"`
	Name     string `gorm:"size:255;not null"`
	Email    string `gorm:"size:255;uniqueIndex"`
	Password string `gorm:"size:255"`
}

func CreateUsersTable(db *database.Database) error {
	return db.DB.AutoMigrate(&User{})
}
