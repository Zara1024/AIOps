import { defineStore } from 'pinia'
import { ref } from 'vue'
import { loginApi, logoutApi, refreshTokenApi, getUserInfoApi, type LoginParams, type UserInfo } from '@/api'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
    // 状态
    const accessToken = ref('')
    const refreshToken = ref('')
    const userInfo = ref<UserInfo | null>(null)

    // 登录
    async function login(params: LoginParams) {
        const res = await loginApi(params)
        const { data } = res as any
        accessToken.value = data.access_token
        refreshToken.value = data.refresh_token
        userInfo.value = data.user
        return data
    }

    // 刷新 Token
    async function refreshAccessToken() {
        if (!refreshToken.value) throw new Error('无 Refresh Token')
        const res = await refreshTokenApi(refreshToken.value)
        const { data } = res as any
        accessToken.value = data.access_token
        refreshToken.value = data.refresh_token
        return data
    }

    // 获取用户信息
    async function fetchUserInfo() {
        const res = await getUserInfoApi()
        const { data } = res as any
        userInfo.value = data
        return data
    }

    // 登出
    async function logout() {
        try {
            await logoutApi()
        } catch { /* 忽略登出接口错误 */ }
        accessToken.value = ''
        refreshToken.value = ''
        userInfo.value = null
        router.push('/login')
    }

    // 是否已登录
    function isLoggedIn() {
        return !!accessToken.value
    }

    return {
        accessToken,
        refreshToken,
        userInfo,
        login,
        refreshAccessToken,
        fetchUserInfo,
        logout,
        isLoggedIn,
    }
}, {
    persist: {
        pick: ['accessToken', 'refreshToken', 'userInfo'],
    },
})
