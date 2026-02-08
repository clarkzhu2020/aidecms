package models

import (
	"time"

	"gorm.io/gorm"
)

// CSSSession 客服会话
type CSSSession struct {
	ID           string    `gorm:"type:varchar(36);primaryKey"`
	UserID       *uint64   `gorm:"type:bigint unsigned;index"`
	Channel      string    `gorm:"type:varchar(50);not null;index"`
	ChannelID    string    `gorm:"type:varchar(100);index"`
	Status       string    `gorm:"type:varchar(20);default:'active';index"`
	AssignedTo   *uint64   `gorm:"type:bigint unsigned;index"`
	Confidence   float32   `gorm:"type:float;default:0.8"`
	MessageCount int       `gorm:"type:int;default:0"`
	LastActiveAt time.Time `gorm:"type:datetime;index"`
	TransferReason string    `gorm:"type:varchar(255)"`
	Metadata     string    `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"type:datetime;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"type:datetime;autoUpdateTime"`
}

// TableName 指定表名
func (CSSSession) TableName() string {
	return "css_sessions"
}

// CSSMessage 客服消息
type CSSMessage struct {
	ID           string  `gorm:"type:varchar(36);primaryKey"`
	SessionID     string  `gorm:"type:varchar(36);not null;index"`
	Role         string  `gorm:"type:varchar(20);not null;index"`
	Content      string  `gorm:"type:text;not null"`
	Tokens       int     `gorm:"type:int;default:0"`
	Duration     int     `gorm:"type:int"`
	Confidence   float32 `gorm:"type:float"`
	DocumentRefs string  `gorm:"type:text"`
	Metadata     string  `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"type:datetime;autoCreateTime;index"`
}

// TableName 指定表名
func (CSSMessage) TableName() string {
	return "css_messages"
}

// KBDocument 知识库文档
type KBDocument struct {
	ID          string  `gorm:"type:varchar(36);primaryKey"`
	Title       string  `gorm:"type:varchar(255);not null"`
	Category    string  `gorm:"type:varchar(100);index"`
	Tags        string  `gorm:"type:text"`
	Content     string  `gorm:"type:text;not null"`
	FileURL     string  `gorm:"type:varchar(500)"`
	FileType    string  `gorm:"type:varchar(50)"`
	FileSize    int64   `gorm:"type:bigint"`
	ChunkCount  int     `gorm:"type:int;default:0"`
	Status      string  `gorm:"type:varchar(20);default:'processing';index"`
	Version     int     `gorm:"type:int;default:1"`
	CreatedBy   uint64  `gorm:"type:bigint unsigned;index"`
	CreatedAt   time.Time `gorm:"type:datetime;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"type:datetime;autoUpdateTime"`
}

// TableName 指定表名
func (KBDocument) TableName() string {
	return "kb_documents"
}

// KBChunk 知识库文档分块
type KBChunk struct {
	ID         string    `gorm:"type:varchar(36);primaryKey"`
	DocumentID string    `gorm:"type:varchar(36);not null;index"`
	ChunkOrder int       `gorm:"type:int;not null"`
	Content    string    `gorm:"type:text;not null"`
	Vector     []float32 `gorm:"type:vector(1536)"`
	TokenCount int       `gorm:"type:int;default:0"`
	Metadata   string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"type:datetime;autoCreateTime"`
}

// TableName 指定表名
func (KBChunk) TableName() string {
	return "kb_chunks"
}

// CSAgent 客服人员
type CSAgent struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement"`
	UserID          uint64    `gorm:"type:bigint unsigned;uniqueIndex"`
	Nickname        string    `gorm:"type:varchar(100)"`
	Avatar          string    `gorm:"type:varchar(500)"`
	Status          string    `gorm:"type:varchar(20);default:'offline';index"`
	MaxConcurrent   int       `gorm:"type:int;default:5"`
	CurrentSessions int       `gorm:"type:int;default:0"`
	Stats           string    `gorm:"type:text"`
	CreatedAt       time.Time `gorm:"type:datetime;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"type:datetime;autoUpdateTime"`
}

// TableName 指定表名
func (CSAgent) TableName() string {
	return "cs_agents"
}

// CSTransfer 转接记录
type CSTransfer struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	SessionID   string    `gorm:"type:varchar(36);not null;index"`
	FromType    string    `gorm:"type:varchar(20);not null"`
	FromID      string    `gorm:"type:varchar(100)"`
	ToType      string    `gorm:"type:varchar(20);not null"`
	ToID        string    `gorm:"type:varchar(100)"`
	Reason      string    `gorm:"type:varchar(255)"`
	TransferredAt time.Time `gorm:"type:datetime;autoCreateTime;index"`
}

// TableName 指定表名
func (CSTransfer) TableName() string {
	return "cs_transfers"
}

// CSSFeedback 用户反馈
type CSSFeedback struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	SessionID   string    `gorm:"type:varchar(36);not null;index"`
	MessageID   string    `gorm:"type:varchar(36);index"`
	Rating      int       `gorm:"type:int;not null"`
	Comment     string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"type:datetime;autoCreateTime"`
}

// TableName 指定表名
func (CSSFeedback) TableName() string {
	return "css_feedback"
}
