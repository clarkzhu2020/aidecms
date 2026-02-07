package config

import (
	"log"

	"github.com/clarkzhu2020/aidecms/pkg/stripe"
)

// InitStripe 初始化Stripe支付服务
func InitStripe() error {
	env := GetEnv()

	// 检查是否启用Stripe
	enabled := env.GetBool("STRIPE_ENABLED", false)
	if !enabled {
		log.Println("Stripe is disabled")
		return nil
	}

	// 获取Stripe配置
	secretKey := env.Get("STRIPE_SECRET_KEY", "")
	webhookSecret := env.Get("STRIPE_WEBHOOK_SECRET", "")
	mode := env.Get("STRIPE_MODE", "test")

	if secretKey == "" {
		log.Println("Stripe secret key not configured, skipping initialization")
		return nil
	}

	// 判断是否为测试模式
	isTest := mode == "test"

	// 初始化Stripe客户端
	if err := stripe.Init(secretKey, webhookSecret, isTest); err != nil {
		log.Printf("Failed to initialize Stripe: %v", err)
		return err
	}

	log.Printf("Stripe initialized successfully in %s mode", mode)
	return nil
}
