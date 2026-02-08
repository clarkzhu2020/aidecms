<template>
  <div class="customer-service-chat">
    <!-- 聊天窗口 -->
    <div class="chat-container">
      <!-- 头部 -->
      <div class="chat-header">
        <div class="agent-info">
          <div class="agent-avatar">
            <svg v-if="!isTransferred" width="40" height="40" viewBox="0 0 40 40">
              <circle cx="20" cy="20" r="20" fill="#1890ff"/>
              <text x="20" y="28" text-anchor="middle" fill="white" font-size="20">🤖</text>
            </svg>
            <svg v-else width="40" height="40" viewBox="0 0 40 40">
              <circle cx="20" cy="20" r="20" fill="#52c41a"/>
              <text x="20" y="28" text-anchor="middle" fill="white" font-size="20">👤</text>
            </svg>
          </div>
          <div class="agent-details">
            <div class="agent-name">{{ isTransferred ? '人工客服' : 'AI智能助手' }}</div>
            <div class="agent-status">
              <span class="status-dot" :class="{ online: connected, offline: !connected }"></span>
              {{ connected ? '在线' : '离线' }}
            </div>
          </div>
        </div>
        <div class="header-actions">
          <button class="icon-btn" @click="toggleHistory" title="历史记录">
            📜
          </button>
          <button class="icon-btn" @click="closeChat" title="关闭">
            ✕
          </button>
        </div>
      </div>

      <!-- 消息列表 -->
      <div class="messages-container" ref="messagesContainer">
        <div
          v-for="(msg, index) in messages"
          :key="index"
          class="message"
          :class="msg.role"
        >
          <!-- 用户消息 -->
          <div v-if="msg.role === 'user'" class="message-content user">
            <div class="message-text">{{ msg.content }}</div>
            <div class="message-time">{{ formatTime(msg.timestamp) }}</div>
          </div>

          <!-- AI/客服消息 -->
          <div v-else class="message-content assistant">
            <div class="message-bubble">
              <div class="message-text" v-html="formatMessage(msg.content)"></div>

              <!-- 置信度指示器 -->
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

              <!-- 知识来源 -->
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

              <!-- 建议操作 -->
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

              <!-- 转人工提示 -->
              <div v-if="msg.transferTo" class="transfer-notice">
                <span class="transfer-icon">🔄</span>
                正在为您转接人工客服...
              </div>

              <div class="message-time">{{ formatTime(msg.timestamp) }}</div>
            </div>
          </div>

          <!-- 正在输入提示 -->
          <div v-if="isTyping" class="message assistant typing">
            <div class="typing-indicator">
              <span></span>
              <span></span>
              <span></span>
            </div>
          </div>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="input-container">
        <textarea
          v-model="inputMessage"
          placeholder="输入您的问题..."
          rows="1"
          ref="inputField"
          @keydown.enter.prevent="handleEnter"
          @input="autoResize"
        ></textarea>
        <button
          class="send-btn"
          @click="sendMessage"
          :disabled="!inputMessage.trim() || isTyping"
        >
          发送
        </button>
      </div>
    </div>

    <!-- 历史记录面板 -->
    <div v-if="showHistory" class="history-panel">
      <div class="history-header">
        <h3>对话历史</h3>
        <button class="icon-btn" @click="toggleHistory">✕</button>
      </div>
      <div class="history-list">
        <div
          v-for="(session, index) in historySessions"
          :key="index"
          class="history-item"
          @click="loadHistory(session.id)"
        >
          <div class="history-time">{{ formatDate(session.createdAt) }}</div>
          <div class="history-preview">{{ session.preview }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'

// 类型定义
interface Message {
  role: 'user' | 'assistant'
  content: string
  timestamp: number
  confidence?: number
  sources?: Source[]
  actions?: string[]
  transferTo?: string
}

interface Source {
  document_id: string
  title: string
  relevance: number
  snippet: string
}

interface HistorySession {
  id: string
  createdAt: number
  preview: string
}

// 响应式数据
const messages = ref<Message[]>([])
const inputMessage = ref('')
const isTyping = ref(false)
const connected = ref(false)
const isTransferred = ref(false)
const showHistory = ref(false)
const historySessions = ref<HistorySession[]>([])
const sessionId = ref(generateSessionId())

// WebSocket连接
let ws: WebSocket | null = null

// DOM引用
const messagesContainer = ref<HTMLElement | null>(null)
const inputField = ref<HTMLTextAreaElement | null>(null)

// 生成会话ID
function generateSessionId(): string {
  return `sess_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
}

// 连接WebSocket
function connectWebSocket() {
  const wsUrl = `ws://localhost:8888/api/css/ws?session_id=${sessionId.value}&channel=web`

  ws = new WebSocket(wsUrl)

  ws.onopen = () => {
    console.log('[CSS] WebSocket connected')
    connected.value = true

    // 欢迎消息
    addMessage({
      role: 'assistant',
      content: '您好！我是AI智能助手,有什么可以帮助您的吗?',
      timestamp: Date.now(),
      confidence: 1.0
    })
  }

  ws.onmessage = (event) => {
    const data = JSON.parse(event.data)
    console.log('[CSS] Received:', data)

    switch (data.type) {
      case 'connected':
        connected.value = true
        break

      case 'message':
        isTyping.value = false
        addMessage({
          role: data.role,
          content: data.content,
          timestamp: Date.now(),
          confidence: data.confidence,
          sources: data.sources,
          actions: data.actions,
          transferTo: data.transfer_to
        })

        // 检查是否转接
        if (data.transfer_to) {
          isTransferred.value = true
        }

        break

      case 'error':
        isTyping.value = false
        console.error('[CSS] Error:', data.message)
        break

      case 'ping':
        // 响应心跳
        ws?.send(JSON.stringify({ type: 'pong' }))
        break
    }

    scrollToBottom()
  }

  ws.onerror = (error) => {
    console.error('[CSS] WebSocket error:', error)
    connected.value = false
  }

  ws.onclose = () => {
    console.log('[CSS] WebSocket closed')
    connected.value = false

    // 5秒后重连
    setTimeout(() => {
      if (sessionId.value) {
        connectWebSocket()
      }
    }, 5000)
  }
}

// 发送消息
function sendMessage() {
  if (!inputMessage.value.trim()) return

  const question = inputMessage.value.trim()

  // 添加用户消息
  addMessage({
    role: 'user',
    content: question,
    timestamp: Date.now()
  })

  // 清空输入
  inputMessage.value = ''

  // 显示正在输入
  isTyping.value = true

  // 通过WebSocket发送
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      type: 'message',
      content: question
    }))
  } else {
    // WebSocket未连接,使用HTTP API
    sendViaHTTP(question)
  }
}

// 通过HTTP API发送
async function sendViaHTTP(question: string) {
  try {
    const response = await fetch('http://localhost:8888/api/css/question', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        session_id: sessionId.value,
        question: question,
        channel: 'web'
      })
    })

    const result = await response.json()

    if (result.success) {
      isTyping.value = false
      addMessage({
        role: 'assistant',
        content: result.data.answer,
        timestamp: result.data.timestamp,
        confidence: result.data.confidence,
        sources: result.data.sources,
        actions: result.data.actions,
        transferTo: result.data.transfer_to
      })

      if (result.data.transfer_to) {
        isTransferred.value = true
      }
    } else {
      isTyping.value = false
      addMessage({
        role: 'assistant',
        content: '抱歉,处理您的问题时出错了。请稍后再试。',
        timestamp: Date.now()
      })
    }

    scrollToBottom()
  } catch (error) {
    console.error('[CSS] Send message error:', error)
    isTyping.value = false
    addMessage({
      role: 'assistant',
      content: '网络连接失败,请检查网络后重试。',
      timestamp: Date.now()
    })
  }
}

// 添加消息
function addMessage(msg: Message) {
  messages.value.push(msg)
  nextTick(() => scrollToBottom())
}

// 滚动到底部
function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

// 格式化消息(支持Markdown)
function formatMessage(content: string): string {
  return content
    .replace(/\n/g, '<br>')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
}

// 格式化时间
function formatTime(timestamp: number): string {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

// 格式化日期
function formatDate(timestamp: number): string {
  const date = new Date(timestamp)
  return date.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取置信度样式
function getConfidenceClass(confidence: number): string {
  if (confidence >= 0.8) return 'high'
  if (confidence >= 0.6) return 'medium'
  return 'low'
}

// 处理建议操作
function handleAction(action: string) {
  inputMessage.value = action
  sendMessage()
}

// 处理回车
function handleEnter(event: KeyboardEvent) {
  if (!event.shiftKey) {
    sendMessage()
  }
}

// 自动调整文本框高度
function autoResize() {
  if (inputField.value) {
    inputField.value.style.height = 'auto'
    inputField.value.style.height = Math.min(inputField.value.scrollHeight, 120) + 'px'
  }
}

// 切换历史记录
function toggleHistory() {
  showHistory.value = !showHistory.value
  if (showHistory.value) {
    loadHistoryList()
  }
}

// 加载历史记录列表
async function loadHistoryList() {
  // 这里应该调用API获取历史会话列表
  // 暂时使用模拟数据
  historySessions.value = [
    {
      id: 'sess_1',
      createdAt: Date.now() - 3600000,
      preview: '如何使用产品A的功能...'
    },
    {
      id: 'sess_2',
      createdAt: Date.now() - 86400000,
      preview: '退款流程是什么...'
    }
  ]
}

// 加载历史会话
async function loadHistory(sessionId: string) {
  try {
    const response = await fetch(`http://localhost:8888/api/css/history/${sessionId}`)
    const result = await response.json()

    if (result.success) {
      messages.value = result.data.map((msg: any) => ({
        role: msg.role,
        content: msg.content,
        timestamp: msg.created_at
      }))
      showHistory.value = false
      scrollToBottom()
    }
  } catch (error) {
    console.error('[CSS] Load history error:', error)
  }
}

// 关闭聊天
function closeChat() {
  if (ws) {
    ws.close()
  }

  // 关闭会话
  fetch(`http://localhost:8888/api/css/session/${sessionId.value}/close`, {
    method: 'POST'
  })

  emit('close')
}

// 组件挂载
onMounted(() => {
  connectWebSocket()
  inputField.value?.focus()
})

// 组件卸载
onUnmounted(() => {
  if (ws) {
    ws.close()
  }
})

// 暴露方法
defineEmits(['close'])
</script>

<style scoped>
.customer-service-chat {
  position: fixed;
  bottom: 20px;
  right: 20px;
  width: 400px;
  height: 600px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  z-index: 1000;
}

.chat-header {
  padding: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px 12px 0 0;
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.agent-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.agent-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
}

.agent-details {
  display: flex;
  flex-direction: column;
}

.agent-name {
  font-weight: 600;
  font-size: 14px;
}

.agent-status {
  font-size: 12px;
  opacity: 0.9;
  display: flex;
  align-items: center;
  gap: 4px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-dot.online {
  background: #52c41a;
}

.status-dot.offline {
  background: #ff4d4f;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.icon-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.2s;
}

.icon-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message {
  display: flex;
}

.message.user {
  justify-content: flex-end;
}

.message.assistant {
  justify-content: flex-start;
}

.message-content {
  max-width: 80%;
  display: flex;
  flex-direction: column;
}

.message-content.user .message-text {
  background: #667eea;
  color: white;
  padding: 10px 14px;
  border-radius: 12px 12px 0 12px;
}

.message-content.assistant .message-bubble {
  background: #f5f5f5;
  padding: 12px;
  border-radius: 12px 12px 12px 0;
}

.message-text {
  font-size: 14px;
  line-height: 1.5;
  word-wrap: break-word;
}

.message-time {
  font-size: 11px;
  color: #999;
  margin-top: 4px;
  text-align: right;
}

.message-content.assistant .message-time {
  text-align: left;
}

/* 置信度指示器 */
.confidence-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  padding: 8px;
  background: white;
  border-radius: 6px;
  font-size: 12px;
}

.confidence-label {
  color: #666;
}

.confidence-bar {
  flex: 1;
  height: 6px;
  background: #e5e5e5;
  border-radius: 3px;
  overflow: hidden;
}

.confidence-fill {
  height: 100%;
  transition: width 0.3s ease;
}

.confidence-fill.high {
  background: #52c41a;
}

.confidence-fill.medium {
  background: #faad14;
}

.confidence-fill.low {
  background: #ff4d4f;
}

.confidence-value {
  font-weight: 600;
  min-width: 40px;
}

/* 知识来源 */
.sources {
  margin-top: 12px;
  padding: 12px;
  background: white;
  border-radius: 6px;
  border-left: 3px solid #1890ff;
}

.sources-title {
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #333;
}

.source-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 0;
  font-size: 12px;
}

.source-title {
  color: #666;
}

.source-relevance {
  color: #1890ff;
  font-weight: 500;
}

/* 建议操作 */
.suggested-actions {
  margin-top: 12px;
}

.actions-title {
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #333;
}

.action-btn {
  display: inline-block;
  padding: 6px 12px;
  margin: 4px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  transition: background 0.2s;
}

.action-btn:hover {
  background: #40a9ff;
}

/* 转接提示 */
.transfer-notice {
  margin-top: 12px;
  padding: 10px;
  background: #fff7e6;
  border: 1px solid #ffd591;
  border-radius: 6px;
  color: #d46b08;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.transfer-icon {
  font-size: 16px;
}

/* 正在输入 */
.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 12px;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  background: #999;
  border-radius: 50%;
  animation: typing 1.4s infinite;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% {
    transform: translateY(0);
  }
  30% {
    transform: translateY(-10px);
  }
}

/* 输入区域 */
.input-container {
  display: flex;
  gap: 8px;
  padding: 16px;
  border-top: 1px solid #e5e5e5;
}

.input-container textarea {
  flex: 1;
  resize: none;
  padding: 10px;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  font-size: 14px;
  font-family: inherit;
  outline: none;
  transition: border 0.2s;
}

.input-container textarea:focus {
  border-color: #1890ff;
}

.send-btn {
  padding: 10px 20px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: background 0.2s;
}

.send-btn:hover:not(:disabled) {
  background: #40a9ff;
}

.send-btn:disabled {
  background: #d9d9d9;
  cursor: not-allowed;
}

/* 历史记录面板 */
.history-panel {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: white;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
}

.history-header {
  padding: 16px;
  background: #f5f5f5;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-radius: 12px 12px 0 0;
}

.history-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.history-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.history-item {
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
  margin-bottom: 4px;
}

.history-item:hover {
  background: #f5f5f5;
}

.history-time {
  font-size: 12px;
  color: #999;
  margin-bottom: 4px;
}

.history-preview {
  font-size: 14px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 响应式 */
@media (max-width: 480px) {
  .customer-service-chat {
    width: 100%;
    height: 100%;
    bottom: 0;
    right: 0;
    border-radius: 0;
  }

  .chat-header {
    border-radius: 0;
  }
}
</style>
