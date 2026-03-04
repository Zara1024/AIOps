import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { RouteLocationNormalized } from 'vue-router'

export interface TagView {
    path: string
    name: string
    title: string
    affix?: boolean // 是否固定不可关闭
}

export const useTagsStore = defineStore('tags', () => {
    const visitedViews = ref<TagView[]>([
        { path: '/dashboard', name: 'Dashboard', title: '仪表盘', affix: true },
    ])

    const cachedViews = ref<string[]>(['Dashboard'])

    // 添加标签
    function addView(route: RouteLocationNormalized) {
        const exists = visitedViews.value.find(v => v.path === route.path)
        if (exists) return

        visitedViews.value.push({
            path: route.path,
            name: route.name as string,
            title: (route.meta?.title as string) || '未命名',
        })

        // 添加到缓存列表
        if (route.name && !cachedViews.value.includes(route.name as string)) {
            cachedViews.value.push(route.name as string)
        }
    }

    // 关闭标签
    function removeView(path: string) {
        const index = visitedViews.value.findIndex(v => v.path === path)
        if (index === -1) return

        const view = visitedViews.value[index]
        if (!view || view.affix) return // 固定标签不可关闭

        visitedViews.value.splice(index, 1)
        // 从缓存中移除
        const cachedIndex = cachedViews.value.indexOf(view.name)
        if (cachedIndex > -1) {
            cachedViews.value.splice(cachedIndex, 1)
        }
    }

    // 关闭其他标签
    function removeOtherViews(path: string) {
        visitedViews.value = visitedViews.value.filter(v => v.affix || v.path === path)
        const currentView = visitedViews.value.find(v => v.path === path)
        cachedViews.value = visitedViews.value
            .filter(v => v.affix || v.path === path)
            .map(v => v.name)
    }

    // 关闭所有标签
    function removeAllViews() {
        visitedViews.value = visitedViews.value.filter(v => v.affix)
        cachedViews.value = visitedViews.value.map(v => v.name)
    }

    return {
        visitedViews,
        cachedViews,
        addView,
        removeView,
        removeOtherViews,
        removeAllViews,
    }
}, {
    persist: {
        storage: sessionStorage,
        pick: ['visitedViews', 'cachedViews'],
    },
})
