<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <div>
        <h1 class="page-title">Dashboard</h1>
        <p class="page-subtitle">Welcome back! Here's what's happening today.</p>
      </div>
      <el-button type="primary" :icon="Refresh" @click="refreshData">Refresh</el-button>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card" v-for="(stat, index) in stats" :key="index" :class="`stat-card-${index + 1}`">
        <div class="stat-icon">
          <el-icon :size="32"><component :is="stat.icon" /></el-icon>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ stat.label }}</div>
          <div class="stat-value">{{ stat.value }}</div>
          <div class="stat-change" :class="stat.trend">
            <el-icon><component :is="stat.trend === 'up' ? TrendCharts : Bottom" /></el-icon>
            {{ stat.change }}
          </div>
        </div>
      </div>
    </div>

    <!-- Charts and Activity -->
    <el-row :gutter="24" class="dashboard-row">
      <el-col :xs="24" :lg="16">
        <div class="dashboard-card">
          <div class="card-header">
            <h3>Activity Overview</h3>
            <el-radio-group v-model="chartPeriod" size="small">
              <el-radio-button label="week">Week</el-radio-button>
              <el-radio-button label="month">Month</el-radio-button>
              <el-radio-button label="year">Year</el-radio-button>
            </el-radio-group>
          </div>
          <div class="chart-placeholder">
            <el-icon :size="64" color="#e5e7eb"><TrendCharts /></el-icon>
            <p>Chart visualization would go here</p>
            <p class="chart-hint">Integrate Chart.js or ECharts for data visualization</p>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :lg="8">
        <div class="dashboard-card">
          <div class="card-header">
            <h3>Recent Activity</h3>
          </div>
          <el-timeline class="activity-timeline">
            <el-timeline-item
              v-for="activity in recentActivities"
              :key="activity.id"
              :timestamp="activity.time"
              :color="activity.color"
            >
              <div class="activity-item">
                <strong>{{ activity.title }}</strong>
                <p>{{ activity.description }}</p>
              </div>
            </el-timeline-item>
          </el-timeline>
        </div>
      </el-col>
    </el-row>

    <!-- Quick Actions and System Status -->
    <el-row :gutter="24" class="dashboard-row">
      <el-col :xs="24" :lg="12">
        <div class="dashboard-card">
          <div class="card-header">
            <h3>Quick Actions</h3>
          </div>
          <div class="quick-actions">
            <el-button type="primary" :icon="Plus" size="large">New Post</el-button>
            <el-button type="success" :icon="UserFilled" size="large">Add User</el-button>
            <el-button type="warning" :icon="Delete" size="large">Clear Cache</el-button>
            <el-button type="info" :icon="Setting" size="large">Settings</el-button>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :lg="12">
        <div class="dashboard-card">
          <div class="card-header">
            <h3>System Status</h3>
          </div>
          <div class="system-status">
            <div class="status-item" v-for="item in systemStatus" :key="item.name">
              <div class="status-info">
                <span class="status-name">{{ item.name }}</span>
                <span class="status-value">{{ item.value }}</span>
              </div>
              <el-progress
                :percentage="item.percentage"
                :color="item.color"
                :show-text="false"
              />
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  User,
  Document,
  ChatDotRound,
  TrendCharts,
  Bottom,
  Refresh,
  Plus,
  UserFilled,
  Delete,
  Setting
} from '@element-plus/icons-vue'

const chartPeriod = ref('week')

const stats = ref([
  {
    label: 'Total Users',
    value: '1,234',
    change: '+12.5%',
    trend: 'up',
    icon: User
  },
  {
    label: 'Total Posts',
    value: '567',
    change: '+8.2%',
    trend: 'up',
    icon: Document
  },
  {
    label: 'Comments',
    value: '890',
    change: '+15.3%',
    trend: 'up',
    icon: ChatDotRound
  },
  {
    label: 'System Load',
    value: '45%',
    change: '-3.1%',
    trend: 'down',
    icon: TrendCharts
  }
])

const recentActivities = ref([
  {
    id: 1,
    title: 'New user registered',
    description: 'john@example.com joined the platform',
    time: '2 minutes ago',
    color: '#67c23a'
  },
  {
    id: 2,
    title: 'Post published',
    description: '"Getting Started with Vue 3" was published',
    time: '15 minutes ago',
    color: '#409eff'
  },
  {
    id: 3,
    title: 'Comment received',
    description: 'New comment on "Introduction to TypeScript"',
    time: '1 hour ago',
    color: '#e6a23c'
  },
  {
    id: 4,
    title: 'System update',
    description: 'Database backup completed successfully',
    time: '3 hours ago',
    color: '#909399'
  }
])

const systemStatus = ref([
  { name: 'CPU Usage', value: '45%', percentage: 45, color: '#67c23a' },
  { name: 'Memory', value: '62%', percentage: 62, color: '#e6a23c' },
  { name: 'Disk Space', value: '78%', percentage: 78, color: '#f56c6c' },
  { name: 'Network', value: '23%', percentage: 23, color: '#409eff' }
])

const refreshData = () => {
  console.log('Refreshing dashboard data...')
}
</script>

<style scoped>
.dashboard {
  animation: fadeIn 0.3s ease-out;
  height: 100%;
  display: flex;
  flex-direction: column;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* ========================================
   Dashboard Header
   ======================================== */

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
  flex-shrink: 0;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--gray-900);
  margin: 0 0 0.5rem 0;
}

.page-subtitle {
  color: var(--gray-600);
  margin: 0;
}

/* ========================================
   Stats Grid
   ======================================== */

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  margin-bottom: 1.5rem;
  flex-shrink: 0;
}

.stat-card {
  background: white;
  border-radius: 16px;
  padding: 1.5rem;
  display: flex;
  gap: 1rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--gradient-color-1), var(--gradient-color-2));
}

.stat-card-1 {
  --gradient-color-1: #667eea;
  --gradient-color-2: #764ba2;
}

.stat-card-2 {
  --gradient-color-1: #f093fb;
  --gradient-color-2: #f5576c;
}

.stat-card-3 {
  --gradient-color-1: #4facfe;
  --gradient-color-2: #00f2fe;
}

.stat-card-4 {
  --gradient-color-1: #43e97b;
  --gradient-color-2: #38f9d7;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--gradient-color-1), var(--gradient-color-2));
  color: white;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.stat-label {
  font-size: 0.875rem;
  color: var(--gray-600);
  margin-bottom: 0.25rem;
}

.stat-value {
  font-size: 1.875rem;
  font-weight: 700;
  color: var(--gray-900);
  margin-bottom: 0.25rem;
}

.stat-change {
  font-size: 0.875rem;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.stat-change.up {
  color: #67c23a;
}

.stat-change.down {
  color: #f56c6c;
}

/* ========================================
   Dashboard Cards
   ======================================== */

.dashboard-row {
  margin-bottom: 1.5rem;
  flex-shrink: 0;
}

.dashboard-card {
  background: white;
  border-radius: 16px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  height: 100%;
  display: flex;
  flex-direction: column;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--gray-200);
}

.card-header h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--gray-900);
  margin: 0;
}

/* ========================================
   Chart Placeholder
   ======================================== */

.chart-placeholder {
  height: 300px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  border-radius: 12px;
  color: var(--gray-600);
}

.chart-hint {
  font-size: 0.875rem;
  color: var(--gray-500);
  margin-top: 0.5rem;
}

/* ========================================
   Activity Timeline
   ======================================== */

.activity-timeline {
  padding: 0;
}

.activity-item strong {
  color: var(--gray-900);
  display: block;
  margin-bottom: 0.25rem;
}

.activity-item p {
  color: var(--gray-600);
  font-size: 0.875rem;
  margin: 0;
}

/* ========================================
   Quick Actions
   ======================================== */

.quick-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 1rem;
}

.quick-actions .el-button {
  width: 100%;
}

/* ========================================
   System Status
   ======================================== */

.system-status {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.status-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.status-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-name {
  font-size: 0.875rem;
  color: var(--gray-700);
  font-weight: 500;
}

.status-value {
  font-size: 0.875rem;
  color: var(--gray-600);
}

/* ========================================
   Responsive Design
   ======================================== */

@media (max-width: 768px) {
  .dashboard-header {
    flex-direction: column;
    gap: 1rem;
  }

  .page-title {
    font-size: 1.5rem;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .quick-actions {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
