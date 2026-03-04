<script setup lang="ts">
// 角色管理页面
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const tableData = ref([
  { id: 1, role_name: '超级管理员', role_key: 'super_admin', description: '拥有所有权限', sort_order: 1, status: 1, created_at: '2026-03-04' },
  { id: 2, role_name: '管理员', role_key: 'admin', description: '管理权限', sort_order: 2, status: 1, created_at: '2026-03-04' },
  { id: 3, role_name: '只读用户', role_key: 'viewer', description: '只读权限', sort_order: 3, status: 1, created_at: '2026-03-04' },
])
const total = ref(3)
const queryParams = ref({ page: 1, page_size: 20, role_name: '' })
const dialogVisible = ref(false)
const dialogTitle = ref('新增角色')
const formData = ref({ role_name: '', role_key: '', description: '', sort_order: 0, status: 1 })

function handleAdd() {
  dialogTitle.value = '新增角色'
  formData.value = { role_name: '', role_key: '', description: '', sort_order: 0, status: 1 }
  dialogVisible.value = true
}

function handleEdit(row: any) {
  dialogTitle.value = '编辑角色'
  formData.value = { ...row }
  dialogVisible.value = true
}

async function handleDelete(row: any) {
  await ElMessageBox.confirm(`确定删除角色「${row.role_name}」？`, '警告', { type: 'warning' })
  ElMessage.success('删除成功')
}

function handleSubmit() {
  ElMessage.success(dialogTitle.value === '新增角色' ? '创建成功' : '更新成功')
  dialogVisible.value = false
}
</script>

<template>
  <div class="page-container fade-in-up">
    <div class="search-bar glass-card">
      <el-form :inline="true" :model="queryParams">
        <el-form-item label="角色名称">
          <el-input v-model="queryParams.role_name" placeholder="请输入角色名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary">搜索</el-button>
          <el-button>重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="table-card glass-card">
      <div class="table-toolbar">
        <el-button type="primary" @click="handleAdd">新增角色</el-button>
      </div>

      <el-table :data="tableData" stripe>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="role_name" label="角色名称" width="140" />
        <el-table-column prop="role_key" label="角色标识" width="140" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="120" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="table-pagination">
        <el-pagination v-model:current-page="queryParams.page" :total="total" background layout="total, prev, pager, next" />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500" destroy-on-close>
      <el-form :model="formData" label-width="80px">
        <el-form-item label="角色名称" required><el-input v-model="formData.role_name" /></el-form-item>
        <el-form-item label="角色标识" required><el-input v-model="formData.role_key" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="formData.description" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="formData.sort_order" :min="0" /></el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
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
.search-bar { padding: 16px 20px 0; }
.table-card { padding: 16px 20px; }
.table-toolbar { margin-bottom: 16px; }
.table-pagination { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
