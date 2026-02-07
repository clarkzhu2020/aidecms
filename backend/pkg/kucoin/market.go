package kucoin

import (
	"fmt"
	"net/url"
)

// Symbol 交易对信息
type Symbol struct {
	Symbol          string `json:"symbol"`
	Name            string `json:"name"`
	BaseCurrency    string `json:"baseCurrency"`
	QuoteCurrency   string `json:"quoteCurrency"`
	BaseMinSize     string `json:"baseMinSize"`
	QuoteMinSize    string `json:"quoteMinSize"`
	BaseMaxSize     string `json:"baseMaxSize"`
	QuoteMaxSize    string `json:"quoteMaxSize"`
	BaseIncrement   string `json:"baseIncrement"`
	QuoteIncrement string `json:"quoteIncrement"`
	PriceIncrement  string `json:"priceIncrement"`
	FeeCurrency     string `json:"feeCurrency"`
	EnableTrading   bool   `json:"enableTrading"`
	Market          string `json:"market"`
}

// GetSymbols 获取交易对列表
func (c *Client) GetSymbols(market string) ([]Symbol, error) {
	endpoint := "/api/v1/symbols"

	if market != "" {
		endpoint += "?market=" + market
	}

	respBody, err := c.httpClient.Get(endpoint, false)
	if err != nil {
		return nil, err
	}

	var symbols []Symbol
	if err := c.httpClient.unmarshalJSON(respBody, &symbols); err != nil {
		return nil, err
	}

	return symbols, nil
}

// Ticker Ticker信息
type Ticker struct {
	Symbol          string `json:"symbol"`
	SymbolName      string `json:"symbolName"`
	Buy             string `json:"buy"`
	Sell            string `json:"sell"`
	ChangeRate      string `json:"changeRate"`
	ChangePrice     string `json:"changePrice"`
	High            string `json:"high"`
	Low             string `json:"low"`
	Vol             string `json:"vol"`
	VolValue        string `json:"volValue"`
	Last            string `json:"last"`
	Open            string `json:"open"`
	Open24H         string `json:"open24h"`
	AveragePrice    string `json:"averagePrice"`
	UsdIndexPrice   string `json:"usdIndexPrice"`
	MarkPrice       string `json:"markPrice"`
	FeeRate         string `json:"feeRate"`
	BestBid         string `json:"bestBid"`
	BestAsk         string `json:"bestAsk"`
	BestBidSize     string `json:"bestBidSize"`
	BestAskSize     string `json:"bestAskSize"`
	OpenInterest    string `json:"openInterest"`
	OpenInterestValue string `json:"openInterestValue"`
}

// GetTicker 获取单个交易对的Ticker
func (c *Client) GetTicker(symbol string) (*Ticker, error) {
	endpoint := fmt.Sprintf("/api/v1/market/orderbook/level1?symbol=%s", url.QueryEscape(symbol))

	respBody, err := c.httpClient.Get(endpoint, false)
	if err != nil {
		return nil, err
	}

	var ticker Ticker
	if err := c.httpClient.unmarshalJSON(respBody, &ticker); err != nil {
		return nil, err
	}

	return &ticker, nil
}

// GetAllTickers 获取所有交易对的Ticker
func (c *Client) GetAllTickers(market string) (*Ticker, error) {
	endpoint := "/api/v1/market/allTickers"

	if market != "" {
		endpoint += "?market=" + market
	}

	respBody, err := c.httpClient.Get(endpoint, false)
	if err != nil {
		return nil, err
	}

	var response struct {
		Time     int64    `json:"time"`
		Timer    int64    `json:"timer"`
		Ticker   Ticker  `json:"ticker"`
		Tickers  []Ticker `json:"tickers"`
	}

	if err := c.httpClient.unmarshalJSON(respBody, &response); err != nil {
		return nil, err
	}

	// 返回概览ticker（包含总交易量等信息）
	return &response.Ticker, nil
}

// GetAllTickersDetail 获取所有交易对的Ticker详情列表
func (c *Client) GetAllTickersDetail(market string) ([]Ticker, error) {
	endpoint := "/api/v1/market/allTickers"

	if market != "" {
		endpoint += "?market=" + market
	}

	respBody, err := c.httpClient.Get(endpoint, false)
	if err != nil {
		return nil, err
	}

	var response struct {
		Time     int64    `json:"time"`
		Timer    int64    `json:"timer"`
		Ticker   Ticker  `json:"ticker"`
		Tickers  []Ticker `json:"tickers"`
	}

	if err := c.httpClient.unmarshalJSON(respBody, &response); err != nil {
		return nil, err
	}

	return response.Tickers, nil
}

// OrderBookLevel 订单簿级别
type OrderBookLevel struct [2]string
// OrderBookLevel[0] - 价格
// OrderBookLevel[1] - 数量

// OrderBook 订单簿
type OrderBook struct {
	Sequence string           `json:"sequence"`
	Bids     []OrderBookLevel `json:"bids"`
	Asks     []OrderBookLevel `json:"asks"`
}

// GetOrderBook 获取订单簿（部分）
func (c *Client) GetOrderBook(symbol string) (*OrderBook, error) {
	endpoint := fmt.Sprintf("/api/v1/market/orderbook/level2_20?symbol=%s", url.QueryEscape(symbol))

	respBody, err := c.httpClient.Get(endpoint, false)
	if err != nil {
		return nil, err
	}

	var orderBook OrderBook
	if err := c.httpClient.unmarshalJSON(respBody, &orderBook); err != nil {
		return nil, err
	}

	return &orderBook, nil
}

// GetOrderBookFull 获取完整订单簿
func (c *Client) GetOrderBookFull(symbol string) (*OrderBook, error) {
	endpoint := fmt.Sprintf("/api/v1/market/orderbook/level2?symbol=%s", url.QueryEscape(symbol))

	respBody, err := c.httpClient.Get(endpoint, false)
	if err != nil {
		return nil, err
	}

	var orderBook OrderBook
	if err := c.httpClient.unmarshalJSON(respBody, &orderBook); err != nil {
		return nil, err
	}

	return &orderBook, nil
}

// GetOrderBookByDepth 按深度获取订单簿
func (c *Client) GetOrderBookByDepth(symbol string, depth int) (*OrderBook, error) {
	endpoint := fmt.Sprintf("/api/v1/market/orderbook/level2_%d?symbol=%s", depth, url.QueryEscape(symbol))

	respBody, err := c.httpClient.Get(endpoint, false)
	if err != nil {
		return nil, err
	}

	var orderBook OrderBook
	if err := c.httpClient.unmarshalJSON(respBody, &orderBook); err != nil {
		return nil, err
	}

	return &orderBook, nil
}

// MarketTrade 市场成交记录
type MarketTrade struct {
	Sequence     string `json:"sequence"`
	Price        string `json:"price"`
	Size         string `json:"size"`
	Side         string `json:"side"`
	Time         int64  `json:"time"`
	TradeId      string `json:"tradeId"`
}

// GetMarketTrades 获取市场成交记录
func (c *Client) GetMarketTrades(symbol string) ([]MarketTrade, error) {
	endpoint := fmt.Sprintf("/api/v1/market/histories?symbol=%s", url.QueryEscape(symbol))

	respBody, err := c.httpClient.Get(endpoint, false)
	if err != nil {
		return nil, err
	}

	var trades []MarketTrade
	if err := c.httpClient.unmarshalJSON(respBody, &trades); err != nil {
		return nil, err
	}

	return trades, nil
}

// Kline K线数据
type Kline struct {
	Time     int64  `json:"0"`
	Open     string `json:"1"`
	Close    string `json:"2"`
	High     string `json:"3"`
	Low      string `json:"4"`
	Volume   string `json:"5"`
	Turnover string `json:"6"`
}

// GetKlines 获取K线数据
func (c *Client) GetKlines(symbol, klineType string, startAt, endAt int64) ([]Kline, error) {
	endpoint := "/api/v1/market/candles"

	params := url.Values{}
	params.Add("symbol", symbol)
	if klineType != "" {
		params.Add("type", klineType)
	}
	if startAt > 0 {
		params.Add("startAt", fmt.Sprintf("%d", startAt))
	}
	if endAt > 0 {
		params.Add("endAt", fmt.Sprintf("%d", endAt))
	}

	endpoint += "?" + params.Encode()

	respBody, err := c.httpClient.Get(endpoint, false)
	if err != nil {
		return nil, err
	}

	var klines []Kline
	if err := c.httpClient.unmarshalJSON(respBody, &klines); err != nil {
		return nil, err
	}

	return klines, nil
}

// Get24HStats 获取24小时统计数据
func (c *Client) Get24HStats(symbol string) (*Ticker, error) {
	endpoint := fmt.Sprintf("/api/v1/market/stats?symbol=%s", url.QueryEscape(symbol))

	respBody, err := c.httpClient.Get(endpoint, false)
	if err != nil {
		return nil, err
	}

	var ticker Ticker
	if err := c.httpClient.unmarshalJSON(respBody, &ticker); err != nil {
		return nil, err
	}

	return &ticker, nil
}
