<template>
  <div class="login-container">
    <div class="login-background">
      <div class="gradient-orb orb-1"></div>
      <div class="gradient-orb orb-2"></div>
      <div class="gradient-orb orb-3"></div>
    </div>

    <div class="login-content">
      <div class="login-card">
        <div class="login-header">
          <div class="logo-container">
            <span class="logo-icon">🚀</span>
            <h1>AideCMS</h1>
          </div>
          <p class="subtitle">Sign in to your account</p>
        </div>

        <el-form :model="form" class="login-form" @submit.prevent="handleLogin">
          <el-form-item>
            <el-input
              v-model="form.username"
              placeholder="Username"
              size="large"
              :prefix-icon="User"
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item>
            <el-input
              v-model="form.password"
              type="password"
              placeholder="Password"
              size="large"
              show-password
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>

          <el-form-item>
            <div class="form-options">
              <el-checkbox v-model="rememberMe">Remember me</el-checkbox>
              <el-link type="primary" :underline="false">Forgot password?</el-link>
            </div>
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              size="large"
              class="login-button"
              @click="handleLogin"
              :loading="loading"
            >
              Sign In
            </el-button>
          </el-form-item>
        </el-form>

        <div class="login-footer">
          <p>Default credentials: <strong>admin / admin</strong></p>
        </div>
      </div>

      <div class="features-grid">
        <div class="feature-item" v-for="(feature, index) in features" :key="index">
          <el-icon :size="24"><component :is="feature.icon" /></el-icon>
          <span>{{ feature.text }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { ElMessage } from 'element-plus'
import { User, Lock, TrendCharts, Wallet, ChatDotRound, Setting } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const rememberMe = ref(false)

const form = ref({
  username: '',
  password: ''
})

const features = ref([
  { icon: TrendCharts, text: 'Analytics Dashboard' },
  { icon: Wallet, text: 'Web3 Integration' },
  { icon: ChatDotRound, text: 'AI Assistant' },
  { icon: Setting, text: 'Full Control' }
])

const handleLogin = async () => {
  loading.value = true
  // Mock login
  setTimeout(() => {
    if (form.value.username === 'admin' && form.value.password === 'admin') {
      userStore.login({ username: 'admin', role: 'admin' }, 'mock-token')
      ElMessage.success('Login successful')
      router.push('/')
    } else {
      ElMessage.error('Invalid credentials (try admin/admin)')
    }
    loading.value = false
  }, 1000)
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

/* ========================================
   Animated Background
   ======================================== */

.login-background {
  position: absolute;
  inset: 0;
  overflow: hidden;
}

.gradient-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.6;
  animation: float 20s infinite ease-in-out;
}

.orb-1 {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.8), transparent);
  top: -250px;
  left: -250px;
  animation-delay: 0s;
}

.orb-2 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(139, 92, 246, 0.8), transparent);
  bottom: -200px;
  right: -200px;
  animation-delay: 7s;
}

.orb-3 {
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(236, 72, 153, 0.8), transparent);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: 14s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  33% {
    transform: translate(30px, -30px) scale(1.1);
  }
  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }
}

/* ========================================
   Login Content
   ======================================== */

.login-content {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 440px;
  padding: 2rem;
  animation: slideUp 0.6s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* ========================================
   Login Card
   ======================================== */

.login-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  padding: 3rem 2.5rem;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.login-header {
  text-align: center;
  margin-bottom: 2.5rem;
}

.logo-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  margin-bottom: 0.5rem;
}

.logo-icon {
  font-size: 3rem;
}

.login-header h1 {
  font-size: 2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
}

.subtitle {
  color: var(--gray-600);
  font-size: 0.9375rem;
  margin: 0;
}

/* ========================================
   Login Form
   ======================================== */

.login-form {
  margin-bottom: 1.5rem;
}

.login-form :deep(.el-input__wrapper) {
  padding: 12px 16px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  transition: all 0.3s;
}

.login-form :deep(.el-input__wrapper:hover) {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.login-form :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.login-button {
  width: 100%;
  height: 48px;
  font-size: 1rem;
  font-weight: 600;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
}

.login-button:active {
  transform: translateY(0);
}

/* ========================================
   Login Footer
   ======================================== */

.login-footer {
  text-align: center;
  padding-top: 1.5rem;
  border-top: 1px solid var(--gray-200);
}

.login-footer p {
  color: var(--gray-600);
  font-size: 0.875rem;
  margin: 0;
}

.login-footer strong {
  color: var(--primary);
}

/* ========================================
   Features Grid
   ======================================== */

.features-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
  margin-top: 2rem;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: white;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.3s;
}

.feature-item:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

/* ========================================
   Responsive Design
   ======================================== */

@media (max-width: 640px) {
  .login-content {
    padding: 1rem;
  }

  .login-card {
    padding: 2rem 1.5rem;
  }

  .login-header h1 {
    font-size: 1.75rem;
  }

  .features-grid {
    grid-template-columns: 1fr;
  }
}
</style>
