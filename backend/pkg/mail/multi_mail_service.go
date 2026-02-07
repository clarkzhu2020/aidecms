package mail

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"aidecms/config"
)

// MultiMailService 多邮件服务器服务
type MultiMailService struct {
	config *config.MultiMailConfig
}

// NewMultiMailService 创建多邮件服务器服务
func NewMultiMailService() (*MultiMailService, error) {
	mailConfig, err := config.LoadMultiMailConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load mail config: %w", err)
	}

	return &MultiMailService{
		config: mailConfig,
	}, nil
}

// SendMail 发送邮件（自动故障切换）
func (s *MultiMailService) SendMail(mail *Mail) error {
	var lastError error
	var triedServers []string

	// 尝试所有可用的服务器
	for i := 0; i < s.getMaxRetries(); i++ {
		server := s.config.GetNextServer()
		if server == nil {
			return fmt.Errorf("no available mail servers")
		}

		// 避免重复尝试同一服务器
		if contains(triedServers, server.ID) && len(triedServers) < len(s.config.Servers) {
			continue
		}
		triedServers = append(triedServers, server.ID)

		// 尝试发送邮件
		err := s.sendWithServer(mail, server)
		if err == nil {
			// 发送成功，更新服务器状态
			s.config.UpdateServerStatus(server.ID, true, nil)
			return nil
		}

		lastError = err
		s.config.UpdateServerStatus(server.ID, false, err)

		// 如果启用了故障切换，继续尝试下一个服务器
		if s.config.FailoverEnabled {
			continue
		} else {
			// 未启用故障切换，直接返回错误
			return err
		}
	}

	return fmt.Errorf("failed to send email after trying %d servers: %w", len(triedServers), lastError)
}

// sendWithServer 使用指定服务器发送邮件
func (s *MultiMailService) sendWithServer(mail *Mail, server *config.SMTPServerConfig) error {
	// 构建邮件内容
	message, err := s.buildMessage(mail, server)
	if err != nil {
		return fmt.Errorf("failed to build message: %w", err)
	}

	// 获取认证
	auth := server.GetSMTPAuth()

	// 获取所有收件人
	recipients := append(mail.To, mail.Cc...)
	recipients = append(recipients, mail.Bcc...)

	// 根据加密类型发送
	switch server.Encryption {
	case "tls", "starttls":
		return s.sendWithTLS(server, auth, recipients, message)
	case "ssl":
		return s.sendWithSSL(server, auth, recipients, message)
	default:
		return s.sendPlain(server, auth, recipients, message)
	}
}

// sendWithTLS 使用TLS发送
func (s *MultiMailService) sendWithTLS(server *config.SMTPServerConfig, auth smtp.Auth, recipients []string, message []byte) error {
	// 连接到服务器
	client, err := smtp.Dial(server.GetAddr())
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server %s: %w", server.ID, err)
	}
	defer client.Close()

	// 启动TLS
	if err = client.StartTLS(server.GetTLSConfig()); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// 认证
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	// 发送邮件
	if err = client.Mail(server.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	for _, recipient := range recipients {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data writer: %w", err)
	}
	defer w.Close()

	if _, err = w.Write(message); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// sendWithSSL 使用SSL发送
func (s *MultiMailService) sendWithSSL(server *config.SMTPServerConfig, auth smtp.Auth, recipients []string, message []byte) error {
	// 创建TLS连接
	tlsConfig := server.GetTLSConfig()
	conn, err := tls.Dial("tcp", server.GetAddr(), tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server with SSL %s: %w", server.ID, err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, server.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// 认证
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	// 发送邮件
	if err = client.Mail(server.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	for _, recipient := range recipients {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data writer: %w", err)
	}
	defer w.Close()

	if _, err = w.Write(message); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// sendPlain 普通发送
func (s *MultiMailService) sendPlain(server *config.SMTPServerConfig, auth smtp.Auth, recipients []string, message []byte) error {
	return smtp.SendMail(server.GetAddr(), auth, server.FromEmail, recipients, message)
}

// buildMessage 构建邮件消息
func (s *MultiMailService) buildMessage(mail *Mail, server *config.SMTPServerConfig) ([]byte, error) {
	var buf bytes.Buffer

	// 添加基本头部
	buf.WriteString(fmt.Sprintf("From: %s <%s>\r\n", server.FromName, server.FromEmail))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ", ")))

	if len(mail.Cc) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(mail.Cc, ", ")))
	}

	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", mail.Subject))
	buf.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	buf.WriteString("MIME-Version: 1.0\r\n")

	// 添加自定义头部
	for key, value := range mail.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// 如果有附件或HTML内容，使用multipart
	if len(mail.Attachments) > 0 || mail.HTMLBody != "" {
		boundary := "boundary_" + fmt.Sprintf("%d", time.Now().Unix())
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n\r\n", boundary))

		// 添加文本内容
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		if mail.HTMLBody != "" {
			altBoundary := "alt_boundary_" + fmt.Sprintf("%d", time.Now().Unix())
			buf.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n\r\n", altBoundary))

			// 纯文本版本
			buf.WriteString(fmt.Sprintf("--%s\r\n", altBoundary))
			buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
			buf.WriteString(mail.Body)
			buf.WriteString("\r\n\r\n")

			// HTML版本
			buf.WriteString(fmt.Sprintf("--%s\r\n", altBoundary))
			buf.WriteString("Content-Type: text/html; charset=utf-8\r\n\r\n")
			buf.WriteString(mail.HTMLBody)
			buf.WriteString("\r\n\r\n")
			buf.WriteString(fmt.Sprintf("--%s--\r\n", altBoundary))
		} else {
			buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
			buf.WriteString(mail.Body)
			buf.WriteString("\r\n\r\n")
		}

		// 添加附件
		for _, attachment := range mail.Attachments {
			buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\r\n", attachment.MimeType))
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\r\n\r\n", attachment.Filename))

			// Base64编码附件内容
			encoded := base64.StdEncoding.EncodeToString(attachment.Content)
			buf.WriteString(encoded)
			buf.WriteString("\r\n\r\n")
		}

		buf.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else {
		// 简单文本邮件
		buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
		buf.WriteString(mail.Body)
		buf.WriteString("\r\n")
	}

	return buf.Bytes(), nil
}

// SendTemplate 发送模板邮件
func (s *MultiMailService) SendTemplate(templateName string, data interface{}, mail *Mail) error {
	// 加载模板
	tmpl, err := template.ParseFiles(fmt.Sprintf("templates/email/%s.html", templateName))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// 渲染模板
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	mail.HTMLBody = buf.String()
	return s.SendMail(mail)
}

// ValidateEmail 验证邮箱地址
func (s *MultiMailService) ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

// GetConfig 获取邮件配置
func (s *MultiMailService) GetConfig() *config.MultiMailConfig {
	return s.config
}

// GetServerStatus 获取所有服务器状态
func (s *MultiMailService) GetServerStatus() map[string]*config.MailServerStatus {
	return s.config.GetServerStatus()
}

// BanServer 封禁服务器
func (s *MultiMailService) BanServer(serverID, reason string) {
	s.config.BanServer(serverID, reason)
}

// UnbanServer 解封服务器
func (s *MultiMailService) UnbanServer(serverID string) {
	s.config.UnbanServer(serverID)
}

// getMaxRetries 获取最大重试次数
func (s *MultiMailService) getMaxRetries() int {
	if s.config.FailoverEnabled {
		return len(s.config.Servers) * 2
	}
	return 1
}

// contains 检查字符串是否包含在切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// SendBulkMail 批量发送邮件（自动负载均衡）
func (s *MultiMailService) SendBulkMail(emails []Mail) ([]SendResult, error) {
	results := make([]SendResult, len(emails))

	// 为每封邮件分配服务器
	for i, email := range emails {
		server := s.config.GetNextServer()
		if server == nil {
			results[i] = SendResult{
				Index:   i,
				Success: false,
				Error:   fmt.Errorf("no available mail servers"),
			}
			continue
		}

		// 尝试发送
		err := s.sendWithServer(&email, server)
		if err != nil {
			results[i] = SendResult{
				Index:   i,
				Success: false,
				ServerID: server.ID,
				Error:   err,
			}
			s.config.UpdateServerStatus(server.ID, false, err)
		} else {
			results[i] = SendResult{
				Index:    i,
				Success:  true,
				ServerID: server.ID,
				To:       email.To,
			}
			s.config.UpdateServerStatus(server.ID, true, nil)
		}
	}

	return results, nil
}

// SendResult 发送结果
type SendResult struct {
	Index    int
	Success  bool
	ServerID string
	To       []string
	Error    error
}
