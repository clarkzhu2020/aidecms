package moonpay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient MoonPay HTTP客户端
type HTTPClient struct {
	apiKey      string
	secretKey   string
	baseURL     string
	httpClient  *http.Client
}

// NewHTTPClient 创建HTTP客户端
func NewHTTPClient(apiKey, secretKey, baseURL string) *HTTPClient {
	return &HTTPClient{
		apiKey:    apiKey,
		secretKey: secretKey,
		baseURL:   baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Do 执行HTTP请求
func (c *HTTPClient) Do(method, endpoint string, body []byte) ([]byte, error) {
	url := c.baseURL + endpoint

	// 创建请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	// 执行请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 检查状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiError map[string]interface{}
		if err := json.Unmarshal(respBody, &apiError); err == nil {
			return nil, fmt.Errorf("API error (status %d): %v", resp.StatusCode, apiError)
		}
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Get 执行GET请求
func (c *HTTPClient) Get(endpoint string) ([]byte, error) {
	return c.Do("GET", endpoint, nil)
}

// Post 执行POST请求
func (c *HTTPClient) Post(endpoint string, body []byte) ([]byte, error) {
	return c.Do("POST", endpoint, body)
}

// Put 执行PUT请求
func (c *HTTPClient) Put(endpoint string, body []byte) ([]byte, error) {
	return c.Do("PUT", endpoint, body)
}

// Delete 执行DELETE请求
func (c *HTTPClient) Delete(endpoint string) ([]byte, error) {
	return c.Do("DELETE", endpoint, nil)
}
