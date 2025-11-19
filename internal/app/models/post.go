package models

import (
	"time"

	"gorm.io/gorm"
)

// Post 文章模型
type Post struct {
	gorm.Model
	Title         string     `gorm:"size:200;not null" json:"title"`
	Slug          string     `gorm:"size:200;uniqueIndex;not null" json:"slug"`
	Content       string     `gorm:"type:longtext" json:"content"`
	Excerpt       string     `gorm:"type:text" json:"excerpt"`
	FeaturedImage string     `gorm:"size:500" json:"featured_image"`
	Status        string     `gorm:"size:20;default:'draft';index" json:"status"` // draft, published, archived
	AuthorID      uint       `gorm:"index;not null" json:"author_id"`
	CategoryID    uint       `gorm:"index" json:"category_id"`
	ViewCount     int        `gorm:"default:0" json:"view_count"`
	LikeCount     int        `gorm:"default:0" json:"like_count"`
	CommentCount  int        `gorm:"default:0" json:"comment_count"`
	PublishedAt   *time.Time `json:"published_at"`

	// SEO字段
	MetaTitle       string `gorm:"size:200" json:"meta_title"`
	MetaDescription string `gorm:"type:text" json:"meta_description"`
	MetaKeywords    string `gorm:"size:500" json:"meta_keywords"`

	// 关联
	Author   *User     `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags     []Tag     `gorm:"many2many:post_tags;" json:"tags,omitempty"`
}

// TableName 指定表名
func (Post) TableName() string {
	return "posts"
}

// BeforeCreate 创建前钩子
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.Status == "" {
		p.Status = "draft"
	}
	return nil
}

// Publish 发布文章
func (p *Post) Publish() {
	p.Status = "published"
	now := time.Now()
	p.PublishedAt = &now
}

// IsPublished 检查是否已发布
func (p *Post) IsPublished() bool {
	return p.Status == "published" && p.PublishedAt != nil
}

// Category 分类模型
type Category struct {
	gorm.Model
	Name        string `gorm:"size:100;not null" json:"name"`
	Slug        string `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description"`
	ParentID    *uint  `gorm:"index" json:"parent_id"`
	Sort        int    `gorm:"default:0" json:"sort"`
	Image       string `gorm:"size:500" json:"image"`

	// SEO字段
	MetaTitle       string `gorm:"size:200" json:"meta_title"`
	MetaDescription string `gorm:"type:text" json:"meta_description"`

	// 关联
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Posts    []Post     `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

// Tag 标签模型
type Tag struct {
	gorm.Model
	Name  string `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Slug  string `gorm:"size:50;uniqueIndex;not null" json:"slug"`
	Count int    `gorm:"default:0" json:"count"` // 使用该标签的文章数

	// 关联
	Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}

// PostTag 文章标签关联表
type PostTag struct {
	PostID uint `gorm:"primaryKey"`
	TagID  uint `gorm:"primaryKey"`
}

// TableName 指定表名
func (PostTag) TableName() string {
	return "post_tags"
}
