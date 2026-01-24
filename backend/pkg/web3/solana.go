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

// SolanaClient Solana 客户端
type SolanaClient struct {
	rpcURL     string
	httpClient *http.Client
}

// SolanaRPCRequest Solana RPC 请求
type SolanaRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params,omitempty"`
}

// SolanaRPCResponse Solana RPC 响应
type SolanaRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *RPCError       `json:"error,omitempty"`
}

// NewSolanaClient 创建 Solana 客户端
func NewSolanaClient(rpcURL string) *SolanaClient {
	return &SolanaClient{
		rpcURL: rpcURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// call RPC 调用
func (c *SolanaClient) call(ctx context.Context, method string, params []interface{}) (json.RawMessage, error) {
	req := SolanaRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
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

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rpcResp SolanaRPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	return rpcResp.Result, nil
}

// GetBalance 获取地址余额（单位：lamports）
func (c *SolanaClient) GetBalance(ctx context.Context, address string) (string, error) {
	if err := ValidateAddress(Solana, address); err != nil {
		return "", err
	}

	result, err := c.call(ctx, "getBalance", []interface{}{address})
	if err != nil {
		return "", err
	}

	var resp struct {
		Value uint64 `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse balance: %w", err)
	}

	return fmt.Sprintf("%d", resp.Value), nil
}

// GetBalanceInSOL 获取余额（以 SOL 为单位）
func (c *SolanaClient) GetBalanceInSOL(ctx context.Context, address string) (string, error) {
	balance, err := c.GetBalance(ctx, address)
	if err != nil {
		return "", err
	}

	var lamports uint64
	fmt.Sscanf(balance, "%d", &lamports)

	// Convert lamports to SOL (1 SOL = 10^9 lamports)
	sol := float64(lamports) / 1e9
	return fmt.Sprintf("%.9f", sol), nil
}

// GetBlockNumber 获取最新区块高度（slot）
func (c *SolanaClient) GetBlockNumber(ctx context.Context) (uint64, error) {
	result, err := c.call(ctx, "getSlot", []interface{}{})
	if err != nil {
		return 0, err
	}

	var slot uint64
	if err := json.Unmarshal(result, &slot); err != nil {
		return 0, fmt.Errorf("failed to parse slot: %w", err)
	}

	return slot, nil
}

// GetTransaction 获取交易信息
func (c *SolanaClient) GetTransaction(ctx context.Context, signature string) (*Transaction, error) {
	if err := ValidateTxHash(Solana, signature); err != nil {
		return nil, err
	}

	params := []interface{}{
		signature,
		map[string]interface{}{
			"encoding": "json",
		},
	}

	result, err := c.call(ctx, "getTransaction", params)
	if err != nil {
		return nil, err
	}

	var solTx struct {
		Slot      uint64 `json:"slot"`
		BlockTime int64  `json:"blockTime"`
		Meta      *struct {
			Err          interface{} `json:"err"`
			Fee          uint64      `json:"fee"`
			PreBalances  []uint64    `json:"preBalances"`
			PostBalances []uint64    `json:"postBalances"`
		} `json:"meta"`
		Transaction struct {
			Message struct {
				AccountKeys []string `json:"accountKeys"`
			} `json:"message"`
			Signatures []string `json:"signatures"`
		} `json:"transaction"`
	}

	if err := json.Unmarshal(result, &solTx); err != nil {
		return nil, fmt.Errorf("failed to parse transaction: %w", err)
	}

	tx := &Transaction{
		Hash:        signature,
		BlockNumber: solTx.Slot,
		Timestamp:   solTx.BlockTime,
		Status:      "confirmed",
		Extra:       make(map[string]interface{}),
	}

	if solTx.Meta != nil {
		if solTx.Meta.Err != nil {
			tx.Status = "failed"
		} else {
			tx.Status = "success"
		}

		tx.GasUsed = solTx.Meta.Fee

		// Calculate value transferred
		if len(solTx.Meta.PreBalances) > 0 && len(solTx.Meta.PostBalances) > 0 {
			if solTx.Meta.PreBalances[0] > solTx.Meta.PostBalances[0] {
				value := solTx.Meta.PreBalances[0] - solTx.Meta.PostBalances[0]
				tx.Value = fmt.Sprintf("%d", value)
			}
		}
	}

	// Get from/to addresses
	if len(solTx.Transaction.Message.AccountKeys) > 0 {
		tx.From = solTx.Transaction.Message.AccountKeys[0]
	}
	if len(solTx.Transaction.Message.AccountKeys) > 1 {
		tx.To = solTx.Transaction.Message.AccountKeys[1]
	}

	// Store signature
	if len(solTx.Transaction.Signatures) > 0 {
		tx.Extra["signature"] = solTx.Transaction.Signatures[0]
	}

	return tx, nil
}

// SendTransaction 发送交易
func (c *SolanaClient) SendTransaction(ctx context.Context, tx *TransactionRequest) (string, error) {
	// Note: Solana requires serialized transaction
	return "", fmt.Errorf("sendTransaction not implemented: requires transaction serialization")
}

// GetChain 获取链类型
func (c *SolanaClient) GetChain() Chain {
	return Solana
}

// Close 关闭连接
func (c *SolanaClient) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

// GetVersion 获取 Solana 版本
func (c *SolanaClient) GetVersion(ctx context.Context) (map[string]interface{}, error) {
	result, err := c.call(ctx, "getVersion", []interface{}{})
	if err != nil {
		return nil, err
	}

	var version map[string]interface{}
	if err := json.Unmarshal(result, &version); err != nil {
		return nil, fmt.Errorf("failed to parse version: %w", err)
	}

	return version, nil
}

// GetBlockHeight 获取区块高度
func (c *SolanaClient) GetBlockHeight(ctx context.Context) (uint64, error) {
	result, err := c.call(ctx, "getBlockHeight", []interface{}{})
	if err != nil {
		return 0, err
	}

	var height uint64
	if err := json.Unmarshal(result, &height); err != nil {
		return 0, fmt.Errorf("failed to parse block height: %w", err)
	}

	return height, nil
}

// GetRecentBlockhash 获取最近的区块哈希
func (c *SolanaClient) GetRecentBlockhash(ctx context.Context) (string, error) {
	result, err := c.call(ctx, "getRecentBlockhash", []interface{}{})
	if err != nil {
		return "", err
	}

	var resp struct {
		Value struct {
			Blockhash string `json:"blockhash"`
		} `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse blockhash: %w", err)
	}

	return resp.Value.Blockhash, nil
}

// GetAccountInfo 获取账户信息
func (c *SolanaClient) GetAccountInfo(ctx context.Context, address string) (map[string]interface{}, error) {
	if err := ValidateAddress(Solana, address); err != nil {
		return nil, err
	}

	params := []interface{}{
		address,
		map[string]interface{}{
			"encoding": "jsonParsed",
		},
	}

	result, err := c.call(ctx, "getAccountInfo", params)
	if err != nil {
		return nil, err
	}

	var info map[string]interface{}
	if err := json.Unmarshal(result, &info); err != nil {
		return nil, fmt.Errorf("failed to parse account info: %w", err)
	}

	return info, nil
}

// GetTokenBalance 获取 SPL Token 余额
func (c *SolanaClient) GetTokenBalance(ctx context.Context, tokenAccount string) (string, error) {
	result, err := c.call(ctx, "getTokenAccountBalance", []interface{}{tokenAccount})
	if err != nil {
		return "", err
	}

	var resp struct {
		Value struct {
			Amount         string `json:"amount"`
			Decimals       int    `json:"decimals"`
			UIAmountString string `json:"uiAmountString"`
		} `json:"value"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse token balance: %w", err)
	}

	return resp.Value.UIAmountString, nil
}
