<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <RouterView />
</template>

<style>
/* 初始化暗色模式 */
html {
  color-scheme: dark;
}
html.dark {
  color-scheme: dark;
}
</style>
