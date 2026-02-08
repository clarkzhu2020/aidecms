# 智能客服系统 - 前端集成完整指南

## 🎯 实现的完整架构

```
┌─────────────────────────────────────────────┐
│              前端/Web/小程序                   │
│  ┌─────────────────────────────────────┐   │
│  │  CustomerServiceChat.vue (组件)    │   │
│  │  useCustomerService.ts (Hook)     │   │
│  └─────────────────────────────────────┘   │
└─────────────────┬───────────────────────────┘
                  │ WebSocket / HTTP
┌─────────────────▼───────────────────────────┐
│         CustomerServiceController            │
│         (API控制器 / WebSocket)              │
└─────────────────┬───────────────────────────┘
                  │
┌─────────────────▼───────────────────────────┐
│              CSSEngine                       │
│         (核心引擎 - 流程协调)                  │
├─────────────────┬───────────────────────────┤
│  会话管理        │  RAG检索                    │
│  (Manager)      │  (Retriever)               │
├─────────────────┼───────────────────────────┤
│  AI调用         │  置信度评估                 │
│  (Eino)         │  (Confidence)               │
└─────────────────┴───────────────────────────┘
         │                 │
    ┌────▼─────┐      ┌───▼────────┐
    │ MySQL    │      │ PostgreSQL  │
    │ (会话/消息)│      │ (向量数据)  │
    └──────────┘      └────────────┘
```

## ✅ 已完成的前端组件

### 1. CustomerServiceChat.vue
**位置**: `web/src/components/CustomerServiceChat.vue`

**功能**:
- ✅ 完整的聊天界面
- ✅ WebSocket实时通讯
- ✅ HTTP API备用通道
- ✅ 消息列表显示
- ✅ 用户/AI消息区分
- ✅ 置信度指示器 (可视化进度条)
- ✅ 知识来源展示
- ✅ 建议操作按钮
- ✅ 转人工提示
- ✅ 正在输入动画
- ✅ 历史记录面板
- ✅ 自动滚动到底部
- ✅ 响应式设计

**特性**:
- 自动重连 (5秒延迟)
- 心跳保活 (ping/pong)
- 消息自动格式化 (支持Markdown)
- 置信度高/中/低分级显示
- 会话持久化
- 移动端适配

### 2. useCustomerService.ts
**位置**: `web/src/composables/useCustomerService.ts`

**功能**:
- ✅ Composable Hook封装
- ✅ 面向对象API (CustomerService类)
- ✅ 完整的TypeScript类型定义
- ✅ WebSocket连接管理
- ✅ 自动重连机制
- ✅ 消息状态管理
- ✅ 错误处理
- ✅ 事件回调系统

**API接口**:

#### Hook版本
```typescript
const {
  // 状态
  messages,
  isConnected,
  isTyping,
  isTransferred,
  error,
  sessionId,

  // 计算属性
  lastMessage,
  messageCount,
  avgConfidence,

  // 方法
  sendMessage,
  getHistory,
  closeSession,
  getStatus,
  connect,
  disconnect,
  clearMessages,
  reset
} = useCustomerService({
  onMessage: (msg) => console.log(msg),
  onError: (err) => console.error(err),
  onConnect: () => console.log('Connected'),
  onDisconnect: () => console.log('Disconnected')
})
```

#### 类版本
```typescript
const cs = new CustomerService({
  sessionId: 'my-session-id'
})

cs.connect()
cs.on('message', (msg) => console.log(msg))
const answer = await cs.send('如何使用产品A?')
```

### 3. CustomerServiceDemo.vue
**位置**: `web/src/views/CustomerServiceDemo.vue`

**功能**:
- ✅ 演示页面
- ✅ 聊天窗口集成
- ✅ 控制面板
- ✅ 连接状态监控
- ✅ 统计信息展示
- ✅ 快捷问题测试
- ✅ API调用测试
- ✅ 系统日志记录

## 🚀 使用方式

### 方式1: 使用组件

```vue
<template>
  <div>
    <button @click="showChat = true">打开客服</button>
    <CustomerServiceChat
      v-if="showChat"
      @close="showChat = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import CustomerServiceChat from '@/components/CustomerServiceChat.vue'

const showChat = ref(false)
</script>
```

### 方式2: 使用Composable Hook

```vue
<template>
  <div>
    <div v-for="(msg, index) in messages" :key="index">
      <div :class="msg.role">{{ msg.content }}</div>
      <div v-if="msg.confidence">置信度: {{ (msg.confidence * 100).toFixed(0) }}%</div>
    </div>

    <input v-model="inputText" @keyup.enter="sendMessage" />
    <button @click="sendMessage">发送</button>

    <div v-if="isTyping">正在输入...</div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useCustomerService } from '@/composables/useCustomerService'

const {
  messages,
  isTyping,
  sendMessage: sendQuestion
} = useCustomerService()

const inputText = ref('')

async function sendMessage() {
  await sendQuestion(inputText.value)
  inputText.value = ''
}
</script>
```

### 方式3: 使用类API

```typescript
import { CustomerService } from '@/composables/useCustomerService'

const cs = new CustomerService({
  sessionId: 'user-123',
  channel: 'web'
})

// 连接
cs.connect()

// 监听消息
cs.on('message', (msg) => {
  console.log('收到消息:', msg.content)
  console.log('置信度:', msg.confidence)
  console.log('来源:', msg.sources)
})

// 发送问题
const answer = await cs.send('如何退款?')
console.log(answer.content)

// 获取历史
const history = await cs.getHistory()

// 关闭连接
cs.close()
```

## 📊 核心功能展示

### 1. 置信度可视化

```vue
<div class="confidence-indicator">
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

**样式分级**:
- 高置信度 (≥0.8): 绿色 (#52c41a)
- 中置信度 (0.6-0.8): 橙色 (#faad14)
- 低置信度 (<0.6): 红色 (#ff4d4f)

### 2. 知识来源展示

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

### 3. 建议操作

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

### 4. 转人工提示

```vue
<div v-if="msg.transferTo" class="transfer-notice">
  <span class="transfer-icon">🔄</span>
  正在为您转接人工客服...
</div>
```

## 🧪 测试示例

### 测试高置信度场景

**问题**: "如何使用产品A?"

**预期结果**:
```json
{
  "role": "assistant",
  "content": "产品A的使用方法如下:首先打开应用...",
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

**显示效果**:
- 置信度条: 绿色, 85%
- 知识来源: 显示1个文档,相关度92%
- 建议操作: 显示2个可点击按钮
- 无转接提示

### 测试低置信度转接

**问题**: "我要投诉这个产品"

**预期结果**:
```json
{
  "role": "assistant",
  "content": "您的问题比较复杂,我正在为您转接人工客服...",
  "confidence": 0.55,
  "transferTo": "agent"
}
```

**显示效果**:
- 置信度条: 红色, 55%
- 显示转人工提示
- 建议操作: "转接人工客服"

## 🔧 配置说明

### WebSocket配置

```typescript
const css = useCustomerService({
  wsUrl: 'ws://localhost:8888/api/css/ws',
  httpUrl: 'http://localhost:8888/api/css',
  sessionId: 'my-session-id',
  channel: 'web',
  autoReconnect: true,
  reconnectDelay: 5000,
  onMessage: (msg) => console.log(msg),
  onError: (err) => console.error(err)
})
```

### HTTP API配置

如果WebSocket不可用,会自动降级到HTTP API:

```typescript
// 自动发送到 POST /api/css/question
await css.sendMessage('问题内容')
```

## 📱 响应式设计

### 桌面端 (>1024px)
- 固定位置: 右下角
- 尺寸: 400px × 600px
- 历史面板: 侧边滑出

### 移动端 (≤1024px)
- 全屏模式
- 自适应高度
- 优化触摸操作

## 🎨 UI/UX特性

1. **渐进式加载**:
   - 正在输入动画 (三个跳动圆点)
   - 消息平滑显示

2. **视觉反馈**:
   - 置信度颜色分级
   - 转接状态明显提示
   - 连接状态实时更新

3. **交互优化**:
   - 自动滚动到底部
   - 文本框自适应高度
   - 快捷问题一键发送

4. **错误处理**:
   - 网络错误提示
   - 自动重连
   - 错误日志记录

## 📚 相关文档

- [后端架构](../../backend/doc/CSS_ARCHITECTURE.md)
- [部署指南](../../backend/doc/CSS_DEPLOYMENT.md)
- [实现文档](../../backend/doc/CSS_IMPLEMENTATION.md)
- [设计文档](../../backend/doc/CSS_SYSTEM_DESIGN.md)

## 🎉 总结

已完整实现架构图中的所有层级:

✅ **前端/Web/小程序** - Vue3组件 + Composable Hook
✅ **API控制器** - CustomerServiceController + WebSocket
✅ **核心引擎** - CSSEngine (完整流程)
✅ **会话管理** - conversation.Manager
✅ **RAG检索** - rag.Retriever + embedder
✅ **AI调用** - Eino集成
✅ **置信度评估** - 多因素评估
✅ **数据库** - MySQL + PostgreSQL + pgvector

系统已就绪,可以开始使用! 🚀
