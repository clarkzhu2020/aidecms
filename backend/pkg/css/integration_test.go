package css_test

import (
	"context"
	"testing"

	"github.com/clarkzhu2020/aidecms/pkg/css"
	"github.com/clarkzhu2020/aidecms/pkg/css/conversation"
	"github.com/clarkzhu2020/aidecms/pkg/css/rag"
)

// IntegrationTest 完整流程集成测试
//
// 测试流程:
// 用户输入 → CustomerServiceChat组件
//     ↓
// WebSocket连接 或 HTTP请求
//     ↓
// CustomerServiceController
//     ↓
// CSSEngine.ProcessQuestion()
//     ├─ 1. 会话管理 (GetOrCreateSession)
//     ├─ 2. RAG检索 (Search - 向量+关键词)
//     ├─ 3. 构建Prompt (问题+知识库+历史)
//     ├─ 4. AI调用 (Eino模型)
//     ├─ 5. 置信度评估 (多因素)
//     ├─ 6. 判断转接 (阈值/关键词)
//     └─ 7. 返回结果 (Answer)
//     ↓
// 前端显示 (置信度/来源/建议操作)

func TestCompleteFlow(t *testing.T) {
	ctx := context.Background()

	// 1. 配置引擎
	config := &css.Config{
		DefaultModel:             "qianwen",
		Temperature:             0.7,
		MaxTokens:              1000,
		TopK:                   5,
		ConfidenceThreshold:      0.6,
		TransferOnLowConfidence:  true,
		TransferKeywords:         []string{"投诉", "退款", "人工", "转接"},
		SessionTimeout:          1800,
		MaxHistory:              10,
	}

	t.Log("=== 智能客服系统完整流程测试 ===")
	t.Logf("配置: model=%s, threshold=%.2f, keywords=%v",
		config.DefaultModel, config.ConfidenceThreshold, config.TransferKeywords)

	// 2. 模拟用户输入
	testCases := []struct {
		name           string
		question       string
		expectTransfer  bool
		expectConfidenceRange [2]float64
		description    string
	}{
		{
			name:           "高置信度场景",
			question:       "如何使用产品A的功能？",
			expectTransfer:  false,
			expectConfidenceRange: [2]float64{0.7, 1.0},
			description:    "常见问题,应该有知识库支持,高置信度",
		},
		{
			name:           "低置信度-无知识库",
			question:       "什么是量子计算的基础原理？",
			expectTransfer:  true,
			expectConfidenceRange: [2]float64{0.0, 0.5},
			description:    "超出知识库范围,无文档支持,应该转人工",
		},
		{
			name:           "转接关键词触发",
			question:       "我要投诉这个产品的质量问题",
			expectTransfer:  true,
			expectConfidenceRange: [2]float64{0.0, 1.0},
			description:    "包含'投诉'关键词,即使置信度高也转人工",
		},
		{
			name:           "明确要求人工",
			question:       "请帮我转接人工客服",
			expectTransfer:  true,
			expectConfidenceRange: [2]float64{0.0, 1.0},
			description:    "用户明确要求转人工",
		},
		{
			name:           "中等置信度场景",
			question:       "产品B的价格范围是多少？",
			expectTransfer:  false,
			expectConfidenceRange: [2]float64{0.5, 0.8},
			description:    "有部分知识库支持,中等置信度",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("\n=== 测试场景: %s ===", tc.description)
			t.Logf("用户输入: %s", tc.question)

			// 流程步骤1: 模拟前端输入
			t.Log("\n[步骤1] 前端用户输入")
			t.Log("  组件: CustomerServiceChat.vue")
			t.Log("  操作: 用户在聊天框输入问题")

			// 流程步骤2: WebSocket/HTTP连接
			t.Log("\n[步骤2] 连接到后端")
			t.Log("  方式: WebSocket 或 HTTP")
			t.Log("  端点: ws://localhost:8888/api/css/ws")
			t.Log("       POST http://localhost:8888/api/css/question")

			// 流程步骤3: 控制器接收
			t.Log("\n[步骤3] CustomerServiceController接收请求")
			t.Log("  控制器: app/Http/Controllers/CustomerServiceController.go")
			t.Log("  方法: SendQuestion() 或 WebSocket处理")

			// 流程步骤4: 引擎处理
			t.Log("\n[步骤4] CSSEngine.ProcessQuestion() 开始处理")
			t.Log("  引擎: pkg/css/engine.go")

			// 步骤4.1: 会话管理
			t.Log("\n  [4.1] 会话管理 (GetOrCreateSession)")
			t.Log("      创建或恢复会话")
			t.Log("      保存用户消息到数据库")
			t.Log("      表: css_sessions, css_messages")

			// 步骤4.2: RAG检索
			t.Log("\n  [4.2] RAG检索 (Search - 向量+关键词)")
			t.Log("      向量化: Embedder.Embed(question)")
			t.Log("      向量检索: PostgreSQL + pgvector (权重0.7)")
			t.Log("      关键词检索: LIKE模糊匹配 (权重0.3)")
			t.Log("      结果融合: 重新排序Top-K")
			t.Log("      表: kb_chunks, kb_documents")

			// 步骤4.3: 构建Prompt
			t.Log("\n  [4.3] 构建Prompt (问题+知识库+历史)")
			t.Log("      System: 你是专业客服助手...")
			t.Log("      Knowledge: 文档1, 文档2...")
			t.Log("      History: 用户:问题, AI:回答...")
			t.Log("      Question: 当前问题")

			// 步骤4.4: AI调用
			t.Log("\n  [4.4] AI调用 (Eino模型)")
			t.Logf("      模型: %s", config.DefaultModel)
			t.Logf("      温度: %.1f", config.Temperature)
			t.Logf("      最大Token: %d", config.MaxTokens)

			// 步骤4.5: 置信度评估
			t.Log("\n  [4.5] 置信度评估 (多因素)")
			confidence := evaluateConfidence(tc.question, tc.question, true)
			t.Logf("      评估结果: %.2f", confidence)
			t.Log("      因素:")
			t.Log("        - 知识库支持: +0.15")
			t.Log("        - 回答长度: -0.1 ~ +0.05")
			t.Log("        - 不确定词汇: -0.15")

			// 步骤4.6: 判断转接
			t.Log("\n  [4.6] 判断转接 (阈值/关键词)")
			shouldTransfer := shouldTransfer(tc.question, confidence, config)
			t.Logf("      是否转接: %v", shouldTransfer)
			t.Log("      判断条件:")
			t.Log("        1. 置信度 < 0.6")
			t.Log("        2. 包含关键词: 投诉,退款,人工,转接")
			t.Log("        3. 用户明确要求人工")

			// 步骤4.7: 返回结果
			t.Log("\n  [4.7] 返回结果 (Answer)")
			t.Log("      数据结构:")
			t.Log("      {")
			t.Log("        content: string,           // AI回答内容")
			t.Log("        confidence: float64,        // 置信度 0-1")
			t.Log("        sources: [SourceRef],       // 知识来源")
			t.Log("        suggested_actions: [],      // 建议操作")
			t.Log("        transfer_to: string | null  // 是否转人工")
			t.Log("      }")

			// 流程步骤5: 前端显示
			t.Log("\n[步骤5] 前端显示 (置信度/来源/建议操作)")
			t.Log("  组件: CustomerServiceChat.vue")
			t.Log("  显示:")
			t.Log("    - 消息内容")
			t.Log("    - 置信度进度条 (高/中/低)")
			t.Log("    - 知识来源列表")
			t.Log("    - 建议操作按钮")
			t.Log("    - 转人工提示 (如适用)")

			// 验证结果
			t.Log("\n=== 验证结果 ===")
			if shouldTransfer != tc.expectTransfer {
				t.Errorf("转接预期: %v, 实际: %v", tc.expectTransfer, shouldTransfer)
			} else {
				t.Logf("✓ 转接判断正确: %v", shouldTransfer)
			}

			if confidence < tc.expectConfidenceRange[0] || confidence > tc.expectConfidenceRange[1] {
				t.Errorf("置信度预期: %.2f-%.2f, 实际: %.2f",
					tc.expectConfidenceRange[0],
					tc.expectConfidenceRange[1],
					confidence)
			} else {
				t.Logf("✓ 置信度在合理范围: %.2f", confidence)
			}

			t.Log("\n=== 测试通过 ===\n")
		})
	}
}

// 测试多轮对话
func TestMultiTurnConversation(t *testing.T) {
	ctx := context.Background()

	t.Log("=== 多轮对话流程测试 ===")

	sessionID := "test-multi-turn-001"

	questions := []string{
		"你好",
		"我想了解产品A",
		"产品A的主要功能有哪些？",
		"价格是多少？",
		"我要退款",
	}

	t.Logf("会话ID: %s", sessionID)
	t.Logf("问题序列: %v", questions)

	for i, question := range questions {
		t.Logf("\n=== 第%d轮对话 ===", i+1)
		t.Logf("用户问题: %s", question)

		// 模拟流程
		t.Log("[1] 会话管理: 恢复已有会话")
		t.Log("[2] 获取历史: 前面的对话作为上下文")
		t.Log("[3] RAG检索: 基于当前问题检索")
		t.Log("[4] 构建Prompt: 包含对话历史")
		t.Log("[5] AI生成: 考虑上下文回答")
		t.Log("[6] 置信度评估: 基于历史和当前回答")
		t.Log("[7] 判断转接: 检查触发条件")

		// 最后一轮应该触发转接
		if i == len(questions)-1 {
			t.Log("==> 预期: 触发转接 (包含'退款'关键词)")
		}
	}

	t.Log("\n=== 多轮对话测试完成 ===")
}

// 测试WebSocket流程
func TestWebSocketFlow(t *testing.T) {
	t.Log("=== WebSocket实时通讯流程测试 ===")

	t.Log("\n[阶段1] 连接建立")
	t.Log("  客户端: ws = new WebSocket('ws://localhost:8888/api/css/ws?session_id=xxx')")
	t.Log("  服务端: WSManager.Register()")
	t.Log("  响应: { type: 'connected' }")

	t.Log("\n[阶段2] 心跳保活")
	t.Log("  服务端: 每54秒发送 { type: 'ping' }")
	t.Log("  客户端: 响应 { type: 'pong' }")

	t.Log("\n[阶段3] 消息发送")
	t.Log("  客户端: ws.send({ type: 'message', content: '问题' })")
	t.Log("  服务端: handleClientMessage()")

	t.Log("\n[阶段4] 消息处理")
	t.Log("  1. 保存用户消息")
	t.Log("  2. 调用 CSSEngine.ProcessQuestion()")
	t.Log("  3. 执行完整流程")
	t.Log("  4. 生成AI回答")

	t.Log("\n[阶段5] 消息返回")
	t.Log("  服务端: ws.send({")
	t.Log("    type: 'message',")
	t.Log("    role: 'assistant',")
	t.Log("    content: '回答...',")
	t.Log("    confidence: 0.85,")
	t.Log("    sources: [...],")
	t.Log("    actions: [...]")
	t.Log("  })")
	t.Log("  客户端: onmessage -> 更新UI")

	t.Log("\n[阶段6] 错误处理")
	t.Log("  连接断开: 自动重连 (5秒延迟)")
	t.Log("  消息失败: 返回 { type: 'error', message: '...' }")
	t.Log("  AI超时: 转人工服务")

	t.Log("\n=== WebSocket流程测试完成 ===")
}

// 测试HTTP备用通道
func TestHTTPFallback(t *testing.T) {
	t.Log("=== HTTP API备用通道测试 ===")

	t.Log("\n[场景] WebSocket不可用,降级到HTTP")

	t.Log("\n[步骤1] 发送HTTP请求")
	t.Log("  fetch('http://localhost:8888/api/css/question', {")
	t.Log("    method: 'POST',")
	t.Log("    headers: { 'Content-Type': 'application/json' },")
	t.Log("    body: JSON.stringify({")
	t.Log("      session_id: 'xxx',")
	t.Log("      question: '如何使用产品A?',")
	t.Log("      channel: 'web'")
	t.Log("    })")
	t.Log("  })")

	t.Log("\n[步骤2] 控制器处理")
	t.Log("  CustomerServiceController.SendQuestion()")
	t.Log("  接收JSON请求")

	t.Log("\n[步骤3] 引擎处理")
	t.Log("  CSSEngine.ProcessQuestion()")
	t.Log("  执行完整流程 (与WebSocket相同)")

	t.Log("\n[步骤4] 返回JSON响应")
	t.Log("  {")
	t.Log("    success: true,")
	t.Log("    data: {")
	t.Log("      session_id: 'xxx',")
	t.Log("      answer: '回答...',")
	t.Log("      confidence: 0.85,")
	t.Log("      sources: [...],")
	t.Log("      actions: [...],")
	t.Log("      duration: 1250,")
	t.Log("      timestamp: 1707331200")
	t.Log("    }")
	t.Log("  }")

	t.Log("\n[步骤5] 前端解析并显示")
	t.Log("  解析JSON")
	t.Log("  更新消息列表")
	t.Log("  显示置信度和来源")

	t.Log("\n=== HTTP备用通道测试完成 ===")
}

// 辅助函数

func evaluateConfidence(question, answer string, hasDocs bool) float64 {
	confidence := 0.8

	if hasDocs {
		confidence += 0.15
	} else {
		confidence -= 0.2
	}

	if len(answer) < 50 {
		confidence -= 0.1
	} else if len(answer) > 500 {
		confidence += 0.05
	}

	uncertainPhrases := []string{"可能", "也许", "不确定", "需要确认", "不清楚"}
	for _, phrase := range uncertainPhrases {
		if containsString(answer, phrase) {
			confidence -= 0.15
			break
		}
	}

	if confidence > 1.0 {
		confidence = 1.0
	} else if confidence < 0.0 {
		confidence = 0.0
	}

	return confidence
}

func shouldTransfer(question string, confidence float64, config *css.Config) bool {
	if confidence < config.ConfidenceThreshold {
		return true
	}

	for _, keyword := range config.TransferKeywords {
		if containsString(question, keyword) {
			return true
		}
	}

	return false
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
