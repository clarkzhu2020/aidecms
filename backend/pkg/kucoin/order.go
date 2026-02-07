package kucoin

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// OrderType 订单类型
type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"
	OrderTypeStop   OrderType = "stop"
	OrderTypeStopLimit OrderType = "stop_limit"
)

// OrderSide 订单方向
type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

// OrderTimeInForce 订单有效期
type OrderTimeInForce string

const (
	TimeInForceGTC OrderTimeInForce = "GTC" // Good Till Cancel
	TimeInForceIOC OrderTimeInForce = "IOC" // Immediate or Cancel
	TimeInForceFOK OrderTimeInForce = "FOK" // Fill or Kill
)

// OrderCreateRequest 创建订单请求
type OrderCreateRequest struct {
	ClientOid    string           `json:"clientOid"`    // 客户端订单ID
	Symbol       string           `json:"symbol"`       // 交易对
	Side         OrderSide        `json:"side"`         // 买卖方向
	Type         OrderType        `json:"type"`         // 订单类型
	Price        string           `json:"price"`        // 价格（limit订单）
	Size         string           `json:"size"`         // 数量
	TimeInForce  OrderTimeInForce `json:"timeInForce"`  // 订单有效期
	Stop         string           `json:"stop"`         // 止损类型 (loss, entry)
	StopPrice    string           `json:"stopPrice"`    // 止损价格
	TradeType    string           `json:"tradeType"`    // 交易类型 (TRADE-SPOT, TRADE_MARGIN)
	Stp          string           `json:"stp"`          // 自我交易保护
	SelfTradePrevention string    `json:"selfTradePrevention"` // ST, CN, CB, DC
	Remark       string           `json:"remark"`       // 备注
}

// OrderResponse 订单响应
type OrderResponse struct {
	OrderId     string `json:"orderId"`
	ClientOid   string `json:"clientOid"`
	Symbol      string `json:"symbol"`
	Side        string `json:"side"`
	Type        string `json:"type"`
	Price       string `json:"price"`
	Size        string `json:"size"`
	DealSize    string `json:"dealSize"`
	DealFunds   string `json:"dealFunds"`
	Fee         string `json:"fee"`
	FeeCurrency string `json:"feeCurrency"`
	Stp         string `json:"stp"`
	SelfTradePrevention string `json:"selfTradePrevention"`
	Stop        string `json:"stop"`
	StopTrigger string `json:"stopTrigger"`
	StopPrice   string `json:"stopPrice"`
	TimeInForce string `json:"timeInForce"`
	PostOnly    bool   `json:"postOnly"`
	Hidden      bool   `json:"hidden"`
	Iceberg     bool   `json:"iceberg"`
	VisibleSize string `json:"visibleSize"`
	CancelAfter int64  `json:"cancelAfter"`
	Channel     string `json:"channel"`
	ReturnTrade bool   `json:"returnTrade"`
	TradeType   string `json:"tradeType"`
	Remark      string `json:"remark"`
	CreatedAt   int64  `json:"createdAt"`
	TradeId     string `json:"tradeId"`
}

// CreateOrder 创建订单
func (c *Client) CreateOrder(req *OrderCreateRequest) (*OrderResponse, error) {
	endpoint := "/api/v1/orders"

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	respBody, err := c.httpClient.Post(endpoint, body, true)
	if err != nil {
		return nil, err
	}

	var response OrderResponse
	if err := c.httpClient.unmarshalJSON(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateOrderSync 同步创建订单（返回订单详情）
func (c *Client) CreateOrderSync(req *OrderCreateRequest) (*OrderResponse, error) {
	endpoint := "/api/v1/orders/sync"

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	respBody, err := c.httpClient.Post(endpoint, body, true)
	if err != nil {
		return nil, err
	}

	var response OrderResponse
	if err := c.httpClient.unmarshalJSON(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// TestOrder 测试订单（不实际下单）
func (c *Client) TestOrder(req *OrderCreateRequest) error {
	endpoint := "/api/v1/orders/test"

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	_, err = c.httpClient.Post(endpoint, body, true)
	if err != nil {
		return err
	}

	return nil
}

// CancelOrderById 通过订单ID取消订单
func (c *Client) CancelOrderById(orderId string) error {
	endpoint := fmt.Sprintf("/api/v1/orders/%s", orderId)

	_, err := c.httpClient.Delete(endpoint, true)
	if err != nil {
		return err
	}

	return nil
}

// CancelOrderByClientOid 通过客户端订单ID取消订单
func (c *Client) CancelOrderByClientOid(clientOid string) error {
	endpoint := fmt.Sprintf("/api/v1/orders/client-order/%s", clientOid)

	_, err := c.httpClient.Delete(endpoint, true)
	if err != nil {
		return err
	}

	return nil
}

// CancelAllOrdersBySymbol 取消指定交易对的所有订单
func (c *Client) CancelAllOrdersBySymbol(symbol string) error {
	endpoint := fmt.Sprintf("/api/v1/orders?symbol=%s", url.QueryEscape(symbol))

	_, err := c.httpClient.Delete(endpoint, true)
	if err != nil {
		return err
	}

	return nil
}

// CancelAllOrders 取消所有订单
func (c *Client) CancelAllOrders() error {
	endpoint := "/api/v1/orders"

	_, err := c.httpClient.Delete(endpoint, true)
	if err != nil {
		return err
	}

	return nil
}

// OrderDetail 订单详情
type OrderDetail struct {
	OrderId     string `json:"id"`
	ClientOid   string `json:"clientOid"`
	Symbol      string `json:"symbol"`
	Side        string `json:"side"`
	Type        string `json:"type"`
	KType       string `json:"kType"`
	Price       string `json:"price"`
	Size        string `json:"size"`
	DealFunds   string `json:"dealFunds"`
	DealSize    string `json:"dealSize"`
	Fee         string `json:"fee"`
	FeeCurrency string `json:"feeCurrency"`
	Stp         string `json:"stp"`
	Stop        string `json:"stop"`
	StopTrigger string `json:"stopTrigger"`
	StopPrice   string `json:"stopPrice"`
	TimeInForce string `json:"timeInForce"`
	PostOnly    bool   `json:"postOnly"`
	Hidden      bool   `json:"hidden"`
	Iceberg     bool   `json:"iceberg"`
	VisibleSize string `json:"visibleSize"`
	CancelAfter int64  `json:"cancelAfter"`
	Channel     string `json:"channel"`
	TradeType   string `json:"tradeType"`
	Remark      string `json:"remark"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
}

// GetOrderById 通过订单ID获取订单详情
func (c *Client) GetOrderById(orderId string) (*OrderDetail, error) {
	endpoint := fmt.Sprintf("/api/v1/orders/%s", orderId)

	respBody, err := c.httpClient.Get(endpoint, true)
	if err != nil {
		return nil, err
	}

	var order OrderDetail
	if err := c.httpClient.unmarshalJSON(respBody, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOrderByClientOid 通过客户端订单ID获取订单详情
func (c *Client) GetOrderByClientOid(clientOid string) (*OrderDetail, error) {
	endpoint := fmt.Sprintf("/api/v1/orders/client-order/%s", clientOid)

	respBody, err := c.httpClient.Get(endpoint, true)
	if err != nil {
		return nil, err
	}

	var order OrderDetail
	if err := c.httpClient.unmarshalJSON(respBody, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOpenOrders 获取未成交订单
func (c *Client) GetOpenOrders(symbol string) ([]OrderDetail, error) {
	endpoint := "/api/v1/orders/open"

	if symbol != "" {
		endpoint += "?symbol=" + url.QueryEscape(symbol)
	}

	respBody, err := c.httpClient.Get(endpoint, true)
	if err != nil {
		return nil, err
	}

	var orders []OrderDetail
	if err := c.httpClient.unmarshalJSON(respBody, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetClosedOrders 获取已成交订单历史
func (c *Client) GetClosedOrders(symbol, status string, currentPage, pageSize int) ([]OrderDetail, error) {
	endpoint := "/api/v1/orders/closed"

	params := url.Values{}
	if symbol != "" {
		params.Add("symbol", symbol)
	}
	if status != "" {
		params.Add("status", status)
	}
	if currentPage > 0 {
		params.Add("currentPage", fmt.Sprintf("%d", currentPage))
	}
	if pageSize > 0 {
		params.Add("pageSize", fmt.Sprintf("%d", pageSize))
	}

	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	respBody, err := c.httpClient.Get(endpoint, true)
	if err != nil {
		return nil, err
	}

	var response struct {
		CurrentPage int          `json:"currentPage"`
		PageSize    int          `json:"pageSize"`
		TotalNum    int          `json:"totalNum"`
		TotalPage   int          `json:"totalPage"`
		Items       []OrderDetail `json:"items"`
	}

	if err := c.httpClient.unmarshalJSON(respBody, &response); err != nil {
		return nil, err
	}

	return response.Items, nil
}

// Trade 成交记录
type Trade struct {
	TradeId     string `json:"tradeId"`
	OrderId     string `json:"orderId"`
	Symbol      string `json:"symbol"`
	Side        string `json:"side"`
	Price       string `json:"price"`
	Size        string `json:"size"`
	Fee         string `json:"fee"`
	FeeCurrency string `json:"feeCurrency"`
	CounterOrderId string `json:"counterOrderId"`
	ForceTaker  bool   `json:"forceTaker"`
	Liquidity   string `json:"liquidity"`
	CreatedAt   int64  `json:"createdAt"`
}

// GetTradeHistory 获取成交历史
func (c *Client) GetTradeHistory(orderId string) ([]Trade, error) {
	endpoint := "/api/v1/orders/trade/history"

	if orderId != "" {
		endpoint += "?orderId=" + orderId
	}

	respBody, err := c.httpClient.Get(endpoint, true)
	if err != nil {
		return nil, err
	}

	var trades []Trade
	if err := c.httpClient.unmarshalJSON(respBody, &trades); err != nil {
		return nil, err
	}

	return trades, nil
}
