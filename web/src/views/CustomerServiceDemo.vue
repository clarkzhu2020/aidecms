<template>
  <div class="customer-service-demo">
    <div class="demo-header">
      <h1>智能客服系统演示</h1>
      <p>AI驱动的智能问答系统,支持RAG检索、置信度评估和人工转接</p>
    </div>

    <div class="demo-layout">
      <!-- 聊天窗口 -->
      <div class="demo-chat">
        <CustomerServiceChat @close="showChat = false" />
      </div>

      <!-- 控制面板 -->
      <div class="demo-controls">
        <h2>控制面板</h2>

        <!-- 连接状态 -->
        <div class="control-section">
          <h3>连接状态</h3>
          <div class="status-indicators">
            <div class="status-item">
              <span class="status-label">WebSocket:</span>
              <span class="status-value" :class="{ online: css.isConnected, offline: !css.isConnected }">
                {{ css.isConnected ? '已连接' : '未连接' }}
              </span>
            </div>
            <div class="status-item">
              <span class="status-label">会话ID:</span>
              <span class="status-value">{{ css.sessionId }}</span>
            </div>
            <div class="status-item">
              <span class="status-label">转接状态:</span>
              <span class="status-value" :class="{ transferred: css.isTransferred }">
                {{ css.isTransferred ? '已转人工' : 'AI服务中' }}
              </span>
            </div>
          </div>
        </div>

        <!-- 统计信息 -->
        <div class="control-section">
          <h3>统计信息</h3>
          <div class="stats-grid">
            <div class="stat-card">
              <div class="stat-value">{{ css.messageCount }}</div>
              <div class="stat-label">消息数</div>
            </div>
            <div class="stat-card">
              <div class="stat-value">{{ (css.avgConfidence * 100).toFixed(0) }}%</div>
              <div class="stat-label">平均置信度</div>
            </div>
          </div>
        </div>

        <!-- 快捷问题 -->
        <div class="control-section">
          <h3>快捷问题</h3>
          <div class="quick-questions">
            <button
              v-for="(question, index) in quickQuestions"
              :key="index"
              class="quick-btn"
              @click="sendQuickQuestion(question)"
            >
              {{ question }}
            </button>
          </div>
        </div>

        <!-- API调用 -->
        <div class="control-section">
          <h3>API调用</h3>
          <div class="api-buttons">
            <button class="api-btn" @click="loadHistory">加载历史</button>
            <button class="api-btn" @click="checkStatus">系统状态</button>
            <button class="api-btn" @click="closeSession">关闭会话</button>
            <button class="api-btn danger" @click="resetSession">重置会话</button>
          </div>
        </div>

        <!-- 系统日志 -->
        <div class="control-section">
          <h3>系统日志</h3>
          <div class="log-container">
            <div
              v-for="(log, index) in logs"
              :key="index"
              class="log-entry"
              :class="log.type"
            >
              <span class="log-time">{{ formatLogTime(log.timestamp) }}</span>
              <span class="log-message">{{ log.message }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import CustomerServiceChat from '../components/CustomerServiceChat.vue'
import { useCustomerService } from '../composables/useCustomerService'

// 使用客服系统
const css = useCustomerService({
  onMessage: (message) => {
    addLog('info', `收到消息: ${message.content.substring(0, 50)}...`)
  },
  onError: (error) => {
    addLog('error', `错误: ${error.message}`)
  },
  onConnect: () => {
    addLog('success', 'WebSocket已连接')
  },
  onDisconnect: () => {
    addLog('warning', 'WebSocket已断开')
  }
})

// 快捷问题
const quickQuestions = [
  '如何使用产品A?',
  '退款流程是什么?',
  '如何联系客服?',
  '我要投诉这个产品',
  '产品B的价格是多少?'
]

// 日志
interface LogEntry {
  timestamp: number
  type: 'info' | 'success' | 'warning' | 'error'
  message: string
}

const logs = ref<LogEntry[]>([])

// 添加日志
function addLog(type: LogEntry['type'], message: string): void {
  logs.value.push({
    timestamp: Date.now(),
    type,
    message
  })

  // 保持最多100条日志
  if (logs.value.length > 100) {
    logs.value.shift()
  }
}

// 格式化日志时间
function formatLogTime(timestamp: number): string {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 发送快捷问题
async function sendQuickQuestion(question: string): Promise<void> {
  addLog('info', `发送问题: ${question}`)
  const message = await css.sendMessage(question)
  if (message) {
    addLog('success', '消息已发送')
  }
}

// 加载历史
async function loadHistory(): Promise<void> {
  addLog('info', '加载对话历史...')
  const history = await css.getHistory()
  addLog('success', `加载了 ${history.length} 条历史消息`)
}

// 检查系统状态
async function checkStatus(): Promise<void> {
  addLog('info', '获取系统状态...')
  const status = await css.getStatus()
  if (status) {
    addLog('success', `AI模型: ${status.engine.ai_model}, 会话数: ${status.engine.session_count}`)
  } else {
    addLog('error', '获取系统状态失败')
  }
}

// 关闭会话
async function closeSession(): Promise<void> {
  addLog('info', '关闭当前会话...')
  const success = await css.closeSession()
  if (success) {
    addLog('success', '会话已关闭')
  } else {
    addLog('error', '关闭会话失败')
  }
}

// 重置会话
function resetSession(): void {
  addLog('warning', '重置会话状态...')
  css.reset()
  addLog('success', '会话已重置')
}

// 组件挂载
onMounted(() => {
  addLog('info', '客服系统演示已启动')
})
</script>

<style scoped>
.customer-service-demo {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 40px 20px;
}

.demo-header {
  text-align: center;
  margin-bottom: 40px;
}

.demo-header h1 {
  font-size: 36px;
  font-weight: 700;
  color: #2c3e50;
  margin: 0 0 12px 0;
}

.demo-header p {
  font-size: 18px;
  color: #5a6c7d;
  margin: 0;
}

.demo-layout {
  max-width: 1400px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 24px;
  align-items: start;
}

.demo-chat {
  position: fixed;
  bottom: 20px;
  right: 20px;
  z-index: 1000;
}

.demo-controls {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.demo-controls h2 {
  padding: 20px;
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  background: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
}

.control-section {
  padding: 20px;
  border-bottom: 1px solid #e9ecef;
}

.control-section:last-child {
  border-bottom: none;
}

.control-section h3 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: #2c3e50;
}

/* 状态指示器 */
.status-indicators {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-label {
  font-size: 14px;
  color: #5a6c7d;
}

.status-value {
  font-size: 14px;
  font-weight: 500;
  color: #2c3e50;
}

.status-value.online {
  color: #52c41a;
}

.status-value.offline {
  color: #ff4d4f;
}

.status-value.transferred {
  color: #1890ff;
}

/* 统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.stat-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 16px;
  text-align: center;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1890ff;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: #5a6c7d;
}

/* 快捷问题 */
.quick-questions {
  display: grid;
  grid-template-columns: 1fr;
  gap: 8px;
}

.quick-btn {
  padding: 10px 14px;
  background: #1890ff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  text-align: left;
  transition: background 0.2s;
}

.quick-btn:hover {
  background: #40a9ff;
}

/* API按钮 */
.api-buttons {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.api-btn {
  padding: 10px;
  background: #52c41a;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: background 0.2s;
}

.api-btn:hover {
  background: #73d13d;
}

.api-btn.danger {
  background: #ff4d4f;
}

.api-btn.danger:hover {
  background: #ff7875;
}

/* 日志容器 */
.log-container {
  max-height: 300px;
  overflow-y: auto;
  background: #1e1e1e;
  border-radius: 6px;
  padding: 12px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

.log-entry {
  display: flex;
  gap: 12px;
  margin-bottom: 4px;
  line-height: 1.6;
}

.log-time {
  color: #858585;
  flex-shrink: 0;
}

.log-message {
  color: #d4d4d4;
}

.log-entry.info .log-message {
  color: #4fc1ff;
}

.log-entry.success .log-message {
  color: #4ade80;
}

.log-entry.warning .log-message {
  color: #fbbf24;
}

.log-entry.error .log-message {
  color: #f87171;
}

/* 响应式 */
@media (max-width: 1024px) {
  .demo-layout {
    grid-template-columns: 1fr;
  }

  .demo-chat {
    position: static;
  }
}
</style>
