package models

import (
	"time"
)

// Payment 支付模型
type Payment struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// 订单信息
	OrderID       string    `json:"order_id" gorm:"size:100;uniqueIndex;not null;comment:内部订单ID"`
	PayPalOrderID string    `json:"paypal_order_id" gorm:"size:100;index;comment:PayPal订单ID"`

	// 支付金额
	Amount        float64   `json:"amount" gorm:"type:decimal(10,2);not null;comment:支付金额"`
	Currency      string    `json:"currency" gorm:"size:10;default:'USD';comment:货币类型"`

	// 支付状态
	Status        string    `json:"status" gorm:"size:20;index;comment:支付状态 pending/paid/failed/cancelled/refunded"`
	PaymentStatus string    `json:"payment_status" gorm:"size:50;default:'CREATED';comment:PayPal支付状态"`

	// 支付详情
	PayerID       string    `json:"payer_id" gorm:"size:100;comment:支付者ID"`
	PayerEmail    string    `json:"payer_email" gorm:"size:255;comment:支付者邮箱"`
	CaptureID     string    `json:"capture_id" gorm:"size:100;index;comment:捕获支付ID"`

	// 订单信息
	Description   string    `json:"description" gorm:"type:text;comment:订单描述"`
	ReferenceID   string    `json:"reference_id" gorm:"size:100;index;comment:参考ID"`

	// 回调信息
	ApprovalURL   string    `json:"approval_url" gorm:"size:500;comment:支付审批URL"`
	ReturnURL     string    `json:"return_url" gorm:"size:500;comment:成功回调URL"`
	CancelURL     string    `json:"cancel_url" gorm:"size:500;comment:取消回调URL"`

	// 元数据
	Metadata      string    `json:"metadata" gorm:"type:json;comment:额外元数据"`

	// 用户关联 (可选)
	UserID        *uint     `json:"user_id" gorm:"index;comment:关联用户ID"`
	User          *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// PaymentRefund 退款记录模型
type PaymentRefund struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// 关联支付
	PaymentID       uint      `json:"payment_id" gorm:"not null;index;comment:支付ID"`
	Payment         *Payment  `json:"payment,omitempty" gorm:"foreignKey:PaymentID"`

	// 退款信息
	RefundID        string    `json:"refund_id" gorm:"size:100;uniqueIndex;comment:PayPal退款ID"`
	Amount          float64   `json:"amount" gorm:"type:decimal(10,2);not null;comment:退款金额"`
	Currency        string    `json:"currency" gorm:"size:10;default:'USD';comment:货币类型"`
	Status          string    `json:"status" gorm:"size:20;default:'pending';comment:退款状态"`

	// 退款详情
	Reason          string    `json:"reason" gorm:"type:text;comment:退款原因"`
	Note            string    `json:"note" gorm:"type:text;comment:备注"`
}

// PaymentWebhook Webhook记录模型
type PaymentWebhook struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Webhook信息
	EventID         string    `json:"event_id" gorm:"size:100;uniqueIndex;comment:PayPal事件ID"`
	EventType       string    `json:"event_type" gorm:"size:100;index;comment:事件类型"`
	ResourceType    string    `json:"resource_type" gorm:"size:100;comment:资源类型"`
	Summary         string    `json:"summary" gorm:"type:text;comment:事件摘要"`

	// 关联资源
	ResourceID      string    `json:"resource_id" gorm:"size:100;index;comment:资源ID"`

	// 状态
	Status          string    `json:"status" gorm:"size:20;default:'pending';index;comment:处理状态"`

	// 原始数据
	RawData         string    `json:"raw_data" gorm:"type:json;comment:原始JSON数据"`
	RawHeaders      string    `json:"raw_headers" gorm:"type:json;comment:原始请求头"`
}

// TableName 指定表名
func (Payment) TableName() string {
	return "payments"
}

func (PaymentRefund) TableName() string {
	return "payment_refunds"
}

func (PaymentWebhook) TableName() string {
	return "payment_webhooks"
}

// IsPaid 检查支付是否已完成
func (p *Payment) IsPaid() bool {
	return p.Status == "paid"
}

// IsPending 检查支付是否待处理
func (p *Payment) IsPending() bool {
	return p.Status == "pending"
}

// IsFailed 检查支付是否失败
func (p *Payment) IsFailed() bool {
	return p.Status == "failed"
}

// IsRefunded 检查是否已退款
func (p *Payment) IsRefunded() bool {
	return p.Status == "refunded"
}
