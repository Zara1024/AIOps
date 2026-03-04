<script setup lang="ts">
// 用户管理页面
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// 模拟数据（后续对接后端 API）
const loading = ref(false)
const tableData = ref([
  { id: 1, username: 'admin', nickname: '超级管理员', email: 'admin@cloudops.local', phone: '13800138000', status: 1, roles: '超级管理员', created_at: '2026-03-04' },
  { id: 2, username: 'zhangsan', nickname: '张三', email: 'zhangsan@example.com', phone: '13900139000', status: 1, roles: '管理员', created_at: '2026-03-04' },
  { id: 3, username: 'lisi', nickname: '李四', email: 'lisi@example.com', phone: '13700137000', status: 0, roles: '只读用户', created_at: '2026-03-04' },
])
const total = ref(3)
const queryParams = ref({ page: 1, page_size: 20, username: '' })

// 对话框控制
const dialogVisible = ref(false)
const dialogTitle = ref('新增用户')
const formData = ref({ username: '', password: '', nickname: '', email: '', phone: '', status: 1 })

function handleAdd() {
  dialogTitle.value = '新增用户'
  formData.value = { username: '', password: '', nickname: '', email: '', phone: '', status: 1 }
  dialogVisible.value = true
}

function handleEdit(row: any) {
  dialogTitle.value = '编辑用户'
  formData.value = { ...row, password: '' }
  dialogVisible.value = true
}

async function handleDelete(row: any) {
  await ElMessageBox.confirm(`确定删除用户「${row.nickname || row.username}」？`, '警告', { type: 'warning' })
  ElMessage.success('删除成功')
}

function handleSubmit() {
  ElMessage.success(dialogTitle.value === '新增用户' ? '创建成功' : '更新成功')
  dialogVisible.value = false
}
</script>

<template>
  <div class="page-container fade-in-up">
    <!-- 搜索栏 -->
    <div class="search-bar glass-card">
      <el-form :inline="true" :model="queryParams">
        <el-form-item label="用户名">
          <el-input v-model="queryParams.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search">搜索</el-button>
          <el-button :icon="Refresh">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 工具栏 + 表格 -->
    <div class="table-card glass-card">
      <div class="table-toolbar">
        <el-button type="primary" :icon="Plus" @click="handleAdd">新增用户</el-button>
      </div>

      <el-table :data="tableData" :loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="nickname" label="昵称" width="120" />
        <el-table-column prop="email" label="邮箱" min-width="180" />
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column prop="roles" label="角色" width="120" />
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
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.page_size"
          :total="total"
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          background
        />
      </div>
    </div>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500" destroy-on-close>
      <el-form :model="formData" label-width="80px">
        <el-form-item label="用户名" required>
          <el-input v-model="formData.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" :required="dialogTitle === '新增用户'">
          <el-input v-model="formData.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="formData.nickname" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="formData.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" active-text="正常" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { Search, Refresh, Plus } from '@element-plus/icons-vue'
export default { components: { Search, Refresh, Plus } }
</script>

<style scoped>
.page-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.search-bar {
  padding: 16px 20px 0;
}

.table-card {
  padding: 16px 20px;
}

.table-toolbar {
  margin-bottom: 16px;
}

.table-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
