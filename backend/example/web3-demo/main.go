package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/clarkzhu2020/aidecms/pkg/web3"
)

func main() {
	fmt.Println("=== AideCMS Web3 Integration Demo ===\n")

	// 初始化 Web3 客户端
	fmt.Println("1. Initializing Web3 clients...")
	if err := web3.InitializeClients(); err != nil {
		log.Fatalf("Failed to initialize Web3 clients: %v", err)
	}

	manager := web3.GetManager()
	chains := manager.GetSupportedChains()
	fmt.Printf("✅ Initialized %d chain(s): %v\n\n", len(chains), chains)

	// 演示 1: 查询 Ethereum 余额
	demoEthereumBalance(manager)

	// 演示 2: 查询 Bitcoin 区块高度
	demoBitcoinBlockHeight(manager)

	// 演示 3: 查询 Solana 余额
	demoSolanaBalance(manager)

	// 演示 4: 地址验证
	demoAddressValidation()

	// 演示 5: 多链余额查询
	demoMultiChainBalance()

	// 演示 6: Ethereum 高级功能
	demoEthereumAdvanced()

	fmt.Println("\n=== Demo Completed ===")
}

func demoEthereumBalance(manager *web3.Manager) {
	fmt.Println("2. Querying Ethereum Balance...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 使用 Vitalik Buterin 的地址作为示例
	address := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"

	balance, err := manager.GetBalance(ctx, web3.Ethereum, address)
	if err != nil {
		log.Printf("❌ Error: %v\n\n", err)
		return
	}

	fmt.Printf("Address: %s\n", address)
	fmt.Printf("Balance: %s wei\n", balance)

	// 获取以 ETH 为单位的余额
	client, _ := manager.GetClient(web3.Ethereum)
	if ethClient, ok := client.(*web3.EthereumClient); ok {
		balanceEth, err := ethClient.GetBalanceInEther(ctx, address)
		if err == nil {
			fmt.Printf("Balance: %s ETH\n", balanceEth)
		}
	}

	fmt.Println("✅ Ethereum balance query completed\n")
}

func demoBitcoinBlockHeight(manager *web3.Manager) {
	fmt.Println("3. Querying Bitcoin Block Height...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := manager.GetClient(web3.Bitcoin)
	if err != nil {
		log.Printf("❌ Error: %v\n\n", err)
		return
	}

	blockHeight, err := client.GetBlockNumber(ctx)
	if err != nil {
		log.Printf("❌ Error: %v\n\n", err)
		return
	}

	fmt.Printf("Bitcoin Block Height: %d\n", blockHeight)
	fmt.Println("✅ Bitcoin block height query completed\n")
}

func demoSolanaBalance(manager *web3.Manager) {
	fmt.Println("4. Querying Solana Balance...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Solana Foundation 地址
	address := "7Np41oeYqPefeNQEHSv1UDhYrehxin3NStELsSKCT4K2"

	balance, err := manager.GetBalance(ctx, web3.Solana, address)
	if err != nil {
		log.Printf("❌ Error: %v\n\n", err)
		return
	}

	fmt.Printf("Address: %s\n", address)
	fmt.Printf("Balance: %s lamports\n", balance)

	// 获取以 SOL 为单位的余额
	client, _ := manager.GetClient(web3.Solana)
	if solClient, ok := client.(*web3.SolanaClient); ok {
		balanceSOL, err := solClient.GetBalanceInSOL(ctx, address)
		if err == nil {
			fmt.Printf("Balance: %s SOL\n", balanceSOL)
		}
	}

	fmt.Println("✅ Solana balance query completed\n")
}

func demoAddressValidation() {
	fmt.Println("5. Address Validation Demo...")

	testCases := []struct {
		chain   web3.Chain
		address string
	}{
		{web3.Ethereum, "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"},
		{web3.Ethereum, "invalid-address"},
		{web3.Bitcoin, "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"},
		{web3.Bitcoin, "invalid-btc-address"},
		{web3.Solana, "7EqQdEULxWcraVx3mXKFjc84LhCkMGZCkRuDpvcMwJeK"},
		{web3.BSC, "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"},
	}

	for _, tc := range testCases {
		err := web3.ValidateAddress(tc.chain, tc.address)
		if err != nil {
			fmt.Printf("❌ %s: %s - Invalid (%v)\n", tc.chain, tc.address, err)
		} else {
			fmt.Printf("✅ %s: %s - Valid\n", tc.chain, tc.address)
		}
	}

	fmt.Println()
}

func demoMultiChainBalance() {
	fmt.Println("6. Multi-Chain Balance Query...")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	addresses := web3.MultiChainAddress{
		Ethereum: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
		BSC:      "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
		Solana:   "7Np41oeYqPefeNQEHSv1UDhYrehxin3NStELsSKCT4K2",
	}

	balances, err := addresses.GetAllBalances(ctx)
	if err != nil {
		log.Printf("❌ Error: %v\n\n", err)
		return
	}

	fmt.Println("Balances across multiple chains:")
	for chain, balance := range balances {
		fmt.Printf("  %s: %s\n", chain, balance)
	}

	fmt.Println("✅ Multi-chain query completed\n")
}

func demoEthereumAdvanced() {
	fmt.Println("7. Ethereum Advanced Features...")

	client, err := web3.NewEthereumClient("https://rpc.ankr.com/eth")
	if err != nil {
		log.Printf("❌ Error: %v\n\n", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 获取链 ID
	chainID, err := client.GetChainID(ctx)
	if err == nil {
		fmt.Printf("Chain ID: %s\n", chainID.String())
	}

	// 获取 Gas 价格
	gasPrice, err := client.GetGasPrice(ctx)
	if err == nil {
		fmt.Printf("Current Gas Price: %s wei\n", gasPrice)
	}

	// 获取最新区块
	blockNumber, err := client.GetBlockNumber(ctx)
	if err == nil {
		fmt.Printf("Latest Block: %d\n", blockNumber)
	}

	// 检查是否为合约地址
	usdtAddress := "0xdAC17F958D2ee523a2206206994597C13D831ec7" // USDT Contract
	isContract, err := client.IsContract(ctx, usdtAddress)
	if err == nil {
		fmt.Printf("Is %s a contract? %v\n", usdtAddress, isContract)
	}

	fmt.Println("✅ Ethereum advanced features demo completed\n")
}
