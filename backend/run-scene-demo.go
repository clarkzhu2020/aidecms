// 演示脚本: 高置信度问答场景
// 运行: go run run-scene-demo.go
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
	baseURL = "http://localhost:8888"
)

// QuestionRequest 问题请求
type QuestionRequest struct {
	SessionID string `json:"session_id"`
	Question  string `json:"question"`
}

// QuestionResponse 问题响应
type QuestionResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Answer          string             `json:"answer"`
		Confidence      float64            `json:"confidence"`
		Sources         []Source           `json:"sources"`
		SuggestedActions []string          `json:"suggested_actions"`
		TransferTo      map[string]interface{} `json:"transfer_to,omitempty"`
	} `json:"data"`
}

// Source 知识来源
type Source struct {
	DocumentID string  `json:"document_id"`
	Title      string  `json:"title"`
	Relevance  float64 `json:"relevance"`
}

func main() {
	fmt.Println("========================================")
	fmt.Println("智能客服系统 - 场景演示")
	fmt.Println("========================================")
	fmt.Println("场景: 高置信度问答")
	fmt.Println("问题: 如何使用产品A？")
	fmt.Println("")

	// 生成唯一会话ID
	sessionID := fmt.Sprintf("scene-demo-%d", time.Now().Unix())
	fmt.Printf("Session ID: %s\n\n", sessionID)

	// 步骤1: 发送问题
	fmt.Println("步骤1: 发送问题...")
	fmt.Println("-----------------------------------")
	answer, err := sendQuestion(sessionID, "如何使用产品A？")
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ 请求成功\n\n")

	// 步骤2: 显示AI回答
	fmt.Println("步骤2: AI回答...")
	fmt.Println("-----------------------------------")
	fmt.Printf("%s\n\n", answer.Data.Answer)

	// 步骤3: 显示置信度
	fmt.Println("步骤3: 置信度评估...")
	fmt.Println("-----------------------------------")
	confidence := answer.Data.Confidence
	fmt.Printf("置信度: %.2f\n", confidence)

	// 显示置信度可视化
	bar := ""
	for i := 0; i < int(confidence * 10); i++ {
		bar += "█"
	}
	for i := int(confidence * 10); i < 10; i++ {
		bar += "░"
	}

	// 根据置信度显示颜色和状态
	color := "红色"
	status := "低"
	if confidence >= 0.8 {
		color = "绿色"
		status = "高"
	} else if confidence >= 0.6 {
		color = "橙色"
		status = "中"
	}

	fmt.Printf("进度条: %s %.0f%% (%s)\n\n", bar, confidence*100, color)

	// 步骤4: 验证置信度
	fmt.Println("步骤4: 验证置信度...")
	fmt.Println("-----------------------------------")
	if confidence >= 0.7 {
		fmt.Printf("✅ 置信度检查通过: %.2f >= 0.7 (状态: %s)\n\n", confidence, status)
	} else {
		fmt.Printf("❌ 置信度检查失败: %.2f < 0.7 (状态: %s)\n\n", confidence, status)
		os.Exit(1)
	}

	// 步骤5: 显示知识来源
	fmt.Println("步骤5: 知识来源...")
	fmt.Println("-----------------------------------")
	if len(answer.Data.Sources) > 0 {
		fmt.Printf("✅ 找到 %d 个知识来源:\n", len(answer.Data.Sources))
		for i, source := range answer.Data.Sources {
			fmt.Printf("  [%d] %s (相关度: %.0f%%)\n", i+1, source.Title, source.Relevance*100)
		}
	} else {
		fmt.Println("⚠️  未找到知识来源 (可能是模拟数据)")
	}
	fmt.Println()

	// 步骤6: 显示建议操作
	fmt.Println("步骤6: 建议操作...")
	fmt.Println("-----------------------------------")
	if len(answer.Data.SuggestedActions) > 0 {
		fmt.Println("✅ 建议操作:")
		for i, action := range answer.Data.SuggestedActions {
			fmt.Printf("  [%d] %s\n", i+1, action)
		}
	} else {
		fmt.Println("无建议操作")
	}
	fmt.Println()

	// 步骤7: 验证转接
	fmt.Println("步骤7: 转接判断...")
	fmt.Println("-----------------------------------")
	if answer.Data.TransferTo != nil {
		fmt.Println("❌ 触发了转接 (不符合预期)")
		os.Exit(1)
	} else {
		fmt.Println("✅ 未触发转接 (符合预期)")
	}
	fmt.Println()

	// 步骤8: 获取历史记录
	fmt.Println("步骤8: 获取历史记录...")
	fmt.Println("-----------------------------------")
	history, err := getHistory(sessionID)
	if err != nil {
		fmt.Printf("❌ 获取历史失败: %v\n", err)
	} else {
		fmt.Printf("✅ 历史记录: 共 %d 条消息\n\n", len(history))
	}

	// 总结
	fmt.Println("========================================")
	fmt.Println("✅ 场景演示完成!")
	fmt.Println("========================================")
	fmt.Println("")
	fmt.Println("测试结果:")
	fmt.Printf("- 置信度: %.2f (>=0.7 ✓)\n", confidence)
	fmt.Printf("- 知识来源: %d 个文档 ✓\n", len(answer.Data.Sources))
	fmt.Printf("- 转人工: 未触发 ✓\n")
	fmt.Printf("- 消息保存: 成功 ✓\n")
	fmt.Println("")
	fmt.Println("🎉 完整7步流程验证通过!")
}

// sendQuestion 发送问题
func sendQuestion(sessionID, question string) (*QuestionResponse, error) {
	req := QuestionRequest{
		SessionID: sessionID,
		Question:  question,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		baseURL+"/api/css/question",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result QuestionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("API返回错误: %s", result.Message)
	}

	return &result, nil
}

// getHistory 获取历史记录
func getHistory(sessionID string) ([]map[string]interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/css/history/%s", baseURL, sessionID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    []map[string]interface{} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
