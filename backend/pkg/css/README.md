# 智能客服系统 (Customer Service System)

## 🎯 核心流程

```
用户提问 → 会话管理 → RAG检索知识库 → 构建Prompt → AI生成
        → 置信度评估 → (高置信度) 返回回答 / (低置信度) 转人工
```

## 📦 快速开始

### 1. 启用系统

在 `backend/.env` 中设置:
```bash
CSS_ENABLED=true
```

### 2. 运行迁移

```bash
cd backend
go run cmd/artisan/main.go artisan migrate
```

### 3. 启动服务

```bash
go run main.go
```

### 4. 测试API

```bash
curl -X POST http://localhost:8888/api/css/question \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-001",
    "question": "如何使用产品A？"
  }'
```

## 📁 项目结构

```
pkg/css/
├── engine.go              # 核心引擎
├── init.go               # 初始化模块
├── websocket.go          # WebSocket服务
├── engine_test.go        # 单元测试
├── conversation/
│   └── manager.go        # 对话管理
├── rag/
│   ├── retriever.go      # RAG检索器
│   ├── embedder.go       # 向量化器
│   └── vector_store.go   # 向量存储
└── kb/
    └── service.go        # 知识库服务
```

## 🔧 核心组件

### CSSEngine (engine.go)

智能客服核心引擎,负责协调所有组件。

**主要方法:**
- `ProcessQuestion()` - 处理用户提问的完整流程
- `evaluateConfidence()` - 评估回答置信度
- `shouldTransfer()` - 判断是否转接人工

### Manager (conversation/manager.go)

会话和消息管理。

**主要方法:**
- `GetOrCreateSession()` - 获取或创建会话
- `SaveMessage()` - 保存消息
- `GetHistory()` - 获取对话历史

### Retriever (rag/retriever.go)

RAG检索引擎,混合使用向量检索和关键词检索。

**主要方法:**
- `Search()` - 执行混合检索

### WSManager (websocket.go)

WebSocket实时通讯管理。

**主要功能:**
- 连接管理
- 消息泵
- 心跳保活

## 📊 配置说明

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| CSS_ENABLED | 启用开关 | false |
| CSS_DEFAULT_MODEL | AI模型 | qianwen |
| CSS_TEMPERATURE | 温度参数 | 0.7 |
| CSS_MAX_TOKENS | 最大Token数 | 1000 |
| CSS_TOP_K | 检索Top-K | 5 |
| CSS_CONFIDENCE_THRESHOLD | 置信度阈值 | 0.6 |
| CSS_TRANSFER_ON_LOW_CONFIDENCE | 低置信度转接 | true |
| CSS_TRANSFER_KEYWORDS | 转接关键词 | 投诉,退款,人工,转接 |
| CSS_SESSION_TIMEOUT | 会话超时(秒) | 1800 |
| CSS_MAX_HISTORY | 最大历史消息数 | 10 |

## 🎨 API接口

### REST API

| 端点 | 方法 | 功能 |
|------|------|------|
| `/api/css/ws` | GET | WebSocket连接 |
| `/api/css/question` | POST | 发送问题 |
| `/api/css/history/:session_id` | GET | 获取历史 |
| `/api/css/session/:session_id/close` | POST | 关闭会话 |
| `/api/css/status` | GET | 系统状态 |

### WebSocket

连接地址: `ws://localhost:8888/api/css/ws?session_id=xxx`

消息格式:
```json
{
  "type": "message",
  "content": "如何使用产品A？"
}
```

响应格式:
```json
{
  "type": "message",
  "role": "assistant",
  "content": "产品A的使用方法如下...",
  "confidence": 0.85,
  "sources": [...]
}
```

## 🧪 测试

运行单元测试:
```bash
go test ./pkg/css/... -v
```

## 📚 文档

- [设计文档](../../doc/CSS_SYSTEM_DESIGN.md)
- [实现文档](../../doc/CSS_IMPLEMENTATION.md)
- [部署指南](../../doc/CSS_DEPLOYMENT.md)

## ⚠️ 注意事项

1. 确保AI服务已启用并配置正确
2. 如需使用RAG功能,需要PostgreSQL + pgvector扩展
3. 首次使用需要迁移数据库创建表结构
4. 建议先在测试环境验证配置

## 🚀 后续开发

- [ ] 知识库文档上传
- [ ] PDF/Word解析
- [ ] 客服队列系统
- [ ] 统计分析接口
- [ ] 前端集成组件

## 📝 示例

### 高置信度回答

```bash
curl -X POST http://localhost:8888/api/css/question \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-001",
    "question": "如何使用产品A？"
  }'
```

**响应:**
```json
{
  "success": true,
  "data": {
    "session_id": "test-001",
    "answer": "产品A的使用方法如下...",
    "confidence": 0.85,
    "sources": [...],
    "actions": ["查看文档", "联系支持"]
  }
}
```

### 低置信度转接

```bash
curl -X POST http://localhost:8888/api/css/question \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-002",
    "question": "我要投诉这个产品"
  }'
```

**响应:**
```json
{
  "success": true,
  "data": {
    "session_id": "test-002",
    "answer": "您的问题比较复杂,我正在为您转接人工客服...",
    "confidence": 0.55,
    "transfer_to": "agent"
  }
}
```

---

**系统已就绪,开始使用智能客服吧!** 🎉
