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

	"github.com/clarkgo/clarkgo/config"
)

// MailService 邮件服务
type MailService struct {
	config *config.MailConfig
}

// NewMailService 创建邮件服务
func NewMailService() (*MailService, error) {
	mailConfig, err := config.LoadMailConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load mail config: %w", err)
	}

	return &MailService{
		config: mailConfig,
	}, nil
}

// Mail 邮件结构
type Mail struct {
	To          []string          `json:"to"`
	Cc          []string          `json:"cc,omitempty"`
	Bcc         []string          `json:"bcc,omitempty"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	HTMLBody    string            `json:"html_body,omitempty"`
	Attachments []Attachment      `json:"attachments,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

// Attachment 附件
type Attachment struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
	MimeType string `json:"mime_type"`
}

// SendMail 发送邮件
func (s *MailService) SendMail(mail *Mail) error {
	switch s.config.Driver {
	case "smtp":
		return s.sendSMTP(mail)
	default:
		return fmt.Errorf("unsupported mail driver: %s", s.config.Driver)
	}
}

// sendSMTP 通过SMTP发送邮件
func (s *MailService) sendSMTP(mail *Mail) error {
	// 构建邮件内容
	message, err := s.buildMessage(mail)
	if err != nil {
		return fmt.Errorf("failed to build message: %w", err)
	}

	// 获取认证
	auth := s.config.GetSMTPAuth()

	// 获取所有收件人
	recipients := append(mail.To, mail.Cc...)
	recipients = append(recipients, mail.Bcc...)

	// 根据加密类型发送
	switch s.config.Encryption {
	case "tls", "starttls":
		return s.sendWithTLS(auth, recipients, message)
	case "ssl":
		return s.sendWithSSL(auth, recipients, message)
	default:
		return s.sendPlain(auth, recipients, message)
	}
}

// sendWithTLS 使用TLS发送
func (s *MailService) sendWithTLS(auth smtp.Auth, recipients []string, message []byte) error {
	// 连接到服务器
	client, err := smtp.Dial(s.config.GetAddr())
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	// 启动TLS
	if err = client.StartTLS(s.config.GetTLSConfig()); err != nil {
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// 认证
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	// 发送邮件
	if err = client.Mail(s.config.FromEmail); err != nil {
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
func (s *MailService) sendWithSSL(auth smtp.Auth, recipients []string, message []byte) error {
	// 创建TLS连接
	tlsConfig := s.config.GetTLSConfig()
	conn, err := tls.Dial("tcp", s.config.GetAddr(), tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server with SSL: %w", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, s.config.Host)
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
	if err = client.Mail(s.config.FromEmail); err != nil {
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
func (s *MailService) sendPlain(auth smtp.Auth, recipients []string, message []byte) error {
	return smtp.SendMail(s.config.GetAddr(), auth, s.config.FromEmail, recipients, message)
}

// buildMessage 构建邮件消息
func (s *MailService) buildMessage(mail *Mail) ([]byte, error) {
	var buf bytes.Buffer

	// 添加基本头部
	buf.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.config.FromName, s.config.FromEmail))
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
			buf.WriteString("Content-Type: multipart/alternative; boundary=alt_boundary\r\n\r\n")

			// 纯文本版本
			buf.WriteString("--alt_boundary\r\n")
			buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
			buf.WriteString(mail.Body)
			buf.WriteString("\r\n\r\n")

			// HTML版本
			buf.WriteString("--alt_boundary\r\n")
			buf.WriteString("Content-Type: text/html; charset=utf-8\r\n\r\n")
			buf.WriteString(mail.HTMLBody)
			buf.WriteString("\r\n\r\n")
			buf.WriteString("--alt_boundary--\r\n")
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
func (s *MailService) SendTemplate(templateName string, data interface{}, mail *Mail) error {
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
func (s *MailService) ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

// GetConfig 获取邮件配置
func (s *MailService) GetConfig() (*config.MailConfig, error) {
	return s.config, nil
}
