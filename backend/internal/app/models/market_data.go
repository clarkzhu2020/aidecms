package models

import "time"

// MarketData represents market price data for ClickHouse
type MarketData struct {
	Exchange  string    `json:"exchange"`
	Pair      string    `json:"pair"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}
