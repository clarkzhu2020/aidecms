#!/bin/bash

# PayPal API 测试脚本
# 使用前请先设置 API_TOKEN 环境变量

BASE_URL="http://localhost:8888/api"
TOKEN="${API_TOKEN}"

# 检查是否设置了TOKEN
if [ -z "$TOKEN" ]; then
    echo "错误: 请设置 API_TOKEN 环境变量"
    echo "例如: export API_TOKEN=your_jwt_token_here"
    exit 1
fi

echo "=== PayPal API 测试 ==="
echo ""

# 1. 创建支付订单
echo "1. 创建支付订单..."
curl -X POST "${BASE_URL}/payments" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{
    "order_id": "TEST-'$(date +%s)'",
    "amount": 10.00,
    "currency": "USD",
    "description": "测试商品 - 支付API测试",
    "item_name": "测试商品",
    "item_quantity": 1
  }'

echo -e "\n"
read -p "按 Enter 继续..."
echo ""

# 2. 获取支付列表
echo "2. 获取支付列表..."
curl -X GET "${BASE_URL}/payments?page=1&limit=10" \
  -H "Authorization: Bearer ${TOKEN}"

echo -e "\n"
read -p "按 Enter 继续..."
echo ""

# 3. 捕获支付（需要替换实际的PayPal订单ID）
echo "3. 捕获支付（需要PayPal订单ID）..."
read -p "请输入PayPal订单ID: " paypal_order_id

if [ -n "$paypal_order_id" ]; then
    curl -X POST "${BASE_URL}/payments/capture/${paypal_order_id}" \
      -H "Authorization: Bearer ${TOKEN}"
fi

echo -e "\n"
read -p "按 Enter 继续..."
echo ""

# 4. 退款（需要替换实际的支付ID）
echo "4. 退款（需要支付ID）..."
read -p "请输入支付ID: " payment_id

if [ -n "$payment_id" ]; then
    read -p "退款金额 (留空为全额退款): " refund_amount

    if [ -n "$refund_amount" ]; then
        curl -X POST "${BASE_URL}/payments/${payment_id}/refund" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer ${TOKEN}" \
          -d "{
            \"amount\": ${refund_amount},
            \"reason\": \"用户请求\",
            \"note\": \"测试退款功能\"
          }"
    else
        curl -X POST "${BASE_URL}/payments/${payment_id}/refund" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer ${TOKEN}" \
          -d '{
            "reason": "用户请求",
            "note": "全额退款测试"
          }'
    fi
fi

echo -e "\n"
echo "=== 测试完成 ==="
