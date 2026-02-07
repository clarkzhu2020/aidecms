# Coinbase 支付和交易集成文档

## 概述

AideCMS 已集成 Coinbase 的 Payment Link API 和 Trade API，支持：
- **Payment Link API**: 创建加密货币支付链接，接受加密货币付款
- **Trade API**: 执行加密货币交易，获取市场数据

## 功能特性

### Payment Link API
- ✅ 创建和管理支付链接
- ✅ 自定义支付金额和货币
- ✅ 设置回调URL
- ✅ 支付状态跟踪
- ✅ Webhook 事件处理

### Trade API
- ✅ 创建市价单和限价单
- ✅ 查询订单状态
- ✅ 取消订单
- ✅ 获取产品列表
- ✅ 获取实时行情
- ✅ 查询账户余额

## 快速开始

### 1. 获取 Coinbase API 密钥

1. 访问 [Coinbase Business](https://www.coinbase.com/business) 注册账户
2. 完成企业验证和 KYB 流程
3. 在 Settings > API 中创建 API Key
4. 记录 API Key 和 API Secret
5. 配置 Webhook 端点 URL

### 2. 环境变量配置

在 `.env` 文件中添加以下配置：

```bash
# Coinbase 支付配置
COINBASE_ENABLED=true
COINBASE_API_KEY=your_api_key_here
COINBASE_API_SECRET=your_api_secret_here
COINBASE_WEBHOOK_KEY=your_webhook_key_here
COINBASE_ACCOUNT_ID=your_account_id_here
COINBASE_SANDBOX=true
```

**配置说明：**
- `COINBASE_ENABLED`: 是否启用 Coinbase（true/false）
- `COINBASE_API_KEY`: Coinbase API Key（必填）
- `COINBASE_API_SECRET`: Coinbase API Secret（必填）
- `COINBASE_WEBHOOK_KEY`: Webhook 签名验证密钥（可选）
- `COINBASE_ACCOUNT_ID`: Coinbase 账户ID（必填）
- `COINBASE_SANDBOX`: 是否使用沙盒环境（true/false）

### 3. 数据库迁移

运行数据库迁移脚本以创建 Coinbase 相关表：

```bash
cd backend
go run database/migrations/create_coinbase_tables.go
```

这将创建以下表：
- `coinbase_payment_links`: Coinbase 支付链接记录
- `coinbase_orders`: Coinbase 交易订单记录
- `coinbase_webhooks`: Coinbase Webhook 事件记录

## Payment Link API

### 创建支付链接

**POST** `/api/coinbase/payment-links`

创建一个新的 Coinbase 加密货币支付链接。

**请求体：**
```json
{
  "amount": "100.00",
  "currency": "USD",
  "title": "Order Payment",
  "description": "Payment for order #12345",
  "redirectUrl": "http://localhost:8888/api/coinbase/success",
  "cancelUrl": "http://localhost:8888/api/coinbase/cancel",
  "name": "John Doe",
  "email": "john@example.com",
  "metadata": {
    "orderId": "12345"
  }
}
```

**响应示例：**
```json
{
  "code": 201,
  "message": "Payment link created successfully",
  "data": {
    "id": 1,
    "link_id": "CL_abc123xyz",
    "payment_url": "https://commerce.coinbase.com/checkout/abc123",
    "amount": "100.00",
    "currency": "USD",
    "status": "active",
    "payment_status": "pending"
  }
}
```

### 获取支付链接详情

**GET** `/api/coinbase/payment-links/:id`

获取指定支付链接的详细信息。

**响应示例：**
```json
{
  "code": 200,
  "message": "Payment link fetched successfully",
  "data": {
    "id": 1,
    "link_id": "CL_abc123xyz",
    "amount": "100.00",
    "currency": "USD",
    "payment_url": "https://commerce.coinbase.com/checkout/abc123",
    "status": "active",
    "payment_status": "completed"
  }
}
```

### 列出支付链接

**GET** `/api/coinbase/payment-links`

获取支付链接列表，支持分页和状态筛选。

**查询参数：**
- `status`: 支付链接状态
- `page`: 页码（默认：1）
- `limit`: 每页数量（默认：20）

### 删除支付链接

**DELETE** `/api/coinbase/payment-links/:id`

删除指定的支付链接。

## Trade API

### 创建交易订单

**POST** `/api/coinbase/orders`

创建一个新的 Coinbase 交易订单。

**请求体：**
```json
{
  "product_id": "BTC-USD",
  "side": "buy",
  "order_type": "limit",
  "size": "0.001",
  "limit_price": "50000.00",
  "time_in_force": "GTC",
  "client_order_id": "ORDER-12345"
}
```

**参数说明：**
- `product_id`: 交易对，如 "BTC-USD", "ETH-USD" 等
- `side`: 交易方向，"buy"（买入）或 "sell"（卖出）
- `order_type`: 订单类型
  - `market`: 市价单
  - `limit`: 限价单
  - `stop`: 止损单
- `size`: 交易数量
- `funds`: 使用金额（仅限市价买单）
- `limit_price`: 限价（仅限限价单）
- `stop_price`: 止损价格（仅限止损单）
- `time_in_force`: 有效期类型
  - `GTC`: Good Till Canceled（取消前有效）
  - `IOC`: Immediate or Cancel（立即成交或取消）
  - `FOK`: Fill or Kill（全部成交或取消）
  - `GTD`: Good Till Date（指定日期前有效）
- `post_only`: 仅作为 maker（bool）
- `self_trade_prevention`: 自成交预防
  - `dc`: Decrease and Cancel（减少并取消）
  - `co`: Cancel Oldest（取消最旧的）
  - `cn`: Cancel Newest（取消最新的）
  - `cb`: Cancel Both（取消两者）

**响应示例：**
```json
{
  "code": 201,
  "message": "Order created successfully",
  "data": {
    "id": 1,
    "order_id": "ORD_abc123xyz",
    "product_id": "BTC-USD",
    "side": "buy",
    "order_type": "limit",
    "size": "0.001",
    "limit_price": "50000.00",
    "status": "open",
    "filled_size": "0",
    "settled": false
  }
}
```

### 获取订单详情

**GET** `/api/coinbase/orders/:id`

获取指定订单的详细信息。

**响应示例：**
```json
{
  "code": 200,
  "message": "Order fetched successfully",
  "data": {
    "id": 1,
    "order_id": "ORD_abc123xyz",
    "product_id": "BTC-USD",
    "side": "buy",
    "status": "filled",
    "filled_size": "0.001",
    "average_fill_price": "49900.00",
    "fill_fees": "2.50",
    "settled": true
  }
}
```

### 列出订单

**GET** `/api/coinbase/orders`

获取订单列表，支持分页和筛选。

**查询参数：**
- `productId`: 产品ID
- `status`: 订单状态
- `page`: 页码（默认：1）
- `limit`: 每页数量（默认：20）

### 取消订单

**POST** `/api/coinbase/orders/:id/cancel`

取消指定的交易订单。

### 获取产品列表

**GET** `/api/coinbase/products`

获取所有可交易的产品列表。

**响应示例：**
```json
{
  "code": 200,
  "message": "Products fetched successfully",
  "data": [
    {
      "product_id": "BTC-USD",
      "base_currency": "BTC",
      "quote_currency": "USD",
      "display_name": "BTC/USD",
      "status": "online",
      "base_min_size": "0.00001",
      "base_max_size": "1000",
      "is_market_orderable": true,
      "is_limit_orderable": true
    }
  ]
}
```

### 获取产品详情

**GET** `/api/coinbase/products/:productId`

获取指定产品的详细信息。

### 获取产品行情

**GET** `/api/coinbase/products/:productId/ticker`

获取指定产品的实时行情信息。

**响应示例：**
```json
{
  "code": 200,
  "message": "Ticker fetched successfully",
  "data": {
    "product_id": "BTC-USD",
    "price": "49950.00",
    "open_24h": "48500.00",
    "high_24h": "50200.00",
    "low_24h": "48000.00",
    "volume_24h": "1234.56",
    "percentage_change": "2.99"
  }
}
```

## Webhook 事件

### Payment Link Webhook 事件

Coinbase 会发送以下 Webhook 事件：

| 事件类型 | 说明 |
|----------|------|
| `payment_link.completed` | 支付已完成 |
| `payment_link.created` | 支付链接已创建 |
| `payment_link.cancelled` | 支付已取消 |

### Trade Webhook 事件

| 事件类型 | 说明 |
|----------|------|
| `order.filled` | 订单已成交 |
| `order.cancelled` | 订单已取消 |
| `order.rejected` | 订单被拒绝 |

### Webhook 处理

所有 Webhook 事件发送到：`POST /api/coinbase/webhook`

**请求头：**
- `X-CC-Webhook-Signature`: Webhook 签名（用于验证请求来源）
- `X-CC-Webhook-Timestamp`: 时间戳

**示例请求：**
```json
{
  "id": "evt_abc123",
  "type": "payment_link.completed",
  "payment_link_id": "CL_abc123xyz",
  "summary": "Payment completed successfully",
  "created_at": "2024-01-01T12:00:00Z"
}
```

## 订单状态

Coinbase 订单支持以下状态：

| 状态 | 说明 |
|------|------|
| `open` | 订单已创建，等待成交 |
| `filled` | 订单已完全成交 |
| `partially_filled` | 订单部分成交 |
| `cancelled` | 订单已取消 |
| `rejected` | 订单被拒绝 |
| `expired` | 订单已过期 |

## JWT 认证

Coinbase API 使用 JWT 进行认证。系统会自动生成 JWT Token。

### JWT 结构

```go
claims := jwt.MapClaims{
    "sub": apiKey,      // 主体
    "iss": apiKey,      // 签发者
    "aud": "https://business.coinbase.com", // 受众
    "exp": 60,         // 过期时间（秒）
    "nbf": 0,          // 生效时间
    "iat": 0,          // 签发时间
}
```

### 使用 JWT 请求

```bash
curl -X GET https://business.coinbase.com/api/v3/brokerage/products \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "CB-ACCESS-KEY: YOUR_API_KEY"
```

## 支持的交易对

Coinbase 支持以下主要交易对：

### 法币交易对
- BTC-USD, BTC-EUR
- ETH-USD, ETH-EUR
- USDT-USD, USDC-USD
- SOL-USD, SOL-USDT
- MATIC-USD, MATIC-USDT
- 等等...

### 加密货币交易对
- BTC-ETH, ETH-BTC
- 等等...

完整的产品列表请调用 `/api/coinbase/products` 获取。

## 测试

### 沙盒环境测试

使用沙盒环境进行测试：

```bash
COINBASE_SANDBOX=true
```

沙盒环境 API 地址：
- Payment Links: `https://business-sandbox.coinbase.com`
- Trade API: `https://api-sandbox.coinbase.com`

### 测试 Webhook

使用 ngrok 或类似工具将本地服务器暴露到公网：

```bash
ngrok http 8888
```

然后在 Coinbase Dashboard 中配置 Webhook URL：
```
https://your-ngrok-url.com/api/coinbase/webhook
```

## 错误处理

API 返回标准错误格式：

```json
{
  "code": 400,
  "message": "Invalid request data"
}
```

常见错误：

| 错误代码 | 说明 |
|----------|------|
| `400` | 请求参数错误 |
| `401` | 未授权（JWT Token无效或过期） |
| `403` | 权限不足 |
| `404` | 资源未找到 |
| `429` | 请求过于频繁 |
| `500` | 服务器内部错误 |

## 安全建议

1. **保护 API 密钥**：不要将 API 密钥提交到版本控制系统
2. **使用 Webhook 签名验证**：验证所有 Webhook 请求的来源
3. **HTTPS**：生产环境必须使用 HTTPS
4. **环境隔离**：沙盒和生产环境使用不同的 API 密钥
5. **日志记录**：记录所有交易和 Webhook 事件
6. **JWT Token 管理**：及时刷新过期的 Token

## 架构说明

### 目录结构

```
backend/
├── pkg/coinbase/              # Coinbase 核心包
│   ├── client.go            # 客户端初始化和JWT生成
│   ├── http_client.go       # HTTP 客户端
│   ├── payment_link.go      # Payment Link API服务
│   └── trade.go            # Trade API服务
├── internal/app/models/     # 数据模型
│   └── coinbase_transaction.go
├── app/Http/Controllers/    # 控制器
│   └── CoinbaseController.go
├── config/                  # 配置
│   └── coinbase.go
├── database/migrations/     # 数据库迁移
│   └── create_coinbase_tables.go
└── doc/                     # 文档
    └── COINBASE_INTEGRATION.md
```

### 数据流程

#### Payment Link 流程
```
1. 用户请求创建支付链接
   ↓
2. 后端调用 Coinbase API 创建链接
   ↓
3. 返回支付链接 URL 给前端
   ↓
4. 用户访问 Coinbase 支付页面
   ↓
5. 用户完成支付
   ↓
6. Coinbase 发送 Webhook 通知
   ↓
7. 后端更新支付状态
   ↓
8. 系统执行业务逻辑
```

#### Trade 流程
```
1. 用户发起交易请求
   ↓
2. 后端生成 JWT Token
   ↓
3. 调用 Coinbase Trade API 创建订单
   ↓
4. 保存订单记录到数据库
   ↓
5. 监听订单状态变化
   ↓
6. 订单成交后更新状态
   ↓
7. 系统执行业务逻辑
```

## 参考资源

- [Coinbase Business API 文档](https://docs.cdp.coinbase.com/)
- [Coinbase Payment Link API](https://docs.cdp.coinbase.com/coinbase-business/payment-link-apis/overview)
- [Coinbase Trade API](https://docs.cdp.coinbase.com/docs/trade-api-beta)
- [Coinbase Dashboard](https://dashboard.coinbase.com/)
- [Coinbase Support](https://help.coinbase.com/)

## 故障排除

### 问题：Coinbase 服务未初始化

**解决方案：**
- 检查 `.env` 文件中 `COINBASE_API_KEY` 是否已配置
- 确认环境变量已正确加载

### 问题：JWT Token 生成失败

**解决方案：**
- 检查 API Secret 是否正确
- 确认 API Key 没有过期

### 问题：订单创建失败

**解决方案：**
- 检查交易对是否正确（如 "BTC-USD"）
- 确认账户余额充足
- 检查订单参数是否符合规则

### 问题：Webhook 未接收

**解决方案：**
- 检查 Coinbase Dashboard 中的 Webhook URL 配置
- 确认服务器可以从公网访问
- 查看 Coinbase Dashboard 的 Webhook 日志

### 问题：订单状态未更新

**解决方案：**
- 检查 Webhook 端点是否正常工作
- 查看应用日志确认 Webhook 是否被正确处理
- 尝试手动调用 Coinbase API 查询订单状态

## 使用示例

### 前端集成示例

#### 创建支付链接
```javascript
// 创建支付链接
const response = await fetch('/api/coinbase/payment-links', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    amount: '100.00',
    currency: 'USD',
    title: 'Order Payment',
    description: 'Payment for order #12345',
    name: 'John Doe',
    email: 'john@example.com'
  })
});

const data = await response.json();

// 打开支付链接
window.open(data.data.payment_url, '_blank');
```

#### 创建交易订单
```javascript
// 创建买单
const response = await fetch('/api/coinbase/orders', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    product_id: 'BTC-USD',
    side: 'buy',
    order_type: 'limit',
    size: '0.001',
    limit_price: '50000.00',
    time_in_force: 'GTC'
  })
});

const data = await response.json();
console.log('Order created:', data.data);
```

#### 获取实时行情
```javascript
// 获取 BTC/USD 行情
const response = await fetch('/api/coinbase/products/BTC-USD/ticker');
const data = await response.json();

console.log('Current price:', data.data.price);
console.log('24h change:', data.data.percentage_change);
```

## 联系支持

如有问题或建议，请联系：
- Coinbase 官方支持：https://help.coinbase.com/
- 项目问题反馈：GitHub Issues
