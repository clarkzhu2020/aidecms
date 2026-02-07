package mail

import (
	"context"
	"fmt"
	"log"
	"time"

	"aidecms/config"
)

// HealthChecker 邮件服务器健康检查器
type HealthChecker struct {
	service *MultiMailService
	stopCh  chan struct{}
}

// NewHealthChecker 创建健康检查器
func NewHealthChecker(service *MultiMailService) *HealthChecker {
	return &HealthChecker{
		service: service,
		stopCh:  make(chan struct{}),
	}
}

// Start 启动健康检查
func (h *HealthChecker) Start(interval time.Duration) {
	go h.run(interval)
	log.Printf("Mail server health checker started with interval: %v", interval)
}

// Stop 停止健康检查
func (h *HealthChecker) Stop() {
	close(h.stopCh)
	log.Println("Mail server health checker stopped")
}

// run 运行健康检查循环
func (h *HealthChecker) run(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.checkAllServers()
		case <-h.stopCh:
			return
		}
	}
}

// checkAllServers 检查所有服务器健康状态
func (h *HealthChecker) checkAllServers() {
	mailConfig := h.service.GetConfig()

	for _, server := range mailConfig.Servers {
		if !server.Enabled {
			continue
		}

		// 检查服务器健康状态
		isHealthy := h.checkServerHealth(server)
		status := mailConfig.GetServerStatus()[server.ID]

		if isHealthy {
			// 服务器健康，如果之前被封禁或限流，尝试恢复
			if !status.IsHealthy || status.RateLimitHit || status.Banned {
				log.Printf("Mail server %s (%s:%d) is healthy again, unblocking", server.ID, server.Host, server.Port)
				mailConfig.UnbanServer(server.ID)
			}
		} else {
			// 服务器不健康，记录状态
			if status.IsHealthy {
				log.Printf("Mail server %s (%s:%d) is unhealthy", server.ID, server.Host, server.Port)
			}
		}
	}
}

// checkServerHealth 检查单个服务器健康状态
func (h *HealthChecker) checkServerHealth(server *config.SMTPServerConfig) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 简单的连接测试
	switch server.Encryption {
	case "tls", "starttls":
		return h.checkWithTLS(ctx, server)
	case "ssl":
		return h.checkWithSSL(ctx, server)
	default:
		return h.checkPlain(ctx, server)
	}
}

// checkWithTLS 检查TLS连接
func (h *HealthChecker) checkWithTLS(ctx context.Context, server *config.SMTPServerConfig) bool {
	// 实现连接测试逻辑
	// 实际应用中应该尝试建立TLS连接并验证
	return true // 简化实现
}

// checkWithSSL 检查SSL连接
func (h *HealthChecker) checkWithSSL(ctx context.Context, server *config.SMTPServerConfig) bool {
	// 实现连接测试逻辑
	return true // 简化实现
}

// checkPlain 检查普通连接
func (h *HealthChecker) checkPlain(ctx context.Context, server *config.SMTPServerConfig) bool {
	// 实现连接测试逻辑
	return true // 简化实现
}

// CheckServerStatus 手动检查服务器状态
func (h *HealthChecker) CheckServerStatus(serverID string) (map[string]interface{}, error) {
	mailConfig := h.service.GetConfig()
	server := mailConfig.GetServerByID(serverID)
	if server == nil {
		return nil, fmt.Errorf("server not found: %s", serverID)
	}

	status := mailConfig.GetServerStatus()[serverID]
	if status == nil {
		return nil, fmt.Errorf("server status not found: %s", serverID)
	}

	isHealthy := h.checkServerHealth(server)

	return map[string]interface{}{
		"server_id":    server.ID,
		"host":         server.Host,
		"port":         server.Port,
		"is_healthy":   isHealthy,
		"is_banned":    status.Banned,
		"fail_count":   status.FailCount,
		"total_sent":   status.TotalSent,
		"total_failed": status.TotalFailed,
		"last_error":   status.LastError,
		"enabled":      server.Enabled,
	}, nil
}

// GetAllServerStatus 获取所有服务器状态
func (h *HealthChecker) GetAllServerStatus() map[string]interface{} {
	mailConfig := h.service.GetConfig()
	serversStatus := mailConfig.GetServerStatus()

	result := make(map[string]interface{})
	for serverID, status := range serversStatus {
		server := mailConfig.GetServerByID(serverID)
		if server != nil {
			result[serverID] = map[string]interface{}{
				"host":         server.Host,
				"port":         server.Port,
				"is_healthy":   status.IsHealthy,
				"is_banned":    status.Banned,
				"rate_limit_hit": status.RateLimitHit,
				"fail_count":   status.FailCount,
				"total_sent":   status.TotalSent,
				"total_failed": status.TotalFailed,
				"last_error":   status.LastError,
				"enabled":      server.Enabled,
				"priority":     server.Priority,
				"weight":       server.Weight,
			}
		}
	}

	return result
}
