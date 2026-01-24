package models

import (
	"gorm.io/gorm"
)

// Menu 菜单模型
type Menu struct {
	gorm.Model
	Name        string `gorm:"size:100;not null" json:"name"`
	Title       string `gorm:"size:200" json:"title"`
	URL         string `gorm:"size:500" json:"url"`
	Icon        string `gorm:"size:100" json:"icon"`
	Target      string `gorm:"size:20;default:_self" json:"target"` // _self, _blank
	Position    string `gorm:"size:50;not null" json:"position"`    // header, footer, sidebar, mobile
	ParentID    uint   `gorm:"default:0" json:"parent_id"`
	Sort        int    `gorm:"default:0" json:"sort"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	CSSClass    string `gorm:"size:100" json:"css_class"`
	Description string `gorm:"size:500" json:"description"`
	Children    []Menu `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Parent      *Menu  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}

// MenuSwagger 菜单Swagger模型
// @Description 菜单信息
type MenuSwagger struct {
	SwaggerBase
	Name        string        `json:"name" example:"主菜单"`
	Title       string        `json:"title" example:"首页"`
	URL         string        `json:"url" example:"/home"`
	Icon        string        `json:"icon" example:"icon-home"`
	Target      string        `json:"target" example:"_self" enums:"_self,_blank"`
	Position    string        `json:"position" example:"header" enums:"header,footer,sidebar,mobile"`
	ParentID    uint          `json:"parent_id" example:"0"`
	Sort        int           `json:"sort" example:"1"`
	IsActive    bool          `json:"is_active" example:"true"`
	CSSClass    string        `json:"css_class" example:"nav-item"`
	Description string        `json:"description" example:"菜单描述"`
	Children    []MenuSwagger `json:"children,omitempty"`
}
