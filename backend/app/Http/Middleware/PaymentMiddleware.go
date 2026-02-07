package middleware

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/plutov/paypal/v4"
)

// PayPalWebhookSignatureMiddleware PayPal Webhook签名验证中间件
func PayPalWebhookSignatureMiddleware(webhookID string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 这里可以添加Webhook签名验证逻辑
		// 注意: PayPal的Webhook验证需要特定的处理

		// 提取Webhook ID（可以从环境变量或配置中获取）
		if webhookID == "" {
			c.JSON(400, map[string]interface{}{
				"error": "Webhook ID not configured",
			})
			c.Abort()
			return
		}

		// 将webhookID存入context供后续使用
		c.Set("webhook_id", webhookID)

		c.Next(ctx)
	}
}

// PayPalPaymentRequiredMiddleware PayPal支付必需中间件
// 用于确保某些操作只有在PayPal支付完成后才能执行
func PayPalPaymentRequiredMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从context或JWT中获取支付ID
		paymentID := c.Param("payment_id")
		if paymentID == "" {
			c.JSON(400, map[string]interface{}{
				"error": "Payment ID is required",
			})
			c.Abort()
			return
		}

		// TODO: 验证支付是否存在且状态为paid
		// 这里需要注入Repository来验证

		// 将paymentID存入context
		c.Set("payment_id", paymentID)

		c.Next(ctx)
	}
}

// PayPalConfiguredMiddleware PayPal配置检查中间件
// 用于确保PayPal已正确配置
func PayPalConfiguredMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 检查PayPal是否已初始化
		// 这里需要从全局变量或依赖注入中检查

		// 将PayPal客户端存入context
		var paypalClient *paypal.Client
		// paypalClient = paypal.GetClient()

		if paypalClient == nil {
			c.JSON(503, map[string]interface{}{
				"error": "PayPal service is not configured",
			})
			c.Abort()
			return
		}

		c.Set("paypal_client", paypalClient)

		c.Next(ctx)
	}
}

// PayPalIPWhitelistMiddleware PayPal IP白名单中间件
// 用于验证请求是否来自PayPal服务器
func PayPalIPWhitelistMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// PayPal官方IP地址段
		// 根据环境（沙盒/生产）使用不同的IP段

		// 获取客户端IP
		clientIP := c.ClientIP()

		// TODO: 验证IP是否在PayPal白名单中
		_ = clientIP

		c.Next(ctx)
	}
}

// ValidatePayPalWebhookEvent 验证PayPal Webhook事件类型
func ValidatePayPalWebhookEvent(validEvents []string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 解析请求体获取事件类型
		var webhookData map[string]interface{}
		if err := c.BindJSON(&webhookData); err != nil {
			c.JSON(400, map[string]interface{}{
				"error": "Invalid webhook data",
			})
			c.Abort()
			return
		}

		eventType, ok := webhookData["event_type"].(string)
		if !ok {
			c.JSON(400, map[string]interface{}{
				"error": "Event type not found in webhook data",
			})
			c.Abort()
			return
		}

		// 验证事件类型是否在允许的列表中
		isValid := false
		for _, validEvent := range validEvents {
			if eventType == validEvent {
				isValid = true
				break
			}
		}

		if !isValid {
			c.JSON(400, map[string]interface{}{
				"error": fmt.Sprintf("Invalid event type: %s", eventType),
			})
			c.Abort()
			return
		}

		// 将事件类型存入context
		c.Set("event_type", eventType)

		c.Next(ctx)
	}
}

// LogPayPalWebhookMiddleware PayPal Webhook日志记录中间件
func LogPayPalWebhookMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 记录Webhook请求的详细信息
		// 包括: 事件ID、事件类型、时间戳、IP地址等

		c.Next(ctx)
	}
}
