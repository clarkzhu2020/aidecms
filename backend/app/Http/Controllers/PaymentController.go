package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"github.com/clarkzhu2020/aidecms/pkg/database"
	"github.com/clarkzhu2020/aidecms/pkg/paypal"
	"github.com/clarkzhu2020/aidecms/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

// PaymentController 支付控制器
type PaymentController struct{}

// NewPaymentController 创建支付控制器
func NewPaymentController() *PaymentController {
	return &PaymentController{}
}

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	OrderID      string  `json:"order_id" validate:"required"`
	Amount       float64 `json:"amount" validate:"required,gt=0"`
	Currency     string  `json:"currency" validate:"required,default=USD"`
	Description  string  `json:"description" validate:"required"`
	ItemName     string  `json:"item_name"`
	ItemQuantity int     `json:"item_quantity" default="1"`
}

// CreatePayment 创建支付订单
// @Summary      创建PayPal支付订单
// @Description  创建一个新的PayPal支付订单
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        payment body CreatePaymentRequest true "支付信息"
// @Success      201 {object} response.Response{data=models.Payment}
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /payments [post]
func (c *PaymentController) CreatePayment(ctx context.Context, hCtx *app.RequestContext) {
	var req CreatePaymentRequest
	if err := hCtx.BindJSON(&req); err != nil {
		response.BadRequest(hCtx, "Invalid request data")
		return
	}

	// 检查PayPal客户端是否初始化
	if !paypal.IsInitialized() {
		response.ServerError(hCtx, "PayPal service not configured")
		return
	}

	// 获取当前用户ID（从JWT中间件设置）
	userID, _ := hCtx.Get("user_id")
	var userIDPtr *uint
	if uid, ok := userID.(uint); ok {
		userIDPtr = &uid
	}

	// 构建回调URL
	baseURL := "http://localhost:8888/api"
	returnURL := fmt.Sprintf("%s/payments/success", baseURL)
	cancelURL := fmt.Sprintf("%s/payments/cancel", baseURL)

	// 生成参考ID
	referenceID := fmt.Sprintf("AIDECMS-%s", req.OrderID)

	// 调用PayPal服务创建订单
	orderService := paypal.NewOrderService()
	orderReq := &paypal.CreateOrderRequest{
		OrderID:      req.OrderID,
		Amount:       req.Amount,
		Currency:     req.Currency,
		Description:  req.Description,
		ReturnURL:    returnURL,
		CancelURL:    cancelURL,
		ReferenceID:  referenceID,
		ItemName:     req.ItemName,
		ItemQuantity: req.ItemQuantity,
	}

	order, err := orderService.CreateOrder(ctx, orderReq)
	if err != nil {
		log.Printf("Failed to create PayPal order: %v", err)
		response.ServerError(hCtx, "Failed to create payment order")
		return
	}

	// 提取审批链接
	var approvalURL string
	for _, link := range order.Links {
		if link.Rel == "approve" {
			approvalURL = link.Href
			break
		}
	}

	// 保存支付记录到数据库
	db := database.GetDB()
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
		ReturnURL:     returnURL,
		CancelURL:     cancelURL,
		UserID:        userIDPtr,
	}

	if err := db.Create(payment).Error; err != nil {
		log.Printf("Failed to save payment record: %v", err)
		response.ServerError(hCtx, "Failed to save payment record")
		return
	}

	// 返回支付信息，包含PayPal审批链接
	response.Created(hCtx, map[string]interface{}{
		"payment_id":   payment.ID,
		"order_id":     payment.OrderID,
		"paypal_order_id": payment.PayPalOrderID,
		"approval_url": payment.ApprovalURL,
		"amount":       payment.Amount,
		"currency":     payment.Currency,
		"status":       payment.Status,
	}, "Payment order created successfully")
}

// GetPayment 获取支付详情
// @Summary      获取支付详情
// @Description  根据ID获取支付详细信息
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id path int true "支付ID"
// @Success      200 {object} response.Response{data=models.Payment}
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /payments/{id} [get]
func (c *PaymentController) GetPayment(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")
	paymentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.BadRequest(hCtx, "Invalid payment ID")
		return
	}

	db := database.GetDB()
	var payment models.Payment
	if err := db.First(&payment, paymentID).Error; err != nil {
		response.NotFound(hCtx, "Payment not found")
		return
	}

	response.Success(hCtx, payment, "Payment fetched successfully")
}

// CapturePayment 捕获支付（确认支付）
// @Summary      捕获支付
// @Description  捕获PayPal订单支付
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        orderID path string true "PayPal订单ID"
// @Success      200 {object} response.Response{data=models.Payment}
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /payments/capture/{orderID} [post]
func (c *PaymentController) CapturePayment(ctx context.Context, hCtx *app.RequestContext) {
	orderID := hCtx.Param("orderID")

	// 检查PayPal客户端是否初始化
	if !paypal.IsInitialized() {
		response.ServerError(hCtx, "PayPal service not configured")
		return
	}

	// 查找支付记录
	db := database.GetDB()
	var payment models.Payment
	if err := db.Where("paypal_order_id = ?", orderID).First(&payment).Error; err != nil {
		response.NotFound(hCtx, "Payment not found")
		return
	}

	// 检查支付状态
	if payment.Status == "paid" {
		response.BadRequest(hCtx, "Payment already captured")
		return
	}

	// 调用PayPal服务捕获支付
	orderService := paypal.NewOrderService()
	capture, err := orderService.CaptureOrder(ctx, orderID)
	if err != nil {
		log.Printf("Failed to capture PayPal order: %v", err)
		response.ServerError(hCtx, "Failed to capture payment")
		return
	}

	// 更新支付记录
	updates := map[string]interface{}{
		"status":         "paid",
		"payment_status": capture.Status,
	}

	// 提取capture ID
	if capture.PurchaseUnits != nil && len(capture.PurchaseUnits) > 0 {
		if captures := capture.PurchaseUnits[0].Payments.Captures; len(captures) > 0 {
			updates["capture_id"] = captures[0].ID
		}
	}

	if err := db.Model(&payment).Updates(updates).Error; err != nil {
		log.Printf("Failed to update payment record: %v", err)
		response.ServerError(hCtx, "Failed to update payment record")
		return
	}

	// 重新加载支付记录
	db.First(&payment, payment.ID)

	response.Success(hCtx, payment, "Payment captured successfully")
}

// RefundPayment 退款
// @Summary      退款
// @Description  对已支付的订单进行退款
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        id path int true "支付ID"
// @Param        amount body map[string]interface{} false "退款金额 {\"amount\": 10.00}"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /payments/{id}/refund [post]
func (c *PaymentController) RefundPayment(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")
	paymentID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.BadRequest(hCtx, "Invalid payment ID")
		return
	}

	// 解析退款金额（可选）
	var req struct {
		Amount   *float64 `json:"amount"`
		Reason   string   `json:"reason"`
		Note     string   `json:"note"`
	}
	if err := hCtx.BindJSON(&req); err != nil {
		// 如果没有body，说明是全额退款
	}

	// 查找支付记录
	db := database.GetDB()
	var payment models.Payment
	if err := db.First(&payment, paymentID).Error; err != nil {
		response.NotFound(hCtx, "Payment not found")
		return
	}

	// 检查支付状态
	if payment.Status != "paid" {
		response.BadRequest(hCtx, "Payment is not eligible for refund")
		return
	}

	// 检查PayPal客户端是否初始化
	if !paypal.IsInitialized() {
		response.ServerError(hCtx, "PayPal service not configured")
		return
	}

	// 调用PayPal服务退款
	orderService := paypal.NewOrderService()
	refund, err := orderService.RefundPayment(ctx, payment.CaptureID, req.Amount, payment.Currency)
	if err != nil {
		log.Printf("Failed to refund payment: %v", err)
		response.ServerError(hCtx, "Failed to process refund")
		return
	}

	// 保存退款记录
	refundRecord := &models.PaymentRefund{
		PaymentID: payment.ID,
		RefundID:  refund.ID,
		Amount:    payment.Amount,
		Currency:  payment.Currency,
		Status:    refund.Status,
		Reason:    req.Reason,
		Note:      req.Note,
	}

	if err := db.Create(refundRecord).Error; err != nil {
		log.Printf("Failed to save refund record: %v", err)
		response.ServerError(hCtx, "Failed to save refund record")
		return
	}

	// 如果是全额退款，更新支付状态
	if req.Amount == nil || *req.Amount >= payment.Amount {
		db.Model(&payment).Update("status", "refunded")
	}

	response.Success(hCtx, refundRecord, "Refund processed successfully")
}

// ListPayments 获取支付列表
// @Summary      获取支付列表
// @Description  获取支付列表，支持分页和筛选
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        status query string false "支付状态" Enums(pending, paid, failed, cancelled, refunded)
// @Param        page query int false "页码" default(1)
// @Param        limit query int false "每页数量" default(20)
// @Success      200 {object} response.Response{data=[]models.Payment}
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /payments [get]
func (c *PaymentController) ListPayments(ctx context.Context, hCtx *app.RequestContext) {
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

	var payments []models.Payment
	query := db.Model(&models.Payment{})

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

// HandleWebhook 处理PayPal Webhook通知
// @Summary      处理Webhook
// @Description  接收并处理PayPal的Webhook通知
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        webhook body map[string]interface{} true "Webhook数据"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /payments/webhook [post]
func (c *PaymentController) HandleWebhook(ctx context.Context, hCtx *app.RequestContext) {
	// 读取请求体
	body := hCtx.Request.Body()
	var webhookData map[string]interface{}
	if err := json.Unmarshal(body, &webhookData); err != nil {
		log.Printf("Failed to parse webhook data: %v", err)
		response.BadRequest(hCtx, "Invalid webhook data")
		return
	}

	// 提取Webhook信息
	eventID := ""
	eventType := ""
	resourceType := ""
	resourceID := ""
	summary := ""

	if id, ok := webhookData["id"].(string); ok {
		eventID = id
	}
	if et, ok := webhookData["event_type"].(string); ok {
		eventType = et
	}
	if rt, ok := webhookData["resource_type"].(string); ok {
		resourceType = rt
	}
	if s, ok := webhookData["summary"].(string); ok {
		summary = s
	}

	// 提取资源ID
	if resource, ok := webhookData["resource"].(map[string]interface{}); ok {
		if id, ok := resource["id"].(string); ok {
			resourceID = id
		}
	}

	// 保存Webhook记录
	db := database.GetDB()
	webhook := &models.PaymentWebhook{
		EventID:      eventID,
		EventType:    eventType,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Summary:      summary,
		Status:       "processed",
	}

	// 序列化原始数据
	if jsonData, err := json.Marshal(webhookData); err == nil {
		webhook.RawData = string(jsonData)
	}

	// 序列化请求头
	headers := make(map[string]string)
	hCtx.Request.Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})
	if headerData, err := json.Marshal(headers); err == nil {
		webhook.RawHeaders = string(headerData)
	}

	if err := db.Create(webhook).Error; err != nil {
		log.Printf("Failed to save webhook record: %v", err)
	}

	// 根据事件类型处理
	switch eventType {
	case "PAYMENT.CAPTURE.COMPLETED":
		c.handlePaymentCaptureCompleted(ctx, resourceID, webhookData)
	case "PAYMENT.CAPTURE.DENIED":
		c.handlePaymentCaptureDenied(ctx, resourceID, webhookData)
	}

	response.Success(hCtx, map[string]interface{}{
		"status": "processed",
	}, "Webhook processed successfully")
}

// handlePaymentCaptureCompleted 处理支付完成事件
func (c *PaymentController) handlePaymentCaptureCompleted(ctx context.Context, resourceID string, webhookData map[string]interface{}) {
	db := database.GetDB()

	// 查找对应的支付记录
	var payment models.Payment
	if err := db.Where("capture_id = ?", resourceID).First(&payment).Error; err != nil {
		log.Printf("Payment not found for capture ID: %s", resourceID)
		return
	}

	// 更新支付状态
	db.Model(&payment).Update("status", "paid")
	log.Printf("Payment %d marked as paid via webhook", payment.ID)
}

// handlePaymentCaptureDenied 处理支付拒绝事件
func (c *PaymentController) handlePaymentCaptureDenied(ctx context.Context, resourceID string, webhookData map[string]interface{}) {
	db := database.GetDB()

	// 查找对应的支付记录
	var payment models.Payment
	if err := db.Where("capture_id = ?", resourceID).First(&payment).Error; err != nil {
		log.Printf("Payment not found for capture ID: %s", resourceID)
		return
	}

	// 更新支付状态
	db.Model(&payment).Update("status", "failed")
	log.Printf("Payment %d marked as failed via webhook", payment.ID)
}
