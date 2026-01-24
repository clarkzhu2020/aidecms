package web3

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// CoinbaseClient Coinbase Exchange API 客户端
type CoinbaseClient struct {
	apiKey     string
	apiSecret  string
	baseURL    string
	httpClient *http.Client
}

// CoinbaseAccount 账户信息
type CoinbaseAccount struct {
	ID        string `json:"id"`
	Currency  string `json:"currency"`
	Balance   string `json:"balance"`
	Available string `json:"available"`
	Hold      string `json:"hold"`
}

// CoinbaseProduct 交易对信息
type CoinbaseProduct struct {
	ID             string `json:"id"`
	BaseCurrency   string `json:"base_currency"`
	QuoteCurrency  string `json:"quote_currency"`
	BaseIncrement  string `json:"base_increment"`
	QuoteIncrement string `json:"quote_increment"`
	DisplayName    string `json:"display_name"`
	Status         string `json:"status"`
}

// CoinbaseTicker 行情信息
type CoinbaseTicker struct {
	TradeID int64  `json:"trade_id"`
	Price   string `json:"price"`
	Size    string `json:"size"`
	Time    string `json:"time"`
	Bid     string `json:"bid"`
	Ask     string `json:"ask"`
	Volume  string `json:"volume"`
}

// CoinbaseOrder 订单信息
type CoinbaseOrder struct {
	ID            string `json:"id"`
	Price         string `json:"price"`
	Size          string `json:"size"`
	ProductID     string `json:"product_id"`
	Side          string `json:"side"`
	Type          string `json:"type"`
	TimeInForce   string `json:"time_in_force"`
	PostOnly      bool   `json:"post_only"`
	CreatedAt     string `json:"created_at"`
	FillFees      string `json:"fill_fees"`
	FilledSize    string `json:"filled_size"`
	ExecutedValue string `json:"executed_value"`
	Status        string `json:"status"`
	Settled       bool   `json:"settled"`
}

// NewCoinbaseClient 创建 Coinbase 客户端
func NewCoinbaseClient(apiKey, apiSecret string) *CoinbaseClient {
	return &CoinbaseClient{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   "https://api.exchange.coinbase.com",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// generateSignature 生成签名
func (c *CoinbaseClient) generateSignature(timestamp, method, requestPath, body string) string {
	message := timestamp + method + requestPath + body
	h := hmac.New(sha256.New, []byte(c.apiSecret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// request 发送请求
func (c *CoinbaseClient) request(ctx context.Context, method, path string, body string) ([]byte, error) {
	url := c.baseURL + path
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := c.generateSignature(timestamp, method, path, body)

	var reqBody io.Reader
	if body != "" {
		reqBody = bytes.NewBufferString(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CB-ACCESS-KEY", c.apiKey)
	req.Header.Set("CB-ACCESS-SIGN", signature)
	req.Header.Set("CB-ACCESS-TIMESTAMP", timestamp)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("coinbase API error: %s", string(data))
	}

	return data, nil
}

// GetAccounts 获取账户列表
func (c *CoinbaseClient) GetAccounts(ctx context.Context) ([]CoinbaseAccount, error) {
	data, err := c.request(ctx, "GET", "/accounts", "")
	if err != nil {
		return nil, err
	}

	var accounts []CoinbaseAccount
	if err := json.Unmarshal(data, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

// GetAccount 获取指定账户信息
func (c *CoinbaseClient) GetAccount(ctx context.Context, accountID string) (*CoinbaseAccount, error) {
	data, err := c.request(ctx, "GET", "/accounts/"+accountID, "")
	if err != nil {
		return nil, err
	}

	var account CoinbaseAccount
	if err := json.Unmarshal(data, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// GetProducts 获取交易对列表
func (c *CoinbaseClient) GetProducts(ctx context.Context) ([]CoinbaseProduct, error) {
	data, err := c.request(ctx, "GET", "/products", "")
	if err != nil {
		return nil, err
	}

	var products []CoinbaseProduct
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// GetTicker 获取行情
func (c *CoinbaseClient) GetTicker(ctx context.Context, productID string) (*CoinbaseTicker, error) {
	data, err := c.request(ctx, "GET", "/products/"+productID+"/ticker", "")
	if err != nil {
		return nil, err
	}

	var ticker CoinbaseTicker
	if err := json.Unmarshal(data, &ticker); err != nil {
		return nil, err
	}

	return &ticker, nil
}

// GetOrders 获取订单列表
func (c *CoinbaseClient) GetOrders(ctx context.Context, status string) ([]CoinbaseOrder, error) {
	path := "/orders"
	if status != "" {
		path += "?status=" + status
	}

	data, err := c.request(ctx, "GET", path, "")
	if err != nil {
		return nil, err
	}

	var orders []CoinbaseOrder
	if err := json.Unmarshal(data, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// PlaceOrder 下单
func (c *CoinbaseClient) PlaceOrder(ctx context.Context, productID, side, orderType, size, price string) (*CoinbaseOrder, error) {
	orderData := map[string]interface{}{
		"product_id": productID,
		"side":       side,
		"type":       orderType,
		"size":       size,
	}

	if orderType == "limit" && price != "" {
		orderData["price"] = price
	}

	body, err := json.Marshal(orderData)
	if err != nil {
		return nil, err
	}

	data, err := c.request(ctx, "POST", "/orders", string(body))
	if err != nil {
		return nil, err
	}

	var order CoinbaseOrder
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

// CancelOrder 取消订单
func (c *CoinbaseClient) CancelOrder(ctx context.Context, orderID string) error {
	_, err := c.request(ctx, "DELETE", "/orders/"+orderID, "")
	return err
}

// GetBalance 获取指定币种余额
func (c *CoinbaseClient) GetBalance(ctx context.Context, currency string) (string, error) {
	accounts, err := c.GetAccounts(ctx)
	if err != nil {
		return "", err
	}

	for _, account := range accounts {
		if account.Currency == currency {
			return account.Balance, nil
		}
	}

	return "0", nil
}

// GetPrice 获取价格
func (c *CoinbaseClient) GetPrice(ctx context.Context, pair string) (string, error) {
	ticker, err := c.GetTicker(ctx, pair)
	if err != nil {
		return "", err
	}
	return ticker.Price, nil
}

// GetBalances 获取所有余额
func (c *CoinbaseClient) GetBalances(ctx context.Context) (map[string]string, error) {
	accounts, err := c.GetAccounts(ctx)
	if err != nil {
		return nil, err
	}

	balances := make(map[string]string)
	for _, account := range accounts {
		balances[account.Currency] = account.Balance
	}

	return balances, nil
}
