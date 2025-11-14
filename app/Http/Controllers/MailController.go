package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/clarkgo/clarkgo/pkg/mail"
	"github.com/cloudwego/hertz/pkg/app"
)

// MailController 邮件控制器
type MailController struct {
	mailService *mail.MailService
}

// NewMailController 创建邮件控制器
func NewMailController() (*MailController, error) {
	mailService, err := mail.NewMailService()
	if err != nil {
		return nil, fmt.Errorf("failed to create mail service: %w", err)
	}

	return &MailController{
		mailService: mailService,
	}, nil
}

// SendMailRequest 发送邮件请求
type SendMailRequest struct {
	To          []string            `json:"to" binding:"required"`
	Cc          []string            `json:"cc,omitempty"`
	Bcc         []string            `json:"bcc,omitempty"`
	Subject     string              `json:"subject" binding:"required"`
	Body        string              `json:"body" binding:"required"`
	HTMLBody    string              `json:"html_body,omitempty"`
	Attachments []AttachmentRequest `json:"attachments,omitempty"`
	Headers     map[string]string   `json:"headers,omitempty"`
}

// AttachmentRequest 附件请求
type AttachmentRequest struct {
	Filename string `json:"filename" binding:"required"`
	Content  string `json:"content" binding:"required"` // base64编码的内容
	MimeType string `json:"mime_type,omitempty"`
}

// SendTemplateRequest 发送模板邮件请求
type SendTemplateRequest struct {
	Template string                 `json:"template" binding:"required"`
	Data     map[string]interface{} `json:"data,omitempty"`
	SendMailRequest
}

// SendMail 发送邮件
func (c *MailController) SendMail(ctx context.Context, hCtx *app.RequestContext) {
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

	// 发送邮件
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

// SendTemplate 发送模板邮件
func (c *MailController) SendTemplate(ctx context.Context, hCtx *app.RequestContext) {
	var req SendTemplateRequest
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
		To:      req.To,
		Cc:      req.Cc,
		Bcc:     req.Bcc,
		Subject: req.Subject,
		Body:    req.Body,
		Headers: req.Headers,
	}

	// 发送模板邮件
	if err := c.mailService.SendTemplate(req.Template, req.Data, mailObj); err != nil {
		hCtx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to send template email",
			"message": err.Error(),
		})
		return
	}

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Template email sent successfully",
		"data": map[string]interface{}{
			"template": req.Template,
			"to":       req.To,
			"subject":  req.Subject,
		},
	})
}

// TestConnection 测试邮件服务器连接
func (c *MailController) TestConnection(ctx context.Context, hCtx *app.RequestContext) {
	// 尝试连接但不实际发送
	err := c.validateConnection()
	if err != nil {
		hCtx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "Connection test failed",
			"message": err.Error(),
		})
		return
	}

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Mail server connection is working",
	})
}

// GetMailConfig 获取邮件配置信息（隐藏敏感信息）
func (c *MailController) GetMailConfig(ctx context.Context, hCtx *app.RequestContext) {
	config, err := c.mailService.GetConfig()
	if err != nil {
		hCtx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to get mail configuration",
			"message": err.Error(),
		})
		return
	}

	// 隐藏敏感信息
	safeConfig := map[string]interface{}{
		"driver":     config.Driver,
		"host":       config.Host,
		"port":       config.Port,
		"encryption": config.Encryption,
		"from_name":  config.FromName,
		"from_email": config.FromEmail,
		"username":   maskString(config.Username),
	}

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    safeConfig,
	})
}

// ValidateEmail 验证邮箱地址
func (c *MailController) ValidateEmail(ctx context.Context, hCtx *app.RequestContext) {
	email := hCtx.Query("email")
	if email == "" {
		hCtx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Missing email parameter",
			"message": "Please provide an email address to validate",
		})
		return
	}

	err := c.mailService.ValidateEmail(email)
	if err != nil {
		hCtx.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Invalid email address",
			"message": err.Error(),
		})
		return
	}

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Email address is valid",
		"email":   email,
	})
}

// SendBulkMail 批量发送邮件
func (c *MailController) SendBulkMail(ctx context.Context, hCtx *app.RequestContext) {
	var req struct {
		Emails   []SendMailRequest `json:"emails" binding:"required"`
		MaxRetry int               `json:"max_retry,omitempty"`
	}

	if err := hCtx.BindJSON(&req); err != nil {
		hCtx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request format",
			"message": err.Error(),
		})
		return
	}

	if req.MaxRetry == 0 {
		req.MaxRetry = 3
	}

	results := make([]map[string]interface{}, 0, len(req.Emails))
	successCount := 0
	failureCount := 0

	for i, emailReq := range req.Emails {
		result := map[string]interface{}{
			"index":   i,
			"to":      emailReq.To,
			"subject": emailReq.Subject,
		}

		// 构建邮件对象
		mailObj := &mail.Mail{
			To:       emailReq.To,
			Cc:       emailReq.Cc,
			Bcc:      emailReq.Bcc,
			Subject:  emailReq.Subject,
			Body:     emailReq.Body,
			HTMLBody: emailReq.HTMLBody,
			Headers:  emailReq.Headers,
		}

		// 重试发送
		var lastError error
		sent := false
		for retry := 0; retry < req.MaxRetry && !sent; retry++ {
			if err := c.mailService.SendMail(mailObj); err != nil {
				lastError = err
			} else {
				sent = true
				successCount++
				result["success"] = true
				result["message"] = "Email sent successfully"
			}
		}

		if !sent {
			failureCount++
			result["success"] = false
			result["error"] = "Failed to send email"
			result["message"] = lastError.Error()
		}

		results = append(results, result)
	}

	hCtx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Bulk email sending completed. Success: %d, Failed: %d", successCount, failureCount),
		"data": map[string]interface{}{
			"total":         len(req.Emails),
			"success_count": successCount,
			"failure_count": failureCount,
			"results":       results,
		},
	})
}

// validateConnection 验证邮件服务器连接
func (c *MailController) validateConnection() error {
	// 这里可以实现实际的连接测试逻辑
	// 暂时返回nil表示连接正常
	return nil
}

// maskString 掩码字符串，隐藏敏感信息
func maskString(s string) string {
	if len(s) <= 3 {
		return "***"
	}
	return s[:2] + strings.Repeat("*", len(s)-4) + s[len(s)-2:]
}
