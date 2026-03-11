<template>
  <div>
    <div style="display: flex; justify-content: space-between; margin-bottom: 16px">
      <h3 style="margin: 0">{{ $t('logos.title') }}</h3>
      <el-upload :action="uploadUrl" :headers="uploadHeaders" :on-success="onUploadSuccess" :on-error="onUploadError"
                 :show-file-list="false" accept="image/*" multiple>
        <el-button type="primary">{{ $t('logos.upload') }}</el-button>
      </el-upload>
    </div>

    <el-table :data="logos" v-loading="loading" border stripe>
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column :label="$t('common.name')" min-width="180">
        <template #default="{ row }">
          <div v-if="row.isEditing" style="display: flex; gap: 8px; align-items: center">
            <el-input v-model="row.editName" size="small" @keyup.enter="saveName(row)" />
            <el-button size="small" type="success" :icon="Check" circle @click="saveName(row)" />
            <el-button size="small" type="info" :icon="Close" circle @click="cancelEdit(row)" />
          </div>
          <div v-else style="display: flex; align-items: center; justify-content: space-between;">
            <span style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap;" :title="row.name">{{ row.name }}</span>
            <el-button size="small" type="primary" link :icon="Edit" @click="startEdit(row)">{{ $t('common.edit') }}</el-button>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('logos.col_preview')" width="100" align="center">
        <template #default="{ row }">
          <el-image
            :src="row.url_path"
            style="width: 40px; height: 40px; cursor: pointer"
            fit="contain"
            :preview-src-list="[row.url_path]"
            :z-index="3000"
            preview-teleported
            hide-on-click-modal
          />
        </template>
      </el-table-column>
      <el-table-column prop="url_path" :label="$t('logos.col_url_path')" min-width="250" show-overflow-tooltip />
      <el-table-column prop="created_at" :label="$t('logos.col_upload_time')" width="180">
        <template #default="{ row }">{{ new Date(row.created_at).toLocaleString() }}</template>
      </el-table-column>
      <el-table-column :label="$t('common.operations')" width="100" fixed="right" align="center">
        <template #default="{ row }">
          <el-tooltip :content="$t('common.delete')" placement="top" :show-after="500">
            <el-button :icon="Delete" size="small" circle type="danger" @click="handleDelete(row)" />
          </el-tooltip>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Delete, Edit, Check, Close } from '@element-plus/icons-vue'
import api from '../api'

const { t } = useI18n()

const logos = ref([])
const loading = ref(false)

const uploadUrl = '/api/logos/upload'
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${localStorage.getItem('token')}`,
}))

onMounted(() => loadLogos())

async function loadLogos() {
  loading.value = true
  try {
    const { data } = await api.get('/logos')
    // 初始化每个元素的编辑状态
    logos.value = (data || []).map(item => ({
      ...item,
      isEditing: false,
      editName: ''
    }))
  } finally { loading.value = false }
}

function onUploadSuccess(response) {
  // Check if it's actually an error wrapped in a 200/success response
  if (response && response.error) {
    ElMessage.error(response.error)
  } else {
    ElMessage.success(t('logos.upload_success'))
  }
  loadLogos()
}

function onUploadError(err) {
  let msg = t('logos.upload_failed')
  if (err && err.message) {
    try {
      const parsed = JSON.parse(err.message)
      if (parsed.error) msg = parsed.error
    } catch {
      msg = err.message
    }
  }
  ElMessage.error(msg)
}

// 编辑台标名称逻辑
function startEdit(row) {
  row.editName = row.name
  row.isEditing = true
}

function cancelEdit(row) {
  row.isEditing = false
  row.editName = ''
}

async function saveName(row) {
  if (!row.editName || !row.editName.trim()) {
    ElMessage.warning(t('logos.name_empty'))
    return
  }
  
  const newName = row.editName.trim()
  
  if (newName === row.name) {
    cancelEdit(row)
    return
  }

  try {
    await api.put(`/logos/${row.id}`, { name: newName })
    row.name = newName
    ElMessage.success(t('logos.rename_success'))
    row.isEditing = false
  } catch (e) {
    // 错误在api拦截器已处理
  }
}

async function handleDelete(row) {
  await ElMessageBox.confirm(t('logos.delete_confirm', { name: row.name }), t('common.confirm_delete'), { type: 'warning', confirmButtonText: t('common.confirm'), cancelButtonText: t('common.cancel') })
  await api.delete(`/logos/${row.id}`)
  ElMessage.success(t('common.delete_success'))
  await loadLogos()
}
</script>
