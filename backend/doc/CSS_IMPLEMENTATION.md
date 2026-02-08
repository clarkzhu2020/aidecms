# 智能客服系统实现文档

## 📦 已实现模块

### 1. 核心引擎 (`pkg/css/engine.go`)

**文件路径**: `backend/pkg/css/engine.go`

**功能**: 实现智能客服核心引擎，整合AI、RAG和对话管理

**核心流程**:
```
用户提问 → 会话管理 → RAG检索知识库 → 构建Prompt → AI生成
        → 置信度评估 → (高置信度) 返回回答 / (低置信度) 转人工
```

**主要方法**:
- `ProcessQuestion()` - 处理用户提问的完整流程
- `buildPrompt()` - 构建AI Prompt（问题+上下文+知识库）
- `callAI()` - 调用AI模型生成回答
- `evaluateConfidence()` - 评估AI回答置信度
- `shouldTransfer()` - 判断是否需要转接人工
- `transferToHuman()` - 转接人工客服
- `suggestActions()` - 基于问题类型建议操作

**置信度评估因素**:
1. 是否有知识库支持 (+0.15)
2. 回答长度是否合理 (太短-0.1, 详细+0.05)
3. 回答是否包含不确定词汇 (-0.15)

---

### 2. RAG检索模块 (`pkg/css/rag/`)

#### 2.1 检索器 (`retriever.go`)

**文件路径**: `backend/pkg/css/rag/retriever.go`

**功能**: RAG检索引擎，混合使用向量检索和关键词检索

**核心方法**:
- `Search()` - 执行RAG检索
  1. 将查询向量化
  2. 向量相似度检索 (权重 0.7)
  3. 关键词全文检索 (权重 0.3)
  4. 结果融合和重排序

#### 2.2 向量化器 (`embedder.go`)

**文件路径**: `backend/pkg/css/rag/embedder.go`

**功能**: 将文本转换为嵌入向量

**核心方法**:
- `Embed()` - 单文本向量化
- `EmbedBatch()` - 批量向量化
- `GetDimension()` - 获取向量维度 (1536)

#### 2.3 向量存储 (`vector_store.go`)

**文件路径**: `backend/pkg/css/rag/vector_store.go`

**功能**: 存储和检索文档向量

**核心方法**:
- `Store()` - 存储向量 (支持pgvector)
- `Search()` - 向量相似度搜索 (余弦距离)
- `CreateVectorIndex()` - 创建向量索引 (IVFFlat)
- `GetChunkCount()` - 获取分块总数

**技术依赖**: PostgreSQL + pgvector扩展

---

### 3. 对话管理模块 (`pkg/css/conversation/`)

**文件路径**: `backend/pkg/css/conversation/manager.go`

**功能**: 管理会话和消息历史

**核心方法**:
- `GetOrCreateSession()` - 获取或创建会话
- `SaveMessage()` - 保存消息
- `GetHistory()` - 获取对话历史
- `UpdateSessionStatus()` - 更新会话状态
- `AssignSessionToAgent()` - 分配会话给客服
- `GetActiveSessions()` - 获取活跃会话列表

**会话缓存**: 使用map缓存活跃会话，提高性能

---

### 4. 知识库管理 (`pkg/css/kb/`)

**文件路径**: `backend/pkg/css/kb/service.go`

**功能**: 知识库文档管理

**核心方法**:
- `CreateDocument()` - 创建文档
- `GetDocuments()` - 分页获取文档列表
- `FullTextSearch()` - 全文搜索
- `GetCategories()` - 获取分类列表
- `UpdateDocumentStatus()` - 更新文档状态 (processing/indexed/error)

**支持的查询**: LIKE模糊搜索、分类筛选、分页

---

### 5. WebSocket服务 (`pkg/css/websocket.go`)

**文件路径**: `backend/pkg/css/websocket.go`

**功能**: 实时双向通讯

**核心功能**:
- 连接管理 (注册/注销)
- 消息泵 (readPump/writePump)
- 心跳保活 (54秒间隔)
- 会话隔离 (session_id)

**消息类型**:
- `connected` - 连接建立
- `message` - 普通消息
- `typing` - 正在输入
- `ping` - 心跳
- `error` - 错误消息

**处理流程**:
```
客户端连接 → 注册 → readPump(接收) → handleClientMessage → 
ProcessQuestion → writePump(发送)
```

---

### 6. 数据库模型 (`internal/app/models/css_session.go`)

**文件路径**: `backend/internal/app/models/css_session.go`

**数据表**:
1. `css_sessions` - 会话表
2. `css_messages` - 消息表
3. `kb_documents` - 知识库文档表
4. `kb_chunks` - 文档分块表 (含vector)
5. `cs_agents` - 客服人员表
6. `cs_transfers` - 转接记录表
7. `css_feedback` - 用户反馈表

**索引策略**:
- 会话: (user_id), (channel, channel_id), (status), (status, last_active_at)
- 消息: (session_id, created_at), (role)
- 文档: (category), (status)
- 分块: (document_id, chunk_order)

---

### 7. 控制器 (`app/Http/Controllers/CustomerServiceController.go`)

**文件路径**: `backend/app/Http/Controllers/CustomerServiceController.go`

**API端点**:

| 端点 | 方法 | 功能 |
|--------|------|------|
| `/api/css/ws` | GET | WebSocket连接 |
| `/api/css/question` | POST | 发送问题 |
| `/api/css/history/:session_id` | GET | 获取对话历史 |
| `/api/css/session/:session_id/close` | POST | 关闭会话 |
| `/api/css/status` | GET | 获取系统状态 |

---

### 8. 路由配置 (`routes/api.go`)

**已添加路由**:
```go
// Customer Service (智能客服系统)
r.GET("/api/css/ws", adapters.HertzToFramework(csController.WebSocket))
r.POST("/api/css/question", adapters.HertzToFramework(csController.SendQuestion))
r.GET("/api/css/history/:session_id", adapters.HertzToFramework(csController.GetHistory))
r.POST("/api/css/session/:session_id/close", adapters.HertzToFramework(csController.CloseSession))
r.GET("/api/css/status", adapters.HertzToFramework(csController.GetStatus))
```

---

### 9. 环境变量配置 (`.env`)

**已添加配置项**:
```bash
# 客服系统配置
CSS_ENABLED=false
CSS_DEFAULT_MODEL=qianwen
CSS_TEMPERATURE=0.7
CSS_MAX_TOKENS=1000
CSS_TOP_K=5
CSS_CONFIDENCE_THRESHOLD=0.6
CSS_TRANSFER_ON_LOW_CONFIDENCE=true
CSS_MAX_RETRIES=3
CSS_TRANSFER_KEYWORDS=投诉,退款,人工,转接
CSS_SESSION_TIMEOUT=1800  # 30分钟（秒）
CSS_MAX_HISTORY=10
```

### 10. 初始化模块 (`pkg/css/init.go`)

**文件路径**: `backend/pkg/css/init.go`

**功能**: 统一初始化所有客服系统组件

**核心方法**:
- `InitCustomerServiceSystem()` - 初始化客服系统
  - 初始化对话管理器
  - 初始化RAG检索器
  - 初始化知识库服务
  - 初始化核心引擎
  - 初始化WebSocket管理器
  - 自动加载环境变量配置

**在路由中的使用**:
```go
// routes/api.go
if getConfigBool("CSS_ENABLED", false) {
    engine, wsManager := css.InitCustomerServiceSystem(app.DB, manager)
    csController := controllers.NewCustomerServiceController()
    csController.Init(engine, wsManager)
}
```

---

## 🚀 项目结构

```
backend/
├── pkg/
│   ├── css/
│   │   ├── engine.go              # 核心引擎
│   │   ├── websocket.go           # WebSocket服务
│   │   ├── conversation/
│   │   │   └── manager.go       # 对话管理
│   │   ├── rag/
│   │   │   ├── retriever.go    # RAG检索器
│   │   │   ├── embedder.go     # 向量化器
│   │   │   └── vector_store.go # 向量存储
│   │   └── kb/
│   │       └── service.go        # 知识库服务
│   └── ai/                       # 已有AI模块
├── internal/app/models/
│   └── css_session.go          # 数据库模型
├── app/Http/Controllers/
│   └── CustomerServiceController.go # 客服控制器
├── routes/
│   └── api.go                    # 路由配置 (已更新)
└── .env                          # 环境变量 (已更新)
```

---

## 📊 核心流程实现

### 用户提问完整流程

```go
func (e *CSSEngine) ProcessQuestion(ctx context.Context, sessionID, question string) (*Answer, error) {
    // 1. 会话管理
    session, err := e.conversationMgr.GetOrCreateSession(ctx, sessionID)
    
    // 2. 保存用户消息
    e.conversationMgr.SaveMessage(ctx, sessionID, "user", question)
    
    // 3. RAG检索知识库
    docs, err := e.ragRetriever.Search(ctx, question, e.config.TopK)
    
    // 4. 获取对话历史
    history, err := e.conversationMgr.GetHistory(ctx, sessionID, e.config.MaxHistory)
    
    // 5. 构建Prompt
    prompt := e.buildPrompt(question, docs, history)
    
    // 6. AI生成
    aiAnswer, err := e.callAI(ctx, prompt)
    
    // 7. 置信度评估
    confidence := e.evaluateConfidence(question, aiAnswer, docs)
    
    // 8. 判断转接条件
    if e.shouldTransfer(ctx, session, question, confidence) {
        // 转人工
        return e.transferToHuman(ctx, session, "low_confidence", question)
    }
    
    // 9. 保存AI回答
    e.conversationMgr.SaveMessage(ctx, sessionID, "assistant", aiAnswer)
    
    // 10. 返回回答
    return &Answer{
        Content: aiAnswer,
        Confidence: confidence,
        Sources: buildSourceRefs(docs),
    }, nil
}
```

---

## 🔧 使用说明

### 1. 启用客服系统

在 `.env` 文件中设置：
```bash
CSS_ENABLED=true
```

### 2. 配置AI模型

确保AI配置正确：
```bash
AI_ENABLED=true
AI_DEFAULT_PROVIDER=qianwen
```

### 3. 配置向量数据库

需要安装PostgreSQL和pgvector扩展：
```sql
CREATE EXTENSION IF NOT EXISTS vector;
```

### 4. 测试API

```bash
# WebSocket连接
ws://localhost:8888/api/css/ws?session_id=test123

# 发送问题
curl -X POST http://localhost:8888/api/css/question \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test123",
    "question": "如何使用产品A？"
  }'

# 获取历史
curl http://localhost:8888/api/css/history/test123
```

---

## ⚠️ 待完成功能

### Phase 1: 基础功能 (已完成 ✅)
- [x] 核心引擎实现
- [x] RAG检索模块
- [x] 对话管理系统
- [x] 知识库服务
- [x] WebSocket服务
- [x] 数据库模型
- [x] API控制器
- [x] 路由配置
- [x] 环境变量

### Phase 2: 扩展功能 (待实现 ⏳)
- [ ] 知识库文档上传接口
- [ ] 文档解析器 (PDF/Word)
- [ ] 文本分块和向量化流程
- [ ] 客服人员管理
- [ ] 转接队列系统
- [ ] 统计分析接口
- [ ] 满意度反馈
- [ ] 管理端知识库管理界面
- [ ] 客服工作台界面

### Phase 3: 优化和完善 (待实现 ⏳)
- [ ] 性能优化（缓存、连接池）
- [ ] 安全加固（敏感词过滤、权限控制）
- [ ] 多渠道支持（微信、小程序）
- [ ] 监控和日志系统
- [ ] 单元测试和集成测试

---

## 📝 技术要点

### 1. 并发安全

- WebSocket管理器使用 `sync.RWMutex`
- 会话缓存使用读写锁
- 消息通道缓冲大小为256

### 2. 错误处理

- 所有关键操作都有错误处理
- 使用 `hlog` 记录详细日志
- WebSocket连接异常时自动清理

### 3. 性能优化

- 活跃会话缓存
- 向量检索使用索引
- 消息批量保存

### 4. 扩展性设计

- 模块化设计，易于扩展
- 接口抽象，支持替换实现
- 配置驱动，无需修改代码

---

## 🎯 下一步行动

### 立即执行
1. **启用客服系统**
   ```bash
   # 在 .env 中设置
   CSS_ENABLED=true
   ```

2. **数据库迁移**
   ```bash
   go run cmd/artisan/main.go artisan migrate
   ```

3. **启动服务**
   ```bash
   go run main.go
   ```

4. **测试核心流程**
   - 测试REST API: `POST /api/css/question`
   - 测试WebSocket连接
   - 测试转接逻辑

5. **完善知识库**
   - 实现文档上传
   - 实现文档解析
   - 实现向量化流程

6. **前端开发**
   - 聊天窗口组件
   - 知识库管理界面
   - 客服工作台

### 测试示例

#### 测试问题提交
```bash
curl -X POST http://localhost:8888/api/css/question \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-001",
    "question": "如何使用产品A？"
  }'
```

#### 测试转接逻辑
```bash
curl -X POST http://localhost:8888/api/css/question \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-002",
    "question": "我要投诉这个产品"
  }'
```

预期: 由于包含"投诉"关键词,会触发转人工逻辑。

---

## 📚 相关文档

- [设计文档](./CSS_SYSTEM_DESIGN.md)
- [部署指南](./CSS_DEPLOYMENT.md)
- [AI集成文档](./ai.md)
- [邮件系统文档](./mail.md)
- [KuCoin集成](./KUCOIN_INTEGRATION.md)

---

**文档版本**: v1.0  
**创建日期**: 2025-02-07  
**更新日期**: 2025-02-07  
**状态**: 核心功能实现完成
