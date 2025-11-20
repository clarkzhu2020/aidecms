<template>
  <div class="home">
    <!-- Hero Section -->
    <section class="hero">
      <div class="hero-bg">
        <div class="glow glow-1"></div>
        <div class="glow glow-2"></div>
        <div class="grid-pattern"></div>
      </div>
      <div class="container">
        <div class="hero-content" data-aos="fade-up">
          <div class="hero-badge">
            <span class="pulse"></span>
            v1.0.0 正式发布
          </div>
          <h1 class="hero-title">
            构建下一代
            <span class="text-gradient">AI 驱动</span>
            <br>内容管理系统
          </h1>
          <p class="hero-subtitle">
            基于 CloudWeGo Hertz 的企业级 Go CMS 平台框架。
            <br>融合 AI 智能、Web3 支付与交易所集成，为您的业务赋能。
          </p>
          <div class="hero-actions">
            <a href="https://github.com/chenyusolar/aidecms" class="btn btn-primary" target="_blank">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
              </svg>
              开始使用
            </a>
            <router-link to="/docs" class="btn btn-secondary">
              阅读文档
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M5 12h14M12 5l7 7-7 7" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </router-link>
          </div>
          <div class="hero-stats">
            <div class="stat">
              <div class="stat-value">100%</div>
              <div class="stat-label">Go 语言</div>
            </div>
            <div class="stat-divider"></div>
            <div class="stat">
              <div class="stat-value">MIT</div>
              <div class="stat-label">开源协议</div>
            </div>
            <div class="stat-divider"></div>
            <div class="stat">
              <div class="stat-value">AI</div>
              <div class="stat-label">原生驱动</div>
            </div>
          </div>
        </div>
        <div class="hero-image" data-aos="fade-left" data-aos-delay="200">
          <div class="code-window glass">
            <div class="window-header">
              <div class="window-dots">
                <span></span><span></span><span></span>
              </div>
              <div class="window-title">main.go</div>
            </div>
            <pre class="code-content"><code>package main

import (
    "github.com/chenyusolar/aidecms/pkg/framework"
    "github.com/chenyusolar/aidecms/routes"
)

func main() {
    // 初始化应用
    app := framework.NewApplication()
    
    // 启动核心组件
    app.Boot()
    
    // 注册路由
    routes.RegisterRoutes(app)
    
    // 运行服务 🚀
    app.Run()
}</code></pre>
          </div>
        </div>
      </div>
    </section>

    <!-- Innovation Highlights Section -->
    <section class="highlights">
      <div class="container">
        <div class="section-header" data-aos="fade-up">
          <h2>创新亮点</h2>
          <p>融合前沿技术，打造极致体验</p>
        </div>
        <div class="highlights-grid">
          <div class="highlight-card" data-aos="fade-up" data-aos-delay="0">
            <div class="card-bg"></div>
            <div class="highlight-icon ai">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2a10 10 0 1 0 10 10A10 10 0 0 0 12 2zm0 14a4 4 0 1 1 4-4 4 4 0 0 1-4 4z"/></svg>
            </div>
            <h3>AI 智能驱动</h3>
            <p>内置多模型 AI 能力，支持智能内容生成、自动摘要、智能标签、智能推荐等，助力内容创作与管理效率飞跃。</p>
          </div>
          <div class="highlight-card" data-aos="fade-up" data-aos-delay="100">
            <div class="card-bg"></div>
            <div class="highlight-icon web3">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M20 12V8H6a2 2 0 0 1-2-2c0-1.1.9-2 2-2h12v4"/><path d="M4 6v12c0 1.1.9 2 2 2h14v-4"/><path d="M18 12a2 2 0 0 0-2 2c0 1.1.9 2 2 2h4v-4h-4z"/></svg>
            </div>
            <h3>Web3 原生集成</h3>
            <p>支持钱包登录、NFT 内容确权、链上数据同步，轻松打造去中心化内容平台，拥抱未来互联网。</p>
          </div>
          <div class="highlight-card" data-aos="fade-up" data-aos-delay="200">
            <div class="card-bg"></div>
            <div class="highlight-icon exchange">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 3v18h18"/><path d="M18.7 8l-5.1 5.2-2.8-2.7L7 14.3"/></svg>
            </div>
            <h3>交易所级集成</h3>
            <p>内置主流交易所 API 对接，支持行情展示、资产管理、自动化内容驱动交易等创新场景。</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="features">
      <div class="container">
        <div class="section-header" data-aos="fade-up">
          <h2>强大的功能特性</h2>
          <p>为现代 Web 应用提供完整的解决方案</p>
        </div>

        <div class="features-grid">
          <div class="feature-card" v-for="(feature, index) in features" :key="index" data-aos="fade-up" :data-aos-delay="index * 50">
            <div class="feature-icon">
              <component :is="feature.icon" />
            </div>
            <h3>{{ feature.title }}</h3>
            <p>{{ feature.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="cta">
      <div class="cta-bg">
        <div class="glow glow-3"></div>
      </div>
      <div class="container">
        <div class="cta-content" data-aos="zoom-in">
          <h2>准备好开始了吗？</h2>
          <p>立即使用 AideCMS 构建您的下一个项目</p>
          <div class="cta-actions">
            <a href="https://github.com/chenyusolar/aidecms" class="btn btn-primary btn-large" target="_blank">
              查看 GitHub
            </a>
            <router-link to="/docs" class="btn btn-secondary btn-large">
              阅读文档
            </router-link>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { h } from 'vue'

const features = [
  {
    title: '高性能架构',
    desc: '基于 CloudWeGo Hertz 框架，提供出色的并发性能和低延迟响应',
    icon: h('svg', { width: 24, height: 24, viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': 2 }, [
      h('path', { d: 'M13 2L3 14h9l-1 8 10-12h-9l1-8z' })
    ])
  },
  {
    title: 'AI 智能集成',
    desc: '内置多种 AI 模型支持，轻松实现智能内容生成和分析',
    icon: h('svg', { width: 24, height: 24, viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': 2 }, [
      h('path', { d: 'M12 2v20M2 12h20' }),
      h('circle', { cx: 12, cy: 12, r: 3 })
    ])
  },
  {
    title: '完整 CMS 系统',
    desc: '文章、分类、标签、媒体管理等核心功能开箱即用',
    icon: h('svg', { width: 24, height: 24, viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': 2 }, [
      h('rect', { x: 3, y: 3, width: 18, height: 18, rx: 2 }),
      h('path', { d: 'M9 3v18M3 9h18M3 15h18' })
    ])
  },
  {
    title: '任务调度',
    desc: '强大的 Cron 调度系统，支持并发执行和任务统计',
    icon: h('svg', { width: 24, height: 24, viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': 2 }, [
      h('circle', { cx: 12, cy: 12, r: 10 }),
      h('polyline', { points: '12 6 12 12 16 14' })
    ])
  },
  {
    title: '队列系统',
    desc: '支持 Redis 和内存驱动，提供延迟任务和失败重试机制',
    icon: h('svg', { width: 24, height: 24, viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': 2 }, [
      h('path', { d: 'M4 6h16M4 12h16M4 18h16' })
    ])
  },
  {
    title: '安全可靠',
    desc: 'JWT 认证、权限管理、限流保护，全方位保障系统安全',
    icon: h('svg', { width: 24, height: 24, viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': 2 }, [
      h('path', { d: 'M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z' })
    ])
  }
]
</script>

<style scoped>
/* Hero Section */
.hero {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  padding: 120px 0 80px;
  overflow: hidden;
}

.hero-bg {
  position: absolute;
  inset: 0;
  z-index: -1;
  overflow: hidden;
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
  opacity: 0.4;
  animation: float 10s ease-in-out infinite;
}

.glow-1 {
  top: -10%;
  right: -10%;
  width: 500px;
  height: 500px;
  background: var(--primary-color);
}

.glow-2 {
  bottom: -10%;
  left: -10%;
  width: 400px;
  height: 400px;
  background: var(--secondary-color);
  animation-delay: -5s;
}

.hero .container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 4rem;
  align-items: center;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: rgba(99, 102, 241, 0.1);
  color: var(--primary-color);
  border-radius: 20px;
  font-size: 0.875rem;
  font-weight: 600;
  margin-bottom: 1.5rem;
  border: 1px solid rgba(99, 102, 241, 0.2);
}

.pulse {
  width: 8px;
  height: 8px;
  background: var(--primary-color);
  border-radius: 50%;
  box-shadow: 0 0 0 rgba(99, 102, 241, 0.4);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { box-shadow: 0 0 0 0 rgba(99, 102, 241, 0.4); }
  70% { box-shadow: 0 0 0 10px rgba(99, 102, 241, 0); }
  100% { box-shadow: 0 0 0 0 rgba(99, 102, 241, 0); }
}

.hero-title {
  font-size: 4rem;
  line-height: 1.1;
  margin-bottom: 1.5rem;
  background: linear-gradient(to bottom, var(--text-primary), var(--text-secondary));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-subtitle {
  font-size: 1.25rem;
  color: var(--text-secondary);
  margin-bottom: 2.5rem;
  line-height: 1.6;
}

.hero-actions {
  display: flex;
  gap: 1rem;
  margin-bottom: 3rem;
}

.hero-stats {
  display: flex;
  align-items: center;
  gap: 2rem;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-primary);
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-tertiary);
}

.stat-divider {
  width: 1px;
  height: 40px;
  background: var(--border-color);
}

/* Code Window */
.code-window {
  background: #1e293b;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: var(--shadow-xl);
  border: 1px solid rgba(255, 255, 255, 0.1);
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

.code-content {
  padding: 24px;
  color: #e2e8f0;
  font-family: var(--font-mono);
  font-size: 0.875rem;
  line-height: 1.6;
  overflow-x: auto;
}

/* Highlights Section */
.highlights {
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

.highlights-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 2rem;
}

.highlight-card {
  position: relative;
  background: var(--bg-color);
  padding: 2.5rem;
  border-radius: 16px;
  box-shadow: var(--shadow-md);
  transition: all 0.3s ease;
  overflow: hidden;
  border: 1px solid var(--border-color);
}

.highlight-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-xl);
  border-color: var(--primary-color);
}

.highlight-icon {
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  margin-bottom: 1.5rem;
  color: white;
}

.highlight-icon.ai { background: linear-gradient(135deg, #818cf8, #6366f1); }
.highlight-icon.web3 { background: linear-gradient(135deg, #38bdf8, #0ea5e9); }
.highlight-icon.exchange { background: linear-gradient(135deg, #fbbf24, #f59e0b); }

.highlight-card h3 {
  margin-bottom: 1rem;
  font-size: 1.25rem;
}

.highlight-card p {
  color: var(--text-secondary);
  font-size: 0.95rem;
}

/* Features Section */
.features {
  padding: 100px 0;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 2rem;
}

.feature-card {
  padding: 2rem;
  border-radius: 12px;
  background: var(--bg-secondary);
  transition: all 0.3s ease;
}

.feature-card:hover {
  background: white;
  box-shadow: var(--shadow-lg);
  transform: translateY(-3px);
}

.feature-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  border-radius: 12px;
  color: var(--primary-color);
  margin-bottom: 1rem;
  box-shadow: var(--shadow-sm);
}

.feature-card h3 {
  font-size: 1.125rem;
  margin-bottom: 0.5rem;
}

.feature-card p {
  color: var(--text-secondary);
  font-size: 0.9rem;
}

/* CTA Section */
.cta {
  position: relative;
  padding: 100px 0;
  overflow: hidden;
}

.cta-bg {
  position: absolute;
  inset: 0;
  background: var(--bg-secondary);
  z-index: -1;
}

.glow-3 {
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 600px;
  height: 600px;
  background: var(--accent-color);
  opacity: 0.1;
}

.cta-content {
  text-align: center;
  max-width: 600px;
  margin: 0 auto;
}

.cta-content h2 {
  margin-bottom: 1rem;
}

.cta-content p {
  color: var(--text-secondary);
  margin-bottom: 2.5rem;
  font-size: 1.25rem;
}

.cta-actions {
  display: flex;
  justify-content: center;
  gap: 1rem;
}

.btn-large {
  padding: 1rem 2rem;
  font-size: 1.125rem;
}

/* Responsive */
@media (max-width: 968px) {
  .hero .container {
    grid-template-columns: 1fr;
    text-align: center;
  }

  .hero-actions, .hero-stats {
    justify-content: center;
  }

  .hero-image {
    margin-top: 3rem;
  }

  .highlights-grid, .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .hero-title {
    font-size: 2.5rem;
  }

  .highlights-grid, .features-grid {
    grid-template-columns: 1fr;
  }

  .cta-actions {
    flex-direction: column;
  }
}
</style>
