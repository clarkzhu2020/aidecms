package paypal

import (
	"context"
	"fmt"
	"log"

	"github.com/clarkzhu2020/aidecms/internal/app/models"
)

// PaymentService 支付业务服务
type PaymentService struct {
	orderService *OrderService
	repository   *Repository
}

// NewPaymentService 创建支付业务服务
func NewPaymentService() *PaymentService {
	return &PaymentService{
		orderService: NewOrderService(),
		repository:   NewRepository(),
	}
}

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	OrderID      string  `json:"order_id" validate:"required"`
	Amount       float64 `json:"amount" validate:"required,gt=0"`
	Currency     string  `json:"currency" validate:"required"`
	Description  string  `json:"description" validate:"required"`
	ReturnURL    string  `json:"return_url" validate:"required"`
	CancelURL    string  `json:"cancel_url" validate:"required"`
	UserID       *uint   `json:"user_id"`
	ItemName     string  `json:"item_name"`
	ItemQuantity int     `json:"item_quantity"`
}

// CreatePayment 创建支付（完整的业务流程）
func (s *PaymentService) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*models.Payment, error) {
	// 生成参考ID
	referenceID := fmt.Sprintf("AIDECMS-%s", req.OrderID)

	// 调用PayPal API创建订单
	orderReq := &CreateOrderRequest{
		OrderID:      req.OrderID,
		Amount:       req.Amount,
		Currency:     req.Currency,
		Description:  req.Description,
		ReturnURL:    req.ReturnURL,
		CancelURL:    req.CancelURL,
		ReferenceID:  referenceID,
		ItemName:     req.ItemName,
		ItemQuantity: req.ItemQuantity,
	}

	order, err := s.orderService.CreateOrder(ctx, orderReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create PayPal order: %w", err)
	}

	// 提取审批链接
	var approvalURL string
	for _, link := range order.Links {
		if link.Rel == "approve" {
			approvalURL = link.Href
			break
		}
	}

	// 创建支付记录
	payment := &models.Payment{
		OrderID:       req.OrderID,
		PayPalOrderID: order.ID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Status:        "pending",
		PaymentStatus: order.Status,
		Description:   req.Description,
		ReferenceID:   referenceID,
		ApprovalURL:   approvalURL,
		ReturnURL:     req.ReturnURL,
		CancelURL:     req.CancelURL,
		UserID:        req.UserID,
	}

	// 保存到数据库
	if err := s.repository.CreatePayment(payment); err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	log.Printf("Payment created successfully: %d", payment.ID)
	return payment, nil
}

// ProcessPayment 处理支付（捕获）
func (s *PaymentService) ProcessPayment(ctx context.Context, paypalOrderID string) (*models.Payment, error) {
	// 查找支付记录
	payment, err := s.repository.GetPaymentByPayPalOrderID(paypalOrderID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	// 检查支付状态
	if payment.Status == "paid" {
		return payment, nil // 已经支付过
	}

	// 调用PayPal API捕获支付
	capture, err := s.orderService.CaptureOrder(ctx, paypalOrderID)
	if err != nil {
		// 更新支付状态为失败
		s.repository.UpdatePaymentStatus(payment.ID, "failed")
		return nil, fmt.Errorf("failed to capture PayPal order: %w", err)
	}

	// 提取capture ID
	var captureID string
	if capture.PurchaseUnits != nil && len(capture.PurchaseUnits) > 0 {
		if captures := capture.PurchaseUnits[0].Payments.Captures; len(captures) > 0 {
			captureID = captures[0].ID
		}
	}

	// 更新支付记录
	updates := map[string]interface{}{
		"status":         "paid",
		"payment_status": capture.Status,
		"capture_id":     captureID,
	}

	// 提取支付者信息
	if capture.Payer != nil {
		if capture.Payer.PayerID != "" {
			updates["payer_id"] = capture.Payer.PayerID
		}
		if capture.Payer.EmailAddress != "" {
			updates["payer_email"] = capture.Payer.EmailAddress
		}
	}

	payment.Status = "paid"
	payment.PaymentStatus = capture.Status
	payment.CaptureID = captureID

	if err := s.repository.UpdatePayment(payment); err != nil {
		log.Printf("Failed to update payment record: %v", err)
		return nil, fmt.Errorf("payment captured but failed to update record: %w", err)
	}

	log.Printf("Payment processed successfully: %d", payment.ID)
	return payment, nil
}

// RefundPayment 退款
func (s *PaymentService) RefundPayment(ctx context.Context, paymentID uint, amount *float64, reason, note string) (*models.PaymentRefund, error) {
	// 查找支付记录
	payment, err := s.repository.GetPaymentByID(paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment not found: %w", err)
	}

	// 检查支付状态
	if payment.Status != "paid" {
		return nil, fmt.Errorf("payment is not eligible for refund")
	}

	// 检查是否已全额退款
	if payment.Status == "refunded" {
		return nil, fmt.Errorf("payment already refunded")
	}

	// 计算退款金额
	var refundAmount float64
	if amount != nil {
		refundAmount = *amount
		if refundAmount > payment.Amount {
			return nil, fmt.Errorf("refund amount exceeds payment amount")
		}
	} else {
		refundAmount = payment.Amount // 全额退款
	}

	// 调用PayPal API退款
	refund, err := s.orderService.RefundPayment(ctx, payment.CaptureID, amount, payment.Currency)
	if err != nil {
		return nil, fmt.Errorf("failed to refund payment: %w", err)
	}

	// 创建退款记录
	refundRecord := &models.PaymentRefund{
		PaymentID: payment.ID,
		RefundID:  refund.ID,
		Amount:    refundAmount,
		Currency:  payment.Currency,
		Status:    refund.Status,
		Reason:    reason,
		Note:      note,
	}

	if err := s.repository.CreateRefund(refundRecord); err != nil {
		return nil, fmt.Errorf("refund processed but failed to create record: %w", err)
	}

	// 如果是全额退款，更新支付状态
	if amount == nil || *amount >= payment.Amount {
		s.repository.UpdatePaymentStatus(payment.ID, "refunded")
	}

	log.Printf("Refund processed successfully: %s", refund.ID)
	return refundRecord, nil
}

// CancelPayment 取消支付
func (s *PaymentService) CancelPayment(ctx context.Context, paypalOrderID string) error {
	// 查找支付记录
	payment, err := s.repository.GetPaymentByPayPalOrderID(paypalOrderID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	// 检查支付状态
	if payment.Status != "pending" {
		return fmt.Errorf("payment cannot be cancelled in current status: %s", payment.Status)
	}

	// 更新支付状态为已取消
	if err := s.repository.UpdatePaymentStatus(payment.ID, "cancelled"); err != nil {
		return fmt.Errorf("failed to cancel payment: %w", err)
	}

	log.Printf("Payment cancelled successfully: %d", payment.ID)
	return nil
}

// GetPaymentStats 获取支付统计信息
func (s *PaymentService) GetPaymentStats() (map[string]interface{}, error) {
	return s.repository.GetPaymentStats()
}

// SearchPayments 搜索支付记录
func (s *PaymentService) SearchPayments(keyword string, page, limit int) ([]models.Payment, int64, error) {
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	return s.repository.SearchPayments(keyword, offset, limit)
}

// GetPaymentHistory 获取支付历史
func (s *PaymentService) GetPaymentHistory(userID uint, page, limit int) ([]models.Payment, int64, error) {
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	return s.repository.ListPaymentsByUserID(userID, offset, limit)
}

// GetRefundHistory 获取退款历史
func (s *PaymentService) GetRefundHistory(paymentID uint) ([]models.PaymentRefund, error) {
	return s.repository.ListRefundsByPaymentID(paymentID)
}
