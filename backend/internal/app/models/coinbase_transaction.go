package models

import (
	"time"
)

// CoinbasePaymentLink Coinbase支付链接模型
type CoinbasePaymentLink struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Coinbase链接ID
	LinkID      string `json:"link_id" gorm:"size:100;uniqueIndex;not null;comment:Coinbase链接ID"`

	// 外部订单ID
	ExternalID  string `json:"external_id" gorm:"size:100;index;comment:外部订单ID"`

	// 支付金额
	Amount      string `json:"amount" gorm:"size:50;not null;comment:支付金额"`
	Currency    string `json:"currency" gorm:"size:10;not null;comment:货币类型"`

	// 支付信息
	Title       string `json:"title" gorm:"size:255;comment:标题"`
	Description string `json:"description" gorm:"type:text;comment:描述"`

	// 链接状态
	Status      string `json:"status" gorm:"size:20;index;comment:链接状态 active/expired/completed"`

	// 支付状态
	PaymentStatus string `json:"payment_status" gorm:"size:20;default:'pending';index;comment:支付状态 pending/completed/failed/cancelled"`

	// URL信息
	PaymentURL  string `json:"payment_url" gorm:"size:500;comment:支付URL"`
	RedirectURL string `json:"redirect_url" gorm:"size:500;comment:成功回调URL"`
	CancelURL   string `json:"cancel_url" gorm:"size:500;comment:取消回调URL"`

	// 过期时间
	ExpiresAt   *time.Time `json:"expires_at" gorm:"comment:过期时间"`

	// 客户信息
	Name        string `json:"name" gorm:"size:255;comment:客户名称"`
	Email       string `json:"email" gorm:"size:255;comment:客户邮箱"`

	// 元数据
	Metadata    string `json:"metadata" gorm:"type:json;comment:额外元数据"`
	RawData     string `json:"raw_data" gorm:"type:json;comment:原始JSON数据"`

	// 用户关联
	UserID      *uint `json:"user_id" gorm:"index;comment:关联用户ID"`
	User        *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// CoinbaseOrder Coinbase交易订单模型
type CoinbaseOrder struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Coinbase订单ID
	OrderID     string `json:"order_id" gorm:"size:100;uniqueIndex;not null;comment:Coinbase订单ID"`

	// 客户端订单ID
	ClientOrderID string `json:"client_order_id" gorm:"size:100;index;comment:客户端订单ID"`

	// 外部订单ID
	ExternalID  string `json:"external_id" gorm:"size:100;index;comment:外部订单ID"`

	// 产品信息
	ProductID   string `json:"product_id" gorm:"size:20;index;not null;comment:产品ID 如BTC-USD"`

	// 订单信息
	Side        string `json:"side" gorm:"size:10;not null;comment:买卖方向 buy/sell"`
	OrderType   string `json:"order_type" gorm:"size:20;not null;comment:订单类型 market/limit/stop"`

	// 数量和价格
	Size        string `json:"size" gorm:"size:50;comment:数量"`
	Funds       string `json:"funds" gorm:"size:50;comment:金额"`
	LimitPrice  string `json:"limit_price" gorm:"size:50;comment:限价"`
	StopPrice   string `json:"stop_price" gorm:"size:50;comment:止损价"`

	// 执行信息
	FilledSize         string  `json:"filled_size" gorm:"size:50;default:'0';comment:已成交数量"`
	AverageFillPrice   string  `json:"average_fill_price" gorm:"size:50;comment:平均成交价"`
	FillFees          string  `json:"fill_fees" gorm:"size:50;default:'0';comment:手续费"`

	// 订单状态
	Status      string `json:"status" gorm:"size:20;index;not null;comment:订单状态 open/filled/rejected/cancelled"`
	Settled     bool   `json:"settled" gorm:"default:false;comment:是否已结算"`

	// 时间信息
	ExpiredAt   *time.Time `json:"expired_at" gorm:"comment:过期时间"`

	// 配置信息
	TimeInForce string `json:"time_in_force" gorm:"size:10;comment:有效期类型 GTC/IOC/FOK/GTD"`
	PostOnly    bool   `json:"post_only" gorm:"default:false;comment:仅作为maker"`

	// 元数据
	Metadata    string `json:"metadata" gorm:"type:json;comment:额外元数据"`
	RawData     string `json:"raw_data" gorm:"type:json;comment:原始JSON数据"`

	// 用户关联
	UserID      *uint `json:"user_id" gorm:"index;comment:关联用户ID"`
	User        *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// CoinbaseWebhook Coinbase Webhook记录模型
type CoinbaseWebhook struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Webhook信息
	WebhookID   string `json:"webhook_id" gorm:"size:100;uniqueIndex;comment:Webhook事件ID"`
	EventType   string `json:"event_type" gorm:"size:100;index;comment:事件类型"`
	Summary     string `json:"summary" gorm:"type:text;comment:事件摘要"`

	// 关联资源
	LinkID      string `json:"link_id" gorm:"size:100;index;comment:支付链接ID"`
	Link        *CoinbasePaymentLink `json:"link,omitempty" gorm:"foreignKey:LinkID;references:LinkID"`

	OrderID     string `json:"order_id" gorm:"size:100;index;comment:交易订单ID"`
	Order       *CoinbaseOrder `json:"order,omitempty" gorm:"foreignKey:OrderID;references:OrderID"`

	// 状态
	Status      string `json:"status" gorm:"size:20;default:'pending';index;comment:处理状态"`

	// 原始数据
	RawData     string `json:"raw_data" gorm:"type:json;comment:原始JSON数据"`
	RawHeaders  string `json:"raw_headers" gorm:"type:json;comment:原始请求头"`
}

// TableName 指定表名
func (CoinbasePaymentLink) TableName() string {
	return "coinbase_payment_links"
}

func (CoinbaseOrder) TableName() string {
	return "coinbase_orders"
}

func (CoinbaseWebhook) TableName() string {
	return "coinbase_webhooks"
}

// IsCompleted 检查支付链接是否已完成
func (p *CoinbasePaymentLink) IsCompleted() bool {
	return p.PaymentStatus == "completed"
}

// IsPending 检查支付链接是否待处理
func (p *CoinbasePaymentLink) IsPending() bool {
	return p.PaymentStatus == "pending"
}

// IsExpired 检查支付链接是否已过期
func (p *CoinbasePaymentLink) IsExpired() bool {
	return p.Status == "expired"
}

// IsFilled 检查订单是否已成交
func (o *CoinbaseOrder) IsFilled() bool {
	return o.Status == "filled"
}

// IsOpen 检查订单是否开启
func (o *CoinbaseOrder) IsOpen() bool {
	return o.Status == "open"
}

// IsRejected 检查订单是否被拒绝
func (o *CoinbaseOrder) IsRejected() bool {
	return o.Status == "rejected"
}

// IsCancelled 检查订单是否已取消
func (o *CoinbaseOrder) IsCancelled() bool {
	return o.Status == "cancelled"
}

// IsBuyOrder 检查是否为买单
func (o *CoinbaseOrder) IsBuyOrder() bool {
	return o.Side == "buy"
}

// IsSellOrder 检查是否为卖单
func (o *CoinbaseOrder) IsSellOrder() bool {
	return o.Side == "sell"
}
