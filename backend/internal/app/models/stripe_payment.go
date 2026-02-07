package models

import (
	"time"
)

// StripePayment Stripe支付模型
type StripePayment struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// 订单信息
	OrderID       string    `json:"order_id" gorm:"size:100;uniqueIndex;not null;comment:内部订单ID"`
	PaymentIntentID string   `json:"payment_intent_id" gorm:"size:100;uniqueIndex;comment:Stripe PaymentIntent ID"`
	ChargeID      string    `json:"charge_id" gorm:"size:100;index;comment:Stripe Charge ID"`

	// 支付金额
	Amount        int64     `json:"amount" gorm:"not null;comment:支付金额(最小单位)"`
	Currency      string    `json:"currency" gorm:"size:10;not null;default:'usd';comment:货币类型"`

	// 支付状态
	Status        string    `json:"status" gorm:"size:20;index;comment:支付状态"`
	PaymentStatus string    `json:"payment_status" gorm:"size:50;comment:Stripe支付状态"`

	// 客户信息
	CustomerID    string    `json:"customer_id" gorm:"size:100;index;comment:Stripe客户ID"`
	CustomerEmail string    `json:"customer_email" gorm:"size:255;comment:客户邮箱"`

	// 支付信息
	Description   string    `json:"description" gorm:"type:text;comment:支付描述"`
	ReceiptURL    string    `json:"receipt_url" gorm:"size:500;comment:收据URL"`

	// 元数据
	Metadata      string    `json:"metadata" gorm:"type:json;comment:元数据"`

	// 用户关联 (可选)
	UserID        *uint     `json:"user_id" gorm:"index;comment:关联用户ID"`
	User          *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// StripeRefund Stripe退款记录模型
type StripeRefund struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// 关联支付
	PaymentID       uint      `json:"payment_id" gorm:"not null;index;comment:支付ID"`
	Payment         *StripePayment `json:"payment,omitempty" gorm:"foreignKey:PaymentID"`

	// 退款信息
	RefundID        string    `json:"refund_id" gorm:"size:100;uniqueIndex;comment:Stripe退款ID"`
	ChargeID        string    `json:"charge_id" gorm:"size:100;index;comment:关联的Charge ID"`
	Amount          int64     `json:"amount" gorm:"not null;comment:退款金额(最小单位)"`
	Currency        string    `json:"currency" gorm:"size:10;default:'usd';comment:货币类型"`
	Status          string    `json:"status" gorm:"size:20;default:'pending';comment:退款状态"`
	Reason          string    `json:"reason" gorm:"size:50;comment:退款原因"`

	// 退款详情
	Description     string    `json:"description" gorm:"type:text;comment:退款描述"`
}

// StripeWebhook Webhook记录模型
type StripeWebhook struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Webhook信息
	StripeID        string    `json:"stripe_id" gorm:"size:100;uniqueIndex;comment:Stripe事件ID"`
	Type            string    `json:"type" gorm:"size:100;index;comment:事件类型"`
	Data            string    `json:"data" gorm:"type:json;comment:事件数据"`

	// 状态
	Status          string    `json:"status" gorm:"size:20;default:'pending';index;comment:处理状态"`
	Error           string    `json:"error" gorm:"type:text;comment:错误信息"`
}

// TableName 指定表名
func (StripePayment) TableName() string {
	return "stripe_payments"
}

func (StripeRefund) TableName() string {
	return "stripe_refunds"
}

func (StripeWebhook) TableName() string {
	return "stripe_webhooks"
}

// IsSucceeded 检查支付是否成功
func (p *StripePayment) IsSucceeded() bool {
	return p.Status == "succeeded"
}

// IsPending 检查支付是否待处理
func (p *StripePayment) IsPending() bool {
	return p.Status == "pending" || p.Status == "processing"
}

// IsFailed 检查支付是否失败
func (p *StripePayment) IsFailed() bool {
	return p.Status == "failed"
}

// IsCanceled 检查支付是否已取消
func (p *StripePayment) IsCanceled() bool {
	return p.Status == "canceled"
}

// GetAmountInDollars 获取美元金额
func (p *StripePayment) GetAmountInDollars() float64 {
	return float64(p.Amount) / 100.0
}
