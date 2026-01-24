package web3

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Chain 区块链类型
type Chain string

const (
	Bitcoin  Chain = "bitcoin"
	Ethereum Chain = "ethereum"
	BSC      Chain = "bsc" // Binance Smart Chain
	Solana   Chain = "solana"
)

// Client Web3 客户端接口
type Client interface {
	// GetBalance 获取地址余额
	GetBalance(ctx context.Context, address string) (string, error)

	// GetBlockNumber 获取最新区块高度
	GetBlockNumber(ctx context.Context) (uint64, error)

	// GetTransaction 获取交易信息
	GetTransaction(ctx context.Context, txHash string) (*Transaction, error)

	// SendTransaction 发送交易
	SendTransaction(ctx context.Context, tx *TransactionRequest) (string, error)

	// GetChain 获取链类型
	GetChain() Chain

	// Close 关闭连接
	Close() error
}

// Transaction 交易信息
type Transaction struct {
	Hash        string                 `json:"hash"`
	From        string                 `json:"from"`
	To          string                 `json:"to"`
	Value       string                 `json:"value"`
	BlockNumber uint64                 `json:"block_number"`
	BlockHash   string                 `json:"block_hash"`
	Status      string                 `json:"status"`
	GasUsed     uint64                 `json:"gas_used,omitempty"`
	GasPrice    string                 `json:"gas_price,omitempty"`
	Nonce       uint64                 `json:"nonce,omitempty"`
	Data        string                 `json:"data,omitempty"`
	Timestamp   int64                  `json:"timestamp"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
}

// TransactionRequest 交易请求
type TransactionRequest struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data,omitempty"`
	GasLimit uint64 `json:"gas_limit,omitempty"`
	GasPrice string `json:"gas_price,omitempty"`
	Nonce    uint64 `json:"nonce,omitempty"`
}

// Manager Web3 管理器
type Manager struct {
	clients map[Chain]Client
	mu      sync.RWMutex
}

var (
	globalManager *Manager
	once          sync.Once
)

// GetManager 获取全局管理器
func GetManager() *Manager {
	once.Do(func() {
		globalManager = &Manager{
			clients: make(map[Chain]Client),
		}
	})
	return globalManager
}

// RegisterClient 注册客户端
func (m *Manager) RegisterClient(chain Chain, client Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[chain] = client
}

// GetClient 获取客户端
func (m *Manager) GetClient(chain Chain) (Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, exists := m.clients[chain]
	if !exists {
		return nil, fmt.Errorf("client for chain %s not registered", chain)
	}
	return client, nil
}

// GetBalance 获取余额
func (m *Manager) GetBalance(ctx context.Context, chain Chain, address string) (string, error) {
	client, err := m.GetClient(chain)
	if err != nil {
		return "", err
	}
	return client.GetBalance(ctx, address)
}

// GetTransaction 获取交易
func (m *Manager) GetTransaction(ctx context.Context, chain Chain, txHash string) (*Transaction, error) {
	client, err := m.GetClient(chain)
	if err != nil {
		return nil, err
	}
	return client.GetTransaction(ctx, txHash)
}

// SendTransaction 发送交易
func (m *Manager) SendTransaction(ctx context.Context, chain Chain, tx *TransactionRequest) (string, error) {
	client, err := m.GetClient(chain)
	if err != nil {
		return "", err
	}
	return client.SendTransaction(ctx, tx)
}

// Close 关闭所有客户端
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for chain, client := range m.clients {
		if err := client.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close %s client: %w", chain, err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// GetSupportedChains 获取支持的链
func (m *Manager) GetSupportedChains() []Chain {
	m.mu.RLock()
	defer m.mu.RUnlock()

	chains := make([]Chain, 0, len(m.clients))
	for chain := range m.clients {
		chains = append(chains, chain)
	}
	return chains
}

// ValidateAddress 验证地址格式
func ValidateAddress(chain Chain, address string) error {
	if address == "" {
		return errors.New("address is required")
	}

	switch chain {
	case Bitcoin:
		// Bitcoin 地址验证（简化版）
		// Legacy (1...), SegWit (3...), Bech32 (bc1...)
		if len(address) < 26 || len(address) > 90 {
			return errors.New("invalid bitcoin address length")
		}
	case Ethereum, BSC:
		// Ethereum/BSC 地址验证
		if len(address) != 42 || address[:2] != "0x" {
			return errors.New("invalid ethereum address format")
		}
	case Solana:
		// Solana 地址验证（Base58）
		if len(address) < 32 || len(address) > 44 {
			return errors.New("invalid solana address length")
		}
	default:
		return fmt.Errorf("unsupported chain: %s", chain)
	}

	return nil
}

// ValidateTxHash 验证交易哈希格式
func ValidateTxHash(chain Chain, txHash string) error {
	if txHash == "" {
		return errors.New("transaction hash is required")
	}

	switch chain {
	case Bitcoin:
		if len(txHash) != 64 {
			return errors.New("invalid bitcoin transaction hash length")
		}
	case Ethereum, BSC:
		if len(txHash) != 66 || txHash[:2] != "0x" {
			return errors.New("invalid ethereum transaction hash format")
		}
	case Solana:
		if len(txHash) < 80 || len(txHash) > 90 {
			return errors.New("invalid solana transaction hash length")
		}
	default:
		return fmt.Errorf("unsupported chain: %s", chain)
	}

	return nil
}
