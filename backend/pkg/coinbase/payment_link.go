package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// PaymentLinkRequest 创建支付链接请求
type PaymentLinkRequest struct {
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Description string `json:"description,omitempty"`
	Title       string `json:"title,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`
	CancelURL   string `json:"cancel_url,omitempty"`
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// PaymentLink 支付链接
type PaymentLink struct {
	ID            string                 `json:"id"`
	PaymentURL    string                 `json:"payment_url"`
	Amount        string                 `json:"amount"`
	Currency      string                 `json:"currency"`
	Description   string                 `json:"description"`
	Title         string                 `json:"title"`
	Status        string                 `json:"status"`
	CreatedAt     string                 `json:"created_at"`
	UpdatedAt     string                 `json:"updated_at"`
	ExpiresAt     string                 `json:"expires_at,omitempty"`
	PaymentStatus string                 `json:"payment_status"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	RawData       map[string]interface{} `json:"-"`
}

// PaymentLinkList 支付链接列表
type PaymentLinkList struct {
	PaymentLinks []PaymentLink `json:"payment_links"`
	TotalCount   int           `json:"total_count"`
	NextPageToken string       `json:"next_page_token,omitempty"`
}

// PaymentLinkService 支付链接服务
type PaymentLinkService struct {
	client *Client
}

// NewPaymentLinkService 创建支付链接服务
func NewPaymentLinkService() *PaymentLinkService {
	return &PaymentLinkService{
		client: GetClient(),
	}
}

// CreatePaymentLink 创建支付链接
func (s *PaymentLinkService) CreatePaymentLink(ctx context.Context, req *PaymentLinkRequest) (*PaymentLink, error) {
	if !IsInitialized() {
		return nil, fmt.Errorf("Coinbase client not initialized")
	}

	// 生成JWT Token
	token, err := GenerateJWTWithDuration(s.client.config.APIKey, s.client.config.APISecret, 300) // 5分钟
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	// 构建请求头
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	// 构建请求体
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 调用Coinbase API创建支付链接
	resp, err := s.client.httpClient.Post("/payment-links", body, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment link: %w", err)
	}

	// 解析响应
	var paymentLink PaymentLink
	if err := json.Unmarshal(resp, &paymentLink); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	log.Printf("Coinbase payment link created: ID=%s, URL=%s", paymentLink.ID, paymentLink.PaymentURL)
	return &paymentLink, nil
}

// GetPaymentLink 获取支付链接详情
func (s *PaymentLinkService) GetPaymentLink(ctx context.Context, linkID string) (*PaymentLink, error) {
	if !IsInitialized() {
		return nil, fmt.Errorf("Coinbase client not initialized")
	}

	// 生成JWT Token
	token, err := GenerateJWT(s.client.config.APIKey, s.client.config.APISecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	// 构建请求头
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	endpoint := fmt.Sprintf("/payment-links/%s", linkID)
	resp, err := s.client.httpClient.Get(endpoint, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment link: %w", err)
	}

	var paymentLink PaymentLink
	if err := json.Unmarshal(resp, &paymentLink); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &paymentLink, nil
}

// ListPaymentLinks 列出支付链接
func (s *PaymentLinkService) ListPaymentLinks(ctx context.Context, limit, offset int) (*PaymentLinkList, error) {
	if !IsInitialized() {
		return nil, fmt.Errorf("Coinbase client not initialized")
	}

	// 生成JWT Token
	token, err := GenerateJWT(s.client.config.APIKey, s.client.config.APISecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	// 构建请求头
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	endpoint := fmt.Sprintf("/payment-links?limit=%d&offset=%d", limit, offset)
	resp, err := s.client.httpClient.Get(endpoint, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to list payment links: %w", err)
	}

	var paymentLinks PaymentLinkList
	if err := json.Unmarshal(resp, &paymentLinks); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &paymentLinks, nil
}

// UpdatePaymentLink 更新支付链接
func (s *PaymentLinkService) UpdatePaymentLink(ctx context.Context, linkID string, req *PaymentLinkRequest) (*PaymentLink, error) {
	if !IsInitialized() {
		return nil, fmt.Errorf("Coinbase client not initialized")
	}

	// 生成JWT Token
	token, err := GenerateJWTWithDuration(s.client.config.APIKey, s.client.config.APISecret, 300)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	// 构建请求头
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	// 构建请求体
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("/payment-links/%s", linkID)
	resp, err := s.client.httpClient.Put(endpoint, body, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to update payment link: %w", err)
	}

	var paymentLink PaymentLink
	if err := json.Unmarshal(resp, &paymentLink); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &paymentLink, nil
}

// DeletePaymentLink 删除支付链接
func (s *PaymentLinkService) DeletePaymentLink(ctx context.Context, linkID string) error {
	if !IsInitialized() {
		return fmt.Errorf("Coinbase client not initialized")
	}

	// 生成JWT Token
	token, err := GenerateJWT(s.client.config.APIKey, s.client.config.APISecret)
	if err != nil {
		return fmt.Errorf("failed to generate JWT: %w", err)
	}

	// 构建请求头
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	endpoint := fmt.Sprintf("/payment-links/%s", linkID)
	_, err = s.client.httpClient.Delete(endpoint, headers)
	if err != nil {
		return fmt.Errorf("failed to delete payment link: %w", err)
	}

	return nil
}

// VerifyWebhook 验证Webhook签名
func (s *PaymentLinkService) VerifyWebhook(payload []byte, signature, timestamp string) bool {
	// Coinbase使用X-CC-Webhook-Signature和X-CC-Webhook-Timestamp headers
	// 实现HMAC-SHA256验证
	// 注意: 这里需要根据Coinbase实际的签名验证逻辑实现
	// 由于Coinbase SDK可能提供验证方法，这里预留接口
	return true
}
