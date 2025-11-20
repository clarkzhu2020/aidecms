<template>
  <div class="features-page">
    <section class="page-hero">
      <div class="hero-bg">
        <div class="glow glow-1"></div>
        <div class="glow glow-2"></div>
        <div class="grid-pattern"></div>
      </div>
      <div class="container">
        <h1 data-aos="fade-up">功能特性</h1>
        <p data-aos="fade-up" data-aos-delay="100">
          全面了解 AideCMS 提供的强大功能
        </p>
      </div>
    </section>

    <section class="features-list">
      <div class="container">
        <div v-for="(feature, index) in features" :key="index" class="feature-item" :class="{ 'reverse': index % 2 === 1 }" data-aos="fade-up">
          <div class="feature-content">
            <span class="badge">{{ feature.badge }}</span>
            <h3>{{ feature.title }}</h3>
            <p>{{ feature.description }}</p>
            <ul class="feature-points">
              <li v-for="(item, i) in feature.items" :key="i">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M5 13l4 4L19 7" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                {{ item }}
              </li>
            </ul>
          </div>
          <div class="feature-code">
            <div class="code-window">
              <div class="window-header">
                <div class="window-dots">
                  <span></span><span></span><span></span>
                </div>
                <div class="window-title">example.go</div>
              </div>
              <pre><code>{{ feature.code }}</code></pre>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="comparison">
      <div class="container">
        <div class="section-header" data-aos="fade-up">
          <h2>为什么选择 AideCMS？</h2>
          <p>与传统框架的对比</p>
        </div>

        <div class="comparison-table" data-aos="fade-up" data-aos-delay="200">
          <div class="comparison-row header">
            <div class="comparison-cell">特性</div>
            <div class="comparison-cell highlight">AideCMS</div>
            <div class="comparison-cell">传统框架</div>
          </div>
          <div class="comparison-row" v-for="item in comparisonData" :key="item.feature">
            <div class="comparison-cell">{{ item.feature }}</div>
            <div class="comparison-cell highlight">
              <svg v-if="item.aidecms" width="24" height="24" viewBox="0 0 24 24" fill="currentColor" style="color: #10b981">
                <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <span v-else style="color: #ef4444">×</span>
            </div>
            <div class="comparison-cell">
              <svg v-if="item.traditional" width="24" height="24" viewBox="0 0 24 24" fill="currentColor" style="color: #10b981">
                <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <span v-else style="color: #ef4444">×</span>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
const features = [
  {
    badge: '高性能',
    title: 'CloudWeGo Hertz 框架',
    description: '基于字节跳动开源的 Hertz 框架，提供极致的性能和稳定性',
    items: [
      '高性能网络库，支持百万级并发',
      '自动服务发现和负载均衡',
      '完善的中间件生态',
      '企业级生产环境验证'
    ],
    code: `app := framework.NewApplication()
app.Boot()
app.Run()`
  },
  {
    badge: 'AI 驱动',
    title: 'AI 智能集成',
    description: '内置多种 AI 模型支持，轻松实现智能内容生成和分析',
    items: [
      '支持 OpenAI、Azure、字节跳动豆包等多家AI服务',
      '智能内容生成和优化',
      '语义分析和关键词提取',
      '流式响应支持'
    ],
    code: `manager := ai.NewManager()
response, _ := manager.Chat(ctx, "帮我写一篇关于Go语言的文章")

stream, _ := manager.ChatStream(ctx, prompt)
for chunk := range stream {
  fmt.Print(chunk)
}`
  },
  {
    badge: '任务调度',
    title: '强大的任务调度系统',
    description: '完整的任务调度解决方案，支持 Cron 表达式和便捷方法',
    items: [
      '标准 Cron 表达式支持',
      '15+ 便捷调度方法',
      '并发任务执行',
      '任务统计和监控'
    ],
    code: `scheduler := schedule.NewScheduler()

scheduler.EveryMinute(func() {
  log.Println("Task executed!")
})

scheduler.Daily().At("00:00", cleanupTask)
scheduler.WeeklyOn(1, "09:00", reportTask)`
  },
  {
    badge: '队列系统',
    title: '高性能消息队列',
    description: '支持 Redis 和内存驱动的消息队列，确保任务可靠执行',
    items: [
      'Redis 和内存双驱动',
      '延迟任务支持',
      '自动失败重试',
      '死信队列处理',
      'Worker Pool 并发控制'
    ],
    code: `queue := queue.NewQueue("redis")

queue.Push(&queue.Job{
  Name: "send-email",
  Data: emailData,
  Delay: 5 * time.Minute,
})

queue.Process("send-email", func(job *Job) {
  sendEmail(job.Data)
})`
  },
  {
    badge: 'Web3 集成',
    title: '区块链与 Web3 支持',
    description: '内置多链支持，轻松实现链上数据交互',
    items: [
      '多链支持 (ETH, BTC, etc)',
      '链上数据同步与查询',
      '智能合约交互',
      '钱包集成'
    ],
    code: `manager := web3.GetManager()
balance, _ := manager.GetBalance(ctx, web3.Ethereum, address)
transaction, _ := manager.GetTransaction(ctx, web3.Bitcoin, txHash)`
  },
  {
    badge: '交易所集成',
    title: '主流交易所 API 对接',
    description: '内置 Coinbase、KuCoin、Hyperliquid 等交易所集成，支持行情、资产、自动化内容驱动交易',
    items: [
      '多交易所支持',
      '统一 API 接口',
      '实时行情获取',
      '自动化交易支持'
    ],
    code: `client := exchange.NewClient("coinbase", apiKey, apiSecret)
balance, _ := client.GetBalance(ctx)
price, _ := client.GetTicker(ctx, "BTC-USD")`
  },
  {
    badge: '事件系统',
    title: '事件驱动架构',
    description: '灵活的事件系统，支持同步/异步监听器和优先级',
    items: [
      '同步和异步监听器',
      '优先级支持',
      '事件统计分析'
    ],
    code: `// 注册监听器
event.Listen("user.created", func(e Event) {
  sendWelcomeEmail(e.Data)
}, event.WithPriority(10))

event.Dispatch("user.created", userData)
event.DispatchAsync("user.created", userData)`
  }
]

const comparisonData = [
  { feature: '高性能架构', aidecms: true, traditional: false },
  { feature: 'AI 原生集成', aidecms: true, traditional: false },
  { feature: 'Web3 & 交易所支持', aidecms: true, traditional: true },
  { feature: '模块化设计', aidecms: true, traditional: true },
  { feature: '开箱即用', aidecms: true, traditional: false },
]
</script>

<style scoped>
.page-hero {
  position: relative;
  padding: 140px 0 80px;
  text-align: center;
  overflow: hidden;
}

.hero-bg {
  position: absolute;
  inset: 0;
  z-index: -1;
  background: var(--bg-secondary);
}

.grid-pattern {
  position: absolute;
  inset: 0;
  background-image: 
    linear-gradient(rgba(99, 102, 241, 0.05) 1px, transparent 1px),
    linear-gradient(90deg, rgba(99, 102, 241, 0.05) 1px, transparent 1px);
  background-size: 40px 40px;
  mask-image: radial-gradient(circle at center, black 40%, transparent 80%);
}

.glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.3;
}

.glow-1 {
  top: -20%;
  left: 20%;
  width: 500px;
  height: 500px;
  background: var(--primary-color);
}

.glow-2 {
  bottom: -20%;
  right: 20%;
  width: 400px;
  height: 400px;
  background: var(--secondary-color);
}

.page-hero h1 {
  font-size: 3rem;
  font-weight: 800;
  margin-bottom: 1rem;
  background: linear-gradient(to bottom, var(--text-primary), var(--text-secondary));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-hero p {
  font-size: 1.25rem;
  color: var(--text-secondary);
}

.features-list {
  padding: 80px 0;
}

.feature-item {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 4rem;
  align-items: center;
  margin-bottom: 8rem;
}

.feature-item:last-child {
  margin-bottom: 0;
}

.feature-item.reverse {
  direction: rtl;
}

.feature-item.reverse .feature-content {
  direction: ltr;
}

.feature-content h3 {
  font-size: 2rem;
  margin-bottom: 1rem;
}

.feature-content p {
  color: var(--text-secondary);
  font-size: 1.125rem;
  margin-bottom: 2rem;
  line-height: 1.7;
}

.badge {
  display: inline-block;
  padding: 4px 12px;
  background: rgba(99, 102, 241, 0.1);
  color: var(--primary-color);
  border-radius: 20px;
  font-size: 0.875rem;
  font-weight: 600;
  margin-bottom: 1rem;
}

.feature-points {
  list-style: none;
  padding: 0;
}

.feature-points li {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 1rem;
  color: var(--text-primary);
  font-weight: 500;
}

.feature-points svg {
  color: var(--primary-color);
}

/* Code Window */
.code-window {
  background: #1e293b;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: var(--shadow-xl);
  border: 1px solid rgba(255, 255, 255, 0.1);
  direction: ltr;
}

.window-header {
  background: rgba(15, 23, 42, 0.5);
  padding: 12px 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.window-dots {
  display: flex;
  gap: 6px;
}

.window-dots span {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.window-dots span:nth-child(1) { background: #ef4444; }
.window-dots span:nth-child(2) { background: #f59e0b; }
.window-dots span:nth-child(3) { background: #10b981; }

.window-title {
  color: #94a3b8;
  font-size: 0.875rem;
  font-family: var(--font-mono);
}

.code-window pre {
  margin: 0;
  padding: 24px;
  color: #e2e8f0;
  font-family: var(--font-mono);
  font-size: 0.875rem;
  line-height: 1.6;
  overflow-x: auto;
}

/* Comparison Section */
.comparison {
  padding: 100px 0;
  background: var(--bg-secondary);
}

.section-header {
  text-align: center;
  margin-bottom: 4rem;
}

.section-header h2 {
  margin-bottom: 1rem;
}

.section-header p {
  color: var(--text-secondary);
  font-size: 1.125rem;
}

.comparison-table {
  background: var(--bg-color);
  border-radius: 16px;
  box-shadow: var(--shadow-lg);
  overflow: hidden;
  max-width: 900px;
  margin: 0 auto;
}

.comparison-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr;
  border-bottom: 1px solid var(--border-color);
}

.comparison-row:last-child {
  border-bottom: none;
}

.comparison-row.header {
  background: var(--bg-secondary);
  font-weight: 600;
  color: var(--text-primary);
}

.comparison-cell {
  padding: 1.5rem;
  display: flex;
  align-items: center;
}

.comparison-cell.highlight {
  justify-content: center;
  background: rgba(99, 102, 241, 0.02);
}

.comparison-row:not(.header) .comparison-cell:not(:first-child) {
  justify-content: center;
}

@media (max-width: 968px) {
  .feature-item {
    grid-template-columns: 1fr;
    gap: 2rem;
    margin-bottom: 6rem;
  }

  .feature-item.reverse {
    direction: ltr;
  }

  .feature-item.reverse .feature-content {
    direction: ltr;
  }

  .feature-code {
    order: -1;
  }
}
</style>
