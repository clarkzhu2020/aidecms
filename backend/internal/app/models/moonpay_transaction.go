package models

import (
	"time"
)

// MoonPayTransaction MoonPay交易模型
type MoonPayTransaction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// MoonPay交易ID
	TransactionID string `json:"transaction_id" gorm:"size:100;uniqueIndex;not null;comment:MoonPay交易ID"`

	// 外部订单ID
	ExternalID    string `json:"external_id" gorm:"size:100;index;comment:外部订单ID"`

	// 交易类型
	TransactionType string `json:"transaction_type" gorm:"size:20;comment:交易类型 buy/sell/swap"`

	// 支付金额
	BaseCurrencyAmount float64 `json:"base_currency_amount" gorm:"type:decimal(10,2);not null;comment:法币金额"`
	BaseCurrencyCode  string  `json:"base_currency_code" gorm:"size:10;not null;comment:法币类型"`

	// 加密货币
	QuoteCurrencyAmount float64 `json:"quote_currency_amount" gorm:"type:decimal(20,8);comment:加密货币数量"`
	QuoteCurrencyCode  string  `json:"quote_currency_code" gorm:"size:20;comment:加密货币类型"`
	CurrencyCode        string  `json:"currency_code" gorm:"size:20;not null;comment:目标货币类型"`

	// 钱包地址
	WalletAddress string `json:"wallet_address" gorm:"size:255;not null;comment:钱包地址"`

	// 交易状态
	Status string `json:"status" gorm:"size:20;index;not null;comment:交易状态 pending/waiting_payment/pending_approval/completed/failed"`

	// KYC状态
	KycStatus string `json:"kyc_status" gorm:"size:20;default:'not_started';comment:KYC状态"`

	// 支付方式
	PaymentMethodType string `json:"payment_method_type" gorm:"size:20;comment:支付方式类型"`

	// 费用
	FeeAmount        float64 `json:"fee_amount" gorm:"type:decimal(10,2);comment:基础费用"`
	ExtraFeeAmount   float64 `json:"extra_fee_amount" gorm:"type:decimal(10,2);comment:额外费用"`
	NetworkFeeAmount float64 `json:"network_fee_amount" gorm:"type:decimal(20,8);comment:网络费用"`

	// 客户信息
	CustomerID    string `json:"customer_id" gorm:"size:100;index;comment:MoonPay客户ID"`
	CustomerEmail string `json:"customer_email" gorm:"size:255;comment:客户邮箱"`
	FirstName     string `json:"first_name" gorm:"size:100;comment:名"`
	LastName      string `json:"last_name" gorm:"size:100;comment:姓"`

	// 回调URL
	RedirectURL string `json:"redirect_url" gorm:"size:500;comment:重定向URL"`
	WidgetURL   string `json:"widget_url" gorm:"size:500;comment:Widget URL"`

	// 元数据
	Metadata string `json:"metadata" gorm:"type:json;comment:额外元数据"`
	RawData  string `json:"raw_data" gorm:"type:json;comment:原始JSON数据"`

	// 用户关联 (可选)
	UserID *uint `json:"user_id" gorm:"index;comment:关联用户ID"`
	User   *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// MoonPayWebhook MoonPay Webhook记录模型
type MoonPayWebhook struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Webhook信息
	WebhookID string `json:"webhook_id" gorm:"size:100;uniqueIndex;comment:MoonPay事件ID"`
	EventType string `json:"event_type" gorm:"size:100;index;comment:事件类型"`
	Summary   string `json:"summary" gorm:"type:text;comment:事件摘要"`

	// 关联交易
	TransactionID   string `json:"transaction_id" gorm:"size:100;index;comment:交易ID"`
	Transaction     *MoonPayTransaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID;references:TransactionID"`

	// 状态
	Status string `json:"status" gorm:"size:20;default:'pending';index;comment:处理状态"`

	// 原始数据
	RawData    string `json:"raw_data" gorm:"type:json;comment:原始JSON数据"`
	RawHeaders string `json:"raw_headers" gorm:"type:json;comment:原始请求头"`
}

// TableName 指定表名
func (MoonPayTransaction) TableName() string {
	return "moonpay_transactions"
}

func (MoonPayWebhook) TableName() string {
	return "moonpay_webhooks"
}

// IsCompleted 检查交易是否已完成
func (t *MoonPayTransaction) IsCompleted() bool {
	return t.Status == "completed"
}

// IsPending 检查交易是否待处理
func (t *MoonPayTransaction) IsPending() bool {
	return t.Status == "pending" || t.Status == "waiting_payment" || t.Status == "pending_approval"
}

// IsFailed 检查交易是否失败
func (t *MoonPayTransaction) IsFailed() bool {
	return t.Status == "failed"
}

// IsBuyTransaction 检查是否为购买交易
func (t *MoonPayTransaction) IsBuyTransaction() bool {
	return t.TransactionType == "buy"
}

// IsSellTransaction 检查是否为出售交易
func (t *MoonPayTransaction) IsSellTransaction() bool {
	return t.TransactionType == "sell"
}

// GetStatusText 获取状态文本描述
func (t *MoonPayTransaction) GetStatusText() string {
	statusMap := map[string]string{
		"pending":          "待处理",
		"waiting_payment":   "等待支付",
		"pending_approval": "待审批",
		"completed":         "已完成",
		"failed":           "失败",
		"cancelled":        "已取消",
	}
	if text, ok := statusMap[t.Status]; ok {
		return text
	}
	return t.Status
}
