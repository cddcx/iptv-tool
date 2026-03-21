<template>
  <div class="log-container">
    <div class="log-header">
      <h2>{{ $t('logs.runtime_title') }}</h2>
      <div class="log-actions">
        <el-select v-model="levelFilter" size="small" style="width: 120px" :placeholder="$t('logs.filter_all')">
          <el-option :label="$t('logs.filter_all')" value="" />
          <el-option label="DEBUG" value="DEBUG" />
          <el-option label="INFO" value="INFO" />
          <el-option label="WARN" value="WARN" />
          <el-option label="ERROR" value="ERROR" />
        </el-select>
        <el-button
            :type="isPaused ? 'success' : 'warning'"
            @click="isPaused = !isPaused"
            :icon="isPaused ? VideoPlay : VideoPause"
            size="small"
        >
          {{ isPaused ? $t('logs.resume') : $t('logs.pause') }}
        </el-button>
        <el-button type="danger" :icon="Delete" size="small" plain @click="confirmClear">
          {{ $t('logs.clear') }}
        </el-button>
        <el-button type="primary" :icon="Download" size="small" plain @click="downloadLogs">
          {{ $t('logs.download') }}
        </el-button>
      </div>
    </div>
    <div class="log-area">
      <div v-if="filteredEntries.length === 0" class="log-empty">{{ $t('logs.no_logs') }}</div>
      <div v-for="entry in filteredEntries" :key="entry.id" class="log-line">
        <span class="log-time">{{ entry.time }}</span>
        <span :class="['log-level', `level-${entry.level.toLowerCase()}`]">{{ entry.level }}</span>
        <span class="log-content">{{ entry.content }}</span>
      </div>
    </div>
    <div class="log-status-bar">
      <span>{{ $t('logs.total_lines', { count: filteredEntries.length }) }}</span>
      <span :class="['status-dot', isPaused ? 'paused' : 'live']"></span>
      <span>{{ isPaused ? $t('logs.status_paused') : $t('logs.status_live') }}</span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { Delete, Download, VideoPlay, VideoPause } from '@element-plus/icons-vue'
import api from '../api'

const { t } = useI18n()
const entries = ref([])
const isPaused = ref(false)
const lastID = ref(0)
const levelFilter = ref('')
let timer = null

const filteredEntries = computed(() => {
  if (!levelFilter.value) return entries.value
  return entries.value.filter(e => e.level === levelFilter.value)
})

async function fetchLogs() {
  if (isPaused.value) return
  try {
    const { data } = await api.get(`/logs/runtime?since=${lastID.value}`)
    if (data.entries && data.entries.length > 0) {
      let maxID = lastID.value
      for (const e of data.entries) {
        if (e.id > maxID) maxID = e.id
      }
      lastID.value = maxID
      entries.value.unshift(...data.entries)
      if (entries.value.length > 10000) {
        entries.value = entries.value.slice(0, 5000)
      }
    }
  } catch {
    // Silently ignore polling errors
  }
}

function confirmClear() {
  ElMessageBox.confirm(t('logs.clear_confirm'), t('logs.clear'), {
    confirmButtonText: t('common.confirm'),
    cancelButtonText: t('common.cancel'),
    type: 'warning',
  }).then(() => clearLogs()).catch(() => {})
}

async function clearLogs() {
  try {
    await api.delete('/logs/runtime')
    entries.value = []
    lastID.value = 0
    ElMessage.success(t('logs.clear_success'))
  } catch {
    // error handled by interceptor
  }
}

async function downloadLogs() {
  try {
    const response = await api.get('/logs/runtime/download', { responseType: 'blob' })
    const url = URL.createObjectURL(response.data)
    const a = document.createElement('a')
    a.href = url
    a.download = response.headers['content-disposition']?.split('filename=')[1] || 'runtime.log'
    a.click()
    URL.revokeObjectURL(url)
  } catch {
    // error handled by interceptor
  }
}

onMounted(() => {
  fetchLogs()
  timer = setInterval(fetchLogs, 2000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.log-container {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 108px);
  background: #1e1e1e;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.2);
}
.log-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  background: #252526;
  border-bottom: 1px solid #333;
  flex-shrink: 0;
}
.log-header h2 {
  color: #e0e0e0;
  font-size: 16px;
  font-weight: 600;
  margin: 0;
}
.log-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}
.log-area {
  flex: 1;
  overflow-y: auto;
  padding: 12px 16px;
  font-family: 'Cascadia Code', 'Fira Code', 'JetBrains Mono', 'Consolas', monospace;
  font-size: 12px;
  line-height: 1.6;
}
.log-area::-webkit-scrollbar {
  width: 8px;
}
.log-area::-webkit-scrollbar-track {
  background: #1e1e1e;
}
.log-area::-webkit-scrollbar-thumb {
  background: #555;
  border-radius: 4px;
}
.log-area::-webkit-scrollbar-thumb:hover {
  background: #777;
}
.log-empty {
  color: #666;
  text-align: center;
  padding: 40px;
  font-size: 14px;
}
.log-line {
  white-space: pre-wrap;
  word-break: break-all;
  padding: 1px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.03);
}
.log-time {
  color: #6a9955;
  margin-right: 8px;
}
.log-level {
  display: inline-block;
  min-width: 44px;
  text-align: center;
  font-weight: 600;
  font-size: 11px;
  margin-right: 8px;
  padding: 0 4px;
  border-radius: 3px;
}
.level-debug { color: #858585; }
.level-info { color: #4ec9b0; }
.level-warn { color: #ce9178; background: rgba(206,145,120,0.12); }
.level-error { color: #f44747; background: rgba(244,71,71,0.12); }
.log-content {
  color: #d4d4d4;
}
.log-status-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 20px;
  background: #007acc;
  color: #fff;
  font-size: 12px;
  flex-shrink: 0;
}
.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}
.status-dot.live {
  background: #4ec9b0;
  animation: pulse 1.5s infinite;
}
.status-dot.paused {
  background: #ce9178;
}
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}
</style>
