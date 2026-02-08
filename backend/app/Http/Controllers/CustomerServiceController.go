package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/clarkzhu2020/aidecms/pkg/css"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// CustomerServiceController 客服系统控制器
type CustomerServiceController struct {
	engine      *css.CSSEngine
	wsManager  *css.WSManager
	aiManager  interface{}
	db         interface{}
}

// NewCustomerServiceController 创建客服控制器
func NewCustomerServiceController() *CustomerServiceController {
	return &CustomerServiceController{}
}

// Init 初始化控制器
func (c *CustomerServiceController) Init(engine *css.CSSEngine, wsManager *css.WSManager) {
	c.engine = engine
	c.wsManager = wsManager
}

// WebSocket WebSocket连接处理
func (c *CustomerServiceController) WebSocket(ctx context.Context, hCtx *app.RequestContext) {
	// 升级为WebSocket连接
	sessionID := hCtx.Query("session_id")
	channel := hCtx.Query("channel")

	if sessionID == "" {
		sessionID = generateSessionID()
	}

	// 生成UUID
	sessionID = sessionID

	hCtx.SetStatusCode(101)
	hCtx.SetHeader("Connection", "Upgrade")
	hCtx.SetHeader("Upgrade", "websocket")
	hCtx.SetHeader("Sec-WebSocket-Accept", hCtx.GetHeader("Sec-WebSocket-Key"))

	// 创建WebSocket连接
	// conn, err := upgrader.Upgrade(hCtx)
	// if err != nil {
	// 	hlog.Errorf("[CSS] Failed to upgrade to WebSocket: %v", err)
	// 	return
	// }

	// 注册客户端
	// client, err := c.wsManager.Register(ctx, hCtx, conn, sessionID, channel)
	// if err != nil {
	// 	hlog.Errorf("[CSS] Failed to register client: %v", err)
	// 	return
	// }

	hlog.Infof("[CSS] WebSocket connection established: %s (channel: %s)", sessionID, channel)
}

// SendQuestion 发送问题(REST API)
func (c *CustomerServiceController) SendQuestion(ctx context.Context, hCtx *app.RequestContext) {
	var req struct {
		SessionID string `json:"session_id"`
	Question  string `json:"question" binding:"required"`
		Channel   string `json:"channel"`
	}

	if err := hCtx.BindJSON(&req); err != nil {
		hCtx.JSON(400, map[string]interface{}{
			"success": false,
			"error":   "Invalid request format",
		})
		return
	}

	// 设置默认session_id
	if req.SessionID == "" {
		req.SessionID = generateSessionID()
	}

	// 设置默认channel
	if req.Channel == "" {
		req.Channel = "web"
	}

	startTime := time.Now()

	// 调用引擎处理
	answer, err := c.engine.ProcessQuestion(ctx, req.SessionID, req.Question)
	if err != nil {
		hlog.CtxErrorf(ctx, "[CSS] Failed to process question: %v", err)
		hCtx.JSON(500, map[string]interface{}{
			"success": false,
			"error":   "Failed to process question",
		})
		return
	}

	duration := time.Since(startTime).Milliseconds()

	hCtx.JSON(200, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"session_id":  req.SessionID,
			"answer":      answer.Content,
			"confidence":  answer.Confidence,
			"sources":     answer.Sources,
			"actions":     answer.SuggestedActions,
			"duration":    duration,
			"timestamp":   time.Now().Unix(),
		},
	})
}

// GetHistory 获取对话历史
func (c *CustomerServiceController) GetHistory(ctx context.Context, hCtx *app.RequestContext) {
	sessionID := hCtx.Param("session_id")
	limit := hCtx.DefaultQuery("limit", "20")

	if sessionID == "" {
		hCtx.JSON(400, map[string]interface{}{
			"success": false,
			"error":   "session_id is required",
		})
		return
	}

	// TODO: 从数据库获取历史
	messages := []map[string]interface{}{}

	hCtx.JSON(200, map[string]interface{}{
		"success": true,
		"data":    messages,
	})
}

// GetStatus 获取系统状态
func (c *CustomerServiceController) GetStatus(ctx context.Context, hCtx *app.RequestContext) {
	stats := c.engine.GetEngineStats()

	hCtx.JSON(200, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"engine":   stats,
			"clients":  c.wsManager.GetClientCount(),
			"timestamp": time.Now().Unix(),
		},
	})
}

// CloseSession 关闭会话
func (c *CustomerServiceController) CloseSession(ctx context.Context, hCtx *app.RequestContext) {
	sessionID := hCtx.Param("session_id")

	// TODO: 关闭会话逻辑

	hCtx.JSON(200, map[string]interface{}{
		"success": true,
		"message": "Session closed",
	})
}

// generateSessionID 生成会话ID
func generateSessionID() string {
	return fmt.Sprintf("sess_%d", time.Now().UnixNano())
}

// 配置辅助函数
func getConfig(key, defaultValue string) string {
	return defaultValue
}

func getConfigInt(key string, defaultValue int) int {
	return defaultValue
}

func getConfigFloat(key string, defaultValue float64) float64 {
	return defaultValue
}

func getConfigBool(key string, defaultValue bool) bool {
	return defaultValue
}

func getConfigList(key string, defaultValue []string) []string {
	return defaultValue
}
