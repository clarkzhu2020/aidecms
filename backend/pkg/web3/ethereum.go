package web3

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// EthereumClient Ethereum/BSC 客户端
type EthereumClient struct {
	client *ethclient.Client
	rpc    *rpc.Client
	chain  Chain
}

// NewEthereumClient 创建 Ethereum 客户端
func NewEthereumClient(rpcURL string) (*EthereumClient, error) {
	return newEVMClient(rpcURL, Ethereum)
}

// NewBSCClient 创建 BSC 客户端
func NewBSCClient(rpcURL string) (*EthereumClient, error) {
	return newEVMClient(rpcURL, BSC)
}

func newEVMClient(rpcURL string, chain Chain) (*EthereumClient, error) {
	rpcClient, err := rpc.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", chain, err)
	}

	client := ethclient.NewClient(rpcClient)
	return &EthereumClient{
		client: client,
		rpc:    rpcClient,
		chain:  chain,
	}, nil
}

// GetBalance 获取地址余额
func (c *EthereumClient) GetBalance(ctx context.Context, address string) (string, error) {
	if err := ValidateAddress(c.chain, address); err != nil {
		return "", err
	}

	addr := common.HexToAddress(address)
	balance, err := c.client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get balance: %w", err)
	}

	return balance.String(), nil
}

// GetBalanceInEther 获取余额（以 Ether 为单位）
func (c *EthereumClient) GetBalanceInEther(ctx context.Context, address string) (string, error) {
	balance, err := c.GetBalance(ctx, address)
	if err != nil {
		return "", err
	}

	wei := new(big.Int)
	wei.SetString(balance, 10)

	// Convert wei to ether (1 ether = 10^18 wei)
	ether := new(big.Float).SetInt(wei)
	ether.Quo(ether, big.NewFloat(1e18))

	return ether.String(), nil
}

// GetBlockNumber 获取最新区块高度
func (c *EthereumClient) GetBlockNumber(ctx context.Context) (uint64, error) {
	blockNumber, err := c.client.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get block number: %w", err)
	}
	return blockNumber, nil
}

// GetTransaction 获取交易信息
func (c *EthereumClient) GetTransaction(ctx context.Context, txHash string) (*Transaction, error) {
	if err := ValidateTxHash(c.chain, txHash); err != nil {
		return nil, err
	}

	hash := common.HexToHash(txHash)
	tx, isPending, err := c.client.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	receipt, err := c.client.TransactionReceipt(ctx, hash)
	if err != nil && !isPending {
		return nil, fmt.Errorf("failed to get transaction receipt: %w", err)
	}

	result := &Transaction{
		Hash:  txHash,
		Value: tx.Value().String(),
		Nonce: tx.Nonce(),
		Data:  common.Bytes2Hex(tx.Data()),
		Extra: make(map[string]interface{}),
	}

	if to := tx.To(); to != nil {
		result.To = to.Hex()
	}

	if tx.GasPrice() != nil {
		result.GasPrice = tx.GasPrice().String()
	}

	if receipt != nil {
		result.BlockNumber = receipt.BlockNumber.Uint64()
		result.BlockHash = receipt.BlockHash.Hex()
		result.GasUsed = receipt.GasUsed

		// Get From address from transaction sender
		msg, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err == nil {
			result.From = msg.Hex()
		}

		if receipt.Status == types.ReceiptStatusSuccessful {
			result.Status = "success"
		} else {
			result.Status = "failed"
		} // Get block timestamp
		block, err := c.client.BlockByHash(ctx, receipt.BlockHash)
		if err == nil {
			result.Timestamp = int64(block.Time())
		}
	} else {
		result.Status = "pending"
	}

	return result, nil
}

// SendTransaction 发送交易
func (c *EthereumClient) SendTransaction(ctx context.Context, tx *TransactionRequest) (string, error) {
	// Note: This is a placeholder. Actual implementation requires:
	// 1. Private key for signing
	// 2. Proper transaction construction
	// 3. Gas estimation
	return "", fmt.Errorf("sendTransaction not implemented: requires private key integration")
}

// GetChain 获取链类型
func (c *EthereumClient) GetChain() Chain {
	return c.chain
}

// Close 关闭连接
func (c *EthereumClient) Close() error {
	c.client.Close()
	c.rpc.Close()
	return nil
}

// GetChainID 获取链 ID
func (c *EthereumClient) GetChainID(ctx context.Context) (*big.Int, error) {
	return c.client.ChainID(ctx)
}

// GetGasPrice 获取当前 Gas 价格
func (c *EthereumClient) GetGasPrice(ctx context.Context) (string, error) {
	gasPrice, err := c.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %w", err)
	}
	return gasPrice.String(), nil
}

// EstimateGas 估算 Gas 用量
func (c *EthereumClient) EstimateGas(ctx context.Context, from, to, data string, value *big.Int) (uint64, error) {
	msg := map[string]interface{}{
		"from": from,
		"to":   to,
	}

	if data != "" {
		msg["data"] = data
	}

	if value != nil && value.Sign() > 0 {
		msg["value"] = fmt.Sprintf("0x%x", value)
	}

	var result string
	err := c.rpc.CallContext(ctx, &result, "eth_estimateGas", msg)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas: %w", err)
	}

	gas := new(big.Int)
	gas.SetString(result[2:], 16)
	return gas.Uint64(), nil
}

// GetTransactionCount 获取地址的交易计数（nonce）
func (c *EthereumClient) GetTransactionCount(ctx context.Context, address string) (uint64, error) {
	if err := ValidateAddress(c.chain, address); err != nil {
		return 0, err
	}

	addr := common.HexToAddress(address)
	nonce, err := c.client.PendingNonceAt(ctx, addr)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction count: %w", err)
	}
	return nonce, nil
}

// GetCode 获取合约代码
func (c *EthereumClient) GetCode(ctx context.Context, contractAddress string) ([]byte, error) {
	if err := ValidateAddress(c.chain, contractAddress); err != nil {
		return nil, err
	}

	addr := common.HexToAddress(contractAddress)
	code, err := c.client.CodeAt(ctx, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract code: %w", err)
	}
	return code, nil
}

// IsContract 检查地址是否为合约
func (c *EthereumClient) IsContract(ctx context.Context, address string) (bool, error) {
	code, err := c.GetCode(ctx, address)
	if err != nil {
		return false, err
	}
	return len(code) > 0, nil
}
