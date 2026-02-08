/**
 * 客服系统 Composable
 *
 * 提供客服系统的完整功能封装
 */

import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'

// 类型定义
export interface CSSMessage {
  role: 'user' | 'assistant'
  content: string
  timestamp: number
  confidence?: number
  sources?: CSSSource[]
  actions?: string[]
  transferTo?: string
}

export interface CSSSource {
  document_id: string
  title: string
  relevance: number
  snippet: string
}

export interface CSSConfig {
  wsUrl?: string
  httpUrl?: string
  sessionId?: string
  channel?: string
  autoReconnect?: boolean
  reconnectDelay?: number
  onMessage?: (message: CSSMessage) => void
  onError?: (error: Error) => void
  onConnect?: () => void
  onDisconnect?: () => void
}

export interface CSSEngineStats {
  config: {
    default_model: string
    temperature: number
    max_tokens: number
    top_k: number
    confidence_threshold: number
  }
  ai_model: string
  rag_enabled: boolean
  session_count: number
}

/**
 * 客服系统 Hook
 *
 * 使用示例:
 * ```ts
 * const {
 *   messages,
 *   isConnected,
 *   isTyping,
 *   isTransferred,
 *   sendMessage,
 *   getHistory,
 *   getStatus
 * } = useCustomerService()
 * ```
 */
export function useCustomerService(config: CSSConfig = {}) {
  // 默认配置
  const defaultConfig: CSSConfig = {
    wsUrl: 'ws://localhost:8888/api/css/ws',
    httpUrl: 'http://localhost:8888/api/css',
    sessionId: generateSessionId(),
    channel: 'web',
    autoReconnect: true,
    reconnectDelay: 5000,
    ...config
  }

  // 响应式状态
  const messages = ref<CSSMessage[]>([])
  const isConnected = ref(false)
  const isTyping = ref(false)
  const isTransferred = ref(false)
  const sessionId = ref(defaultConfig.sessionId!)
  const error = ref<Error | null>(null)

  // WebSocket实例
  let ws: WebSocket | null = null
  let reconnectTimer: number | null = null

  /**
   * 生成会话ID
   */
  function generateSessionId(): string {
    return `sess_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
  }

  /**
   * 连接WebSocket
   */
  function connect(): void {
    const wsUrl = `${defaultConfig.wsUrl}?session_id=${sessionId.value}&channel=${defaultConfig.channel}`

    try {
      ws = new WebSocket(wsUrl)

      ws.onopen = () => {
        console.log('[CSS] WebSocket connected')
        isConnected.value = true
        error.value = null
        defaultConfig.onConnect?.()
      }

      ws.onmessage = (event) => {
        const data = JSON.parse(event.data)
        handleWebSocketMessage(data)
      }

      ws.onerror = (event) => {
        console.error('[CSS] WebSocket error:', event)
        const err = new Error('WebSocket connection error')
        error.value = err
        defaultConfig.onError?.(err)
      }

      ws.onclose = () => {
        console.log('[CSS] WebSocket closed')
        isConnected.value = false
        defaultConfig.onDisconnect?.()

        // 自动重连
        if (defaultConfig.autoReconnect) {
          reconnectTimer = window.setTimeout(() => {
            connect()
          }, defaultConfig.reconnectDelay)
        }
      }
    } catch (err) {
      const errorObj = err instanceof Error ? err : new Error('Failed to connect')
      error.value = errorObj
      defaultConfig.onError?.(errorObj)
    }
  }

  /**
   * 处理WebSocket消息
   */
  function handleWebSocketMessage(data: any): void {
    console.log('[CSS] Received:', data)

    switch (data.type) {
      case 'connected':
        isConnected.value = true
        break

      case 'message':
        isTyping.value = false
        const message: CSSMessage = {
          role: data.role,
          content: data.content,
          timestamp: Date.now(),
          confidence: data.confidence,
          sources: data.sources,
          actions: data.actions,
          transferTo: data.transfer_to
        }

        messages.value.push(message)

        // 检查是否转接
        if (data.transfer_to) {
          isTransferred.value = true
        }

        defaultConfig.onMessage?.(message)
        break

      case 'typing':
        isTyping.value = true
        break

      case 'error':
        isTyping.value = false
        const err = new Error(data.message || 'Unknown error')
        error.value = err
        defaultConfig.onError?.(err)
        break

      case 'ping':
        // 响应心跳
        if (ws?.readyState === WebSocket.OPEN) {
          ws.send(JSON.stringify({ type: 'pong' }))
        }
        break
    }
  }

  /**
   * 发送消息
   */
  async function sendMessage(question: string): Promise<CSSMessage | null> {
    if (!question.trim()) {
      console.warn('[CSS] Empty message')
      return null
    }

    // 添加用户消息
    const userMessage: CSSMessage = {
      role: 'user',
      content: question.trim(),
      timestamp: Date.now()
    }
    messages.value.push(userMessage)

    // 显示正在输入
    isTyping.value = true

    // 通过WebSocket发送
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({
        type: 'message',
        content: question.trim()
      }))
    } else {
      // WebSocket未连接,使用HTTP API
      return await sendViaHTTP(question)
    }

    return userMessage
  }

  /**
   * 通过HTTP API发送消息
   */
  async function sendViaHTTP(question: string): Promise<CSSMessage | null> {
    try {
      const response = await fetch(`${defaultConfig.httpUrl}/question`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          session_id: sessionId.value,
          question: question.trim(),
          channel: defaultConfig.channel
        })
      })

      const result = await response.json()

      if (result.success) {
        isTyping.value = false

        const assistantMessage: CSSMessage = {
          role: 'assistant',
          content: result.data.answer,
          timestamp: result.data.timestamp,
          confidence: result.data.confidence,
          sources: result.data.sources,
          actions: result.data.actions,
          transferTo: result.data.transfer_to
        }

        messages.value.push(assistantMessage)

        if (result.data.transfer_to) {
          isTransferred.value = true
        }

        defaultConfig.onMessage?.(assistantMessage)
        return assistantMessage
      } else {
        throw new Error(result.error || 'Failed to send message')
      }
    } catch (err) {
      isTyping.value = false
      const errorObj = err instanceof Error ? err : new Error('Network error')
      error.value = errorObj
      defaultConfig.onError?.(errorObj)
      return null
    }
  }

  /**
   * 获取对话历史
   */
  async function getHistory(sessionIdParam?: string): Promise<CSSMessage[]> {
    const sid = sessionIdParam || sessionId.value

    try {
      const response = await fetch(`${defaultConfig.httpUrl}/history/${sid}`)
      const result = await response.json()

      if (result.success) {
        return result.data.map((msg: any) => ({
          role: msg.role,
          content: msg.content,
          timestamp: new Date(msg.created_at).getTime(),
          confidence: msg.confidence
        }))
      } else {
        throw new Error(result.error || 'Failed to get history')
      }
    } catch (err) {
      const errorObj = err instanceof Error ? err : new Error('Network error')
      error.value = errorObj
      defaultConfig.onError?.(errorObj)
      return []
    }
  }

  /**
   * 关闭会话
   */
  async function closeSession(): Promise<boolean> {
    try {
      const response = await fetch(`${defaultConfig.httpUrl}/session/${sessionId.value}/close`, {
        method: 'POST'
      })
      const result = await response.json()
      return result.success
    } catch (err) {
      const errorObj = err instanceof Error ? err : new Error('Network error')
      error.value = errorObj
      return false
    }
  }

  /**
   * 获取系统状态
   */
  async function getStatus(): Promise<{ engine: CSSEngineStats; clients: number; timestamp: number } | null> {
    try {
      const response = await fetch(`${defaultConfig.httpUrl}/status`)
      const result = await response.json()

      if (result.success) {
        return result.data
      } else {
        throw new Error(result.error || 'Failed to get status')
      }
    } catch (err) {
      const errorObj = err instanceof Error ? err : new Error('Network error')
      error.value = errorObj
      defaultConfig.onError?.(errorObj)
      return null
    }
  }

  /**
   * 断开连接
   */
  function disconnect(): void {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }

    if (ws) {
      ws.close()
      ws = null
    }

    isConnected.value = false
  }

  /**
   * 清空消息
   */
  function clearMessages(): void {
    messages.value = []
  }

  /**
   * 重置状态
   */
  function reset(): void {
    disconnect()
    clearMessages()
    isTyping.value = false
    isTransferred.value = false
    error.value = null
    sessionId.value = generateSessionId()
  }

  // 计算属性
  const lastMessage = computed(() => {
    return messages.value[messages.value.length - 1] || null
  })

  const messageCount = computed(() => messages.value.length)

  const avgConfidence = computed(() => {
    const assistantMessages = messages.value.filter(m => m.role === 'assistant' && m.confidence !== undefined)
    if (assistantMessages.length === 0) return 0
    const sum = assistantMessages.reduce((acc, m) => acc + (m.confidence || 0), 0)
    return sum / assistantMessages.length
  })

  // 组件挂载时连接
  onMounted(() => {
    connect()
  })

  // 组件卸载时断开连接
  onUnmounted(() => {
    disconnect()
  })

  return {
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
  }
}

/**
 * 客服服务类 (面向对象API)
 */
export class CustomerService {
  private config: CSSConfig
  private ws: WebSocket | null = null
  private callbacks: Map<string, Function[]> = new Map()

  constructor(config: CSSConfig = {}) {
    this.config = {
      wsUrl: 'ws://localhost:8888/api/css/ws',
      httpUrl: 'http://localhost:8888/api/css',
      sessionId: `sess_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      channel: 'web',
      autoReconnect: true,
      reconnectDelay: 5000,
      ...config
    }
  }

  /**
   * 连接到服务器
   */
  connect(): void {
    const wsUrl = `${this.config.wsUrl}?session_id=${this.config.sessionId}&channel=${this.config.channel}`
    this.ws = new WebSocket(wsUrl)

    this.ws.onopen = () => {
      this.emit('connected')
    }

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      this.emit('message', data)
    }

    this.ws.onerror = (error) => {
      this.emit('error', error)
    }

    this.ws.onclose = () => {
      this.emit('disconnected')
    }
  }

  /**
   * 发送消息
   */
  send(question: string): Promise<CSSMessage> {
    return new Promise((resolve, reject) => {
      fetch(`${this.config.httpUrl}/question`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          session_id: this.config.sessionId,
          question,
          channel: this.config.channel
        })
      })
        .then(res => res.json())
        .then(result => {
          if (result.success) {
            resolve({
              role: 'assistant',
              content: result.data.answer,
              timestamp: result.data.timestamp,
              confidence: result.data.confidence,
              sources: result.data.sources,
              actions: result.data.actions,
              transferTo: result.data.transfer_to
            })
          } else {
            reject(new Error(result.error))
          }
        })
        .catch(reject)
    })
  }

  /**
   * 获取历史
   */
  async getHistory(): Promise<CSSMessage[]> {
    const response = await fetch(`${this.config.httpUrl}/history/${this.config.sessionId}`)
    const result = await response.json()
    return result.success ? result.data : []
  }

  /**
   * 事件监听
   */
  on(event: string, callback: Function): void {
    if (!this.callbacks.has(event)) {
      this.callbacks.set(event, [])
    }
    this.callbacks.get(event)!.push(callback)
  }

  /**
   * 移除监听
   */
  off(event: string, callback?: Function): void {
    if (!callback) {
      this.callbacks.delete(event)
    } else {
      const callbacks = this.callbacks.get(event)
      if (callbacks) {
        const index = callbacks.indexOf(callback)
        if (index > -1) {
          callbacks.splice(index, 1)
        }
      }
    }
  }

  /**
   * 触发事件
   */
  private emit(event: string, data?: any): void {
    const callbacks = this.callbacks.get(event)
    if (callbacks) {
      callbacks.forEach(callback => callback(data))
    }
  }

  /**
   * 关闭连接
   */
  close(): void {
    this.ws?.close()
    this.callbacks.clear()
  }
}
