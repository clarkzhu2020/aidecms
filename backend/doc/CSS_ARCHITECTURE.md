# 智能客服系统 - 架构图说明

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────┐
│              前端/Web/小程序                   │
│         (Vue3 / React / 微信小程序)          │
│                                           │
│  ┌─────────────────────────────────────┐   │
│  │     CustomerServiceChat.vue        │   │
│  │     (聊天窗口组件)                   │   │
│  │  ┌───────────────────────────────┐  │   │
│  │  │   useCustomerService()       │  │   │
│  │  │   (Composable Hook)          │  │   │
│  │  └───────────────────────────────┘  │   │
│  │     │                              │   │
│  │     ├─ WebSocket (实时通讯)          │   │
│  │     └─ HTTP (备用通道)              │   │
│  └─────────────────────────────────────┘   │
└─────────────────┬───────────────────────────┘
                  │ WebSocket / HTTP
                  │ (JSON格式)
┌─────────────────▼───────────────────────────┐
│         CustomerServiceController            │
│         (API控制器 / WebSocket处理)            │
│                                           │
│  ┌─────────────────────────────────────┐   │
│  │  WebSocket Handler               │   │
│  │  - 连接管理                       │   │
│  │  - 消息路由                       │   │
│  │  - 心跳保活                       │   │
│  └─────────────────────────────────────┘   │
│                                           │
│  ┌─────────────────────────────────────┐   │
│  │  REST API Endpoints              │   │
│  │  POST /api/css/question          │   │
│  │  GET  /api/css/history/:id       │   │
│  │  POST /api/css/session/:id/close │   │
│  │  GET  /api/css/status            │   │
│  └─────────────────────────────────────┘   │
└─────────────────┬───────────────────────────┘
                  │ 调用引擎
┌─────────────────▼───────────────────────────┐
│              CSSEngine                       │
│         (核心引擎 - 流程协调)                  │
│                                           │
│  ProcessQuestion() 完整流程:                 │
│                                           │
│  1. 会话管理 ─────────────────────────┐      │
│     │                             │      │
│     ▼                             │      │
│  ┌──────────────────┐              │      │
│  │ conversation     │              │      │
│  │  .Manager       │              │      │
│  │                 │              │      │
│  │ - GetOrCreate   │              │      │
│  │ - SaveMessage   │              │      │
│  │ - GetHistory    │              │      │
│  └────────┬─────────┘              │      │
│           │                        │      │
│           ▼                        │      │
│  ┌──────────────────┐              │      │
│  │ MySQL / SQLite  │              │      │
│  │ (会话/消息表)     │              │      │
│  └──────────────────┘              │      │
│                                    │      │
│  2. RAG检索 ────────────────────┐  │      │
│     │                          │  │      │
│     ▼                          │  │      │
│  ┌──────────────────┐          │  │      │
│  │ rag.Retriever   │          │  │      │
│  │                 │          │  │      │
│  │ Search()        │          │  │      │
│  │   ├─ 向量检索(0.7)│          │  │      │
│  │   └─ 关键词(0.3)│          │  │      │
│  └────────┬─────────┘          │  │      │
│           │                    │  │      │
│           ├────────────────────┼┘      │
│           ▼                    │        │
│  ┌──────────────────┐          │        │
│  │ rag.Embedder    │          │        │
│  │                 │          │        │
│  │ Embed(text)     │          │        │
│  │   ──► [0.1, 0.2]│          │        │
│  └────────┬─────────┘          │        │
│           │                    │        │
│           ▼                    ▼        │
│  ┌──────────────────┐  ┌──────────────┐│
│  │ PostgreSQL       │  │ AI Service   ││
│  │ + pgvector      │  │ (Eino)       ││
│  │                 │  │              ││
│  │ kb_chunks       │  │ - 向量化      ││
│  │ (向量数据)        │  │ - 模型推理     ││
│  │                 │  │              ││
│  │ IVFFlat索引     │  │ qianwen/    ││
│  │                 │  │ gpt/etc     ││
│  └──────────────────┘  └──────────────┘│
│                                             │
│  3. 构建Prompt ─────────────────────────│    │
│     │                                  │    │
│     ▼                                  │    │
│  ┌────────────────────────────────────┐  │    │
│  │ Prompt =                          │  │    │
│  │   - System: 你是专业客服助手...    │  │    │
│  │   - Knowledge: 文档1, 文档2...    │  │    │
│  │   - History: 用户:问题, AI:回答... │  │    │
│  │   - Question: 当前问题            │  │    │
│  └────────────┬───────────────────────┘  │    │
│               │                             │    │
│               ▼                             │    │
│  ┌──────────────────┐                     │    │
│  │ AI Call         │                     │    │
│  │                 │                     │    │
│  │ callAI(prompt)  │─────────────────────┘    │
│  │                 │                          │
│  │ temperature:0.7 │                          │
│  │ max_tokens:1000 │                          │
│  └────────┬─────────┘                          │
│           │                                      │
│           ▼                                      │
│     Answer (文本回答)                             │
│                                             │
│  4. 置信度评估 ────────────────────────┐      │
│     │                                 │      │
│     ▼                                 │      │
│  ┌──────────────────┐                  │      │
│  │ evaluateConfidence()               │      │
│  │                 │                  │      │
│  │ 因素:           │                  │      │
│  │ - 知识库支持    │                  │      │
│  │   +0.15 / -0.2 │                  │      │
│  │ - 回答长度      │                  │      │
│  │   -0.1 / +0.05 │                  │      │
│  │ - 不确定词汇    │                  │      │
│  │   -0.15        │                  │      │
│  └────────┬─────────┘                  │      │
│           │                            │      │
│           ▼                            │      │
│     confidence score (0-1)              │      │
│                                          │
│  5. 判断转接 ─────────────────────────┐  │    │
│     │                                 │  │    │
│     ├─ confidence < 0.6? ──────────┐ │  │    │
│     │    Yes ──► 转人工            │ │  │    │
│     │    No  ──► 返回AI回答         │ │  │    │
│     │                                 │ │  │    │
│     ├─ 关键词匹配? ────────────────┐ │ │  │    │
│     │    (投诉/退款/人工/转接)       │ │ │  │    │
│     │    Yes ──► 转人工            │ │ │  │    │
│     │    No  ──► 继续             │ │ │  │    │
│     │                                 │ │  │    │
│     └─ 用户明确要求? ──────────────┐ │ │  │    │
│        ("我要人工服务")              │ │ │  │    │
│        Yes ──► 转人工             │ │ │  │    │
│        No  ──► 继续              │ │ │  │    │
│                                       │ │  │    │
│  6. 返回结果 ───────────────────────┼─┘ │    │
│     │                                 │      │
│     ▼                                 │      │
│  ┌──────────────────┐                 │      │
│  │ Answer          │                 │      │
│  │                 │                 │      │
│  │ {               │                 │      │
│  │   content,      │                 │      │
│  │   confidence,   │                 │      │
│  │   sources [],   │                 │      │
│  │   actions [],   │                 │      │
│  │   transfer_to   │                 │      │
│  │ }               │                 │      │
│  └────────┬─────────┘                 │      │
│           │                            │      │
│           └────────────────────────────┘      │
│                                           │
└───────────────────────────────────────────┘
         │                                       │
         │ 存储消息                               │
         ▼                                       ▼
    ┌──────────┐                         ┌──────────────┐
    │  MySQL   │                         │ PostgreSQL  │
    │          │                         │             │
    │ 会话数据  │                         │ 向量数据     │
    │ 消息历史  │                         │ 知识库       │
    │ 转接记录  │                         │ 文档分块     │
    └──────────┘                         └──────────────┘
```

## 🔄 数据流向

### 1. 用户提问流程

```
用户输入 → 前端组件 → WebSocket/HTTP
    ↓
CustomerServiceController
    ↓
CSSEngine.ProcessQuestion()
    ↓
├─ 1. 会话管理 (Manager)
│     └─ MySQL (会话/消息表)
├─ 2. RAG检索 (Retriever)
│     ├─ Embedder (向量化)
│     └─ PostgreSQL + pgvector (向量检索)
├─ 3. 构建Prompt (问题+知识库+历史)
├─ 4. AI调用 (Eino)
├─ 5. 置信度评估 (多因素)
├─ 6. 判断转接 (阈值/关键词)
└─ 7. 返回结果 (Answer)
    ↓
前端显示
```

### 2. WebSocket实时通讯

```
┌─────────────┐       WebSocket       ┌─────────────┐
│   前端      │ ◄──────────────────► │   后端      │
│             │                       │             │
│  发送消息    │  {                  │  接收消息    │
│  {          │    type: 'message',  │  {          │
│    type:    │    content: '...'    │    type:    │
│    'message',│  }                  │    'message',│
│    content:  │                     │    role:    │
│    '...'     │  ◄────────────────► │    'assistant'│
│  }          │  {                  │    content:  │
│             │    type: 'message',  │    '...'     │
│             │    role: 'assistant' │    confidence│
│             │  }                  │    sources:[]│
└─────────────┘                       └─────────────┘
```

### 3. RAG检索流程

```
用户问题
    ↓
文本向量化 (Embedder)
    ↓
[向量: 0.12, 0.45, 0.78, ...]  (1536维)
    ↓
向量相似度检索 (PostgreSQL + pgvector)
    │
    ├─ 向量检索 (权重 0.7)
    │   └─ IVFFlat索引加速
    │
    └─ 关键词检索 (权重 0.3)
        └─ LIKE模糊匹配
    ↓
结果融合 + 重排序
    ↓
Top-K 相关文档
    ↓
返回给引擎
```

## 🗄️ 数据库设计

### MySQL数据库

```sql
-- 会话表
CREATE TABLE css_sessions (
    id VARCHAR(36) PRIMARY KEY,
    user_id BIGINT UNSIGNED,
    channel VARCHAR(50) NOT NULL,
    channel_id VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active',
    assigned_to BIGINT UNSIGNED,
    confidence FLOAT DEFAULT 0.8,
    message_count INT DEFAULT 0,
    last_active_at DATETIME,
    transfer_reason VARCHAR(255),
    metadata TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    INDEX idx_user (user_id),
    INDEX idx_channel (channel, channel_id),
    INDEX idx_status (status),
    INDEX idx_status_active (status, last_active_at)
);

-- 消息表
CREATE TABLE css_messages (
    id VARCHAR(36) PRIMARY KEY,
    session_id VARCHAR(36) NOT NULL,
    role VARCHAR(20) NOT NULL,
    content TEXT NOT NULL,
    tokens INT DEFAULT 0,
    duration INT,
    confidence FLOAT,
    document_refs TEXT,
    metadata TEXT,
    created_at DATETIME,
    INDEX idx_session_created (session_id, created_at),
    INDEX idx_role (role)
);
```

### PostgreSQL数据库 (向量存储)

```sql
-- 启用pgvector扩展
CREATE EXTENSION IF NOT EXISTS vector;

-- 知识库文档表
CREATE TABLE kb_documents (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    category VARCHAR(100),
    status VARCHAR(20) DEFAULT 'processing',
    metadata JSONB,
    created_at DATETIME,
    updated_at DATETIME,
    INDEX idx_category (category),
    INDEX idx_status (status)
);

-- 文档分块表 (含向量)
CREATE TABLE kb_chunks (
    id VARCHAR(36) PRIMARY KEY,
    document_id VARCHAR(36) NOT NULL,
    chunk_order INT NOT NULL,
    content TEXT NOT NULL,
    embedding vector(1536),  -- pgvector类型
    metadata JSONB,
    created_at DATETIME,
    INDEX idx_document (document_id, chunk_order),
    INDEX idx_embedding (embedding vector_cosine_ops)  -- 向量索引
);

-- 创建向量索引 (IVFFlat)
CREATE INDEX idx_embedding_ivfflat ON kb_chunks
USING ivfflat (embedding vector_cosine_ops)
WITH (lists = 100);
```

## 🔌 API接口

### REST API

| 端点 | 方法 | 功能 | 请求 | 响应 |
|------|------|------|------|------|
| `/api/css/question` | POST | 发送问题 | `{ session_id, question, channel }` | `{ answer, confidence, sources, actions }` |
| `/api/css/history/:id` | GET | 获取历史 | - | `[{ role, content, created_at }]` |
| `/api/css/session/:id/close` | POST | 关闭会话 | - | `{ success, message }` |
| `/api/css/status` | GET | 系统状态 | - | `{ engine, clients, timestamp }` |

### WebSocket协议

#### 客户端 → 服务器

```json
{
  "type": "message",
  "content": "如何使用产品A？"
}
```

#### 服务器 → 客户端

```json
{
  "type": "message",
  "role": "assistant",
  "content": "产品A的使用方法如下...",
  "confidence": 0.85,
  "sources": [
    {
      "document_id": "doc-001",
      "title": "产品A使用指南",
      "relevance": 0.92,
      "snippet": "产品A是一款智能办公软件..."
    }
  ],
  "actions": ["查看文档", "联系支持"],
  "transfer_to": null
}
```

## ⚙️ 配置参数

| 参数 | 说明 | 默认值 | 影响 |
|------|------|--------|------|
| `CSS_ENABLED` | 启用开关 | false | 系统是否初始化 |
| `CSS_DEFAULT_MODEL` | AI模型 | qianwen | 使用的AI模型 |
| `CSS_TEMPERATURE` | 温度参数 | 0.7 | AI创造性 |
| `CSS_MAX_TOKENS` | 最大Token | 1000 | 回答长度限制 |
| `CSS_TOP_K` | 检索数量 | 5 | RAG返回文档数 |
| `CSS_CONFIDENCE_THRESHOLD` | 置信度阈值 | 0.6 | 转人工触发点 |
| `CSS_TRANSFER_ON_LOW_CONFIDENCE` | 低置信度转接 | true | 自动转接 |
| `CSS_TRANSFER_KEYWORDS` | 转接关键词 | 投诉,退款,人工,转接 | 关键词触发 |
| `CSS_SESSION_TIMEOUT` | 会话超时 | 1800秒 | 会话保活时间 |
| `CSS_MAX_HISTORY` | 最大历史 | 10条 | 上下文长度 |

## 📊 性能指标

### 关键指标

- **响应时间**: < 2秒 (RAG + AI)
- **准确率**: > 85% (基于置信度)
- **转接率**: < 15% (正常情况)
- **并发支持**: 1000+ 会话/分钟

### 优化策略

1. **会话缓存**: 活跃会话内存缓存
2. **向量索引**: IVFFlat加速检索
3. **消息队列**: 缓冲区256条
4. **连接池**: 数据库连接复用

## 🔐 安全考虑

1. **数据加密**: HTTPS/WSS传输
2. **权限控制**: JWT认证
3. **敏感词过滤**: 输入输出过滤
4. **审计日志**: 所有操作记录

---

**架构图完整说明了从前端到后端,从用户提问到AI回答的完整流程!**
