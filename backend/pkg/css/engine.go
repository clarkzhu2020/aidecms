package css

import (
	"context"
	"fmt"
	"time"

	"github.com/clarkzhu2020/aidecms/pkg/ai"
	"github.com/clarkzhu2020/aidecms/pkg/css/rag"
	"github.com/clarkzhu2020/aidecms/pkg/css/conversation"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// CSSEngine 客服系统核心引擎
type CSSEngine struct {
	aiManager       *ai.Manager
	ragRetriever    *rag.Retriever
	conversationMgr *conversation.Manager
	config         *Config
}

// Config 客服系统配置
type Config struct {
	// AI 配置
	DefaultModel     string  `json:"default_model"`     // 默认AI模型
	Temperature     float64 `json:"temperature"`      // 温度参数
	MaxTokens      int     `json:"max_tokens"`       // 最大Token数

	// RAG 配置
	TopK              int     `json:"top_k"`              // 检索Top-K文档
	ConfidenceThreshold float64 `json:"confidence_threshold"` // 置信度阈值

	// 转接配置
	TransferOnLowConfidence bool    `json:"transfer_on_low_confidence"`
	MaxRetries             int       `json:"max_retries"`
	TransferKeywords        []string  `json:"transfer_keywords"`

	// 会话配置
	SessionTimeout time.Duration `json:"session_timeout"`
	MaxHistory     int           `json:"max_history"`
}

// NewCSSEngine 创建客服系统引擎
func NewCSSEngine(aiManager *ai.Manager, ragRetriever *rag.Retriever, convMgr *conversation.Manager, config *Config) *CSSEngine {
	return &CSSEngine{
		aiManager:       aiManager,
		ragRetriever:    ragRetriever,
		conversationMgr: convMgr,
		config:         config,
	}
}

// Answer AI回答
type Answer struct {
	Content     string         `json:"content"`
	Confidence  float64        `json:"confidence"`
	Sources     []SourceRef    `json:"sources"`
	SuggestedActions []string       `json:"suggested_actions,omitempty"`
	TransferTo  *string        `json:"transfer_to,omitempty"`  // 如果需要转人工，设置客服ID
}

// SourceRef 知识来源引用
type SourceRef struct {
	DocumentID  string  `json:"document_id"`
	Title       string  `json:"title"`
	Relevance   float64 `json:"relevance"`
	Snippet     string  `json:"snippet"`
}

// ProcessQuestion 处理用户提问 - 核心流程实现
//
// 流程:
// 用户提问 → 会话管理 → RAG检索知识库 → 构建Prompt → AI生成
//        → 置信度评估 → (高置信度) 返回回答 / (低置信度) 转人工
func (e *CSSEngine) ProcessQuestion(ctx context.Context, sessionID, question string) (*Answer, error) {
	hlog.CtxInfof(ctx, "[CSS] Processing question for session %s: %s", sessionID, question)

	// 1. 会话管理：创建或恢复会话
	session, err := e.conversationMgr.GetOrCreateSession(ctx, sessionID)
	if err != nil {
		hlog.CtxErrorf(ctx, "[CSS] Failed to get session: %v", err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// 2. 保存用户消息到历史
	if err := e.conversationMgr.SaveMessage(ctx, sessionID, "user", question); err != nil {
		hlog.CtxErrorf(ctx, "[CSS] Failed to save user message: %v", err)
	}

	// 3. RAG检索：获取相关知识库文档
	docs, err := e.ragRetriever.Search(ctx, question, e.config.TopK)
	if err != nil {
		hlog.CtxErrorf(ctx, "[CSS] RAG search failed: %v", err)
		// 即使检索失败也继续，使用纯AI回答
		docs = nil
	}

	// 4. 获取对话历史（用于上下文）
	history, err := e.conversationMgr.GetHistory(ctx, sessionID, e.config.MaxHistory)
	if err != nil {
		hlog.CtxErrorf(ctx, "[CSS] Failed to get history: %v", err)
		history = nil
	}

	// 5. 构建Prompt（问题 + 上下文 + 知识库）
	prompt := e.buildPrompt(question, docs, history)

	// 6. AI生成：调用AI模型生成回答
	aiAnswer, err := e.callAI(ctx, prompt)
	if err != nil {
		hlog.CtxErrorf(ctx, "[CSS] AI generation failed: %v", err)
		
		// AI调用失败，尝试转人工
		if e.config.TransferOnLowConfidence {
			return e.transferToHuman(ctx, session, "ai_failed", "AI服务暂时不可用，已为您转接人工客服")
		}
		return nil, fmt.Errorf("ai generation failed: %w", err)
	}

	// 7. 置信度评估
	confidence := e.evaluateConfidence(question, aiAnswer, docs)
	
	// 8. 检查是否需要转接人工
	if e.shouldTransfer(ctx, session, question, confidence) {
		answer := &Answer{
			Content:    "您的问题比较复杂，我正在为您转接人工客服，请稍候...",
			Confidence: confidence,
			TransferTo:  strPtr("agent"),
		}
		
		// 异步转接（不阻塞用户）
		go func() {
			if err := e.transferToHuman(ctx, session, "low_confidence", question); err != nil {
				hlog.Errorf("[CSS] Transfer to human failed: %v", err)
			}
		}()
		
		return answer, nil
	}

	// 9. 高置信度：保存并返回AI回答
	sourceRefs := e.buildSourceRefs(docs)
	answer := &Answer{
		Content:          aiAnswer,
		Confidence:       confidence,
		Sources:          sourceRefs,
		SuggestedActions: e.suggestActions(question, aiAnswer, docs),
	}

	// 保存AI回答到历史
	if err := e.conversationMgr.SaveMessage(ctx, sessionID, "assistant", aiAnswer); err != nil {
		hlog.CtxErrorf(ctx, "[CSS] Failed to save assistant message: %v", err)
	}

	// 更新会话最后活跃时间
	if err := e.conversationMgr.UpdateSessionActivity(ctx, sessionID); err != nil {
		hlog.CtxErrorf(ctx, "[CSS] Failed to update session activity: %v", err)
	}

	hlog.CtxInfof(ctx, "[CSS] Answer generated with confidence %.2f", confidence)
	return answer, nil
}

// buildPrompt 构建AI Prompt
func (e *CSSEngine) buildPrompt(question string, docs []rag.Document, history []conversation.Message) string {
	prompt := "你是一个专业的客服助手。请基于以下知识库内容回答用户问题。\n\n"
	
	// 添加知识库上下文
	if len(docs) > 0 {
		prompt += "知识库内容：\n"
		for i, doc := range docs {
			prompt += fmt.Sprintf("\n【文档%d】%s\n内容：%s\n", i+1, doc.Title, doc.Content)
		}
		prompt += "\n"
	}
	
	// 添加对话历史
	if len(history) > 0 {
		prompt += "对话历史：\n"
		for _, msg := range history {
			role := "用户"
			if msg.Role == "assistant" {
				role = "客服"
			}
			prompt += fmt.Sprintf("%s：%s\n", role, msg.Content)
		}
		prompt += "\n"
	}
	
	prompt += fmt.Sprintf("用户问题：%s\n\n", question)
	prompt += "请基于知识库提供准确、专业的回答。如果知识库中没有相关信息，请诚实告知用户。"
	
	return prompt
}

// callAI 调用AI模型
func (e *CSSEngine) callAI(ctx context.Context, prompt string) (string, error) {
	client, err := e.aiManager.GetClient(e.config.DefaultModel)
	if err != nil {
		return "", fmt.Errorf("failed to get AI client: %w", err)
	}

	chatReq := &ai.ChatRequest{
		Messages: []*ai.Message{
			{
				Role:    "system",
				Content: "你是一个专业的客服助手，善于基于知识库回答用户问题，语言友好、专业。",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: &e.config.Temperature,
		MaxTokens:   &e.config.MaxTokens,
	}

	response, err := client.Chat(ctx, chatReq)
	if err != nil {
		return "", err
	}

	return response.Message.Content, nil
}

// evaluateConfidence 评估AI回答的置信度
func (e *CSSEngine) evaluateConfidence(question string, answer string, docs []rag.Document) float64 {
	confidence := 0.8 // 基础置信度

	// 因素1：是否有知识库支持
	if len(docs) > 0 {
		confidence += 0.15
	} else {
		confidence -= 0.2 // 没有知识库支持降低置信度
	}

	// 因素2：回答长度是否合理
	if len(answer) < 50 {
		confidence -= 0.1 // 回答太短可能不完整
	} else if len(answer) > 500 {
		confidence += 0.05 // 详细回答更可信
	}

	// 因素3：回答是否包含不确定词汇
	uncertainPhrases := []string{"可能", "也许", "不确定", "需要确认", "不清楚"}
	for _, phrase := range uncertainPhrases {
		if containsString(answer, phrase) {
			confidence -= 0.15
			break
		}
	}

	// 确保在[0, 1]范围内
	if confidence > 1.0 {
		confidence = 1.0
	} else if confidence < 0.0 {
		confidence = 0.0
	}

	return confidence
}

// shouldTransfer 判断是否需要转接人工
func (e *CSSEngine) shouldTransfer(ctx context.Context, session *conversation.Session, question string, confidence float64) bool {
	// 条件1：置信度低于阈值
	if e.config.TransferOnLowConfidence && confidence < e.config.ConfidenceThreshold {
		hlog.CtxInfof(ctx, "[CSS] Low confidence (%.2f < %.2f), transferring to human", confidence, e.config.ConfidenceThreshold)
		return true
	}

	// 条件2：检测转接关键词
	for _, keyword := range e.config.TransferKeywords {
		if containsString(question, keyword) {
			hlog.CtxInfof(ctx, "[CSS] Transfer keyword detected: %s", keyword)
			return true
		}
	}

	// 条件3：用户明确要求人工
	if containsString(question, "人工") || containsString(question, "客服") || containsString(question, "转接") {
		hlog.CtxInfof(ctx, "[CSS] User requested human agent")
		return true
	}

	return false
}

// transferToHuman 转接人工客服
func (e *CSSEngine) transferToHuman(ctx context.Context, session *conversation.Session, reason, userMessage string) error {
	hlog.CtxInfof(ctx, "[CSS] Transferring session %s to human (reason: %s)", session.ID, reason)

	// 1. 更新会话状态
	if err := e.conversationMgr.UpdateSessionStatus(ctx, session.ID, "transferred", nil); err != nil {
		return fmt.Errorf("failed to update session status: %w", err)
	}

	// 2. 查找可用客服
	agentID := e.findAvailableAgent()
	if agentID != nil {
		if err := e.conversationMgr.AssignSessionToAgent(ctx, session.ID, *agentID); err != nil {
			return fmt.Errorf("failed to assign session to agent: %w", err)
		}
	}

	// 3. 记录转接日志
	if err := e.conversationMgr.SaveTransferRecord(ctx, session.ID, "ai", "", "agent", fmt.Sprintf("%v", agentID), reason, userMessage); err != nil {
		hlog.CtxErrorf(ctx, "[CSS] Failed to save transfer record: %v", err)
	}

	// 4. 发布转接事件（通知客服工作台）
	// event.Publish("css.transfer", map[string]interface{}{
	//     "session_id": session.ID,
	//     "agent_id": agentID,
	//     "reason": reason,
	// })

	return nil
}

// findAvailableAgent 查找可用客服
func (e *CSSEngine) findAvailableAgent() *uint64 {
	// TODO: 实现客服队列和可用性管理
	// 简单实现：返回第一个可用客服ID
	return nil
}

// buildSourceRefs 构建知识来源引用
func (e *CSSEngine) buildSourceRefs(docs []rag.Document) []SourceRef {
	if len(docs) == 0 {
		return []SourceRef{}
	}

	refs := make([]SourceRef, len(docs))
	for i, doc := range docs {
		snippet := doc.Content
		if len(snippet) > 200 {
			snippet = snippet[:200] + "..."
		}

		refs[i] = SourceRef{
			DocumentID: doc.ID,
			Title:      doc.Title,
			Relevance:  doc.Relevance,
			Snippet:    snippet,
		}
	}
	return refs
}

// suggestActions 基于问题类型建议操作
func (e *CSSEngine) suggestActions(question string, answer string, docs []rag.Document) []string {
	actions := []string{}

	// 根据问题内容建议操作
	if containsString(question, "退款") {
		actions = append(actions, "查看退款政策")
		actions = append(actions, "发起退款申请")
	} else if containsString(question, "价格") {
		actions = append(actions, "查看价格表")
		actions = append(actions, "联系销售顾问")
	} else if containsString(question, "投诉") {
		actions = append(actions, "提交投诉单")
		actions = append(actions, "拨打投诉热线")
	} else if containsString(answer, "建议联系客服") {
		actions = append(actions, "转接人工客服")
	}

	return actions
}

// GetEngineStats 获取引擎统计信息
func (e *CSSEngine) GetEngineStats() map[string]interface{} {
	return map[string]interface{}{
		"config": e.config,
		"ai_model": e.config.DefaultModel,
		"rag_enabled": e.ragRetriever != nil,
		"session_count": e.conversationMgr.GetSessionCount(),
	}
}

// 辅助函数
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func strPtr(s string) *string {
	return &s
}
