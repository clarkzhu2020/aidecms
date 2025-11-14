package config

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strconv"

	envConfig "github.com/clarkgo/clarkgo/pkg/config"
)

// MailConfig 邮件配置
type MailConfig struct {
	Driver     string `json:"driver"`     // smtp, sendmail, mailgun, postmark
	Host       string `json:"host"`       // SMTP 服务器地址
	Port       int    `json:"port"`       // SMTP 端口
	Username   string `json:"username"`   // 用户名
	Password   string `json:"password"`   // 密码
	Encryption string `json:"encryption"` // none, tls, ssl
	FromName   string `json:"from_name"`  // 发件人姓名
	FromEmail  string `json:"from_email"` // 发件人邮箱
}

// LoadMailConfig 加载邮件配置
func LoadMailConfig() (*MailConfig, error) {
	// 首先加载 .env 文件
	if err := envConfig.LoadEnv(".env"); err != nil {
		envConfig.LoadEnv(".env.example")
	}

	port, _ := strconv.Atoi(envConfig.GetEnv("MAIL_PORT", "587"))

	config := &MailConfig{
		Driver:     envConfig.GetEnv("MAIL_MAILER", "smtp"),
		Host:       envConfig.GetEnv("MAIL_HOST", "localhost"),
		Port:       port,
		Username:   envConfig.GetEnv("MAIL_USERNAME", ""),
		Password:   envConfig.GetEnv("MAIL_PASSWORD", ""),
		Encryption: envConfig.GetEnv("MAIL_ENCRYPTION", "tls"),
		FromName:   envConfig.GetEnv("MAIL_FROM_NAME", "ClarkGo"),
		FromEmail:  envConfig.GetEnv("MAIL_FROM_ADDRESS", "noreply@example.com"),
	}

	return config, nil
}

// GetSMTPAuth 获取SMTP认证
func (c *MailConfig) GetSMTPAuth() smtp.Auth {
	if c.Username == "" {
		return nil
	}
	return smtp.PlainAuth("", c.Username, c.Password, c.Host)
}

// GetTLSConfig 获取TLS配置
func (c *MailConfig) GetTLSConfig() *tls.Config {
	return &tls.Config{
		ServerName:         c.Host,
		InsecureSkipVerify: false,
	}
}

// GetAddr 获取SMTP服务器地址
func (c *MailConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
