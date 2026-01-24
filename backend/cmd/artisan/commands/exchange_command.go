package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/clarkzhu2020/aidecms/pkg/web3"
	"github.com/spf13/cobra"
)

// ExchangeCommand Exchange命令
var ExchangeCommand = &cobra.Command{
	Use:   "exchange",
	Short: "Cryptocurrency exchange operations",
	Long:  `Query balances, prices, and other data from cryptocurrency exchanges like Coinbase and KuCoin`,
}

// exchangeBalanceCmd 查询交易所余额
var exchangeBalanceCmd = &cobra.Command{
	Use:   "balance [exchange] [currency]",
	Short: "Get balance for a currency on an exchange",
	Long:  `Query the balance of a specific currency on a cryptocurrency exchange`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		exchange := web3.Exchange(args[0])
		currency := args[1]

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		fmt.Printf("Querying %s balance on %s...\n", currency, exchange)

		balance, err := web3.GetExchangeManager().GetBalance(ctx, exchange, currency)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Printf("✅ Balance: %s %s\n", balance, currency)
	},
}

// exchangeBalancesCmd 查询交易所所有余额
var exchangeBalancesCmd = &cobra.Command{
	Use:   "balances [exchange]",
	Short: "Get all balances on an exchange",
	Long:  `Query all currency balances on a cryptocurrency exchange`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exchange := web3.Exchange(args[0])

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		fmt.Printf("Querying all balances on %s...\n", exchange)

		balances, err := web3.GetExchangeManager().GetBalances(ctx, exchange)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Println("✅ Balances:")
		for currency, balance := range balances {
			fmt.Printf("  %s: %s\n", currency, balance)
		}
	},
}

// exchangePriceCmd 查询交易对价格
var exchangePriceCmd = &cobra.Command{
	Use:   "price [exchange] [pair]",
	Short: "Get price for a trading pair",
	Long:  `Query the current price of a trading pair on an exchange (e.g., BTC-USD)`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		exchange := web3.Exchange(args[0])
		pair := args[1]

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		fmt.Printf("Querying %s price on %s...\n", pair, exchange)

		price, err := web3.GetExchangeManager().GetPrice(ctx, exchange, pair)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Printf("✅ Price: %s\n", price)
	},
}

// exchangeListCmd 列出支持的交易所
var exchangeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List supported exchanges",
	Long:  `Show all cryptocurrency exchanges supported by the platform`,
	Run: func(cmd *cobra.Command, args []string) {
		exchanges := web3.GetExchangeManager().GetSupportedExchanges()

		fmt.Println("Supported Exchanges:")
		for _, exchange := range exchanges {
			fmt.Printf("  - %s\n", exchange)
		}
		fmt.Printf("\nTotal: %d exchanges\n", len(exchanges))
	},
}

// exchangeCompareCmd 比较多个交易所的价格
var exchangeCompareCmd = &cobra.Command{
	Use:   "compare [pair]",
	Short: "Compare prices across all exchanges",
	Long:  `Query and compare the price of a trading pair across all supported exchanges`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pair := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		fmt.Printf("Comparing %s prices across exchanges...\n", pair)

		prices, err := web3.GetAllExchangePrices(ctx, pair)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Println("✅ Prices:")
		for exchange, price := range prices {
			fmt.Printf("  %s: %s\n", exchange, price)
		}
	},
}

// exchangeBalanceAllCmd 查询所有交易所的指定币种余额
var exchangeBalanceAllCmd = &cobra.Command{
	Use:   "balance-all [currency]",
	Short: "Get balance across all exchanges",
	Long:  `Query the balance of a specific currency across all supported exchanges`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		currency := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		fmt.Printf("Querying %s balance across all exchanges...\n", currency)

		balances, err := web3.GetAllExchangeBalances(ctx, currency)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Println("✅ Balances:")
		total := "0"
		for exchange, balance := range balances {
			fmt.Printf("  %s: %s\n", exchange, balance)
			// TODO: Add balances (would need decimal math library)
		}
		fmt.Printf("\nNote: Total balance calculation requires decimal math (sum: %s)\n", total)
	},
}

func init() {
	// 注册子命令
	ExchangeCommand.AddCommand(exchangeBalanceCmd)
	ExchangeCommand.AddCommand(exchangeBalancesCmd)
	ExchangeCommand.AddCommand(exchangePriceCmd)
	ExchangeCommand.AddCommand(exchangeListCmd)
	ExchangeCommand.AddCommand(exchangeCompareCmd)
	ExchangeCommand.AddCommand(exchangeBalanceAllCmd)
}

// ExchangeCommandWrapper 命令入口（用于 Artisan）
func ExchangeCommandWrapper(args []string) {
	// 将命令行参数转换为 Cobra 格式
	ExchangeCommand.SetArgs(args)
	if err := ExchangeCommand.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
