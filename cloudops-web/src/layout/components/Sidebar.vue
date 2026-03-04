<script setup lang="ts">
// 侧边栏组件 - 深空科技风
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'

const router = useRouter()
const route = useRoute()
const appStore = useAppStore()

// 菜单数据
const menuList = [
  {
    path: '/dashboard',
    title: '仪表盘',
    icon: 'Monitor',
  },
  {
    path: '/system',
    title: '系统管理',
    icon: 'Setting',
    children: [
      { path: '/system/users', title: '用户管理', icon: 'User' },
      { path: '/system/roles', title: '角色管理', icon: 'UserFilled' },
      { path: '/system/menus', title: '菜单管理', icon: 'Menu' },
      { path: '/system/departments', title: '部门管理', icon: 'OfficeBuilding' },
    ],
  },
]

const activePath = computed(() => route.path)

function handleMenuSelect(path: string) {
  router.push(path)
}
</script>

<template>
  <div class="sidebar" :class="{ collapsed: appStore.sidebarCollapsed }">
    <!-- Logo -->
    <div class="sidebar-logo" @click="router.push('/')">
      <div class="logo-icon">
        <el-icon :size="28"><Cpu /></el-icon>
      </div>
      <transition name="fade">
        <span v-if="!appStore.sidebarCollapsed" class="logo-text">智维云枢</span>
      </transition>
    </div>

    <!-- 菜单 -->
    <el-scrollbar class="sidebar-menu-wrap">
      <el-menu
        :default-active="activePath"
        :collapse="appStore.sidebarCollapsed"
        :collapse-transition="false"
        background-color="transparent"
        text-color="var(--co-text-secondary)"
        active-text-color="#667eea"
        @select="handleMenuSelect"
      >
        <template v-for="menu in menuList" :key="menu.path">
          <!-- 有子菜单 -->
          <el-sub-menu v-if="menu.children" :index="menu.path">
            <template #title>
              <el-icon><component :is="menu.icon" /></el-icon>
              <span>{{ menu.title }}</span>
            </template>
            <el-menu-item
              v-for="child in menu.children"
              :key="child.path"
              :index="child.path"
            >
              <el-icon><component :is="child.icon" /></el-icon>
              <span>{{ child.title }}</span>
            </el-menu-item>
          </el-sub-menu>

          <!-- 无子菜单 -->
          <el-menu-item v-else :index="menu.path">
            <el-icon><component :is="menu.icon" /></el-icon>
            <span>{{ menu.title }}</span>
          </el-menu-item>
        </template>
      </el-menu>
    </el-scrollbar>
  </div>
</template>

<style scoped>
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: var(--co-sidebar-width);
  background: var(--co-bg-sidebar);
  border-right: 1px solid var(--co-border);
  backdrop-filter: blur(20px);
  z-index: 1001;
  display: flex;
  flex-direction: column;
  transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.sidebar.collapsed {
  width: var(--co-sidebar-collapsed);
}

.sidebar-logo {
  height: var(--co-header-height);
  display: flex;
  align-items: center;
  padding: 0 16px;
  cursor: pointer;
  border-bottom: 1px solid var(--co-border);
  gap: 10px;
  flex-shrink: 0;
}

.logo-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--co-gradient);
  border-radius: 10px;
  color: white;
  flex-shrink: 0;
}

.logo-text {
  font-size: 16px;
  font-weight: 700;
  background: var(--co-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  white-space: nowrap;
}

.sidebar-menu-wrap {
  flex: 1;
  overflow: hidden;
}

/* 菜单项样式覆盖 */
:deep(.el-menu) {
  border: none;
  padding: 8px;
}

:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  border-radius: 8px;
  margin-bottom: 2px;
  height: 44px;
  line-height: 44px;
}

:deep(.el-menu-item:hover),
:deep(.el-sub-menu__title:hover) {
  background: rgba(102, 126, 234, 0.1) !important;
}

:deep(.el-menu-item.is-active) {
  background: rgba(102, 126, 234, 0.15) !important;
  color: #667eea !important;
  font-weight: 600;
}

:deep(.el-menu-item.is-active::before) {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 20px;
  background: var(--co-gradient);
  border-radius: 0 3px 3px 0;
}

/* 折叠状态下隐藏文字 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
