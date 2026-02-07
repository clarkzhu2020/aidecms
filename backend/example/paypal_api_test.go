package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	baseURL = "http://localhost:8888/api"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	// 测试Token (需要先登录获取)
	token := os.Getenv("TEST_TOKEN")
	if token == "" {
		fmt.Println("请设置环境变量 TEST_TOKEN")
		fmt.Println("例如: export TEST_TOKEN=your_jwt_token_here")
		return
	}

	fmt.Println("=== PayPal API 测试 ===\n")

	// 1. 创建支付订单
	fmt.Println("1. 创建支付订单...")
	paymentID, paypalOrderID := createPayment(token)
	if paymentID == 0 {
		fmt.Println("创建支付失败")
		return
	}

	// 2. 获取支付详情
	fmt.Println("\n2. 获取支付详情...")
	getPayment(token, paymentID)

	// 3. 获取支付列表
	fmt.Println("\n3. 获取支付列表...")
	listPayments(token)

	// 4. 捕获支付（需要实际的PayPal订单ID和用户授权）
	// 这个测试需要用户在PayPal页面上完成支付后才能执行
	// fmt.Println("\n4. 捕获支付...")
	// capturePayment(token, paypalOrderID)

	// 5. 退款测试（需要先有已支付的订单）
	// fmt.Println("\n5. 退款测试...")
	// refundPayment(token, paymentID)

	fmt.Println("\n=== 测试完成 ===")
}

// 创建支付订单
func createPayment(token string) (uint, string) {
	reqBody := map[string]interface{}{
		"order_id":      fmt.Sprintf("TEST-%d", time.Now().Unix()),
		"amount":        10.00,
		"currency":      "USD",
		"description":   "测试商品 - 支付API测试",
		"item_name":     "测试商品",
		"item_quantity": 1,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", baseURL+"/payments", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return 0, ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response Response
	json.Unmarshal(body, &response)

	if resp.StatusCode == 201 {
		fmt.Printf("✓ 创建成功\n")
		if data, ok := response.Data.(map[string]interface{}); ok {
			fmt.Printf("  - 支付ID: %.0f\n", data["payment_id"].(float64))
			fmt.Printf("  - 订单ID: %s\n", data["order_id"].(string))
			fmt.Printf("  - PayPal订单ID: %s\n", data["paypal_order_id"].(string))
			fmt.Printf("  - 金额: %.2f %s\n", data["amount"].(float64), data["currency"].(string))
			fmt.Printf("  - 审批URL: %s\n", data["approval_url"].(string))

			return uint(data["payment_id"].(float64)), data["paypal_order_id"].(string)
		}
	} else {
		fmt.Printf("✗ 创建失败: %s\n", response.Message)
		fmt.Printf("  响应: %s\n", string(body))
	}

	return 0, ""
}

// 获取支付详情
func getPayment(token string, paymentID uint) {
	url := fmt.Sprintf("%s/payments/%d", baseURL, paymentID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response Response
	json.Unmarshal(body, &response)

	if resp.StatusCode == 200 {
		fmt.Printf("✓ 获取成功\n")
		fmt.Printf("  响应: %s\n", string(body))
	} else {
		fmt.Printf("✗ 获取失败: %s\n", response.Message)
	}
}

// 获取支付列表
func listPayments(token string) {
	url := fmt.Sprintf("%s/payments?page=1&limit=10", baseURL)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response Response
	json.Unmarshal(body, &response)

	if resp.StatusCode == 200 {
		fmt.Printf("✓ 获取成功\n")
		fmt.Printf("  响应: %s\n", string(body))
	} else {
		fmt.Printf("✗ 获取失败: %s\n", response.Message)
	}
}

// 捕获支付
func capturePayment(token, paypalOrderID string) {
	url := fmt.Sprintf("%s/payments/capture/%s", baseURL, paypalOrderID)
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response Response
	json.Unmarshal(body, &response)

	if resp.StatusCode == 200 {
		fmt.Printf("✓ 捕获成功\n")
		fmt.Printf("  响应: %s\n", string(body))
	} else {
		fmt.Printf("✗ 捕获失败: %s\n", response.Message)
	}
}

// 退款
func refundPayment(token string, paymentID uint) {
	reqBody := map[string]interface{}{
		"amount": 5.00, // 部分退款
		"reason": "用户请求",
		"note":   "测试退款功能",
	}

	jsonData, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("%s/payments/%d/refund", baseURL, paymentID)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response Response
	json.Unmarshal(body, &response)

	if resp.StatusCode == 200 {
		fmt.Printf("✓ 退款成功\n")
		fmt.Printf("  响应: %s\n", string(body))
	} else {
		fmt.Printf("✗ 退款失败: %s\n", response.Message)
	}
}
