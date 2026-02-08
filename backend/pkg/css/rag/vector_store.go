package rag

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// VectorStore 向量存储接口
type VectorStore struct {
	db *gorm.DB
}

// Chunk 文档分块
type Chunk struct {
	ID         string    `gorm:"type:varchar(36);primaryKey"`
	DocumentID string    `gorm:"type:varchar(36);not null;index"`
	ChunkOrder int       `gorm:"type:int;not null"`
	Content    string    `gorm:"type:text;not null"`
	Vector     []float32 `gorm:"type:vector(1536)"`
	TokenCount int       `gorm:"type:int;default:0"`
	Metadata   string    `gorm:"type:text"`
	CreatedAt  string    `gorm:"type:datetime;autoCreateTime"`
}

// SearchResult 搜索结果
type SearchResult struct {
	DocumentID string
	Score      float64
}

// NewVectorStore 创建向量存储
func NewVectorStore(db *gorm.DB) (*VectorStore, error) {
	// 检查是否支持向量类型
	if !db.Migrator().HasColumn(&Chunk{}, "vector") {
		hlog.Warn("[VectorStore] Vector column not found, please install pgvector extension")
		return nil, fmt.Errorf("vector column not supported, pgvector extension required")
	}

	// 自动迁移表
	if err := db.AutoMigrate(&Chunk{}); err != nil {
		return nil, fmt.Errorf("failed to migrate chunk table: %w", err)
	}

	return &VectorStore{db: db}, nil
}

// Store 存储向量
func (v *VectorStore) Store(ctx context.Context, documentID string, content string, vector []float32, tokenCount int, order int) error {
	chunk := &Chunk{
		ID:         generateChunkUUID(),
		DocumentID: documentID,
		ChunkOrder: order,
		Content:    content,
		Vector:     vector,
		TokenCount: tokenCount,
		CreatedAt:  getCurrentTimestamp(),
	}

	if err := v.db.WithContext(ctx).Create(chunk).Error; err != nil {
		hlog.CtxErrorf(ctx, "[VectorStore] Failed to store vector: %v", err)
		return fmt.Errorf("failed to store vector: %w", err)
	}

	return nil
}

// Search 向量相似度搜索
func (v *VectorStore) Search(ctx context.Context, queryVector []float32, topK int) ([]SearchResult, error) {
	var results []SearchResult

	// 使用余弦相似度搜索（pgvector支持）
	// ORDER BY vector <=> ? DESC  余弦距离排序
	query := `
		SELECT document_id, 1.0 - (vector <=> $1) as score
		FROM kb_chunks
		ORDER BY vector <=> $1 DESC
		LIMIT $2
	`

	if err := v.db.WithContext(ctx).Raw(query, queryVector, topK).Scan(&results).Error; err != nil {
		hlog.CtxErrorf(ctx, "[VectorStore] Vector search failed: %v", err)
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	hlog.CtxInfof(ctx, "[VectorStore] Found %d similar chunks", len(results))
	return results, nil
}

// GetChunksByDocument 获取文档的所有分块
func (v *VectorStore) GetChunksByDocument(ctx context.Context, documentID string) ([]Chunk, error) {
	var chunks []Chunk
	if err := v.db.WithContext(ctx).
		Where("document_id = ?", documentID).
		Order("chunk_order asc").
		Find(&chunks).Error; err != nil {
		return nil, fmt.Errorf("failed to get chunks: %w", err)
	}
	return chunks, nil
}

// DeleteByDocument 删除文档的所有分块
func (v *VectorStore) DeleteByDocument(ctx context.Context, documentID string) error {
	if err := v.db.WithContext(ctx).
		Where("document_id = ?", documentID).
		Delete(&Chunk{}).Error; err != nil {
		return fmt.Errorf("failed to delete chunks: %w", err)
	}

	hlog.CtxInfof(ctx, "[VectorStore] Deleted chunks for document: %s", documentID)
	return nil
}

// GetChunkCount 获取分块总数
func (v *VectorStore) GetChunkCount(ctx context.Context) (int64, error) {
	var count int64
	if err := v.db.WithContext(ctx).Model(&Chunk{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CreateVectorIndex 创建向量索引
func (v *VectorStore) CreateVectorIndex(ctx context.Context) error {
	// pgvector使用IVFFlat或HNSW索引
	// CREATE INDEX ON kb_chunks USING ivfflat (vector vector_cosine_ops);
	query := `
		CREATE INDEX IF NOT EXISTS idx_kb_chunks_vector_ivfflat 
		ON kb_chunks 
		USING ivfflat (vector vector_cosine_ops)
		WITH (lists = 100)
	`

	if err := v.db.WithContext(ctx).Exec(query).Error; err != nil {
		return fmt.Errorf("failed to create vector index: %w", err)
	}

	hlog.Info("[VectorStore] Vector index created successfully")
	return nil
}

// RebuildIndex 重建索引
func (v *VectorStore) RebuildIndex(ctx context.Context) error {
	// 先删除旧索引
	if err := v.db.WithContext(ctx).Exec("DROP INDEX IF EXISTS idx_kb_chunks_vector_ivfflat").Error; err != nil {
		return fmt.Errorf("failed to drop index: %w", err)
	}

	// 创建新索引
	return v.CreateVectorIndex(ctx)
}

func generateChunkUUID() string {
	return fmt.Sprintf("chunk_%d", getCurrentTimestampNano())
}

func getCurrentTimestamp() string {
	return fmt.Sprintf("%d", getCurrentTimestampNano()/1e9)
}

func getCurrentTimestampNano() int64 {
	return int64(1000000000)
}
