<script setup lang="ts">
// 部门管理页面（树形表格）
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const tableData = ref([
  {
    id: 1, dept_name: '智维云枢', leader: 'CEO', phone: '13800138000', email: 'ceo@cloudops.local',
    sort_order: 1, status: 1,
    children: [
      { id: 2, dept_name: '技术部', leader: 'CTO', phone: '13900139000', email: 'cto@cloudops.local', sort_order: 1, status: 1,
        children: [
          { id: 3, dept_name: '运维组', leader: '张三', phone: '', email: '', sort_order: 1, status: 1 },
          { id: 4, dept_name: '开发组', leader: '李四', phone: '', email: '', sort_order: 2, status: 1 },
        ],
      },
      { id: 5, dept_name: '运营部', leader: '王五', phone: '', email: '', sort_order: 2, status: 1 },
    ],
  },
])

const dialogVisible = ref(false)
const dialogTitle = ref('新增部门')
const formData = ref({ dept_name: '', leader: '', phone: '', email: '', parent_id: 0, sort_order: 0, status: 1 })

function handleAdd() {
  dialogTitle.value = '新增部门'
  formData.value = { dept_name: '', leader: '', phone: '', email: '', parent_id: 0, sort_order: 0, status: 1 }
  dialogVisible.value = true
}

function handleEdit(row: any) {
  dialogTitle.value = '编辑部门'
  formData.value = { ...row }
  dialogVisible.value = true
}

async function handleDelete(row: any) {
  await ElMessageBox.confirm(`确定删除部门「${row.dept_name}」？`, '警告', { type: 'warning' })
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
        <el-button type="primary" @click="handleAdd">新增部门</el-button>
      </div>

      <el-table :data="tableData" row-key="id" :tree-props="{ children: 'children' }" default-expand-all stripe>
        <el-table-column prop="dept_name" label="部门名称" min-width="200" />
        <el-table-column prop="leader" label="负责人" width="120" />
        <el-table-column prop="phone" label="联系电话" width="140" />
        <el-table-column prop="email" label="邮箱" min-width="180" />
        <el-table-column prop="sort_order" label="排序" width="70" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '正常' : '停用' }}
            </el-tag>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500" destroy-on-close>
      <el-form :model="formData" label-width="80px">
        <el-form-item label="部门名称" required><el-input v-model="formData.dept_name" /></el-form-item>
        <el-form-item label="负责人"><el-input v-model="formData.leader" /></el-form-item>
        <el-form-item label="联系电话"><el-input v-model="formData.phone" /></el-form-item>
        <el-form-item label="邮箱"><el-input v-model="formData.email" /></el-form-item>
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
.table-card { padding: 16px 20px; }
.table-toolbar { margin-bottom: 16px; }
</style>
