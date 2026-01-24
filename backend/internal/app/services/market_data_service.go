package services

import (
	"context"
	"strconv"
	"time"

	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"github.com/clarkzhu2020/aidecms/pkg/database"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type MarketDataService struct {
	clickhouse *database.ClickHouseManager
}

func NewMarketDataService(ch *database.ClickHouseManager) *MarketDataService {
	return &MarketDataService{
		clickhouse: ch,
	}
}

// SavePrice records a price point to ClickHouse
func (s *MarketDataService) SavePrice(exchange, pair, priceStr string) {
	if s.clickhouse == nil || s.clickhouse.Conn == nil {
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		hlog.Errorf("Invalid price format for %s %s: %s", exchange, pair, priceStr)
		return
	}

	data := &models.MarketData{
		Exchange:  exchange,
		Pair:      pair,
		Price:     price,
		Timestamp: time.Now().UTC(),
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.clickhouse.Conn.Exec(ctx, "INSERT INTO market_data (exchange, pair, price, timestamp) VALUES (?, ?, ?, ?)",
			data.Exchange, data.Pair, data.Price, data.Timestamp)
		
		if err != nil {
			hlog.Errorf("Failed to save market data to ClickHouse: %v", err)
		}
	}()
}
