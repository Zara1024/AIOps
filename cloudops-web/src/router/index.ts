import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { useUserStore } from '@/stores/user'
import { useTagsStore } from '@/stores/tags'

NProgress.configure({ showSpinner: false })

// 静态路由
const constantRoutes: RouteRecordRaw[] = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/login/index.vue'),
        meta: { title: '登录', hidden: true },
    },
    {
        path: '/',
        name: 'Layout',
        component: () => import('@/layout/index.vue'),
        redirect: '/dashboard',
        children: [
            {
                path: 'dashboard',
                name: 'Dashboard',
                component: () => import('@/views/dashboard/index.vue'),
                meta: { title: '仪表盘', icon: 'Monitor', affix: true },
            },
            // 系统管理
            {
                path: 'system/users',
                name: 'SystemUsers',
                component: () => import('@/views/system/users/index.vue'),
                meta: { title: '用户管理', icon: 'User' },
            },
            {
                path: 'system/roles',
                name: 'SystemRoles',
                component: () => import('@/views/system/roles/index.vue'),
                meta: { title: '角色管理', icon: 'UserFilled' },
            },
            {
                path: 'system/menus',
                name: 'SystemMenus',
                component: () => import('@/views/system/menus/index.vue'),
                meta: { title: '菜单管理', icon: 'Menu' },
            },
            {
                path: 'system/departments',
                name: 'SystemDepartments',
                component: () => import('@/views/system/departments/index.vue'),
                meta: { title: '部门管理', icon: 'OfficeBuilding' },
            },
        ],
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/views/error/404.vue'),
        meta: { title: '404', hidden: true },
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes: constantRoutes,
})

// 白名单路径
const whiteList = ['/login']

// 路由守卫
router.beforeEach(async (to, from, next) => {
    NProgress.start()
    document.title = `${to.meta.title || ''} - 智维云枢`

    const userStore = useUserStore()

    if (userStore.isLoggedIn()) {
        if (to.path === '/login') {
            next({ path: '/' })
        } else {
            // 如果没有用户信息，获取一次
            if (!userStore.userInfo) {
                try {
                    await userStore.fetchUserInfo()
                } catch {
                    userStore.logout()
                    next(`/login?redirect=${to.path}`)
                    return
                }
            }
            // 添加标签
            const tagsStore = useTagsStore()
            if (to.name && !to.meta.hidden) {
                tagsStore.addView(to)
            }
            next()
        }
    } else {
        if (whiteList.includes(to.path)) {
            next()
        } else {
            next(`/login?redirect=${to.path}`)
        }
    }
})

router.afterEach(() => {
    NProgress.done()
})

export default router
