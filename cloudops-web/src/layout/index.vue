<script setup lang="ts">
// 布局主框架：侧边栏 + 头部 + 标签页 + 内容区
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'
import { useTagsStore } from '@/stores/tags'
import AppSidebar from './components/Sidebar.vue'
import AppHeader from './components/Header.vue'
import TagsView from './components/TagsView.vue'

const appStore = useAppStore()
const tagsStore = useTagsStore()

const mainStyle = computed(() => ({
  marginLeft: appStore.sidebarCollapsed ? 'var(--co-sidebar-collapsed)' : 'var(--co-sidebar-width)',
  transition: 'margin-left 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
}))
</script>

<template>
  <div class="app-layout">
    <!-- 侧边栏 -->
    <AppSidebar />

    <!-- 主内容区 -->
    <div class="main-container" :style="mainStyle">
      <!-- 头部 -->
      <AppHeader />

      <!-- 标签页 -->
      <TagsView />

      <!-- 内容区域 -->
      <main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade-slide" mode="out-in">
            <keep-alive :include="tagsStore.cachedViews">
              <component :is="Component" />
            </keep-alive>
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<style scoped>
.app-layout {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
  background: var(--co-bg-deep);
}

.main-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background: var(--co-bg-deep);
}

/* 路由切换动画 */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.25s ease-out;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}
</style>
