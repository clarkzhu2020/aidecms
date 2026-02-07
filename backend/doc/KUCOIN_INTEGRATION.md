# KuCoin 交易所集成文档

## 概述

本文档说明了如何使用AideCMS与KuCoin交易所进行交互，包括订单管理、账户查询和市场数据获取。

## 目录

- [配置](#配置)
- [API认证](#api认证)
- [数据库表结构](#数据库表结构)
- [API端点](#api端点)
- [使用示例](#使用示例)
- [错误处理](#错误处理)
- [注意事项](#注意事项)

## 配置

### 环境变量

在 `.env` 文件中配置以下环境变量：

```bash
# KuCoin 交易配置
KUCOIN_ENABLED=true                    # 是否启用KuCoin集成
KUCOIN_API_KEY=your_api_key           # KuCoin API Key
KUCOIN_API_SECRET=your_api_secret     # KuCoin API Secret
KUCOIN_PASSPHRASE=your_passphrase     # KuCoin API Passphrase
KUCOIN_SANDBOX=true                   # 是否使用沙盒环境
```

### 获取API凭证

1. 登录 [KuCoin官网](https://www.kucoin.com/)
2. 进入"API管理"页面
3. 创建新的API Key
4. 设置权限：
   - **General** - 查询账户信息
   - **Spot** - 现货交易
   - **Transfer** - 内部转账（可选）
5. 保存API Key、API Secret和Passphrase

### IP白名单（推荐）

为了安全起见，建议设置IP白名单限制API访问。

## API认证

KuCoin使用HMAC SHA256签名进行API认证。

### 认证流程

1. 生成时间戳（毫秒）
2. 构造签名字符串：`timestamp + method + endpoint + body`
3. 使用HMAC SHA256和API Secret进行签名
4. Base64编码签名结果
5. 在请求头中添加认证信息

### 请求头

| 请求头 | 描述 | 示例 |
|--------|------|------|
| `KC-API-KEY` | API密钥 | `60d6a5c7a5c7a5c7a5c7a5c7` |
| `KC-API-SIGN` | 签名 | `base64编码的签名` |
| `KC-API-TIMESTAMP` | 时间戳（毫秒） | `1629080342000` |
| `KC-API-PASSPHRASE` | Passphrase签名 | `base64编码的passphrase签名` |
| `KC-API-KEY-VERSION` | API密钥版本 | `2` |

## 数据库表结构

### KuCoin订单表 (kucoin_orders)

```go
type KuCoinOrder struct {
    ID              uint      `gorm:"primaryKey"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
    DeletedAt       gorm.DeletedAt
    OrderID         string    `gorm:"type:varchar(100);uniqueIndex;not null"`
    ClientOrderID   string    `gorm:"type:varchar(100);index"`
    Symbol          string    `gorm:"type:varchar(50);not null;index"`
    Side            string    `gorm:"type:varchar(10);not null"` // buy, sell
    Type            string    `gorm:"type:varchar(20);not null"` // limit, market, stop
    Price           string    `gorm:"type:decimal(30,18)"`
    Size            string    `gorm:"type:decimal(30,18);not null"`
    DealSize        string    `gorm:"type:decimal(30,18)"`
    DealFunds       string    `gorm:"type:decimal(30,18)"`
    Fee             string    `gorm:"type:decimal(30,18)"`
    FeeCurrency     string    `gorm:"type:varchar(20)"`
    StopPrice       string    `gorm:"type:decimal(30,18)"`
    TimeInForce     string    `gorm:"type:varchar(10)"` // GTC, IOC, FOK
    Status          string    `gorm:"type:varchar(20);index"` // open, done, match, canceled
    KuCoinCreatedAt int64
    KuCoinUpdatedAt int64
    Remark          string    `gorm:"type:text"`
}
```

### KuCoin成交记录表 (kucoin_trades)

```go
type KuCoinTrade struct {
    ID              uint      `gorm:"primaryKey"`
    CreatedAt       time.Time
    TradeID         string    `gorm:"type:varchar(100);uniqueIndex;not null"`
    OrderID         string    `gorm:"type:varchar(100);index;not null"`
    Symbol          string    `gorm:"type:varchar(50);not null;index"`
    Side            string    `gorm:"type:varchar(10);not null"` // buy, sell
    Price           string    `gorm:"type:decimal(30,18);not null"`
    Size            string    `gorm:"type:decimal(30,18);not null"`
    Fee             string    `gorm:"type:decimal(30,18)"`
    FeeCurrency     string    `gorm:"type:varchar(20)"`
    CounterOrderID  string    `gorm:"type:varchar(100)"`
    ForceTaker      bool      `gorm:"default:false"`
    Liquidity       string    `gorm:"type:varchar(20)"`
    KuCoinCreatedAt int64
}
```

### KuCoin账户表 (kucoin_accounts)

```go
type KuCoinAccount struct {
    ID             string    `gorm:"primaryKey"`
    CreatedAt      time.Time
    UpdatedAt      time.Time
    AccountID      string    `gorm:"type:varchar(100);uniqueIndex;not null"`
    Currency       string    `gorm:"type:varchar(20);not null;index"`
    Type           string    `gorm:"type:varchar(20);not null"` // main, trade, margin
    Balance        string    `gorm:"type:decimal(30,18)"`
    Available      string    `gorm:"type:decimal(30,18)"`
    Holds          string    `gorm:"type:decimal(30,18)"`
    LastSyncedAt   time.Time `gorm:"index"`
}
```

### KuCoin余额快照表 (kucoin_balance_snapshots)

```go
type KuCoinBalanceSnapshot struct {
    ID            string    `gorm:"primaryKey"`
    CreatedAt     time.Time
    SnapshotID    string    `gorm:"type:varchar(100);uniqueIndex;not null"`
    AccountID     string    `gorm:"type:varchar(100);index"`
    Currency      string    `gorm:"type:varchar(20);not null;index"`
    Balance       string    `gorm:"type:decimal(30,18)"`
    Available     string    `gorm:"type:decimal(30,18)"`
    Holds         string    `gorm:"type:decimal(30,18)"`
    SnapshotAt    time.Time `gorm:"not null;index"`
}
```

## API端点

### 服务器时间

#### 获取服务器时间
```
GET /api/kucoin/server-time
```

**响应示例：**
```json
{
  "message": "Server time retrieved successfully",
  "data": {
    "timestamp": 1629080342000,
    "datetime": "2021-08-15T12:45:42Z"
  }
}
```

### 市场数据

#### 获取交易对列表
```
GET /api/kucoin/symbols?market=BTC
```

**参数：**
- `market` (可选) - 市场类型，如 `BTC`, `ETH`

**响应示例：**
```json
{
  "message": "Symbols retrieved successfully",
  "data": [
    {
      "symbol": "BTC-USDT",
      "baseCurrency": "BTC",
      "quoteCurrency": "USDT",
      "baseMinSize": "0.000001",
      "quoteMinSize": "0.01",
      "enableTrading": true
    }
  ]
}
```

#### 获取Ticker
```
GET /api/kucoin/ticker?symbol=BTC-USDT
```

**参数：**
- `symbol` (必填) - 交易对，如 `BTC-USDT`

**响应示例：**
```json
{
  "message": "Ticker retrieved successfully",
  "data": {
    "symbol": "BTC-USDT",
    "buy": "45000.00",
    "sell": "45001.00",
    "changeRate": "0.025",
    "high": "45500.00",
    "low": "44500.00",
    "vol": "1234.56",
    "last": "45000.50"
  }
}
```

#### 获取订单簿
```
GET /api/kucoin/orderbook?symbol=BTC-USDT
```

**参数：**
- `symbol` (必填) - 交易对

**响应示例：**
```json
{
  "message": "Order book retrieved successfully",
  "data": {
    "sequence": "123456789",
    "bids": [
      ["45000.00", "1.5"],
      ["44999.00", "2.3"]
    ],
    "asks": [
      ["45001.00", "1.2"],
      ["45002.00", "3.4"]
    ]
  }
}
```

#### 获取市场成交记录
```
GET /api/kucoin/trades?symbol=BTC-USDT
```

**参数：**
- `symbol` (必填) - 交易对

**响应示例：**
```json
{
  "message": "Market trades retrieved successfully",
  "data": [
    {
      "sequence": "123456789",
      "price": "45000.00",
      "size": "0.5",
      "side": "buy",
      "time": 1629080342000,
      "tradeId": "123456"
    }
  ]
}
```

#### 获取K线数据
```
GET /api/kucoin/klines?symbol=BTC-USDT&type=1hour
```

**参数：**
- `symbol` (必填) - 交易对
- `type` (可选) - K线类型，默认 `1hour`，支持：`1min`, `5min`, `15min`, `30min`, `1hour`, `4hour`, `1day`, `1week`

**响应示例：**
```json
{
  "message": "Klines retrieved successfully",
  "data": [
    {
      "Time": 1629080342000,
      "Open": "44950.00",
      "Close": "45000.00",
      "High": "45050.00",
      "Low": "44900.00",
      "Volume": "123.45",
      "Turnover": "5543210.00"
    }
  ]
}
```

#### 获取24小时统计数据
```
GET /api/kucoin/24h-stats?symbol=BTC-USDT
```

**参数：**
- `symbol` (必填) - 交易对

### 账户管理

#### 获取账户列表
```
GET /api/kucoin/accounts?currency=USDT&type=trade
```

**参数：**
- `currency` (可选) - 币种
- `type` (可选) - 账户类型：`main`, `trade`, `margin`

**响应示例：**
```json
{
  "message": "Accounts retrieved successfully",
  "data": [
    {
      "id": "123456789",
      "currency": "USDT",
      "type": "trade",
      "balance": "10000.00",
      "available": "9500.00",
      "holds": "500.00"
    }
  ]
}
```

#### 获取账户详情
```
GET /api/kucoin/accounts/:accountId
```

#### 同步账户信息
```
POST /api/kucoin/accounts/sync
```

**功能：** 同步KuCoin账户信息到本地数据库

**响应示例：**
```json
{
  "message": "Accounts synced successfully",
  "data": {
    "synced_count": 10,
    "total_count": 10
  }
}
```

#### 创建余额快照
```
POST /api/kucoin/balance/snapshot
```

**功能：** 创建当前余额快照，用于历史记录和追踪

**响应示例：**
```json
{
  "message": "Balance snapshot created successfully",
  "data": {
    "snapshot_id": "uuid",
    "snapshot_at": "2024-01-01T00:00:00Z",
    "snapshot_count": 15
  }
}
```

### 订单管理

#### 创建订单
```
POST /api/kucoin/orders
```

**请求体：**
```json
{
  "clientOid": "unique-client-order-id",
  "symbol": "BTC-USDT",
  "side": "buy",
  "type": "limit",
  "price": "45000.00",
  "size": "0.001",
  "timeInForce": "GTC",
  "remark": "Test order"
}
```

**参数说明：**
- `clientOid` (必填) - 客户端订单ID，必须唯一
- `symbol` (必填) - 交易对
- `side` (必填) - 买卖方向：`buy`, `sell`
- `type` (必填) - 订单类型：`limit`, `market`, `stop`
- `price` (限价订单必填) - 价格
- `size` (必填) - 数量
- `timeInForce` (可选) - 订单有效期：`GTC`, `IOC`, `FOK`
- `stopPrice` (止损订单必填) - 止损价格
- `remark` (可选) - 备注

**响应示例：**
```json
{
  "message": "Order created successfully",
  "data": {
    "orderId": "1234567890",
    "clientOid": "unique-client-order-id",
    "symbol": "BTC-USDT",
    "side": "buy",
    "type": "limit",
    "price": "45000.00",
    "size": "0.001",
    "createdAt": 1629080342000
  }
}
```

#### 取消订单
```
DELETE /api/kucoin/orders/:orderId
```

#### 获取订单详情
```
GET /api/kucoin/orders/:orderId
```

#### 获取未成交订单
```
GET /api/kucoin/orders?symbol=BTC-USDT
```

**参数：**
- `symbol` (可选) - 交易对

#### 获取已成交订单
```
GET /api/kucoin/orders/closed?symbol=BTC-USDT&status=done&currentPage=1&pageSize=20
```

**参数：**
- `symbol` (可选) - 交易对
- `status` (可选) - 订单状态
- `currentPage` (可选) - 当前页码
- `pageSize` (可选) - 每页数量

## 使用示例

### 创建限价买单

```bash
curl -X POST http://localhost:8888/api/kucoin/orders \
  -H "Content-Type: application/json" \
  -d '{
    "clientOid": "buy-btc-001",
    "symbol": "BTC-USDT",
    "side": "buy",
    "type": "limit",
    "price": "45000.00",
    "size": "0.001",
    "timeInForce": "GTC"
  }'
```

### 创建市价卖单

```bash
curl -X POST http://localhost:8888/api/kucoin/orders \
  -H "Content-Type: application/json" \
  -d '{
    "clientOid": "sell-btc-001",
    "symbol": "BTC-USDT",
    "side": "sell",
    "type": "market",
    "size": "0.001"
  }'
```

### 获取账户余额

```bash
curl http://localhost:8888/api/kucoin/accounts?currency=USDT&type=trade
```

### 获取BTC-USDT的Ticker

```bash
curl http://localhost:8888/api/kucoin/ticker?symbol=BTC-USDT
```

### 同步账户到数据库

```bash
curl -X POST http://localhost:8888/api/kucoin/accounts/sync
```

## 错误处理

### 错误响应格式

```json
{
  "error": "Error message description"
}
```

### 常见错误

| 错误代码 | 描述 | 解决方案 |
|----------|------|----------|
| `KuCoin client not initialized` | KuCoin客户端未初始化 | 检查环境变量配置，确保 `KUCOIN_ENABLED=true` |
| `Invalid request` | 请求参数无效 | 检查请求参数格式和必填字段 |
| `Failed to create order` | 创建订单失败 | 检查账户余额是否充足，参数是否正确 |
| `API error` | KuCoin API返回错误 | 查看具体错误信息，可能是权限或限流问题 |

## 注意事项

### 安全性

1. **保护API凭证**：永远不要在客户端代码中暴露API Secret和Passphrase
2. **使用IP白名单**：限制API访问的IP地址
3. **最小权限原则**：只为API Key分配必要的权限
4. **定期轮换密钥**：定期更换API Key

### 限流规则

KuCoin API有严格的限流规则，不同VIP级别有不同的请求限制：
- 免费账户：每30秒有限制
- VIP1-9：限流更高

建议：
- 合理设置请求频率
- 使用批量操作代替多次单次操作
- 实现请求重试机制

### 时间同步

KuCoin API要求服务器时间与请求时间戳误差在30秒内：
- 系统会自动同步服务器时间
- 如果误差过大，会返回错误

### 订单状态

订单状态说明：
- `open` - 未成交
- `done` - 已完成
- `match` - 部分成交
- `canceled` - 已取消
- `invalid` - 无效订单

### 测试环境

使用沙盒环境测试：
```bash
KUCOIN_SANDBOX=true
```

沙盒环境API地址：`https://openapi-sandbox.kucoin.com`

生产环境API地址：`https://api.kucoin.com`

## 数据库迁移

运行数据库迁移创建KuCoin相关表：

```go
import "aidecms/database/migrations"

// 在初始化数据库后运行
migrations.CreateKuCoinTables(db)
```

## 支持

如有问题，请参考：
- [KuCoin API文档](https://www.kucoin.com/docs-new/)
- [KuCoin开发者社区](https://community.kucoin.com/)
