package ai

import (
	"context"
	"fmt"
)

// ProviderType 提供商类型
type ProviderType string

const (
	ProviderOpenAI    ProviderType = "openai"
	ProviderAnthropic ProviderType = "anthropic"
	ProviderDoubao    ProviderType = "doubao"
	ProviderQianwen   ProviderType = "qianwen"
	ProviderChatGLM   ProviderType = "chatglm"
	ProviderBaichuan  ProviderType = "baichuan"
	ProviderMiniMax   ProviderType = "minimax"
)

// ModelConfig 模型配置
type ModelConfig struct {
	Name        string                 `json:"name"`
	Provider    ProviderType           `json:"provider"`
	APIKey      string                 `json:"api_key"`
	APISecret   string                 `json:"api_secret,omitempty"`
	BaseURL     string                 `json:"base_url,omitempty"`
	Model       string                 `json:"model"`
	Temperature float64                `json:"temperature"`
	MaxTokens   int                    `json:"max_tokens"`
	TopP        float64                `json:"top_p,omitempty"`
	TopK        int                    `json:"top_k,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// EinoClient Eino框架客户端
type EinoClient struct {
	config   *ModelConfig
	aiClient *Client
	// einoModel  model.Model // 将在实际集成时添加
}

// NewEinoClient 创建Eino客户端
func NewEinoClient(config *ModelConfig) (*EinoClient, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// 转换配置格式
	clientConfig := &Config{
		Provider:    string(config.Provider),
		APIKey:      config.APIKey,
		BaseURL:     config.BaseURL,
		Model:       config.Model,
		Temperature: config.Temperature,
		MaxTokens:   config.MaxTokens,
		Options:     config.Options,
	}

	aiClient, err := NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create ai client: %w", err)
	}

	return &EinoClient{
		config:   config,
		aiClient: aiClient,
	}, nil
}

// ChatCompletion 聊天补全
func (e *EinoClient) ChatCompletion(ctx context.Context, messages []*Message, options ...ChatOption) (*ChatResponse, error) {
	req := &ChatRequest{
		Messages: messages,
	}

	// 应用配置默认值
	if e.config.Temperature > 0 {
		req.Temperature = &e.config.Temperature
	}
	if e.config.MaxTokens > 0 {
		req.MaxTokens = &e.config.MaxTokens
	}

	// 应用选项
	for _, opt := range options {
		opt(req)
	}

	return e.aiClient.Chat(ctx, req)
}

// TextCompletion 文本补全
func (e *EinoClient) TextCompletion(ctx context.Context, prompt string, options ...ChatOption) (string, error) {
	messages := []*Message{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := e.ChatCompletion(ctx, messages, options...)
	if err != nil {
		return "", err
	}

	return response.Message.Content, nil
}

// StreamCompletion 流式补全
func (e *EinoClient) StreamCompletion(ctx context.Context, prompt string, callback func(*StreamResponse) error) error {
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

	responseCh, errorCh := e.aiClient.StreamChat(ctx, req)

	for {
		select {
		case response, ok := <-responseCh:
			if !ok {
				return nil
			}

			// 假设 response 是字符串类型的简化版本
			content := fmt.Sprintf("%v", response)
			streamResp := &StreamResponse{
				Delta:   content,
				Message: content,
				Done:    false,
			}

			if err := callback(streamResp); err != nil {
				return err
			}

			// 简化的完成检查
			if len(content) == 0 {
				streamResp.Done = true
				if err := callback(streamResp); err != nil {
					return err
				}
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

// Embedding 生成嵌入向量
func (e *EinoClient) Embedding(ctx context.Context, texts []string) ([][]float64, error) {
	return e.aiClient.CreateEmbedding(ctx, texts)
}

// ChatOption 聊天选项
type ChatOption func(*ChatRequest)

// WithChatTemperature 设置温度
func WithChatTemperature(temperature float64) ChatOption {
	return func(req *ChatRequest) {
		req.Temperature = &temperature
	}
}

// WithChatMaxTokens 设置最大令牌数
func WithChatMaxTokens(maxTokens int) ChatOption {
	return func(req *ChatRequest) {
		req.MaxTokens = &maxTokens
	}
}

// WithSystemPrompt 添加系统提示
func WithSystemPrompt(systemPrompt string) ChatOption {
	return func(req *ChatRequest) {
		// 在消息开头插入系统消息
		systemMessage := &Message{
			Role:    "system",
			Content: systemPrompt,
		}
		req.Messages = append([]*Message{systemMessage}, req.Messages...)
	}
}

// ConversationContext 对话上下文
type ConversationContext struct {
	Messages []*Message `json:"messages"`
	MaxLen   int        `json:"max_len"`
}

// NewConversationContext 创建对话上下文
func NewConversationContext(maxLen int) *ConversationContext {
	if maxLen <= 0 {
		maxLen = 100 // 默认保持100条消息
	}
	return &ConversationContext{
		Messages: make([]*Message, 0),
		MaxLen:   maxLen,
	}
}

// AddMessage 添加消息
func (c *ConversationContext) AddMessage(role, content string) {
	message := &Message{
		Role:    role,
		Content: content,
	}

	c.Messages = append(c.Messages, message)

	// 保持消息数量不超过限制
	if len(c.Messages) > c.MaxLen {
		// 保留系统消息，删除最早的用户/助手消息
		systemMessages := make([]*Message, 0)
		otherMessages := make([]*Message, 0)

		for _, msg := range c.Messages {
			if msg.Role == "system" {
				systemMessages = append(systemMessages, msg)
			} else {
				otherMessages = append(otherMessages, msg)
			}
		}

		// 保留最近的消息
		keepCount := c.MaxLen - len(systemMessages)
		if keepCount > 0 && len(otherMessages) > keepCount {
			otherMessages = otherMessages[len(otherMessages)-keepCount:]
		}

		c.Messages = append(systemMessages, otherMessages...)
	}
}

// GetMessages 获取所有消息
func (c *ConversationContext) GetMessages() []*Message {
	return c.Messages
}

// Clear 清空消息
func (c *ConversationContext) Clear() {
	c.Messages = make([]*Message, 0)
}

// ConversationClient 对话客户端
type ConversationClient struct {
	client  *EinoClient
	context *ConversationContext
}

// NewConversationClient 创建对话客户端
func NewConversationClient(client *EinoClient, maxContextLen int) *ConversationClient {
	return &ConversationClient{
		client:  client,
		context: NewConversationContext(maxContextLen),
	}
}

// Chat 对话聊天（保持上下文）
func (c *ConversationClient) Chat(ctx context.Context, userMessage string, options ...ChatOption) (string, error) {
	// 添加用户消息到上下文
	c.context.AddMessage("user", userMessage)

	// 获取响应
	response, err := c.client.ChatCompletion(ctx, c.context.GetMessages(), options...)
	if err != nil {
		return "", err
	}

	// 添加助手响应到上下文
	c.context.AddMessage("assistant", response.Message.Content)

	return response.Message.Content, nil
}

// ClearHistory 清空对话历史
func (c *ConversationClient) ClearHistory() {
	c.context.Clear()
}

// GetHistory 获取对话历史
func (c *ConversationClient) GetHistory() []*Message {
	return c.context.GetMessages()
}

// Close 关闭客户端
func (e *EinoClient) Close() error {
	return e.aiClient.Close()
}
