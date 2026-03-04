<script setup lang="ts">
// 菜单管理页面（树形表格）
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const tableData = ref([
  {
    id: 1, menu_name: '仪表盘', menu_type: 2, path: '/dashboard', component: 'views/dashboard/index',
    icon: 'Monitor', permission: '', sort_order: 1, visible: true, status: 1,
  },
  {
    id: 2, menu_name: '系统管理', menu_type: 1, path: '/system', component: '',
    icon: 'Setting', permission: '', sort_order: 100, visible: true, status: 1,
    children: [
      { id: 3, menu_name: '用户管理', menu_type: 2, path: '/system/users', component: 'views/system/users/index', icon: 'User', permission: 'system:user:list', sort_order: 1, visible: true, status: 1 },
      { id: 4, menu_name: '角色管理', menu_type: 2, path: '/system/roles', component: 'views/system/roles/index', icon: 'UserFilled', permission: 'system:role:list', sort_order: 2, visible: true, status: 1 },
      { id: 5, menu_name: '菜单管理', menu_type: 2, path: '/system/menus', component: 'views/system/menus/index', icon: 'Menu', permission: 'system:menu:list', sort_order: 3, visible: true, status: 1 },
      { id: 6, menu_name: '部门管理', menu_type: 2, path: '/system/departments', component: 'views/system/departments/index', icon: 'OfficeBuilding', permission: 'system:dept:list', sort_order: 4, visible: true, status: 1 },
    ],
  },
])

const dialogVisible = ref(false)
const dialogTitle = ref('新增菜单')
const formData = ref({ menu_name: '', menu_type: 2, path: '', component: '', icon: '', permission: '', parent_id: 0, sort_order: 0, visible: true, status: 1 })

const menuTypeMap: Record<number, string> = { 1: '目录', 2: '菜单', 3: '按钮' }

function handleAdd() {
  dialogTitle.value = '新增菜单'
  formData.value = { menu_name: '', menu_type: 2, path: '', component: '', icon: '', permission: '', parent_id: 0, sort_order: 0, visible: true, status: 1 }
  dialogVisible.value = true
}

function handleEdit(row: any) {
  dialogTitle.value = '编辑菜单'
  formData.value = { ...row }
  dialogVisible.value = true
}

async function handleDelete(row: any) {
  await ElMessageBox.confirm(`确定删除菜单「${row.menu_name}」？`, '警告', { type: 'warning' })
  ElMessage.success('删除成功')
}

function handleSubmit() {
  ElMessage.success('保存成功')
  dialogVisible.value = false
}
</script>

<template>
  <div class="page-container fade-in-up">
    <div class="table-card glass-card">
      <div class="table-toolbar">
        <el-button type="primary" @click="handleAdd">新增菜单</el-button>
      </div>

      <el-table :data="tableData" row-key="id" :tree-props="{ children: 'children' }" default-expand-all stripe>
        <el-table-column prop="menu_name" label="菜单名称" min-width="160" />
        <el-table-column prop="icon" label="图标" width="80">
          <template #default="{ row }">
            <el-icon v-if="row.icon"><component :is="row.icon" /></el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="menu_type" label="类型" width="80">
          <template #default="{ row }">
            <el-tag :type="row.menu_type === 1 ? 'warning' : row.menu_type === 2 ? 'primary' : 'info'" size="small">
              {{ menuTypeMap[row.menu_type] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路由路径" min-width="160" />
        <el-table-column prop="permission" label="权限标识" min-width="160" />
        <el-table-column prop="sort_order" label="排序" width="70" />
        <el-table-column prop="visible" label="可见" width="70">
          <template #default="{ row }">
            <el-tag :type="row.visible ? 'success' : 'info'" size="small">{{ row.visible ? '是' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="550" destroy-on-close>
      <el-form :model="formData" label-width="80px">
        <el-form-item label="菜单类型">
          <el-radio-group v-model="formData.menu_type">
            <el-radio :value="1">目录</el-radio>
            <el-radio :value="2">菜单</el-radio>
            <el-radio :value="3">按钮</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="菜单名称" required><el-input v-model="formData.menu_name" /></el-form-item>
        <el-form-item label="路由路径" v-if="formData.menu_type !== 3"><el-input v-model="formData.path" /></el-form-item>
        <el-form-item label="组件路径" v-if="formData.menu_type === 2"><el-input v-model="formData.component" /></el-form-item>
        <el-form-item label="图标" v-if="formData.menu_type !== 3"><el-input v-model="formData.icon" /></el-form-item>
        <el-form-item label="权限标识" v-if="formData.menu_type === 3"><el-input v-model="formData.permission" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="formData.sort_order" :min="0" /></el-form-item>
        <el-form-item label="是否可见" v-if="formData.menu_type !== 3"><el-switch v-model="formData.visible" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page-container { display: flex; flex-direction: column; gap: 16px; }
.table-card { padding: 16px 20px; }
.table-toolbar { margin-bottom: 16px; }
</style>
