package config

import (
	"log"
	"os"

	"github.com/clarkzhu2020/aidecms/pkg/moonpay"
)

// InitMoonPay 初始化MoonPay配置
func InitMoonPay() {
	// 获取环境变量
	apiKey := os.Getenv("MOONPAY_API_KEY")
	secretKey := os.Getenv("MOONPAY_SECRET_KEY")
	webhookKey := os.Getenv("MOONPAY_WEBHOOK_KEY")
	isSandbox := os.Getenv("MOONPAY_SANDBOX") == "true"

	// 检查必需的环境变量
	if apiKey == "" {
		log.Println("Warning: MOONPAY_API_KEY not set, MoonPay integration disabled")
		return
	}

	// 初始化MoonPay客户端
	if err := moonpay.Init(apiKey, secretKey, webhookKey, isSandbox); err != nil {
		log.Printf("Failed to initialize MoonPay client: %v", err)
		return
	}

	log.Println("MoonPay integration initialized successfully")
}
