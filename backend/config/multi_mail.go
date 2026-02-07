package config

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/smtp"
	"strconv"
	"strings"
	"sync"

	envConfig "github.com/clarkzhu2020/aidecms/pkg/config"
)

// SMTPServerConfig 单个SMTP服务器配置
type SMTPServerConfig struct {
	ID         string `json:"id"`         // 服务器唯一标识
	Host       string `json:"host"`       // SMTP服务器地址
	Port       int    `json:"port"`       // SMTP端口
	Username   string `json:"username"`   // 用户名
	Password   string `json:"password"`   // 密码
	Encryption string `json:"encryption"` // none, tls, ssl
	FromName   string `json:"from_name"` // 发件人姓名
	FromEmail  string `json:"from_email"` // 发件人邮箱
	Enabled    bool   `json:"enabled"`    // 是否启用
	Priority   int    `json:"priority"`   // 优先级，数字越小优先级越高
	Weight     int    `json:"weight"`     // 权重，用于负载均衡
}

// MailServerStatus 服务器状态
type MailServerStatus struct {
	ID            string    `json:"id"`
	IsHealthy     bool      `json:"is_healthy"`
	LastError     string    `json:"last_error,omitempty"`
	FailCount     int       `json:"fail_count"`
	LastCheckAt   int64     `json:"last_check_at"`
	LastSuccessAt int64     `json:"last_success_at"`
	TotalSent     int64     `json:"total_sent"`
	TotalFailed   int64     `json:"total_failed"`
	RateLimitHit  bool      `json:"rate_limit_hit"`
	Banned        bool      `json:"banned"`
	mu            sync.RWMutex
}

// IsAvailable 检查服务器是否可用
func (s *MailServerStatus) IsAvailable() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Enabled() && s.IsHealthy && !s.RateLimitHit && !s.Banned
}

// Enabled 检查服务器是否启用
func (s *MailServerStatus) Enabled() bool {
	return !s.Banned
}

// MarkSuccess 标记发送成功
func (s *MailServerStatus) MarkSuccess() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.IsHealthy = true
	s.TotalSent++
	s.LastSuccessAt = getCurrentTimestamp()
	s.LastCheckAt = s.LastSuccessAt
	s.RateLimitHit = false
	s.FailCount = 0
}

// MarkFailure 标记发送失败
func (s *MailServerStatus) MarkFailure(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.FailCount++
	s.TotalFailed++
	s.LastError = err.Error()
	s.LastCheckAt = getCurrentTimestamp()

	// 如果连续失败超过阈值，标记为不健康
	if s.FailCount >= 5 {
		s.IsHealthy = false
	}

	// 检查是否触发限流
	if strings.Contains(s.LastError, "rate limit") || 
	   strings.Contains(s.LastError, "too many") ||
	   strings.Contains(s.LastError, "quota exceeded") {
		s.RateLimitHit = true
	}

	// 检查是否被封禁
	if strings.Contains(s.LastError, "blocked") || 
	   strings.Contains(s.LastError, "banned") ||
	   strings.Contains(s.LastError, "suspended") {
		s.Banned = true
	}
}

// Ban 封禁服务器
func (s *MailServerStatus) Ban(reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Banned = true
	s.IsHealthy = false
	s.LastError = reason
	s.LastCheckAt = getCurrentTimestamp()
}

// Unban 解封服务器
func (s *MailServerStatus) Unban() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Banned = false
	s.IsHealthy = true
	s.FailCount = 0
	s.LastCheckAt = getCurrentTimestamp()
}

// MultiMailConfig 多邮件服务器配置
type MultiMailConfig struct {
	Driver          string                      `json:"driver"`
	Servers         []*SMTPServerConfig         `json:"servers"`
	ServerStatus    map[string]*MailServerStatus `json:"server_status"`
	FailoverEnabled bool                        `json:"failover_enabled"` // 是否启用故障切换
	LoadBalanceMode string                      `json:"load_balance_mode"` // 负载均衡模式: round_robin, weighted, priority
	CurrentIndex    int                         `json:"current_index"`    // 当前服务器索引（用于轮询）
	mu              sync.RWMutex
}

// LoadMultiMailConfig 加载多邮件服务器配置
func LoadMultiMailConfig() (*MultiMailConfig, error) {
	// 加载 .env 文件
	if err := envConfig.LoadEnv(".env"); err != nil {
		envConfig.LoadEnv(".env.example")
	}

	config := &MultiMailConfig{
		Driver:          envConfig.GetEnv("MAIL_DRIVER", "smtp"),
		FailoverEnabled:  envConfig.GetEnvBool("MAIL_FAILOVER_ENABLED", true),
		LoadBalanceMode: envConfig.GetEnv("MAIL_LOAD_BALANCE_MODE", "round_robin"),
		Servers:         make([]*SMTPServerConfig, 0),
		ServerStatus:    make(map[string]*MailServerStatus),
	}

	// 解析多个服务器配置
	serverConfigs := envConfig.GetEnv("MAIL_SERVERS", "")
	if serverConfigs == "" {
		// 兼容旧版本，使用单个服务器配置
		return loadLegacyConfig(config)
	}

	// 解析JSON格式的服务器配置
	var servers []SMTPServerConfig
	if err := json.Unmarshal([]byte(serverConfigs), &servers); err != nil {
		return nil, fmt.Errorf("failed to parse mail servers config: %w", err)
	}

	// 初始化服务器配置和状态
	for i := range servers {
		if servers[i].ID == "" {
			servers[i].ID = fmt.Sprintf("server_%d", i)
		}
		if servers[i].Priority == 0 {
			servers[i].Priority = i
		}
		if servers[i].Weight == 0 {
			servers[i].Weight = 1
		}
		servers[i].Enabled = true

		config.Servers = append(config.Servers, &servers[i])
		config.ServerStatus[servers[i].ID] = &MailServerStatus{
			ID:        servers[i].ID,
			IsHealthy: true,
		}
	}

	return config, nil
}

// loadLegacyConfig 加载旧版本的单服务器配置
func loadLegacyConfig(config *MultiMailConfig) (*MultiMailConfig, error) {
	port, _ := strconv.Atoi(envConfig.GetEnv("MAIL_PORT", "587"))

	server := &SMTPServerConfig{
		ID:         "default",
		Host:       envConfig.GetEnv("MAIL_HOST", "localhost"),
		Port:       port,
		Username:   envConfig.GetEnv("MAIL_USERNAME", ""),
		Password:   envConfig.GetEnv("MAIL_PASSWORD", ""),
		Encryption: envConfig.GetEnv("MAIL_ENCRYPTION", "tls"),
		FromName:   envConfig.GetEnv("MAIL_FROM_NAME", "AideCMS"),
		FromEmail:  envConfig.GetEnv("MAIL_FROM_ADDRESS", "noreply@example.com"),
		Enabled:    true,
		Priority:   0,
		Weight:     1,
	}

	config.Servers = append(config.Servers, server)
	config.ServerStatus[server.ID] = &MailServerStatus{
		ID:        server.ID,
		IsHealthy: true,
	}

	return config, nil
}

// GetNextServer 获取下一个可用的服务器
func (c *MultiMailConfig) GetNextServer() *SMTPServerConfig {
	c.mu.Lock()
	defer c.mu.Unlock()

	availableServers := c.getAvailableServers()
	if len(availableServers) == 0 {
		return nil
	}

	var server *SMTPServerConfig

	switch c.LoadBalanceMode {
	case "priority":
		server = c.getNextServerByPriority(availableServers)
	case "weighted":
		server = c.getNextServerByWeight(availableServers)
	case "round_robin":
		fallthrough
	default:
		server = c.getNextServerByRoundRobin(availableServers)
	}

	return server
}

// getAvailableServers 获取所有可用的服务器
func (c *MultiMailConfig) getAvailableServers() []*SMTPServerConfig {
	available := make([]*SMTPServerConfig, 0)
	for _, server := range c.Servers {
		if status, ok := c.ServerStatus[server.ID]; ok {
			if status.IsAvailable() && server.Enabled {
				available = append(available, server)
			}
		}
	}
	return available
}

// getNextServerByRoundRobin 按轮询获取服务器
func (c *MultiMailConfig) getNextServerByRoundRobin(servers []*SMTPServerConfig) *SMTPServerConfig {
	if len(servers) == 0 {
		return nil
	}

	// 找到当前索引对应的服务器
	for i, server := range servers {
		if server.ID == c.getServerByIndex(c.CurrentIndex).ID {
			c.CurrentIndex = (i + 1) % len(servers)
			return server
		}
	}

	c.CurrentIndex = 0
	return servers[0]
}

// getServerByIndex 根据索引获取服务器
func (c *MultiMailConfig) getServerByIndex(index int) *SMTPServerConfig {
	if len(c.Servers) == 0 {
		return nil
	}
	return c.Servers[index%len(c.Servers)]
}

// getNextServerByPriority 按优先级获取服务器
func (c *MultiMailConfig) getNextServerByPriority(servers []*SMTPServerConfig) *SMTPServerConfig {
	if len(servers) == 0 {
		return nil
	}

	// 按优先级排序
	priority := servers[0].Priority
	for _, server := range servers {
		if server.Priority < priority {
			priority = server.Priority
		}
	}

	// 返回最高优先级的第一个可用服务器
	for _, server := range servers {
		if server.Priority == priority {
			return server
		}
	}

	return servers[0]
}

// getNextServerByWeight 按权重获取服务器
func (c *MultiMailConfig) getNextServerByWeight(servers []*SMTPServerConfig) *SMTPServerConfig {
	if len(servers) == 0 {
		return nil
	}

	// 计算总权重
	totalWeight := 0
	for _, server := range servers {
		totalWeight += server.Weight
	}

	// 简单实现：轮询时按权重分配
	// TODO: 实现更精确的加权轮询算法
	for _, server := range servers {
		if server.Weight > 0 {
			return server
		}
	}

	return servers[0]
}

// GetServerStatus 获取所有服务器状态
func (c *MultiMailConfig) GetServerStatus() map[string]*MailServerStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	statusMap := make(map[string]*MailServerStatus)
	for id, status := range c.ServerStatus {
		statusCopy := *status
		statusMap[id] = &statusCopy
	}
	return statusMap
}

// UpdateServerStatus 更新服务器状态
func (c *MultiMailConfig) UpdateServerStatus(serverID string, success bool, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if status, ok := c.ServerStatus[serverID]; ok {
		if success {
			status.MarkSuccess()
		} else {
			status.MarkFailure(err)
		}
	}
}

// BanServer 封禁服务器
func (c *MultiMailConfig) BanServer(serverID, reason string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if status, ok := c.ServerStatus[serverID]; ok {
		status.Ban(reason)
	}
}

// UnbanServer 解封服务器
func (c *MultiMailConfig) UnbanServer(serverID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if status, ok := c.ServerStatus[serverID]; ok {
		status.Unban()
	}
}

// GetServerByID 根据ID获取服务器配置
func (c *MultiMailConfig) GetServerByID(id string) *SMTPServerConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, server := range c.Servers {
		if server.ID == id {
			return server
		}
	}
	return nil
}

// GetSMTPAuth 获取SMTP认证
func (s *SMTPServerConfig) GetSMTPAuth() smtp.Auth {
	if s.Username == "" {
		return nil
	}
	return smtp.PlainAuth("", s.Username, s.Password, s.Host)
}

// GetTLSConfig 获取TLS配置
func (s *SMTPServerConfig) GetTLSConfig() *tls.Config {
	return &tls.Config{
		ServerName:         s.Host,
		InsecureSkipVerify: false,
	}
}

// GetAddr 获取SMTP服务器地址
func (s *SMTPServerConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// getCurrentTimestamp 获取当前时间戳（秒）
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}
