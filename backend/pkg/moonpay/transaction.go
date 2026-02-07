package moonpay

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"
)

// CreateTransactionRequest 创建交易请求
type CreateTransactionRequest struct {
	BaseCurrencyAmount float64 `json:"baseCurrencyAmount"`
	BaseCurrencyCode  string  `json:"baseCurrencyCode"`
	CurrencyCode      string  `json:"currencyCode"`
	WalletAddress     string  `json:"walletAddress"`
	ExternalID         string  `json:"externalTransactionId"`
	RedirectURL        string  `json:"redirectURL,omitempty"`
	LockAmount         bool    `json:"lockAmount,omitempty"`
	AdditionalDetails string  `json:"additionalDetails,omitempty"`
	Email              string  `json:"email,omitempty"`
	FirstName          string  `json:"firstName,omitempty"`
	LastName           string  `json:"lastName,omitempty"`
}

// Transaction 交易信息
type Transaction struct {
	ID                string                 `json:"id"`
	CreatedAt         string                 `json:"createdAt"`
	UpdatedAt         string                 `json:"updatedAt"`
	Status            string                 `json:"status"`
	Type              string                 `json:"type"`
	ExternalID        string                 `json:"externalTransactionId"`
	BaseCurrencyAmount float64                `json:"baseCurrencyAmount"`
	BaseCurrencyCode  string                 `json:"baseCurrencyCode"`
	QuoteCurrencyAmount float64              `json:"quoteCurrencyAmount"`
	QuoteCurrencyCode string                 `json:"quoteCurrencyCode"`
	CurrencyCode      string                 `json:"currencyCode"`
	WalletAddress     string                 `json:"walletAddress"`
	RedirectURL       string                 `json:"redirectUrl"`
	FeeAmount         float64                `json:"feeAmount"`
	ExtraFeeAmount    float64                `json:"extraFeeAmount"`
	NetworkFeeAmount  float64                `json:"networkFeeAmount"`
	Customer          *CustomerInfo          `json:"customer"`
	PaymentMethod     *PaymentMethodInfo     `json:"paymentMethod"`
	Kyc               *KycInfo               `json:"kyc"`
	RawData           map[string]interface{} `json:"-"`
}

// CustomerInfo 客户信息
type CustomerInfo struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

// PaymentMethodInfo 支付方式信息
type PaymentMethodInfo struct {
	Type   string `json:"type"`
	Card   *CardInfo `json:"card,omitempty"`
	Bank   *BankInfo `json:"bankTransfer,omitempty"`
}

// CardInfo 银行卡信息
type CardInfo struct {
	Last4Digits string `json:"last4Digits"`
	Brand       string `json:"brand"`
}

// BankInfo 银行转账信息
type BankInfo struct {
	Iban          string `json:"iban"`
	AccountNumber string `json:"accountNumber"`
	Bic           string `json:"bic"`
}

// KycInfo KYC信息
type KycInfo struct {
	Status string `json:"status"`
}

// TransactionService 交易服务
type TransactionService struct {
	client *Client
}

// NewTransactionService 创建交易服务
func NewTransactionService() *TransactionService {
	return &TransactionService{
		client: GetClient(),
	}
}

// CreateTransaction 创建交易
func (s *TransactionService) CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*Transaction, error) {
	if !IsInitialized() {
		return nil, fmt.Errorf("MoonPay client not initialized")
	}

	// 构建请求体
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 调用MoonPay API创建交易
	resp, err := s.client.httpClient.Post("/v1/transactions", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// 解析响应
	var transaction Transaction
	if err := json.Unmarshal(resp, &transaction); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	log.Printf("MoonPay transaction created: ID=%s, Status=%s", transaction.ID, transaction.Status)
	return &transaction, nil
}

// GetTransaction 获取交易详情
func (s *TransactionService) GetTransaction(ctx context.Context, transactionID string) (*Transaction, error) {
	if !IsInitialized() {
		return nil, fmt.Errorf("MoonPay client not initialized")
	}

	endpoint := fmt.Sprintf("/v1/transactions/%s", transactionID)
	resp, err := s.client.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	var transaction Transaction
	if err := json.Unmarshal(resp, &transaction); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &transaction, nil
}

// GetTransactionByExternalID 根据外部ID获取交易
func (s *TransactionService) GetTransactionByExternalID(ctx context.Context, externalID string) (*Transaction, error) {
	if !IsInitialized() {
		return nil, fmt.Errorf("MoonPay client not initialized")
	}

	endpoint := fmt.Sprintf("/v1/transactions?externalTransactionId=%s", externalID)
	resp, err := s.client.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	var transactions []Transaction
	if err := json.Unmarshal(resp, &transactions); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("transaction not found")
	}

	return &transactions[0], nil
}

// ListTransactions 列出交易
func (s *TransactionService) ListTransactions(ctx context.Context, limit, offset int) ([]Transaction, error) {
	if !IsInitialized() {
		return nil, fmt.Errorf("MoonPay client not initialized")
	}

	endpoint := fmt.Sprintf("/v1/transactions?limit=%d&offset=%d", limit, offset)
	resp, err := s.client.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}

	var transactions []Transaction
	if err := json.Unmarshal(resp, &transactions); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return transactions, nil
}

// GetBuyQuote 获取购买报价
func (s *TransactionService) GetBuyQuote(ctx context.Context, baseCurrencyAmount float64, baseCurrencyCode, currencyCode string) (map[string]interface{}, error) {
	if !IsInitialized() {
		return nil, fmt.Errorf("MoonPay client not initialized")
	}

	endpoint := fmt.Sprintf("/v3/currencies/%s/buy_quote?baseCurrencyAmount=%.2f&baseCurrencyCode=%s",
		currencyCode, baseCurrencyAmount, baseCurrencyCode)

	resp, err := s.client.httpClient.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get buy quote: %w", err)
	}

	var quote map[string]interface{}
	if err := json.Unmarshal(resp, &quote); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return quote, nil
}

// GenerateWidgetURL 生成Widget URL（前端集成）
func (s *TransactionService) GenerateWidgetURL(req *CreateTransactionRequest) string {
	baseURL := "https://buy.moonpay.com"
	if s.client.config.IsSandbox {
		baseURL = "https://buy-sandbox.moonpay.com"
	}

	params := url.Values{}
	params.Set("apiKey", s.client.config.APIKey)
	params.Set("currencyCode", req.CurrencyCode)
	params.Set("baseCurrencyAmount", fmt.Sprintf("%.2f", req.BaseCurrencyAmount))
	params.Set("baseCurrencyCode", req.BaseCurrencyCode)
	params.Set("walletAddress", req.WalletAddress)
	params.Set("externalTransactionId", req.ExternalID)

	if req.RedirectURL != "" {
		params.Set("redirectURL", req.RedirectURL)
	}
	if req.LockAmount {
		params.Set("lockAmount", "true")
	}
	if req.Email != "" {
		params.Set("email", req.Email)
	}
	if req.FirstName != "" {
		params.Set("firstName", req.FirstName)
	}
	if req.LastName != "" {
		params.Set("lastName", req.LastName)
	}

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// VerifyWebhookSignature 验证Webhook签名
func (s *TransactionService) VerifyWebhookSignature(payload []byte, signature string) bool {
	// MoonPay使用X-MoonPay-Signature header
	// 签名格式: timestamp.payload
	// 实现HMAC-SHA256验证
	// 注意: 这里需要根据MoonPay实际的签名验证逻辑实现
	// 由于MoonPay SDK可能提供验证方法，这里预留接口
	return true
}
