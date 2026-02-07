package stripe

import (
	"log"

	"github.com/stripe/stripe-go/v84"
)

// Config Stripe配置
type Config struct {
	SecretKey      string
	PublishableKey string
	WebhookSecret  string
	IsTest        bool
}

// Client Stripe客户端封装
type Client struct {
	*stripe.APIBackend
	config *Config
}

var (
	// 全局Stripe客户端实例
	stripeClient *Client
)

// NewClient 创建新的Stripe客户端
func NewClient(config *Config) *Client {
	stripe.Key = config.SecretKey

	// 设置Stripe配置
	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "AideCMS",
		Version: "1.0.0",
		URL:     "https://github.com/clarkzhu2020/aidecms",
	})

	return &Client{
		config: config,
	}
}

// Init 初始化Stripe客户端
func Init(secretKey, webhookSecret string, isTest bool) error {
	config := &Config{
		SecretKey:     secretKey,
		WebhookSecret: webhookSecret,
		IsTest:        isTest,
	}

	client := NewClient(config)
	stripeClient = client

	log.Printf("Stripe client initialized successfully (test mode: %v)", isTest)
	return nil
}

// GetClient 获取Stripe客户端实例
func GetClient() *Client {
	return stripeClient
}

// IsInitialized 检查客户端是否已初始化
func IsInitialized() bool {
	return stripeClient != nil
}

// GetConfig 获取Stripe配置
func (c *Client) GetConfig() *Config {
	return c.config
}

// GetSecretKey 获取密钥
func (c *Client) GetSecretKey() string {
	return c.config.SecretKey
}

// GetPublishableKey 获取公钥
func (c *Client) GetPublishableKey() string {
	return c.config.PublishableKey
}

// GetWebhookSecret 获取Webhook密钥
func (c *Client) GetWebhookSecret() string {
	return c.config.WebhookSecret
}

// IsTestMode 检查是否为测试模式
func (c *Client) IsTestMode() bool {
	return c.config.IsTest
}
