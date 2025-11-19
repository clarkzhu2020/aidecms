package seo

import (
	"gorm.io/gorm"
)

// Redirect URL 重定向模型
type Redirect struct {
	gorm.Model
	FromURL    string `gorm:"uniqueIndex;size:500" json:"from_url"`
	ToURL      string `gorm:"size:500" json:"to_url"`
	StatusCode int    `gorm:"default:301" json:"status_code"` // 301 永久重定向, 302 临时重定向
	HitCount   int    `gorm:"default:0" json:"hit_count"`     // 重定向次数统计
	IsActive   bool   `gorm:"default:true" json:"is_active"`  // 是否启用
}

// RedirectManager 重定向管理器
type RedirectManager struct {
	db *gorm.DB
}

// NewRedirectManager 创建重定向管理器
func NewRedirectManager(db *gorm.DB) *RedirectManager {
	return &RedirectManager{db: db}
}

// Add 添加重定向规则
func (m *RedirectManager) Add(fromURL, toURL string, statusCode int) error {
	redirect := &Redirect{
		FromURL:    fromURL,
		ToURL:      toURL,
		StatusCode: statusCode,
		IsActive:   true,
	}
	return m.db.Create(redirect).Error
}

// Get 获取重定向规则
func (m *RedirectManager) Get(fromURL string) (*Redirect, error) {
	var redirect Redirect
	err := m.db.Where("from_url = ? AND is_active = ?", fromURL, true).First(&redirect).Error
	if err != nil {
		return nil, err
	}
	return &redirect, nil
}

// Update 更新重定向规则
func (m *RedirectManager) Update(id uint, toURL string, statusCode int) error {
	return m.db.Model(&Redirect{}).Where("id = ?", id).Updates(map[string]interface{}{
		"to_url":      toURL,
		"status_code": statusCode,
	}).Error
}

// Delete 删除重定向规则
func (m *RedirectManager) Delete(id uint) error {
	return m.db.Delete(&Redirect{}, id).Error
}

// List 获取所有重定向规则
func (m *RedirectManager) List(page, perPage int) ([]Redirect, int64, error) {
	var redirects []Redirect
	var total int64

	offset := (page - 1) * perPage

	if err := m.db.Model(&Redirect{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := m.db.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&redirects).Error; err != nil {
		return nil, 0, err
	}

	return redirects, total, nil
}

// IncrementHitCount 增加重定向计数
func (m *RedirectManager) IncrementHitCount(id uint) error {
	return m.db.Model(&Redirect{}).Where("id = ?", id).UpdateColumn("hit_count", gorm.Expr("hit_count + ?", 1)).Error
}

// Toggle 切换启用状态
func (m *RedirectManager) Toggle(id uint) error {
	var redirect Redirect
	if err := m.db.First(&redirect, id).Error; err != nil {
		return err
	}
	return m.db.Model(&redirect).Update("is_active", !redirect.IsActive).Error
}

// Migrate 执行数据库迁移
func (m *RedirectManager) Migrate() error {
	return m.db.AutoMigrate(&Redirect{})
}
