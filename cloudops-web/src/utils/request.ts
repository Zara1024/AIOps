import axios, { type AxiosInstance, type AxiosRequestConfig, type InternalAxiosRequestConfig, type AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import router from '@/router'

// 创建 Axios 实例
const service: AxiosInstance = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
    timeout: 15000,
    headers: {
        'Content-Type': 'application/json',
    },
})

// 请求拦截器
service.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
        const userStore = useUserStore()
        if (userStore.accessToken) {
            config.headers.Authorization = `Bearer ${userStore.accessToken}`
        }
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// 响应拦截器
service.interceptors.response.use(
    (response: AxiosResponse) => {
        const { code, message, data } = response.data

        // 业务成功
        if (code === 200) {
            return response.data
        }

        // 业务错误
        ElMessage.error(message || '请求失败')
        return Promise.reject(new Error(message || '请求失败'))
    },
    async (error) => {
        const { response } = error

        if (response) {
            switch (response.status) {
                case 401:
                    // Token 过期，尝试刷新
                    const userStore = useUserStore()
                    if (userStore.refreshToken) {
                        try {
                            await userStore.refreshAccessToken()
                            // 重试原请求
                            return service(error.config)
                        } catch {
                            userStore.logout()
                            router.push('/login')
                            ElMessage.error('登录已过期，请重新登录')
                        }
                    } else {
                        userStore.logout()
                        router.push('/login')
                        ElMessage.error('请先登录')
                    }
                    break
                case 403:
                    ElMessage.error('无权限访问')
                    break
                case 404:
                    ElMessage.error('请求资源不存在')
                    break
                case 429:
                    ElMessage.error('请求过于频繁，请稍后重试')
                    break
                case 500:
                    ElMessage.error('服务器内部错误')
                    break
                default:
                    ElMessage.error(response.data?.message || '请求失败')
            }
        } else {
            ElMessage.error('网络连接失败，请检查网络')
        }

        return Promise.reject(error)
    }
)

// 封装请求方法
export const request = {
    get<T = any>(url: string, params?: any, config?: AxiosRequestConfig): Promise<T> {
        return service.get(url, { params, ...config }) as any
    },
    post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
        return service.post(url, data, config) as any
    },
    put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
        return service.put(url, data, config) as any
    },
    delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
        return service.delete(url, config) as any
    },
}

export default service
