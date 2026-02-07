package config

import (
	"log"
	"os"

	"github.com/clarkzhu2020/aidecms/pkg/coinbase"
)

// InitCoinbase 初始化Coinbase配置
func InitCoinbase() {
	// 获取环境变量
	apiKey := os.Getenv("COINBASE_API_KEY")
	apiSecret := os.Getenv("COINBASE_API_SECRET")
	webhookKey := os.Getenv("COINBASE_WEBHOOK_KEY")
	accountID := os.Getenv("COINBASE_ACCOUNT_ID")
	isSandbox := os.Getenv("COINBASE_SANDBOX") == "true"

	// 检查必需的环境变量
	if apiKey == "" {
		log.Println("Warning: COINBASE_API_KEY not set, Coinbase integration disabled")
		return
	}

	// 初始化Coinbase客户端
	if err := coinbase.Init(apiKey, apiSecret, webhookKey, accountID, isSandbox); err != nil {
		log.Printf("Failed to initialize Coinbase client: %v", err)
		return
	}

	log.Println("Coinbase integration initialized successfully")
}
