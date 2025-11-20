package web3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BitcoinClient Bitcoin 客户端
type BitcoinClient struct {
	rpcURL     string
	apiKey     string
	httpClient *http.Client
}

// BitcoinRPCRequest RPC 请求
type BitcoinRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// BitcoinRPCResponse RPC 响应
type BitcoinRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error,omitempty"`
}

// RPCError RPC 错误
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewBitcoinClient 创建 Bitcoin 客户端
func NewBitcoinClient(rpcURL, apiKey string) *BitcoinClient {
	return &BitcoinClient{
		rpcURL: rpcURL,
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// call RPC 调用
func (c *BitcoinClient) call(ctx context.Context, method string, params []interface{}) (json.RawMessage, error) {
	req := BitcoinRPCRequest{
		JSONRPC: "2.0",
		ID:      "aidecms",
		Method:  method,
		Params:  params,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.rpcURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rpcResp BitcoinRPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	return rpcResp.Result, nil
}

// GetBalance 获取地址余额
func (c *BitcoinClient) GetBalance(ctx context.Context, address string) (string, error) {
	if err := ValidateAddress(Bitcoin, address); err != nil {
		return "", err
	}

	// Note: Bitcoin Core doesn't have a direct address balance query
	// This would typically require a block explorer API or indexer
	result, err := c.call(ctx, "getbalance", []interface{}{})
	if err != nil {
		return "", err
	}

	var balance float64
	if err := json.Unmarshal(result, &balance); err != nil {
		return "", fmt.Errorf("failed to parse balance: %w", err)
	}

	return fmt.Sprintf("%.8f", balance), nil
}

// GetBlockNumber 获取最新区块高度
func (c *BitcoinClient) GetBlockNumber(ctx context.Context) (uint64, error) {
	result, err := c.call(ctx, "getblockcount", []interface{}{})
	if err != nil {
		return 0, err
	}

	var height uint64
	if err := json.Unmarshal(result, &height); err != nil {
		return 0, fmt.Errorf("failed to parse block height: %w", err)
	}

	return height, nil
}

// GetTransaction 获取交易信息
func (c *BitcoinClient) GetTransaction(ctx context.Context, txHash string) (*Transaction, error) {
	if err := ValidateTxHash(Bitcoin, txHash); err != nil {
		return nil, err
	}

	result, err := c.call(ctx, "getrawtransaction", []interface{}{txHash, true})
	if err != nil {
		return nil, err
	}

	var btcTx struct {
		TxID          string `json:"txid"`
		Hash          string `json:"hash"`
		Size          int    `json:"size"`
		VSize         int    `json:"vsize"`
		Version       int    `json:"version"`
		Locktime      int64  `json:"locktime"`
		Confirmations int    `json:"confirmations"`
		BlockHash     string `json:"blockhash"`
		BlockHeight   uint64 `json:"blockheight"`
		Time          int64  `json:"time"`
		Vin           []struct {
			TxID string `json:"txid"`
			Vout int    `json:"vout"`
		} `json:"vin"`
		Vout []struct {
			Value        float64 `json:"value"`
			N            int     `json:"n"`
			ScriptPubKey struct {
				Address string `json:"address"`
			} `json:"scriptPubKey"`
		} `json:"vout"`
	}

	if err := json.Unmarshal(result, &btcTx); err != nil {
		return nil, fmt.Errorf("failed to parse transaction: %w", err)
	}

	tx := &Transaction{
		Hash:        btcTx.TxID,
		BlockHash:   btcTx.BlockHash,
		BlockNumber: btcTx.BlockHeight,
		Timestamp:   btcTx.Time,
		Status:      "confirmed",
		Extra:       make(map[string]interface{}),
	}

	if btcTx.Confirmations == 0 {
		tx.Status = "pending"
	}

	// Get first output address and value
	if len(btcTx.Vout) > 0 {
		tx.To = btcTx.Vout[0].ScriptPubKey.Address
		tx.Value = fmt.Sprintf("%.8f", btcTx.Vout[0].Value)
	}

	// Store full transaction data
	tx.Extra["vsize"] = btcTx.VSize
	tx.Extra["confirmations"] = btcTx.Confirmations
	tx.Extra["vin_count"] = len(btcTx.Vin)
	tx.Extra["vout_count"] = len(btcTx.Vout)

	return tx, nil
}

// SendTransaction 发送交易
func (c *BitcoinClient) SendTransaction(ctx context.Context, tx *TransactionRequest) (string, error) {
	// Note: Bitcoin requires raw transaction hex
	// This is a placeholder implementation
	return "", fmt.Errorf("sendTransaction not implemented: requires raw transaction construction")
}

// GetChain 获取链类型
func (c *BitcoinClient) GetChain() Chain {
	return Bitcoin
}

// Close 关闭连接
func (c *BitcoinClient) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

// GetBlockHash 获取区块哈希
func (c *BitcoinClient) GetBlockHash(ctx context.Context, height uint64) (string, error) {
	result, err := c.call(ctx, "getblockhash", []interface{}{height})
	if err != nil {
		return "", err
	}

	var hash string
	if err := json.Unmarshal(result, &hash); err != nil {
		return "", fmt.Errorf("failed to parse block hash: %w", err)
	}

	return hash, nil
}

// GetBlock 获取区块信息
func (c *BitcoinClient) GetBlock(ctx context.Context, blockHash string) (map[string]interface{}, error) {
	result, err := c.call(ctx, "getblock", []interface{}{blockHash})
	if err != nil {
		return nil, err
	}

	var block map[string]interface{}
	if err := json.Unmarshal(result, &block); err != nil {
		return nil, fmt.Errorf("failed to parse block: %w", err)
	}

	return block, nil
}

// GetMempoolInfo 获取内存池信息
func (c *BitcoinClient) GetMempoolInfo(ctx context.Context) (map[string]interface{}, error) {
	result, err := c.call(ctx, "getmempoolinfo", []interface{}{})
	if err != nil {
		return nil, err
	}

	var info map[string]interface{}
	if err := json.Unmarshal(result, &info); err != nil {
		return nil, fmt.Errorf("failed to parse mempool info: %w", err)
	}

	return info, nil
}
