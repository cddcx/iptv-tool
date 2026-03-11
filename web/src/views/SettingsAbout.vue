<template>
  <div class="settings-page">
    <h3 style="margin: 0 0 20px">{{ $t('settings_about.title') }}</h3>

    <el-card shadow="hover" class="settings-card">
      <template #header>
        <div style="display: flex; align-items: center; gap: 8px">
          <el-icon :size="18"><InfoFilled /></el-icon>
          <span>{{ $t('settings_about.title') }}</span>
        </div>
      </template>
      <el-descriptions :column="1" border size="small">
        <el-descriptions-item :label="$t('settings_about.system_name')">IPTV Tool</el-descriptions-item>
        <el-descriptions-item :label="$t('settings_about.system_version')">
          <el-tag v-if="appVersion" size="small" type="primary">{{ appVersion }}</el-tag>
          <span v-else style="color: #909399; font-size: 12px">{{ $t('settings_about.fetching') }}</span>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('settings_about.tech_stack')">Vue 3 + Element Plus / Go + Gin + SQLite</el-descriptions-item>
        <el-descriptions-item :label="$t('settings_about.runtime_mode')">{{ $t('settings_about.runtime_value') }}</el-descriptions-item>
        <el-descriptions-item :label="$t('settings_about.repo')">
          <div style="display: flex; align-items: center; gap: 16px">
            <el-link href="https://github.com/super321/iptv-tool" target="_blank" :underline="false" style="display: flex; align-items: center; gap: 4px; color: #303133">
              <svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor">
                <path d="M12 2C6.477 2 2 6.477 2 12c0 4.418 2.865 8.166 6.839 9.489.5.092.682-.217.682-.482 0-.237-.009-.866-.013-1.699-2.782.603-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.462-1.11-1.462-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.564 9.564 0 0112 6.844c.85.004 1.705.114 2.504.336 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.138 20.161 22 16.416 22 12c0-5.523-4.477-10-10-10z"/>
              </svg>
              <span>GitHub</span>
            </el-link>
          </div>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card shadow="hover" class="settings-card sponsor-card">
      <template #header>
        <div style="display: flex; align-items: center; gap: 8px">
          <el-icon :size="18"><Coffee /></el-icon>
          <span>{{ $t('settings_about.sponsor') }}</span>
        </div>
      </template>

      <p class="sponsor-desc">{{ $t('settings_about.sponsor_desc') }}</p>

      <div class="sponsor-list">
        <!-- Ko-fi -->
        <div class="sponsor-item">
          <div class="sponsor-info">
            <div class="sponsor-icon kofi-icon">
              <svg viewBox="0 0 24 24" width="22" height="22" fill="currentColor">
                <path d="M23.881 8.948c-.773-4.085-4.859-4.593-4.859-4.593H.723c-.604 0-.679.798-.679.798s-.082 7.324-.022 11.822c.164 2.424 2.586 2.672 2.586 2.672s8.267-.023 11.966-.049c2.438-.426 2.683-2.566 2.658-3.734 4.352.24 7.422-2.831 6.649-6.916zm-11.062 3.511c-1.246 1.453-4.011 3.976-4.011 3.976s-.121.119-.31.023c-.076-.057-.108-.09-.108-.09-.443-.441-3.368-3.049-4.034-3.954-.709-.965-1.041-2.7-.091-3.71.951-1.01 3.005-1.086 4.363.407 0 0 1.565-1.782 3.468-.963 1.904.82 1.832 3.011.723 4.311zm6.173.478c-.928.116-1.682.028-1.682.028V7.284h1.77s1.971.551 1.971 2.638c0 1.913-.985 2.667-2.059 3.015z"/>
              </svg>
            </div>
            <div class="sponsor-text">
              <span class="sponsor-name">Ko-fi</span>
              <span class="sponsor-sub">Buy me a coffee</span>
            </div>
          </div>
          <el-button type="warning" class="kofi-btn" @click="openLink('https://ko-fi.com/super321')">
            <el-icon style="margin-right: 4px"><Link /></el-icon>
            {{ $t('settings_about.go_sponsor') }}
          </el-button>
        </div>

        <!-- 爱发电 -->
        <div class="sponsor-item">
          <div class="sponsor-info">
            <div class="sponsor-icon afdian-icon">
              <svg viewBox="0 0 24 24" width="22" height="22" fill="currentColor">
                <path d="M12 2L4 7v10l8 5 8-5V7l-8-5zm0 2.18L18 8v8l-6 3.82L6 16V8l6-3.82z"/>
                <path d="M12 8l-4 2.5v5L12 18l4-2.5v-5L12 8zm0 2.18l2 1.25v2.5L12 15.18l-2-1.25v-2.5L12 10.18z"/>
              </svg>
            </div>
            <div class="sponsor-text">
              <span class="sponsor-name">{{ $t('settings_about.afdian') }}</span>
              <span class="sponsor-sub">{{ $t('settings_about.afdian_sub') }}</span>
            </div>
          </div>
          <el-button type="primary" class="afdian-btn" @click="openLink('https://afdian.com/a/super321')">
            <el-icon style="margin-right: 4px"><Link /></el-icon>
            {{ $t('settings_about.go_sponsor') }}
          </el-button>
        </div>

        <!-- Ethereum -->
        <div class="sponsor-item eth-item">
          <div class="sponsor-info">
            <div class="sponsor-icon eth-icon">
              <svg viewBox="0 0 24 24" width="22" height="22" fill="currentColor">
                <path d="M12 1.5l-7 10.657L12 16.5l7-4.343L12 1.5zM12 17.657l-7-4.343L12 22.5l7-9.186-7 4.343z"/>
              </svg>
            </div>
            <div class="sponsor-text">
              <span class="sponsor-name">Ethereum</span>
              <span class="sponsor-sub">{{ $t('settings_about.eth_sub') }}</span>
            </div>
          </div>
          <el-button type="info" @click="showEthDialog = true">
            <el-icon style="margin-right: 4px"><Wallet /></el-icon>
            {{ $t('settings_about.view_address') }}
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- Ethereum Detail Dialog -->
    <el-dialog v-model="showEthDialog" :title="$t('settings_about.eth_title')" width="400px" align-center>
      <div class="eth-dialog-content">
        <img :src="ethQrCode" alt="Ethereum QR Code" class="eth-qrcode" />
        <p class="eth-note">{{ $t('settings_about.eth_note') }}</p>
        <div class="eth-address-box">
          <span class="eth-label">{{ $t('settings_about.wallet_address') }}</span>
          <el-input
            v-model="ethAddress"
            readonly
            size="small"
          >
            <template #append>
              <el-button @click="copyAddress">
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </template>
          </el-input>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { InfoFilled, Link, CopyDocument, Wallet, Star } from '@element-plus/icons-vue'
import { Coffee } from '@element-plus/icons-vue'
import ethQrCode from '../assets/eth_qrcode.jpg'
import api from '../api'

const { t } = useI18n()

const showEthDialog = ref(false)
const ethAddress = ref('0x6989acE6Eb2CC196fAFce3cEcAEC6b6b63716C83')
const appVersion = ref('')

onMounted(async () => {
  try {
    const { data } = await api.get('/system/version')
    appVersion.value = data.version || 'unknown'
  } catch {
    appVersion.value = 'unknown'
  }
})

function openLink(url) {
  window.open(url, '_blank')
}

async function copyAddress() {
  try {
    await navigator.clipboard.writeText(ethAddress.value)
    ElMessage.success(t('settings_about.copied'))
  } catch {
    ElMessage.error(t('settings_about.copy_failed'))
  }
}
</script>

<style scoped>
.settings-page {
  width: 100%;
}
.settings-card {
  max-width: 600px;
}
.sponsor-card {
  margin-top: 20px;
}
.sponsor-desc {
  margin: 0 0 20px;
  color: #606266;
  font-size: 14px;
  line-height: 1.6;
}

/* Sponsor list */
.sponsor-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.sponsor-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-radius: 8px;
  background: #f9fafb;
  transition: background 0.2s;
}
.sponsor-item:hover {
  background: #f0f2f5;
}
.sponsor-info {
  display: flex;
  align-items: center;
  gap: 12px;
}
.sponsor-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
}
.kofi-icon {
  background: linear-gradient(135deg, #ff5e5b, #ff4081);
}
.afdian-icon {
  background: linear-gradient(135deg, #946ce6, #7c4dff);
}
.eth-icon {
  background: linear-gradient(135deg, #627eea, #3c3c3d);
}
.sponsor-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.sponsor-name {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}
.sponsor-sub {
  font-size: 12px;
  color: #909399;
}

/* Buttons */
.kofi-btn {
  background: linear-gradient(135deg, #ff5e5b, #ff4081) !important;
  border: none !important;
  color: #fff !important;
}
.kofi-btn:hover {
  opacity: 0.9;
}
.afdian-btn {
  background: linear-gradient(135deg, #946ce6, #7c4dff) !important;
  border: none !important;
  color: #fff !important;
}
.afdian-btn:hover {
  opacity: 0.9;
}

/* Ethereum Dialog */
.eth-dialog-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}
.eth-qrcode {
  width: 220px;
  height: 220px;
  object-fit: contain;
  border-radius: 8px;
}
.eth-note {
  margin: 0;
  font-size: 13px;
  color: #909399;
}
.eth-address-box {
  width: 100%;
  text-align: center;
}
.eth-label {
  display: block;
  margin-bottom: 8px;
  font-size: 13px;
  color: #606266;
}
</style>
