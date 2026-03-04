import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
    // 侧边栏折叠状态
    const sidebarCollapsed = ref(false)
    // 暗色模式（默认开启）
    const isDark = ref(true)

    function toggleSidebar() {
        sidebarCollapsed.value = !sidebarCollapsed.value
    }

    function toggleDark() {
        isDark.value = !isDark.value
        // 切换 Element Plus 暗色模式
        if (isDark.value) {
            document.documentElement.classList.add('dark')
        } else {
            document.documentElement.classList.remove('dark')
        }
    }

    return {
        sidebarCollapsed,
        isDark,
        toggleSidebar,
        toggleDark,
    }
}, {
    persist: {
        pick: ['sidebarCollapsed', 'isDark'],
    },
})
