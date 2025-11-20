<template>
  <footer class="footer">
    <div class="container">
      <div class="footer-content">
        <div class="footer-section">
          <div class="logo">
            <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
              <rect width="32" height="32" rx="8" fill="url(#gradient)" />
              <path d="M10 16L14 20L22 12" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              <defs>
                <linearGradient id="gradient" x1="0" y1="0" x2="32" y2="32">
                  <stop offset="0%" stop-color="#667eea"/>
                  <stop offset="100%" stop-color="#764ba2"/>
                </linearGradient>
              </defs>
            </svg>
            <span>AideCMS</span>
          </div>
          <p class="description">
            企业级 Go CMS 平台框架<br>
            基于 CloudWeGo Hertz 的高性能解决方案
          </p>
          <div class="social-links">
            <a href="https://github.com/clarkzhu2020/aidecms" target="_blank" aria-label="GitHub">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
              </svg>
            </a>
          </div>
        </div>

        <div class="footer-section">
          <h3>产品</h3>
          <ul>
            <li><router-link to="/features">功能特性</router-link></li>
            <li><router-link to="/docs">文档</router-link></li>
            <li><a href="https://github.com/clarkzhu2020/aidecms/releases" target="_blank">版本发布</a></li>
            <li><a href="https://github.com/clarkzhu2020/aidecms/blob/main/README.md" target="_blank">快速开始</a></li>
          </ul>
        </div>

        <div class="footer-section">
          <h3>资源</h3>
          <ul>
            <li><a href="https://github.com/clarkzhu2020/aidecms/tree/main/doc" target="_blank">API 文档</a></li>
            <li><a href="https://github.com/clarkzhu2020/aidecms/tree/main/example" target="_blank">示例代码</a></li>
            <li><a href="https://github.com/clarkzhu2020/aidecms/issues" target="_blank">问题反馈</a></li>
            <li><a href="https://github.com/clarkzhu2020/aidecms/discussions" target="_blank">社区讨论</a></li>
          </ul>
        </div>

        <div class="footer-section">
          <h3>关于</h3>
          <ul>
            <li><a href="https://github.com/clarkzhu2020/aidecms" target="_blank">GitHub</a></li>
            <li><a href="https://github.com/clarkzhu2020/aidecms/blob/main/LICENSE" target="_blank">开源协议</a></li>
            <li><a href="https://github.com/clarkzhu2020/aidecms/graphs/contributors" target="_blank">贡献者</a></li>
          </ul>
        </div>
      </div>

      <div class="footer-bottom">
        <p>&copy; 2025 AideCMS. All rights reserved.</p>
        <p>Built with ❤️ using CloudWeGo Hertz & Go</p>
      </div>
    </div>
  </footer>
</template>

<style scoped>
.footer {
  background: var(--dark-bg);
  color: #94a3b8;
  padding: 4rem 0 2rem;
  margin-top: auto;
}

.footer-content {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 3rem;
  margin-bottom: 3rem;
}

.footer-section .logo {
  display: flex;
  align-items: center;
  gap: 12px;
  color: white;
  font-weight: 700;
  font-size: 1.5rem;
  margin-bottom: 1rem;
}

.footer-section .description {
  margin-bottom: 1.5rem;
  line-height: 1.6;
}

.footer-section h3 {
  color: white;
  font-size: 1rem;
  font-weight: 600;
  margin-bottom: 1rem;
}

.footer-section ul {
  list-style: none;
}

.footer-section ul li {
  margin-bottom: 0.75rem;
}

.footer-section a {
  color: #94a3b8;
  text-decoration: none;
  transition: color 0.3s ease;
}

.footer-section a:hover {
  color: var(--accent-color);
}

.social-links {
  display: flex;
  gap: 1rem;
}

.social-links a {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  transition: all 0.3s ease;
}

.social-links a:hover {
  background: var(--gradient-primary);
  transform: translateY(-2px);
}

.footer-bottom {
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  padding-top: 2rem;
  text-align: center;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

@media (max-width: 968px) {
  .footer-content {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .footer-content {
    grid-template-columns: 1fr;
  }

  .footer-bottom {
    flex-direction: column;
    gap: 1rem;
  }
}
</style>
