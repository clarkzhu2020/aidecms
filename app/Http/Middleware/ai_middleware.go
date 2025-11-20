package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chenyusolar/aidecms/pkg/ai"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

const (
	AIManagerKey = "ai_manager"
	AIClientKey  = "ai_client"
)

// AIMiddleware AI中间件配置
type AIMiddleware struct {
	manager       *ai.Manager
	defaultClient string
}

// NewAIMiddleware 创建AI中间件
func NewAIMiddleware(manager *ai.Manager) *AIMiddleware {
	return &AIMiddleware{
		manager:       manager,
		defaultClient: manager.GetDefault(),
	}
}

// Handler AI中间件处理器
func (m *AIMiddleware) Handler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 将AI管理器注入到请求上下文
		c.Set(AIManagerKey, m.manager)

		// 如果有默认客户端，也注入默认客户端
		if m.defaultClient != "" {
			if client, err := m.manager.GetClient(m.defaultClient); err == nil {
				c.Set(AIClientKey, client)
			}
		}

		c.Next(ctx)
	}
}

// GetAIManager 从请求上下文获取AI管理器
func GetAIManager(c *app.RequestContext) (*ai.Manager, error) {
	manager, exists := c.Get(AIManagerKey)
	if !exists {
		return nil, fmt.Errorf("ai manager not found in context")
	}

	aiManager, ok := manager.(*ai.Manager)
	if !ok {
		return nil, fmt.Errorf("invalid ai manager type")
	}

	return aiManager, nil
}

// GetAIClient 从请求上下文获取AI客户端
func GetAIClient(c *app.RequestContext) (*ai.Client, error) {
	client, exists := c.Get(AIClientKey)
	if !exists {
		return nil, fmt.Errorf("ai client not found in context")
	}

	aiClient, ok := client.(*ai.Client)
	if !ok {
		return nil, fmt.Errorf("invalid ai client type")
	}

	return aiClient, nil
}

// ChatMiddleware 聊天中间件
type ChatMiddleware struct {
	manager      *ai.Manager
	clientName   string
	systemPrompt string
	autoResponse bool
}

// ChatMiddlewareConfig 聊天中间件配置
type ChatMiddlewareConfig struct {
	ClientName   string `json:"client_name"`
	SystemPrompt string `json:"system_prompt"`
	AutoResponse bool   `json:"auto_response"`
}

// NewChatMiddleware 创建聊天中间件
func NewChatMiddleware(manager *ai.Manager, config *ChatMiddlewareConfig) *ChatMiddleware {
	if config == nil {
		config = &ChatMiddlewareConfig{}
	}

	return &ChatMiddleware{
		manager:      manager,
		clientName:   config.ClientName,
		systemPrompt: config.SystemPrompt,
		autoResponse: config.AutoResponse,
	}
}

// Handler 聊天中间件处理器
func (m *ChatMiddleware) Handler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 检查是否是聊天请求
		if !m.isChatRequest(c) {
			c.Next(ctx)
			return
		}

		if m.autoResponse {
			// 自动处理聊天请求
			if err := m.handleChatRequest(ctx, c); err != nil {
				hlog.CtxErrorf(ctx, "chat middleware error: %v", err)
				c.JSON(500, map[string]interface{}{
					"error": "Internal server error",
				})
				return
			}
			return
		}

		// 只注入聊天相关的上下文，继续处理
		m.injectChatContext(c)
		c.Next(ctx)
	}
}

// isChatRequest 检查是否是聊天请求
func (m *ChatMiddleware) isChatRequest(c *app.RequestContext) bool {
	path := string(c.Path())
	return strings.Contains(path, "/chat") ||
		strings.Contains(path, "/ai") ||
		string(c.ContentType()) == "application/json" &&
			(string(c.Method()) == "POST" && strings.Contains(path, "/api"))
}

// handleChatRequest 处理聊天请求
func (m *ChatMiddleware) handleChatRequest(ctx context.Context, c *app.RequestContext) error {
	// 解析请求
	var req struct {
		Message string `json:"message"`
		Model   string `json:"model,omitempty"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "Invalid request format",
		})
		return nil
	}

	if req.Message == "" {
		c.JSON(400, map[string]interface{}{
			"error": "Message is required",
		})
		return nil
	}

	// 选择客户端
	clientName := req.Model
	if clientName == "" {
		clientName = m.clientName
	}

	client, err := m.manager.GetClient(clientName)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": fmt.Sprintf("Model not available: %v", err),
		})
		return nil
	}

	// 构建消息
	messages := []*ai.Message{}
	if m.systemPrompt != "" {
		messages = append(messages, &ai.Message{
			Role:    "system",
			Content: m.systemPrompt,
		})
	}
	messages = append(messages, &ai.Message{
		Role:    "user",
		Content: req.Message,
	})

	// 执行聊天
	response, err := client.Chat(ctx, &ai.ChatRequest{
		Messages: messages,
	})
	if err != nil {
		return err
	}

	// 返回响应
	c.JSON(200, map[string]interface{}{
		"message": response.Message.Content,
		"usage":   response.Usage,
	})

	return nil
}

// injectChatContext 注入聊天上下文
func (m *ChatMiddleware) injectChatContext(c *app.RequestContext) {
	// 注入管理器
	c.Set(AIManagerKey, m.manager)

	// 注入客户端
	if m.clientName != "" {
		if client, err := m.manager.GetClient(m.clientName); err == nil {
			c.Set(AIClientKey, client)
		}
	}

	// 注入系统提示
	if m.systemPrompt != "" {
		c.Set("ai_system_prompt", m.systemPrompt)
	}
}

// StreamMiddleware 流式响应中间件
type StreamMiddleware struct {
	manager    *ai.Manager
	clientName string
}

// NewStreamMiddleware 创建流式响应中间件
func NewStreamMiddleware(manager *ai.Manager, clientName string) *StreamMiddleware {
	return &StreamMiddleware{
		manager:    manager,
		clientName: clientName,
	}
}

// Handler 流式响应中间件处理器
func (m *StreamMiddleware) Handler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 检查是否请求流式响应
		if !m.isStreamRequest(c) {
			c.Next(ctx)
			return
		}

		// 处理流式请求
		if err := m.handleStreamRequest(ctx, c); err != nil {
			hlog.CtxErrorf(ctx, "stream middleware error: %v", err)
			c.JSON(500, map[string]interface{}{
				"error": "Internal server error",
			})
		}
	}
}

// isStreamRequest 检查是否是流式请求
func (m *StreamMiddleware) isStreamRequest(c *app.RequestContext) bool {
	return string(c.GetHeader("Accept")) == "text/event-stream" ||
		c.Query("stream") == "true"
}

// handleStreamRequest 处理流式请求
func (m *StreamMiddleware) handleStreamRequest(ctx context.Context, c *app.RequestContext) error {
	// 解析请求
	var req struct {
		Prompt string `json:"prompt"`
		Model  string `json:"model,omitempty"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"error": "Invalid request format",
		})
		return nil
	}

	// 选择客户端
	clientName := req.Model
	if clientName == "" {
		clientName = m.clientName
	}

	client, err := m.manager.GetClient(clientName)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": fmt.Sprintf("Model not available: %v", err),
		})
		return nil
	}

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// 创建流式写入器
	writer := &SSEWriter{ctx: c}

	// 执行流式补全
	if err := client.StreamCompletion(ctx, req.Prompt, writer); err != nil {
		return err
	}

	return nil
}

// SSEWriter Server-Sent Events写入器
type SSEWriter struct {
	ctx *app.RequestContext
}

// WriteResponse 写入响应
func (w *SSEWriter) WriteResponse(response *ai.StreamResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	// 写入SSE格式数据
	w.ctx.Write([]byte(fmt.Sprintf("data: %s\n\n", data)))
	w.ctx.Flush()

	return nil
}

// Write 实现io.Writer接口
func (w *SSEWriter) Write(data []byte) (int, error) {
	w.ctx.Write([]byte(fmt.Sprintf("data: %s\n\n", data)))
	w.ctx.Flush()
	return len(data), nil
}
