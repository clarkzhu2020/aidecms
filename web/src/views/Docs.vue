<template>
  <div class="docs-page">
    <section class="page-hero">
      <div class="container">
        <h1 data-aos="fade-up">开发文档</h1>
        <p data-aos="fade-up" data-aos-delay="100">
          快速开始使用 AideCMS
        </p>
      </div>
    </section>

    <section class="docs-content">
      <div class="container">
        <div class="docs-layout">
          <aside class="docs-sidebar" data-aos="fade-right">
            <nav class="docs-nav">
              <div class="nav-section" v-for="section in docsSections" :key="section.title">
                <h3>{{ section.title }}</h3>
                <ul>
                  <li v-for="item in section.items" :key="item.id">
                    <a href="#" :class="{ active: item.id === activeDocId }" @click.prevent="scrollToDoc(item.id)">
                      {{ item.title }}
                    </a>
                  </li>
                </ul>
              </div>
            </nav>
          </aside>

          <main class="docs-main" data-aos="fade-left">
            <div class="doc-content">
              
              <!-- 快速开始 -->
              <div id="quick-start" class="doc-section">
                <h2>快速开始</h2>
                <h3>安装</h3>
                <p>使用 Go modules 安装 AideCMS：</p>
                <div class="code-block">
                  <pre><code>go get github.com/zhuclark2020/aidecms</code></pre>
                </div>

                <h3>创建项目</h3>
                <p>创建一个新的 Go 项目并初始化：</p>
                <div class="code-block">
                  <pre><code>mkdir myproject && cd myproject
go mod init myproject
go get github.com/zhuclark2020/aidecms</code></pre>
                </div>

                <h3>基础示例</h3>
                <p>创建 <code>main.go</code> 文件：</p>
                <div class="code-block">
                  <pre><code>package main

import (
    "context"
    "github.com/zhuclark2020/aidecms/pkg/framework"
)

func main() {
    // 创建应用实例
    app := framework.NewApplication()
    
    // 启动应用（自动初始化所有组件）
    app.Boot()
    
    // 定义路由
    app.RegisterRoutes(func(router *framework.Router) {
        router.GET("/", func(ctx context.Context, c *framework.RequestContext) {
            c.JSON(200, map[string]interface{}{
                "message": "Hello AideCMS!",
            })
        })
    })
    
    // 运行应用
    app.Run()
}</code></pre>
                </div>
              </div>

              <!-- AI 集成 -->
              <div id="ai-integration" class="doc-section">
                <h2>AI 智能集成</h2>
                <p>AideCMS 集成了 CloudWeGo Eino 框架，提供统一的 AI 接口，支持多种主流大模型。</p>
                
                <h3>核心特性</h3>
                <ul>
                  <li>支持 OpenAI, Anthropic, 豆包, 通义千问等多种模型</li>
                  <li>统一的 AI 客户端接口</li>
                  <li>支持流式响应 (SSE)</li>
                  <li>对话上下文管理</li>
                </ul>

                <h3>使用示例</h3>
                <div class="code-block">
                  <pre><code>import "github.com/zhuclark2020/aidecms/pkg/ai"

// 创建 AI 客户端
config := &ai.Config{
    Provider: "openai",
    APIKey:   "sk-your-key",
    Model:    "gpt-4",
}
client, _ := ai.NewClient(config)

// 简单对话
response, _ := client.CreateCompletion(ctx, "Hello, AI!")
fmt.Println(response)

// 流式对话
stream, _ := client.ChatStream(ctx, "写一首诗")
for chunk := range stream {
    fmt.Print(chunk)
}</code></pre>
                </div>
              </div>

              <!-- Web3 集成 -->
              <div id="web3-integration" class="doc-section">
                <h2>Web3 区块链集成</h2>
                <p>支持 Bitcoin, Ethereum, BSC 和 Solana 四大主流区块链网络。</p>

                <h3>功能支持</h3>
                <ul>
                  <li>多链地址余额查询</li>
                  <li>交易信息查询</li>
                  <li>区块高度查询</li>
                  <li>智能合约交互 (Ethereum/BSC)</li>
                </ul>

                <h3>使用示例</h3>
                <div class="code-block">
                  <pre><code>import "github.com/zhuclark2020/aidecms/pkg/web3"

manager := web3.GetManager()

// 查询 Ethereum 余额
balance, _ := manager.GetBalance(ctx, web3.Ethereum, "0xd8dA6...")

// 查询 Bitcoin 交易
tx, _ := manager.GetTransaction(ctx, web3.Bitcoin, "tx_hash...")

// 多链查询
addresses := web3.MultiChainAddress{
    Bitcoin:  "1A1zP...",
    Ethereum: "0xd8dA6...",
}
balances, _ := addresses.GetAllBalances(ctx)</code></pre>
                </div>
              </div>

              <!-- 交易所集成 -->
              <div id="exchange-integration" class="doc-section">
                <h2>交易所集成</h2>
                <p>内置 Coinbase, KuCoin, Hyperliquid 等主流交易所 API 对接。</p>

                <h3>功能支持</h3>
                <ul>
                  <li>账户余额查询</li>
                  <li>实时行情获取</li>
                  <li>多交易所价格比较</li>
                  <li>自动化交易支持</li>
                </ul>

                <h3>使用示例</h3>
                <div class="code-block">
                  <pre><code>import "github.com/zhuclark2020/aidecms/pkg/web3"

manager := web3.GetExchangeManager()

// 查询 Coinbase BTC 余额
balance, _ := manager.GetBalance(ctx, web3.ExchangeCoinbase, "BTC")

// 查询价格
price, _ := manager.GetPrice(ctx, web3.ExchangeCoinbase, "BTC-USD")

// 比较全网价格
prices, _ := web3.GetAllExchangePrices(ctx, "BTC-USD")
for exchange, price := range prices {
    fmt.Printf("%s: $%s\n", exchange, price)
}</code></pre>
                </div>
              </div>

              <!-- 任务调度 -->
              <div id="task-schedule" class="doc-section">
                <h2>任务调度</h2>
                <p>强大的 Cron 调度系统，支持 Cron 表达式和链式调用。</p>

                <h3>使用示例</h3>
                <div class="code-block">
                  <pre><code>import "github.com/zhuclark2020/aidecms/pkg/schedule"

scheduler := schedule.NewScheduler()

// 每分钟执行
scheduler.NewTask("cleanup").EveryMinute().Do(func() error {
    return cleanup()
})

// 每天凌晨 2:00 执行
scheduler.NewTask("backup").DailyAt(2, 0).Do(func() error {
    return backupDB()
})

// Cron 表达式
scheduler.NewTask("custom").Cron("0 */2 * * *").Do(func() error {
    return customJob()
})</code></pre>
                </div>
              </div>

              <!-- 队列系统 -->
              <div id="queue-system" class="doc-section">
                <h2>队列系统</h2>
                <p>支持 Redis 和内存驱动的消息队列，处理异步任务。</p>

                <h3>使用示例</h3>
                <div class="code-block">
                  <pre><code>import "github.com/zhuclark2020/aidecms/pkg/queue"

q := queue.NewQueue("redis")

// 推送任务
q.Push(&queue.Job{
    Name: "send-email",
    Data: map[string]interface{}{
        "to": "user@example.com",
        "subject": "Welcome",
    },
    Priority: 1, // 高优先级
})

// 处理任务
q.Process("send-email", func(job *queue.Job) error {
    // 发送邮件逻辑
    return nil
})</code></pre>
                </div>
              </div>

            </div>
          </main>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const activeDocId = ref('quick-start')

const docsSections = [
  {
    title: '入门',
    items: [
      { title: '快速开始', id: 'quick-start' }
    ]
  },
  {
    title: '核心功能',
    items: [
      { title: 'AI 智能集成', id: 'ai-integration' },
      { title: 'Web3 集成', id: 'web3-integration' },
      { title: '交易所集成', id: 'exchange-integration' }
    ]
  },
  {
    title: '系统组件',
    items: [
      { title: '任务调度', id: 'task-schedule' },
      { title: '队列系统', id: 'queue-system' }
    ]
  }
]

const scrollToDoc = (id: string) => {
  activeDocId.value = id
  const element = document.getElementById(id)
  if (element) {
    element.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}
</script>

<style scoped>
.page-hero {
  padding: 140px 0 80px;
  background: var(--gradient-primary);
  color: white;
  text-align: center;
}

.page-hero h1 {
  font-size: 3rem;
  font-weight: 800;
  margin-bottom: 1rem;
}

.page-hero p {
  font-size: 1.25rem;
  opacity: 0.9;
}

.docs-content {
  padding: 80px 0;
  background: var(--bg-secondary);
}

.docs-layout {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 4rem;
  align-items: start;
}

.docs-sidebar {
  position: sticky;
  top: 100px;
  background: var(--bg-color);
  border-radius: 16px;
  padding: 2rem;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-color);
}

.nav-section {
  margin-bottom: 2rem;
}

.nav-section:last-child {
  margin-bottom: 0;
}

.nav-section h3 {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-tertiary);
  text-transform: uppercase;
  margin-bottom: 1rem;
  letter-spacing: 0.05em;
}

.nav-section ul {
  list-style: none;
}

.nav-section li {
  margin-bottom: 0.5rem;
}

.nav-section a {
  display: block;
  padding: 0.5rem 1rem;
  color: var(--text-secondary);
  text-decoration: none;
  border-radius: 8px;
  transition: all 0.3s ease;
  font-size: 0.95rem;
}

.nav-section a:hover {
  background: var(--bg-secondary);
  color: var(--primary-color);
}

.nav-section a.active {
  background: var(--gradient-primary);
  color: white;
}

.docs-main {
  background: var(--bg-color);
  border-radius: 16px;
  padding: 3rem;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-color);
}

.doc-content h2 {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--border-color);
}

.doc-section {
  margin-bottom: 4rem;
  scroll-margin-top: 100px;
}

.doc-section h3 {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 2rem 0 1rem;
}

.doc-section p {
  color: var(--text-secondary);
  line-height: 1.8;
  margin-bottom: 1rem;
}

.doc-section ul {
  margin-bottom: 1.5rem;
  padding-left: 1.5rem;
  color: var(--text-secondary);
}

.doc-section li {
  margin-bottom: 0.5rem;
  line-height: 1.6;
}

.doc-section code {
  background: var(--bg-secondary);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: var(--font-mono);
  font-size: 0.875rem;
  color: var(--primary-color);
}

.code-block {
  background: #1e293b;
  border-radius: 12px;
  padding: 1.5rem;
  margin: 1.5rem 0;
  overflow-x: auto;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.code-block pre {
  margin: 0;
  color: #e2e8f0;
  font-family: var(--font-mono);
  font-size: 0.875rem;
  line-height: 1.7;
}

@media (max-width: 968px) {
  .docs-layout {
    grid-template-columns: 1fr;
  }

  .docs-sidebar {
    position: static;
    margin-bottom: 2rem;
  }

  .docs-main {
    padding: 2rem;
  }
}
</style>
