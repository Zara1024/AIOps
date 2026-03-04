<script setup lang="ts">
// 标签页导航
import { useRouter, useRoute } from 'vue-router'
import { useTagsStore } from '@/stores/tags'

const router = useRouter()
const route = useRoute()
const tagsStore = useTagsStore()

function handleClick(path: string) {
  router.push(path)
}

function handleClose(path: string) {
  tagsStore.removeView(path)
  // 如果关闭的是当前标签，跳转到最后一个标签
  if (path === route.path) {
    const lastView = tagsStore.visitedViews[tagsStore.visitedViews.length - 1]
    router.push(lastView?.path || '/dashboard')
  }
}

function handleContextMenu(action: string) {
  switch (action) {
    case 'closeOthers':
      tagsStore.removeOtherViews(route.path)
      break
    case 'closeAll':
      tagsStore.removeAllViews()
      router.push('/dashboard')
      break
  }
}
</script>

<template>
  <div class="tags-view">
    <el-scrollbar class="tags-scroll">
      <div class="tags-wrap">
        <div
          v-for="tag in tagsStore.visitedViews"
          :key="tag.path"
          class="tag-item"
          :class="{ active: route.path === tag.path }"
          @click="handleClick(tag.path)"
        >
          <span class="tag-dot" v-if="route.path === tag.path"></span>
          <span class="tag-title">{{ tag.title }}</span>
          <el-icon
            v-if="!tag.affix"
            class="tag-close"
            :size="12"
            @click.stop="handleClose(tag.path)"
          >
            <Close />
          </el-icon>
        </div>
      </div>
    </el-scrollbar>

    <!-- 操作按钮 -->
    <el-dropdown trigger="click" @command="handleContextMenu">
      <div class="tags-action">
        <el-icon :size="14"><ArrowDown /></el-icon>
      </div>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="closeOthers">关闭其他</el-dropdown-item>
          <el-dropdown-item command="closeAll">关闭全部</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<style scoped>
.tags-view {
  height: var(--co-tabs-height);
  background: var(--co-bg-dark);
  border-bottom: 1px solid var(--co-border);
  display: flex;
  align-items: center;
  padding: 0 8px;
  flex-shrink: 0;
}

.tags-scroll {
  flex: 1;
  overflow: hidden;
}

.tags-wrap {
  display: flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
  padding: 4px 0;
}

.tag-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 12px;
  color: var(--co-text-secondary);
  cursor: pointer;
  transition: var(--co-transition);
  border: 1px solid transparent;
}

.tag-item:hover {
  background: rgba(102, 126, 234, 0.08);
  color: var(--co-text-primary);
}

.tag-item.active {
  background: rgba(102, 126, 234, 0.12);
  color: var(--co-primary);
  border-color: rgba(102, 126, 234, 0.2);
}

.tag-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--co-primary);
}

.tag-close {
  border-radius: 50%;
  padding: 1px;
  transition: var(--co-transition);
}

.tag-close:hover {
  background: rgba(102, 126, 234, 0.3);
  color: white;
}

.tags-action {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  cursor: pointer;
  color: var(--co-text-secondary);
  transition: var(--co-transition);
  margin-left: 4px;
}

.tags-action:hover {
  background: rgba(102, 126, 234, 0.1);
  color: var(--co-primary);
}
</style>
