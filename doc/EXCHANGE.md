# 加密货币交易所集成

AideCMS 提供了与主流加密货币交易所的集成功能，支持查询账户余额、交易对价格等操作。

## 目录

- [特性](#特性)
- [支持的交易所](#支持的交易所)
- [快速开始](#快速开始)
- [配置](#配置)
- [API 文档](#api-文档)
- [CLI 命令](#cli-命令)
- [代码示例](#代码示例)
- [API 密钥安全](#api-密钥安全)

## 特性

- ✅ 支持多个主流交易所（Coinbase, KuCoin）
- ✅ 统一的接口设计
- ✅ 账户余额查询
- ✅ 交易对价格查询
- ✅ 多交易所价格比较
- ✅ 安全的 API 签名认证
- ✅ HTTP API 和 CLI 命令
- ✅ 超时控制和错误处理

## 支持的交易所

| 交易所 | 标识符 | 类型 | API 版本 | 认证方式 |
|--------|--------|------|----------|----------|
| Coinbase | `coinbase` | CEX | v2 | HMAC-SHA256 |
| KuCoin | `kucoin` | CEX | v2 | HMAC-SHA256 + Base64 |
| Hyperliquid | `hyperliquid` | DEX | v1 | EIP-712签名 |

### 计划支持

- Binance
- Kraken
- Bitfinex
- Huobi

## 快速开始

### 1. 配置环境变量

在 `.env` 文件中添加交易所 API 凭证：

```env
# Coinbase Exchange
EXCHANGE_COINBASE_API_KEY=your_coinbase_api_key
EXCHANGE_COINBASE_API_SECRET=your_coinbase_api_secret

# KuCoin Exchange
EXCHANGE_KUCOIN_API_KEY=your_kucoin_api_key
EXCHANGE_KUCOIN_API_SECRET=your_kucoin_api_secret
EXCHANGE_KUCOIN_PASSPHRASE=your_kucoin_passphrase
```

### 2. 初始化客户端

```go
import "github.com/chenyusolar/aidecms/pkg/web3"

// 自动从环境变量加载配置
config := web3.LoadConfig()
config.InitializeClients()

// 获取交易所管理器
manager := web3.GetExchangeManager()
```

### 3. 查询余额

```go
ctx := context.Background()
balance, err := manager.GetBalance(ctx, web3.ExchangeCoinbase, "BTC")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("BTC Balance: %s\n", balance)
```

## 配置

### 环境变量

#### Coinbase

- `EXCHANGE_COINBASE_API_KEY` - Coinbase API Key
- `EXCHANGE_COINBASE_API_SECRET` - Coinbase API Secret

#### KuCoin

- `EXCHANGE_KUCOIN_API_KEY` - KuCoin API Key
- `EXCHANGE_KUCOIN_API_SECRET` - KuCoin API Secret
- `EXCHANGE_KUCOIN_PASSPHRASE` - KuCoin API Passphrase

### 获取 API 密钥

#### Coinbase

1. 登录 [Coinbase Pro](https://pro.coinbase.com/)
2. 进入 API 设置页面
3. 创建新的 API Key
4. 权限设置：至少需要 "View" 权限查询余额

#### KuCoin

1. 登录 [KuCoin](https://www.kucoin.com/)
2. 进入 API Management
3. 创建 API
4. 设置 API Passphrase（自定义密码）
5. 权限设置：勾选 "General" 和 "Trade"

## API 文档

### HTTP API 端点

#### 1. 获取单个币种余额

```http
GET /api/exchange/:exchange/balance/:currency
```

**参数：**
- `exchange` - 交易所标识（coinbase, kucoin）
- `currency` - 币种代码（BTC, ETH, USD等）

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "exchange": "coinbase",
    "currency": "BTC",
    "balance": "0.12345678"
  }
}
```

#### 2. 获取所有余额

```http
GET /api/exchange/:exchange/balances
```

**参数：**
- `exchange` - 交易所标识

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "exchange": "coinbase",
    "balances": {
      "BTC": "0.12345678",
      "ETH": "2.5",
      "USD": "10000.50"
    }
  }
}
```

#### 3. 获取交易对价格

```http
GET /api/exchange/:exchange/price/:pair
```

**参数：**
- `exchange` - 交易所标识
- `pair` - 交易对（BTC-USD, ETH-USDT等）

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "exchange": "coinbase",
    "pair": "BTC-USD",
    "price": "45678.90"
  }
}
```

#### 4. 列出支持的交易所

```http
GET /api/exchange/supported
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "exchanges": ["coinbase", "kucoin"],
    "count": 2
  }
}
```

#### 5. 跨交易所余额查询

```http
GET /api/exchange/all/balance/:currency
```

**参数：**
- `currency` - 币种代码

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "currency": "BTC",
    "balances": {
      "coinbase": "0.5",
      "kucoin": "0.3"
    }
  }
}
```

#### 6. 跨交易所价格比较

```http
GET /api/exchange/all/price/:pair
```

**参数：**
- `pair` - 交易对

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "pair": "BTC-USD",
    "prices": {
      "coinbase": "45678.90",
      "kucoin": "45680.50"
    }
  }
}
```

## CLI 命令

### 列出支持的交易所

```bash
go run . artisan exchange list
```

**输出：**
```
Supported Exchanges:
  - coinbase
  - kucoin

Total: 2 exchanges
```

### 查询余额

```bash
go run . artisan exchange balance coinbase BTC
```

**输出：**
```
Querying BTC balance on coinbase...
✅ Balance: 0.12345678 BTC
```

### 查询所有余额

```bash
go run . artisan exchange balances coinbase
```

**输出：**
```
Querying all balances on coinbase...
✅ Balances:
  BTC: 0.12345678
  ETH: 2.5
  USD: 10000.50
```

### 查询价格

```bash
go run . artisan exchange price coinbase BTC-USD
```

**输出：**
```
Querying BTC-USD price on coinbase...
✅ Price: 45678.90
```

### 比较价格

```bash
go run . artisan exchange compare BTC-USD
```

**输出：**
```
Comparing BTC-USD prices across exchanges...
✅ Prices:
  coinbase: 45678.90
  kucoin: 45680.50
```

### 跨交易所余额查询

```bash
go run . artisan exchange balance-all BTC
```

**输出：**
```
Querying BTC balance across all exchanges...
✅ Balances:
  coinbase: 0.5
  kucoin: 0.3
```

## 代码示例

### 基本用法

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/chenyusolar/aidecms/pkg/web3"
)

func main() {
    // 加载配置
    config := web3.LoadConfig()
    config.InitializeClients()

    manager := web3.GetExchangeManager()
    ctx := context.Background()

    // 查询余额
    balance, err := manager.GetBalance(ctx, web3.ExchangeCoinbase, "BTC")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("BTC Balance: %s\n", balance)

    // 查询价格
    price, err := manager.GetPrice(ctx, web3.ExchangeCoinbase, "BTC-USD")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("BTC Price: $%s\n", price)
}
```

### 多交易所查询

```go
// 比较所有交易所的BTC价格
prices, err := web3.GetAllExchangePrices(ctx, "BTC-USD")
if err != nil {
    log.Fatal(err)
}

fmt.Println("BTC Prices:")
for exchange, price := range prices {
    fmt.Printf("%s: $%s\n", exchange, price)
}

// 查询所有交易所的BTC余额
balances, err := web3.GetAllExchangeBalances(ctx, "BTC")
if err != nil {
    log.Fatal(err)
}

fmt.Println("BTC Balances:")
total := 0.0
for exchange, balance := range balances {
    fmt.Printf("%s: %s BTC\n", exchange, balance)
    // 计算总余额（需要字符串转数字）
}
```

### 错误处理

```go
balance, err := manager.GetBalance(ctx, web3.ExchangeCoinbase, "BTC")
if err != nil {
    switch {
    case errors.Is(err, web3.ErrExchangeNotFound):
        fmt.Println("交易所不支持")
    case errors.Is(err, web3.ErrInvalidCredentials):
        fmt.Println("API 密钥无效")
    case errors.Is(err, context.DeadlineExceeded):
        fmt.Println("请求超时")
    default:
        fmt.Printf("未知错误: %v\n", err)
    }
    return
}
```

### 自定义超时

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

balance, err := manager.GetBalance(ctx, web3.ExchangeCoinbase, "BTC")
```

## API 密钥安全

### 最佳实践

1. **权限最小化**
   - 只授予必要的权限（只读权限足以查询余额和价格）
   - 不要授予提现权限

2. **密钥保护**
   - 永远不要在代码中硬编码 API 密钥
   - 使用环境变量或密钥管理服务
   - 不要将 `.env` 文件提交到版本控制

3. **IP 白名单**
   - 在交易所设置中配置 IP 白名单
   - 限制 API 访问来源

4. **定期轮换**
   - 定期更换 API 密钥
   - 删除不再使用的密钥

5. **监控**
   - 监控 API 使用情况
   - 设置异常告警

### 示例：使用密钥管理服务

```go
// 从 AWS Secrets Manager 加载密钥
import "github.com/aws/aws-sdk-go/service/secretsmanager"

func loadAPIKeys() (*web3.Config, error) {
    // 从密钥管理服务获取凭证
    secret, err := getSecretFromAWS("clarkgo/exchange-keys")
    if err != nil {
        return nil, err
    }

    config := &web3.Config{
        CoinbaseAPIKey:    secret["COINBASE_API_KEY"],
        CoinbaseAPISecret: secret["COINBASE_API_SECRET"],
        // ...
    }
    return config, nil
}
```

## 故障排查

### 常见错误

#### 1. 认证失败

**错误信息：**
```
Error: invalid signature
```

**解决方案：**
- 检查 API Key 和 Secret 是否正确
- 确认 API Key 有足够的权限
- 检查系统时间是否同步

#### 2. 请求超时

**错误信息：**
```
Error: context deadline exceeded
```

**解决方案：**
- 增加超时时间
- 检查网络连接
- 确认交易所 API 服务状态

#### 3. 交易所不支持

**错误信息：**
```
Error: exchange not found
```

**解决方案：**
- 检查交易所标识符是否正确
- 确认该交易所已配置
- 运行 `exchange list` 查看支持的交易所

#### 4. 币种不存在

**错误信息：**
```
Error: currency not found
```

**解决方案：**
- 检查币种代码是否正确（BTC, ETH, USD等）
- 确认该币种在交易所账户中存在

## 技术细节

### 认证机制

#### Coinbase

使用 HMAC-SHA256 签名：

```
signature = HMAC-SHA256(timestamp + method + path + body, api_secret)
```

请求头：
- `CB-ACCESS-KEY`: API Key
- `CB-ACCESS-SIGN`: 签名
- `CB-ACCESS-TIMESTAMP`: Unix 时间戳
- `CB-ACCESS-PASSPHRASE`: Passphrase（如果有）

#### KuCoin

使用 HMAC-SHA256 + Base64 编码：

```
signature = Base64(HMAC-SHA256(timestamp + method + path + body, api_secret))
passphrase = Base64(HMAC-SHA256(passphrase, api_secret))
```

请求头：
- `KC-API-KEY`: API Key
- `KC-API-SIGN`: 签名
- `KC-API-TIMESTAMP`: 毫秒时间戳
- `KC-API-PASSPHRASE`: 加密后的 Passphrase
- `KC-API-KEY-VERSION`: "2"

### 架构设计

```
ExchangeManager (单例)
    ├── CoinbaseClient (实现 ExchangeClient 接口)
    │   ├── GetBalance()
    │   ├── GetBalances()
    │   ├── GetPrice()
    │   └── signRequest()
    │
    └── KuCoinClient (实现 ExchangeClient 接口)
        ├── GetBalance()
        ├── GetBalances()
        ├── GetPrice()
        └── signRequest()
```

### 接口定义

```go
type ExchangeClient interface {
    GetBalance(ctx context.Context, currency string) (string, error)
    GetBalances(ctx context.Context) (map[string]string, error)
    GetPrice(ctx context.Context, pair string) (string, error)
}
```

## 扩展开发

### 添加新交易所

1. **创建客户端文件**

```go
// pkg/web3/binance.go
package web3

type BinanceClient struct {
    apiKey    string
    apiSecret string
    baseURL   string
}

func NewBinanceClient(apiKey, apiSecret string) *BinanceClient {
    return &BinanceClient{
        apiKey:    apiKey,
        apiSecret: apiSecret,
        baseURL:   "https://api.binance.com",
    }
}

func (b *BinanceClient) GetBalance(ctx context.Context, currency string) (string, error) {
    // 实现余额查询
}
```

2. **在配置中添加**

```go
// pkg/web3/config.go
type Config struct {
    // ... 现有配置
    BinanceAPIKey    string
    BinanceAPISecret string
}
```

3. **注册到管理器**

```go
// pkg/web3/config.go - InitializeClients()
if cfg.BinanceAPIKey != "" {
    binance := NewBinanceClient(cfg.BinanceAPIKey, cfg.BinanceAPISecret)
    GetExchangeManager().RegisterExchange(ExchangeBinance, binance)
}
```

## 参考链接

- [Coinbase Pro API 文档](https://docs.cloud.coinbase.com/exchange/reference/)
- [KuCoin API 文档](https://docs.kucoin.com/)
- [AideCMS 框架文档](./getting-started.md)
- [Web3 集成文档](./web3.md)

## 更新日志

### v1.0.0 (2024-01-xx)
- ✅ 初始版本
- ✅ 支持 Coinbase 和 KuCoin
- ✅ HTTP API 和 CLI 命令
- ✅ 多交易所查询功能

## 许可证

MIT License
