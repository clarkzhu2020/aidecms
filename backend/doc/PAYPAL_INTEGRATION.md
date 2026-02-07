# PayPal 集成文档

## 概述

本系统已集成 PayPal 支付功能,提供完整的支付、退款和 Webhook 处理能力。

## 功能特性

- ✅ 创建支付订单
- ✅ 捕获支付
- ✅ 退款处理
- ✅ 支付状态查询
- ✅ Webhook 通知处理
- ✅ 支付历史记录
- ✅ 沙盒/生产环境支持

## 安装步骤

### 1. 安装依赖

```bash
cd backend
go mod tidy
```

### 2. 配置环境变量

在 `.env` 文件中添加以下配置:

```env
# PayPal 支付配置
PAYPAL_ENABLED=true
PAYPAL_CLIENT_ID=your_client_id
PAYPAL_CLIENT_SECRET=your_client_secret
PAYPAL_MODE=sandbox  # sandbox 或 live
PAYPAL_WEBHOOK_ID=your_webhook_id  # 可选
```

### 3. 获取 PayPal 凭据

1. 登录 [PayPal Developer Dashboard](https://developer.paypal.com/dashboard/)
2. 创建新的应用程序
3. 获取 Client ID 和 Client Secret
4. 配置 Webhook (可选)

## API 接口

### 1. 创建支付订单

```http
POST /api/payments
Content-Type: application/json
Authorization: Bearer {token}

{
  "order_id": "ORDER-123456",
  "amount": 99.99,
  "currency": "USD",
  "description": "商品描述",
  "item_name": "商品名称",
  "item_quantity": 1
}
```

响应示例:

```json
{
  "code": 201,
  "message": "Payment order created successfully",
  "data": {
    "payment_id": 1,
    "order_id": "ORDER-123456",
    "paypal_order_id": "O-4J082351X3132253H",
    "approval_url": "https://www.sandbox.paypal.com/checkoutnow?token=...",
    "amount": 99.99,
    "currency": "USD",
    "status": "pending"
  }
}
```

### 2. 获取支付详情

```http
GET /api/payments/{id}
Authorization: Bearer {token}
```

### 3. 捕获支付

```http
POST /api/payments/capture/{paypal_order_id}
```

### 4. 退款

```http
POST /api/payments/{payment_id}/refund
Content-Type: application/json

{
  "amount": 50.00,  // 可选,不填则为全额退款
  "reason": "客户要求",
  "note": "备注信息"
}
```

### 5. 获取支付列表

```http
GET /api/payments?status=paid&page=1&limit=20
Authorization: Bearer {token}
```

### 6. Webhook 端点

```http
POST /api/payments/webhook
Content-Type: application/json
```

## 支付状态

| 状态 | 说明 |
|------|------|
| `pending` | 待支付 |
| `paid` | 已支付 |
| `failed` | 支付失败 |
| `cancelled` | 已取消 |
| `refunded` | 已退款 |

## 数据库表结构

### payments 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| order_id | string | 内部订单ID |
| paypal_order_id | string | PayPal订单ID |
| amount | decimal | 支付金额 |
| currency | string | 货币类型 |
| status | string | 支付状态 |
| payment_status | string | PayPal支付状态 |
| payer_id | string | 支付者ID |
| payer_email | string | 支付者邮箱 |
| capture_id | string | 捕获ID |
| description | text | 订单描述 |
| approval_url | string | 支付审批URL |
| user_id | uint | 关联用户ID |

### payment_refunds 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| payment_id | uint | 支付ID |
| refund_id | string | PayPal退款ID |
| amount | decimal | 退款金额 |
| currency | string | 货币类型 |
| status | string | 退款状态 |
| reason | text | 退款原因 |

### payment_webhooks 表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| event_id | string | 事件ID |
| event_type | string | 事件类型 |
| resource_type | string | 资源类型 |
| resource_id | string | 资源ID |
| status | string | 处理状态 |
| raw_data | json | 原始数据 |

## 使用示例

### Go 代码示例

```go
package main

import (
    "context"
    "fmt"
    "github.com/clarkzhu2020/aidecms/pkg/paypal"
)

func main() {
    // 初始化PayPal客户端
    err := paypal.Init(
        "your_client_id",
        "your_client_secret",
        true, // 沙盒模式
    )
    if err != nil {
        panic(err)
    }

    // 创建订单服务
    orderService := paypal.NewOrderService()

    // 创建订单
    req := &paypal.CreateOrderRequest{
        OrderID:      "ORDER-001",
        Amount:       99.99,
        Currency:     "USD",
        Description:  "测试商品",
        ReturnURL:    "http://localhost:8888/api/payments/success",
        CancelURL:    "http://localhost:8888/api/payments/cancel",
        ReferenceID:  "REF-001",
    }

    order, err := orderService.CreateOrder(context.Background(), req)
    if err != nil {
        panic(err)
    }

    fmt.Printf("订单创建成功: %s\n", order.ID)
}
```

## 支持的货币

- USD - 美元
- EUR - 欧元
- GBP - 英镑
- JPY - 日元
- CAD - 加元
- AUD - 澳元
- CNY - 人民币 (部分支持)

## 测试

### 沙盒测试

1. 使用沙盒账户进行测试
2. 使用测试信用卡号: `4111111111111111`
3. 过期日期任意,验证码任意

### 测试端点

```bash
# 测试创建支付
curl -X POST http://localhost:8888/api/payments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "order_id": "TEST-001",
    "amount": 10.00,
    "currency": "USD",
    "description": "测试支付"
  }'
```

## 注意事项

1. **安全性**: 不要在前端暴露 Client Secret
2. **Webhook**: 建议配置 Webhook 接收支付状态通知
3. **错误处理**: 妥善处理网络错误和 API 异常
4. **日志记录**: 记录所有支付相关操作以便审计
5. **生产环境**: 部署前切换到 live 模式

## 常见问题

### Q: 如何从沙盒切换到生产环境?

A: 修改 `.env` 文件中的 `PAYPAL_MODE=live` 并使用生产环境的 Client ID 和 Secret。

### Q: 支付失败如何处理?

A: 检查支付记录的状态,可以使用 `GetPayment` API 查询详细错误信息。

### Q: 如何处理部分退款?

A: 在退款请求中指定 `amount` 参数,不指定则全额退款。

## 相关链接

- [PayPal API 官方文档](https://developer.paypal.com/api/rest/)
- [PayPal Go SDK](https://github.com/plutov/paypal)
- [PayPal 开发者仪表板](https://developer.paypal.com/dashboard/)

## 更新日志

- v1.0.0 (2026-02-07) - 初始版本,支持基本的支付、退款和 Webhook 功能
