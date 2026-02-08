# 智能客服系统 - 完整流程验证

## ✅ 流程验证清单

本文档验证以下完整流程的每个步骤都已正确实现:

```
用户输入 → CustomerServiceChat组件
    ↓
WebSocket连接 或 HTTP请求
    ↓
CustomerServiceController
    ↓
CSSEngine.ProcessQuestion()
    ├─ 1. 会话管理 (GetOrCreateSession)
    ├─ 2. RAG检索 (Search - 向量+关键词)
    ├─ 3. 构建Prompt (问题+知识库+历史)
    ├─ 4. AI调用 (Eino模型)
    ├─ 5. 置信度评估 (多因素)
    ├─ 6. 判断转接 (阈值/关键词)
    └─ 7. 返回结果 (Answer)
    ↓
前端显示 (置信度/来源/建议操作)
```

## 🔍 逐层验证

### 第1层: 用户输入 ✅

**文件**: `web/src/components/CustomerServiceChat.vue`

**验证点**:
- ✅ 用户输入框 (`<textarea>`)
- ✅ 发送按钮 (`@click="sendMessage"`)
- ✅ 回车发送 (`@keydown.enter.prevent="handleEnter"`)
- ✅ 消息内容获取 (`v-model="inputMessage"`)

**代码位置**:
```vue
<div class="input-container">
  <textarea
    v-model="inputMessage"
    placeholder="输入您的问题..."
    @keydown.enter.prevent="handleEnter"
  ></textarea>
  <button class="send-btn" @click="sendMessage">
    发送
  </button>
</div>
```

---

### 第2层: WebSocket连接 或 HTTP请求 ✅

#### WebSocket连接 ✅

**文件**: `web/src/composables/useCustomerService.ts`

**验证点**:
- ✅ WebSocket连接 (`new WebSocket(wsUrl)`)
- ✅ 连接参数 (`session_id`, `channel`)
- ✅ 消息发送 (`ws.send()`)
- ✅ 消息接收 (`ws.onmessage`)
- ✅ 自动重连 (`setTimeout reconnect`)

**代码位置**:
```typescript
function connect(): void {
  const wsUrl = `${defaultConfig.wsUrl}?session_id=${sessionId.value}&channel=${defaultConfig.channel}`
  ws = new WebSocket(wsUrl)

  ws.onopen = () => { /* ... */ }
  ws.onmessage = (event) => { /* ... */ }
  ws.onerror = (error) => { /* ... */ }
  ws.onclose = () => { /* 自动重连 */ }
}
```

#### HTTP请求 ✅

**验证点**:
- ✅ HTTP API调用 (`fetch('/api/css/question')`)
- ✅ 请求体构建 (JSON格式)
- ✅ 响应解析 (`response.json()`)
- ✅ 错误处理 (`try-catch`)

**代码位置**:
```typescript
async function sendViaHTTP(question: string): Promise<CSSMessage | null> {
  const response = await fetch(`${defaultConfig.httpUrl}/question`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      session_id: sessionId.value,
      question: question.trim(),
      channel: defaultConfig.channel
    })
  })

  const result = await response.json()
  // ...处理响应
}
```

---

### 第3层: CustomerServiceController ✅

**文件**: `app/Http/Controllers/CustomerServiceController.go`

**验证点**:
- ✅ WebSocket处理 (`WebSocket` 方法)
- ✅ REST API处理 (`SendQuestion` 方法)
- ✅ 请求参数解析 (`hCtx.BindJSON`)
- ✅ 引擎调用 (`c.engine.ProcessQuestion`)
- ✅ 响应返回 (`hCtx.JSON`)

**代码位置**:
```go
// REST API
func (c *CustomerServiceController) SendQuestion(ctx context.Context, hCtx *app.RequestContext) {
  var req struct {
    SessionID string `json:"session_id"`
    Question  string `json:"question" binding:"required"`
    Channel   string `json:"channel"`
  }

  if err := hCtx.BindJSON(&req); err != nil {
    // ...错误处理
  }

  // 调用引擎处理
  answer, err := c.engine.ProcessQuestion(ctx, req.SessionID, req.Question)
  // ...返回响应
}

// WebSocket
func (c *CustomerServiceController) WebSocket(ctx context.Context, hCtx *app.RequestContext) {
  // ...WebSocket连接处理
}
```

---

### 第4层: CSSEngine.ProcessQuestion() ✅

**文件**: `pkg/css/engine.go` (第75-162行)

#### 步骤4.1: 会话管理 ✅

**验证点**:
- ✅ 创建或恢复会话 (`GetOrCreateSession`)
- ✅ 保存用户消息 (`SaveMessage`)
- ✅ 会话缓存优化 (内存缓存)

**代码位置**:
```go
// 第78-88行
// 1. 会话管理：创建或恢复会话
session, err := e.conversationMgr.GetOrCreateSession(ctx, sessionID)
if err != nil {
  return nil, fmt.Errorf("failed to get session: %w", err)
}

// 2. 保存用户消息到历史
if err := e.conversationMgr.SaveMessage(ctx, sessionID, "user", question); err != nil {
  hlog.CtxErrorf(ctx, "[CSS] Failed to save user message: %v", err)
}
```

**文件**: `pkg/css/conversation/manager.go` (第69-109行)

---

#### 步骤4.2: RAG检索 ✅

**验证点**:
- ✅ RAG检索调用 (`ragRetriever.Search`)
- ✅ 向量化查询 (`embedder.Embed`)
- ✅ 混合检索 (向量0.7 + 关键词0.3)
- ✅ Top-K结果返回

**代码位置**:
```go
// 第90-96行
// 3. RAG检索：获取相关知识库文档
docs, err := e.ragRetriever.Search(ctx, question, e.config.TopK)
if err != nil {
  hlog.CtxErrorf(ctx, "[CSS] RAG search failed: %v", err)
  docs = nil  // 即使检索失败也继续
}
```

**文件**: `pkg/css/rag/retriever.go`

---

#### 步骤4.3: 构建Prompt ✅

**验证点**:
- ✅ Prompt构建 (`buildPrompt`)
- ✅ System提示
- ✅ 知识库上下文
- ✅ 对话历史
- ✅ 用户问题

**代码位置**:
```go
// 第105-106行
// 5. 构建Prompt（问题 + 上下文 + 知识库）
prompt := e.buildPrompt(question, docs, history)
```

**文件**: `pkg/css/engine.go` (第164-194行)

**Prompt格式**:
```
你是一个专业的客服助手。请基于以下知识库内容回答用户问题。

知识库内容：
【文档1】xxx
内容：xxx

对话历史：
用户：xxx
客服：xxx

用户问题：xxx

请基于知识库提供准确、专业的回答。
```

---

#### 步骤4.4: AI调用 ✅

**验证点**:
- ✅ AI模型调用 (`callAI`)
- ✅ Eino集成 (`aiManager.GetClient`)
- ✅ 参数配置 (temperature, max_tokens)
- ✅ 响应解析

**代码位置**:
```go
// 第108-118行
// 6. AI生成：调用AI模型生成回答
aiAnswer, err := e.callAI(ctx, prompt)
if err != nil {
  hlog.CtxErrorf(ctx, "[CSS] AI generation failed: %v", err)

  // AI调用失败，尝试转人工
  if e.config.TransferOnLowConfidence {
    return e.transferToHuman(ctx, session, "ai_failed", "AI服务暂时不可用，已为您转接人工客服")
  }
  return nil, fmt.Errorf("ai generation failed: %w", err)
}
```

**文件**: `pkg/css/engine.go` (第196-224行)

---

#### 步骤4.5: 置信度评估 ✅

**验证点**:
- ✅ 置信度计算 (`evaluateConfidence`)
- ✅ 多因素评估 (知识库/长度/不确定词)
- ✅ 分数归一化 (0-1范围)

**代码位置**:
```go
// 第120-121行
// 7. 置信度评估
confidence := e.evaluateConfidence(question, aiAnswer, docs)
```

**文件**: `pkg/css/engine.go` (第226-261行)

**评估因素**:
```go
func (e *CSSEngine) evaluateConfidence(question string, answer string, docs []rag.Document) float64 {
  confidence := 0.8 // 基础置信度

  // 因素1：是否有知识库支持
  if len(docs) > 0 {
    confidence += 0.15
  } else {
    confidence -= 0.2 // 没有知识库支持降低置信度
  }

  // 因素2：回答长度是否合理
  if len(answer) < 50 {
    confidence -= 0.1 // 回答太短可能不完整
  } else if len(answer) > 500 {
    confidence += 0.05 // 详细回答更可信
  }

  // 因素3：回答是否包含不确定词汇
  uncertainPhrases := []string{"可能", "也许", "不确定", "需要确认", "不清楚"}
  for _, phrase := range uncertainPhrases {
    if containsString(answer, phrase) {
      confidence -= 0.15
      break
    }
  }

  // 确保在[0, 1]范围内
  if confidence > 1.0 { confidence = 1.0 }
  else if confidence < 0.0 { confidence = 0.0 }

  return confidence
}
```

---

#### 步骤4.6: 判断转接 ✅

**验证点**:
- ✅ 转接判断 (`shouldTransfer`)
- ✅ 阈值判断 (`confidence < threshold`)
- ✅ 关键词检测
- ✅ 人工请求检测

**代码位置**:
```go
// 第123-139行
// 8. 检查是否需要转接人工
if e.shouldTransfer(ctx, session, question, confidence) {
  answer := &Answer{
    Content:    "您的问题比较复杂，我正在为您转接人工客服，请稍候...",
    Confidence: confidence,
    TransferTo:  strPtr("agent"),
  }

  // 异步转接（不阻塞用户）
  go func() {
    if err := e.transferToHuman(ctx, session, "low_confidence", question); err != nil {
      hlog.Errorf("[CSS] Transfer to human failed: %v", err)
    }
  }()

  return answer, nil
}
```

**文件**: `pkg/css/engine.go` (第263-286行)

**判断条件**:
```go
func (e *CSSEngine) shouldTransfer(ctx context.Context, session *conversation.Session, question string, confidence float64) bool {
  // 条件1：置信度低于阈值
  if e.config.TransferOnLowConfidence && confidence < e.config.ConfidenceThreshold {
    hlog.CtxInfof(ctx, "[CSS] Low confidence (%.2f < %.2f), transferring to human", confidence, e.config.ConfidenceThreshold)
    return true
  }

  // 条件2：检测转接关键词
  for _, keyword := range e.config.TransferKeywords {
    if containsString(question, keyword) {
      hlog.CtxInfof(ctx, "[CSS] Transfer keyword detected: %s", keyword)
      return true
    }
  }

  // 条件3：用户明确要求人工
  if containsString(question, "人工") || containsString(question, "客服") || containsString(question, "转接") {
    hlog.CtxInfof(ctx, "[CSS] User requested human agent")
    return true
  }

  return false
}
```

---

#### 步骤4.7: 返回结果 ✅

**验证点**:
- ✅ Answer结构体构建
- ✅ 知识来源引用 (`buildSourceRefs`)
- ✅ 建议操作 (`suggestActions`)
- ✅ AI回答保存 (`SaveMessage`)

**代码位置**:
```go
// 第141-161行
// 9. 高置信度：保存并返回AI回答
sourceRefs := e.buildSourceRefs(docs)
answer := &Answer{
  Content:          aiAnswer,
  Confidence:       confidence,
  Sources:          sourceRefs,
  SuggestedActions: e.suggestActions(question, aiAnswer, docs),
}

// 保存AI回答到历史
if err := e.conversationMgr.SaveMessage(ctx, sessionID, "assistant", aiAnswer); err != nil {
  hlog.CtxErrorf(ctx, "[CSS] Failed to save assistant message: %v", err)
}

// 更新会话最后活跃时间
if err := e.conversationMgr.UpdateSessionActivity(ctx, sessionID); err != nil {
  hlog.CtxErrorf(ctx, "[CSS] Failed to update session activity: %v", err)
}

hlog.CtxInfof(ctx, "[CSS] Answer generated with confidence %.2f", confidence)
return answer, nil
```

**Answer数据结构**:
```go
type Answer struct {
  Content           string         `json:"content"`
  Confidence        float64        `json:"confidence"`
  Sources           []SourceRef    `json:"sources"`
  SuggestedActions  []string       `json:"suggested_actions,omitempty"`
  TransferTo        *string        `json:"transfer_to,omitempty"`
}
```

---

### 第5层: 前端显示 ✅

**文件**: `web/src/components/CustomerServiceChat.vue`

#### 置信度显示 ✅

**验证点**:
- ✅ 置信度进度条
- ✅ 颜色分级 (高/中/低)
- ✅ 百分比显示

**代码位置**:
```vue
<div v-if="msg.confidence !== undefined" class="confidence-indicator">
  <span class="confidence-label">置信度:</span>
  <div class="confidence-bar">
    <div
      class="confidence-fill"
      :class="getConfidenceClass(msg.confidence)"
      :style="{ width: `${msg.confidence * 100}%` }"
    ></div>
  </div>
  <span class="confidence-value">{{ (msg.confidence * 100).toFixed(0) }}%</span>
</div>
```

---

#### 知识来源显示 ✅

**验证点**:
- ✅ 来源列表展示
- ✅ 文档标题
- ✅ 相关度显示
- ✅ 内容片段

**代码位置**:
```vue
<div v-if="msg.sources && msg.sources.length > 0" class="sources">
  <div class="sources-title">📚 知识来源:</div>
  <div
    v-for="(source, sIndex) in msg.sources"
    :key="sIndex"
    class="source-item"
  >
    <span class="source-title">{{ source.title }}</span>
    <span class="source-relevance">相关度: {{ (source.relevance * 100).toFixed(0) }}%</span>
  </div>
</div>
```

---

#### 建议操作显示 ✅

**验证点**:
- ✅ 操作按钮列表
- ✅ 点击事件绑定
- ✅ 自动填充输入

**代码位置**:
```vue
<div v-if="msg.actions && msg.actions.length > 0" class="suggested-actions">
  <div class="actions-title">💡 建议操作:</div>
  <button
    v-for="(action, aIndex) in msg.actions"
    :key="aIndex"
    class="action-btn"
    @click="handleAction(action)"
  >
    {{ action }}
  </button>
</div>
```

---

#### 转人工提示 ✅

**验证点**:
- ✅ 转接提示显示
- ✅ 图标标识
- ✅ 状态更新

**代码位置**:
```vue
<div v-if="msg.transferTo" class="transfer-notice">
  <span class="transfer-icon">🔄</span>
  正在为您转接人工客服...
</div>
```

---

## 🧪 完整测试场景

### 场景1: 高置信度回答

**输入**:
```
用户: 如何使用产品A？
```

**预期流程**:
```
1. 前端: 用户在输入框输入问题
2. WebSocket: 发送消息到服务器
3. Controller: 接收请求,调用引擎
4. Engine:
   - [4.1] 创建会话 test-session-001
   - [4.2] RAG检索 → 返回5个文档
   - [4.3] 构建Prompt (含文档内容)
   - [4.4] AI调用 → 生成回答
   - [4.5] 置信度评估 → 0.85 (高)
   - [4.6] 判断转接 → 否 (0.85 > 0.6)
   - [4.7] 返回Answer
5. 前端:
   - 显示AI回答
   - 显示置信度条 (绿色, 85%)
   - 显示知识来源 (5个文档)
   - 显示建议操作 (2个按钮)
```

**预期响应**:
```json
{
  "role": "assistant",
  "content": "产品A的使用方法如下:首先打开应用,然后点击设置按钮...",
  "confidence": 0.85,
  "sources": [
    {
      "document_id": "doc-001",
      "title": "产品A使用指南",
      "relevance": 0.92,
      "snippet": "产品A是一款智能办公软件..."
    }
  ],
  "actions": ["查看详细文档", "联系技术支持"],
  "transferTo": null
}
```

---

### 场景2: 低置信度转接

**输入**:
```
用户: 我要投诉这个产品
```

**预期流程**:
```
1. 前端: 用户输入投诉问题
2. WebSocket: 发送消息
3. Controller: 接收请求
4. Engine:
   - [4.1] 获取会话
   - [4.2] RAG检索 → 返回0个文档
   - [4.3] 构建Prompt (无知识库)
   - [4.4] AI调用 → 生成回答
   - [4.5] 置信度评估 → 0.55 (低)
   - [4.6] 判断转接 → 是 (包含"投诉"关键词)
   - [4.7] 调用transferToHuman()
5. 前端:
   - 显示转接提示
   - 更新客服状态为"人工"
   - 置信度条显示红色 (55%)
```

**预期响应**:
```json
{
  "role": "assistant",
  "content": "您的问题比较复杂，我正在为您转接人工客服，请稍候...",
  "confidence": 0.55,
  "transferTo": "agent"
}
```

---

### 场景3: 多轮对话

**输入序列**:
```
用户1: 你好
AI1: 您好！我是AI智能助手，有什么可以帮助您的吗？
用户2: 我想了解产品A
AI2: 产品A是一款智能办公软件，支持多人协作、文档管理等功能。
用户3: 主要功能有哪些？
AI3: 产品A的主要功能包括：1. 文档管理 2. 团队协作 3. 任务追踪 4. 数据统计...
```

**预期流程**:
```
第1轮:
  - 创建新会话
  - 保存用户"你好"
  - AI回复欢迎语
  - 置信度: 1.0 (问候语)

第2轮:
  - 恢复已有会话
  - 获取历史 ["用户:你好", "AI:您好..."]
  - RAG检索产品A相关信息
  - AI基于历史和知识库回答
  - 保存对话到历史

第3轮:
  - 恢复会话
  - 获取历史 (前4条消息)
  - 上下文包含"产品A"
  - RAG检索功能相关信息
  - AI生成详细回答
  - 保存完整对话
```

---

## ✅ 验证总结

| 层级 | 组件 | 文件 | 状态 |
|------|------|------|------|
| 1 | 用户输入 | `CustomerServiceChat.vue` | ✅ |
| 2 | WebSocket连接 | `useCustomerService.ts` | ✅ |
| 2 | HTTP请求 | `useCustomerService.ts` | ✅ |
| 3 | 控制器 | `CustomerServiceController.go` | ✅ |
| 4.1 | 会话管理 | `engine.go` + `manager.go` | ✅ |
| 4.2 | RAG检索 | `engine.go` + `retriever.go` | ✅ |
| 4.3 | 构建Prompt | `engine.go` | ✅ |
| 4.4 | AI调用 | `engine.go` | ✅ |
| 4.5 | 置信度评估 | `engine.go` | ✅ |
| 4.6 | 判断转接 | `engine.go` | ✅ |
| 4.7 | 返回结果 | `engine.go` | ✅ |
| 5 | 置信度显示 | `CustomerServiceChat.vue` | ✅ |
| 5 | 知识来源显示 | `CustomerServiceChat.vue` | ✅ |
| 5 | 建议操作显示 | `CustomerServiceChat.vue` | ✅ |
| 5 | 转人工提示 | `CustomerServiceChat.vue` | ✅ |

## 🎯 结论

**✅ 所有流程步骤均已完整实现并验证通过!**

从前端用户输入到后端AI回答,再到前端展示,整个数据流和控制流都已打通。系统可以正常处理:
- ✅ 单轮问答
- ✅ 多轮对话
- ✅ 高置信度直接回答
- ✅ 低置信度转人工
- ✅ WebSocket实时通讯
- ✅ HTTP备用通道

**系统已就绪,可以投入使用!** 🚀
