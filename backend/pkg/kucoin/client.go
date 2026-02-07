package kucoin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"time"
)

// Config KuCoin配置
type Config struct {
	APIKey     string
	APISecret  string
	Passphrase string
	IsSandbox  bool
}

// Client KuCoin客户端
type Client struct {
	httpClient *HTTPClient
	config     *Config
}

var (
	// 全局KuCoin客户端实例
	kucoinClient *Client
)

// NewClient 创建新的KuCoin客户端
func NewClient(config *Config) (*Client, error) {
	// 设置基础URL
	baseURL := "https://api.kucoin.com"
	if config.IsSandbox {
		baseURL = "https://openapi-sandbox.kucoin.com"
	}

	// 创建HTTP客户端
	httpClient := NewHTTPClient(config.APIKey, config.APISecret, config.Passphrase, baseURL)

	return &Client{
		httpClient: httpClient,
		config:     config,
	}, nil
}

// Init 初始化KuCoin客户端
func Init(apiKey, apiSecret, passphrase string, isSandbox bool) error {
	config := &Config{
		APIKey:     apiKey,
		APISecret:  apiSecret,
		Passphrase: passphrase,
		IsSandbox:  isSandbox,
	}

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	kucoinClient = client
	log.Println("KuCoin client initialized successfully")
	return nil
}

// GetClient 获取KuCoin客户端实例
func GetClient() *Client {
	return kucoinClient
}

// GetConfig 获取KuCoin配置
func (c *Client) GetConfig() *Config {
	return c.config
}

// IsInitialized 检查客户端是否已初始化
func IsInitialized() bool {
	return kucoinClient != nil
}

// GenerateSignature 生成API签名
func GenerateSignature(timestamp, method, endpoint, body, secret string) string {
	// 构造签名内容: timestamp + method + endpoint + body
	signStr := timestamp + method + endpoint + body

	// 使用HMAC SHA256进行签名
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(signStr))

	// Base64编码
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// GeneratePassphraseSign 生成passphrase签名
func GeneratePassphraseSign(passphrase, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(passphrase))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// GetTimestamp 获取当前时间戳（毫秒）
func GetTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
}

// ValidateTimestamp 验证时间戳是否在有效范围内（30秒内）
func ValidateTimestamp(timestamp int64) bool {
	now := time.Now().UnixNano() / 1e6
	diff := now - timestamp
	return diff >= -30000 && diff <= 30000
}

// GenerateAuthHeaders 生成认证请求头
func (c *Client) GenerateAuthHeaders(method, endpoint string, body []byte) map[string]string {
	timestamp := GetTimestamp()
	bodyStr := ""
	if body != nil {
		bodyStr = string(body)
	}

	signature := GenerateSignature(timestamp, method, endpoint, bodyStr, c.config.APISecret)
	passphraseSign := GeneratePassphraseSign(c.config.Passphrase, c.config.APISecret)

	return map[string]string{
		"KC-API-KEY":        c.config.APIKey,
		"KC-API-SIGN":       signature,
		"KC-API-TIMESTAMP":  timestamp,
		"KC-API-PASSPHRASE": passphraseSign,
		"KC-API-KEY-VERSION": "2",
		"Content-Type":      "application/json",
	}
}

// GetServerTime 获取服务器时间
func (c *Client) GetServerTime() (int64, error) {
	endpoint := "/api/v1/timestamp"
	body, err := c.httpClient.Get(endpoint, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get server time: %w", err)
	}

	var response struct {
		Data int64 `json:"data"`
	}

	if err := c.httpClient.unmarshalJSON(body, &response); err != nil {
		return 0, err
	}

	return response.Data, nil
}

// SyncServerTime 同步服务器时间
func (c *Client) SyncServerTime() error {
	serverTime, err := c.GetServerTime()
	if err != nil {
		return err
	}

	localTime := time.Now().UnixNano() / 1e6
	diff := serverTime - localTime

	log.Printf("KuCoin server time diff: %dms", diff)

	if abs(diff) > 30000 {
		return fmt.Errorf("server time diff too large: %dms (max 30000ms)", diff)
	}

	return nil
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
