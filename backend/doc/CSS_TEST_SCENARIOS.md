# 智能客服系统 - 测试场景文档

## 场景1: 高置信度问答

### 基本信息
- **场景ID**: scene-high-confidence
- **测试问题**: "如何使用产品A？"
- **预期行为**: 置信度>=0.7, 显示知识来源, 不转接

### 完整流程
```
步骤1: 用户输入
  问题: "如何使用产品A？"

步骤2: WebSocket连接
  连接到: ws://localhost:8888/api/css/ws
  Session ID: scene-high-conf-001

步骤3: CustomerServiceController
  接收请求, 调用引擎

步骤4: CSSEngine.ProcessQuestion()
  4.1 会话管理: GetOrCreateSession("scene-high-conf-001")
  4.2 保存消息: SaveMessage(sessionID, "user", question)

  4.3 RAG检索: Search("如何使用产品A?", topK=5)
       检索结果:
       - [1] 产品A使用指南 (相关度: 0.92)
       - [2] 快速入门指南 (相关度: 0.85)

  4.4 获取历史: GetHistory(sessionID, maxHistory=10)
       历史消息: 0 条(首次对话)

  4.5 构建Prompt
       System: 你是一个智能客服助手...
       Context:
         知识库:
         1. 产品A使用指南: 产品A是一款智能办公软件。使用方法...
         2. 快速入门指南: 产品A提供快速入门功能...
       History: (无)
       User: 如何使用产品A？

  4.6 AI调用: callAI(prompt)
       AI回答: "产品A的使用方法如下：首先打开应用，然后点击设置按钮，
               选择产品A，按照屏幕提示完成配置即可。如果需要帮助，
               可以查看快速入门指南。"

  4.7 置信度评估: evaluateConfidence(question, answer, docs)
       基础分: 0.8
       + 有知识库支持: +0.15
       + 回答长度适中: +0.05
       - 无不确定词汇: 0.00
       = 最终置信度: 0.92

  4.8 判断转接: shouldTransfer(question, confidence)
       检查转接关键词: 无 (✓)
       检查置信度: 0.92 >= 0.6 (✓)
       结果: 不转接

  4.9 返回结果
       Answer {
         Content: "产品A的使用方法如下...",
         Confidence: 0.92,
         Sources: [
           {
             DocumentID: "doc-product-a",
             Title: "产品A使用指南",
             Relevance: 0.92
           }
         ],
         SuggestedActions: [
           "查看详细文档",
           "联系技术支持"
         ],
         TransferTo: null
       }

步骤5: 前端显示
  - 显示AI回答
  - 显示置信度进度条: ██████████ 92% (绿色)
  - 显示知识来源: "产品A使用指南"
  - 显示建议操作: "查看详细文档" "联系技术支持"
  - 无转人工提示
```

### 验证清单
- [x] 会话已创建
- [x] 用户消息已保存
- [x] RAG检索返回相关文档
- [x] AI回答有内容
- [x] 置信度 >= 0.7
- [x] 显示知识来源
- [x] 不触发转接
- [x] AI消息已保存

### 测试方法

#### 方法1: 运行单元测试
```bash
cd backend
go test -v ./pkg/css -run SceneHighConfidenceTest
```

#### 方法2: 运行Shell脚本
```bash
cd backend
chmod +x test-scene-high-confidence.sh
./test-scene-high-confidence.sh
```

#### 方法3: 运行Windows批处理
```cmd
cd backend
test-scene-high-confidence.bat
```

#### 方法4: 手动测试
```bash
# 1. 启动服务
CSS_ENABLED=true
go run main.go

# 2. 发送请求
curl -X POST http://localhost:8888/api/css/question \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-001",
    "question": "如何使用产品A？"
  }'

# 3. 验证响应
# 预期:
{
  "answer": "产品A的使用方法如下...",
  "confidence": 0.92,
  "sources": [...],
  "actions": [...]
}
```

#### 方法5: 使用前端界面
```bash
cd web
npm install
npm run dev
# 访问 http://localhost:5173/customer-service-demo
```

### 预期响应示例
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "answer": "产品A的使用方法如下：首先打开应用，然后点击设置按钮，选择产品A，按照屏幕提示完成配置即可。如果需要帮助，可以查看快速入门指南。",
    "confidence": 0.92,
    "sources": [
      {
        "document_id": "doc-product-a",
        "title": "产品A使用指南",
        "content": "产品A是一款智能办公软件。使用方法...",
        "relevance": 0.92
      }
    ],
    "suggested_actions": [
      "查看详细文档",
      "联系技术支持"
    ]
  }
}
```

---

## 场景2: 低置信度转接

### 基本信息
- **场景ID**: scene-low-confidence
- **测试问题**: "我要投诉这个产品"
- **预期行为**: 触发转人工

### 触发条件
- 包含关键词: "投诉" ✓
- 置信度: 较低 (0.45)
- 转接原因: 关键词匹配

### 预期响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "answer": "您的问题比较复杂，我正在为您转接人工客服，请稍候...",
    "confidence": 0.45,
    "transfer_to": {
      "agent_id": "agent-001",
      "agent_name": "客服专员",
      "estimated_wait_time": "30秒"
    }
  }
}
```

---

## 场景3: 多轮对话

### 基本信息
- **场景ID**: scene-multi-turn
- **测试流程**: 连续多轮问答
- **预期行为**: 每轮都能基于历史上下文回答

### 对话流程
```
第1轮:
  用户: "你好，我想了解产品A"
  AI: "您好！产品A是一款智能办公软件，支持多人协作..."

第2轮:
  用户: "产品A的价格是多少？"
  AI: "基于您对产品A的询问，产品A的定价为..."
  (历史上下文: 用户对产品A感兴趣)

第3轮:
  用户: "怎么购买？"
  AI: "您可以通过以下方式购买产品A..."
  (历史上下文: 产品A + 价格询问)
```

---

## 运行所有测试

### 运行所有场景测试
```bash
cd backend
go test -v ./pkg/css -run Scene
```

### 运行所有测试
```bash
cd backend
go test -v ./pkg/css
```

---

## 测试环境要求

### 数据库
- MySQL: 存储会话和消息
- PostgreSQL + pgvector: 存储知识库和向量

### 环境变量
```bash
CSS_ENABLED=true
CSS_DEFAULT_MODEL=qianwen
CSS_CONFIDENCE_THRESHOLD=0.6
CSS_TRANSFER_ON_LOW_CONFIDENCE=true
CSS_TRANSFER_KEYWORDS=投诉,退款,人工,转接
```

### 知识库
- 至少包含1个测试文档
- 文档已向量化并存储到PostgreSQL

---

## 调试技巧

### 查看日志
```bash
# 查看后端日志
tail -f logs/app.log

# 查看CSS模块日志
tail -f logs/css.log
```

### 查看数据库
```sql
-- 查看会话
SELECT * FROM css_sessions WHERE session_id = 'test-001';

-- 查看消息
SELECT * FROM css_messages WHERE session_id = 'test-001';

-- 查看知识库
SELECT * FROM kb_documents LIMIT 10;
```

### WebSocket调试
```javascript
// 浏览器控制台
const ws = new WebSocket('ws://localhost:8888/api/css/ws');

ws.onmessage = (event) => {
  console.log('收到消息:', JSON.parse(event.data));
};

ws.send(JSON.stringify({
  type: 'message',
  data: {
    session_id: 'test-001',
    question: '如何使用产品A？'
  }
}));
```

---

## 故障排查

### 问题: 置信度总是很低
**原因**: 知识库为空或文档不相关
**解决**:
```bash
# 添加测试文档到知识库
curl -X POST http://localhost:8888/api/kb/documents \
  -H "Content-Type: application/json" \
  -d '{"title": "产品A使用指南", "content": "..."}'
```

### 问题: 始终转接人工
**原因**: 转接关键词匹配过于宽泛
**解决**:
```bash
# 调整CSS_TRANSFER_KEYWORDS
CSS_TRANSFER_KEYWORDS=投诉,严重投诉,人工客服
```

### 问题: WebSocket连接失败
**原因**: 后端未启动或端口被占用
**解决**:
```bash
# 检查后端状态
curl http://localhost:8888/api/css/status

# 检查端口占用
netstat -an | grep 8888
```
