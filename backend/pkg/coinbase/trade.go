package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// OrderRequest 订单请求
type OrderRequest struct {
	ProductID       string  `json:"product_id"`        // 交易对，如 "BTC-USD"
	Side            string  `json:"side"`              // "buy" 或 "sell"
	OrderType       string  `json:"order_type"`        // "market" 或 "limit"
	Size            string  `json:"size,omitempty"`     // 买入或卖出的数量
	Funds           string  `json:"funds,omitempty"`    // 使用的金额（仅限市价单）
	LimitPrice      string  `json:"limit_price,omitempty"` // 限价（仅限限价单）
	StopPrice       string  `json:"stop_price,omitempty"`  // 止损价格
	Stop            string  `json:"stop,omitempty"`      // "loss" 或 "entry"
	TimeInForce     string  `json:"time_in_force,omitempty"` // "GTC", "IOC", "FOK", "GTD"
	ClientOid       string  `json:"client_oid,omitempty"`     // 客户端订单ID
	PostOnly        bool    `json:"post_only,omitempty"`       // 仅作为maker
	SelfTradePrevention string `json:"self_trade_prevention,omitempty"` // "dc", "co", "cn", "cb"
}

// Order 订单信息
type Order struct {
	ID                string                 `json:"id"`
	ClientOid         string                 `json:"client_oid,omitempty"`
	ProductID         string                 `json:"product_id"`
	Side              string                 `json:"side"`
	OrderType         string                 `json:"order_type"`
	Size              string                 `json:"size"`
	Funds             string                 `json:"funds,omitempty"`
	LimitPrice        string                 `json:"limit_price,omitempty"`
	StopPrice         string                 `json:"stop_price,omitempty"`
	Stop              string                 `json:"stop,omitempty"`
	TimeInForce       string                 `json:"time_in_force,omitempty"`
	CreatedAt         string                 `json:"created_at"`
	UpdatedAt         string                 `json:"updated_at"`
	ExpiredAt         string                 `json:"expired_at,omitempty"`
	FilledSize        string                 `json:"filled_size"`
	AverageFillPrice  string                 `json:"average_fill_price,omitempty"`
	FillFees         string                 `json:"fill_fees"`
	Status            string                 `json:"status"` // "open", "filled", "rejected", "cancelled"
	Settled           bool                   `json:"settled"`
	PostOnly          bool                   `json:"post_only"`
	CreatedBy         string                 `json:"created_by,omitempty"`
	SelfTradePrevention string               `json:"self_trade_prevention,omitempty"`
	RawData           map[string]interface{} `json:"-"`
}

// Account 账户信息
type Account struct {
	AccountID      string                 `json:"account_id"`
	AvailableBalance map[string]float64   `json:"available_balance"`
	DefaultCurrency string                 `json:"default_currency"`
	Name           string                 `json:"name"`
	Type           string                 `json:"type"` // "exchange", "wallet"
	CreatedAt      string                 `json:"created_at"`
	UpdatedAt      string                 `json:"updated_at"`
	RawData        map[string]interface{} `json:"-"`
}

// Product 产品信息
type Product struct {
	ProductID           string  `json:"product_id"`
	BaseCurrency        string  `json:"base_currency"`
	QuoteCurrency       string  `json:"quote_currency"`
	BaseMinSize         string  `json:"base_min_size"`
	BaseMaxSize         string  `json:"base_max_size"`
	QuoteIncrement      string  `json:"quote_increment"`
	BaseIncrement       string  `json:"base_increment"`
	DisplayName         string  `json:"display_name"`
	Status              string  `json:"status"` // "online", "delisted", "suspend_only"
	MarginEnabled       bool    `json:"margin_enabled"`
	IsMarketOrderable   bool    `json:"is_market_orderable"`
	IsLimitOrderable    bool    `json:"is_limit_orderable"`
	IsStopOrderable     bool    `json:"is_stop_orderable"`
}

// Ticker 行情信息
type Ticker struct {
	ProductID      string  `json:"product_id"`
	Price          string  `json:"price"`
	Open24H        string  `json:"open_24h"`
	High24H        string  `json:"high_24h"`
	Low24H         string  `json:"low_24h"`
	Volume24H      string  `json:"volume_24h"`
	Amount24H      string  `json:"amount_24h"`
	Volume30D      string  `json:"volume_30d"`
	HighestBid     string  `json:"highest_bid"`
	LowestAsk      string  `json:"lowest_ask"`
	BidAskSpread   string  `json:"bid_ask_spread"`
	PercentageChange string `json:"percentage_change"`
}

// TradeService 交易服务
type TradeService struct {
	client *Client
}

// NewTradeService 创建交易服务
func NewTradeService() *TradeService {
	return &TradeService{
		client: GetClient(),
	}
}

// CreateOrder 创建订单
func (s *TradeService) CreateOrder(ctx context.Context, req *OrderRequest) (*Order, error) {
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
		"CB-ACCESS-KEY": s.client.config.APIKey,
	}

	// 构建请求体
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 调用Coinbase Trade API创建订单
	resp, err := s.client.httpClient.Post("/api/v3/brokerage/orders", body, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// 解析响应
	var order Order
	if err := json.Unmarshal(resp, &order); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	log.Printf("Coinbase order created: ID=%s, ProductID=%s, Side=%s", order.ID, order.ProductID, order.Side)
	return &order, nil
}

// GetOrder 获取订单详情
func (s *TradeService) GetOrder(ctx context.Context, orderID string) (*Order, error) {
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
		"CB-ACCESS-KEY": s.client.config.APIKey,
	}

	endpoint := fmt.Sprintf("/api/v3/brokerage/orders/historical/%s", orderID)
	resp, err := s.client.httpClient.Get(endpoint, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	var order Order
	if err := json.Unmarshal(resp, &order); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &order, nil
}

// ListOrders 列出订单
func (s *TradeService) ListOrders(ctx context.Context, productID string, limit, offset int) ([]Order, error) {
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
		"CB-ACCESS-KEY": s.client.config.APIKey,
	}

	endpoint := fmt.Sprintf("/api/v3/brokerage/orders/historical?limit=%d&offset=%d", limit, offset)
	if productID != "" {
		endpoint += "&product_id=" + productID
	}

	resp, err := s.client.httpClient.Get(endpoint, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	var orders []Order
	if err := json.Unmarshal(resp, &orders); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return orders, nil
}

// CancelOrder 取消订单
func (s *TradeService) CancelOrder(ctx context.Context, orderID string) error {
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
		"CB-ACCESS-KEY": s.client.config.APIKey,
	}

	endpoint := fmt.Sprintf("/api/v3/brokerage/orders/%s", orderID)
	_, err = s.client.httpClient.Delete(endpoint, headers)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}

// GetAccounts 获取账户列表
func (s *TradeService) GetAccounts(ctx context.Context) ([]Account, error) {
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
		"CB-ACCESS-KEY": s.client.config.APIKey,
	}

	resp, err := s.client.httpClient.Get("/api/v3/brokerage/accounts", headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	var accounts []Account
	if err := json.Unmarshal(resp, &accounts); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return accounts, nil
}

// GetProducts 获取产品列表
func (s *TradeService) GetProducts(ctx context.Context) ([]Product, error) {
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
		"CB-ACCESS-KEY": s.client.config.APIKey,
	}

	resp, err := s.client.httpClient.Get("/api/v3/brokerage/products", headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	var products []Product
	if err := json.Unmarshal(resp, &products); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return products, nil
}

// GetProduct 获取产品详情
func (s *TradeService) GetProduct(ctx context.Context, productID string) (*Product, error) {
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
		"CB-ACCESS-KEY": s.client.config.APIKey,
	}

	endpoint := fmt.Sprintf("/api/v3/brokerage/products/%s", productID)
	resp, err := s.client.httpClient.Get(endpoint, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	var product Product
	if err := json.Unmarshal(resp, &product); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &product, nil
}

// GetTicker 获取产品行情
func (s *TradeService) GetTicker(ctx context.Context, productID string) (*Ticker, error) {
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
		"CB-ACCESS-KEY": s.client.config.APIKey,
	}

	endpoint := fmt.Sprintf("/api/v3/brokerage/products/%s/ticker", productID)
	resp, err := s.client.httpClient.Get(endpoint, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticker: %w", err)
	}

	var ticker Ticker
	if err := json.Unmarshal(resp, &ticker); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &ticker, nil
}

// GetOrderBook 获取订单簿
func (s *TradeService) GetOrderBook(ctx context.Context, productID string) (map[string]interface{}, error) {
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
		"CB-ACCESS-KEY": s.client.config.APIKey,
	}

	endpoint := fmt.Sprintf("/api/v3/brokerage/product_book?product_id=%s", productID)
	resp, err := s.client.httpClient.Get(endpoint, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get order book: %w", err)
	}

	var orderBook map[string]interface{}
	if err := json.Unmarshal(resp, &orderBook); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return orderBook, nil
}
