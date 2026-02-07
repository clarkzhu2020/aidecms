package Controllers

import (
	"aidecms/internal/app/models"
	"aidecms/pkg/kucoin"
	"time"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// KuCoinController KuCoin控制器
type KuCoinController struct{}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	ClientOrderID string `json:"clientOid" binding:"required"`
	Symbol        string `json:"symbol" binding:"required"`
	Side          string `json:"side" binding:"required,oneof=buy sell"`
	Type          string `json:"type" binding:"required,oneof=limit market stop"`
	Price         string `json:"price"`
	Size          string `json:"size" binding:"required"`
	TimeInForce   string `json:"timeInForce" binding:"omitempty,oneof=GTC IOC FOK"`
	StopPrice     string `json:"stopPrice"`
	Remark        string `json:"remark"`
}

// CreateOrder 创建订单
func (c *KuCoinController) CreateOrder(ctx *app.RequestContext) {
	var req CreateOrderRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]interface{}{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// 检查KuCoin客户端是否初始化
	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 创建KuCoin订单请求
	kucoinReq := &kucoin.OrderCreateRequest{
		ClientOid:   req.ClientOrderID,
		Symbol:      req.Symbol,
		Side:        kucoin.OrderSide(req.Side),
		Type:        kucoin.OrderType(req.Type),
		Size:        req.Size,
		TimeInForce: kucoin.OrderTimeInForce(req.TimeInForce),
		Remark:      req.Remark,
	}

	if req.Price != "" {
		kucoinReq.Price = req.Price
	}

	if req.StopPrice != "" {
		kucoinReq.StopPrice = req.StopPrice
	}

	// 调用KuCoin API
	orderResp, err := kucoin.GetClient().CreateOrder(kucoinReq)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to create order: " + err.Error(),
		})
		return
	}

	// 保存订单到数据库
	order := &models.KuCoinOrder{
		OrderID:        orderResp.OrderId,
		ClientOrderID:  orderResp.ClientOid,
		Symbol:         orderResp.Symbol,
		Side:           orderResp.Side,
		Type:           orderResp.Type,
		Price:          orderResp.Price,
		Size:           orderResp.Size,
		DealSize:       orderResp.DealSize,
		DealFunds:      orderResp.DealFunds,
		Fee:            orderResp.Fee,
		FeeCurrency:    orderResp.FeeCurrency,
		TimeInForce:    orderResp.TimeInForce,
		KuCoinCreatedAt: orderResp.CreatedAt,
		Remark:         orderResp.Remark,
		Status:         "open",
	}

	// 保存到数据库
	db := GetDB()
	if err := db.Create(order).Error; err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to save order: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Order created successfully",
		"data":    orderResp,
	})
}

// CancelOrder 取消订单
func (c *KuCoinController) CancelOrder(ctx *app.RequestContext) {
	orderID := ctx.Param("orderId")

	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 调用KuCoin API
	err := kucoin.GetClient().CancelOrderById(orderID)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to cancel order: " + err.Error(),
		})
		return
	}

	// 更新数据库中的订单状态
	db := GetDB()
	err = db.Model(&models.KuCoinOrder{}).
		Where("order_id = ?", orderID).
		Update("status", "canceled").Error

	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to update order status: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Order canceled successfully",
	})
}

// GetOrder 获取订单详情
func (c *KuCoinController) GetOrder(ctx *app.RequestContext) {
	orderID := ctx.Param("orderId")

	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 调用KuCoin API
	orderResp, err := kucoin.GetClient().GetOrderById(orderID)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to get order: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Order retrieved successfully",
		"data":    orderResp,
	})
}

// GetOpenOrders 获取未成交订单
func (c *KuCoinController) GetOpenOrders(ctx *app.RequestContext) {
	symbol := ctx.Query("symbol")

	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 调用KuCoin API
	orders, err := kucoin.GetClient().GetOpenOrders(symbol)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to get open orders: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Open orders retrieved successfully",
		"data":    orders,
	})
}

// GetClosedOrders 获取已成交订单
func (c *KuCoinController) GetClosedOrders(ctx *app.RequestContext) {
	symbol := ctx.Query("symbol")
	status := ctx.Query("status")
	currentPage := ctx.Query("currentPage")
	pageSize := ctx.Query("pageSize")

	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 转换分页参数
	var page, size int
	if currentPage != "" {
		page = 1
	}
	if pageSize != "" {
		size = 20
	}

	// 调用KuCoin API
	orders, err := kucoin.GetClient().GetClosedOrders(symbol, status, page, size)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to get closed orders: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Closed orders retrieved successfully",
		"data":    orders,
	})
}

// GetAccounts 获取账户列表
func (c *KuCoinController) GetAccounts(ctx *app.RequestContext) {
	currency := ctx.Query("currency")
	accountType := ctx.Query("type")

	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 调用KuCoin API
	accounts, err := kucoin.GetClient().GetAccounts(currency, accountType)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to get accounts: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Accounts retrieved successfully",
		"data":    accounts,
	})
}

// GetAccountDetail 获取账户详情
func (c *KuCoinController) GetAccountDetail(ctx *app.RequestContext) {
	accountID := ctx.Param("accountId")

	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 调用KuCoin API
	account, err := kucoin.GetClient().GetAccountDetail(accountID)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to get account detail: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Account detail retrieved successfully",
		"data":    account,
	})
}

// GetSymbols 获取交易对列表
func (c *KuCoinController) GetSymbols(ctx *app.RequestContext) {
	market := ctx.Query("market")

	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 调用KuCoin API
	symbols, err := kucoin.GetClient().GetSymbols(market)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to get symbols: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Symbols retrieved successfully",
		"data":    symbols,
	})
}

// GetTicker 获取Ticker
func (c *KuCoinController) GetTicker(ctx *app.RequestContext) {
	symbol := ctx.Query("symbol")

	if symbol == "" {
		ctx.JSON(400, map[string]interface{}{
			"error": "symbol parameter is required",
		})
		return
	}

	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 调用KuCoin API
	ticker, err := kucoin.GetClient().GetTicker(symbol)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to get ticker: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Ticker retrieved successfully",
		"data":    ticker,
	})
}

// GetOrderBook 获取订单簿
func (c *KuCoinController) GetOrderBook(ctx *app.RequestContext) {
	symbol := ctx.Query("symbol")

	if symbol == "" {
		ctx.JSON(400, map[string]interface{}{
			"error": "symbol parameter is required",
		})
		return
	}

	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 调用KuCoin API
	orderBook, err := kucoin.GetClient().GetOrderBook(symbol)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to get order book: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Order book retrieved successfully",
		"data":    orderBook,
	})
}

// GetMarketTrades 获取市场成交记录
func (c *KuCoinController) GetMarketTrades(ctx *app.RequestContext) {
	symbol := ctx.Query("symbol")

	if symbol == "" {
		ctx.JSON(400, map[string]interface{}{
			"error": "symbol parameter is required",
		})
		return
	}

	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 调用KuCoin API
	trades, err := kucoin.GetClient().GetMarketTrades(symbol)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to get market trades: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Market trades retrieved successfully",
		"data":    trades,
	})
}

// GetKlines 获取K线数据
func (c *KuCoinController) GetKlines(ctx *app.RequestContext) {
	symbol := ctx.Query("symbol")
	klineType := ctx.DefaultQuery("type", "1hour")

	if symbol == "" {
		ctx.JSON(400, map[string]interface{}{
			"error": "symbol parameter is required",
		})
		return
	}

	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 调用KuCoin API
	klines, err := kucoin.GetClient().GetKlines(symbol, klineType, 0, 0)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to get klines: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Klines retrieved successfully",
		"data":    klines,
	})
}

// Get24HStats 获取24小时统计数据
func (c *KuCoinController) Get24HStats(ctx *app.RequestContext) {
	symbol := ctx.Query("symbol")

	if symbol == "" {
		ctx.JSON(400, map[string]interface{}{
			"error": "symbol parameter is required",
		})
		return
	}

	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 调用KuCoin API
	stats, err := kucoin.GetClient().Get24HStats(symbol)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to get 24h stats: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "24h stats retrieved successfully",
		"data":    stats,
	})
}

// GetServerTime 获取服务器时间
func (c *KuCoinController) GetServerTime(ctx *app.RequestContext) {
	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 调用KuCoin API
	timestamp, err := kucoin.GetClient().GetServerTime()
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to get server time: " + err.Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Server time retrieved successfully",
		"data": map[string]interface{}{
			"timestamp": timestamp,
			"datetime":  time.Unix(0, timestamp*1e6).Format(time.RFC3339),
		},
	})
}

// SyncAccounts 同步账户信息到数据库
func (c *KuCoinController) SyncAccounts(ctx *app.RequestContext) {
	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 调用KuCoin API获取账户列表
	accounts, err := kucoin.GetClient().GetAccounts("", "")
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to get accounts: " + err.Error(),
		})
		return
	}

	db := GetDB()
	syncedCount := 0

	// 同步每个账户到数据库
	for _, account := range accounts {
		var dbAccount models.KuCoinAccount
		err := db.Where("account_id = ?", account.Id).First(&dbAccount).Error

		if err == gorm.ErrRecordNotFound {
			// 创建新记录
			dbAccount = models.KuCoinAccount{
				AccountID:    account.Id,
				Currency:     account.Currency,
				Type:         account.Type,
				Balance:      account.Balance,
				Available:    account.Available,
				Holds:        account.Holds,
				LastSyncedAt: time.Now(),
			}
			if err := db.Create(&dbAccount).Error; err != nil {
				continue
			}
			syncedCount++
		} else if err == nil {
			// 更新现有记录
			dbAccount.Balance = account.Balance
			dbAccount.Available = account.Available
			dbAccount.Holds = account.Holds
			dbAccount.LastSyncedAt = time.Now()
			if err := db.Save(&dbAccount).Error; err != nil {
				continue
			}
			syncedCount++
		}
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Accounts synced successfully",
		"data": map[string]interface{}{
			"synced_count": syncedCount,
			"total_count":  len(accounts),
		},
	})
}

// CreateBalanceSnapshot 创建余额快照
func (c *KuCoinController) CreateBalanceSnapshot(ctx *app.RequestContext) {
	if !kucoin.IsInitialized() {
		ctx.JSON(500, map[string]interface{}{
			"error": "KuCoin client not initialized",
		})
		return
	}

	// 调用KuCoin API获取账户列表
	accounts, err := kucoin.GetClient().GetAccounts("", "")
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to get accounts: " + err.Error(),
		})
		return
	}

	db := GetDB()
	snapshotID := uuid.New().String()
	snapshotTime := time.Now()
	snapshotCount := 0

	// 创建余额快照
	for _, account := range accounts {
		snapshot := &models.KuCoinBalanceSnapshot{
			SnapshotID: snapshotID,
			AccountID:  account.Id,
			Currency:   account.Currency,
			Balance:    account.Balance,
			Available:  account.Available,
			Holds:      account.Holds,
			SnapshotAt: snapshotTime,
		}

		if err := db.Create(snapshot).Error; err != nil {
			continue
		}
		snapshotCount++
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Balance snapshot created successfully",
		"data": map[string]interface{}{
			"snapshot_id":  snapshotID,
			"snapshot_at":  snapshotTime,
			"snapshot_count": snapshotCount,
		},
	})
}
