package stripe

import (
	"context"
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

// PaymentService Stripe支付服务
type PaymentService struct {
	client *Client
}

// NewPaymentService 创建支付服务
func NewPaymentService() *PaymentService {
	return &PaymentService{
		client: GetClient(),
	}
}

// CreatePaymentIntentRequest 创建支付意图请求
type CreatePaymentIntentRequest struct {
	Amount           int64   `json:"amount"`            // 金额(最小单位,如 cents)
	Currency         string  `json:"currency"`          // 货币代码(USD, EUR, CNY等)
	CustomerID       string  `json:"customer_id"`       // 客户ID(可选)
	Description      string  `json:"description"`       // 描述
	ReceiptEmail     string  `json:"receipt_email"`     // 收据邮箱
	Metadata         map[string]string `json:"metadata"` // 元数据
	PaymentMethodTypes []string `json:"payment_method_types"` // 支付方式
}

// CreatePaymentIntent 创建支付意图
func (s *PaymentService) CreatePaymentIntent(ctx context.Context, req *CreatePaymentIntentRequest) (*stripe.PaymentIntent, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Stripe client not initialized")
	}

	// 设置默认支付方式
	if req.PaymentMethodTypes == nil || len(req.PaymentMethodTypes) == 0 {
		req.PaymentMethodTypes = []string{"card"}
	}

	params := &stripe.PaymentIntentParams{
		Amount:           stripe.Int64(req.Amount),
		Currency:         stripe.String(req.Currency),
		Description:      stripe.String(req.Description),
	}

	// 可选参数
	if req.CustomerID != "" {
		params.Customer = stripe.String(req.CustomerID)
	}
	if req.ReceiptEmail != "" {
		params.ReceiptEmail = stripe.String(req.ReceiptEmail)
	}
	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}
	if req.PaymentMethodTypes != nil {
		params.PaymentMethodTypes = stripe.StringSlice(req.PaymentMethodTypes)
	}

	// 创建支付意图
	pi, err := stripe.NewClient(s.client.GetSecretKey()).PaymentIntents.New(ctx, params)
	if err != nil {
		log.Printf("Failed to create payment intent: %v", err)
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	log.Printf("PaymentIntent created successfully: %s", pi.ID)
	return pi, nil
}

// GetPaymentIntent 获取支付意图
func (s *PaymentService) GetPaymentIntent(ctx context.Context, paymentIntentID string) (*stripe.PaymentIntent, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Stripe client not initialized")
	}

	pi, err := stripe.NewClient(s.client.GetSecretKey()).PaymentIntents.Get(ctx, paymentIntentID, nil)
	if err != nil {
		log.Printf("Failed to get payment intent %s: %v", paymentIntentID, err)
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}

	return pi, nil
}

// ConfirmPaymentIntent 确认支付意图
func (s *PaymentService) ConfirmPaymentIntent(ctx context.Context, paymentIntentID string, paymentMethodID string) (*stripe.PaymentIntent, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Stripe client not initialized")
	}

	params := &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String(paymentMethodID),
	}

	pi, err := stripe.NewClient(s.client.GetSecretKey()).PaymentIntents.Confirm(ctx, paymentIntentID, params)
	if err != nil {
		log.Printf("Failed to confirm payment intent %s: %v", paymentIntentID, err)
		return nil, fmt.Errorf("failed to confirm payment intent: %w", err)
	}

	log.Printf("PaymentIntent confirmed successfully: %s", paymentIntentID)
	return pi, nil
}

// CancelPaymentIntent 取消支付意图
func (s *PaymentService) CancelPaymentIntent(ctx context.Context, paymentIntentID string) (*stripe.PaymentIntent, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Stripe client not initialized")
	}

	pi, err := stripe.NewClient(s.client.GetSecretKey()).PaymentIntents.Cancel(ctx, paymentIntentID, nil)
	if err != nil {
		log.Printf("Failed to cancel payment intent %s: %v", paymentIntentID, err)
		return nil, fmt.Errorf("failed to cancel payment intent: %w", err)
	}

	log.Printf("PaymentIntent cancelled successfully: %s", paymentIntentID)
	return pi, nil
}

// ListPaymentIntents 列出支付意图
func (s *PaymentService) ListPaymentIntents(ctx context.Context, limit int64, startingAfter string) (*stripe.PaymentIntentList, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Stripe client not initialized")
	}

	params := &stripe.PaymentIntentListParams{
		Limit: stripe.Int64(limit),
	}

	if startingAfter != "" {
		params.StartingAfter = stripe.String(startingAfter)
	}

	piList, err := stripe.NewClient(s.client.GetSecretKey()).PaymentIntents.List(ctx, params)
	if err != nil {
		log.Printf("Failed to list payment intents: %v", err)
		return nil, fmt.Errorf("failed to list payment intents: %w", err)
	}

	return piList, nil
}

// CreateChargeRequest 创建收费请求
type CreateChargeRequest struct {
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	CustomerID  string `json:"customer_id"`
	Source      string `json:"source"`
	Description string `json:"description"`
}

// CreateCharge 创建直接收费
func (s *PaymentService) CreateCharge(ctx context.Context, req *CreateChargeRequest) (*stripe.Charge, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Stripe client not initialized")
	}

	params := &stripe.ChargeParams{
		Amount:      stripe.Int64(req.Amount),
		Currency:    stripe.String(req.Currency),
		Customer:    stripe.String(req.CustomerID),
		Source:      stripe.String(req.Source),
		Description: stripe.String(req.Description),
	}

	charge, err := stripe.NewClient(s.client.GetSecretKey()).Charges.New(ctx, params)
	if err != nil {
		log.Printf("Failed to create charge: %v", err)
		return nil, fmt.Errorf("failed to create charge: %w", err)
	}

	log.Printf("Charge created successfully: %s", charge.ID)
	return charge, nil
}

// GetCharge 获取收费详情
func (s *PaymentService) GetCharge(ctx context.Context, chargeID string) (*stripe.Charge, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Stripe client not initialized")
	}

	charge, err := stripe.NewClient(s.client.GetSecretKey()).Charges.Get(ctx, chargeID, nil)
	if err != nil {
		log.Printf("Failed to get charge %s: %v", chargeID, err)
		return nil, fmt.Errorf("failed to get charge: %w", err)
	}

	return charge, nil
}

// RefundChargeRequest 退款请求
type RefundChargeRequest struct {
	ChargeID    string  `json:"charge_id"`
	Amount      int64   `json:"amount"`      // 可选,不指定则全额退款
	Reason      string  `json:"reason"`      // 退款原因
	Metadata    map[string]string `json:"metadata"`
}

// RefundCharge 退款
func (s *PaymentService) RefundCharge(ctx context.Context, req *RefundChargeRequest) (*stripe.Refund, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Stripe client not initialized")
	}

	params := &stripe.RefundParams{
		Charge: stripe.String(req.ChargeID),
	}

	if req.Amount > 0 {
		params.Amount = stripe.Int64(req.Amount)
	}
	if req.Reason != "" {
		params.Reason = stripe.String(req.Reason)
	}
	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}

	refund, err := stripe.NewClient(s.client.GetSecretKey()).Refunds.New(ctx, params)
	if err != nil {
		log.Printf("Failed to refund charge %s: %v", req.ChargeID, err)
		return nil, fmt.Errorf("failed to refund charge: %w", err)
	}

	log.Printf("Refund processed successfully: %s", refund.ID)
	return refund, nil
}

// CreateCustomerRequest 创建客户请求
type CreateCustomerRequest struct {
	Email       string            `json:"email"`
	Name        string            `json:"name"`
	Phone       string            `json:"phone"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
}

// CreateCustomer 创建客户
func (s *PaymentService) CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (*stripe.Customer, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Stripe client not initialized")
	}

	params := &stripe.CustomerParams{
		Email:       stripe.String(req.Email),
		Description: stripe.String(req.Description),
	}

	if req.Name != "" {
		params.Name = stripe.String(req.Name)
	}
	if req.Phone != "" {
		params.Phone = stripe.String(req.Phone)
	}
	if req.Metadata != nil {
		params.Metadata = req.Metadata
	}

	customer, err := stripe.NewClient(s.client.GetSecretKey()).Customers.New(ctx, params)
	if err != nil {
		log.Printf("Failed to create customer: %v", err)
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	log.Printf("Customer created successfully: %s", customer.ID)
	return customer, nil
}

// GetCustomer 获取客户
func (s *PaymentService) GetCustomer(ctx context.Context, customerID string) (*stripe.Customer, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Stripe client not initialized")
	}

	customer, err := stripe.NewClient(s.client.GetSecretKey()).Customers.Get(ctx, customerID, nil)
	if err != nil {
		log.Printf("Failed to get customer %s: %v", customerID, err)
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return customer, nil
}

// ConstructWebhookEvent 构建Webhook事件
func (s *PaymentService) ConstructWebhookEvent(payload []byte, sigHeader string) (stripe.Event, error) {
	if s.client == nil {
		return stripe.Event{}, fmt.Errorf("Stripe client not initialized")
	}

	event, err := webhook.NewEvent(payload, sigHeader, s.client.GetWebhookSecret())
	if err != nil {
		log.Printf("Failed to construct webhook event: %v", err)
		return stripe.Event{}, fmt.Errorf("failed to construct webhook event: %w", err)
	}

	return event, nil
}
