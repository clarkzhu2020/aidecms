package config

import (
	"log"

	"github.com/clarkzhu2020/aidecms/pkg/paypal"
)

// InitPayPal 初始化PayPal支付服务
func InitPayPal() error {
	env := GetEnv()
	
	// 检查是否启用PayPal
	enabled := env.GetBool("PAYPAL_ENABLED", false)
	if !enabled {
		log.Println("PayPal is disabled")
		return nil
	}

	// 获取PayPal配置
	clientID := env.Get("PAYPAL_CLIENT_ID", "")
	clientSecret := env.Get("PAYPAL_CLIENT_SECRET", "")
	mode := env.Get("PAYPAL_MODE", "sandbox")
	
	if clientID == "" || clientSecret == "" {
		log.Println("PayPal credentials not configured, skipping initialization")
		return nil
	}

	// 判断是否为沙盒环境
	isSandbox := mode == "sandbox"

	// 初始化PayPal客户端
	if err := paypal.Init(clientID, clientSecret, isSandbox); err != nil {
		log.Printf("Failed to initialize PayPal: %v", err)
		return err
	}

	log.Printf("PayPal initialized successfully in %s mode", mode)
	return nil
}
