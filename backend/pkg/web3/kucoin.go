package web3

import (
	"bytes"
	"context"
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

// KuCoinClient KuCoin Exchange API 客户端
type KuCoinClient struct {
	apiKey     string
	apiSecret  string
	passphrase string
	baseURL    string
	httpClient *http.Client
}

// KuCoinResponse 通用响应
type KuCoinResponse struct {
	Code string          `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

// KuCoinAccount 账户信息
type KuCoinAccount struct {
	ID        string `json:"id"`
	Currency  string `json:"currency"`
	Type      string `json:"type"`
	Balance   string `json:"balance"`
	Available string `json:"available"`
	Holds     string `json:"holds"`
}

// KuCoinTicker 行情信息
type KuCoinTicker struct {
	Symbol       string `json:"symbol"`
	Buy          string `json:"buy"`
	Sell         string `json:"sell"`
	ChangeRate   string `json:"changeRate"`
	ChangePrice  string `json:"changePrice"`
	High         string `json:"high"`
	Low          string `json:"low"`
	Vol          string `json:"vol"`
	VolValue     string `json:"volValue"`
	Last         string `json:"last"`
	AveragePrice string `json:"averagePrice"`
	Time         int64  `json:"time"`
}

// KuCoinOrder 订单信息
type KuCoinOrder struct {
	ID            string `json:"id"`
	Symbol        string `json:"symbol"`
	OpType        string `json:"opType"`
	Type          string `json:"type"`
	Side          string `json:"side"`
	Price         string `json:"price"`
	Size          string `json:"size"`
	Funds         string `json:"funds"`
	DealFunds     string `json:"dealFunds"`
	DealSize      string `json:"dealSize"`
	Fee           string `json:"fee"`
	FeeCurrency   string `json:"feeCurrency"`
	Stop          string `json:"stop"`
	StopTriggered bool   `json:"stopTriggered"`
	StopPrice     string `json:"stopPrice"`
	TimeInForce   string `json:"timeInForce"`
	PostOnly      bool   `json:"postOnly"`
	Hidden        bool   `json:"hidden"`
	Iceberg       bool   `json:"iceberg"`
	VisibleSize   string `json:"visibleSize"`
	CancelAfter   int64  `json:"cancelAfter"`
	Channel       string `json:"channel"`
	ClientOid     string `json:"clientOid"`
	Remark        string `json:"remark"`
	Tags          string `json:"tags"`
	IsActive      bool   `json:"isActive"`
	CancelExist   bool   `json:"cancelExist"`
	CreatedAt     int64  `json:"createdAt"`
	TradeType     string `json:"tradeType"`
}

// KuCoinSymbol 交易对信息
type KuCoinSymbol struct {
	Symbol          string `json:"symbol"`
	Name            string `json:"name"`
	BaseCurrency    string `json:"baseCurrency"`
	QuoteCurrency   string `json:"quoteCurrency"`
	BaseMinSize     string `json:"baseMinSize"`
	QuoteMinSize    string `json:"quoteMinSize"`
	BaseMaxSize     string `json:"baseMaxSize"`
	QuoteMaxSize    string `json:"quoteMaxSize"`
	BaseIncrement   string `json:"baseIncrement"`
	QuoteIncrement  string `json:"quoteIncrement"`
	PriceIncrement  string `json:"priceIncrement"`
	FeeCurrency     string `json:"feeCurrency"`
	EnableTrading   bool   `json:"enableTrading"`
	IsMarginEnabled bool   `json:"isMarginEnabled"`
	PriceLimitRate  string `json:"priceLimitRate"`
}

// NewKuCoinClient 创建 KuCoin 客户端
func NewKuCoinClient(apiKey, apiSecret, passphrase string) *KuCoinClient {
	return &KuCoinClient{
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		passphrase: passphrase,
		baseURL:    "https://api.kucoin.com",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// generateSignature 生成签名
func (k *KuCoinClient) generateSignature(timestamp, method, endpoint, body string) string {
	strToSign := timestamp + method + endpoint + body
	h := hmac.New(sha256.New, []byte(k.apiSecret))
	h.Write([]byte(strToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// generatePassphrase 生成加密的 passphrase
func (k *KuCoinClient) generatePassphrase() string {
	h := hmac.New(sha256.New, []byte(k.apiSecret))
	h.Write([]byte(k.passphrase))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// request 发送请求
func (k *KuCoinClient) request(ctx context.Context, method, endpoint string, body string) ([]byte, error) {
	url := k.baseURL + endpoint
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := k.generateSignature(timestamp, method, endpoint, body)
	passphrase := k.generatePassphrase()

	var reqBody io.Reader
	if body != "" {
		reqBody = bytes.NewBufferString(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("KC-API-KEY", k.apiKey)
	req.Header.Set("KC-API-SIGN", signature)
	req.Header.Set("KC-API-TIMESTAMP", timestamp)
	req.Header.Set("KC-API-PASSPHRASE", passphrase)
	req.Header.Set("KC-API-KEY-VERSION", "2")

	resp, err := k.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp KuCoinResponse
	if err := json.Unmarshal(data, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.Code != "200000" {
		return nil, fmt.Errorf("kucoin API error: %s - %s", apiResp.Code, apiResp.Msg)
	}

	return apiResp.Data, nil
}

// GetAccounts 获取账户列表
func (k *KuCoinClient) GetAccounts(ctx context.Context) ([]KuCoinAccount, error) {
	data, err := k.request(ctx, "GET", "/api/v1/accounts", "")
	if err != nil {
		return nil, err
	}

	var accounts []KuCoinAccount
	if err := json.Unmarshal(data, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

// GetAccount 获取指定账户
func (k *KuCoinClient) GetAccount(ctx context.Context, accountID string) (*KuCoinAccount, error) {
	data, err := k.request(ctx, "GET", "/api/v1/accounts/"+accountID, "")
	if err != nil {
		return nil, err
	}

	var account KuCoinAccount
	if err := json.Unmarshal(data, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

// GetTicker 获取行情
func (k *KuCoinClient) GetTicker(ctx context.Context, symbol string) (*KuCoinTicker, error) {
	endpoint := "/api/v1/market/orderbook/level1?symbol=" + symbol
	data, err := k.request(ctx, "GET", endpoint, "")
	if err != nil {
		return nil, err
	}

	var ticker KuCoinTicker
	if err := json.Unmarshal(data, &ticker); err != nil {
		return nil, err
	}

	return &ticker, nil
}

// GetAllTickers 获取所有交易对行情
func (k *KuCoinClient) GetAllTickers(ctx context.Context) ([]KuCoinTicker, error) {
	data, err := k.request(ctx, "GET", "/api/v1/market/allTickers", "")
	if err != nil {
		return nil, err
	}

	var resp struct {
		Ticker []KuCoinTicker `json:"ticker"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.Ticker, nil
}

// GetSymbols 获取交易对列表
func (k *KuCoinClient) GetSymbols(ctx context.Context) ([]KuCoinSymbol, error) {
	data, err := k.request(ctx, "GET", "/api/v1/symbols", "")
	if err != nil {
		return nil, err
	}

	var symbols []KuCoinSymbol
	if err := json.Unmarshal(data, &symbols); err != nil {
		return nil, err
	}

	return symbols, nil
}

// PlaceOrder 下单
func (k *KuCoinClient) PlaceOrder(ctx context.Context, clientOid, side, symbol, orderType, size, price string) (*KuCoinOrder, error) {
	orderData := map[string]interface{}{
		"clientOid": clientOid,
		"side":      side,
		"symbol":    symbol,
		"type":      orderType,
		"size":      size,
	}

	if orderType == "limit" && price != "" {
		orderData["price"] = price
	}

	body, err := json.Marshal(orderData)
	if err != nil {
		return nil, err
	}

	data, err := k.request(ctx, "POST", "/api/v1/orders", string(body))
	if err != nil {
		return nil, err
	}

	var result struct {
		OrderID string `json:"orderId"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	// Get order details
	return k.GetOrder(ctx, result.OrderID)
}

// GetOrder 获取订单详情
func (k *KuCoinClient) GetOrder(ctx context.Context, orderID string) (*KuCoinOrder, error) {
	data, err := k.request(ctx, "GET", "/api/v1/orders/"+orderID, "")
	if err != nil {
		return nil, err
	}

	var order KuCoinOrder
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOrders 获取订单列表
func (k *KuCoinClient) GetOrders(ctx context.Context, status string) ([]KuCoinOrder, error) {
	endpoint := "/api/v1/orders"
	if status != "" {
		endpoint += "?status=" + status
	}

	data, err := k.request(ctx, "GET", endpoint, "")
	if err != nil {
		return nil, err
	}

	var resp struct {
		Items []KuCoinOrder `json:"items"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.Items, nil
}

// CancelOrder 取消订单
func (k *KuCoinClient) CancelOrder(ctx context.Context, orderID string) error {
	_, err := k.request(ctx, "DELETE", "/api/v1/orders/"+orderID, "")
	return err
}

// GetBalance 获取指定币种余额
func (k *KuCoinClient) GetBalance(ctx context.Context, currency string) (string, error) {
	accounts, err := k.GetAccounts(ctx)
	if err != nil {
		return "", err
	}

	for _, account := range accounts {
		if account.Currency == currency && account.Type == "trade" {
			return account.Balance, nil
		}
	}

	return "0", nil
}

// GetBalances 获取所有余额
func (k *KuCoinClient) GetBalances(ctx context.Context) (map[string]string, error) {
	accounts, err := k.GetAccounts(ctx)
	if err != nil {
		return nil, err
	}

	balances := make(map[string]string)
	for _, account := range accounts {
		if account.Type == "trade" {
			balances[account.Currency] = account.Balance
		}
	}

	return balances, nil
}

// GetPrice 获取价格
func (k *KuCoinClient) GetPrice(ctx context.Context, symbol string) (string, error) {
	ticker, err := k.GetTicker(ctx, symbol)
	if err != nil {
		return "", err
	}
	return ticker.Last, nil
}
