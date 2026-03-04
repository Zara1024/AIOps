<script setup lang="ts">
// CMDB 主机管理页面
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const loading = ref(false)
const queryParams = ref({ page: 1, page_size: 20, ip: '', hostname: '', status: 0, group_id: 0 })
const total = ref(5)

// 模拟数据
const tableData = ref([
  { id: 1, hostname: 'web-server-01', ip: '192.168.1.10', port: 22, os_type: 'linux', os_version: 'CentOS 7.9', cpu: '4核', memory: '8GB', status: 1, agent_status: 1, group: { group_name: '生产环境' }, created_at: '2026-03-04' },
  { id: 2, hostname: 'web-server-02', ip: '192.168.1.11', port: 22, os_type: 'linux', os_version: 'Ubuntu 22.04', cpu: '8核', memory: '16GB', status: 1, agent_status: 1, group: { group_name: '生产环境' }, created_at: '2026-03-04' },
  { id: 3, hostname: 'db-master', ip: '192.168.1.20', port: 22, os_type: 'linux', os_version: 'CentOS 8', cpu: '16核', memory: '64GB', status: 1, agent_status: 0, group: { group_name: '数据库' }, created_at: '2026-03-04' },
  { id: 4, hostname: 'k8s-node-01', ip: '10.0.1.10', port: 22, os_type: 'linux', os_version: 'Rocky 9', cpu: '8核', memory: '32GB', status: 2, agent_status: 2, group: null, created_at: '2026-03-04' },
  { id: 5, hostname: 'win-ad-01', ip: '10.0.2.10', port: 3389, os_type: 'windows', os_version: 'Windows Server 2022', cpu: '4核', memory: '16GB', status: 1, agent_status: 0, group: { group_name: '办公网络' }, created_at: '2026-03-04' },
])

const dialogVisible = ref(false)
const dialogTitle = ref('新增主机')
const formData = ref({
  hostname: '', ip: '', port: 22, os_type: 'linux', auth_type: 'password',
  username: 'root', password: '', private_key: '', group_id: null as number | null,
  description: '', labels: '',
})

const statusMap: Record<number, { text: string; type: string }> = {
  0: { text: '禁用', type: 'info' },
  1: { text: '在线', type: 'success' },
  2: { text: '离线', type: 'danger' },
  3: { text: '未知', type: 'warning' },
}

const agentStatusMap: Record<number, { text: string; type: string }> = {
  0: { text: '未安装', type: 'info' },
  1: { text: '正常', type: 'success' },
  2: { text: '异常', type: 'danger' },
}

const multipleSelection = ref<any[]>([])

function handleSelectionChange(val: any[]) {
  multipleSelection.value = val
}

function handleAdd() {
  dialogTitle.value = '新增主机'
  formData.value = { hostname: '', ip: '', port: 22, os_type: 'linux', auth_type: 'password', username: 'root', password: '', private_key: '', group_id: null, description: '', labels: '' }
  dialogVisible.value = true
}

function handleEdit(row: any) {
  dialogTitle.value = '编辑主机'
  formData.value = { ...row, password: '', private_key: '' }
  dialogVisible.value = true
}

async function handleDelete(row: any) {
  await ElMessageBox.confirm(`确定删除主机「${row.hostname}」(${row.ip})？`, '警告', { type: 'warning' })
  ElMessage.success('删除成功')
}

async function handleBatchDelete() {
  if (multipleSelection.value.length === 0) {
    ElMessage.warning('请选择要删除的主机')
    return
  }
  await ElMessageBox.confirm(`确定删除选中的 ${multipleSelection.value.length} 台主机？`, '批量删除', { type: 'warning' })
  ElMessage.success('批量删除成功')
}

function handleTerminal(row: any) {
  ElMessage.info(`SSH终端连接 ${row.ip} 功能将在后续版本实现`)
}

function handleSubmit() {
  ElMessage.success(dialogTitle.value === '新增主机' ? '创建成功' : '更新成功')
  dialogVisible.value = false
}
</script>

<template>
  <div class="page-container fade-in-up">
    <!-- 搜索栏 -->
    <div class="search-bar glass-card">
      <el-form :inline="true" :model="queryParams">
        <el-form-item label="IP">
          <el-input v-model="queryParams.ip" placeholder="IP地址" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="主机名">
          <el-input v-model="queryParams.hostname" placeholder="主机名" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="queryParams.status" placeholder="全部" clearable style="width: 100px">
            <el-option label="在线" :value="1" />
            <el-option label="离线" :value="2" />
            <el-option label="未知" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary">搜索</el-button>
          <el-button>重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 工具栏 + 表格 -->
    <div class="table-card glass-card">
      <div class="table-toolbar">
        <div class="toolbar-left">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>新增主机
          </el-button>
          <el-button type="danger" plain :disabled="multipleSelection.length === 0" @click="handleBatchDelete">
            <el-icon><Delete /></el-icon>批量删除
          </el-button>
        </div>
        <div class="toolbar-right">
          <span class="host-stat">
            共 <strong>{{ total }}</strong> 台 ·
            <span style="color: #10b981">在线 3</span> /
            <span style="color: #ef4444">离线 1</span> /
            <span style="color: #f59e0b">未知 1</span>
          </span>
        </div>
      </div>

      <el-table :data="tableData" :loading="loading" stripe @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="ip" label="IP 地址" width="140">
          <template #default="{ row }">
            <div class="ip-cell">
              <el-icon v-if="row.os_type === 'linux'" style="color: #f59e0b"><Platform /></el-icon>
              <el-icon v-else style="color: #3b82f6"><Monitor /></el-icon>
              <span>{{ row.ip }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="hostname" label="主机名" min-width="140" />
        <el-table-column prop="os_version" label="操作系统" min-width="150" />
        <el-table-column label="配置" width="120">
          <template #default="{ row }">
            <span>{{ row.cpu }} / {{ row.memory }}</span>
          </template>
        </el-table-column>
        <el-table-column label="分组" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.group" size="small">{{ row.group.group_name }}</el-tag>
            <span v-else style="color: var(--co-text-muted)">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="(statusMap[row.status]?.type || 'info') as any" size="small" effect="dark">
              {{ statusMap[row.status]?.text || '未知' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="agent_status" label="Agent" width="80">
          <template #default="{ row }">
            <el-tag :type="(agentStatusMap[row.agent_status]?.type || 'info') as any" size="small">
              {{ agentStatusMap[row.agent_status]?.text || '未知' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleTerminal(row)">
              <el-icon><Connection /></el-icon>终端
            </el-button>
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600" destroy-on-close>
      <el-form :model="formData" label-width="80px">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="主机名" required>
              <el-input v-model="formData.hostname" placeholder="如 web-server-01" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="IP 地址" required>
              <el-input v-model="formData.ip" placeholder="如 192.168.1.10" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="端口">
              <el-input-number v-model="formData.port" :min="1" :max="65535" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="系统">
              <el-select v-model="formData.os_type" style="width: 100%">
                <el-option label="Linux" value="linux" />
                <el-option label="Windows" value="windows" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="认证">
              <el-select v-model="formData.auth_type" style="width: 100%">
                <el-option label="密码" value="password" />
                <el-option label="密钥" value="key" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="用户名">
              <el-input v-model="formData.username" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item v-if="formData.auth_type === 'password'" label="密码">
              <el-input v-model="formData.password" type="password" show-password />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item v-if="formData.auth_type === 'key'" label="私钥">
          <el-input v-model="formData.private_key" type="textarea" :rows="4" placeholder="粘贴 SSH 私钥内容" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="formData.description" type="textarea" :rows="2" />
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
.table-toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.toolbar-left { display: flex; gap: 8px; }
.host-stat { color: var(--co-text-secondary); font-size: 13px; }
.host-stat strong { color: var(--co-primary); }
.ip-cell { display: flex; align-items: center; gap: 6px; font-family: 'Courier New', monospace; }
.table-pagination { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
