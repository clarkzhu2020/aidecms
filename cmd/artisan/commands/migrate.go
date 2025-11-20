package commands

import (
	"fmt"
	"log"

	"github.com/chenyusolar/aidecms/config"
	"github.com/chenyusolar/aidecms/database/migrations"
)

func Migrate(args []string) {
	// 使用全局DB实例
	db := config.DB

	// 执行迁移
	if err := db.AutoMigrate(&migrations.User{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migration completed successfully")
}
