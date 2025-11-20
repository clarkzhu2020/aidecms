package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/chenyusolar/aidecms/pkg/web3"
	"github.com/spf13/cobra"
)

// Web3Cmd Web3 命令
var Web3Cmd = &cobra.Command{
	Use:   "web3",
	Short: "Web3 blockchain integration commands",
	Long:  `Commands for interacting with Bitcoin, Ethereum, BSC, and Solana blockchains`,
}

var web3BalanceCmd = &cobra.Command{
	Use:   "balance [chain] [address]",
	Short: "Get blockchain address balance",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		chain := web3.Chain(args[0])
		address := args[1]

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		balance, err := web3.GetManager().GetBalance(ctx, chain, address)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Chain: %s\n", chain)
		fmt.Printf("Address: %s\n", address)
		fmt.Printf("Balance: %s\n", balance)
	},
}

var web3TxCmd = &cobra.Command{
	Use:   "transaction [chain] [hash]",
	Short: "Get blockchain transaction details",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		chain := web3.Chain(args[0])
		txHash := args[1]

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		tx, err := web3.GetManager().GetTransaction(ctx, chain, txHash)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Chain: %s\n", chain)
		fmt.Printf("Transaction Hash: %s\n", tx.Hash)
		fmt.Printf("From: %s\n", tx.From)
		fmt.Printf("To: %s\n", tx.To)
		fmt.Printf("Value: %s\n", tx.Value)
		fmt.Printf("Block Number: %d\n", tx.BlockNumber)
		fmt.Printf("Block Hash: %s\n", tx.BlockHash)
		fmt.Printf("Status: %s\n", tx.Status)
		fmt.Printf("Timestamp: %d\n", tx.Timestamp)
		if tx.GasUsed > 0 {
			fmt.Printf("Gas Used: %d\n", tx.GasUsed)
		}
	},
}

var web3BlockCmd = &cobra.Command{
	Use:   "block [chain]",
	Short: "Get latest block number",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		chain := web3.Chain(args[0])

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		client, err := web3.GetManager().GetClient(chain)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		blockNumber, err := client.GetBlockNumber(ctx)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Chain: %s\n", chain)
		fmt.Printf("Latest Block: %d\n", blockNumber)
	},
}

var web3ChainsCmd = &cobra.Command{
	Use:   "chains",
	Short: "List supported blockchain networks",
	Run: func(cmd *cobra.Command, args []string) {
		chains := web3.GetManager().GetSupportedChains()

		if len(chains) == 0 {
			fmt.Println("No chains configured. Please initialize Web3 clients first.")
			return
		}

		fmt.Printf("Supported Chains (%d):\n", len(chains))
		for _, chain := range chains {
			fmt.Printf("  - %s\n", chain)
		}
	},
}

var web3ValidateCmd = &cobra.Command{
	Use:   "validate [chain] [address]",
	Short: "Validate blockchain address format",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		chain := web3.Chain(args[0])
		address := args[1]

		err := web3.ValidateAddress(chain, address)
		if err != nil {
			fmt.Printf("❌ Invalid address: %v\n", err)
			return
		}

		fmt.Printf("✅ Valid %s address: %s\n", chain, address)
	},
}

var web3WalletCmd = &cobra.Command{
	Use:   "wallet [chain] [address]",
	Short: "Get wallet information",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		chain := web3.Chain(args[0])
		address := args[1]

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		info, err := web3.GetWalletInfo(ctx, chain, address)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Chain: %s\n", info.Chain)
		fmt.Printf("Address: %s\n", info.Address)
		fmt.Printf("Balance: %s\n", info.Balance)
		if info.Nonce > 0 {
			fmt.Printf("Nonce: %d\n", info.Nonce)
			fmt.Printf("Transaction Count: %d\n", info.TxCount)
		}
	},
}

var web3InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Web3 clients",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing Web3 clients...")

		if err := web3.InitializeClients(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		chains := web3.GetManager().GetSupportedChains()
		fmt.Printf("✅ Initialized %d chain(s):\n", len(chains))
		for _, chain := range chains {
			fmt.Printf("  - %s\n", chain)
		}
	},
}

func init() {
	// Add sub-commands
	Web3Cmd.AddCommand(web3BalanceCmd)
	Web3Cmd.AddCommand(web3TxCmd)
	Web3Cmd.AddCommand(web3BlockCmd)
	Web3Cmd.AddCommand(web3ChainsCmd)
	Web3Cmd.AddCommand(web3ValidateCmd)
	Web3Cmd.AddCommand(web3WalletCmd)
	Web3Cmd.AddCommand(web3InitCmd)
}

// Web3Command Web3 命令入口（用于 Artisan）
func Web3Command(args []string) {
	// 将命令行参数转换为 Cobra 格式
	Web3Cmd.SetArgs(args)
	if err := Web3Cmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
