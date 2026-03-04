<script setup lang="ts">
// 登录页面 - 深空科技风
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loading = ref(false)
const loginForm = reactive({
  username: '',
  password: '',
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

const formRef = ref()

async function handleLogin() {
  try {
    await formRef.value?.validate()
  } catch { return }

  loading.value = true
  try {
    await userStore.login(loginForm)
    ElMessage.success('登录成功')
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (err: any) {
    ElMessage.error(err?.message || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <!-- 动态背景 -->
    <div class="bg-animation">
      <div class="stars"></div>
      <div class="stars2"></div>
      <div class="stars3"></div>
    </div>

    <!-- 登录卡片 -->
    <div class="login-card fade-in-up">
      <!-- Logo -->
      <div class="login-header">
        <div class="login-logo">
          <el-icon :size="36"><Cpu /></el-icon>
        </div>
        <h1 class="login-title">智维云枢</h1>
        <p class="login-subtitle">CloudOps Hub · AI 驱动智能运维平台</p>
      </div>

      <!-- 表单 -->
      <el-form
        ref="formRef"
        :model="loginForm"
        :rules="rules"
        class="login-form"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="用户名"
            size="large"
            :prefix-icon="User"
            autofocus
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
            size="large"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="login-btn"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登 录' }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="login-footer">
        <span>默认账号: admin / Admin@2026</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { User, Lock, Cpu } from '@element-plus/icons-vue'
export default { components: { User, Lock, Cpu } }
</script>

<style scoped>
.login-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--co-bg-deep);
  position: relative;
  overflow: hidden;
}

/* 星空动画背景 */
.bg-animation {
  position: absolute;
  inset: 0;
  overflow: hidden;
}

.stars, .stars2, .stars3 {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: transparent;
}

.stars {
  background-image:
    radial-gradient(1px 1px at 20% 30%, rgba(255,255,255,0.5), transparent),
    radial-gradient(1px 1px at 40% 70%, rgba(255,255,255,0.3), transparent),
    radial-gradient(1px 1px at 60% 40%, rgba(255,255,255,0.4), transparent),
    radial-gradient(2px 2px at 80% 10%, rgba(102,126,234,0.5), transparent),
    radial-gradient(1px 1px at 10% 60%, rgba(255,255,255,0.3), transparent),
    radial-gradient(1px 1px at 70% 80%, rgba(118,75,162,0.4), transparent),
    radial-gradient(1px 1px at 50% 20%, rgba(255,255,255,0.4), transparent),
    radial-gradient(1px 1px at 30% 90%, rgba(102,126,234,0.3), transparent);
  animation: float 60s linear infinite;
}

.stars2 {
  background-image:
    radial-gradient(1px 1px at 25% 35%, rgba(255,255,255,0.4), transparent),
    radial-gradient(2px 2px at 75% 55%, rgba(102,126,234,0.4), transparent),
    radial-gradient(1px 1px at 15% 75%, rgba(255,255,255,0.3), transparent),
    radial-gradient(1px 1px at 85% 25%, rgba(118,75,162,0.3), transparent),
    radial-gradient(1px 1px at 45% 85%, rgba(255,255,255,0.2), transparent);
  animation: float 90s linear infinite;
}

.stars3 {
  background-image:
    radial-gradient(2px 2px at 35% 45%, rgba(102,126,234,0.3), transparent),
    radial-gradient(1px 1px at 65% 15%, rgba(255,255,255,0.3), transparent),
    radial-gradient(1px 1px at 55% 65%, rgba(118,75,162,0.3), transparent);
  animation: float 120s linear infinite;
}

@keyframes float {
  from { transform: translateY(0); }
  to { transform: translateY(-100vh); }
}

/* 登录卡片 */
.login-card {
  width: 420px;
  padding: 48px 40px;
  background: rgba(17, 22, 49, 0.85);
  backdrop-filter: blur(30px);
  border: 1px solid var(--co-border);
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5), var(--co-glow);
  position: relative;
  z-index: 1;
}

.login-header {
  text-align: center;
  margin-bottom: 36px;
}

.login-logo {
  width: 64px;
  height: 64px;
  margin: 0 auto 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--co-gradient);
  border-radius: 18px;
  color: white;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.3);
  animation: pulse-glow 3s ease-in-out infinite;
}

.login-title {
  font-size: 28px;
  font-weight: 700;
  background: var(--co-gradient);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 8px;
}

.login-subtitle {
  color: var(--co-text-muted);
  font-size: 14px;
}

.login-form {
  margin-bottom: 0;
}

:deep(.el-input__wrapper) {
  padding: 8px 12px;
  border-radius: 10px;
}

.login-btn {
  width: 100%;
  height: 44px;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 2px;
}

.login-footer {
  text-align: center;
  margin-top: 16px;
  color: var(--co-text-muted);
  font-size: 12px;
}
</style>
