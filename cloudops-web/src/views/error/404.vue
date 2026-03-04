<template>
  <div class="not-found-page">
    <div class="error-content fade-in-up">
      <div class="error-code gradient-text">404</div>
      <h2 class="error-title">页面不存在</h2>
      <p class="error-desc">你访问的页面走丢了，请检查地址或返回首页</p>
      <el-button type="primary" size="large" @click="$router.push('/')">
        返回首页
      </el-button>
    </div>
  </div>
</template>

<style scoped>
.not-found-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--co-bg-deep);
}

.error-content {
  text-align: center;
}

.error-code {
  font-size: 120px;
  font-weight: 900;
  line-height: 1;
  margin-bottom: 16px;
}

.error-title {
  font-size: 24px;
  color: var(--co-text-primary);
  margin-bottom: 12px;
}

.error-desc {
  color: var(--co-text-muted);
  margin-bottom: 32px;
  font-size: 14px;
}
</style>
