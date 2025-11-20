# Web3 区块链集成

AideCMS 提供完整的 Web3 区块链集成功能，支持 Bitcoin、Ethereum、BSC（Binance Smart Chain）和 Solana 等主流区块链网络。

## 功能特性

### ✅ 支持的区块链

- **Bitcoin** - 比特币主网和测试网
- **Ethereum** - 以太坊主网、Sepolia、Goerli 等测试网
- **BSC** - 币安智能链主网和测试网
- **Solana** - Solana 主网、测试网和开发网

### ✅ 核心功能

- ✅ 地址余额查询
- ✅ 交易信息查询
- ✅ 最新区块高度查询
- ✅ 钱包信息查询
- ✅ 地址格式验证
- ✅ 多链余额查询
- ✅ Gas 价格查询（EVM 链）
- ✅ 合约代码查询（EVM 链）
- ✅ SPL Token 查询（Solana）

## 快速开始

### 1. 环境配置

在 `.env` 文件中添加 RPC 端点配置：

```env
# Ethereum
WEB3_ETHEREUM_RPC=https://mainnet.infura.io/v3/YOUR-PROJECT-ID

# BSC
WEB3_BSC_RPC=https://bsc-dataseed.binance.org/

# Bitcoin
WEB3_BITCOIN_RPC=https://bitcoin-mainnet.core.chainstack.com
WEB3_BITCOIN_API_KEY=your-api-key

# Solana
WEB3_SOLANA_RPC=https://api.mainnet-beta.solana.com
```

### 2. 初始化客户端

```go
package main

import (
    "github.com/chenyusolar/aidecms/pkg/web3"
)

func main() {
    // 加载配置并初始化所有客户端
    if err := web3.InitializeClients(); err != nil {
        panic(err)
    }
    
    // 获取管理器
    manager := web3.GetManager()
    
    // 查看已初始化的链
    chains := manager.GetSupportedChains()
    fmt.Printf("Initialized chains: %v\n", chains)
}
```

### 3. 查询地址余额

#### Ethereum / BSC

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/chenyusolar/aidecms/pkg/web3"
)

func main() {
    web3.InitializeClients()
    manager := web3.GetManager()
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // 查询 Ethereum 余额
    address := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"
    balance, err := manager.GetBalance(ctx, web3.Ethereum, address)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Balance: %s wei\n", balance)
    
    // 使用客户端获取更多信息
    client, _ := manager.GetClient(web3.Ethereum)
    ethClient := client.(*web3.EthereumClient)
    
    // 以 Ether 为单位获取余额
    balanceEth, _ := ethClient.GetBalanceInEther(ctx, address)
    fmt.Printf("Balance: %s ETH\n", balanceEth)
}
```

#### Bitcoin

```go
func main() {
    web3.InitializeClients()
    manager := web3.GetManager()
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    address := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
    balance, err := manager.GetBalance(ctx, web3.Bitcoin, address)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("BTC Balance: %s\n", balance)
}
```

#### Solana

```go
func main() {
    web3.InitializeClients()
    manager := web3.GetManager()
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    address := "7EqQdEULxWcraVx3mXKFjc84LhCkMGZCkRuDpvcMwJeK"
    balance, err := manager.GetBalance(ctx, web3.Solana, address)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("SOL Balance: %s lamports\n", balance)
    
    // 以 SOL 为单位获取余额
    client, _ := manager.GetClient(web3.Solana)
    solClient := client.(*web3.SolanaClient)
    balanceSOL, _ := solClient.GetBalanceInSOL(ctx, address)
    fmt.Printf("SOL Balance: %s SOL\n", balanceSOL)
}
```

### 4. 查询交易信息

```go
func main() {
    web3.InitializeClients()
    manager := web3.GetManager()
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Ethereum 交易
    txHash := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
    tx, err := manager.GetTransaction(ctx, web3.Ethereum, txHash)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Transaction Hash: %s\n", tx.Hash)
    fmt.Printf("From: %s\n", tx.From)
    fmt.Printf("To: %s\n", tx.To)
    fmt.Printf("Value: %s wei\n", tx.Value)
    fmt.Printf("Block: %d\n", tx.BlockNumber)
    fmt.Printf("Status: %s\n", tx.Status)
    fmt.Printf("Gas Used: %d\n", tx.GasUsed)
}
```

### 5. 获取最新区块

```go
func main() {
    web3.InitializeClients()
    manager := web3.GetManager()
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    client, _ := manager.GetClient(web3.Ethereum)
    blockNumber, err := client.GetBlockNumber(ctx)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Latest block: %d\n", blockNumber)
}
```

### 6. 多链余额查询

```go
func main() {
    web3.InitializeClients()
    
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()
    
    addresses := web3.MultiChainAddress{
        Bitcoin:  "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
        Ethereum: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
        BSC:      "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
        Solana:   "7EqQdEULxWcraVx3mXKFjc84LhCkMGZCkRuDpvcMwJeK",
    }
    
    balances, err := addresses.GetAllBalances(ctx)
    if err != nil {
        panic(err)
    }
    
    for chain, balance := range balances {
        fmt.Printf("%s: %s\n", chain, balance)
    }
}
```

### 7. 钱包信息查询

```go
func main() {
    web3.InitializeClients()
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    address := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"
    info, err := web3.GetWalletInfo(ctx, web3.Ethereum, address)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Address: %s\n", info.Address)
    fmt.Printf("Chain: %s\n", info.Chain)
    fmt.Printf("Balance: %s\n", info.Balance)
    fmt.Printf("Nonce: %d\n", info.Nonce)
    fmt.Printf("Transaction Count: %d\n", info.TxCount)
}
```

## Ethereum 专用功能

### Gas 价格查询

```go
func main() {
    client, _ := web3.NewEthereumClient("https://mainnet.infura.io/v3/YOUR-PROJECT-ID")
    defer client.Close()
    
    ctx := context.Background()
    
    gasPrice, err := client.GetGasPrice(ctx)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Current gas price: %s wei\n", gasPrice)
}
```

### Gas 估算

```go
func main() {
    client, _ := web3.NewEthereumClient("https://mainnet.infura.io/v3/YOUR-PROJECT-ID")
    defer client.Close()
    
    ctx := context.Background()
    
    from := "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"
    to := "0x1234567890abcdef1234567890abcdef12345678"
    value := big.NewInt(1000000000000000000) // 1 ETH
    
    gas, err := client.EstimateGas(ctx, from, to, "", value)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Estimated gas: %d\n", gas)
}
```

### 合约查询

```go
func main() {
    client, _ := web3.NewEthereumClient("https://mainnet.infura.io/v3/YOUR-PROJECT-ID")
    defer client.Close()
    
    ctx := context.Background()
    
    contractAddress := "0xdAC17F958D2ee523a2206206994597C13D831ec7" // USDT
    
    // 检查是否为合约
    isContract, err := client.IsContract(ctx, contractAddress)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Is contract: %v\n", isContract)
    
    // 获取合约代码
    code, err := client.GetCode(ctx, contractAddress)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Contract code length: %d bytes\n", len(code))
}
```

## Solana 专用功能

### SPL Token 余额

```go
func main() {
    client := web3.NewSolanaClient("https://api.mainnet-beta.solana.com")
    defer client.Close()
    
    ctx := context.Background()
    
    tokenAccount := "TokenAccountAddress..."
    balance, err := client.GetTokenBalance(ctx, tokenAccount)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Token balance: %s\n", balance)
}
```

### 账户信息

```go
func main() {
    client := web3.NewSolanaClient("https://api.mainnet-beta.solana.com")
    defer client.Close()
    
    ctx := context.Background()
    
    address := "7EqQdEULxWcraVx3mXKFjc84LhCkMGZCkRuDpvcMwJeK"
    accountInfo, err := client.GetAccountInfo(ctx, address)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Account info: %+v\n", accountInfo)
}
```

## API 路由

在 `routes/api.go` 中注册 Web3 路由：

```go
package routes

import (
    "github.com/chenyusolar/aidecms/app/Http/Controllers"
    "github.com/chenyusolar/aidecms/pkg/framework"
)

func RegisterWeb3Routes(router *framework.RouteGroup) {
    web3Controller := &Controllers.Web3Controller{}
    
    web3 := router.Group("/web3")
    {
        // 获取支持的链
        web3.GET("/chains", web3Controller.GetSupportedChains)
        
        // 地址余额
        web3.GET("/:chain/balance/:address", web3Controller.GetBalance)
        
        // 交易信息
        web3.GET("/:chain/transaction/:hash", web3Controller.GetTransaction)
        
        // 最新区块
        web3.GET("/:chain/block-number", web3Controller.GetBlockNumber)
        
        // 钱包信息
        web3.GET("/:chain/wallet/:address", web3Controller.GetWalletInfo)
        
        // 验证地址
        web3.GET("/:chain/validate/:address", web3Controller.ValidateAddress)
        
        // 多链余额
        web3.POST("/multi-balance", web3Controller.GetMultiChainBalances)
    }
}
```

## Artisan 命令

### 初始化客户端

```bash
go run cmd/artisan/main.go artisan web3 init
```

### 查看支持的链

```bash
go run cmd/artisan/main.go artisan web3 chains
```

### 查询余额

```bash
# Ethereum
go run cmd/artisan/main.go artisan web3 balance ethereum 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0

# Bitcoin
go run cmd/artisan/main.go artisan web3 balance bitcoin 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa

# Solana
go run cmd/artisan/main.go artisan web3 balance solana 7EqQdEULxWcraVx3mXKFjc84LhCkMGZCkRuDpvcMwJeK
```

### 查询交易

```bash
go run cmd/artisan/main.go artisan web3 transaction ethereum 0x1234...
```

### 查询区块

```bash
go run cmd/artisan/main.go artisan web3 block ethereum
```

### 验证地址

```bash
go run cmd/artisan/main.go artisan web3 validate ethereum 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0
```

### 钱包信息

```bash
go run cmd/artisan/main.go artisan web3 wallet ethereum 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0
```

## 测试

运行单元测试：

```bash
# 运行所有测试
go test ./pkg/web3/...

# 跳过集成测试（需要网络连接）
go test -short ./pkg/web3/...

# 运行包含集成测试
go test -v ./pkg/web3/...
```

## RPC 端点推荐

### Ethereum

**主网：**
- Infura: `https://mainnet.infura.io/v3/YOUR-PROJECT-ID`
- Alchemy: `https://eth-mainnet.g.alchemy.com/v2/YOUR-API-KEY`
- Ankr: `https://rpc.ankr.com/eth`
- QuickNode: `https://your-endpoint.quiknode.pro/YOUR-API-KEY/`

**测试网 (Sepolia)：**
- Infura: `https://sepolia.infura.io/v3/YOUR-PROJECT-ID`
- Alchemy: `https://eth-sepolia.g.alchemy.com/v2/YOUR-API-KEY`

### BSC

**主网：**
- Official: `https://bsc-dataseed.binance.org/`
- Ankr: `https://rpc.ankr.com/bsc`

**测试网：**
- Official: `https://data-seed-prebsc-1-s1.binance.org:8545/`

### Bitcoin

**主网：**
- Chainstack: `https://bitcoin-mainnet.core.chainstack.com`
- BlockCypher: `https://api.blockcypher.com/v1/btc/main`

### Solana

**主网：**
- Official: `https://api.mainnet-beta.solana.com`
- QuickNode: `https://your-endpoint.solana-mainnet.quiknode.pro/YOUR-API-KEY/`
- Ankr: `https://rpc.ankr.com/solana`

**测试网：**
- Official: `https://api.testnet.solana.com`

## 最佳实践

### 1. 超时控制

始终为区块链查询设置超时：

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

balance, err := manager.GetBalance(ctx, web3.Ethereum, address)
```

### 2. 错误处理

妥善处理网络错误和 RPC 错误：

```go
balance, err := manager.GetBalance(ctx, web3.Ethereum, address)
if err != nil {
    // 记录错误
    log.Printf("Failed to get balance: %v", err)
    
    // 返回友好的错误信息
    return "", fmt.Errorf("blockchain query failed: %w", err)
}
```

### 3. 地址验证

在查询前验证地址格式：

```go
if err := web3.ValidateAddress(web3.Ethereum, address); err != nil {
    return fmt.Errorf("invalid address: %w", err)
}
```

### 4. 连接池管理

使用单例 Manager 管理客户端连接：

```go
// 获取全局管理器
manager := web3.GetManager()

// 程序退出时关闭所有连接
defer manager.Close()
```

### 5. 缓存查询结果

对于不变的数据（如已确认交易），可以缓存结果：

```go
import "github.com/chenyusolar/aidecms/pkg/cache"

func GetTransactionCached(ctx context.Context, chain web3.Chain, txHash string) (*web3.Transaction, error) {
    cacheKey := fmt.Sprintf("web3:tx:%s:%s", chain, txHash)
    
    // 尝试从缓存获取
    var tx web3.Transaction
    if err := cache.Get(cacheKey, &tx); err == nil {
        return &tx, nil
    }
    
    // 从区块链查询
    tx, err := web3.GetManager().GetTransaction(ctx, chain, txHash)
    if err != nil {
        return nil, err
    }
    
    // 缓存已确认的交易（1小时）
    if tx.Status == "success" || tx.Status == "failed" {
        cache.Set(cacheKey, tx, 1*time.Hour)
    }
    
    return tx, nil
}
```

## 故障排查

### 连接超时

如果遇到连接超时，检查：
1. RPC 端点是否可访问
2. API Key 是否正确
3. 网络防火墙设置
4. 超时时间是否足够

### 余额查询失败

- **Ethereum/BSC**: 确保地址以 `0x` 开头，长度为 42 字符
- **Bitcoin**: 支持 Legacy、SegWit 和 Bech32 地址格式
- **Solana**: 地址为 Base58 编码，32-44 字符

### 交易查询失败

- 确认交易已被打包（不在 mempool 中）
- 检查交易哈希格式是否正确
- 某些旧交易可能需要归档节点支持

## 安全建议

1. **私钥管理**: 永远不要在代码中硬编码私钥
2. **RPC 安全**: 使用 HTTPS 端点，避免使用不受信任的 RPC
3. **API Key**: 将 API Key 存储在环境变量中
4. **Rate Limiting**: 对 RPC 调用实施速率限制
5. **输入验证**: 验证所有用户输入的地址和交易哈希

## 扩展功能

### 交易签名与发送

要实现交易签名和发送功能，需要：

1. 集成密钥管理系统
2. 实现交易构造逻辑
3. 添加 Gas 估算和 Nonce 管理
4. 支持交易重放保护

### ERC20/BEP20 代币

实现代币余额和转账功能需要：

1. ABI 编码/解码
2. 合约调用接口
3. 事件日志解析

### NFT 支持

添加 NFT 功能需要：

1. ERC721/ERC1155 接口实现
2. 元数据查询
3. 所有权验证

## 相关资源

- [Ethereum JSON-RPC API](https://ethereum.org/en/developers/docs/apis/json-rpc/)
- [Bitcoin RPC API](https://developer.bitcoin.org/reference/rpc/)
- [Solana JSON RPC API](https://docs.solana.com/api/http)
- [go-ethereum Documentation](https://geth.ethereum.org/docs)

## 许可证

MIT License
