package controllers

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"github.com/clarkzhu2020/aidecms/pkg/database"
	"github.com/clarkzhu2020/aidecms/pkg/stripe"
	"github.com/clarkzhu2020/aidecms/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

// StripePaymentController Stripe支付控制器
type StripePaymentController struct{}

// NewStripePaymentController 创建Stripe支付控制器
func NewStripePaymentController() *StripePaymentController {
	return &StripePaymentController{}
}

// CreatePaymentIntentRequest 创建支付意图请求
type CreatePaymentIntentRequest struct {
	OrderID          string            `json:"order_id" validate:"required"`
	Amount           int64             `json:"amount" validate:"required,gt=0"`
	Currency         string            `json:"currency" validate:"required"`
	Description      string            `json:"description" validate:"required"`
	CustomerEmail    string            `json:"customer_email"`
	PaymentMethodTypes []string        `json:"payment_method_types"`
	Metadata         map[string]string `json:"metadata"`
}

// CreatePaymentIntent 创建Stripe支付意图
// @Summary      创建Stripe支付意图
// @Description  创建一个新的Stripe PaymentIntent
// @Tags         StripePayments
// @Accept       json
// @Produce      json
// @Param        payment body CreatePaymentIntentRequest true "支付信息"
// @Success      201 {object} response.Response{data=models.StripePayment}
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /stripe/payments/intent [post]
func (c *StripePaymentController) CreatePaymentIntent(ctx context.Context, hCtx *app.RequestContext) {
	var req CreatePaymentIntentRequest
	if err := hCtx.BindJSON(&req); err != nil {
		response.BadRequest(hCtx, "Invalid request data")
		return
	}

	// 检查Stripe客户端是否初始化
	if !stripe.IsInitialized() {
		response.ServerError(hCtx, "Stripe service not configured")
		return
	}

	// 获取当前用户ID（从JWT中间件设置）
	userID, _ := hCtx.Get("user_id")
	var userIDPtr *uint
	if uid, ok := userID.(uint); ok {
		userIDPtr = &uid
	}

	// 创建客户（如果提供了邮箱）
	var customerID string
	if req.CustomerEmail != "" {
		customerReq := &stripe.CreateCustomerRequest{
			Email:       req.CustomerEmail,
			Description: "Customer from AideCMS",
		}
		customer, err := stripe.NewPaymentService().CreateCustomer(ctx, customerReq)
		if err == nil {
			customerID = customer.ID
		}
	}

	// 调用Stripe服务创建支付意图
	intentReq := &stripe.CreatePaymentIntentRequest{
		Amount:           req.Amount,
		Currency:         req.Currency,
		CustomerID:       customerID,
		Description:      req.Description,
		ReceiptEmail:     req.CustomerEmail,
		PaymentMethodTypes: req.PaymentMethodTypes,
		Metadata:         req.Metadata,
	}

	// 添加订单ID到元数据
	if intentReq.Metadata == nil {
		intentReq.Metadata = make(map[string]string)
	}
	intentReq.Metadata["order_id"] = req.OrderID

	intent, err := stripe.NewPaymentService().CreatePaymentIntent(ctx, intentReq)
	if err != nil {
		log.Printf("Failed to create Stripe payment intent: %v", err)
		response.ServerError(hCtx, "Failed to create payment intent")
		return
	}

	// 保存支付记录到数据库
	db := database.GetDB()
	payment := &models.StripePayment{
		OrderID:         req.OrderID,
		PaymentIntentID:  intent.ID,
		Amount:          intent.Amount,
		Currency:        string(intent.Currency),
		Status:          string(intent.Status),
		PaymentStatus:   string(intent.Status),
		CustomerID:      customerID,
		CustomerEmail:   req.CustomerEmail,
		Description:     req.Description,
		UserID:          userIDPtr,
	}

	// 保存元数据
	if intent.Metadata != nil {
		if metadataBytes, err := json.Marshal(intent.Metadata); err == nil {
			payment.Metadata = string(metadataBytes)
		}
	}

	if err := db.Create(payment).Error; err != nil {
		log.Printf("Failed to save payment record: %v", err)
		response.ServerError(hCtx, "Failed to save payment record")
		return
	}

	// 返回支付意图信息
	response.Created(hCtx, map[string]interface{}{
		"payment_id":       payment.ID,
		"order_id":         payment.OrderID,
		"payment_intent_id": payment.PaymentIntentID,
		"client_secret":     intent.ClientSecret,
		"amount":           payment.Amount,
		"currency":         payment.Currency,
		"status":           payment.Status,
	}, "Payment intent created successfully")
}

// GetPayment 获取支付详情
// @Summary      获取Stripe支付详情
// @Description  根据ID获取支付详细信息
// @Tags         StripePayments
// @Accept       json
// @Produce      json
// @Param        id path int true "支付ID"
// @Success      200 {object} response.Response{data=models.StripePayment}
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /stripe/payments/{id} [get]
func (c *StripePaymentController) GetPayment(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")
	paymentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.BadRequest(hCtx, "Invalid payment ID")
		return
	}

	db := database.GetDB()
	var payment models.StripePayment
	if err := db.First(&payment, paymentID).Error; err != nil {
		response.NotFound(hCtx, "Payment not found")
		return
	}

	response.Success(hCtx, payment, "Payment fetched successfully")
}

// ListPayments 获取支付列表
// @Summary      获取Stripe支付列表
// @Description  获取支付列表，支持分页和筛选
// @Tags         StripePayments
// @Accept       json
// @Produce      json
// @Param        status query string false "支付状态"
// @Param        page query int false "页码" default(1)
// @Param        limit query int false "每页数量" default(20)
// @Success      200 {object} response.Response{data=[]models.StripePayment}
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /stripe/payments [get]
func (c *StripePaymentController) ListPayments(ctx context.Context, hCtx *app.RequestContext) {
	db := database.GetDB()

	status := string(hCtx.Query("status"))
	page, _ := strconv.Atoi(hCtx.Query("page", "1"))
	limit, _ := strconv.Atoi(hCtx.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var payments []models.StripePayment
	query := db.Model(&models.StripePayment{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 分页
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit).Order("created_at DESC")

	if err := query.Find(&payments).Error; err != nil {
		response.ServerError(hCtx, "Failed to fetch payments")
		return
	}

	response.Success(hCtx, map[string]interface{}{
		"data":   payments,
		"page":   page,
		"limit":  limit,
		"total":  len(payments),
	}, "Payments fetched successfully")
}

// ConfirmPayment 确认支付
// @Summary      确认支付
// @Description  使用支付方法确认Stripe支付意图
// @Tags         StripePayments
// @Accept       json
// @Produce      json
// @Param        paymentIntentID path string true "支付意图ID"
// @Param        payment_method body map[string]string true "支付方法ID"
// @Success      200 {object} response.Response{data=models.StripePayment}
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /stripe/payments/{paymentIntentID}/confirm [post]
func (c *StripePaymentController) ConfirmPayment(ctx context.Context, hCtx *app.RequestContext) {
	paymentIntentID := hCtx.Param("paymentIntentID")

	var req struct {
		PaymentMethodID string `json:"payment_method_id" validate:"required"`
	}
	if err := hCtx.BindJSON(&req); err != nil {
		response.BadRequest(hCtx, "Invalid request data")
		return
	}

	// 检查Stripe客户端是否初始化
	if !stripe.IsInitialized() {
		response.ServerError(hCtx, "Stripe service not configured")
		return
	}

	// 调用Stripe服务确认支付
	intent, err := stripe.NewPaymentService().ConfirmPaymentIntent(ctx, paymentIntentID, req.PaymentMethodID)
	if err != nil {
		log.Printf("Failed to confirm payment intent: %v", err)
		response.ServerError(hCtx, "Failed to confirm payment")
		return
	}

	// 更新支付记录
	db := database.GetDB()
	var payment models.StripePayment
	if err := db.Where("payment_intent_id = ?", paymentIntentID).First(&payment).Error; err == nil {
		payment.Status = string(intent.Status)
		payment.PaymentStatus = string(intent.Status)

		// 如果成功，保存Charge ID
		if intent.LatestCharge != nil {
			payment.ChargeID = intent.LatestCharge.ID
		}
		if intent.ReceiptURL != "" {
			payment.ReceiptURL = intent.ReceiptURL
		}

		db.Save(&payment)
	}

	response.Success(hCtx, map[string]interface{}{
		"payment_intent_id": paymentIntentID,
		"status":           string(intent.Status),
		"amount":           intent.Amount,
		"currency":         string(intent.Currency),
	}, "Payment confirmed successfully")
}

// RefundPayment 退款
// @Summary      退款
// @Description  对已支付的订单进行退款
// @Tags         StripePayments
// @Accept       json
// @Produce      json
// @Param        id path int true "支付ID"
// @Param        amount body map[string]interface{} false "退款金额 {\"amount\": 1000}"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /stripe/payments/{id}/refund [post]
func (c *StripePaymentController) RefundPayment(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")
	paymentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.BadRequest(hCtx, "Invalid payment ID")
		return
	}

	// 解析退款金额（可选）
	var req struct {
		Amount int64  `json:"amount"` // 可选,不填则全额退款
		Reason string `json:"reason"`
	}
	if err := hCtx.BindJSON(&req); err != nil {
		// 如果没有body，说明是全额退款
	}

	// 查找支付记录
	db := database.GetDB()
	var payment models.StripePayment
	if err := db.First(&payment, paymentID).Error; err != nil {
		response.NotFound(hCtx, "Payment not found")
		return
	}

	// 检查支付状态
	if !payment.IsSucceeded() {
		response.BadRequest(hCtx, "Payment is not eligible for refund")
		return
	}

	// 检查Stripe客户端是否初始化
	if !stripe.IsInitialized() {
		response.ServerError(hCtx, "Stripe service not configured")
		return
	}

	// 调用Stripe服务退款
	refundReq := &stripe.RefundChargeRequest{
		ChargeID: payment.ChargeID,
		Amount:   req.Amount,
		Reason:   req.Reason,
	}

	refund, err := stripe.NewPaymentService().RefundCharge(ctx, refundReq)
	if err != nil {
		log.Printf("Failed to refund payment: %v", err)
		response.ServerError(hCtx, "Failed to process refund")
		return
	}

	// 保存退款记录
	refundRecord := &models.StripeRefund{
		PaymentID:   payment.ID,
		RefundID:    refund.ID,
		ChargeID:    refund.Charge,
		Amount:      refund.Amount,
		Currency:    string(refund.Currency),
		Status:      string(refund.Status),
		Reason:      req.Reason,
		Description: payment.Description,
	}

	if err := db.Create(refundRecord).Error; err != nil {
		log.Printf("Failed to save refund record: %v", err)
		response.ServerError(hCtx, "Failed to save refund record")
		return
	}

	response.Success(hCtx, refundRecord, "Refund processed successfully")
}

// HandleWebhook 处理Stripe Webhook通知
// @Summary      处理Webhook
// @Description  接收并处理Stripe的Webhook通知
// @Tags         StripePayments
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /stripe/webhook [post]
func (c *StripePaymentController) HandleWebhook(ctx context.Context, hCtx *app.RequestContext) {
	// 读取请求体
	body := hCtx.Request.Body()
	sigHeader := string(hCtx.GetHeader("Stripe-Signature"))

	// 验证Webhook签名
	paymentService := stripe.NewPaymentService()
	event, err := paymentService.ConstructWebhookEvent(body, sigHeader)
	if err != nil {
		log.Printf("Failed to verify webhook signature: %v", err)
		response.BadRequest(hCtx, "Invalid webhook signature")
		return
	}

	// 保存Webhook记录
	db := database.GetDB()
	webhookRecord := &models.StripeWebhook{
		StripeID: event.ID,
		Type:     event.Type,
		Status:   "processed",
	}

	// 序列化事件数据
	if jsonData, err := json.Marshal(event.Data); err == nil {
		webhookRecord.Data = string(jsonData)
	}

	if err := db.Create(webhookRecord).Error; err != nil {
		log.Printf("Failed to save webhook record: %v", err)
	}

	// 根据事件类型处理
	switch event.Type {
	case "payment_intent.succeeded":
		c.handlePaymentIntentSucceeded(ctx, event)
	case "payment_intent.payment_failed":
		c.handlePaymentIntentFailed(ctx, event)
	case "charge.refunded":
		c.handleChargeRefunded(ctx, event)
	}

	response.Success(hCtx, map[string]interface{}{
		"status": "processed",
	}, "Webhook processed successfully")
}

// handlePaymentIntentSucceeded 处理支付成功事件
func (c *StripePaymentController) handlePaymentIntentSucceeded(ctx context.Context, event stripe.Event) {
	db := database.GetDB()

	var paymentIntent stripe.PaymentIntent
	if err := event.Data.GetObject(&paymentIntent); err != nil {
		log.Printf("Failed to parse payment intent: %v", err)
		return
	}

	// 更新支付记录
	var payment models.StripePayment
	if err := db.Where("payment_intent_id = ?", paymentIntent.ID).First(&payment).Error; err == nil {
		payment.Status = string(paymentIntent.Status)
		payment.PaymentStatus = string(paymentIntent.Status)

		if paymentIntent.LatestCharge != nil {
			payment.ChargeID = paymentIntent.LatestCharge.ID
		}
		if paymentIntent.ReceiptURL != "" {
			payment.ReceiptURL = paymentIntent.ReceiptURL
		}

		db.Save(&payment)
		log.Printf("Payment %d marked as succeeded via webhook", payment.ID)
	}
}

// handlePaymentIntentFailed 处理支付失败事件
func (c *StripePaymentController) handlePaymentIntentFailed(ctx context.Context, event stripe.Event) {
	db := database.GetDB()

	var paymentIntent stripe.PaymentIntent
	if err := event.Data.GetObject(&paymentIntent); err != nil {
		log.Printf("Failed to parse payment intent: %v", err)
		return
	}

	// 更新支付记录
	var payment models.StripePayment
	if err := db.Where("payment_intent_id = ?", paymentIntent.ID).First(&payment).Error; err == nil {
		payment.Status = string(paymentIntent.Status)
		payment.PaymentStatus = string(paymentIntent.Status)
		db.Save(&payment)
		log.Printf("Payment %d marked as failed via webhook", payment.ID)
	}
}

// handleChargeRefunded 处理退款事件
func (c *StripePaymentController) handleChargeRefunded(ctx context.Context, event stripe.Event) {
	db := database.GetDB()

	var charge stripe.Charge
	if err := event.Data.GetObject(&charge); err != nil {
		log.Printf("Failed to parse charge: %v", err)
		return
	}

	// 查找并更新支付记录
	var payment models.StripePayment
	if err := db.Where("charge_id = ?", charge.ID).First(&payment).Error; err == nil {
		// 这里可以添加退款逻辑
		log.Printf("Refund received for payment %d", payment.ID)
	}
}
