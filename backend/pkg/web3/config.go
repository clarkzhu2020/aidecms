package web3

import (
	"os"
	"sync"
)

// Config Web3 配置
type Config struct {
	// Ethereum
	EthereumRPC string

	// BSC (Binance Smart Chain)
	BSCRPC string

	// Bitcoin
	BitcoinRPC    string
	BitcoinAPIKey string

	// Solana
	SolanaRPC string

	// Exchanges
	CoinbaseAPIKey    string
	CoinbaseAPISecret string
	KuCoinAPIKey      string
	KuCoinAPISecret   string
	KuCoinPassphrase  string

	// Hyperliquid DEX
	HyperliquidPrivateKey string
	HyperliquidAddress    string
}

var (
	config     *Config
	configOnce sync.Once
)

// LoadConfig 加载配置
func LoadConfig() *Config {
	configOnce.Do(func() {
		config = &Config{
			EthereumRPC:       getEnv("WEB3_ETHEREUM_RPC", "https://mainnet.infura.io/v3/YOUR-PROJECT-ID"),
			BSCRPC:            getEnv("WEB3_BSC_RPC", "https://bsc-dataseed.binance.org/"),
			BitcoinRPC:        getEnv("WEB3_BITCOIN_RPC", "https://bitcoin-mainnet.core.chainstack.com"),
			BitcoinAPIKey:     getEnv("WEB3_BITCOIN_API_KEY", ""),
			SolanaRPC:         getEnv("WEB3_SOLANA_RPC", "https://api.mainnet-beta.solana.com"),
			CoinbaseAPIKey:    getEnv("EXCHANGE_COINBASE_API_KEY", ""),
			CoinbaseAPISecret: getEnv("EXCHANGE_COINBASE_API_SECRET", ""),
			KuCoinAPIKey:      getEnv("EXCHANGE_KUCOIN_API_KEY", ""),
			KuCoinAPISecret:   getEnv("EXCHANGE_KUCOIN_API_SECRET", ""),
			KuCoinPassphrase:  getEnv("EXCHANGE_KUCOIN_PASSPHRASE", ""),

			// Hyperliquid DEX
			HyperliquidPrivateKey: getEnv("EXCHANGE_HYPERLIQUID_PRIVATE_KEY", ""),
			HyperliquidAddress:    getEnv("EXCHANGE_HYPERLIQUID_ADDRESS", ""),
		}
	})
	return config
}

// GetConfig 获取配置
func GetConfig() *Config {
	if config == nil {
		return LoadConfig()
	}
	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// InitializeClients 初始化所有客户端
func InitializeClients() error {
	cfg := GetConfig()
	manager := GetManager()

	// Initialize Ethereum
	if cfg.EthereumRPC != "" {
		ethClient, err := NewEthereumClient(cfg.EthereumRPC)
		if err != nil {
			return err
		}
		manager.RegisterClient(Ethereum, ethClient)
	}

	// Initialize BSC
	if cfg.BSCRPC != "" {
		bscClient, err := NewBSCClient(cfg.BSCRPC)
		if err != nil {
			return err
		}
		manager.RegisterClient(BSC, bscClient)
	}

	// Initialize Bitcoin
	if cfg.BitcoinRPC != "" {
		btcClient := NewBitcoinClient(cfg.BitcoinRPC, cfg.BitcoinAPIKey)
		manager.RegisterClient(Bitcoin, btcClient)
	}

	// Initialize Solana
	if cfg.SolanaRPC != "" {
		solClient := NewSolanaClient(cfg.SolanaRPC)
		manager.RegisterClient(Solana, solClient)
	}

	// Initialize Exchanges
	exchangeManager := GetExchangeManager()

	if cfg.CoinbaseAPIKey != "" && cfg.CoinbaseAPISecret != "" {
		coinbaseClient := NewCoinbaseClient(cfg.CoinbaseAPIKey, cfg.CoinbaseAPISecret)
		exchangeManager.RegisterExchange(Coinbase, coinbaseClient)
	}

	if cfg.KuCoinAPIKey != "" && cfg.KuCoinAPISecret != "" && cfg.KuCoinPassphrase != "" {
		kucoinClient := NewKuCoinClient(cfg.KuCoinAPIKey, cfg.KuCoinAPISecret, cfg.KuCoinPassphrase)
		exchangeManager.RegisterExchange(KuCoin, kucoinClient)
	}

	// Initialize Hyperliquid DEX
	if cfg.HyperliquidPrivateKey != "" {
		hyperliquidClient, err := NewHyperliquidClient(cfg.HyperliquidPrivateKey)
		if err == nil {
			exchangeManager.RegisterExchange(Hyperliquid, hyperliquidClient)
		}
	}

	return nil
}
