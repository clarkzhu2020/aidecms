# MoonPay 支付集成文档

## 概述

AideCMS 已集成 MoonPay 支付服务，支持法币购买加密货币的功能。MoonPay 是一个专业的加密货币入金服务，提供简单、快速的加密货币购买体验。

## 功能特性

- ✅ 加密货币购买交易创建
- ✅ 实时报价查询
- ✅ Widget URL 生成（前端集成）
- ✅ 交易状态查询
- ✅ Webhook 事件处理
- ✅ 沙盒环境支持
- ✅ 完整的交易记录管理

## 快速开始

### 1. 获取 MoonPay API 密钥

1. 访问 [MoonPay Dashboard](https://dashboard.moonpay.com/) 注册账户
2. 完成 KYB（了解您的业务）流程
3. 在 "Developers" 选项卡中获取 API 密钥
4. 配置 Webhook 端点 URL

### 2. 环境变量配置

在 `.env` 文件中添加以下配置：

```bash
# MoonPay 支付配置
MOONPAY_ENABLED=true
MOONPAY_API_KEY=your_api_key_here
MOONPAY_SECRET_KEY=your_secret_key_here
MOONPAY_WEBHOOK_KEY=your_webhook_key_here
MOONPAY_SANDBOX=true
```

**配置说明：**
- `MOONPAY_ENABLED`: 是否启用 MoonPay（true/false）
- `MOONPAY_API_KEY`: MoonPay API 密钥（必填）
- `MOONPAY_SECRET_KEY`: MoonPay 密钥（必填，用于服务器端API调用）
- `MOONPAY_WEBHOOK_KEY`: Webhook 签名验证密钥（可选）
- `MOONPAY_SANDBOX`: 是否使用沙盒环境（true/false）

### 3. 数据库迁移

运行数据库迁移脚本以创建 MoonPay 相关表：

```bash
cd backend
go run database/migrations/create_moonpay_tables.go
```

这将创建以下表：
- `moonpay_transactions`: MoonPay 交易记录
- `moonpay_webhooks`: MoonPay Webhook 事件记录

## API 端点

### 创建交易

**POST** `/api/moonpay/transactions`

创建一个新的 MoonPay 加密货币购买交易。

**请求体：**
```json
{
  "baseCurrencyAmount": 100.00,
  "baseCurrencyCode": "USD",
  "currencyCode": "BTC",
  "walletAddress": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
  "externalId": "ORDER-12345",
  "redirectUrl": "http://localhost:8888/api/moonpay/success",
  "lockAmount": true,
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe"
}
```

**响应示例：**
```json
{
  "code": 201,
  "message": "Transaction created successfully",
  "data": {
    "transaction_id": 1,
    "moonpay_id": "mp_abc123xyz",
    "widget_url": "https://buy-sandbox.moonpay.com?apiKey=xxx&currencyCode=BTC...",
    "status": "pending",
    "amount": 100.00,
    "currency": "USD",
    "currency_code": "BTC",
    "wallet_address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
  }
}
```

### 获取交易详情

**GET** `/api/moonpay/transactions/:id`

获取指定交易的详细信息。

**响应示例：**
```json
{
  "code": 200,
  "message": "Transaction fetched successfully",
  "data": {
    "id": 1,
    "transaction_id": "mp_abc123xyz",
    "external_id": "ORDER-12345",
    "status": "completed",
    "base_currency_amount": 100.00,
    "base_currency_code": "USD",
    "quote_currency_amount": 0.0025,
    "quote_currency_code": "BTC",
    "wallet_address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
  }
}
```

### 获取交易列表

**GET** `/api/moonpay/transactions`

获取交易列表，支持分页和状态筛选。

**查询参数：**
- `status`: 交易状态（pending/waiting_payment/pending_approval/completed/failed）
- `page`: 页码（默认：1）
- `limit`: 每页数量（默认：20，最大：100）

**响应示例：**
```json
{
  "code": 200,
  "message": "Transactions fetched successfully",
  "data": {
    "data": [...],
    "page": 1,
    "limit": 20,
    "total": 50
  }
}
```

### 获取购买报价

**POST** `/api/moonpay/quote`

获取加密货币购买的实时报价。

**请求体：**
```json
{
  "baseCurrencyAmount": 100.00,
  "baseCurrencyCode": "USD",
  "currencyCode": "BTC"
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "Quote fetched successfully",
  "data": {
    "baseCurrencyAmount": 100.00,
    "baseCurrencyCode": "USD",
    "quoteCurrencyAmount": 0.0025,
    "quoteCurrencyCode": "BTC",
    "totalFeeAmount": 5.00,
    "networkFeeAmount": 0.0001
  }
}
```

### 生成 Widget URL

**POST** `/api/moonpay/widget-url`

生成 MoonPay Widget URL 用于前端集成。

**请求体：** 同创建交易请求。

**响应示例：**
```json
{
  "code": 200,
  "message": "Widget URL generated successfully",
  "data": {
    "widget_url": "https://buy-sandbox.moonpay.com?apiKey=xxx&currencyCode=BTC&baseCurrencyAmount=100&..."
  }
}
```

### Webhook 处理

**POST** `/api/moonpay/webhook`

接收并处理 MoonPay 的 Webhook 通知。

**请求头：**
- `X-MoonPay-Signature`: Webhook 签名（用于验证请求来源）

**请求体示例：**
```json
{
  "id": "evt_abc123",
  "type": "transaction.completed",
  "transactionId": "mp_abc123xyz",
  "summary": "Transaction completed successfully",
  "createdAt": "2024-01-01T12:00:00Z"
}
```

## 交易状态

MoonPay 交易支持以下状态：

| 状态 | 说明 |
|------|------|
| `pending` | 交易已创建，等待处理 |
| `waiting_payment` | 等待用户支付 |
| `pending_approval` | 等待审批 |
| `completed` | 交易已完成 |
| `failed` | 交易失败 |
| `cancelled` | 交易已取消 |

## Webhook 事件

MoonPay 会发送以下 Webhook 事件：

| 事件类型 | 说明 |
|----------|------|
| `transaction.created` | 交易已创建 |
| `transaction.completed` | 交易已完成 |
| `transaction.failed` | 交易失败 |
| `transaction.cancelled` | 交易已取消 |

## 前端集成

### Widget 方式（推荐）

1. 调用 API 创建交易或生成 Widget URL
2. 将返回的 `widget_url` 在浏览器中打开
3. 用户在 MoonPay 的支付页面完成支付
4. 支付完成后跳转到您设置的 `redirect_url`

示例：
```javascript
// 创建交易
const response = await fetch('/api/moonpay/transactions', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    baseCurrencyAmount: 100,
    baseCurrencyCode: 'USD',
    currencyCode: 'BTC',
    walletAddress: '0x...',
    externalId: 'ORDER-123'
  })
});

const data = await response.json();

// 打开 MoonPay Widget
window.location.href = data.data.widget_url;
```

### SDK 方式

MoonPay 提供 JavaScript SDK，可以在您的应用中嵌入 Widget：

```html
<script>
  const moonpay = new MoonPayWebSdk.WebButton({
    debug: true,
    flow: 'buy',
    environment: 'sandbox',
    params: {
      apiKey: 'your_api_key',
      currencyCode: 'BTC',
      walletAddress: '0x...'
    }
  });

  document.body.appendChild(moonpay);
</script>
```

更多 SDK 信息请参考：[MoonPay 官方文档](https://dev.moonpay.com/docs/on-ramp-overview)

## 测试

### 沙盒环境测试

使用沙盒环境进行测试：

```bash
MOONPAY_SANDBOX=true
```

沙盒环境 API 地址：
- Buy Widget: `https://buy-sandbox.moonpay.com`
- API: `https://api-sandbox.moonpay.com`

### 测试 Webhook

使用 ngrok 或类似工具将本地服务器暴露到公网：

```bash
ngrok http 8888
```

然后在 MoonPay Dashboard 中配置 Webhook URL：
```
https://your-ngrok-url.com/api/moonpay/webhook
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
| `401` | 未授权（API密钥无效） |
| `404` | 资源未找到 |
| `500` | 服务器内部错误 |

## 安全建议

1. **保护 API 密钥**：不要将 API 密钥提交到版本控制系统
2. **使用 Webhook 签名验证**：验证所有 Webhook 请求的来源
3. **HTTPS**：生产环境必须使用 HTTPS
4. **环境隔离**：沙盒和生产环境使用不同的 API 密钥
5. **日志记录**：记录所有交易和 Webhook 事件

## 架构说明

### 目录结构

```
backend/
├── pkg/moonpay/              # MoonPay 核心包
│   ├── client.go            # 客户端初始化
│   ├── http_client.go       # HTTP 客户端
│   └── transaction.go       # 交易服务
├── internal/app/models/     # 数据模型
│   └── moonpay_transaction.go
├── app/Http/Controllers/    # 控制器
│   └── MoonPayController.go
├── config/                  # 配置
│   └── moonpay.go
├── database/migrations/     # 数据库迁移
│   └── create_moonpay_tables.go
└── doc/                     # 文档
    └── MOONPAY_INTEGRATION.md
```

### 数据流程

```
1. 用户发起购买请求
   ↓
2. 后端创建交易并生成 Widget URL
   ↓
3. 用户在 MoonPay 页面完成支付
   ↓
4. MoonPay 发送 Webhook 通知
   ↓
5. 后端更新交易状态
   ↓
6. 系统执行业务逻辑（如发放资产）
```

## 支持的加密货币

MoonPay 支持 80+ 种加密货币，包括但不限于：

- BTC (Bitcoin)
- ETH (Ethereum)
- USDT, USDC (稳定币)
- BNB (Binance Coin)
- SOL (Solana)
- MATIC (Polygon)
- 等等...

完整的货币列表请参考：[MoonPay 官方文档](https://dev.moonpay.com/docs/currencies)

## 参考资源

- [MoonPay 官方文档](https://dev.moonpay.com)
- [MoonPay API Reference](https://dev.moonpay.com/reference)
- [MoonPay Dashboard](https://dashboard.moonpay.com/)
- [MoonPay Support](https://support.moonpay.com/)

## 故障排除

### 问题：MoonPay 服务未初始化

**解决方案：**
- 检查 `.env` 文件中 `MOONPAY_API_KEY` 是否已配置
- 确认环境变量已正确加载

### 问题：Webhook 未接收

**解决方案：**
- 检查 MoonPay Dashboard 中的 Webhook URL 配置
- 确认服务器可以从公网访问
- 查看 MoonPay Dashboard 的 Webhook 日志

### 问题：交易状态未更新

**解决方案：**
- 检查 Webhook 端点是否正常工作
- 查看应用日志确认 Webhook 是否被正确处理
- 尝试手动调用 MoonPay API 查询交易状态

## 联系支持

如有问题或建议，请联系：
- MoonPay 官方支持：https://support.moonpay.com/
- 项目问题反馈：GitHub Issues
