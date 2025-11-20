package migrations

import (
	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"gorm.io/gorm"
)

// CreateMenusTable 创建菜单表
func CreateMenusTable(db *gorm.DB) error {
	return db.AutoMigrate(&models.Menu{})
}

// SeedDefaultMenus 添加默认菜单
func SeedDefaultMenus(db *gorm.DB) error {
	// 检查是否已有菜单
	var count int64
	if err := db.Model(&models.Menu{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil // 已有菜单，跳过
	}

	// Header 菜单
	headerMenus := []models.Menu{
		{
			Name:     "home",
			Title:    "首页",
			URL:      "/",
			Icon:     "icon-home",
			Target:   "_self",
			Position: "header",
			Sort:     1,
			IsActive: true,
		},
		{
			Name:     "blog",
			Title:    "博客",
			URL:      "/blog",
			Icon:     "icon-blog",
			Target:   "_self",
			Position: "header",
			Sort:     2,
			IsActive: true,
		},
		{
			Name:     "about",
			Title:    "关于",
			URL:      "/about",
			Icon:     "icon-info",
			Target:   "_self",
			Position: "header",
			Sort:     3,
			IsActive: true,
		},
		{
			Name:     "contact",
			Title:    "联系我们",
			URL:      "/contact",
			Icon:     "icon-mail",
			Target:   "_self",
			Position: "header",
			Sort:     4,
			IsActive: true,
		},
	}

	// Footer 菜单
	footerMenus := []models.Menu{
		{
			Name:     "privacy",
			Title:    "隐私政策",
			URL:      "/privacy",
			Target:   "_self",
			Position: "footer",
			Sort:     1,
			IsActive: true,
		},
		{
			Name:     "terms",
			Title:    "服务条款",
			URL:      "/terms",
			Target:   "_self",
			Position: "footer",
			Sort:     2,
			IsActive: true,
		},
		{
			Name:     "sitemap",
			Title:    "站点地图",
			URL:      "/sitemap.xml",
			Target:   "_blank",
			Position: "footer",
			Sort:     3,
			IsActive: true,
		},
	}

	// 插入所有菜单
	allMenus := append(headerMenus, footerMenus...)
	for _, menu := range allMenus {
		if err := db.Create(&menu).Error; err != nil {
			return err
		}
	}

	return nil
}
