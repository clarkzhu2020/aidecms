package config

import (
	"aidecms/pkg/kucoin"
	"os"
	"strconv"
)

// InitKuCoin 初始化KuCoin配置
func InitKuCoin() error {
	enabled, _ := strconv.ParseBool(os.Getenv("KUCOIN_ENABLED"))
	if !enabled {
		return nil
	}

	apiKey := os.Getenv("KUCOIN_API_KEY")
	apiSecret := os.Getenv("KUCOIN_API_SECRET")
	passphrase := os.Getenv("KUCOIN_PASSPHRASE")
	isSandbox, _ := strconv.ParseBool(os.Getenv("KUCOIN_SANDBOX"))

	if apiKey == "" || apiSecret == "" || passphrase == "" {
		return nil
	}

	return kucoin.Init(apiKey, apiSecret, passphrase, isSandbox)
}
