package kucoin

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// HTTPClient KuCoin HTTP客户端
type HTTPClient struct {
	apiKey     string
	apiSecret  string
	passphrase string
	baseURL    string
	httpClient *http.Client
}

// NewHTTPClient 创建HTTP客户端
func NewHTTPClient(apiKey, apiSecret, passphrase, baseURL string) *HTTPClient {
	return &HTTPClient{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		passphrase: passphrase,
		baseURL:    baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// generateSignature 生成API签名
func (c *HTTPClient) generateSignature(timestamp, method, endpoint, body string) string {
	// 构造签名内容: timestamp + method + endpoint + body
	signStr := timestamp + method + endpoint + body

	// 使用HMAC SHA256进行签名
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(signStr))

	// Base64编码
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// generatePassphraseSign 生成passphrase签名
func (c *HTTPClient) generatePassphraseSign() string {
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(c.passphrase))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// getTimestamp 获取当前时间戳（毫秒）
func (c *HTTPClient) getTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
}

// Do 执行HTTP请求
func (c *HTTPClient) Do(method, endpoint string, body []byte, useAuth bool) ([]byte, error) {
	url := c.baseURL + endpoint

	// 创建请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 如果需要认证，添加认证头
	if useAuth {
		timestamp := c.getTimestamp()
		bodyStr := ""
		if body != nil {
			bodyStr = string(body)
		}

		signature := c.generateSignature(timestamp, method, endpoint, bodyStr)
		passphraseSign := c.generatePassphraseSign()

		req.Header.Set("KC-API-KEY", c.apiKey)
		req.Header.Set("KC-API-SIGN", signature)
		req.Header.Set("KC-API-TIMESTAMP", timestamp)
		req.Header.Set("KC-API-PASSPHRASE", passphraseSign)
		req.Header.Set("KC-API-KEY-VERSION", "2")
	}

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
		var apiError struct {
			Code    string `json:"code"`
			Message string `json:"msg"`
		}
		if err := json.Unmarshal(respBody, &apiError); err == nil && apiError.Message != "" {
			return nil, fmt.Errorf("API error (status %d): %s - %s", resp.StatusCode, apiError.Code, apiError.Message)
		}
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Get 执行GET请求
func (c *HTTPClient) Get(endpoint string, useAuth bool) ([]byte, error) {
	return c.Do("GET", endpoint, nil, useAuth)
}

// Post 执行POST请求
func (c *HTTPClient) Post(endpoint string, body []byte, useAuth bool) ([]byte, error) {
	return c.Do("POST", endpoint, body, useAuth)
}

// Put 执行PUT请求
func (c *HTTPClient) Put(endpoint string, body []byte, useAuth bool) ([]byte, error) {
	return c.Do("PUT", endpoint, body, useAuth)
}

// Delete 执行DELETE请求
func (c *HTTPClient) Delete(endpoint string, useAuth bool) ([]byte, error) {
	return c.Do("DELETE", endpoint, nil, useAuth)
}

// unmarshalJSON 反序列化JSON响应
func (c *HTTPClient) unmarshalJSON(data []byte, v interface{}) error {
	var response struct {
		Code string          `json:"code"`
		Data json.RawMessage `json:"data"`
		Msg  string         `json:"msg"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	// 检查响应码
	if response.Code != "200000" && response.Code != "" {
		return fmt.Errorf("API error: %s - %s", response.Code, response.Msg)
	}

	// 反序列化data字段
	if err := json.Unmarshal(response.Data, v); err != nil {
		return err
	}

	return nil
}
