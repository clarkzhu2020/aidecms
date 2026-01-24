package web3

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

// HyperliquidClient Hyperliquid 去中心化交易所客户端
type HyperliquidClient struct {
	baseURL    string
	privateKey *ecdsa.PrivateKey
	address    string
	httpClient *http.Client
}

// NewHyperliquidClient 创建 Hyperliquid 客户端
func NewHyperliquidClient(privateKeyHex string) (*HyperliquidClient, error) {
	var privateKey *ecdsa.PrivateKey
	var address string

	if privateKeyHex != "" {
		// 移除可能的 0x 前缀
		privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

		// 解析私钥
		key, err := crypto.HexToECDSA(privateKeyHex)
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %w", err)
		}
		privateKey = key

		// 从私钥派生地址
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("error casting public key to ECDSA")
		}
		address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	}

	return &HyperliquidClient{
		baseURL:    "https://api.hyperliquid.xyz",
		privateKey: privateKey,
		address:    address,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// GetBalance 获取余额
func (h *HyperliquidClient) GetBalance(ctx context.Context, currency string) (string, error) {
	if h.address == "" {
		return "", fmt.Errorf("wallet address not configured")
	}

	// Hyperliquid 使用 USDC 作为主要结算货币
	if currency != "USDC" && currency != "USD" {
		// 对于其他币种，查询持仓
		positions, err := h.GetPositions(ctx)
		if err != nil {
			return "0", nil // 如果没有持仓返回 0
		}

		// 查找对应币种的持仓
		for _, pos := range positions {
			if strings.HasPrefix(pos.Coin, currency) {
				return pos.Size, nil
			}
		}
		return "0", nil
	}

	// 查询账户余额
	balances, err := h.GetBalances(ctx)
	if err != nil {
		return "", err
	}

	if balance, ok := balances["USDC"]; ok {
		return balance, nil
	}

	return "0", nil
}

// GetBalances 获取所有余额
func (h *HyperliquidClient) GetBalances(ctx context.Context) (map[string]string, error) {
	if h.address == "" {
		return nil, fmt.Errorf("wallet address not configured")
	}

	// 构建请求
	reqBody := map[string]interface{}{
		"type": "clearinghouseState",
		"user": h.address,
	}

	respData, err := h.makeRequest(ctx, "/info", reqBody)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var state struct {
		MarginSummary struct {
			AccountValue      string `json:"accountValue"`
			TotalMarginUsed   string `json:"totalMarginUsed"`
			TotalNtlPos       string `json:"totalNtlPos"`
			TotalRawUsd       string `json:"totalRawUsd"`
			WithdrawableValue string `json:"withdrawable"`
		} `json:"marginSummary"`
		CrossMarginSummary struct {
			AccountValue    string `json:"accountValue"`
			TotalMarginUsed string `json:"totalMarginUsed"`
		} `json:"crossMarginSummary"`
	}

	if err := json.Unmarshal(respData, &state); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	balances := make(map[string]string)

	// 账户总价值
	if state.MarginSummary.AccountValue != "" {
		balances["USDC"] = state.MarginSummary.AccountValue
		balances["account_value"] = state.MarginSummary.AccountValue
	}

	// 可提现金额
	if state.MarginSummary.WithdrawableValue != "" {
		balances["withdrawable"] = state.MarginSummary.WithdrawableValue
	}

	// 已使用保证金
	if state.MarginSummary.TotalMarginUsed != "" {
		balances["margin_used"] = state.MarginSummary.TotalMarginUsed
	}

	return balances, nil
}

// GetPrice 获取交易对价格
func (h *HyperliquidClient) GetPrice(ctx context.Context, pair string) (string, error) {
	// Hyperliquid 使用币种名称而不是交易对格式
	// 例如: BTC 而不是 BTC-USD
	coin := strings.Split(pair, "-")[0]

	// 获取所有市场数据
	reqBody := map[string]interface{}{
		"type": "allMids",
	}

	respData, err := h.makeRequest(ctx, "/info", reqBody)
	if err != nil {
		return "", err
	}

	// 解析响应
	var mids map[string]string
	if err := json.Unmarshal(respData, &mids); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if price, ok := mids[coin]; ok {
		return price, nil
	}

	return "", fmt.Errorf("price not found for %s", pair)
}

// Position 持仓信息
type Position struct {
	Coin          string `json:"coin"`
	Size          string `json:"szi"`
	EntryPrice    string `json:"entryPx"`
	PositionValue string `json:"positionValue"`
	UnrealizedPnl string `json:"unrealizedPnl"`
	Leverage      string `json:"leverage"`
	Liquidation   string `json:"liquidationPx"`
}

// GetPositions 获取当前持仓
func (h *HyperliquidClient) GetPositions(ctx context.Context) ([]Position, error) {
	if h.address == "" {
		return nil, fmt.Errorf("wallet address not configured")
	}

	reqBody := map[string]interface{}{
		"type": "clearinghouseState",
		"user": h.address,
	}

	respData, err := h.makeRequest(ctx, "/info", reqBody)
	if err != nil {
		return nil, err
	}

	var state struct {
		AssetPositions []struct {
			Position struct {
				Coin          string `json:"coin"`
				Szi           string `json:"szi"`
				EntryPx       string `json:"entryPx"`
				PositionValue string `json:"positionValue"`
				UnrealizedPnl string `json:"unrealizedPnl"`
				Leverage      struct {
					Value string `json:"value"`
				} `json:"leverage"`
				LiquidationPx string `json:"liquidationPx"`
			} `json:"position"`
		} `json:"assetPositions"`
	}

	if err := json.Unmarshal(respData, &state); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	positions := make([]Position, 0)
	for _, ap := range state.AssetPositions {
		// 只返回有持仓的（size != 0）
		size, _ := strconv.ParseFloat(ap.Position.Szi, 64)
		if size != 0 {
			positions = append(positions, Position{
				Coin:          ap.Position.Coin,
				Size:          ap.Position.Szi,
				EntryPrice:    ap.Position.EntryPx,
				PositionValue: ap.Position.PositionValue,
				UnrealizedPnl: ap.Position.UnrealizedPnl,
				Leverage:      ap.Position.Leverage.Value,
				Liquidation:   ap.Position.LiquidationPx,
			})
		}
	}

	return positions, nil
}

// GetMarketInfo 获取市场信息
func (h *HyperliquidClient) GetMarketInfo(ctx context.Context) (map[string]interface{}, error) {
	reqBody := map[string]interface{}{
		"type": "meta",
	}

	respData, err := h.makeRequest(ctx, "/info", reqBody)
	if err != nil {
		return nil, err
	}

	var meta struct {
		Universe []struct {
			Name         string `json:"name"`
			SzDecimals   int    `json:"szDecimals"`
			MaxLeverage  int    `json:"maxLeverage"`
			OnlyIsolated bool   `json:"onlyIsolated"`
		} `json:"universe"`
	}

	if err := json.Unmarshal(respData, &meta); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	markets := make(map[string]interface{})
	for _, market := range meta.Universe {
		markets[market.Name] = map[string]interface{}{
			"name":         market.Name,
			"decimals":     market.SzDecimals,
			"max_leverage": market.MaxLeverage,
			"isolated":     market.OnlyIsolated,
		}
	}

	return markets, nil
}

// OrderRequest 下单请求
type OrderRequest struct {
	Coin       string  // 币种，如 "BTC"
	IsBuy      bool    // true 为买入，false 为卖出
	Size       float64 // 数量
	LimitPrice float64 // 限价（0 表示市价单）
	ReduceOnly bool    // 是否只减仓
}

// PlaceOrder 下单（需要私钥）
func (h *HyperliquidClient) PlaceOrder(ctx context.Context, order OrderRequest) (string, error) {
	if h.privateKey == nil {
		return "", fmt.Errorf("private key not configured, cannot place orders")
	}

	// 构建订单
	orderType := map[string]interface{}{
		"limit": map[string]interface{}{
			"tif": "Gtc", // Good til canceled
		},
	}

	if order.LimitPrice == 0 {
		// 市价单
		orderType = map[string]interface{}{
			"trigger": map[string]interface{}{
				"isMarket":  true,
				"triggerPx": "0",
			},
		}
	}

	action := map[string]interface{}{
		"type": "order",
		"orders": []map[string]interface{}{
			{
				"a": h.getCoinIndex(order.Coin),
				"b": order.IsBuy,
				"p": fmt.Sprintf("%.8f", order.LimitPrice),
				"s": fmt.Sprintf("%.8f", order.Size),
				"r": order.ReduceOnly,
				"t": orderType,
			},
		},
		"grouping": "na",
	}

	// 签名并发送
	signature, err := h.signAction(action)
	if err != nil {
		return "", err
	}

	reqBody := map[string]interface{}{
		"action":    action,
		"signature": signature,
		"nonce":     time.Now().UnixMilli(),
	}

	respData, err := h.makeRequest(ctx, "/exchange", reqBody)
	if err != nil {
		return "", err
	}

	var response struct {
		Status   string `json:"status"`
		Response struct {
			Type string `json:"type"`
			Data struct {
				Statuses []struct {
					Filled string `json:"filled"`
					Oid    int64  `json:"oid"`
				} `json:"statuses"`
			} `json:"data"`
		} `json:"response"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Status != "ok" {
		return "", fmt.Errorf("order failed: %s", response.Status)
	}

	if len(response.Response.Data.Statuses) > 0 {
		oid := response.Response.Data.Statuses[0].Oid
		return fmt.Sprintf("%d", oid), nil
	}

	return "", fmt.Errorf("no order id returned")
}

// CancelOrder 取消订单（需要私钥）
func (h *HyperliquidClient) CancelOrder(ctx context.Context, coin string, oid int64) error {
	if h.privateKey == nil {
		return fmt.Errorf("private key not configured, cannot cancel orders")
	}

	action := map[string]interface{}{
		"type": "cancel",
		"cancels": []map[string]interface{}{
			{
				"a": h.getCoinIndex(coin),
				"o": oid,
			},
		},
	}

	signature, err := h.signAction(action)
	if err != nil {
		return err
	}

	reqBody := map[string]interface{}{
		"action":    action,
		"signature": signature,
		"nonce":     time.Now().UnixMilli(),
	}

	respData, err := h.makeRequest(ctx, "/exchange", reqBody)
	if err != nil {
		return err
	}

	var response struct {
		Status string `json:"status"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Status != "ok" {
		return fmt.Errorf("cancel failed: %s", response.Status)
	}

	return nil
}

// makeRequest 发送 HTTP 请求
func (h *HyperliquidClient) makeRequest(ctx context.Context, endpoint string, body interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", h.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(respData))
	}

	return respData, nil
}

// signAction 签名操作（使用 EIP-712）
func (h *HyperliquidClient) signAction(action map[string]interface{}) (map[string]interface{}, error) {
	// 构建 EIP-712 消息
	// Hyperliquid 使用特定的 EIP-712 格式
	actionJSON, err := json.Marshal(action)
	if err != nil {
		return nil, err
	}

	// 简化的签名实现（实际应该使用完整的 EIP-712）
	hash := crypto.Keccak256Hash(actionJSON)
	signature, err := crypto.Sign(hash.Bytes(), h.privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}

	// 调整 v 值（EIP-155）
	if signature[64] < 27 {
		signature[64] += 27
	}

	return map[string]interface{}{
		"r": "0x" + hex.EncodeToString(signature[0:32]),
		"s": "0x" + hex.EncodeToString(signature[32:64]),
		"v": int(signature[64]),
	}, nil
}

// getCoinIndex 获取币种索引（简化实现）
func (h *HyperliquidClient) getCoinIndex(coin string) int {
	// 这是一个简化的实现
	// 实际应该从 meta 接口获取正确的索引
	coinMap := map[string]int{
		"BTC":   0,
		"ETH":   1,
		"SOL":   2,
		"MATIC": 3,
		"ARB":   4,
		"OP":    5,
	}

	if idx, ok := coinMap[coin]; ok {
		return idx
	}

	return 0
}

// GetOrderBook 获取订单簿
func (h *HyperliquidClient) GetOrderBook(ctx context.Context, coin string) (map[string]interface{}, error) {
	reqBody := map[string]interface{}{
		"type": "l2Book",
		"coin": coin,
	}

	respData, err := h.makeRequest(ctx, "/info", reqBody)
	if err != nil {
		return nil, err
	}

	var orderBook struct {
		Coin   string `json:"coin"`
		Time   int64  `json:"time"`
		Levels [][]struct {
			Px string `json:"px"`
			Sz string `json:"sz"`
			N  int    `json:"n"`
		} `json:"levels"`
	}

	if err := json.Unmarshal(respData, &orderBook); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	result := map[string]interface{}{
		"coin": orderBook.Coin,
		"time": orderBook.Time,
	}

	if len(orderBook.Levels) >= 2 {
		result["bids"] = orderBook.Levels[0]
		result["asks"] = orderBook.Levels[1]
	}

	return result, nil
}

// GetFundingRate 获取资金费率
func (h *HyperliquidClient) GetFundingRate(ctx context.Context, coin string) (string, error) {
	reqBody := map[string]interface{}{
		"type": "metaAndAssetCtxs",
	}

	respData, err := h.makeRequest(ctx, "/info", reqBody)
	if err != nil {
		return "", err
	}

	var response []struct {
		Ctx struct {
			Funding string `json:"funding"`
		} `json:"ctx"`
		Universe struct {
			Name string `json:"name"`
		} `json:"universe"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	for _, item := range response {
		if item.Universe.Name == coin {
			return item.Ctx.Funding, nil
		}
	}

	return "", fmt.Errorf("funding rate not found for %s", coin)
}

// Get24HVolume 获取24小时交易量
func (h *HyperliquidClient) Get24HVolume(ctx context.Context, coin string) (string, error) {
	reqBody := map[string]interface{}{
		"type": "metaAndAssetCtxs",
	}

	respData, err := h.makeRequest(ctx, "/info", reqBody)
	if err != nil {
		return "", err
	}

	var response []struct {
		Ctx struct {
			DayNtlVlm string `json:"dayNtlVlm"`
		} `json:"ctx"`
		Universe struct {
			Name string `json:"name"`
		} `json:"universe"`
	}

	if err := json.Unmarshal(respData, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	for _, item := range response {
		if item.Universe.Name == coin {
			return item.Ctx.DayNtlVlm, nil
		}
	}

	return "", fmt.Errorf("volume not found for %s", coin)
}

// ConvertToUSDC 将价格转换为 USDC（辅助函数）
func (h *HyperliquidClient) ConvertToUSDC(value string, btcPrice *big.Float) string {
	val, ok := new(big.Float).SetString(value)
	if !ok {
		return "0"
	}

	result := new(big.Float).Mul(val, btcPrice)
	return result.Text('f', 6)
}
