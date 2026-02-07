package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"github.com/clarkzhu2020/aidecms/pkg/database"
	"github.com/clarkzhu2020/aidecms/pkg/moonpay"
	"github.com/clarkzhu2020/aidecms/pkg/response"
	"github.com/cloudwego/hertz/pkg/app"
)

// MoonPayController MoonPay支付控制器
type MoonPayController struct{}

// NewMoonPayController 创建MoonPay控制器
func NewMoonPayController() *MoonPayController {
	return &MoonPayController{}
}

// CreateTransactionRequest 创建交易请求
type CreateTransactionRequest struct {
	BaseCurrencyAmount float64 `json:"baseCurrencyAmount" validate:"required,gt=0"`
	BaseCurrencyCode  string  `json:"baseCurrencyCode" validate:"required"`
	CurrencyCode      string  `json:"currencyCode" validate:"required"`
	WalletAddress     string  `json:"walletAddress" validate:"required"`
	ExternalID        string  `json:"externalId" validate:"required"`
	RedirectURL       string  `json:"redirectUrl"`
	LockAmount        bool    `json:"lockAmount"`
	Email             string  `json:"email"`
	FirstName         string  `json:"firstName"`
	LastName          string  `json:"lastName"`
}

// GetQuoteRequest 获取报价请求
type GetQuoteRequest struct {
	BaseCurrencyAmount float64 `json:"baseCurrencyAmount" validate:"required,gt=0"`
	BaseCurrencyCode  string  `json:"baseCurrencyCode" validate:"required"`
	CurrencyCode      string  `json:"currencyCode" validate:"required"`
}

// CreateTransaction 创建MoonPay交易
// @Summary      创建MoonPay交易
// @Description  创建一个新的MoonPay加密货币购买交易
// @Tags         MoonPay
// @Accept       json
// @Produce      json
// @Param        transaction body CreateTransactionRequest true "交易信息"
// @Success      201 {object} response.Response{data=models.MoonPayTransaction}
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /moonpay/transactions [post]
func (c *MoonPayController) CreateTransaction(ctx context.Context, hCtx *app.RequestContext) {
	var req CreateTransactionRequest
	if err := hCtx.BindJSON(&req); err != nil {
		response.BadRequest(hCtx, "Invalid request data")
		return
	}

	// 检查MoonPay客户端是否初始化
	if !moonpay.IsInitialized() {
		response.ServerError(hCtx, "MoonPay service not configured")
		return
	}

	// 获取当前用户ID（从JWT中间件设置）
	userID, _ := hCtx.Get("user_id")
	var userIDPtr *uint
	if uid, ok := userID.(uint); ok {
		userIDPtr = &uid
	}

	// 构建重定向URL
	if req.RedirectURL == "" {
		baseURL := "http://localhost:8888/api"
		req.RedirectURL = fmt.Sprintf("%s/moonpay/success", baseURL)
	}

	// 生成外部ID（如果没有提供）
	externalID := req.ExternalID
	if externalID == "" {
		externalID = fmt.Sprintf("AIDECMS-%d", userID)
	}

	// 获取MoonPay客户端配置
	client := moonpay.GetClient()
	config := client.GetConfig()

	// 调用MoonPay服务创建交易
	transactionService := moonpay.NewTransactionService()
	transactionReq := &moonpay.CreateTransactionRequest{
		BaseCurrencyAmount: req.BaseCurrencyAmount,
		BaseCurrencyCode:  req.BaseCurrencyCode,
		CurrencyCode:      req.CurrencyCode,
		WalletAddress:     req.WalletAddress,
		ExternalID:        externalID,
		RedirectURL:       req.RedirectURL,
		LockAmount:        req.LockAmount,
		Email:             req.Email,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
	}

	transaction, err := transactionService.CreateTransaction(ctx, transactionReq)
	if err != nil {
		log.Printf("Failed to create MoonPay transaction: %v", err)
		response.ServerError(hCtx, "Failed to create transaction")
		return
	}

	// 生成Widget URL（前端集成）
	widgetURL := transactionService.GenerateWidgetURL(transactionReq)

	// 保存交易记录到数据库
	db := database.GetDB()
	moonpayTransaction := &models.MoonPayTransaction{
		TransactionID:      transaction.ID,
		ExternalID:         externalID,
		TransactionType:    "buy",
		BaseCurrencyAmount: req.BaseCurrencyAmount,
		BaseCurrencyCode:   req.BaseCurrencyCode,
		CurrencyCode:       req.CurrencyCode,
		WalletAddress:      req.WalletAddress,
		Status:            transaction.Status,
		RedirectURL:       req.RedirectURL,
		WidgetURL:         widgetURL,
		UserID:            userIDPtr,
	}

	// 解析客户信息
	if transaction.Customer != nil {
		moonpayTransaction.CustomerID = transaction.Customer.ID
		moonpayTransaction.CustomerEmail = transaction.Customer.Email
		moonpayTransaction.FirstName = transaction.Customer.FirstName
		moonpayTransaction.LastName = transaction.Customer.LastName
	}

	// 序列化原始数据
	if rawData, err := json.Marshal(transaction); err == nil {
		moonpayTransaction.RawData = string(rawData)
	}

	if err := db.Create(moonpayTransaction).Error; err != nil {
		log.Printf("Failed to save transaction record: %v", err)
		response.ServerError(hCtx, "Failed to save transaction record")
		return
	}

	// 返回交易信息
	response.Created(hCtx, map[string]interface{}{
		"transaction_id":  moonpayTransaction.ID,
		"moonpay_id":      moonpayTransaction.TransactionID,
		"widget_url":      moonpayTransaction.WidgetURL,
		"status":          moonpayTransaction.Status,
		"amount":          moonpayTransaction.BaseCurrencyAmount,
		"currency":        moonpayTransaction.BaseCurrencyCode,
		"currency_code":   moonpayTransaction.CurrencyCode,
		"wallet_address":  moonpayTransaction.WalletAddress,
	}, "Transaction created successfully")
}

// GetTransaction 获取交易详情
// @Summary      获取交易详情
// @Description  根据ID获取MoonPay交易详细信息
// @Tags         MoonPay
// @Accept       json
// @Produce      json
// @Param        id path int true "交易ID"
// @Success      200 {object} response.Response{data=models.MoonPayTransaction}
// @Failure      404 {object} response.Response
// @Security     BearerAuth
// @Router       /moonpay/transactions/{id} [get]
func (c *MoonPayController) GetTransaction(ctx context.Context, hCtx *app.RequestContext) {
	id := hCtx.Param("id")
	transactionID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.BadRequest(hCtx, "Invalid transaction ID")
		return
	}

	db := database.GetDB()
	var transaction models.MoonPayTransaction
	if err := db.First(&transaction, transactionID).Error; err != nil {
		response.NotFound(hCtx, "Transaction not found")
		return
	}

	// 可选：从MoonPay API获取最新状态
	if moonpay.IsInitialized() {
		transactionService := moonpay.NewTransactionService()
		if mpTransaction, err := transactionService.GetTransaction(ctx, transaction.TransactionID); err == nil {
			transaction.Status = mpTransaction.Status
			transaction.QuoteCurrencyAmount = mpTransaction.QuoteCurrencyAmount
			transaction.QuoteCurrencyCode = mpTransaction.QuoteCurrencyCode
			transaction.FeeAmount = mpTransaction.FeeAmount
			transaction.ExtraFeeAmount = mpTransaction.ExtraFeeAmount
			transaction.NetworkFeeAmount = mpTransaction.NetworkFeeAmount
			if rawData, err := json.Marshal(mpTransaction); err == nil {
				transaction.RawData = string(rawData)
			}
			db.Save(&transaction)
		}
	}

	response.Success(hCtx, transaction, "Transaction fetched successfully")
}

// ListTransactions 获取交易列表
// @Summary      获取交易列表
// @Description  获取MoonPay交易列表，支持分页和筛选
// @Tags         MoonPay
// @Accept       json
// @Produce      json
// @Param        status query string false "交易状态" Enums(pending, waiting_payment, pending_approval, completed, failed)
// @Param        page query int false "页码" default(1)
// @Param        limit query int false "每页数量" default(20)
// @Success      200 {object} response.Response{data=[]models.MoonPayTransaction}
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /moonpay/transactions [get]
func (c *MoonPayController) ListTransactions(ctx context.Context, hCtx *app.RequestContext) {
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

	var transactions []models.MoonPayTransaction
	query := db.Model(&models.MoonPayTransaction{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 分页
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit).Order("created_at DESC")

	if err := query.Find(&transactions).Error; err != nil {
		response.ServerError(hCtx, "Failed to fetch transactions")
		return
	}

	response.Success(hCtx, map[string]interface{}{
		"data":   transactions,
		"page":   page,
		"limit":  limit,
		"total":  len(transactions),
	}, "Transactions fetched successfully")
}

// GetQuote 获取购买报价
// @Summary      获取购买报价
// @Description  获取加密货币购买的实时报价
// @Tags         MoonPay
// @Accept       json
// @Produce      json
// @Param        quote body GetQuoteRequest true "报价请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /moonpay/quote [post]
func (c *MoonPayController) GetQuote(ctx context.Context, hCtx *app.RequestContext) {
	var req GetQuoteRequest
	if err := hCtx.BindJSON(&req); err != nil {
		response.BadRequest(hCtx, "Invalid request data")
		return
	}

	// 检查MoonPay客户端是否初始化
	if !moonpay.IsInitialized() {
		response.ServerError(hCtx, "MoonPay service not configured")
		return
	}

	// 调用MoonPay服务获取报价
	transactionService := moonpay.NewTransactionService()
	quote, err := transactionService.GetBuyQuote(ctx, req.BaseCurrencyAmount, req.BaseCurrencyCode, req.CurrencyCode)
	if err != nil {
		log.Printf("Failed to get quote: %v", err)
		response.ServerError(hCtx, "Failed to get quote")
		return
	}

	response.Success(hCtx, quote, "Quote fetched successfully")
}

// GenerateWidgetURL 生成Widget URL
// @Summary      生成Widget URL
// @Description  生成MoonPay Widget URL用于前端集成
// @Tags         MoonPay
// @Accept       json
// @Produce      json
// @Param        transaction body CreateTransactionRequest true "交易参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Security     BearerAuth
// @Router       /moonpay/widget-url [post]
func (c *MoonPayController) GenerateWidgetURL(ctx context.Context, hCtx *app.RequestContext) {
	var req CreateTransactionRequest
	if err := hCtx.BindJSON(&req); err != nil {
		response.BadRequest(hCtx, "Invalid request data")
		return
	}

	// 检查MoonPay客户端是否初始化
	if !moonpay.IsInitialized() {
		response.ServerError(hCtx, "MoonPay service not configured")
		return
	}

	// 构建重定向URL
	if req.RedirectURL == "" {
		baseURL := "http://localhost:8888/api"
		req.RedirectURL = fmt.Sprintf("%s/moonpay/success", baseURL)
	}

	// 调用MoonPay服务生成Widget URL
	transactionService := moonpay.NewTransactionService()
	transactionReq := &moonpay.CreateTransactionRequest{
		BaseCurrencyAmount: req.BaseCurrencyAmount,
		BaseCurrencyCode:  req.BaseCurrencyCode,
		CurrencyCode:      req.CurrencyCode,
		WalletAddress:     req.WalletAddress,
		ExternalID:        req.ExternalID,
		RedirectURL:       req.RedirectURL,
		LockAmount:        req.LockAmount,
		Email:             req.Email,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
	}

	widgetURL := transactionService.GenerateWidgetURL(transactionReq)

	response.Success(hCtx, map[string]interface{}{
		"widget_url": widgetURL,
	}, "Widget URL generated successfully")
}

// HandleWebhook 处理MoonPay Webhook通知
// @Summary      处理Webhook
// @Description  接收并处理MoonPay的Webhook通知
// @Tags         MoonPay
// @Accept       json
// @Produce      json
// @Param        webhook body map[string]interface{} true "Webhook数据"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /moonpay/webhook [post]
func (c *MoonPayController) HandleWebhook(ctx context.Context, hCtx *app.RequestContext) {
	// 读取请求体
	body := hCtx.Request.Body()
	var webhookData map[string]interface{}
	if err := json.Unmarshal(body, &webhookData); err != nil {
		log.Printf("Failed to parse webhook data: %v", err)
		response.BadRequest(hCtx, "Invalid webhook data")
		return
	}

	// 验证签名（如果配置了Webhook密钥）
	// signature := hCtx.GetHeader("X-MoonPay-Signature")

	// 提取Webhook信息
	webhookID := ""
	eventType := ""
	transactionID := ""
	summary := ""

	if id, ok := webhookData["id"].(string); ok {
		webhookID = id
	}
	if et, ok := webhookData["type"].(string); ok {
		eventType = et
	}
	if tid, ok := webhookData["transactionId"].(string); ok {
		transactionID = tid
	}
	if s, ok := webhookData["summary"].(string); ok {
		summary = s
	}

	// 保存Webhook记录
	db := database.GetDB()
	webhook := &models.MoonPayWebhook{
		WebhookID:     webhookID,
		EventType:     eventType,
		TransactionID: transactionID,
		Summary:       summary,
		Status:        "processed",
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
	case "transaction.completed":
		c.handleTransactionCompleted(ctx, transactionID, webhookData)
	case "transaction.failed":
		c.handleTransactionFailed(ctx, transactionID, webhookData)
	case "transaction.created":
		c.handleTransactionCreated(ctx, transactionID, webhookData)
	}

	response.Success(hCtx, map[string]interface{}{
		"status": "processed",
	}, "Webhook processed successfully")
}

// handleTransactionCompleted 处理交易完成事件
func (c *MoonPayController) handleTransactionCompleted(ctx context.Context, transactionID string, webhookData map[string]interface{}) {
	db := database.GetDB()

	// 查找对应的交易记录
	var transaction models.MoonPayTransaction
	if err := db.Where("transaction_id = ?", transactionID).First(&transaction).Error; err != nil {
		log.Printf("Transaction not found: %s", transactionID)
		return
	}

	// 更新交易状态
	db.Model(&transaction).Update("status", "completed")
	log.Printf("MoonPay transaction %s marked as completed via webhook", transactionID)
}

// handleTransactionFailed 处理交易失败事件
func (c *MoonPayController) handleTransactionFailed(ctx context.Context, transactionID string, webhookData map[string]interface{}) {
	db := database.GetDB()

	// 查找对应的交易记录
	var transaction models.MoonPayTransaction
	if err := db.Where("transaction_id = ?", transactionID).First(&transaction).Error; err != nil {
		log.Printf("Transaction not found: %s", transactionID)
		return
	}

	// 更新交易状态
	db.Model(&transaction).Update("status", "failed")
	log.Printf("MoonPay transaction %s marked as failed via webhook", transactionID)
}

// handleTransactionCreated 处理交易创建事件
func (c *MoonPayController) handleTransactionCreated(ctx context.Context, transactionID string, webhookData map[string]interface{}) {
	db := database.GetDB()

	// 查找对应的交易记录
	var transaction models.MoonPayTransaction
	if err := db.Where("transaction_id = ?", transactionID).First(&transaction).Error; err != nil {
		// 如果记录不存在，创建新记录
		// 这里可以根据实际需求从webhookData中提取信息创建记录
		log.Printf("Transaction created webhook received for new transaction: %s", transactionID)
		return
	}

	// 更新交易状态为pending
	if transaction.Status == "" {
		db.Model(&transaction).Update("status", "pending")
	}
}
