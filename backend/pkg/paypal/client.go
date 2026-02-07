package paypal

import (
	"log"
	"os"

	"github.com/plutov/paypal/v4"
)

// Config PayPal配置
type Config struct {
	ClientID     string
	ClientSecret string
	IsSandbox    bool
}

// Client PayPal客户端
type Client struct {
	*paypal.Client
	config *Config
}

var (
	// 全局PayPal客户端实例
	paypalClient *Client
)

// NewClient 创建新的PayPal客户端
func NewClient(config *Config) (*Client, error) {
	var apiBase string
	if config.IsSandbox {
		apiBase = paypal.APIBaseSandBox
	} else {
		apiBase = paypal.APIBaseLive
	}

	client, err := paypal.NewClient(config.ClientID, config.ClientSecret, apiBase)
	if err != nil {
		return nil, err
	}

	// 设置日志输出
	client.SetLog(os.Stdout)

	return &Client{
		Client: client,
		config: config,
	}, nil
}

// Init 初始化PayPal客户端
func Init(clientID, clientSecret string, isSandbox bool) error {
	config := &Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		IsSandbox:    isSandbox,
	}

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	paypalClient = client
	log.Println("PayPal client initialized successfully")
	return nil
}

// GetClient 获取PayPal客户端实例
func GetClient() *Client {
	return paypalClient
}

// IsInitialized 检查客户端是否已初始化
func IsInitialized() bool {
	return paypalClient != nil
}
