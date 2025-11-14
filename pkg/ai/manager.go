package ai

import (
	"context"
	"fmt"
	"sync"
)

// Manager AI管理器
type Manager struct {
	clients     map[string]*Client
	configs     map[string]*Config
	mutex       sync.RWMutex
	defaultName string
}

// NewManager 创建AI管理器
func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]*Client),
		configs: make(map[string]*Config),
	}
}

// AddClient 添加AI客户端
func (m *Manager) AddClient(name string, config *Config) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, err := NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create client %s: %w", name, err)
	}

	m.clients[name] = client
	m.configs[name] = config

	// 如果是第一个客户端，设置为默认
	if len(m.clients) == 1 {
		m.defaultName = name
	}

	return nil
}

// GetClient 获取AI客户端
func (m *Manager) GetClient(name string) (*Client, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if name == "" {
		name = m.defaultName
	}

	client, exists := m.clients[name]
	if !exists {
		return nil, fmt.Errorf("client %s not found", name)
	}

	return client, nil
}

// GetDefaultClient 获取默认AI客户端
func (m *Manager) GetDefaultClient() (*Client, error) {
	return m.GetClient("")
}

// SetDefault 设置默认客户端
func (m *Manager) SetDefault(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.clients[name]; !exists {
		return fmt.Errorf("client %s not found", name)
	}

	m.defaultName = name
	return nil
}

// RemoveClient 移除AI客户端
func (m *Manager) RemoveClient(name string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	client, exists := m.clients[name]
	if !exists {
		return fmt.Errorf("client %s not found", name)
	}

	// 关闭客户端
	if err := client.Close(); err != nil {
		return fmt.Errorf("failed to close client %s: %w", name, err)
	}

	delete(m.clients, name)
	delete(m.configs, name)

	// 如果删除的是默认客户端，重新选择默认客户端
	if m.defaultName == name {
		for clientName := range m.clients {
			m.defaultName = clientName
			break
		}
		if len(m.clients) == 0 {
			m.defaultName = ""
		}
	}

	return nil
}

// ListClients 列出所有客户端
func (m *Manager) ListClients() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	names := make([]string, 0, len(m.clients))
	for name := range m.clients {
		names = append(names, name)
	}
	return names
}

// Chat 使用指定客户端进行聊天
func (m *Manager) Chat(ctx context.Context, clientName string, req *ChatRequest) (*ChatResponse, error) {
	client, err := m.GetClient(clientName)
	if err != nil {
		return nil, err
	}
	return client.Chat(ctx, req)
}

// CreateCompletion 使用指定客户端进行文本补全
func (m *Manager) CreateCompletion(ctx context.Context, clientName string, prompt string, options ...Option) (string, error) {
	client, err := m.GetClient(clientName)
	if err != nil {
		return "", err
	}
	return client.CreateCompletion(ctx, prompt, options...)
}

// CreateEmbedding 使用指定客户端创建嵌入向量
func (m *Manager) CreateEmbedding(ctx context.Context, clientName string, texts []string) ([][]float64, error) {
	client, err := m.GetClient(clientName)
	if err != nil {
		return nil, err
	}
	return client.CreateEmbedding(ctx, texts)
}

// Close 关闭所有客户端
func (m *Manager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for name, client := range m.clients {
		if err := client.Close(); err != nil {
			return fmt.Errorf("failed to close client %s: %w", name, err)
		}
	}

	m.clients = make(map[string]*Client)
	m.configs = make(map[string]*Config)
	m.defaultName = ""

	return nil
}

// GetDefault 获取默认客户端名称
func (m *Manager) GetDefault() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.defaultName
}

// GetConfig 获取客户端配置
func (m *Manager) GetConfig(name string) (*Config, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if name == "" {
		name = m.defaultName
	}

	config, exists := m.configs[name]
	if !exists {
		return nil, fmt.Errorf("config for client %s not found", name)
	}

	return config, nil
}
