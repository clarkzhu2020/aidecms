package rag

import (
	"context"
	"fmt"

	"github.com/clarkzhu2020/aidecms/pkg/ai"
	"github.com/clarkzhu2020/aidecms/pkg/css/kb"
	"gorm.io/gorm"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Retriever RAG检索器
type Retriever struct {
	db           *gorm.DB
	embedder     *Embedder
	vectorStore  *VectorStore
	kbService    *kb.Service
}

// NewRetriever 创建RAG检索器
func NewRetriever(db *gorm.DB, aiClient ai.Client, kbService *kb.Service) (*Retriever, error) {
	embedder := NewEmbedder(aiClient)
	vectorStore, err := NewVectorStore(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create vector store: %w", err)
	}

	return &Retriever{
		db:          db,
		embedder:    embedder,
		vectorStore: vectorStore,
		kbService:   kbService,
	}, nil
}

// Document 文档
type Document struct {
	ID       string
	Title    string
	Content  string
	Relevance float64
}

// SearchResult 搜索结果
type SearchResult struct {
	DocumentID string
	Title      string
	Content    string
	Score      float64
	Metadata   map[string]interface{}
}

// Search 执行RAG检索
func (r *Retriever) Search(ctx context.Context, query string, topK int) ([]Document, error) {
	hlog.CtxInfof(ctx, "[RAG] Searching for: %s (topK=%d)", query, topK)

	// 1. 将查询向量化
	queryVector, err := r.embedder.Embed(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "[RAG] Failed to embed query: %v", err)
		return nil, err
	}

	// 2. 向量相似度检索
	vectorResults, err := r.vectorStore.Search(ctx, queryVector, topK)
	if err != nil {
		hlog.CtxErrorf(ctx, "[RAG] Vector search failed: %v", err)
		return nil, err
	}

	// 3. 关键词全文检索（混合检索）
	keywordResults, err := r.kbService.FullTextSearch(ctx, query, topK)
	if err != nil {
		hlog.CtxErrorf(ctx, "[RAG] Full-text search failed: %v", err)
		// 关键词搜索失败不影响向量检索结果
	}

	// 4. 结果融合和重排序
	docs := r.mergeResults(vectorResults, keywordResults, topK)

	hlog.CtxInfof(ctx, "[RAG] Found %d documents", len(docs))
	return docs, nil
}

// mergeResults 融合向量检索和关键词检索结果
func (r *Retriever) mergeResults(vectorResults, keywordResults []SearchResult, topK int) []Document {
	// 创建评分map，避免重复
	scores := make(map[string]float64)

	// 向量检索结果权重: 0.7
	for _, result := range vectorResults {
		scores[result.DocumentID] = result.Score * 0.7
	}

	// 关键词检索结果权重: 0.3
	for _, result := range keywordResults {
		if existingScore, exists := scores[result.DocumentID]; exists {
			scores[result.DocumentID] = existingScore + result.Score*0.3
		} else {
			scores[result.DocumentID] = result.Score * 0.3
		}
	}

	// 获取完整文档内容
	docs := make([]Document, 0, topK)
	for docID, score := range scores {
		doc, err := r.kbService.GetDocumentByID(context.Background(), docID)
		if err != nil {
			continue
		}

		docs = append(docs, Document{
			ID:       docID,
			Title:    doc.Title,
			Content:  doc.Content,
			Relevance: score,
		})

		if len(docs) >= topK {
			break
		}
	}

	// 按相关性排序
	sortDocumentsByRelevance(docs)

	return docs
}

// sortDocumentsByRelevance 按相关性排序
func sortDocumentsByRelevance(docs []Document) {
	n := len(docs)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if docs[j].Relevance < docs[j+1].Relevance {
				docs[j], docs[j+1] = docs[j+1], docs[j]
			}
		}
	}
}

// GetRetrieverStats 获取检索器统计信息
func (r *Retriever) GetRetrieverStats(ctx context.Context) (map[string]interface{}, error) {
	docCount, err := r.kbService.GetDocumentCount(ctx)
	if err != nil {
		return nil, err
	}

	chunkCount, err := r.vectorStore.GetChunkCount(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_documents": docCount,
		"total_chunks":   chunkCount,
		"vector_search_enabled": r.vectorStore != nil,
		"fulltext_search_enabled": true,
	}, nil
}
