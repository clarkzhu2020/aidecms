package controllers

import (
	"context"
	"time"

	"github.com/clarkzhu2020/aidecms/pkg/response"
	"github.com/clarkzhu2020/aidecms/pkg/web3"
	"github.com/cloudwego/hertz/pkg/app"
)

// Web3Controller Web3 控制器
type Web3Controller struct{}

// GetBalance 获取地址余额
// @Summary 获取区块链地址余额
// @Tags Web3
// @Accept json
// @Produce json
// @Param chain path string true "链类型 (bitcoin/ethereum/bsc/solana)"
// @Param address path string true "区块链地址"
// @Success 200 {object} map[string]interface{}
// @Router /api/web3/{chain}/balance/{address} [get]
func (w *Web3Controller) GetBalance(ctx context.Context, c *app.RequestContext) {
	chain := web3.Chain(c.Param("chain"))
	address := c.Param("address")

	if address == "" {
		response.Error(c, 400, "Bad Request", "address is required")
		return
	}

	balance, err := web3.GetManager().GetBalance(ctx, chain, address)
	if err != nil {
		response.Error(c, 500, "Internal Server Error", err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"chain":   chain,
		"address": address,
		"balance": balance,
	}, "success")
}

// GetTransaction 获取交易信息
// @Summary 获取区块链交易信息
// @Tags Web3
// @Accept json
// @Produce json
// @Param chain path string true "链类型 (bitcoin/ethereum/bsc/solana)"
// @Param hash path string true "交易哈希"
// @Success 200 {object} map[string]interface{}
// @Router /api/web3/{chain}/transaction/{hash} [get]
func (w *Web3Controller) GetTransaction(ctx context.Context, c *app.RequestContext) {
	chain := web3.Chain(c.Param("chain"))
	txHash := c.Param("hash")

	if txHash == "" {
		response.Error(c, 400, "Bad Request", "transaction hash is required")
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tx, err := web3.GetManager().GetTransaction(timeoutCtx, chain, txHash)
	if err != nil {
		response.Error(c, 500, "Internal Server Error", err.Error())
		return
	}

	response.Success(c, tx, "success")
}

// GetBlockNumber 获取最新区块高度
// @Summary 获取最新区块高度
// @Tags Web3
// @Accept json
// @Produce json
// @Param chain path string true "链类型 (bitcoin/ethereum/bsc/solana)"
// @Success 200 {object} map[string]interface{}
// @Router /api/web3/{chain}/block-number [get]
func (w *Web3Controller) GetBlockNumber(ctx context.Context, c *app.RequestContext) {
	chain := web3.Chain(c.Param("chain"))

	client, err := web3.GetManager().GetClient(chain)
	if err != nil {
		response.Error(c, 500, "Internal Server Error", err.Error())
		return
	}

	blockNumber, err := client.GetBlockNumber(ctx)
	if err != nil {
		response.Error(c, 500, "Internal Server Error", err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"chain":        chain,
		"block_number": blockNumber,
	}, "success")
}

// GetWalletInfo 获取钱包信息
// @Summary 获取钱包详细信息
// @Tags Web3
// @Accept json
// @Produce json
// @Param chain path string true "链类型 (bitcoin/ethereum/bsc/solana)"
// @Param address path string true "钱包地址"
// @Success 200 {object} map[string]interface{}
// @Router /api/web3/{chain}/wallet/{address} [get]
func (w *Web3Controller) GetWalletInfo(ctx context.Context, c *app.RequestContext) {
	chain := web3.Chain(c.Param("chain"))
	address := c.Param("address")

	if address == "" {
		response.Error(c, 400, "Bad Request", "address is required")
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	info, err := web3.GetWalletInfo(timeoutCtx, chain, address)
	if err != nil {
		response.Error(c, 500, "Internal Server Error", err.Error())
		return
	}

	response.Success(c, info, "success")
}

// GetSupportedChains 获取支持的链列表
// @Summary 获取支持的区块链列表
// @Tags Web3
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/web3/chains [get]
func (w *Web3Controller) GetSupportedChains(ctx context.Context, c *app.RequestContext) {
	chains := web3.GetManager().GetSupportedChains()

	response.Success(c, map[string]interface{}{
		"chains": chains,
		"count":  len(chains),
	}, "success")
}

// ValidateAddress 验证地址格式
// @Summary 验证区块链地址格式
// @Tags Web3
// @Accept json
// @Produce json
// @Param chain path string true "链类型 (bitcoin/ethereum/bsc/solana)"
// @Param address path string true "区块链地址"
// @Success 200 {object} map[string]interface{}
// @Router /api/web3/{chain}/validate/{address} [get]
func (w *Web3Controller) ValidateAddress(ctx context.Context, c *app.RequestContext) {
	chain := web3.Chain(c.Param("chain"))
	address := c.Param("address")

	if address == "" {
		response.Error(c, 400, "Bad Request", "address is required")
		return
	}

	err := web3.ValidateAddress(chain, address)
	valid := err == nil

	result := map[string]interface{}{
		"chain":   chain,
		"address": address,
		"valid":   valid,
	}

	if err != nil {
		result["error"] = err.Error()
	}

	response.Success(c, result, "success")
}

// GetMultiChainBalances 获取多链余额
// @Summary 获取同一地址在多条链上的余额
// @Tags Web3
// @Accept json
// @Produce json
// @Param body body web3.MultiChainAddress true "多链地址"
// @Success 200 {object} map[string]interface{}
// @Router /api/web3/multi-balance [post]
func (w *Web3Controller) GetMultiChainBalances(ctx context.Context, c *app.RequestContext) {
	var addresses web3.MultiChainAddress
	if err := c.BindJSON(&addresses); err != nil {
		response.Error(c, 400, "Bad Request", "invalid request body")
		return
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	balances, err := addresses.GetAllBalances(timeoutCtx)
	if err != nil {
		response.Error(c, 500, "Internal Server Error", err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"addresses": addresses,
		"balances":  balances,
	}, "success")
}
