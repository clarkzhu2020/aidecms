package models

import (
	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	gorm.Model
	PostID      uint      `gorm:"index;not null" json:"post_id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	ParentID    uint      `gorm:"index;default:0" json:"parent_id"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	AuthorName  string    `gorm:"size:100" json:"author_name"`
	AuthorEmail string    `gorm:"size:100" json:"author_email"`
	AuthorURL   string    `gorm:"size:200" json:"author_url"`
	AuthorIP    string    `gorm:"size:45" json:"author_ip"`
	UserAgent   string    `gorm:"size:500" json:"user_agent"`
	Status      string    `gorm:"size:20;default:pending" json:"status"` // pending, approved, spam, trash
	Rating      int       `gorm:"default:0" json:"rating"`               // 1-5 星级评分（可选）
	IsAnonymous bool      `gorm:"default:false" json:"is_anonymous"`
	Children    []Comment `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Parent      *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	User        *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post        *Post     `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}

// IsApproved 检查评论是否已批准
func (c *Comment) IsApproved() bool {
	return c.Status == "approved"
}

// IsSpam 检查评论是否为垃圾评论
func (c *Comment) IsSpam() bool {
	return c.Status == "spam"
}

// Approve 批准评论
func (c *Comment) Approve() {
	c.Status = "approved"
}

// MarkAsSpam 标记为垃圾评论
func (c *Comment) MarkAsSpam() {
	c.Status = "spam"
}

// CommentSwagger 评论Swagger模型
// @Description 评论信息
type CommentSwagger struct {
	SwaggerBase
	PostID      uint             `json:"post_id" example:"1"`
	UserID      uint             `json:"user_id" example:"1"`
	ParentID    uint             `json:"parent_id" example:"0"`
	Content     string           `json:"content" example:"这是一条评论内容"`
	AuthorName  string           `json:"author_name" example:"张三"`
	AuthorEmail string           `json:"author_email" example:"zhangsan@example.com"`
	AuthorURL   string           `json:"author_url" example:"https://zhangsan.com"`
	Status      string           `json:"status" example:"approved" enums:"pending,approved,spam,trash"`
	Rating      int              `json:"rating" example:"5"`
	IsAnonymous bool             `json:"is_anonymous" example:"false"`
	Children    []CommentSwagger `json:"children,omitempty"`
	User        *UserSwagger     `json:"user,omitempty"`
}
