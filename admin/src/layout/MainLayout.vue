<template>
  <div class="admin-layout">
    <!-- Sidebar -->
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="sidebar-header">
        <div class="logo">
          <span class="logo-icon">🚀</span>
          <span v-if="!sidebarCollapsed" class="logo-text">AideCMS</span>
        </div>
        <button class="collapse-btn" @click="toggleSidebar">
          <el-icon><Fold v-if="!sidebarCollapsed" /><Expand v-else /></el-icon>
        </button>
      </div>

      <nav class="sidebar-nav">
        <el-menu
          :default-active="$route.path"
          :collapse="sidebarCollapsed"
          router
          background-color="transparent"
          text-color="rgba(255, 255, 255, 0.8)"
          active-text-color="#ffffff"
        >
          <el-menu-item index="/dashboard" class="menu-item">
            <el-icon><Odometer /></el-icon>
            <template #title>Dashboard</template>
          </el-menu-item>
          
          <el-menu-item index="/users" class="menu-item">
            <el-icon><User /></el-icon>
            <template #title>Users</template>
          </el-menu-item>
          
          <el-menu-item index="/posts" class="menu-item">
            <el-icon><Document /></el-icon>
            <template #title>Posts</template>
          </el-menu-item>
          
          <el-sub-menu index="/web3" class="menu-item">
            <template #title>
              <el-icon><Wallet /></el-icon>
              <span>Web3</span>
            </template>
            <el-menu-item index="/web3/wallet">Wallet</el-menu-item>
            <el-menu-item index="/web3/transactions">Transactions</el-menu-item>
          </el-sub-menu>
        </el-menu>
      </nav>

      <div class="sidebar-footer">
        <div class="user-profile" v-if="!sidebarCollapsed">
          <el-avatar :size="40" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
          <div class="user-info">
            <div class="user-name">Admin</div>
            <div class="user-role">Administrator</div>
          </div>
        </div>
      </div>
    </aside>

    <!-- Main Content -->
    <div class="main-wrapper">
      <!-- Header -->
      <header class="header">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">
              <el-icon><HomeFilled /></el-icon>
              Home
            </el-breadcrumb-item>
            <el-breadcrumb-item v-if="$route.meta.title">
              {{ $route.meta.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="header-right">
          <el-badge :value="3" class="notification-badge">
            <el-button :icon="Bell" circle />
          </el-badge>
          
          <el-dropdown @command="handleCommand" trigger="click">
            <div class="user-dropdown">
              <el-avatar size="small" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
              <span class="username">Admin</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>
                  <el-icon><User /></el-icon>
                  Profile
                </el-dropdown-item>
                <el-dropdown-item>
                  <el-icon><Setting /></el-icon>
                  Settings
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  Logout
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>

      <!-- Page Content -->
      <main class="content">
        <RouterView v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </RouterView>
      </main>
    </div>

    <!-- AI Chat Component -->
    <AIChat />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Bell, HomeFilled, ArrowDown, Setting, SwitchButton, Fold, Expand } from '@element-plus/icons-vue'
// import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import AIChat from '../components/ai/AIChat.vue'

const userStore = useUserStore()
const sidebarCollapsed = ref(false)

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

const handleCommand = (command: string) => {
  if (command === 'logout') {
    userStore.logout()
  }
}
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  width: 100vw;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  overflow: hidden;
}

/* ========================================
   Sidebar Styles
   ======================================== */

.sidebar {
  width: 260px;
  min-width: 260px;
  background: linear-gradient(180deg, rgba(30, 41, 59, 0.95) 0%, rgba(15, 23, 42, 0.95) 100%);
  backdrop-filter: blur(20px);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  flex-direction: column;
  transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  z-index: 100;
  height: 100vh;
  overflow-y: auto;
}

.sidebar.collapsed {
  width: 64px;
  min-width: 64px;
}

.sidebar-header {
  padding: 1.5rem 1rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  color: white;
  font-weight: 600;
  font-size: 1.25rem;
}

.logo-icon {
  font-size: 1.5rem;
}

.logo-text {
  white-space: nowrap;
  overflow: hidden;
}

.collapse-btn {
  background: rgba(255, 255, 255, 0.1);
  border: none;
  color: white;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.collapse-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.sidebar-nav {
  flex: 1;
  padding: 1rem 0;
  overflow-y: auto;
}

.sidebar-nav :deep(.el-menu) {
  border: none;
}

.sidebar-nav :deep(.el-menu-item),
.sidebar-nav :deep(.el-sub-menu__title) {
  margin: 0.25rem 0.75rem;
  border-radius: 8px;
  transition: all 0.2s;
}

.sidebar-nav :deep(.el-menu-item:hover),
.sidebar-nav :deep(.el-sub-menu__title:hover) {
  background: rgba(255, 255, 255, 0.1) !important;
}

.sidebar-nav :deep(.el-menu-item.is-active) {
  background: linear-gradient(90deg, rgba(99, 102, 241, 0.8), rgba(139, 92, 246, 0.8)) !important;
  color: white !important;
}

.sidebar-footer {
  padding: 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name {
  color: white;
  font-weight: 500;
  font-size: 0.875rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-role {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.75rem;
}

/* ========================================
   Main Wrapper
   ======================================== */

.main-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  height: 100vh;
  overflow: hidden;
}

/* ========================================
   Header Styles
   ======================================== */

.header {
  height: 64px;
  min-height: 64px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 2rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  flex-shrink: 0;
}

.header-left :deep(.el-breadcrumb) {
  font-size: 0.875rem;
}

.header-left :deep(.el-breadcrumb__item) {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.notification-badge {
  cursor: pointer;
}

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  padding: 0.5rem 0.75rem;
  border-radius: 8px;
  transition: background 0.2s;
}

.user-dropdown:hover {
  background: rgba(0, 0, 0, 0.05);
}

.username {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--gray-700);
}

/* ========================================
   Content Area
   ======================================== */

.content {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
  background: rgba(249, 250, 251, 0.5);
  height: calc(100vh - 64px);
}

/* ========================================
   Transitions
   ======================================== */

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* ========================================
   Responsive Design
   ======================================== */

@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    z-index: 1000;
  }
  
  .sidebar.collapsed {
    transform: translateX(-100%);
  }
  
  .header {
    padding: 0 1rem;
  }
  
  .content {
    padding: 1rem;
  }
}
</style>
