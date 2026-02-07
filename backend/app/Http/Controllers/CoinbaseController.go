package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"github.com/clarkzhu2020/aidecms/pkg/coinbase"
	"github.com/clarkzhu2020/aidecms/pkg/database"
	"github.com/clarkzhu2020/aidecms/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

// CoinbaseController Coinbase控制器
type CoinbaseController struct{}

// NewCoinbaseController 创建Coinbase控制器
func NewCoinbaseController() *CoinbaseController {
	return &CoinbaseController{}
}

// ========== Payment Link API ==========

// CreatePaymentLinkRequest 创建支付链接请求
type CreatePaymentLinkRequest struct {
	Amount      string                 `json:"amount" validate:"required"`
	Currency    string                 `json:"currency" validate:"required"`
	Description string                 `json:"description"`
	Title       string                 `json:"title"`
	RedirectURL string                 `json:"redirectUrl"`
	CancelURL   string                 `json:"cancelUrl"`
	Name        string                 `json:"name"`
	Email       string                 `json:"email"`
	Metadata    map[string]interface{} `json:"metadata"`
	ExternalID  string                 `json:"externalId"`
}

// CreatePaymentLink 创建Coinbase支付链接
// @Summary      创建支付链接
// @Description  创建一个新的Coinbase加密货币支付链接
// @Tags         Coinbase-Payment
// @Accept       json
// @Produce      json
// @Param        paymentLink body CreatePaymentLinkRequest true "支付链接信息"
// @Success      201 {object} response.Response{data=models.CoinbasePaymentLink}
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /coinbase/payment-links [post]
func (c *CoinbaseController) CreatePaymentLink(ctx context.Context, hCtx *app.RequestContext) {
	var req CreatePaymentLinkRequest
	if err := hCtx.BindJSON(&req); err != nil {
		response.BadRequest(hCtx, "Invalid request data")
		return
	}

	// 检查Coinbase客户端是否初始化
	if !coinbase.IsInitialized() {
		response.ServerError(hCtx, "Coinbase service not configured")
		return
	}

	// 获取当前用户ID
	userID, _ := hCtx.Get("user_id")
	var userIDPtr *uint
	if uid, ok := userID.(uint); ok {
		userIDPtr = &uid
	}

	// 构建回调URL
	if req.RedirectURL == "" {
		baseURL := "http://localhost:8888/api"
		req.RedirectURL = fmt.Sprintf("%s/coinbase/success", baseURL)
	}

	// 调用Coinbase服务创建支付链接
	paymentLinkService := coinbase.NewPaymentLinkService()
	coinbaseReq := &coinbase.PaymentLinkRequest{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Description: req.Description,
		Title:       req.Title,
		RedirectURL: req.RedirectURL,
		CancelURL:   req.CancelURL,
		Name:        req.Name,
		Email:       req.Email,
		Metadata:    req.Metadata,
	}

	paymentLink, err := paymentLinkService.CreatePaymentLink(ctx, coinbaseReq)
	if err != nil {
		log.Printf("Failed to create Coinbase payment link: %v", err)
		response.ServerError(hCtx, "Failed to create payment link")
		return
	}

	// 保存支付链接记录到数据库
	db := database.GetDB()
	paymentLinkRecord := &models.CoinbasePaymentLink{
		LinkID:        paymentLink.ID,
		ExternalID:    req.ExternalID,
		Amount:        paymentLink.Amount,
		Currency:      paymentLink.Currency,
		Title:         paymentLink.Title,
		Description:   paymentLink.Description,
		Status:        paymentLink.Status,
		PaymentStatus: paymentLink.PaymentStatus,
		PaymentURL:    paymentLink.PaymentURL,
		RedirectURL:   req.RedirectURL,
		CancelURL:     req.CancelURL,
		Name:          req.Name,
		Email:         req.Email,
		UserID:        userIDPtr,
	}

	// 序列化原始数据
	if rawData, err := json.Marshal(paymentLink); err == nil {
		paymentLinkRecord.RawData = string(rawData)
	}

	// 序列化metadata
	if req.Metadata != nil {
		if metadataData, err := json.Marshal(req.Metadata); err == nil {
			paymentLinkRecord.Metadata = string(metadataData)
		}
	}

	if err := db.Create(paymentLinkRecord).Error; err != nil {
		log.Printf("Failed to save payment link record: %v", err)
		response.ServerError(hCtx, "Failed to save payment link record")
		return
	}

	response.Created(hCtx, paymentLinkRecord, "Payment link created successfully")
}

// GetPaymentLink 获取支付链接详情
// @Summary      获取支付链接详情
// @Description  根据ID获取Coinbase支付链接详细信息
// @Tags         Coinbase-Payment
// @Accept       json
// @Produce      json
// @Param        id path int true "支付链接ID"
// @Success      200 {object} response.Response{data=models.CoinbasePaymentLink}
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /coinbase/payment-links/{id} [get]
func (c *CoinbaseController) GetPaymentLink(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")
	linkID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.BadRequest(hCtx, "Invalid payment link ID")
		return
	}

	db := database.GetDB()
	var paymentLink models.CoinbasePaymentLink
	if err := db.First(&paymentLink, linkID).Error; err != nil {
		response.NotFound(hCtx, "Payment link not found")
		return
	}

	response.Success(hCtx, paymentLink, "Payment link fetched successfully")
}

// ListPaymentLinks 获取支付链接列表
// @Summary      获取支付链接列表
// @Description  获取Coinbase支付链接列表，支持分页和筛选
// @Tags         Coinbase-Payment
// @Accept       json
// @Produce      json
// @Param        status query string false "状态"
// @Param        page query int false "页码" default(1)
// @Param        limit query int false "每页数量" default(20)
// @Success      200 {object} response.Response{data=[]models.CoinbasePaymentLink}
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /coinbase/payment-links [get]
func (c *CoinbaseController) ListPaymentLinks(ctx context.Context, hCtx *app.RequestContext) {
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

	var paymentLinks []models.CoinbasePaymentLink
	query := db.Model(&models.CoinbasePaymentLink{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit).Order("created_at DESC")

	if err := query.Find(&paymentLinks).Error; err != nil {
		response.ServerError(hCtx, "Failed to fetch payment links")
		return
	}

	response.Success(hCtx, map[string]interface{}{
		"data":  paymentLinks,
		"page":  page,
		"limit": limit,
		"total": len(paymentLinks),
	}, "Payment links fetched successfully")
}

// DeletePaymentLink 删除支付链接
// @Summary      删除支付链接
// @Description  删除指定的Coinbase支付链接
// @Tags         Coinbase-Payment
// @Accept       json
// @Produce      json
// @Param        id path int true "支付链接ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /coinbase/payment-links/{id} [delete]
func (c *CoinbaseController) DeletePaymentLink(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")
	linkID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.BadRequest(hCtx, "Invalid payment link ID")
		return
	}

	db := database.GetDB()
	var paymentLink models.CoinbasePaymentLink
	if err := db.First(&paymentLink, linkID).Error; err != nil {
		response.NotFound(hCtx, "Payment link not found")
		return
	}

	// 调用Coinbase API删除支付链接
	if coinbase.IsInitialized() {
		paymentLinkService := coinbase.NewPaymentLinkService()
		if err := paymentLinkService.DeletePaymentLink(ctx, paymentLink.LinkID); err != nil {
			log.Printf("Failed to delete Coinbase payment link: %v", err)
			// 即使API调用失败，也删除本地记录
		}
	}

	// 删除本地记录
	if err := db.Delete(&paymentLink).Error; err != nil {
		log.Printf("Failed to delete payment link record: %v", err)
		response.ServerError(hCtx, "Failed to delete payment link record")
		return
	}

	response.Success(hCtx, nil, "Payment link deleted successfully")
}

// ========== Trade API ==========

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	ProductID            string  `json:"product_id" validate:"required"`
	Side                 string  `json:"side" validate:"required,oneof=buy sell"`
	OrderType            string  `json:"order_type" validate:"required,oneof=market limit stop"`
	Size                 string  `json:"size"`
	Funds                string  `json:"funds"`
	LimitPrice           string  `json:"limit_price"`
	StopPrice            string  `json:"stop_price"`
	Stop                 string  `json:"stop" validate:"omitempty,oneof=loss entry"`
	TimeInForce          string  `json:"time_in_force" validate:"omitempty,oneof=GTC IOC FOK GTD"`
	ClientOrderID        string  `json:"client_order_id"`
	PostOnly             bool    `json:"post_only"`
	SelfTradePrevention  string  `json:"self_trade_prevention" validate:"omitempty,oneof=dc co cn cb"`
	ExternalID           string  `json:"external_id"`
}

// CreateOrder 创建交易订单
// @Summary      创建交易订单
// @Description  创建一个新的Coinbase交易订单
// @Tags         Coinbase-Trade
// @Accept       json
// @Produce      json
// @Param        order body CreateOrderRequest true "订单信息"
// @Success      201 {object} response.Response{data=models.CoinbaseOrder}
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /coinbase/orders [post]
func (c *CoinbaseController) CreateOrder(ctx context.Context, hCtx *app.RequestContext) {
	var req CreateOrderRequest
	if err := hCtx.BindJSON(&req); err != nil {
		response.BadRequest(hCtx, "Invalid request data")
		return
	}

	// 检查Coinbase客户端是否初始化
	if !coinbase.IsInitialized() {
		response.ServerError(hCtx, "Coinbase service not configured")
		return
	}

	// 获取当前用户ID
	userID, _ := hCtx.Get("user_id")
	var userIDPtr *uint
	if uid, ok := userID.(uint); ok {
		userIDPtr = &uid
	}

	// 调用Coinbase服务创建订单
	tradeService := coinbase.NewTradeService()
	coinbaseReq := &coinbase.OrderRequest{
		ProductID:            req.ProductID,
		Side:                 req.Side,
		OrderType:            req.OrderType,
		Size:                 req.Size,
		Funds:                req.Funds,
		LimitPrice:           req.LimitPrice,
		StopPrice:            req.StopPrice,
		Stop:                 req.Stop,
		TimeInForce:          req.TimeInForce,
		ClientOid:            req.ClientOrderID,
		PostOnly:             req.PostOnly,
		SelfTradePrevention:  req.SelfTradePrevention,
	}

	order, err := tradeService.CreateOrder(ctx, coinbaseReq)
	if err != nil {
		log.Printf("Failed to create Coinbase order: %v", err)
		response.ServerError(hCtx, "Failed to create order")
		return
	}

	// 保存订单记录到数据库
	db := database.GetDB()
	orderRecord := &models.CoinbaseOrder{
		OrderID:            order.ID,
		ClientOrderID:      order.ClientOid,
		ExternalID:         req.ExternalID,
		ProductID:          order.ProductID,
		Side:               order.Side,
		OrderType:          order.OrderType,
		Size:               order.Size,
		Funds:              order.Funds,
		LimitPrice:         order.LimitPrice,
		StopPrice:          order.StopPrice,
		FilledSize:         order.FilledSize,
		AverageFillPrice:   order.AverageFillPrice,
		FillFees:           order.FillFees,
		Status:             order.Status,
		Settled:            order.Settled,
		TimeInForce:        order.TimeInForce,
		PostOnly:           order.PostOnly,
		UserID:             userIDPtr,
	}

	// 序列化原始数据
	if rawData, err := json.Marshal(order); err == nil {
		orderRecord.RawData = string(rawData)
	}

	if err := db.Create(orderRecord).Error; err != nil {
		log.Printf("Failed to save order record: %v", err)
		response.ServerError(hCtx, "Failed to save order record")
		return
	}

	response.Created(hCtx, orderRecord, "Order created successfully")
}

// GetOrder 获取订单详情
// @Summary      获取订单详情
// @Description  根据ID获取Coinbase订单详细信息
// @Tags         Coinbase-Trade
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Success      200 {object} response.Response{data=models.CoinbaseOrder}
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /coinbase/orders/{id} [get]
func (c *CoinbaseController) GetOrder(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")
	orderID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.BadRequest(hCtx, "Invalid order ID")
		return
	}

	db := database.GetDB()
	var order models.CoinbaseOrder
	if err := db.First(&order, orderID).Error; err != nil {
		response.NotFound(hCtx, "Order not found")
		return
	}

	// 可选：从Coinbase API获取最新状态
	if coinbase.IsInitialized() {
		tradeService := coinbase.NewTradeService()
		if coinbaseOrder, err := tradeService.GetOrder(ctx, order.OrderID); err == nil {
			order.FilledSize = coinbaseOrder.FilledSize
			order.AverageFillPrice = coinbaseOrder.AverageFillPrice
			order.FillFees = coinbaseOrder.FillFees
			order.Status = coinbaseOrder.Status
			order.Settled = coinbaseOrder.Settled
			if rawData, err := json.Marshal(coinbaseOrder); err == nil {
				order.RawData = string(rawData)
			}
			db.Save(&order)
		}
	}

	response.Success(hCtx, order, "Order fetched successfully")
}

// ListOrders 获取订单列表
// @Summary      获取订单列表
// @Description  获取Coinbase订单列表，支持分页和筛选
// @Tags         Coinbase-Trade
// @Accept       json
// @Produce      json
// @Param        productId query string false "产品ID"
// @Param        status query string false "订单状态"
// @Param        page query int false "页码" default(1)
// @Param        limit query int false "每页数量" default(20)
// @Success      200 {object} response.Response{data=[]models.CoinbaseOrder}
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /coinbase/orders [get]
func (c *CoinbaseController) ListOrders(ctx context.Context, hCtx *app.RequestContext) {
	db := database.GetDB()

	productID := string(hCtx.Query("productId"))
	status := string(hCtx.Query("status"))
	page, _ := strconv.Atoi(hCtx.Query("page", "1"))
	limit, _ := strconv.Atoi(hCtx.Query("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var orders []models.CoinbaseOrder
	query := db.Model(&models.CoinbaseOrder{})

	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit).Order("created_at DESC")

	if err := query.Find(&orders).Error; err != nil {
		response.ServerError(hCtx, "Failed to fetch orders")
		return
	}

	response.Success(hCtx, map[string]interface{}{
		"data":  orders,
		"page":  page,
		"limit": limit,
		"total": len(orders),
	}, "Orders fetched successfully")
}

// CancelOrder 取消订单
// @Summary      取消订单
// @Description  取消指定的Coinbase交易订单
// @Tags         Coinbase-Trade
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /coinbase/orders/{id} [cancel]
func (c *CoinbaseController) CancelOrder(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")
	orderID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.BadRequest(hCtx, "Invalid order ID")
		return
	}

	db := database.GetDB()
	var order models.CoinbaseOrder
	if err := db.First(&order, orderID).Error; err != nil {
		response.NotFound(hCtx, "Order not found")
		return
	}

	// 调用Coinbase API取消订单
	if coinbase.IsInitialized() {
		tradeService := coinbase.NewTradeService()
		if err := tradeService.CancelOrder(ctx, order.OrderID); err != nil {
			log.Printf("Failed to cancel Coinbase order: %v", err)
			response.ServerError(hCtx, "Failed to cancel order")
			return
		}

		// 更新本地订单状态
		order.Status = "cancelled"
		db.Save(&order)
	}

	response.Success(hCtx, nil, "Order cancelled successfully")
}

// GetProducts 获取产品列表
// @Summary      获取产品列表
// @Description  获取Coinbase可交易的产品列表
// @Tags         Coinbase-Trade
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /coinbase/products [get]
func (c *CoinbaseController) GetProducts(ctx context.Context, hCtx *app.RequestContext) {
	if !coinbase.IsInitialized() {
		response.ServerError(hCtx, "Coinbase service not configured")
		return
	}

	tradeService := coinbase.NewTradeService()
	products, err := tradeService.GetProducts(ctx)
	if err != nil {
		log.Printf("Failed to get products: %v", err)
		response.ServerError(hCtx, "Failed to get products")
		return
	}

	response.Success(hCtx, products, "Products fetched successfully")
}

// GetProduct 获取产品详情
// @Summary      获取产品详情
// @Description  根据产品ID获取Coinbase产品详细信息
// @Tags         Coinbase-Trade
// @Accept       json
// @Produce      json
// @Param        productId path string true "产品ID"
// @Success      200 {object} response.Response
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /coinbase/products/{productId} [get]
func (c *CoinbaseController) GetProduct(ctx context.Context, hCtx *app.RequestContext) {
	if !coinbase.IsInitialized() {
		response.ServerError(hCtx, "Coinbase service not configured")
		return
	}

	productID := hCtx.Param("productId")
	tradeService := coinbase.NewTradeService()
	product, err := tradeService.GetProduct(ctx, productID)
	if err != nil {
		log.Printf("Failed to get product: %v", err)
		response.ServerError(hCtx, "Failed to get product")
		return
	}

	response.Success(hCtx, product, "Product fetched successfully")
}

// GetTicker 获取产品行情
// @Summary      获取产品行情
// @Description  根据产品ID获取Coinbase产品实时行情
// @Tags         Coinbase-Trade
// @Accept       json
// @Produce      json
// @Param        productId path string true "产品ID"
// @Success      200 {object} response.Response
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /coinbase/products/{productId}/ticker [get]
func (c *CoinbaseController) GetTicker(ctx context.Context, hCtx *app.RequestContext) {
	if !coinbase.IsInitialized() {
		response.ServerError(hCtx, "Coinbase service not configured")
		return
	}

	productID := hCtx.Param("productId")
	tradeService := coinbase.NewTradeService()
	ticker, err := tradeService.GetTicker(ctx, productID)
	if err != nil {
		log.Printf("Failed to get ticker: %v", err)
		response.ServerError(hCtx, "Failed to get ticker")
		return
	}

	response.Success(hCtx, ticker, "Ticker fetched successfully")
}

// HandleWebhook 处理Coinbase Webhook通知
// @Summary      处理Webhook
// @Description  接收并处理Coinbase的Webhook通知
// @Tags         Coinbase
// @Accept       json
// @Produce      json
// @Param        webhook body map[string]interface{} true "Webhook数据"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /coinbase/webhook [post]
func (c *CoinbaseController) HandleWebhook(ctx context.Context, hCtx *app.RequestContext) {
	// 读取请求体
	body := hCtx.Request.Body()
	var webhookData map[string]interface{}
	if err := json.Unmarshal(body, &webhookData); err != nil {
		log.Printf("Failed to parse webhook data: %v", err)
		response.BadRequest(hCtx, "Invalid webhook data")
		return
	}

	// 验证签名（如果配置了Webhook密钥）
	// signature := hCtx.GetHeader("X-CC-Webhook-Signature")
	// timestamp := hCtx.GetHeader("X-CC-Webhook-Timestamp")

	// 提取Webhook信息
	webhookID := ""
	eventType := ""
	linkID := ""
	orderID := ""
	summary := ""

	if id, ok := webhookData["id"].(string); ok {
		webhookID = id
	}
	if et, ok := webhookData["type"].(string); ok {
		eventType = et
	}
	if lid, ok := webhookData["payment_link_id"].(string); ok {
		linkID = lid
	}
	if oid, ok := webhookData["order_id"].(string); ok {
		orderID = oid
	}
	if s, ok := webhookData["summary"].(string); ok {
		summary = s
	}

	// 保存Webhook记录
	db := database.GetDB()
	webhook := &models.CoinbaseWebhook{
		WebhookID: webhookID,
		EventType:  eventType,
		LinkID:     linkID,
		OrderID:    orderID,
		Summary:    summary,
		Status:     "processed",
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
	case "payment_link.completed":
		c.handlePaymentLinkCompleted(ctx, linkID, webhookData)
	case "order.filled":
		c.handleOrderFilled(ctx, orderID, webhookData)
	case "order.cancelled":
		c.handleOrderCancelled(ctx, orderID, webhookData)
	}

	response.Success(hCtx, map[string]interface{}{
		"status": "processed",
	}, "Webhook processed successfully")
}

// handlePaymentLinkCompleted 处理支付链接完成事件
func (c *CoinbaseController) handlePaymentLinkCompleted(ctx context.Context, linkID string, webhookData map[string]interface{}) {
	db := database.GetDB()

	// 查找对应的支付链接记录
	var paymentLink models.CoinbasePaymentLink
	if err := db.Where("link_id = ?", linkID).First(&paymentLink).Error; err != nil {
		log.Printf("Payment link not found: %s", linkID)
		return
	}

	// 更新支付链接状态
	db.Model(&paymentLink).Update("payment_status", "completed")
	log.Printf("Coinbase payment link %s marked as completed via webhook", linkID)
}

// handleOrderFilled 处理订单成交事件
func (c *CoinbaseController) handleOrderFilled(ctx context.Context, orderID string, webhookData map[string]interface{}) {
	db := database.GetDB()

	// 查找对应的订单记录
	var order models.CoinbaseOrder
	if err := db.Where("order_id = ?", orderID).First(&order).Error; err != nil {
		log.Printf("Order not found: %s", orderID)
		return
	}

	// 更新订单状态
	db.Model(&order).Updates(map[string]interface{}{
		"status":  "filled",
		"settled": true,
	})
	log.Printf("Coinbase order %s marked as filled via webhook", orderID)
}

// handleOrderCancelled 处理订单取消事件
func (c *CoinbaseController) handleOrderCancelled(ctx context.Context, orderID string, webhookData map[string]interface{}) {
	db := database.GetDB()

	// 查找对应的订单记录
	var order models.CoinbaseOrder
	if err := db.Where("order_id = ?", orderID).First(&order).Error; err != nil {
		log.Printf("Order not found: %s", orderID)
		return
	}

	// 更新订单状态
	db.Model(&order).Update("status", "cancelled")
	log.Printf("Coinbase order %s marked as cancelled via webhook", orderID)
}
