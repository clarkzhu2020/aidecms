package models

import "time"

// SwaggerBase Swagger 基础模型（用于文档生成）
type SwaggerBase struct {
	ID        uint       `json:"id" example:"1"`
	CreatedAt time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" example:"null"`
}

// PostSwagger 文章Swagger模型
// @Description 文章信息
type PostSwagger struct {
	SwaggerBase
	Title           string           `json:"title" example:"文章标题"`
	Slug            string           `json:"slug" example:"article-title"`
	Content         string           `json:"content" example:"文章内容..."`
	Excerpt         string           `json:"excerpt" example:"文章摘要"`
	FeaturedImage   string           `json:"featured_image" example:"/uploads/image.jpg"`
	Status          string           `json:"status" example:"published" enums:"draft,published,archived"`
	ViewCount       int              `json:"view_count" example:"100"`
	AuthorID        uint             `json:"author_id" example:"1"`
	CategoryID      uint             `json:"category_id" example:"1"`
	PublishedAt     *time.Time       `json:"published_at" example:"2023-01-01T00:00:00Z"`
	MetaTitle       string           `json:"meta_title" example:"SEO标题"`
	MetaDescription string           `json:"meta_description" example:"SEO描述"`
	MetaKeywords    string           `json:"meta_keywords" example:"关键词1,关键词2"`
	Author          *UserSwagger     `json:"author,omitempty"`
	Category        *CategorySwagger `json:"category,omitempty"`
	Tags            []TagSwagger     `json:"tags,omitempty"`
}

// CategorySwagger 分类Swagger模型
// @Description 分类信息
type CategorySwagger struct {
	SwaggerBase
	Name            string            `json:"name" example:"技术"`
	Slug            string            `json:"slug" example:"tech"`
	Description     string            `json:"description" example:"技术相关文章"`
	ParentID        uint              `json:"parent_id" example:"0"`
	MetaTitle       string            `json:"meta_title" example:"技术分类"`
	MetaDescription string            `json:"meta_description" example:"技术相关的所有文章"`
	Children        []CategorySwagger `json:"children,omitempty"`
}

// TagSwagger 标签Swagger模型
// @Description 标签信息
type TagSwagger struct {
	SwaggerBase
	Name string `json:"name" example:"Go语言"`
	Slug string `json:"slug" example:"golang"`
}

// UserSwagger 用户Swagger模型
// @Description 用户信息
type UserSwagger struct {
	SwaggerBase
	Username  string `json:"username" example:"admin"`
	Email     string `json:"email" example:"admin@example.com"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
}

// MediaSwagger 媒体Swagger模型
// @Description 媒体文件信息
type MediaSwagger struct {
	SwaggerBase
	FileName     string `json:"file_name" example:"image-uuid.jpg"`
	OriginalName string `json:"original_name" example:"my-photo.jpg"`
	FileURL      string `json:"file_url" example:"/uploads/2023/01/01/image-uuid.jpg"`
	FileSize     int64  `json:"file_size" example:"102400"`
	FileType     string `json:"file_type" example:"image" enums:"image,document,video,audio,other"`
	MimeType     string `json:"mime_type" example:"image/jpeg"`
	Width        int    `json:"width" example:"1920"`
	Height       int    `json:"height" example:"1080"`
	Thumbnails   string `json:"thumbnails" example:"{\"small\":\"/uploads/thumb_small.jpg\"}"`
	Hash         string `json:"hash" example:"abc123def456"`
	UploaderID   uint   `json:"uploader_id" example:"1"`
}

// RoleSwagger 角色Swagger模型
// @Description 角色信息
type RoleSwagger struct {
	SwaggerBase
	Name        string              `json:"name" example:"admin"`
	DisplayName string              `json:"display_name" example:"管理员"`
	Description string              `json:"description" example:"系统管理员角色"`
	Permissions []PermissionSwagger `json:"permissions,omitempty"`
}

// PermissionSwagger 权限Swagger模型
// @Description 权限信息
type PermissionSwagger struct {
	SwaggerBase
	Name        string `json:"name" example:"post.create"`
	DisplayName string `json:"display_name" example:"创建文章"`
	Description string `json:"description" example:"允许创建文章"`
}
