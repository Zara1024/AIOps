import { request } from '@/utils/request'

// ============= 认证接口 =============

export interface LoginParams {
    username: string
    password: string
}

export interface LoginResult {
    access_token: string
    refresh_token: string
    expires_at: number
    user: UserInfo
}

export interface UserInfo {
    id: number
    username: string
    nickname: string
    email: string
    phone: string
    avatar: string
    department_id: number | null
    role_keys: string[]
    roles: string[]
}

// 登录
export const loginApi = (data: LoginParams) =>
    request.post<{ code: number; data: LoginResult }>('/auth/login', data)

// 登出
export const logoutApi = () =>
    request.post('/auth/logout')

// 刷新 Token
export const refreshTokenApi = (refresh_token: string) =>
    request.post<{ code: number; data: { access_token: string; refresh_token: string; expires_at: number } }>(
        '/auth/refresh', { refresh_token }
    )

// 获取用户信息
export const getUserInfoApi = () =>
    request.get<{ code: number; data: UserInfo }>('/auth/userinfo')

// 修改密码
export const changePasswordApi = (data: { old_password: string; new_password: string }) =>
    request.put('/auth/password', data)

// ============= 系统管理接口 =============

// 获取用户菜单（动态路由）
export const getUserMenusApi = () =>
    request.get<{ code: number; data: MenuItem[] }>('/system/menus/user')

export interface MenuItem {
    id: number
    parent_id: number
    menu_name: string
    menu_type: number  // 1目录 2菜单 3按钮
    path: string
    component: string
    icon: string
    permission: string
    sort_order: number
    visible: boolean
    children?: MenuItem[]
}

// 用户管理
export const getUserListApi = (params: any) =>
    request.get('/system/users', params)

export const createUserApi = (data: any) =>
    request.post('/system/users', data)

export const updateUserApi = (id: number, data: any) =>
    request.put(`/system/users/${id}`, data)

export const deleteUserApi = (id: number) =>
    request.delete(`/system/users/${id}`)

// 角色管理
export const getRoleListApi = (params: any) =>
    request.get('/system/roles', params)

export const createRoleApi = (data: any) =>
    request.post('/system/roles', data)

export const updateRoleApi = (id: number, data: any) =>
    request.put(`/system/roles/${id}`, data)

export const deleteRoleApi = (id: number) =>
    request.delete(`/system/roles/${id}`)

// 菜单管理
export const getMenuTreeApi = () =>
    request.get('/system/menus')

export const createMenuApi = (data: any) =>
    request.post('/system/menus', data)

export const updateMenuApi = (id: number, data: any) =>
    request.put(`/system/menus/${id}`, data)

export const deleteMenuApi = (id: number) =>
    request.delete(`/system/menus/${id}`)

// 部门管理
export const getDeptTreeApi = () =>
    request.get('/system/departments')

export const createDeptApi = (data: any) =>
    request.post('/system/departments', data)

export const updateDeptApi = (id: number, data: any) =>
    request.put(`/system/departments/${id}`, data)

export const deleteDeptApi = (id: number) =>
    request.delete(`/system/departments/${id}`)
