<script setup lang="ts">
// 头部组件
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import { useRoute } from 'vue-router'
import { computed } from 'vue'

const appStore = useAppStore()
const userStore = useUserStore()
const route = useRoute()

// 面包屑
const breadcrumbs = computed(() => {
  const matched = route.matched.filter(item => item.meta?.title)
  return matched.map(item => ({
    title: item.meta.title as string,
    path: item.path,
  }))
})

function handleLogout() {
  userStore.logout()
}
</script>

<template>
  <header class="app-header">
    <div class="header-left">
      <!-- 折叠按钮 -->
      <div class="collapse-btn" @click="appStore.toggleSidebar">
        <el-icon :size="18">
          <component :is="appStore.sidebarCollapsed ? 'Expand' : 'Fold'" />
        </el-icon>
      </div>

      <!-- 面包屑 -->
      <el-breadcrumb separator="/">
        <el-breadcrumb-item v-for="item in breadcrumbs" :key="item.path">
          {{ item.title }}
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <div class="header-right">
      <!-- 全局搜索 -->
      <div class="header-action" title="搜索">
        <el-icon :size="18"><Search /></el-icon>
      </div>

      <!-- 全屏 -->
      <div class="header-action" title="全屏">
        <el-icon :size="18"><FullScreen /></el-icon>
      </div>

      <!-- 通知 -->
      <div class="header-action" title="通知">
        <el-badge :value="3" :max="99">
          <el-icon :size="18"><Bell /></el-icon>
        </el-badge>
      </div>

      <!-- 用户信息 -->
      <el-dropdown trigger="click">
        <div class="user-info">
          <el-avatar :size="32" class="user-avatar">
            {{ userStore.userInfo?.nickname?.charAt(0) || 'U' }}
          </el-avatar>
          <span class="user-name">{{ userStore.userInfo?.nickname || userStore.userInfo?.username || '用户' }}</span>
          <el-icon :size="12"><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item>
              <el-icon><User /></el-icon>个人中心
            </el-dropdown-item>
            <el-dropdown-item divided @click="handleLogout">
              <el-icon><SwitchButton /></el-icon>退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>

<style scoped>
.app-header {
  height: var(--co-header-height);
  background: var(--co-bg-header);
  border-bottom: 1px solid var(--co-border);
  backdrop-filter: blur(20px);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.collapse-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  cursor: pointer;
  color: var(--co-text-secondary);
  transition: var(--co-transition);
}

.collapse-btn:hover {
  background: rgba(102, 126, 234, 0.1);
  color: var(--co-primary);
}

:deep(.el-breadcrumb__inner) {
  color: var(--co-text-muted) !important;
}

:deep(.el-breadcrumb__item:last-child .el-breadcrumb__inner) {
  color: var(--co-text-primary) !important;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.header-action {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  cursor: pointer;
  color: var(--co-text-secondary);
  transition: var(--co-transition);
}

.header-action:hover {
  background: rgba(102, 126, 234, 0.1);
  color: var(--co-primary);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: var(--co-transition);
}

.user-info:hover {
  background: rgba(102, 126, 234, 0.1);
}

.user-avatar {
  background: var(--co-gradient) !important;
  color: white !important;
  font-weight: 600;
}

.user-name {
  font-size: 14px;
  color: var(--co-text-primary);
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
