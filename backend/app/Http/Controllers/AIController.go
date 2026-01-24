package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/clarkzhu2020/aidecms/pkg/ai"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// AIController AI控制器
type AIController struct {
	manager *ai.Manager
}

// NewAIController 创建AI控制器
func NewAIController(manager *ai.Manager) *AIController {
	return &AIController{
		manager: manager,
	}
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Message     string                 `json:"message" binding:"required"`
	Model       string                 `json:"model,omitempty"`
	Temperature *float64               `json:"temperature,omitempty"`
	MaxTokens   *int                   `json:"max_tokens,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
	Context     []*ai.Message          `json:"context,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Message string    `json:"message"`
	Usage   *ai.Usage `json:"usage,omitempty"`
	Model   string    `json:"model"`
}

// Chat 聊天接口
func (c *AIController) Chat(ctx context.Context, hCtx *app.RequestContext) {
	var req ChatRequest
	if err := hCtx.BindJSON(&req); err != nil {
		hCtx.JSON(400, map[string]interface{}{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	// 选择客户端
	clientName := req.Model
	if clientName == "" {
		clientName = c.manager.GetDefault()
	}

	client, err := c.manager.GetClient(clientName)
	if err != nil {
		hCtx.JSON(400, map[string]interface{}{
			"error": "Model not available: " + err.Error(),
		})
		return
	}

	// 构建消息
	messages := []*ai.Message{}
	if req.Context != nil {
		messages = append(messages, req.Context...)
	}
	messages = append(messages, &ai.Message{
		Role:    "user",
		Content: req.Message,
	})

	// 构建聊天请求
	chatReq := &ai.ChatRequest{
		Messages:    messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Stream:      req.Stream,
	}

	// 处理流式请求
	if req.Stream {
		c.handleStreamChat(ctx, hCtx, client, chatReq, clientName)
		return
	}

	// 执行聊天
	response, err := client.Chat(ctx, chatReq)
	if err != nil {
		hlog.CtxErrorf(ctx, "chat failed: %v", err)
		hCtx.JSON(500, map[string]interface{}{
			"error": "Chat failed",
		})
		return
	}

	// 返回响应
	hCtx.JSON(200, ChatResponse{
		Message: response.Message.Content,
		Usage:   response.Usage,
		Model:   clientName,
	})
}

// handleStreamChat 处理流式聊天
func (c *AIController) handleStreamChat(ctx context.Context, hCtx *app.RequestContext, client *ai.Client, req *ai.ChatRequest, modelName string) {
	// 设置SSE响应头
	hCtx.Header("Content-Type", "text/event-stream")
	hCtx.Header("Cache-Control", "no-cache")
	hCtx.Header("Connection", "keep-alive")
	hCtx.Header("Access-Control-Allow-Origin", "*")

	// 获取流式响应
	responseCh, errorCh := client.StreamChat(ctx, req)

	for {
		select {
		case response, ok := <-responseCh:
			if !ok {
				// 发送结束事件
				hCtx.Write([]byte("data: [DONE]\n\n"))
				hCtx.Flush()
				return
			}

			// 构建流式响应
			streamResp := map[string]interface{}{
				"message": response.Message, // response.Message 已经是字符串
				"model":   modelName,
			}

			data, _ := json.Marshal(streamResp)
			hCtx.Write([]byte(fmt.Sprintf("data: %s\n\n", data)))
			hCtx.Flush()

		case err := <-errorCh:
			if err != nil {
				hlog.CtxErrorf(ctx, "stream chat error: %v", err)
				errorData, _ := json.Marshal(map[string]string{"error": err.Error()})
				hCtx.Write([]byte(fmt.Sprintf("data: %s\n\n", errorData)))
				hCtx.Flush()
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// CompletionRequest 补全请求
type CompletionRequest struct {
	Prompt      string                 `json:"prompt" binding:"required"`
	Model       string                 `json:"model,omitempty"`
	Temperature *float64               `json:"temperature,omitempty"`
	MaxTokens   *int                   `json:"max_tokens,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// CompletionResponse 补全响应
type CompletionResponse struct {
	Text  string `json:"text"`
	Model string `json:"model"`
}

// Completion 文本补全接口
func (c *AIController) Completion(ctx context.Context, hCtx *app.RequestContext) {
	var req CompletionRequest
	if err := hCtx.BindJSON(&req); err != nil {
		hCtx.JSON(400, map[string]interface{}{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	// 选择客户端
	clientName := req.Model
	if clientName == "" {
		clientName = c.manager.GetDefault()
	}

	client, err := c.manager.GetClient(clientName)
	if err != nil {
		hCtx.JSON(400, map[string]interface{}{
			"error": "Model not available: " + err.Error(),
		})
		return
	}

	// 构建选项
	options := []ai.Option{}
	if req.Temperature != nil {
		options = append(options, ai.WithTemperature(*req.Temperature))
	}
	if req.MaxTokens != nil {
		options = append(options, ai.WithMaxTokens(*req.MaxTokens))
	}

	// 执行补全
	result, err := client.CreateCompletion(ctx, req.Prompt, options...)
	if err != nil {
		hlog.CtxErrorf(ctx, "completion failed: %v", err)
		hCtx.JSON(500, map[string]interface{}{
			"error": "Completion failed",
		})
		return
	}

	// 返回响应
	hCtx.JSON(200, CompletionResponse{
		Text:  result,
		Model: clientName,
	})
}

// EmbeddingRequest 嵌入请求
type EmbeddingRequest struct {
	Input []string `json:"input" binding:"required"`
	Model string   `json:"model,omitempty"`
}

// EmbeddingResponse 嵌入响应
type EmbeddingResponse struct {
	Embeddings [][]float64 `json:"embeddings"`
	Model      string      `json:"model"`
}

// Embedding 生成嵌入向量
func (c *AIController) Embedding(ctx context.Context, hCtx *app.RequestContext) {
	var req EmbeddingRequest
	if err := hCtx.BindJSON(&req); err != nil {
		hCtx.JSON(400, map[string]interface{}{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	// 选择客户端
	clientName := req.Model
	if clientName == "" {
		clientName = c.manager.GetDefault()
	}

	client, err := c.manager.GetClient(clientName)
	if err != nil {
		hCtx.JSON(400, map[string]interface{}{
			"error": "Model not available: " + err.Error(),
		})
		return
	}

	// 生成嵌入向量
	embeddings, err := client.CreateEmbedding(ctx, req.Input)
	if err != nil {
		hlog.CtxErrorf(ctx, "embedding failed: %v", err)
		hCtx.JSON(500, map[string]interface{}{
			"error": "Embedding failed",
		})
		return
	}

	// 返回响应
	hCtx.JSON(200, EmbeddingResponse{
		Embeddings: embeddings,
		Model:      clientName,
	})
}

// Models 获取可用模型列表
func (c *AIController) Models(ctx context.Context, hCtx *app.RequestContext) {
	models := c.manager.ListClients()

	modelList := make([]map[string]interface{}, 0, len(models))
	for _, name := range models {
		config, err := c.manager.GetConfig(name)
		if err != nil {
			continue
		}

		modelInfo := map[string]interface{}{
			"id":       name,
			"provider": config.Provider,
			"model":    config.Model,
		}
		modelList = append(modelList, modelInfo)
	}

	hCtx.JSON(200, map[string]interface{}{
		"models":  modelList,
		"default": c.manager.GetDefault(),
	})
}

// Health AI服务健康检查
func (c *AIController) Health(ctx context.Context, hCtx *app.RequestContext) {
	models := c.manager.ListClients()

	status := map[string]interface{}{
		"status":           "healthy",
		"models":           len(models),
		"available_models": models,
	}

	// 检查默认模型是否可用
	if defaultModel := c.manager.GetDefault(); defaultModel != "" {
		if _, err := c.manager.GetClient(defaultModel); err != nil {
			status["status"] = "degraded"
			status["error"] = "Default model unavailable"
		}
	}

	hCtx.JSON(200, status)
}

// ConversationRequest 对话请求
type ConversationRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Model     string `json:"model,omitempty"`
	MaxLen    *int   `json:"max_len,omitempty"`
	Clear     bool   `json:"clear,omitempty"`
}

// ConversationResponse 对话响应
type ConversationResponse struct {
	Message   string        `json:"message"`
	SessionID string        `json:"session_id"`
	History   []*ai.Message `json:"history,omitempty"`
	Model     string        `json:"model"`
}

// conversations 存储对话上下文（实际应该用数据库或缓存）
var conversations = make(map[string]*ai.ConversationClient)

// Conversation 对话接口（保持上下文）
func (c *AIController) Conversation(ctx context.Context, hCtx *app.RequestContext) {
	var req ConversationRequest
	if err := hCtx.BindJSON(&req); err != nil {
		hCtx.JSON(400, map[string]interface{}{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}

	// 选择客户端
	clientName := req.Model
	if clientName == "" {
		clientName = c.manager.GetDefault()
	}

	// 获取或创建对话客户端
	conversationClient, err := c.getOrCreateConversationClient(req.SessionID, clientName, req.MaxLen)
	if err != nil {
		hCtx.JSON(500, map[string]interface{}{
			"error": "Failed to create conversation client: " + err.Error(),
		})
		return
	}

	// 如果请求清空历史
	if req.Clear {
		conversationClient.ClearHistory()
		hCtx.JSON(200, map[string]interface{}{
			"message":    "History cleared",
			"session_id": req.SessionID,
		})
		return
	}

	// 执行对话
	response, err := conversationClient.Chat(ctx, req.Message)
	if err != nil {
		hlog.CtxErrorf(ctx, "conversation failed: %v", err)
		hCtx.JSON(500, map[string]interface{}{
			"error": "Conversation failed",
		})
		return
	}

	// 返回响应
	resp := ConversationResponse{
		Message:   response,
		SessionID: req.SessionID,
		Model:     clientName,
	}

	// 如果请求包含历史记录
	if strings.ToLower(hCtx.Query("include_history")) == "true" {
		resp.History = conversationClient.GetHistory()
	}

	hCtx.JSON(200, resp)
}

// getOrCreateConversationClient 获取或创建对话客户端
func (c *AIController) getOrCreateConversationClient(sessionID, clientName string, maxLen *int) (*ai.ConversationClient, error) {
	key := fmt.Sprintf("%s:%s", sessionID, clientName)

	if client, exists := conversations[key]; exists {
		return client, nil
	}

	// 获取AI客户端
	manager := c.manager
	_, err := manager.GetClient(clientName)
	if err != nil {
		return nil, err
	}

	// 创建Eino客户端（这里需要根据实际情况创建）
	config, err := manager.GetConfig(clientName)
	if err != nil {
		return nil, err
	}

	modelConfig := &ai.ModelConfig{
		Name:        clientName,
		Provider:    ai.ProviderType(config.Provider),
		APIKey:      config.APIKey,
		Model:       config.Model,
		Temperature: config.Temperature,
		MaxTokens:   config.MaxTokens,
	}

	einoClient, err := ai.NewEinoClient(modelConfig)
	if err != nil {
		return nil, err
	}

	// 设置最大上下文长度
	contextLen := 50 // 默认50条消息
	if maxLen != nil && *maxLen > 0 {
		contextLen = *maxLen
	}

	conversationClient := ai.NewConversationClient(einoClient, contextLen)
	conversations[key] = conversationClient

	return conversationClient, nil
}

// GetConversationHistory 获取对话历史
func (c *AIController) GetConversationHistory(ctx context.Context, hCtx *app.RequestContext) {
	sessionID := hCtx.Param("session_id")
	if sessionID == "" {
		hCtx.JSON(400, map[string]interface{}{
			"error": "Session ID is required",
		})
		return
	}

	clientName := hCtx.Query("model")
	if clientName == "" {
		clientName = c.manager.GetDefault()
	}

	key := fmt.Sprintf("%s:%s", sessionID, clientName)

	if client, exists := conversations[key]; exists {
		history := client.GetHistory()
		hCtx.JSON(200, map[string]interface{}{
			"session_id": sessionID,
			"model":      clientName,
			"history":    history,
			"count":      len(history),
		})
	} else {
		hCtx.JSON(404, map[string]interface{}{
			"error": "Conversation not found",
		})
	}
}

// ClearConversationHistory 清空对话历史
func (c *AIController) ClearConversationHistory(ctx context.Context, hCtx *app.RequestContext) {
	sessionID := hCtx.Param("session_id")
	if sessionID == "" {
		hCtx.JSON(400, map[string]interface{}{
			"error": "Session ID is required",
		})
		return
	}

	clientName := hCtx.Query("model")
	if clientName == "" {
		clientName = c.manager.GetDefault()
	}

	key := fmt.Sprintf("%s:%s", sessionID, clientName)

	if client, exists := conversations[key]; exists {
		client.ClearHistory()
		hCtx.JSON(200, map[string]interface{}{
			"message":    "History cleared",
			"session_id": sessionID,
		})
	} else {
		hCtx.JSON(404, map[string]interface{}{
			"error": "Conversation not found",
		})
	}
}
