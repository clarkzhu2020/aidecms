package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/clarkzhu2020/aidecms/pkg/mail"
	"github.com/cloudwego/hertz/pkg/app"
)

// MultiMailController 多邮件服务器控制器
type MultiMailController struct {
	mailService   *mail.MultiMailService
	healthChecker *mail.HealthChecker
}

// NewMultiMailController 创建多邮件服务器控制器
func NewMultiMailController() (*MultiMailController, error) {
	mailService, err := mail.NewMultiMailService()
	if err != nil {
		return nil, fmt.Errorf("failed to create multi mail service: %w", err)
	}

	healthChecker := mail.NewHealthChecker(mailService)

	return &MultiMailController{
		mailService:   mailService,
		healthChecker: healthChecker,
	}, nil
}

// SendMail 发送邮件
func (c *MultiMailController) SendMail(ctx context.Context, hCtx *app.RequestContext) {
	var req SendMailRequest
	if err := hCtx.BindJSON(&req); err != nil {
		hCtx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request format",
			"message": err.Error(),
		})
		return
	}

	// 验证邮箱地址
	for _, email := range req.To {
		if err := c.mailService.ValidateEmail(email); err != nil {
			hCtx.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "Invalid email address",
				"message": fmt.Sprintf("Invalid email: %s", email),
			})
			return
		}
	}

	// 构建邮件对象
	mailObj := &mail.Mail{
		To:       req.To,
		Cc:       req.Cc,
		Bcc:      req.Bcc,
		Subject:  req.Subject,
		Body:     req.Body,
		HTMLBody: req.HTMLBody,
		Headers:  req.Headers,
	}

	// 处理附件
	if len(req.Attachments) > 0 {
		attachments := make([]mail.Attachment, 0, len(req.Attachments))
		for _, att := range req.Attachments {
			content, err := base64.StdEncoding.DecodeString(att.Content)
			if err != nil {
				hCtx.JSON(http.StatusBadRequest, map[string]interface{}{
					"error":   "Invalid attachment content",
					"message": fmt.Sprintf("Failed to decode attachment %s: %v", att.Filename, err),
				})
				return
			}

			mimeType := att.MimeType
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}

			attachments = append(attachments, mail.Attachment{
				Filename: att.Filename,
				Content:  content,
				MimeType: mimeType,
			})
		}
		mailObj.Attachments = attachments
	}

	// 发送邮件（自动故障切换）
	if err := c.mailService.SendMail(mailObj); err != nil {
		hCtx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to send email",
			"message": err.Error(),
		})
		return
	}

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Email sent successfully",
		"data": map[string]interface{}{
			"to":      req.To,
			"subject": req.Subject,
		},
	})
}

// SendBulkMail 批量发送邮件（支持负载均衡和故障切换）
func (c *MultiMailController) SendBulkMail(ctx context.Context, hCtx *app.RequestContext) {
	var req struct {
		Emails   []SendMailRequest `json:"emails" binding:"required"`
		Distribute string          `json:"distribute"` // round_robin, random, priority
	}

	if err := hCtx.BindJSON(&req); err != nil {
		hCtx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request format",
			"message": err.Error(),
		})
		return
	}

	// 构建邮件列表
	emails := make([]mail.Mail, len(req.Emails))
	for i, emailReq := range req.Emails {
		// 验证邮箱地址
		for _, email := range emailReq.To {
			if err := c.mailService.ValidateEmail(email); err != nil {
				hCtx.JSON(http.StatusBadRequest, map[string]interface{}{
					"error":   "Invalid email address",
					"message": fmt.Sprintf("Invalid email: %s", email),
				})
				return
			}
		}

		emails[i] = mail.Mail{
			To:       emailReq.To,
			Cc:       emailReq.Cc,
			Bcc:      emailReq.Bcc,
			Subject:  emailReq.Subject,
			Body:     emailReq.Body,
			HTMLBody: emailReq.HTMLBody,
			Headers:  emailReq.Headers,
		}
	}

	// 批量发送（自动负载均衡）
	results, err := c.mailService.SendBulkMail(emails)
	if err != nil {
		hCtx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to send bulk emails",
			"message": err.Error(),
		})
		return
	}

	// 统计结果
	successCount := 0
	failureCount := 0
	serverStats := make(map[string]int)

	for _, result := range results {
		if result.Success {
			successCount++
			serverStats[result.ServerID]++
		} else {
			failureCount++
		}
	}

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Bulk email sending completed. Success: %d, Failed: %d", successCount, failureCount),
		"data": map[string]interface{}{
			"total":         len(emails),
			"success_count": successCount,
			"failure_count": failureCount,
			"server_stats":  serverStats,
			"results":       results,
		},
	})
}

// GetServerStatus 获取所有服务器状态
func (c *MultiMailController) GetServerStatus(ctx context.Context, hCtx *app.RequestContext) {
	statusMap := c.healthChecker.GetAllServerStatus()

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Server status retrieved successfully",
		"data":    statusMap,
	})
}

// CheckServerHealth 检查指定服务器健康状态
func (c *MultiMailController) CheckServerHealth(ctx context.Context, hCtx *app.RequestContext) {
	serverID := hCtx.Query("server_id")
	if serverID == "" {
		hCtx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Missing server_id parameter",
			"message": "Please provide server_id parameter",
		})
		return
	}

	status, err := c.healthChecker.CheckServerStatus(serverID)
	if err != nil {
		hCtx.JSON(http.StatusNotFound, map[string]interface{}{
			"error":   "Server not found",
			"message": err.Error(),
		})
		return
	}

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Server health check completed",
		"data":    status,
	})
}

// BanServer 封禁服务器
func (c *MultiMailController) BanServer(ctx context.Context, hCtx *app.RequestContext) {
	serverID := hCtx.Query("server_id")
	if serverID == "" {
		hCtx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Missing server_id parameter",
			"message": "Please provide server_id parameter",
		})
		return
	}

	reason := hCtx.Query("reason")
	if reason == "" {
		reason = "Manually banned via API"
	}

	c.mailService.BanServer(serverID, reason)

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Server %s has been banned", serverID),
		"data": map[string]interface{}{
			"server_id": serverID,
			"reason":    reason,
		},
	})
}

// UnbanServer 解封服务器
func (c *MultiMailController) UnbanServer(ctx context.Context, hCtx *app.RequestContext) {
	serverID := hCtx.Query("server_id")
	if serverID == "" {
		hCtx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Missing server_id parameter",
			"message": "Please provide server_id parameter",
		})
		return
	}

	c.mailService.UnbanServer(serverID)

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Server %s has been unbanned", serverID),
		"data": map[string]interface{}{
			"server_id": serverID,
		},
	})
}

// GetMailConfig 获取邮件配置信息
func (c *MultiMailController) GetMailConfig(ctx context.Context, hCtx *app.RequestContext) {
	config := c.mailService.GetConfig()

	// 隐藏敏感信息
	safeServers := make([]map[string]interface{}, len(config.Servers))
	for i, server := range config.Servers {
		safeServers[i] = map[string]interface{}{
			"id":         server.ID,
			"host":       server.Host,
			"port":       server.Port,
			"encryption": server.Encryption,
			"from_name":  server.FromName,
			"from_email": server.FromEmail,
			"enabled":    server.Enabled,
			"priority":   server.Priority,
			"weight":     server.Weight,
			"username":   maskString(server.Username),
		}
	}

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Mail config retrieved successfully",
		"data": map[string]interface{}{
			"driver":            config.Driver,
			"failover_enabled":  config.FailoverEnabled,
			"load_balance_mode": config.LoadBalanceMode,
			"servers":           safeServers,
		},
	})
}

// TestConnection 测试所有服务器连接
func (c *MultiMailController) TestConnection(ctx context.Context, hCtx *app.RequestContext) {
	config := c.mailService.GetConfig()
	results := make([]map[string]interface{}, len(config.Servers))

	for i, server := range config.Servers {
		status := c.healthChecker.CheckServerHealth(server.ID)
		results[i] = map[string]interface{}{
			"server_id": server.ID,
			"host":      server.Host,
			"port":      server.Port,
			"status":    status,
		}
	}

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Server connection test completed",
		"data": map[string]interface{}{
			"total_servers": len(config.Servers),
			"results":      results,
		},
	})
}

// maskString 掩码字符串，隐藏敏感信息
func maskString(s string) string {
	if len(s) <= 3 {
		return "***"
	}
	return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
}
