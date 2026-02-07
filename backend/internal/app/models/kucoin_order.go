package models

import (
	"time"
	"gorm.io/gorm"
)

// KuCoinOrder KuCoin订单模型
type KuCoinOrder struct {
	gorm.Model
	OrderID       string  `gorm:"type:varchar(100);uniqueIndex;not null" json:"order_id"`
	ClientOrderID string  `gorm:"type:varchar(100);index" json:"client_order_id"`
	Symbol        string  `gorm:"type:varchar(50);not null;index" json:"symbol"`
	Side          string  `gorm:"type:varchar(10);not null" json:"side"` // buy, sell
	Type          string  `gorm:"type:varchar(20);not null" json:"type"` // limit, market, stop
	Price         string  `gorm:"type:decimal(30,18)" json:"price"`
	Size          string  `gorm:"type:decimal(30,18);not null" json:"size"`
	DealSize      string  `gorm:"type:decimal(30,18)" json:"deal_size"`
	DealFunds     string  `gorm:"type:decimal(30,18)" json:"deal_funds"`
	Fee           string  `gorm:"type:decimal(30,18)" json:"fee"`
	FeeCurrency   string  `gorm:"type:varchar(20)" json:"fee_currency"`
	StopPrice     string  `gorm:"type:decimal(30,18)" json:"stop_price"`
	TimeInForce   string  `gorm:"type:varchar(10)" json:"time_in_force"` // GTC, IOC, FOK
	Status        string  `gorm:"type:varchar(20);index" json:"status"` // open, done, match, canceled
	KuCoinCreatedAt int64 `json:"kucoin_created_at"`
	KuCoinUpdatedAt int64 `json:"kucoin_updated_at"`
	Remark        string  `gorm:"type:text" json:"remark"`
}

// TableName 指定表名
func (KuCoinOrder) TableName() string {
	return "kucoin_orders"
}

// KuCoinTrade KuCoin成交记录模型
type KuCoinTrade struct {
	gorm.Model
	TradeID       string  `gorm:"type:varchar(100);uniqueIndex;not null" json:"trade_id"`
	OrderID       string  `gorm:"type:varchar(100);index;not null" json:"order_id"`
	Symbol        string  `gorm:"type:varchar(50);not null;index" json:"symbol"`
	Side          string  `gorm:"type:varchar(10);not null" json:"side"` // buy, sell
	Price         string  `gorm:"type:decimal(30,18);not null" json:"price"`
	Size          string  `gorm:"type:decimal(30,18);not null" json:"size"`
	Fee           string  `gorm:"type:decimal(30,18)" json:"fee"`
	FeeCurrency   string  `gorm:"type:varchar(20)" json:"fee_currency"`
	CounterOrderID string `gorm:"type:varchar(100)" json:"counter_order_id"`
	ForceTaker    bool    `gorm:"default:false" json:"force_taker"`
	Liquidity     string  `gorm:"type:varchar(20)" json:"liquidity"`
	KuCoinCreatedAt int64 `json:"kucoin_created_at"`
}

// TableName 指定表名
func (KuCoinTrade) TableName() string {
	return "kucoin_trades"
}

// KuCoinAccount KuCoin账户模型
type KuCoinAccount struct {
	gorm.Model
	AccountID     string  `gorm:"type:varchar(100);uniqueIndex;not null" json:"account_id"`
	Currency      string  `gorm:"type:varchar(20);not null;index" json:"currency"`
	Type          string  `gorm:"type:varchar(20);not null" json:"type"` // main, trade, margin
	Balance       string  `gorm:"type:decimal(30,18)" json:"balance"`
	Available     string  `gorm:"type:decimal(30,18)" json:"available"`
	Holds         string  `gorm:"type:decimal(30,18)" json:"holds"`
	LastSyncedAt  time.Time `gorm:"index" json:"last_synced_at"`
}

// TableName 指定表名
func (KuCoinAccount) TableName() string {
	return "kucoin_accounts"
}

// KuCoinBalanceSnapshot KuCoin余额快照模型
type KuCoinBalanceSnapshot struct {
	gorm.Model
	SnapshotID    string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"snapshot_id"`
	AccountID     string    `gorm:"type:varchar(100);index" json:"account_id"`
	Currency      string    `gorm:"type:varchar(20);not null;index" json:"currency"`
	Balance       string    `gorm:"type:decimal(30,18)" json:"balance"`
	Available     string    `gorm:"type:decimal(30,18)" json:"available"`
	Holds         string    `gorm:"type:decimal(30,18)" json:"holds"`
	SnapshotAt    time.Time `gorm:"not null;index" json:"snapshot_at"`
}

// TableName 指定表名
func (KuCoinBalanceSnapshot) TableName() string {
	return "kucoin_balance_snapshots"
}
