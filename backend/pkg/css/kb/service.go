package kb

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Service 知识库服务
type Service struct {
	db *gorm.DB
}

// Document 文档
type Document struct {
	ID          string  `gorm:"type:varchar(36);primaryKey"`
	Title       string  `gorm:"type:varchar(255);not null"`
	Category    string  `gorm:"type:varchar(100)"`
	Tags        string  `gorm:"type:text"`
	Content     string  `gorm:"type:text;not null"`
	FileURL     string  `gorm:"type:varchar(500)"`
	FileType    string  `gorm:"type:varchar(50)"`
	FileSize    int64   `gorm:"type:bigint"`
	ChunkCount  int     `gorm:"type:int;default:0"`
	Status      string  `gorm:"type:varchar(20);default:'processing'"`
	Version     int     `gorm:"type:int;default:1"`
	CreatedBy   uint64  `gorm:"type:bigint unsigned"`
	CreatedAt   string  `gorm:"type:datetime;autoCreateTime"`
	UpdatedAt   string  `gorm:"type:datetime;autoUpdateTime"`
}

// SearchResult 搜索结果
type SearchResult struct {
	DocumentID string
	Title      string
	Score      float64
}

// NewService 创建知识库服务
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateDocument 创建文档
func (s *Service) CreateDocument(ctx context.Context, doc *Document) error {
	doc.ID = generateDocUUID()
	doc.CreatedAt = getCurrentTimestamp()
	doc.UpdatedAt = getCurrentTimestamp()

	if err := s.db.WithContext(ctx).Create(doc).Error; err != nil {
		hlog.CtxErrorf(ctx, "[KB] Failed to create document: %v", err)
		return fmt.Errorf("failed to create document: %w", err)
	}

	hlog.CtxInfof(ctx, "[KB] Document created: %s", doc.ID)
	return nil
}

// GetDocumentByID 根据ID获取文档
func (s *Service) GetDocumentByID(ctx context.Context, id string) (*Document, error) {
	doc := &Document{}
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(doc).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	return doc, nil
}

// GetDocuments 获取文档列表
func (s *Service) GetDocuments(ctx context.Context, category, status string, page, pageSize int) ([]Document, int64, error) {
	var docs []Document
	var total int64

	query := s.db.WithContext(ctx).Model(&Document{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&docs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get documents: %w", err)
	}

	return docs, total, nil
}

// UpdateDocument 更新文档
func (s *Service) UpdateDocument(ctx context.Context, id string, updates map[string]interface{}) error {
	updates["updated_at"] = getCurrentTimestamp()

	if err := s.db.WithContext(ctx).Model(&Document{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	hlog.CtxInfof(ctx, "[KB] Document updated: %s", id)
	return nil
}

// DeleteDocument 删除文档
func (s *Service) DeleteDocument(ctx context.Context, id string) error {
	if err := s.db.WithContext(ctx).Where("id = ?", id).Delete(&Document{}).Error; err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	hlog.CtxInfof(ctx, "[KB] Document deleted: %s", id)
	return nil
}

// UpdateDocumentStatus 更新文档状态
func (s *Service) UpdateDocumentStatus(ctx context.Context, id, status string) error {
	return s.UpdateDocument(ctx, id, map[string]interface{}{
		"status":     status,
		"updated_at": getCurrentTimestamp(),
	})
}

// IncrementChunkCount 增加分块计数
func (s *Service) IncrementChunkCount(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Model(&Document{}).
		Where("id = ?", id).
		UpdateColumn("chunk_count", gorm.Expr("chunk_count + 1")).Error
}

// FullTextSearch 全文搜索
func (s *Service) FullTextSearch(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	var results []SearchResult

	// 使用LIKE进行模糊搜索
	query := s.db.WithContext(ctx).Table("kb_documents").
		Select("id as document_id, title, ? as score", "1.0").
		Where("title LIKE ? OR content LIKE ?", "%"+query+"%", "%"+query+"%").
		Where("status = ?", "indexed").
		Order("created_at desc").
		Limit(limit)

	if err := query.Scan(&results).Error; err != nil {
		hlog.CtxErrorf(ctx, "[KB] Full-text search failed: %v", err)
		return nil, fmt.Errorf("full-text search failed: %w", err)
	}

	hlog.CtxInfof(ctx, "[KB] Full-text search found %d results", len(results))
	return results, nil
}

// GetDocumentCount 获取文档总数
func (s *Service) GetDocumentCount(ctx context.Context) (int64, error) {
	var count int64
	if err := s.db.WithContext(ctx).Model(&Document{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetCategories 获取分类列表
func (s *Service) GetCategories(ctx context.Context) ([]string, error) {
	var categories []string
	if err := s.db.WithContext(ctx).Model(&Document{}).
		Where("status = ?", "indexed").
		Distinct("category").
		Pluck("category", &categories).Error; err != nil {
		return nil, err
	}

	// 过滤空值
	filtered := make([]string, 0)
	for _, cat := range categories {
		if cat != "" {
			filtered = append(filtered, cat)
		}
	}

	return filtered, nil
}

// ParseTags 解析标签
func ParseTags(tagsJSON string) []string {
	if tagsJSON == "" {
		return []string{}
	}

	// 简单实现：逗号分隔
	tags := strings.Split(tagsJSON, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	return tags
}

// FormatTags 格式化标签
func FormatTags(tags []string) string {
	return strings.Join(tags, ",")
}

// generateDocUUID 生成文档UUID
func generateDocUUID() string {
	return fmt.Sprintf("doc_%d", getCurrentTimestampNano())
}

func getCurrentTimestamp() string {
	return fmt.Sprintf("%d", getCurrentTimestampNano()/1e9)
}

func getCurrentTimestampNano() int64 {
	return int64(1000000000)
}
