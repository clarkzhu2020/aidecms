package models

import (
	"time"

	"gorm.io/gorm"
)

// Media 媒体文件模型
type Media struct {
	gorm.Model
	UserID       uint      `gorm:"index" json:"user_id"`                   // 上传用户ID
	FileName     string    `gorm:"size:255;not null" json:"file_name"`     // 存储文件名
	OriginalName string    `gorm:"size:255" json:"original_name"`          // 原始文件名
	FilePath     string    `gorm:"size:500;not null" json:"file_path"`     // 文件路径
	FileURL      string    `gorm:"size:500" json:"file_url"`               // 访问URL
	FileSize     int64     `json:"file_size"`                              // 文件大小（字节）
	FileType     string    `gorm:"size:50" json:"file_type"`               // 文件类型（image/document/video等）
	MimeType     string    `gorm:"size:100" json:"mime_type"`              // MIME类型
	Extension    string    `gorm:"size:20" json:"extension"`               // 文件扩展名
	Hash         string    `gorm:"size:64;index" json:"hash"`              // 文件哈希值
	Width        int       `json:"width,omitempty"`                        // 图片宽度
	Height       int       `json:"height,omitempty"`                       // 图片高度
	Thumbnails   string    `gorm:"type:text" json:"thumbnails,omitempty"`  // 缩略图信息（JSON格式）
	Description  string    `gorm:"type:text" json:"description"`           // 描述
	Alt          string    `gorm:"size:255" json:"alt"`                    // 图片Alt文本
	Title        string    `gorm:"size:255" json:"title"`                  // 标题
	Status       string    `gorm:"size:20;default:'active'" json:"status"` // 状态：active/archived/deleted
	UploadedAt   time.Time `json:"uploaded_at"`                            // 上传时间

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Media) TableName() string {
	return "media"
}

// BeforeCreate 创建前钩子
func (m *Media) BeforeCreate(tx *gorm.DB) error {
	m.UploadedAt = time.Now()
	if m.Status == "" {
		m.Status = "active"
	}
	return nil
}

// GetFileType 根据MIME类型获取文件类型
func GetFileType(mimeType string) string {
	switch {
	case len(mimeType) >= 5 && mimeType[:5] == "image":
		return "image"
	case len(mimeType) >= 5 && mimeType[:5] == "video":
		return "video"
	case len(mimeType) >= 5 && mimeType[:5] == "audio":
		return "audio"
	case mimeType == "application/pdf":
		return "pdf"
	case mimeType == "application/zip" || mimeType == "application/x-rar-compressed":
		return "archive"
	case mimeType == "application/msword" || mimeType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return "document"
	case mimeType == "application/vnd.ms-excel" || mimeType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return "spreadsheet"
	default:
		return "other"
	}
}
