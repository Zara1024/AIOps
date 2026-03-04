<script setup lang="ts">
// CMDB 主机分组管理页面（树形表格）
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const tableData = ref([
  {
    id: 1, group_name: '生产环境', sort_order: 1, host_count: 12,
    children: [
      { id: 2, group_name: 'Web 服务器', sort_order: 1, host_count: 5 },
      { id: 3, group_name: '数据库', sort_order: 2, host_count: 3 },
      { id: 4, group_name: '缓存集群', sort_order: 3, host_count: 4 },
    ],
  },
  {
    id: 5, group_name: '测试环境', sort_order: 2, host_count: 8,
    children: [
      { id: 6, group_name: '测试服务器', sort_order: 1, host_count: 5 },
      { id: 7, group_name: '开发服务器', sort_order: 2, host_count: 3 },
    ],
  },
  { id: 8, group_name: '办公网络', sort_order: 3, host_count: 3 },
])

const dialogVisible = ref(false)
const dialogTitle = ref('新增分组')
const formData = ref({ group_name: '', parent_id: 0, sort_order: 0 })

function handleAdd() {
  dialogTitle.value = '新增分组'
  formData.value = { group_name: '', parent_id: 0, sort_order: 0 }
  dialogVisible.value = true
}

function handleEdit(row: any) {
  dialogTitle.value = '编辑分组'
  formData.value = { ...row }
  dialogVisible.value = true
}

async function handleDelete(row: any) {
  if (row.host_count > 0) {
    ElMessage.warning(`分组下还有 ${row.host_count} 台主机，不能删除`)
    return
  }
  await ElMessageBox.confirm(`确定删除分组「${row.group_name}」？`, '警告', { type: 'warning' })
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
        <el-button type="primary" @click="handleAdd">新增分组</el-button>
      </div>

      <el-table :data="tableData" row-key="id" :tree-props="{ children: 'children' }" default-expand-all stripe>
        <el-table-column prop="group_name" label="分组名称" min-width="250" />
        <el-table-column prop="host_count" label="主机数量" width="120">
          <template #default="{ row }">
            <el-tag size="small" type="primary">{{ row.host_count }} 台</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button link type="danger" size="small" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="450" destroy-on-close>
      <el-form :model="formData" label-width="80px">
        <el-form-item label="分组名称" required><el-input v-model="formData.group_name" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="formData.sort_order" :min="0" /></el-form-item>
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
