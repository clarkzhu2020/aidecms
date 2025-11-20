# Hyperliquid 去中心化交易所集成

## 概述

Hyperliquid 是一个高性能的去中心化永续合约交易所（DEX），提供链上订单簿、低延迟交易和高资本效率。AideCMS 框架集成了 Hyperliquid API，支持查询余额、持仓、价格以及交易功能。

## 特性

- ✅ 去中心化交易所（DEX）
- ✅ 链上订单簿
- ✅ 永续合约交易
- ✅ 账户余额查询
- ✅ 持仓信息查询
- ✅ 实时价格查询
- ✅ 资金费率查询
- ✅ 订单簿查询
- ✅ 24小时交易量查询
- ✅ EIP-712 签名认证
- ✅ 支持限价单和市价单

## 配置

### 环境变量

```env
# Hyperliquid DEX
EXCHANGE_HYPERLIQUID_PRIVATE_KEY=your_ethereum_private_key
EXCHANGE_HYPERLIQUID_ADDRESS=your_ethereum_address
```

### 获取私钥

Hyperliquid 使用以太坊钱包进行认证：

1. **使用 MetaMask 或其他钱包**
   - 导出钱包私钥
   - 私钥格式：64位十六进制字符串（不含 0x 前缀）

2. **安全性建议**
   - 为 Hyperliquid 创建专用钱包
   - 不要使用主钱包的私钥
   - 只存入交易所需的资金
   - 定期检查账户活动

3. **权限说明**
   - 私钥可以进行交易操作（下单、取消订单）
   - 如果只需查询功能，可以只配置 ADDRESS
   - 完整功能需要配置 PRIVATE_KEY

## 使用指南

### 初始化客户端

```go
import "github.com/chenyusolar/aidecms/pkg/web3"

// 方法1：使用环境变量自动初始化
config := web3.LoadConfig()
config.InitializeClients()
manager := web3.GetExchangeManager()

// 方法2：手动创建客户端
privateKey := "your_private_key_without_0x"
client, err := web3.NewHyperliquidClient(privateKey)
if err != nil {
    log.Fatal(err)
}
```

### 查询账户余额

```go
ctx := context.Background()

// 获取 USDC 余额（账户总价值）
balance, err := manager.GetBalance(ctx, web3.Hyperliquid, "USDC")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("USDC Balance: %s\n", balance)

// 获取所有余额信息
balances, err := manager.GetBalances(ctx, web3.Hyperliquid)
if err != nil {
    log.Fatal(err)
}
for key, value := range balances {
    fmt.Printf("%s: %s\n", key, value)
}
// 输出示例:
// USDC: 10000.50
// account_value: 10000.50
// withdrawable: 9500.25
// margin_used: 500.25
```

### 查询持仓信息

```go
positions, err := client.GetPositions(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Println("当前持仓:")
for _, pos := range positions {
    fmt.Printf("币种: %s\n", pos.Coin)
    fmt.Printf("  数量: %s\n", pos.Size)
    fmt.Printf("  开仓价: %s\n", pos.EntryPrice)
    fmt.Printf("  持仓价值: %s\n", pos.PositionValue)
    fmt.Printf("  未实现盈亏: %s\n", pos.UnrealizedPnl)
    fmt.Printf("  杠杆: %s\n", pos.Leverage)
    fmt.Printf("  清算价: %s\n", pos.Liquidation)
    fmt.Println()
}
```

### 查询价格

```go
// 查询 BTC 价格
price, err := manager.GetPrice(ctx, web3.Hyperliquid, "BTC-USD")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("BTC Price: $%s\n", price)

// 查询 ETH 价格
price, err := manager.GetPrice(ctx, web3.Hyperliquid, "ETH")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("ETH Price: $%s\n", price)
```

### 查询资金费率

```go
fundingRate, err := client.GetFundingRate(ctx, "BTC")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("BTC Funding Rate: %s\n", fundingRate)
```

### 查询24小时交易量

```go
volume, err := client.Get24HVolume(ctx, "BTC")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("BTC 24H Volume: $%s\n", volume)
```

### 查询订单簿

```go
orderBook, err := client.GetOrderBook(ctx, "BTC")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("币种: %s\n", orderBook["coin"])
fmt.Printf("时间: %v\n", orderBook["time"])

// 买单（Bids）
if bids, ok := orderBook["bids"].([][]struct {
    Px string `json:"px"`
    Sz string `json:"sz"`
    N  int    `json:"n"`
}); ok {
    fmt.Println("买单:")
    for _, bid := range bids[0][:5] { // 显示前5档
        fmt.Printf("  价格: %s, 数量: %s\n", bid.Px, bid.Sz)
    }
}

// 卖单（Asks）
if asks, ok := orderBook["asks"].([][]struct {
    Px string `json:"px"`
    Sz string `json:"sz"`
    N  int    `json:"n"`
}); ok {
    fmt.Println("卖单:")
    for _, ask := range asks[0][:5] { // 显示前5档
        fmt.Printf("  价格: %s, 数量: %s\n", ask.Px, ask.Sz)
    }
}
```

### 下单交易

```go
// 限价买单
order := web3.OrderRequest{
    Coin:       "BTC",
    IsBuy:      true,
    Size:       0.01,
    LimitPrice: 45000.0,
    ReduceOnly: false,
}

orderID, err := client.PlaceOrder(ctx, order)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Order placed: %s\n", orderID)

// 市价卖单
marketOrder := web3.OrderRequest{
    Coin:       "BTC",
    IsBuy:      false,
    Size:       0.01,
    LimitPrice: 0, // 市价单
    ReduceOnly: false,
}

orderID, err = client.PlaceOrder(ctx, marketOrder)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Market order placed: %s\n", orderID)

// 只减仓订单
closeOrder := web3.OrderRequest{
    Coin:       "BTC",
    IsBuy:      false,
    Size:       0.01,
    LimitPrice: 46000.0,
    ReduceOnly: true, // 只减仓，不会开新仓
}

orderID, err = client.PlaceOrder(ctx, closeOrder)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Close order placed: %s\n", orderID)
```

### 取消订单

```go
err := client.CancelOrder(ctx, "BTC", 12345)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Order cancelled successfully")
```

### 获取市场信息

```go
markets, err := client.GetMarketInfo(ctx)
if err != nil {
    log.Fatal(err)
}

for name, info := range markets {
    marketInfo := info.(map[string]interface{})
    fmt.Printf("%s:\n", name)
    fmt.Printf("  最大杠杆: %v\n", marketInfo["max_leverage"])
    fmt.Printf("  数量精度: %v\n", marketInfo["decimals"])
    fmt.Printf("  仅支持逐仓: %v\n", marketInfo["isolated"])
}
```

## HTTP API

### 获取余额

```http
GET /api/exchange/hyperliquid/balance/USDC
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "exchange": "hyperliquid",
    "currency": "USDC",
    "balance": "10000.50"
  }
}
```

### 获取所有余额

```http
GET /api/exchange/hyperliquid/balances
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "exchange": "hyperliquid",
    "balances": {
      "USDC": "10000.50",
      "account_value": "10000.50",
      "withdrawable": "9500.25",
      "margin_used": "500.25"
    }
  }
}
```

### 获取价格

```http
GET /api/exchange/hyperliquid/price/BTC-USD
```

**响应：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "exchange": "hyperliquid",
    "pair": "BTC-USD",
    "price": "45678.90"
  }
}
```

## CLI 命令

### 查询余额

```bash
go run . artisan exchange balance hyperliquid USDC
```

**输出：**
```
Querying USDC balance on hyperliquid...
✅ Balance: 10000.50 USDC
```

### 查询所有余额

```bash
go run . artisan exchange balances hyperliquid
```

**输出：**
```
Querying all balances on hyperliquid...
✅ Balances:
  USDC: 10000.50
  account_value: 10000.50
  withdrawable: 9500.25
  margin_used: 500.25
```

### 查询价格

```bash
go run . artisan exchange price hyperliquid BTC-USD
```

**输出：**
```
Querying BTC-USD price on hyperliquid...
✅ Price: 45678.90
```

### 价格比较

```bash
go run . artisan exchange compare BTC-USD
```

**输出：**
```
Comparing BTC-USD prices across exchanges...
✅ Prices:
  coinbase: 45678.90
  kucoin: 45680.50
  hyperliquid: 45675.20
```

## 技术细节

### 认证机制

Hyperliquid 使用基于以太坊的 EIP-712 签名进行认证：

1. **签名流程**
   - 构建符合 EIP-712 标准的消息
   - 使用以太坊私钥签名
   - 生成 (r, s, v) 签名数据
   - 附加到请求中

2. **请求格式**
   ```json
   {
     "action": {
       "type": "order",
       "orders": [...]
     },
     "signature": {
       "r": "0x...",
       "s": "0x...",
       "v": 27
     },
     "nonce": 1637150400000
   }
   ```

3. **安全性**
   - 私钥本地签名，不传输到服务器
   - 每个请求包含时间戳防重放
   - 签名确保消息完整性

### API 端点

- **信息查询**: `https://api.hyperliquid.xyz/info`
- **交易操作**: `https://api.hyperliquid.xyz/exchange`

### 数据结构

#### 余额数据
```json
{
  "marginSummary": {
    "accountValue": "10000.50",
    "totalMarginUsed": "500.25",
    "withdrawable": "9500.25"
  }
}
```

#### 持仓数据
```json
{
  "assetPositions": [
    {
      "position": {
        "coin": "BTC",
        "szi": "0.5",
        "entryPx": "45000.0",
        "positionValue": "22500.0",
        "unrealizedPnl": "500.0",
        "leverage": {"value": "2"},
        "liquidationPx": "30000.0"
      }
    }
  ]
}
```

## 交易策略示例

### 简单网格交易

```go
func gridTrading(client *web3.HyperliquidClient, coin string, basePrice float64) {
    ctx := context.Background()
    gridSize := 0.01 // 1% 网格
    orderSize := 0.01 // 0.01 BTC

    // 设置买单
    for i := 1; i <= 5; i++ {
        price := basePrice * (1 - float64(i)*gridSize)
        order := web3.OrderRequest{
            Coin:       coin,
            IsBuy:      true,
            Size:       orderSize,
            LimitPrice: price,
        }
        orderID, err := client.PlaceOrder(ctx, order)
        if err != nil {
            log.Printf("Failed to place buy order: %v", err)
            continue
        }
        log.Printf("Buy order placed at %f: %s", price, orderID)
    }

    // 设置卖单
    for i := 1; i <= 5; i++ {
        price := basePrice * (1 + float64(i)*gridSize)
        order := web3.OrderRequest{
            Coin:       coin,
            IsBuy:      false,
            Size:       orderSize,
            LimitPrice: price,
        }
        orderID, err := client.PlaceOrder(ctx, order)
        if err != nil {
            log.Printf("Failed to place sell order: %v", err)
            continue
        }
        log.Printf("Sell order placed at %f: %s", price, orderID)
    }
}
```

### 止损止盈

```go
func stopLossAndTakeProfit(client *web3.HyperliquidClient) {
    ctx := context.Background()
    
    // 获取持仓
    positions, err := client.GetPositions(ctx)
    if err != nil {
        log.Fatal(err)
    }

    for _, pos := range positions {
        size, _ := strconv.ParseFloat(pos.Size, 64)
        entryPrice, _ := strconv.ParseFloat(pos.EntryPrice, 64)
        
        if size > 0 { // 多仓
            // 止损 -2%
            stopLoss := entryPrice * 0.98
            // 止盈 +5%
            takeProfit := entryPrice * 1.05
            
            // 设置止损单
            stopOrder := web3.OrderRequest{
                Coin:       pos.Coin,
                IsBuy:      false,
                Size:       math.Abs(size),
                LimitPrice: stopLoss,
                ReduceOnly: true,
            }
            client.PlaceOrder(ctx, stopOrder)
            
            // 设置止盈单
            profitOrder := web3.OrderRequest{
                Coin:       pos.Coin,
                IsBuy:      false,
                Size:       math.Abs(size),
                LimitPrice: takeProfit,
                ReduceOnly: true,
            }
            client.PlaceOrder(ctx, profitOrder)
        }
    }
}
```

## 常见问题

### 1. 如何获取私钥？

使用 MetaMask：
1. 打开 MetaMask 扩展
2. 点击账户详情
3. 导出私钥
4. 输入密码确认
5. 复制私钥（移除 0x 前缀）

### 2. 私钥安全吗？

- ✅ 私钥仅在本地使用，不会发送到任何服务器
- ✅ 用于本地签名请求
- ⚠️ 建议使用专用钱包，不要使用主钱包
- ⚠️ 不要在代码中硬编码私钥
- ⚠️ 使用环境变量或密钥管理服务

### 3. 支持哪些币种？

Hyperliquid 支持多种永续合约：
- BTC-USD
- ETH-USD
- SOL-USD
- MATIC-USD
- ARB-USD
- OP-USD
- 等更多...

使用 `GetMarketInfo()` 查询完整列表。

### 4. 手续费如何计算？

- Maker 费率: -0.02%（返佣）
- Taker 费率: 0.05%
- 资金费率: 每8小时结算一次

### 5. 最大杠杆是多少？

不同币种的最大杠杆不同：
- BTC: 50x
- ETH: 50x
- 其他币种: 20x-50x

使用 `GetMarketInfo()` 查询具体币种的最大杠杆。

### 6. 如何查看历史交易？

目前实现了实时查询功能，历史交易记录可通过 Hyperliquid 网页界面查看。

### 7. 支持现货交易吗？

Hyperliquid 主要是永续合约交易所，不支持传统现货交易。

## 风险提示

⚠️ **重要警告**：

1. **市场风险**: 加密货币市场波动剧烈，可能导致重大损失
2. **杠杆风险**: 高杠杆交易可能导致快速清算
3. **智能合约风险**: DEX 依赖智能合约，可能存在未知漏洞
4. **私钥安全**: 私钥泄露将导致资金完全损失
5. **测试建议**: 先用小额资金测试，熟悉后再增加投入

## 最佳实践

### 1. 安全管理私钥

```go
// 不要这样做 ❌
privateKey := "1234567890abcdef..."

// 应该这样做 ✅
privateKey := os.Getenv("EXCHANGE_HYPERLIQUID_PRIVATE_KEY")
if privateKey == "" {
    log.Fatal("Private key not configured")
}
```

### 2. 错误处理

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

balance, err := client.GetBalance(ctx, "USDC")
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        log.Println("Request timeout, retrying...")
        // 重试逻辑
    } else {
        log.Printf("Failed to get balance: %v", err)
    }
    return
}
```

### 3. 资金管理

```go
// 计算合理的仓位大小
func calculatePositionSize(accountValue, risk float64) float64 {
    // 风险不超过账户的 2%
    maxRisk := accountValue * 0.02
    stopLoss := 0.05 // 5% 止损
    positionSize := maxRisk / stopLoss
    return positionSize
}
```

### 4. 监控持仓

```go
// 定期检查持仓和风险
func monitorPositions(client *web3.HyperliquidClient) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        positions, err := client.GetPositions(context.Background())
        if err != nil {
            log.Printf("Failed to get positions: %v", err)
            continue
        }

        for _, pos := range positions {
            pnl, _ := strconv.ParseFloat(pos.UnrealizedPnl, 64)
            if pnl < -1000 { // 亏损超过 $1000
                log.Printf("⚠️ Large loss detected: %s %s PnL: %f",
                    pos.Coin, pos.Size, pnl)
                // 发送告警或自动平仓
            }
        }
    }
}
```

## 参考资料

- [Hyperliquid 官网](https://hyperliquid.xyz/)
- [Hyperliquid 文档](https://hyperliquid.gitbook.io/)
- [Hyperliquid API 文档](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api)
- [EIP-712 规范](https://eips.ethereum.org/EIPS/eip-712)

## 更新日志

### v1.0.0 (2024-01-XX)
- ✅ 初始版本
- ✅ 支持余额查询
- ✅ 支持持仓查询
- ✅ 支持价格查询
- ✅ 支持交易功能
- ✅ 支持资金费率查询
- ✅ 支持订单簿查询

## 许可证

MIT License
