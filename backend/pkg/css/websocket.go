package css

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/websocket"
	"github.com/google/uuid"
)

// WebSocketMessage WebSocket消息
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	From    string      `json:"from,omitempty"`
	To      string      `json:"to,omitempty"`
}

// WSClient WebSocket客户端
type WSClient struct {
	ID         string
	SessionID  string
	Channel    string
	UserID     *uint64
	Send       chan WebSocketMessage
	Engine     *CSSEngine
	Connection *websocket.Conn
}

// WSManager WebSocket管理器
type WSManager struct {
	clients    map[*WSClient]bool
	register   chan *WSClient
	unregister chan *WSClient
	broadcast  chan WebSocketMessage
	mu         sync.RWMutex
	Engine     *CSSEngine
}

// NewWSManager 创建WebSocket管理器
func NewWSManager(engine *CSSEngine) *WSManager {
	return &WSManager{
		clients:   make(map[*WSClient]bool),
		register:  make(chan *WSClient, 256),
		unregister: make(chan *WSClient, 256),
		broadcast: make(chan WebSocketMessage, 256),
		Engine:    engine,
	}
}

// Start 启动WebSocket管理器
func (m *WSManager) Start(ctx context.Context) {
	hlog.Info("[WS] WebSocket manager started")

	for {
		select {
		case <-ctx.Done():
			hlog.Info("[WS] WebSocket manager stopped")
			return

		case client := <-m.register:
			m.registerClient(client)
			hlog.Infof("[WS] Client registered: %s (session: %s)", client.ID, client.SessionID)

		case client := <-m.unregister:
			m.unregisterClient(client)
			hlog.Infof("[WS] Client unregistered: %s", client.ID)

		case message := <-m.broadcast:
			m.broadcastMessage(message)
		}
	}
}

// Register 注册新客户端
func (m *WSManager) Register(ctx context.Context, c *app.RequestContext, conn *websocket.Conn, sessionID, channel string) (*WSClient, error) {
	clientID := uuid.New().String()

	client := &WSClient{
		ID:        clientID,
		SessionID:  sessionID,
		Channel:    channel,
		Send:       make(chan WebSocketMessage, 256),
		Engine:     m.Engine,
		Connection: conn,
	}

	// 发送到注册通道
	m.register <- client

	// 启动消息发送协程
	go m.readPump(ctx, client)
	go m.writePump(ctx, client)

	// 发送欢迎消息
	client.Send <- WebSocketMessage{
		Type: "connected",
		Data: map[string]interface{}{
			"session_id": sessionID,
			"client_id":  clientID,
			"timestamp":  time.Now().Unix(),
		},
	}

	return client, nil
}

// Unregister 注销客户端
func (m *WSManager) Unregister(client *WSClient) {
	m.unregister <- client
	close(client.Send)
}

// Broadcast 广播消息
func (m *WSManager) Broadcast(message WebSocketMessage) {
	m.broadcast <- message
}

// SendToSession 发送消息到指定会话
func (m *WSManager) SendToSession(sessionID string, message WebSocketMessage) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for client := range m.clients {
		if client.SessionID == sessionID {
			select {
			case client.Send <- message:
			default:
				hlog.Warnf("[WS] Send channel full for client %s", client.ID)
			}
		}
	}
}

// SendToClient 发送消息到指定客户端
func (m *WSManager) SendToClient(clientID string, message WebSocketMessage) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for client := range m.clients {
		if client.ID == clientID {
			select {
			case client.Send <- message:
			default:
				hlog.Warnf("[WS] Send channel full for client %s", client.ID)
			}
			return
		}
	}
}

// registerClient 注册客户端到map
func (m *WSManager) registerClient(client *WSClient) {
	m.mu.Lock()
	m.clients[client] = true
	m.mu.Unlock()
}

// unregisterClient 从map中移除客户端
func (m *WSManager) unregisterClient(client *WSClient) {
	m.mu.Lock()
	delete(m.clients, client)
	m.mu.Unlock()

	// 通知引擎客户端断开
	// TODO: 通知引擎清理会话
}

// broadcastMessage 广播消息到所有客户端
func (m *WSManager) broadcastMessage(message WebSocketMessage) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for client := range m.clients {
		select {
		case client.Send <- message:
		default:
			hlog.Warnf("[WS] Send channel full for client %s", client.ID)
		}
	}
}

// GetClientCount 获取客户端数量
func (m *WSManager) GetClientCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients)
}

// GetSessionClients 获取指定会话的所有客户端
func (m *WSManager) GetSessionClients(sessionID string) []*WSClient {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clients := make([]*WSClient, 0)
	for client := range m.clients {
		if client.SessionID == sessionID {
			clients = append(clients, client)
		}
	}

	return clients
}

// readPump 读取消息泵
func (m *WSManager) readPump(ctx context.Context, client *WSClient) {
	defer client.Connection.Close()

	for {
		messageType, data, err := client.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				hlog.Infof("[WS] Client %s disconnected unexpectedly: %v", client.ID, err)
			} else {
				hlog.Errorf("[WS] Read error for client %s: %v", client.ID, err)
			}
			m.Unregister(client)
			break
		}

		if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
			var wsMsg WebSocketMessage
			if err := json.Unmarshal(data, &wsMsg); err != nil {
				hlog.Errorf("[WS] Failed to unmarshal message: %v", err)
				continue
			}

			// 处理消息
			m.handleClientMessage(ctx, client, &wsMsg)
		}
	}
}

// writePump 写入消息泵
func (m *WSManager) writePump(ctx context.Context, client *WSClient) {
	ticker := time.NewTicker(54 * time.Second) // 心跳间隔
	defer ticker.Stop()
	defer client.Connection.Close()

	for {
		select {
		case <-ctx.Done():
			return

		case message, ok := <-client.Send:
			if !ok {
				// 通道关闭
				m.Unregister(client)
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				hlog.Errorf("[WS] Failed to marshal message: %v", err)
				continue
			}

			if err := client.Connection.WriteMessage(websocket.TextMessage, data); err != nil {
				hlog.Errorf("[WS] Failed to send message to %s: %v", client.ID, err)
				m.Unregister(client)
				return
			}

		case <-ticker.C:
			// 发送心跳
			pingMsg := WebSocketMessage{
				Type: "ping",
				Data: map[string]interface{}{
					"timestamp": time.Now().Unix(),
				},
			}
			if err := client.Connection.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
				hlog.Errorf("[WS] Failed to send ping to %s: %v", client.ID, err)
				m.Unregister(client)
				return
			}
		}
	}
}

// handleClientMessage 处理客户端消息
func (m *WSManager) handleClientMessage(ctx context.Context, client *WSClient, msg *WebSocketMessage) {
	hlog.CtxInfof(ctx, "[WS] Received message type: %s from %s", msg.Type, client.ID)

	switch msg.Type {
	case "message":
		// 用户提问消息
		m.handleUserMessage(ctx, client, msg)

	case "typing":
		// 用户正在输入
		m.handleTypingMessage(ctx, client)

	case "stop_typing":
		// 用户停止输入
		m.handleStopTypingMessage(ctx, client)

	case "feedback":
		// 用户反馈
		m.handleFeedbackMessage(ctx, client, msg)

	default:
		hlog.Warnf("[WS] Unknown message type: %s", msg.Type)
	}
}

// handleUserMessage 处理用户提问
func (m *WSManager) handleUserMessage(ctx context.Context, client *WSClient, msg *WebSocketMessage) {
	// 提取消息内容
	dataMap, ok := msg.Data.(map[string]interface{})
	if !ok {
		hlog.Warn("[WS] Invalid message data format")
		return
	}

	content, ok := dataMap["content"].(string)
	if !ok || content == "" {
		hlog.Warn("[WS] Empty message content")
		return
	}

	// 发送正在输入指示
	m.SendToSession(client.SessionID, WebSocketMessage{
		Type: "typing",
		From: "assistant",
		Data: map[string]interface{}{"session_id": client.SessionID},
	})

	// 调用引擎处理问题
	startTime := time.Now()
	answer, err := m.Engine.ProcessQuestion(ctx, client.SessionID, content)
	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		// 处理失败，发送错误消息
		m.SendToSession(client.SessionID, WebSocketMessage{
			Type: "error",
			From: "system",
			Data: map[string]interface{}{
				"message": "抱歉，服务暂时不可用，请稍后再试。",
				"error":   err.Error(),
			},
		})
		return
	}

	// 发送AI回答
	response := WebSocketMessage{
		Type: "message",
		From: "assistant",
		Data: map[string]interface{}{
			"session_id": client.SessionID,
			"content":    answer.Content,
			"confidence": answer.Confidence,
			"sources":    answer.Sources,
			"actions":    answer.SuggestedActions,
			"duration":   duration,
			"timestamp":  time.Now().Unix(),
		},
	}

	// 如果需要转接，添加转接信息
	if answer.TransferTo != nil {
		response.Data.(map[string]interface{})["transfer_to"] = *answer.TransferTo
	}

	m.SendToSession(client.SessionID, response)
}

// handleTypingMessage 处理用户正在输入
func (m *WSManager) handleTypingMessage(ctx context.Context, client *WSClient) {
	m.SendToSession(client.SessionID, WebSocketMessage{
		Type: "typing",
		From: "user",
		Data: map[string]interface{}{"session_id": client.SessionID},
	})
}

// handleStopTypingMessage 处理用户停止输入
func (m *WSManager) handleStopTypingMessage(ctx context.Context, client *WSClient) {
	m.SendToSession(client.SessionID, WebSocketMessage{
		Type: "stop_typing",
		From: "user",
		Data: map[string]interface{}{"session_id": client.SessionID},
	})
}

// handleFeedbackMessage 处理用户反馈
func (m *WSManager) handleFeedbackMessage(ctx context.Context, client *WSClient, msg *WebSocketMessage) {
	// TODO: 保存反馈到数据库
	hlog.CtxInfof(ctx, "[WS] User feedback received: %+v", msg.Data)
}
