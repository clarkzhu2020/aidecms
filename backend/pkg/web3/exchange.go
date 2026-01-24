package web3

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Exchange 交易所类型
type Exchange string

const (
	Coinbase    Exchange = "coinbase"
	KuCoin      Exchange = "kucoin"
	Hyperliquid Exchange = "hyperliquid"
)

// ExchangeClient 交易所客户端接口
type ExchangeClient interface {
	// GetBalance 获取余额
	GetBalance(ctx context.Context, currency string) (string, error)

	// GetBalances 获取所有余额
	GetBalances(ctx context.Context) (map[string]string, error)

	// GetPrice 获取价格
	GetPrice(ctx context.Context, pair string) (string, error)
}

// ExchangeManager 交易所管理器
type ExchangeManager struct {
	exchanges map[Exchange]ExchangeClient
	mu        sync.RWMutex
}

var (
	globalExchangeManager *ExchangeManager
	exchangeOnce          sync.Once
)

// GetExchangeManager 获取全局交易所管理器
func GetExchangeManager() *ExchangeManager {
	exchangeOnce.Do(func() {
		globalExchangeManager = &ExchangeManager{
			exchanges: make(map[Exchange]ExchangeClient),
		}
	})
	return globalExchangeManager
}

// RegisterExchange 注册交易所客户端
func (m *ExchangeManager) RegisterExchange(exchange Exchange, client ExchangeClient) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.exchanges[exchange] = client
}

// GetExchange 获取交易所客户端
func (m *ExchangeManager) GetExchange(exchange Exchange) (ExchangeClient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, exists := m.exchanges[exchange]
	if !exists {
		return nil, fmt.Errorf("exchange %s not registered", exchange)
	}
	return client, nil
}

// GetBalance 获取余额
func (m *ExchangeManager) GetBalance(ctx context.Context, exchange Exchange, currency string) (string, error) {
	client, err := m.GetExchange(exchange)
	if err != nil {
		return "", err
	}
	return client.GetBalance(ctx, currency)
}

// GetBalances 获取所有余额
func (m *ExchangeManager) GetBalances(ctx context.Context, exchange Exchange) (map[string]string, error) {
	client, err := m.GetExchange(exchange)
	if err != nil {
		return nil, err
	}
	return client.GetBalances(ctx)
}

// GetPrice 获取价格
func (m *ExchangeManager) GetPrice(ctx context.Context, exchange Exchange, pair string) (string, error) {
	client, err := m.GetExchange(exchange)
	if err != nil {
		return "", err
	}
	return client.GetPrice(ctx, pair)
}

// GetSupportedExchanges 获取支持的交易所
func (m *ExchangeManager) GetSupportedExchanges() []Exchange {
	m.mu.RLock()
	defer m.mu.RUnlock()

	exchanges := make([]Exchange, 0, len(m.exchanges))
	for exchange := range m.exchanges {
		exchanges = append(exchanges, exchange)
	}
	return exchanges
}

// Close 关闭所有交易所客户端连接
func (m *ExchangeManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 交易所客户端使用 HTTP，无需显式关闭
	m.exchanges = make(map[Exchange]ExchangeClient)
	return nil
}

// MultiExchangeBalance 多交易所余额
type MultiExchangeBalance struct {
	Exchange Exchange          `json:"exchange"`
	Balances map[string]string `json:"balances"`
	Error    string            `json:"error,omitempty"`
}

// GetAllExchangeBalances 获取所有交易所的余额
func GetAllExchangeBalances(ctx context.Context, currency string) ([]MultiExchangeBalance, error) {
	manager := GetExchangeManager()
	exchanges := manager.GetSupportedExchanges()

	if len(exchanges) == 0 {
		return nil, errors.New("no exchanges configured")
	}

	results := make([]MultiExchangeBalance, 0, len(exchanges))

	for _, exchange := range exchanges {
		balance, err := manager.GetBalance(ctx, exchange, currency)

		result := MultiExchangeBalance{
			Exchange: exchange,
			Balances: make(map[string]string),
		}

		if err != nil {
			result.Error = err.Error()
		} else {
			result.Balances[currency] = balance
		}

		results = append(results, result)
	}

	return results, nil
}

// ExchangePrice 交易所价格
type ExchangePrice struct {
	Exchange Exchange `json:"exchange"`
	Pair     string   `json:"pair"`
	Price    string   `json:"price"`
	Error    string   `json:"error,omitempty"`
}

// GetAllExchangePrices 获取所有交易所的价格
func GetAllExchangePrices(ctx context.Context, pair string) ([]ExchangePrice, error) {
	manager := GetExchangeManager()
	exchanges := manager.GetSupportedExchanges()

	if len(exchanges) == 0 {
		return nil, errors.New("no exchanges configured")
	}

	results := make([]ExchangePrice, 0, len(exchanges))

	for _, exchange := range exchanges {
		price, err := manager.GetPrice(ctx, exchange, pair)

		result := ExchangePrice{
			Exchange: exchange,
			Pair:     pair,
		}

		if err != nil {
			result.Error = err.Error()
		} else {
			result.Price = price
		}

		results = append(results, result)
	}

	return results, nil
}

// ExchangeConfig 交易所配置
type ExchangeConfig struct {
	Coinbase struct {
		APIKey    string
		APISecret string
	}
	KuCoin struct {
		APIKey     string
		APISecret  string
		Passphrase string
	}
}

// InitializeExchanges 初始化所有交易所客户端
func InitializeExchanges(config *ExchangeConfig) error {
	manager := GetExchangeManager()

	// 初始化 Coinbase
	if config.Coinbase.APIKey != "" && config.Coinbase.APISecret != "" {
		coinbaseClient := NewCoinbaseClient(config.Coinbase.APIKey, config.Coinbase.APISecret)
		manager.RegisterExchange(Coinbase, coinbaseClient)
	}

	// 初始化 KuCoin
	if config.KuCoin.APIKey != "" && config.KuCoin.APISecret != "" && config.KuCoin.Passphrase != "" {
		kucoinClient := NewKuCoinClient(config.KuCoin.APIKey, config.KuCoin.APISecret, config.KuCoin.Passphrase)
		manager.RegisterExchange(KuCoin, kucoinClient)
	}

	return nil
}
