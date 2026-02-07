package paypal

import (
	"context"
	"fmt"
	"log"

	"github.com/plutov/paypal/v4"
)

// OrderService 支付订单服务
type OrderService struct {
	client *Client
}

// NewOrderService 创建支付订单服务
func NewOrderService() *OrderService {
	return &OrderService{
		client: GetClient(),
	}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	OrderID      string  `json:"order_id"`       // 内部订单ID
	Amount       float64 `json:"amount"`         // 支付金额
	Currency     string  `json:"currency"`       // 货币类型 (USD, EUR, CNY等)
	Description  string  `json:"description"`    // 订单描述
	ReturnURL    string  `json:"return_url"`     // 支付成功回调URL
	CancelURL    string  `json:"cancel_url"`     // 取消支付URL
	ReferenceID  string  `json:"reference_id"`   // 参考ID
	ItemName     string  `json:"item_name"`      // 商品名称
	ItemQuantity int     `json:"item_quantity"`  // 商品数量
}

// CreateOrder 创建PayPal订单
func (s *OrderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*paypal.Order, error) {
	if s.client == nil {
		return nil, fmt.Errorf("PayPal client not initialized")
	}

	// 创建购买单元
	units := []paypal.PurchaseUnitRequest{
		{
			ReferenceID: req.ReferenceID,
			Amount: &paypal.PurchaseUnitAmount{
				CurrencyCode: req.Currency,
				Value:        fmt.Sprintf("%.2f", req.Amount),
			},
			Description: req.Description,
		},
	}

	// 设置应用上下文
	appCtx := &paypal.ApplicationContext{
		ReturnURL: req.ReturnURL,
		CancelURL: req.CancelURL,
		UserAction: "PAY_NOW",
		BrandName:  "AideCMS",
	}

	// 创建订单
	order, err := s.client.CreateOrder(ctx, paypal.OrderIntentCapture, units, nil, appCtx)
	if err != nil {
		log.Printf("Failed to create PayPal order: %v", err)
		return nil, fmt.Errorf("failed to create PayPal order: %w", err)
	}

	log.Printf("PayPal order created successfully: %s", order.ID)
	return order, nil
}

// GetOrder 获取订单详情
func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*paypal.Order, error) {
	if s.client == nil {
		return nil, fmt.Errorf("PayPal client not initialized")
	}

	order, err := s.client.GetOrder(ctx, orderID)
	if err != nil {
		log.Printf("Failed to get PayPal order %s: %v", orderID, err)
		return nil, fmt.Errorf("failed to get PayPal order: %w", err)
	}

	return order, nil
}

// CaptureOrder 捕获支付（确认支付）
func (s *OrderService) CaptureOrder(ctx context.Context, orderID string) (*paypal.CaptureOrderResponse, error) {
	if s.client == nil {
		return nil, fmt.Errorf("PayPal client not initialized")
	}

	capture, err := s.client.CaptureOrder(ctx, orderID, paypal.CaptureOrderRequest{})
	if err != nil {
		log.Printf("Failed to capture PayPal order %s: %v", orderID, err)
		return nil, fmt.Errorf("failed to capture PayPal order: %w", err)
	}

	log.Printf("PayPal order captured successfully: %s", orderID)
	return capture, nil
}

// AuthorizeOrder 授权支付
func (s *OrderService) AuthorizeOrder(ctx context.Context, orderID string) (*paypal.AuthorizeOrderResponse, error) {
	if s.client == nil {
		return nil, fmt.Errorf("PayPal client not initialized")
	}

	auth, err := s.client.AuthorizeOrder(ctx, orderID, paypal.AuthorizeOrderRequest{})
	if err != nil {
		log.Printf("Failed to authorize PayPal order %s: %v", orderID, err)
		return nil, fmt.Errorf("failed to authorize PayPal order: %w", err)
	}

	log.Printf("PayPal order authorized successfully: %s", orderID)
	return auth, nil
}

// VoidOrder 取消已授权的订单
func (s *OrderService) VoidOrder(ctx context.Context, orderID string) error {
	if s.client == nil {
		return fmt.Errorf("PayPal client not initialized")
	}

	err := s.client.VoidOrder(ctx, orderID)
	if err != nil {
		log.Printf("Failed to void PayPal order %s: %v", orderID, err)
		return fmt.Errorf("failed to void PayPal order: %w", err)
	}

	log.Printf("PayPal order voided successfully: %s", orderID)
	return nil
}

// RefundPayment 退款
func (s *OrderService) RefundPayment(ctx context.Context, captureID string, amount *float64, currency string) (*paypal.RefundResponse, error) {
	if s.client == nil {
		return nil, fmt.Errorf("PayPal client not initialized")
	}

	refundReq := &paypal.RefundRequest{}

	// 如果指定了金额，则部分退款
	if amount != nil {
		refundReq.Amount = &paypal.PurchaseUnitAmount{
			CurrencyCode: currency,
			Value:        fmt.Sprintf("%.2f", *amount),
		}
	}

	refund, err := s.client.RefundCapture(ctx, captureID, refundReq)
	if err != nil {
		log.Printf("Failed to refund capture %s: %v", captureID, err)
		return nil, fmt.Errorf("failed to refund payment: %w", err)
	}

	log.Printf("Refund processed successfully for capture: %s", captureID)
	return refund, nil
}

// VerifyWebhook 验证Webhook通知
func (s *OrderService) VerifyWebhook(ctx context.Context, webhookID string, headers map[string]string, body []byte) (bool, error) {
	if s.client == nil {
		return false, fmt.Errorf("PayPal client not initialized")
	}

	event, err := s.client.VerifyWebhookSignature(ctx, paypal.VerifyWebhookSignatureRequest{
		WebhookID:         webhookID,
		AuthVersion:       headers["PayPal-Auth-Algo"],
		CertID:            headers["PayPal-Cert-Id"],
		TransmissionID:    headers["PayPal-Transmission-Id"],
		TransmissionSig:   headers["PayPal-Transmission-Sig"],
		TransmissionTime:  headers["PayPal-Transmission-Time"],
		WebhookEvent:      body,
	})

	if err != nil {
		log.Printf("Failed to verify webhook: %v", err)
		return false, fmt.Errorf("failed to verify webhook: %w", err)
	}

	return event.VerificationStatus == "SUCCESS", nil
}
