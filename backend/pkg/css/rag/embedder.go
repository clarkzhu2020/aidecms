package rag

import (
	"context"
	"fmt"

	"github.com/clarkzhu2020/aidecms/pkg/ai"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Embedder 文本向量化器
type Embedder struct {
	aiClient ai.Client
}

// NewEmbedder 创建向量化器
func NewEmbedder(aiClient ai.Client) *Embedder {
	return &Embedder{
		aiClient: aiClient,
	}
}

// Embed 将文本转换为嵌入向量
func (e *Embedder) Embed(ctx context.Context, text string) ([]float32, error) {
	hlog.CtxInfof(ctx, "[Embedder] Embedding text (length: %d)", len(text))

	// 调用AI嵌入接口
	resp, err := e.aiClient.Embedding(ctx, text)
	if err != nil {
		hlog.CtxErrorf(ctx, "[Embedder] Failed to generate embedding: %v", err)
		return nil, fmt.Errorf("embedding failed: %w", err)
	}

	return resp.Embedding, nil
}

// EmbedBatch 批量向量化
func (e *Embedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	hlog.CtxInfof(ctx, "[Embedder] Batch embedding (count: %d)", len(texts))

	vectors := make([][]float32, len(texts))
	for i, text := range texts {
		vector, err := e.Embed(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("batch embedding failed at index %d: %w", i, err)
		}
		vectors[i] = vector
	}

	return vectors, nil
}

// GetDimension 获取向量维度
func (e *Embedder) GetDimension() int {
	// OpenAI text-embedding-ada-002 使用1536维
	// 也可以根据实际AI模型动态获取
	return 1536
}
