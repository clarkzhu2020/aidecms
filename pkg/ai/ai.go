package ai

import (
	"context"
	"fmt"
	"io"
)

// Config AI配置
type Config struct {
	Provider    string                 `json:"provider"` // openai, anthropic, doubao等
	APIKey      string                 `json:"api_key"`
	BaseURL     string                 `json:"base_url"`
	Model       string                 `json:"model"`
	Temperature float64                `json:"temperature"`
	MaxTokens   int                    `json:"max_tokens"`
	Options     map[string]interface{} `json:"options"`
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"` // system, user, assistant
	Content string `json:"content"`
}

// Client AI客户端
type Client struct {
	config *Config
}

// NewClient 创建AI客户端
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("api_key is required")
	}

	if config.Model == "" {
		return nil, fmt.Errorf("model is required")
	}

	client := &Client{
		config: config,
	}

	return client, nil
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Messages    []*Message `json:"messages"`
	Temperature *float64   `json:"temperature,omitempty"`
	MaxTokens   *int       `json:"max_tokens,omitempty"`
	Stream      bool       `json:"stream,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Message *Message `json:"message"`
	Usage   *Usage   `json:"usage,omitempty"`
}

// StreamResponse 流式响应
type StreamResponse struct {
	Delta   string `json:"delta"`
	Message string `json:"message"`
	Done    bool   `json:"done"`
}

// Usage 使用统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Chat 聊天对话
func (c *Client) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// TODO: 集成实际的eino实现
	// 这里先返回一个模拟响应
	return &ChatResponse{
		Message: &Message{
			Role:    "assistant",
			Content: "This is a mock response. Eino integration will be implemented soon.",
		},
		Usage: &Usage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
	}, nil
}

// StreamChat 流式聊天
func (c *Client) StreamChat(ctx context.Context, req *ChatRequest) (<-chan *StreamResponse, <-chan error) {
	responseCh := make(chan *StreamResponse, 100)
	errorCh := make(chan error, 1)

	go func() {
		defer close(responseCh)
		defer close(errorCh)

		// TODO: 实现流式聊天
		// 这需要根据eino的实际API来实现
		responseCh <- &StreamResponse{
			Delta:   "Mock",
			Message: "Mock streaming response",
			Done:    true,
		}
	}()

	return responseCh, errorCh
}

// CreateEmbedding 创建嵌入向量
func (c *Client) CreateEmbedding(ctx context.Context, texts []string) ([][]float64, error) {
	// TODO: 实现嵌入向量生成
	embeddings := make([][]float64, len(texts))
	for i := range texts {
		// 返回模拟的向量
		embeddings[i] = make([]float64, 1536) // 假设向量维度为1536
		for j := range embeddings[i] {
			embeddings[i][j] = 0.1 // 模拟值
		}
	}
	return embeddings, nil
}

// CreateCompletion 文本补全
func (c *Client) CreateCompletion(ctx context.Context, prompt string, options ...Option) (string, error) {
	// 构建消息
	messages := []*Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	req := &ChatRequest{
		Messages: messages,
	}

	// 应用选项
	for _, opt := range options {
		opt(req)
	}

	// 执行聊天
	resp, err := c.Chat(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Message.Content, nil
}

// CompletionWriter 用于流式写入的接口
type CompletionWriter interface {
	io.Writer
	WriteResponse(*StreamResponse) error
}

// StreamCompletion 流式文本补全
func (c *Client) StreamCompletion(ctx context.Context, prompt string, writer CompletionWriter) error {
	messages := []*Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	req := &ChatRequest{
		Messages: messages,
		Stream:   true,
	}

	responseCh, errorCh := c.StreamChat(ctx, req)

	for {
		select {
		case response, ok := <-responseCh:
			if !ok {
				return nil
			}
			if err := writer.WriteResponse(response); err != nil {
				return err
			}
			if response.Done {
				return nil
			}
		case err := <-errorCh:
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Option 选项函数
type Option func(*ChatRequest)

// WithTemperature 设置温度
func WithTemperature(temperature float64) Option {
	return func(req *ChatRequest) {
		req.Temperature = &temperature
	}
}

// WithMaxTokens 设置最大令牌数
func WithMaxTokens(maxTokens int) Option {
	return func(req *ChatRequest) {
		req.MaxTokens = &maxTokens
	}
}

// WithStream 设置流式输出
func WithStream(stream bool) Option {
	return func(req *ChatRequest) {
		req.Stream = stream
	}
}

// GetConfig 获取配置
func (c *Client) GetConfig() *Config {
	return c.config
}

// Close 关闭客户端
func (c *Client) Close() error {
	// 清理资源
	return nil
}
