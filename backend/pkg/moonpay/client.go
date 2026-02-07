package moonpay

import (
	"log"
	"os"
)

// Config MoonPay配置
type Config struct {
	APIKey       string
	SecretKey    string
	WebhookKey   string
	IsSandbox    bool
	BaseURL      string
}

// Client MoonPay客户端
type Client struct {
	httpClient *HTTPClient
	config     *Config
}

var (
	// 全局MoonPay客户端实例
	moonpayClient *Client
)

// NewClient 创建新的MoonPay客户端
func NewClient(config *Config) (*Client, error) {
	// 设置基础URL
	baseURL := config.BaseURL
	if baseURL == "" {
		if config.IsSandbox {
			baseURL = "https://api-sandbox.moonpay.com"
		} else {
			baseURL = "https://api.moonpay.com"
		}
	}

	// 创建HTTP客户端
	httpClient := NewHTTPClient(config.APIKey, config.SecretKey, baseURL)

	return &Client{
		httpClient: httpClient,
		config:     config,
	}, nil
}

// Init 初始化MoonPay客户端
func Init(apiKey, secretKey, webhookKey string, isSandbox bool) error {
	config := &Config{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		WebhookKey: webhookKey,
		IsSandbox:  isSandbox,
	}

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	moonpayClient = client
	log.Println("MoonPay client initialized successfully")
	return nil
}

// GetClient 获取MoonPay客户端实例
func GetClient() *Client {
	return moonpayClient
}

// GetConfig 获取MoonPay配置
func (c *Client) GetConfig() *Config {
	return c.config
}

// IsInitialized 检查客户端是否已初始化
func IsInitialized() bool {
	return moonpayClient != nil
}
