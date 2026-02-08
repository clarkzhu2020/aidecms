# 智能客服系统 - 快速部署指南

## 🎯 系统概览

本文档说明如何部署和使用已实现的智能客服系统核心功能。

### 核心流程

```
用户提问 → 会话管理 → RAG检索知识库 → 构建Prompt → AI生成
        → 置信度评估 → (高置信度) 返回回答 / (低置信度) 转人工
```

## 📦 部署步骤

### 1. 环境配置

在 `backend/.env` 文件中添加或更新以下配置:

```bash
# 客服系统配置
CSS_ENABLED=true                          # 启用客服系统
CSS_DEFAULT_MODEL=qianwen                 # 默认AI模型
CSS_TEMPERATURE=0.7                        # 温度参数
CSS_MAX_TOKENS=1000                        # 最大Token数
CSS_TOP_K=5                                # RAG检索Top-K
CSS_CONFIDENCE_THRESHOLD=0.6               # 置信度阈值
CSS_TRANSFER_ON_LOW_CONFIDENCE=true        # 低置信度时转接
CSS_MAX_RETRIES=3                          # 最大重试次数
CSS_TRANSFER_KEYWORDS=投诉,退款,人工,转接   # 转接关键词
CSS_SESSION_TIMEOUT=1800                   # 会话超时(秒)
CSS_MAX_HISTORY=10                         # 最大历史消息数
```

### 2. 数据库迁移

运行数据库迁移创建所需表结构:

```bash
cd backend
go run cmd/artisan/main.go artisan migrate
```

这将创建以下表:
- `css_sessions` - 会话表
- `css_messages` - 消息表
- `kb_documents` - 知识库文档表
- `kb_chunks` - 文档分块表
- `cs_agents` - 客服人员表
- `cs_transfers` - 转接记录表
- `css_feedback` - 用户反馈表

### 3. 启动服务

```bash
cd backend
go run main.go
```

服务启动后会自动初始化客服系统,你将看到如下日志:

```
[CSS] Initializing customer service system...
[CSS] Config: model=qianwen, top_k=5, confidence=0.60
[CSS] Conversation manager initialized
[CSS] RAG retriever initialized
[CSS] Knowledge base service initialized
[CSS] Core engine initialized
[CSS] WebSocket manager started
[CSS] Customer service system initialized successfully
```

## 🧪 测试API

### 1. 发送问题 (REST API)

```bash
curl -X POST http://localhost:8888/api/css/question \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-session-001",
    "question": "如何使用产品A？",
    "channel": "web"
  }'
```

**响应示例:**

```json
{
  "success": true,
  "data": {
    "session_id": "test-session-001",
    "answer": "产品A的使用方法如下：首先打开应用，然后点击设置按钮，选择产品A，按照提示完成配置即可。",
    "confidence": 0.85,
    "sources": [
      {
        "document_id": "doc-001",
        "title": "产品A使用指南",
        "relevance": 0.92,
        "snippet": "产品A是一款智能办公软件，支持多人协作、文档管理等功能..."
      }
    ],
    "actions": ["查看详细文档", "联系技术支持"],
    "duration": 1250,
    "timestamp": 1707331200
  }
}
```

### 2. 获取对话历史

```bash
curl http://localhost:8888/api/css/history/test-session-001?limit=10
```

### 3. 获取系统状态

```bash
curl http://localhost:8888/api/css/status
```

**响应示例:**

```json
{
  "success": true,
  "data": {
    "engine": {
      "config": {
        "default_model": "qianwen",
        "temperature": 0.7,
        "max_tokens": 1000,
        "top_k": 5,
        "confidence_threshold": 0.6
      },
      "ai_model": "qianwen",
      "rag_enabled": true,
      "session_count": 5
    },
    "clients": 3,
    "timestamp": 1707331200
  }
}
```

### 4. 关闭会话

```bash
curl -X POST http://localhost:8888/api/css/session/test-session-001/close
```

### 5. WebSocket连接 (实时通讯)

使用WebSocket客户端或浏览器JavaScript连接:

```javascript
const ws = new WebSocket('ws://localhost:8888/api/css/ws?session_id=test-session-001&channel=web');

ws.onopen = () => {
  console.log('Connected to CSS');
  // 发送消息
  ws.send(JSON.stringify({
    type: 'message',
    content: '如何退款？'
  }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Received:', data);
  // 数据结构:
  // {
  //   type: 'message',
  //   role: 'assistant',
  //   content: '购买后7天内可申请退款...',
  //   confidence: 0.82,
  //   sources: [...]
  // }
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('Disconnected');
};
```

## 📊 核心流程详解

### ProcessQuestion 完整流程

```go
answer, err := engine.ProcessQuestion(ctx, sessionID, question)
```

#### 步骤说明:

1. **会话管理** (第78-83行)
   ```go
   session, err := e.conversationMgr.GetOrCreateSession(ctx, sessionID)
   ```
   - 创建新会话或恢复已有会话
   - 更新会话最后活跃时间

2. **保存用户消息** (第85-88行)
   ```go
   e.conversationMgr.SaveMessage(ctx, sessionID, "user", question)
   ```
   - 记录用户提问到数据库

3. **RAG检索知识库** (第90-96行)
   ```go
   docs, err := e.ragRetriever.Search(ctx, question, e.config.TopK)
   ```
   - 向量化查询问题
   - 混合检索(向量0.7 + 关键词0.3)
   - 返回Top-K相关文档

4. **获取对话历史** (第98-103行)
   ```go
   history, err := e.conversationMgr.GetHistory(ctx, sessionID, e.config.MaxHistory)
   ```
   - 获取最近N条消息
   - 用于构建对话上下文

5. **构建Prompt** (第105-106行)
   ```go
   prompt := e.buildPrompt(question, docs, history)
   ```
   - 组合问题、知识库、历史
   - 生成AI输入Prompt

6. **AI生成** (第108-118行)
   ```go
   aiAnswer, err := e.callAI(ctx, prompt)
   ```
   - 调用Eino AI模型
   - 生成回答

7. **置信度评估** (第120-121行)
   ```go
   confidence := e.evaluateConfidence(question, aiAnswer, docs)
   ```
   - 多维度评估:
     * 知识库支持 (+0.15)
     * 回答长度 (太短-0.1, 详细+0.05)
     * 不确定词汇 (-0.15)

8. **判断转接** (第123-139行)
   ```go
   if e.shouldTransfer(ctx, session, question, confidence) {
       return e.transferToHuman(ctx, session, "low_confidence", question)
   }
   ```
   - 条件1: 置信度 < 阈值(0.6)
   - 条件2: 包含转接关键词
   - 条件3: 用户明确要求人工

9. **返回回答** (第141-161行)
   ```go
   answer := &Answer{
       Content: aiAnswer,
       Confidence: confidence,
       Sources: sourceRefs,
       SuggestedActions: actions,
   }
   ```
   - 保存AI回答到历史
   - 更新会话活跃时间
   - 返回完整回答信息

## 🔧 配置说明

### 置信度阈值调整

根据业务需求调整 `CSS_CONFIDENCE_THRESHOLD`:

| 阈值 | 说明 | 适用场景 |
|------|------|----------|
| 0.4-0.5 | 较宽松 | AI能力强的场景 |
| 0.6-0.7 | 平衡 | 推荐值 |
| 0.8-0.9 | 严格 | 高质量要求 |

### RAG Top-K调整

调整 `CSS_TOP_K` 控制检索的文档数量:

- **3-5**: 推荐值,平衡准确性和性能
- **1-2**: 简单问题,快速响应
- **10-15**: 复杂问题,更多信息

### 温度参数调整

调整 `CSS_TEMPERATURE` 控制AI创造性:

| 温度 | 说明 | 适用场景 |
|------|------|----------|
| 0.3-0.5 | 保守,回答确定 | FAQ类问题 |
| 0.7 | 平衡 | 推荐值 |
| 0.8-1.0 | 创造性强 | 开放性对话 |

## 📈 监控指标

### 关键性能指标

1. **平均响应时间**: `duration` 字段
2. **平均置信度**: `confidence` 字段
3. **转接率**: 转接人工的对话比例
4. **会话数量**: 活跃会话数

### 日志监控

系统会输出详细日志:

```
[CSS] Processing question for session test-001: 如何退款？
[CSS] RAG search returned 3 documents
[CSS] AI generation completed in 1250ms
[CSS] Answer generated with confidence 0.85
```

## ⚠️ 常见问题

### 1. AI模型无法连接

检查AI配置:
```bash
AI_ENABLED=true
AI_DEFAULT_PROVIDER=qianwen
```

### 2. RAG检索不工作

检查PostgreSQL是否安装pgvector扩展:
```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

### 3. 置信度总是很低

可能原因:
- 知识库为空,需要上传文档
- 配置的阈值太高
- 问题超出知识库范围

### 4. 转接不生效

检查:
```bash
CSS_TRANSFER_ON_LOW_CONFIDENCE=true
CSS_TRANSFER_KEYWORDS=投诉,退款,人工,转接
```

## 🚀 下一步开发

### 待实现功能

1. **知识库管理**
   - [ ] 文档上传API
   - [ ] PDF/Word解析
   - [ ] 自动分块向量化

2. **人工客服**
   - [ ] 客服队列系统
   - [ ] 智能分配算法
   - [ ] 客服工作台API

3. **统计分析**
   - [ ] 对话统计接口
   - [ ] AI效果评估
   - [ ] 满意度统计

4. **前端集成**
   - [ ] Vue聊天组件
   - [ ] 知识库管理界面
   - [ ] 客服工作台

## 📚 相关文档

- [设计文档](./CSS_SYSTEM_DESIGN.md)
- [实现文档](./CSS_IMPLEMENTATION.md)
- [AI集成文档](./ai.md)

---

**部署完成后,系统即可开始处理用户提问并返回智能回答!**
