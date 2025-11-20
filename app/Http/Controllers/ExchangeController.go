package controllers

import (
	"context"
	"time"

	"github.com/clarkzhu2020/aidecms/pkg/response"
	"github.com/clarkzhu2020/aidecms/pkg/web3"
	"github.com/cloudwego/hertz/pkg/app"
)

// ExchangeController 交易所控制器
type ExchangeController struct{}

// GetBalance 获取交易所余额
func (e *ExchangeController) GetBalance(ctx context.Context, c *app.RequestContext) {
	exchange := web3.Exchange(c.Param("exchange"))
	currency := c.Param("currency")

	if currency == "" {
		response.Error(c, 400, "Bad Request", "currency is required")
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	balance, err := web3.GetExchangeManager().GetBalance(timeoutCtx, exchange, currency)
	if err != nil {
		response.Error(c, 500, "Internal Server Error", err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"exchange": exchange,
		"currency": currency,
		"balance":  balance,
	}, "success")
}

// GetBalances 获取交易所所有余额
func (e *ExchangeController) GetBalances(ctx context.Context, c *app.RequestContext) {
	exchange := web3.Exchange(c.Param("exchange"))

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	balances, err := web3.GetExchangeManager().GetBalances(timeoutCtx, exchange)
	if err != nil {
		response.Error(c, 500, "Internal Server Error", err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"exchange": exchange,
		"balances": balances,
	}, "success")
}

// GetPrice 获取交易对价格
func (e *ExchangeController) GetPrice(ctx context.Context, c *app.RequestContext) {
	exchange := web3.Exchange(c.Param("exchange"))
	pair := c.Param("pair")

	if pair == "" {
		response.Error(c, 400, "Bad Request", "pair is required")
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	price, err := web3.GetExchangeManager().GetPrice(timeoutCtx, exchange, pair)
	if err != nil {
		response.Error(c, 500, "Internal Server Error", err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"exchange": exchange,
		"pair":     pair,
		"price":    price,
	}, "success")
}

// GetSupportedExchanges 获取支持的交易所
func (e *ExchangeController) GetSupportedExchanges(ctx context.Context, c *app.RequestContext) {
	exchanges := web3.GetExchangeManager().GetSupportedExchanges()

	response.Success(c, map[string]interface{}{
		"exchanges": exchanges,
		"count":     len(exchanges),
	}, "success")
}

// GetAllBalances 获取所有交易所的指定币种余额
func (e *ExchangeController) GetAllBalances(ctx context.Context, c *app.RequestContext) {
	currency := c.Param("currency")

	if currency == "" {
		response.Error(c, 400, "Bad Request", "currency is required")
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	balances, err := web3.GetAllExchangeBalances(timeoutCtx, currency)
	if err != nil {
		response.Error(c, 500, "Internal Server Error", err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"currency": currency,
		"balances": balances,
	}, "success")
}

// GetAllPrices 获取所有交易所的价格
func (e *ExchangeController) GetAllPrices(ctx context.Context, c *app.RequestContext) {
	pair := c.Param("pair")

	if pair == "" {
		response.Error(c, 400, "Bad Request", "pair is required")
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	prices, err := web3.GetAllExchangePrices(timeoutCtx, pair)
	if err != nil {
		response.Error(c, 500, "Internal Server Error", err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"pair":   pair,
		"prices": prices,
	}, "success")
}
