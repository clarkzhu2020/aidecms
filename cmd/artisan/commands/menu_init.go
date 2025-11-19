package commands

import (
	"fmt"

	"github.com/clarkgo/clarkgo/database/migrations"
	"github.com/clarkgo/clarkgo/pkg/database"
	"github.com/spf13/cobra"
)

// MenuInitCmd 菜单初始化命令
var MenuInitCmd = &cobra.Command{
	Use:   "menu:init",
	Short: "Initialize menu system",
	Long:  `Create menu tables and seed default menus`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing menu system...")

		db := database.GetDB()

		// 创建菜单表
		fmt.Println("Creating menu tables...")
		if err := migrations.CreateMenusTable(db); err != nil {
			fmt.Printf("Error creating menu tables: %v\n", err)
			return
		}
		fmt.Println("✓ Menu tables created successfully")

		// 添加默认菜单
		fmt.Println("\nSeeding default menus...")
		if err := migrations.SeedDefaultMenus(db); err != nil {
			fmt.Printf("Error seeding default menus: %v\n", err)
			return
		}
		fmt.Println("✓ Default menus seeded successfully")

		fmt.Println("\n✓ Menu system initialization completed successfully!")
		fmt.Println("\nDefault menus created:")
		fmt.Println("  Header: 首页, 博客, 关于, 联系我们")
		fmt.Println("  Footer: 隐私政策, 服务条款, 站点地图")
	},
}
