<script setup lang="ts">
// 仪表盘首页
import { ref } from 'vue'

// 概览卡片数据
const overviewCards = ref([
  { title: '主机总数', value: '128', icon: 'Monitor', color: '#667eea', trend: '+5', subtitle: '在线 120 / 离线 8' },
  { title: 'K8s 集群', value: '6', icon: 'Grid', color: '#764ba2', trend: '+1', subtitle: '健康 5 / 异常 1' },
  { title: '活跃告警', value: '23', icon: 'Bell', color: '#f59e0b', trend: '-3', subtitle: 'P1: 2 / P2: 5 / P3: 16' },
  { title: '今日发布', value: '12', icon: 'Upload', color: '#10b981', trend: '+4', subtitle: '成功率 100%' },
])
</script>

<template>
  <div class="dashboard">
    <!-- 欢迎横幅 -->
    <div class="welcome-banner glass-card">
      <div class="welcome-content">
        <h2>欢迎使用 <span class="gradient-text">智维云枢</span></h2>
        <p>AI 驱动的智能运维管理平台，让运维更轻松</p>
      </div>
      <div class="welcome-decoration">
        <el-icon :size="80" style="opacity: 0.1; color: var(--co-primary)"><Cpu /></el-icon>
      </div>
    </div>

    <!-- 概览卡片 -->
    <div class="overview-cards">
      <div
        v-for="card in overviewCards"
        :key="card.title"
        class="overview-card glass-card"
      >
        <div class="card-header">
          <span class="card-title">{{ card.title }}</span>
          <div class="card-icon" :style="{ background: card.color + '20', color: card.color }">
            <el-icon :size="20"><component :is="card.icon" /></el-icon>
          </div>
        </div>
        <div class="card-value">
          {{ card.value }}
          <span class="card-trend" :class="{ positive: card.trend.startsWith('+') }">
            {{ card.trend }}
          </span>
        </div>
        <div class="card-subtitle">{{ card.subtitle }}</div>
      </div>
    </div>

    <!-- 占位图表区域 -->
    <div class="charts-row">
      <div class="chart-panel glass-card">
        <h3 class="panel-title">告警趋势 (近7天)</h3>
        <div class="chart-placeholder">
          <el-icon :size="48" style="color: var(--co-text-muted)"><TrendCharts /></el-icon>
          <p>图表组件将在后续开发中实现</p>
        </div>
      </div>
      <div class="chart-panel glass-card">
        <h3 class="panel-title">资源使用率</h3>
        <div class="chart-placeholder">
          <el-icon :size="48" style="color: var(--co-text-muted)"><PieChart /></el-icon>
          <p>图表组件将在后续开发中实现</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard {
  animation: fadeInUp 0.5s ease-out;
}

.welcome-banner {
  padding: 28px 32px;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  overflow: hidden;
  position: relative;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.12) 0%, rgba(118, 75, 162, 0.08) 100%) !important;
}

.welcome-content h2 {
  font-size: 22px;
  margin-bottom: 8px;
}

.welcome-content p {
  color: var(--co-text-secondary);
  font-size: 14px;
}

.overview-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.overview-card {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.card-title {
  color: var(--co-text-secondary);
  font-size: 13px;
}

.card-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.card-value {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 4px;
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.card-trend {
  font-size: 13px;
  color: #ef4444;
  font-weight: 500;
}

.card-trend.positive {
  color: #10b981;
}

.card-subtitle {
  color: var(--co-text-muted);
  font-size: 12px;
}

.charts-row {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 16px;
}

.chart-panel {
  padding: 20px;
  min-height: 300px;
}

.panel-title {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--co-border);
}

.chart-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  gap: 12px;
  color: var(--co-text-muted);
  font-size: 13px;
}

@media (max-width: 1400px) {
  .overview-cards {
    grid-template-columns: repeat(2, 1fr);
  }
  .charts-row {
    grid-template-columns: 1fr;
  }
}
</style>
