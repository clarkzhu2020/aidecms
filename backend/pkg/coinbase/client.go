package coinbase

import (
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// Config Coinbase配置
type Config struct {
	APIKey      string
	APISecret   string
	WebhookKey  string
	IsSandbox   bool
	AccountID   string
}

// Client Coinbase客户端
type Client struct {
	httpClient *HTTPClient
	config     *Config
}

var (
	// 全局Coinbase客户端实例
	coinbaseClient *Client
)

// NewClient 创建新的Coinbase客户端
func NewClient(config *Config) (*Client, error) {
	// 设置基础URL
	baseURL := "https://business.coinbase.com"
	if config.IsSandbox {
		baseURL = "https://business-sandbox.coinbase.com"
	}

	// 创建HTTP客户端
	httpClient := NewHTTPClient(config.APIKey, config.APISecret, baseURL)

	return &Client{
		httpClient: httpClient,
		config:     config,
	}, nil
}

// Init 初始化Coinbase客户端
func Init(apiKey, apiSecret, webhookKey, accountID string, isSandbox bool) error {
	config := &Config{
		APIKey:     apiKey,
		APISecret:  apiSecret,
		WebhookKey: webhookKey,
		AccountID:  accountID,
		IsSandbox:  isSandbox,
	}

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	coinbaseClient = client
	log.Println("Coinbase client initialized successfully")
	return nil
}

// GetClient 获取Coinbase客户端实例
func GetClient() *Client {
	return coinbaseClient
}

// GetConfig 获取Coinbase配置
func (c *Client) GetConfig() *Config {
	return c.config
}

// IsInitialized 检查客户端是否已初始化
func IsInitialized() bool {
	return coinbaseClient != nil
}

// GenerateJWT 生成JWT Token用于API认证
func GenerateJWT(apiKey, apiSecret string) (string, error) {
	// 创建JWT claims
	claims := jwt.MapClaims{
		"sub":   apiKey,
		"iss":   apiKey,
		"aud":   "https://business.coinbase.com",
		"exp":   60,
		"nbf":   0,
		"iat":   0,
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用API Secret作为签名密钥
	tokenString, err := token.SignedString([]byte(apiSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateJWTWithDuration 生成带过期时间的JWT Token
func GenerateJWTWithDuration(apiKey, apiSecret string, durationSeconds int64) (string, error) {
	// 创建JWT claims
	claims := jwt.MapClaims{
		"sub": apiKey,
		"iss": apiKey,
		"aud": "https://business.coinbase.com",
		"exp": durationSeconds,
		"nbf": 0,
		"iat": 0,
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用API Secret作为签名密钥
	tokenString, err := token.SignedString([]byte(apiSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
