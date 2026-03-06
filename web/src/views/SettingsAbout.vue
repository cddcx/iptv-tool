<template>
  <div class="settings-page">
    <h3 style="margin: 0 0 20px">关于系统</h3>

    <el-card shadow="hover" class="settings-card">
      <template #header>
        <div style="display: flex; align-items: center; gap: 8px">
          <el-icon :size="18"><InfoFilled /></el-icon>
          <span>关于系统</span>
        </div>
      </template>
      <el-descriptions :column="1" border size="small">
        <el-descriptions-item label="系统名称">IPTV Tool</el-descriptions-item>
        <el-descriptions-item label="技术栈">Vue 3 + Element Plus / Go + Gin + SQLite</el-descriptions-item>
        <el-descriptions-item label="运行模式">单文件部署 (go:embed)</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup>
import { InfoFilled } from '@element-plus/icons-vue'
</script>

<style scoped>
.settings-page {
  width: 100%;
}
.settings-card {
  max-width: 600px;
}
</style>
