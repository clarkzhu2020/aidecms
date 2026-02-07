# Stripe 集成文档

## 概述

本系统已集成 Stripe 支付功能,提供完整的支付、退款和 Webhook 处理能力。

## 功能特性

- ✅ 创建支付意图 (Payment Intent)
- ✅ 确认支付
- ✅ 退款处理
- ✅ 支付状态查询
- ✅ Webhook 通知处理
- ✅ 客户管理
- ✅ 支付历史记录
- ✅ 测试/生产环境支持

## 安装步骤

### 1. 安装依赖

```bash
cd backend
go mod tidy
```

### 2. 配置环境变量

在 `.env` 文件中添加以下配置:

```env
# Stripe 支付配置
STRIPE_ENABLED=true
STRIPE_SECRET_KEY=sk_test_your_secret_key_here
STRIPE_PUBLISHABLE_KEY=pk_test_your_publishable_key_here
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret_here
STRIPE_MODE=test
```

### 3. 获取 Stripe 凭据

1. 登录 [Stripe Dashboard](https://dashboard.stripe.com/)
2. 进入 **Developers** > **API keys**
3. 获取 Secret Key 和 Publishable Key
4. 配置 Webhook (可选)

## API 接口

### 1. 创建支付意图

```http
POST /api/stripe/payments/intent
Content-Type: application/json
Authorization: Bearer {token}

{
  "order_id": "ORDER-123456",
  "amount": 1000,
  "currency": "usd",
  "description": "商品描述",
  "customer_email": "customer@example.com",
  "payment_method_types": ["card"],
  "metadata": {
    "custom_field": "value"
  }
}
```

响应示例:

```json
{
  "code": 201,
  "message": "Payment intent created successfully",
  "data": {
    "payment_id": 1,
    "order_id": "ORDER-123456",
    "payment_intent_id": "pi_1234567890",
    "client_secret": "pi_1234567890_secret_abc123",
    "amount": 1000,
    "currency": "usd",
    "status": "requires_payment_method"
  }
}
```

### 2. 确认支付

```http
POST /api/stripe/payments/{paymentIntentID}/confirm
Content-Type: application/json

{
  "payment_method_id": "pm_1234567890"
}
```

### 3. 获取支付详情

```http
GET /api/stripe/payments/{id}
Authorization: Bearer {token}
```

### 4. 退款

```http
POST /api/stripe/payments/{id}/refund
Content-Type: application/json

{
  "amount": 500,  // 可选,不填则为全额退款
  "reason": "requested_by_customer"
}
```

### 5. 获取支付列表

```http
GET /api/stripe/payments?status=succeeded&page=1&limit=20
Authorization: Bearer {token}
```

### 6. Webhook 端点

```http
POST /api/stripe/webhook
Content-Type: application/json
Stripe-Signature: t=1234567890,v1=abc123...
```

## 支付流程

### 完整支付流程

1. **创建支付意图**
   - 前端调用 `/api/stripe/payments/intent` 创建 PaymentIntent
   - 返回 `client_secret` 供前端使用

2. **客户端确认支付**
   - 使用 Stripe Elements 或 Stripe.js 前端SDK
   - 使用返回的 `client_secret` 确认支付

3. **支付完成通知**
   - Stripe 发送 Webhook 到 `/api/stripe/webhook`
   - 系统自动更新支付状态

### 前端集成示例

```html
<script src="https://js.stripe.com/v3/"></script>
<script>
  const stripe = Stripe('pk_test_your_publishable_key');

  // 使用返回的 client_secret
  const { error, paymentIntent } = await stripe.confirmCardPayment(clientSecret, {
    payment_method: {
      card: cardElement,
      billing_details: { name: 'John Doe' }
    }
  });

  if (error) {
    console.error('Payment failed:', error);
  } else {
    console.log('Payment succeeded:', paymentIntent);
  }
</script>
```

## 支付状态

| Stripe状态 | 系统状态 | 说明 |
|-----------|---------|------|
| `requires_payment_method` | `pending` | 需要支付方式 |
| `requires_confirmation` | `pending` | 需要确认 |
| `requires_action` | `pending` | 需要额外操作 |
| `processing` | `processing` | 处理中 |
| `requires_capture` | `pending` | 需要捕获 |
| `canceled` | `canceled` | 已取消 |
| `succeeded` | `succeeded` | 成功 |

## 数据库表结构

### stripe_payments 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| order_id | string | 内部订单ID |
| payment_intent_id | string | Stripe PaymentIntent ID |
| charge_id | string | Stripe Charge ID |
| amount | int64 | 支付金额(最小单位) |
| currency | string | 货币类型 |
| status | string | 支付状态 |
| payment_status | string | Stripe支付状态 |
| customer_id | string | Stripe客户ID |
| customer_email | string | 客户邮箱 |
| description | text | 支付描述 |
| receipt_url | string | 收据URL |
| metadata | json | 元数据 |
| user_id | uint | 关联用户ID |

### stripe_refunds 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| payment_id | uint | 支付ID |
| refund_id | string | Stripe退款ID |
| charge_id | string | 关联的Charge ID |
| amount | int64 | 退款金额(最小单位) |
| currency | string | 货币类型 |
| status | string | 退款状态 |
| reason | string | 退款原因 |
| description | text | 退款描述 |

### stripe_webhooks 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| stripe_id | string | Stripe事件ID |
| type | string | 事件类型 |
| data | json | 事件数据 |
| status | string | 处理状态 |
| error | text | 错误信息 |

## 支持的货币

Stripe 支持大多数国际货币,包括但不限于:

- USD - 美元
- EUR - 欧元
- GBP - 英镑
- JPY - 日元
- CAD - 加元
- AUD - 澳元
- CNY - 人民币

完整货币列表: [Supported Currencies](https://stripe.com/docs/currencies)

## 支付方式

Stripe 支持多种支付方式:

- **信用卡**: `card`
- **Alipay**: `alipay`
- **WeChat Pay**: `wechat_pay`
- **Apple Pay**: `apple_pay`
- **Google Pay**: `google_pay`
- **SEPA Debit**: `sepa_debit`
- **Sofort**: `sofort`

完整列表: [Payment Methods](https://stripe.com/docs/payments/payment-methods)

## 测试

### 测试卡号

Stripe 提供测试卡号用于测试:

```
4242 4242 4242 4242  - 成功支付
4000 0025 0000 3155  - 需要验证
4000 0000 0000 9995  - 拒绝
4000 0000 0000 0002  - 已过期
```

完整测试卡号列表: [Testing Cards](https://stripe.com/docs/testing)

### 测试端点

```bash
# 创建支付意图
curl -X POST http://localhost:8888/api/stripe/payments/intent \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "order_id": "TEST-001",
    "amount": 1000,
    "currency": "usd",
    "description": "测试支付"
  }'
```

## Webhook 事件

系统处理以下 Webhook 事件:

- `payment_intent.succeeded` - 支付成功
- `payment_intent.payment_failed` - 支付失败
- `charge.refunded` - 退款完成

更多事件类型: [Webhook Events](https://stripe.com/docs/api/events/types)

## 安全最佳实践

1. **密钥管理**
   - 永远不要在前端暴露 Secret Key
   - 使用 Publishable Key 进行客户端操作
   - 定期轮换 API 密钥

2. **Webhook 安全**
   - 验证 Webhook 签名
   - 使用 HTTPS 端点
   - 限制 Webhook 来源 IP

3. **支付安全**
   - 实施 3D Secure 验证
   - 使用 Stripe Radar 进行欺诈检测
   - 保存支付方式时需谨慎

## 错误处理

### 常见错误

| 错误代码 | 说明 | 解决方案 |
|---------|------|---------|
| `card_declined` | 卡被拒绝 | 提示用户使用其他卡片 |
| `insufficient_funds` | 余额不足 | 提示用户充值 |
| `expired_card` | 卡已过期 | 提示用户更新卡片信息 |
| `incorrect_cvc` | CVC 错误 | 提示用户重新输入 |
| `processing_error` | 处理错误 | 请用户稍后重试 |

### 错误处理示例

```go
intent, err := paymentService.CreatePaymentIntent(ctx, req)
if err != nil {
  if stripeErr, ok := err.(*stripe.Error); ok {
    switch stripeErr.Code {
    case stripe.ErrorCodeCardDeclined:
      // 处理卡被拒绝
    case stripe.ErrorCodeInsufficientFunds:
      // 处理余额不足
    default:
      // 处理其他错误
    }
  }
}
```

## 注意事项

1. **金额单位**: Stripe 使用最小货币单位(如美分为单位)
2. **异步处理**: 某些支付方式需要异步确认
3. **Webhook 顺序**: Webhook 可能不按顺序到达
4. **幂等性**: 重试请求时应使用相同的 ID
5. **合规性**: 确保遵守 PCI DSS 标准

## 相关链接

- [Stripe API 文档](https://stripe.com/docs/api)
- [Stripe Go SDK](https://github.com/stripe/stripe-go)
- [Stripe Dashboard](https://dashboard.stripe.com/)
- [Payment Intents 指南](https://stripe.com/docs/payments/payment-intents)
- [Webhooks 指南](https://stripe.com/docs/webhooks)
- [测试文档](https://stripe.com/docs/testing)

## 从测试切换到生产

1. 更新 `.env` 配置:
   ```env
   STRIPE_MODE=live
   STRIPE_SECRET_KEY=sk_live_...
   STRIPE_PUBLISHABLE_KEY=pk_live_...
   ```

2. 更新前端公钥

3. 配置生产环境 Webhook

4. 进行真实的测试支付

## 更新日志

- v1.0.0 (2026-02-07) - 初始版本,支持基本的支付、退款和 Webhook 功能
