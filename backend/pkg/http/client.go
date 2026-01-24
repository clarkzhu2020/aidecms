package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Client HTTP客户端
type Client struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

// ClientOption 客户端选项
type ClientOption func(*Client)

// NewClient 创建一个新的HTTP客户端
func NewClient(options ...ClientOption) *Client {
	client := &Client{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
	}

	// 应用选项
	for _, option := range options {
		option(client)
	}

	return client
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.client.Timeout = timeout
	}
}

// WithBaseURL 设置基础URL
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHeader 设置请求头
func WithHeader(key, value string) ClientOption {
	return func(c *Client) {
		c.headers[key] = value
	}
}

// Get 发送GET请求
func (c *Client) Get(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodGet, path, nil, headers)
}

// Post 发送POST请求
func (c *Client) Post(ctx context.Context, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodPost, path, body, headers)
}

// Put 发送PUT请求
func (c *Client) Put(ctx context.Context, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodPut, path, body, headers)
}

// Delete 发送DELETE请求
func (c *Client) Delete(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	return c.Request(ctx, http.MethodDelete, path, nil, headers)
}

// Request 发送请求
func (c *Client) Request(ctx context.Context, method, path string, body interface{}, headers map[string]string) (*http.Response, error) {
	url := path
	if c.baseURL != "" {
		url = c.baseURL + path
	}

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// 设置默认请求头
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 如果有请求体，设置Content-Type
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.client.Do(req)
}

// GetJSON 发送GET请求并解析JSON响应
func (c *Client) GetJSON(ctx context.Context, path string, headers map[string]string, v interface{}) error {
	resp, err := c.Get(ctx, path, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(v)
}

// PostJSON 发送POST请求并解析JSON响应
func (c *Client) PostJSON(ctx context.Context, path string, body interface{}, headers map[string]string, v interface{}) error {
	resp, err := c.Post(ctx, path, body, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(v)
}

// PutJSON 发送PUT请求并解析JSON响应
func (c *Client) PutJSON(ctx context.Context, path string, body interface{}, headers map[string]string, v interface{}) error {
	resp, err := c.Put(ctx, path, body, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(v)
}

// DeleteJSON 发送DELETE请求并解析JSON响应
func (c *Client) DeleteJSON(ctx context.Context, path string, headers map[string]string, v interface{}) error {
	resp, err := c.Delete(ctx, path, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(v)
}
