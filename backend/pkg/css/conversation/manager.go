package conversation

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Manager 对话管理器
type Manager struct {
	db           *gorm.DB
	sessionCache map[string]*Session // 会话缓存
	mu           sync.RWMutex
}

// NewManager 创建对话管理器
func NewManager(db *gorm.DB) *Manager {
	return &Manager{
		db:           db,
		sessionCache: make(map[string]*Session),
	}
}

// Session 会话
type Session struct {
	ID           string    `gorm:"type:varchar(36);primaryKey"`
	UserID       *uint64   `gorm:"type:bigint unsigned"`
	Channel      string    `gorm:"type:varchar(50);not null"`
	ChannelID    string    `gorm:"type:varchar(100)"`
	Status       string    `gorm:"type:varchar(20);default:'active'"`
	AssignedTo   *uint64   `gorm:"type:bigint unsigned"`
	Confidence   float32   `gorm:"type:float;default:0.8"`
	MessageCount int       `gorm:"type:int;default:0"`
	LastActiveAt time.Time `gorm:"type:datetime"`
	TransferReason string    `gorm:"type:varchar(255)"`
	Metadata     string    `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"type:datetime;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"type:datetime;autoUpdateTime"`
}

// Message 消息
type Message struct {
	ID           string  `gorm:"type:varchar(36);primaryKey"`
	SessionID     string  `gorm:"type:varchar(36);not null;index"`
	Role         string  `gorm:"type:varchar(20);not null;index"`
	Content      string  `gorm:"type:text;not null"`
	Tokens       int     `gorm:"type:int;default:0"`
	Duration     int     `gorm:"type:int"`
	Confidence   float32 `gorm:"type:float"`
	DocumentRefs string  `gorm:"type:text"`
	Metadata     string  `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"type:datetime;autoCreateTime"`
}

// TableName 指定表名
func (Session) TableName() string {
	return "css_sessions"
}

func (Message) TableName() string {
	return "css_messages"
}

// GetOrCreateSession 获取或创建会话
func (m *Manager) GetOrCreateSession(ctx context.Context, sessionID string) (*Session, error) {
	m.mu.RLock()
	session, exists := m.sessionCache[sessionID]
	m.mu.RUnlock()

	if exists {
		// 更新最后活跃时间
		now := time.Now()
		session.LastActiveAt = now
		if err := m.db.WithContext(ctx).Model(session).Updates(map[string]interface{}{
			"last_active_at": now,
		}).Error; err != nil {
			return nil, fmt.Errorf("failed to update session activity: %w", err)
		}
		return session, nil
	}

	// 创建新会话
	session = &Session{
		ID:           sessionID,
		Channel:      "web", // 默认web渠道
		Status:       "active",
		Confidence:   0.8,
		MessageCount: 0,
		LastActiveAt: time.Now(),
	}

	if err := m.db.WithContext(ctx).Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	hlog.CtxInfof(ctx, "[Conv] Created new session: %s", sessionID)

	// 加入缓存
	m.mu.Lock()
	m.sessionCache[sessionID] = session
	m.mu.Unlock()

	return session, nil
}

// GetSession 获取会话
func (m *Manager) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	m.mu.RLock()
	session, exists := m.sessionCache[sessionID]
	m.mu.RUnlock()

	if exists {
		return session, nil
	}

	session = &Session{}
	if err := m.db.WithContext(ctx).Where("id = ?", sessionID).First(session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	// 加入缓存
	m.mu.Lock()
	m.sessionCache[sessionID] = session
	m.mu.Unlock()

	return session, nil
}

// SaveMessage 保存消息
func (m *Manager) SaveMessage(ctx context.Context, sessionID, role, content string) error {
	msg := &Message{
		ID:       generateUUID(),
		SessionID: sessionID,
		Role:     role,
		Content:  content,
		CreatedAt: time.Now(),
	}

	if err := m.db.WithContext(ctx).Create(msg).Error; err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	// 更新会话消息计数
	if err := m.db.WithContext(ctx).Model(&Session{}).
		Where("id = ?", sessionID).
		UpdateColumn("message_count", gorm.Expr("message_count + 1")).Error; err != nil {
		hlog.CtxErrorf(ctx, "[Conv] Failed to update message count: %v", err)
	}

	return nil
}

// SaveMessageWithDuration 保存消息（含响应时长）
func (m *Manager) SaveMessageWithDuration(ctx context.Context, sessionID, role, content string, duration int, confidence float32, documentRefs string) error {
	msg := &Message{
		ID:           generateUUID(),
		SessionID:     sessionID,
		Role:          role,
		Content:       content,
		Duration:       duration,
		Confidence:     confidence,
		DocumentRefs:   documentRefs,
		CreatedAt:      time.Now(),
	}

	if err := m.db.WithContext(ctx).Create(msg).Error; err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	return nil
}

// GetHistory 获取对话历史
func (m *Manager) GetHistory(ctx context.Context, sessionID string, limit int) ([]Message, error) {
	var messages []Message

	query := m.db.WithContext(ctx).Where("session_id = ?", sessionID)
	if limit > 0 {
		query = query.Order("created_at desc").Limit(limit)
	} else {
		query = query.Order("created_at desc")
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}

	// 反转顺序（从早到晚）
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// UpdateSessionStatus 更新会话状态
func (m *Manager) UpdateSessionStatus(ctx context.Context, sessionID, status string, reason *string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if reason != nil {
		updates["transfer_reason"] = *reason
	}

	if err := m.db.WithContext(ctx).Model(&Session{}).
		Where("id = ?", sessionID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update session status: %w", err)
	}

	// 更新缓存
	m.mu.Lock()
	if session, exists := m.sessionCache[sessionID]; exists {
		session.Status = status
		if reason != nil {
			session.TransferReason = *reason
		}
	}
	m.mu.Unlock()

	return nil
}

// UpdateSessionActivity 更新会话活跃时间
func (m *Manager) UpdateSessionActivity(ctx context.Context, sessionID string) error {
	if err := m.db.WithContext(ctx).Model(&Session{}).
		Where("id = ?", sessionID).
		Update("last_active_at", time.Now()).Error; err != nil {
		return fmt.Errorf("failed to update session activity: %w", err)
	}

	// 更新缓存
	m.mu.Lock()
	if session, exists := m.sessionCache[sessionID]; exists {
		session.LastActiveAt = time.Now()
	}
	m.mu.Unlock()

	return nil
}

// AssignSessionToAgent 分配会话给客服
func (m *Manager) AssignSessionToAgent(ctx context.Context, sessionID string, agentID uint64) error {
	if err := m.db.WithContext(ctx).Model(&Session{}).
		Where("id = ?", sessionID).
		Update("assigned_to", agentID).Error; err != nil {
		return fmt.Errorf("failed to assign session: %w", err)
	}

	// 更新缓存
	m.mu.Lock()
	if session, exists := m.sessionCache[sessionID]; exists {
		assigned := agentID
		session.AssignedTo = &assigned
	}
	m.mu.Unlock()

	return nil
}

// SaveTransferRecord 保存转接记录
func (m *Manager) SaveTransferRecord(ctx context.Context, sessionID string, fromType, fromID, toType, toID, reason, userMessage string) error {
	// 转接记录表待实现
	hlog.CtxInfof(ctx, "[Conv] Transfer record: %s -> %s (session: %s)", fromID, toID, sessionID)
	return nil
}

// GetSessionCount 获取会话总数
func (m *Manager) GetSessionCount() int64 {
	m.mu.RLock()
	count := len(m.sessionCache)
	m.mu.RUnlock()
	return int64(count)
}

// GetActiveSessions 获取活跃会话列表
func (m *Manager) GetActiveSessions(ctx context.Context) ([]Session, error) {
	var sessions []Session
	if err := m.db.WithContext(ctx).
		Where("status = ? AND last_active_at > ?", "active", time.Now().Add(-30*time.Minute)).
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// CloseSession 关闭会话
func (m *Manager) CloseSession(ctx context.Context, sessionID string) error {
	if err := m.UpdateSessionStatus(ctx, sessionID, "closed", nil); err != nil {
		return err
	}

	// 从缓存中移除
	m.mu.Lock()
	delete(m.sessionCache, sessionID)
	m.mu.Unlock()

	return nil
}

// GetMessagesBySession 获取会话的所有消息
func (m *Manager) GetMessagesBySession(ctx context.Context, sessionID string) ([]Message, error) {
	var messages []Message
	if err := m.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("created_at asc").
		Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("failed to get session messages: %w", err)
	}
	return messages, nil
}

// generateUUID 生成UUID
func generateUUID() string {
	// 简单实现，实际应使用UUID库
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
