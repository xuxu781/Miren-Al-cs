<template>
  <div class="settings-container fade-in">
    <div class="page-header">
      <div class="header-info">
        <div class="icon-wrapper">
          <el-icon :size="24" color="#409eff"><Setting /></el-icon>
        </div>
        <div class="header-title">
          <h2>系统设置</h2>
          <span class="subtitle">管理网站基础信息及相关配置</span>
        </div>
      </div>
    </div>

    <el-card shadow="never" class="settings-card">
      <el-tabs v-model="activeTab" class="custom-tabs">
        <el-tab-pane label="基础设置" name="basic_settings">
          <div class="tab-content">
            <h3 class="section-title">网站基础设置</h3>
            <p class="section-desc">配置网站的基本信息，如名称、Logo、描述等。</p>
            
            <el-form 
              ref="basicFormRef"
              :model="settingsForm" 
              :rules="rules"
              label-width="120px" 
              label-position="top"
              class="settings-form"
              v-loading="loading"
            >
              <el-form-item label="网站名称" prop="site_name">
                <el-input 
                  v-model="settingsForm.site_name" 
                  placeholder="例如: Miren Al" 
                >
                  <template #prefix>
                    <el-icon><Monitor /></el-icon>
                  </template>
                </el-input>
                <div class="form-tip">显示在浏览器标签页及页面顶部等位置。</div>
              </el-form-item>

              <el-form-item label="网站 Logo" prop="site_logo">
                <div class="logo-upload-wrapper">
                  <el-input 
                    v-model="settingsForm.site_logo" 
                    placeholder="例如: /logo.png 或 https://example.com/logo.png" 
                  >
                    <template #prefix>
                      <el-icon><Picture /></el-icon>
                    </template>
                  </el-input>
                  <el-upload
                    class="logo-uploader"
                    action="#"
                    :auto-upload="false"
                    :show-file-list="false"
                    :on-change="handleLogoChange"
                  >
                    <el-button type="primary" plain>选择图片</el-button>
                  </el-upload>
                </div>
                <div class="form-tip">网站的图标地址，支持外链 URL 或相对路径。也可以直接上传。</div>
                <div class="logo-preview" v-if="previewLogoUrl || settingsForm.site_logo">
                  <img :src="previewLogoUrl || formatLogoUrl(settingsForm.site_logo)" alt="Logo Preview" />
                </div>
              </el-form-item>

              <el-form-item label="网站描述 (SEO)" prop="site_description">
                <el-input 
                  v-model="settingsForm.site_description" 
                  type="textarea"
                  :rows="3"
                  placeholder="请输入网站描述..." 
                >
                </el-input>
                <div class="form-tip">用于搜索引擎收录及分享时的摘要展示。</div>
              </el-form-item>

              <el-form-item label="版权信息" prop="site_copyright">
                <el-input 
                  v-model="settingsForm.site_copyright" 
                  placeholder="例如: © 2026 Miren Al All Rights Reserved." 
                >
                </el-input>
                <div class="form-tip">显示在网站底部的版权声明。</div>
              </el-form-item>

              <el-form-item>
                <el-button type="primary" class="save-btn" @click="handleSaveBasic" :loading="saving">
                  保存基础设置
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <el-tab-pane label="运营设置" name="operation_settings">
          <div class="tab-content">
            <h3 class="section-title">运营设置</h3>
            <p class="section-desc">配置平台运营相关参数，例如新用户注册奖励等。</p>
            
            <el-form 
              ref="operationFormRef"
              :model="settingsForm" 
              label-width="120px" 
              label-position="top"
              class="settings-form"
              v-loading="loading"
            >
              <el-form-item label="开启充值功能" prop="enable_recharge">
                <el-switch
                  v-model="settingsForm.enable_recharge"
                  active-text="开启"
                  inactive-text="关闭"
                />
                <div class="form-tip">开启后，用户前台将显示“充值”按钮，点击后跳转至设置的充值链接。</div>
              </el-form-item>

              <el-form-item label="充值跳转链接" prop="recharge_link" v-if="settingsForm.enable_recharge">
                <el-input 
                  v-model="settingsForm.recharge_link" 
                  placeholder="例如: https://pay.example.com" 
                >
                  <template #prefix>
                    <el-icon><Link /></el-icon>
                  </template>
                </el-input>
                <div class="form-tip">用户点击“充值”按钮后将新窗口跳转到此链接。</div>
              </el-form-item>

              <el-form-item label="新用户默认积分" prop="default_register_points">
                <el-input-number 
                  v-model="settingsForm.default_register_points" 
                  :min="0"
                  :step="1"
                  placeholder="例如: 5" 
                >
                </el-input-number>
                <div class="form-tip">新用户注册成功后默认发放的初始积分数量。</div>
              </el-form-item>

              <el-form-item>
                <el-button type="primary" class="save-btn" @click="handleSaveOperation" :loading="saving">
                  保存运营设置
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <el-tab-pane label="系统图片设置" name="image_settings">
          <div class="tab-content">
            <h3 class="section-title">系统图片设置</h3>
            <p class="section-desc">配置图片访问域名及大模型参考图的传参方式。</p>
            
            <el-form 
              ref="imageFormRef"
              :model="settingsForm" 
              label-width="120px" 
              label-position="top"
              class="settings-form"
              v-loading="loading"
            >
              <el-form-item label="图片域名" prop="image_domain">
                <el-input 
                  v-model="settingsForm.image_domain" 
                  placeholder="例如: https://img.example.com" 
                >
                  <template #prefix>
                    <el-icon><Link /></el-icon>
                  </template>
                </el-input>
                <div class="form-tip">配置上传图片的公网访问域名（包含 http/https）。如果留空，将默认使用当前后端服务地址。当你的大模型接口需要公网 URL 才能下载参考图时，请配置此项并确保可公网访问。</div>
              </el-form-item>

              <el-form-item label="参考图传参模式" prop="enable_reference_image">
                <el-switch
                  v-model="settingsForm.enable_reference_image"
                  active-text="URL 模式"
                  inactive-text="本地上传模型"
                />
                <div class="form-tip">开启时，将参考图的公网 URL 发送给大模型（如 Midjourney）；关闭时，后端会直接将图片文件通过 Multipart Form-Data 上传至模型接口（如 OpenAI 的 /v1/images/edits）。</div>
              </el-form-item>

              <el-form-item label="自动清理历史图片" prop="enable_auto_clean_images">
                <div style="display: flex; align-items: center; justify-content: space-between; width: 100%;">
                  <el-switch
                    v-model="settingsForm.enable_auto_clean_images"
                    active-text="开启"
                    inactive-text="关闭"
                  />
                  <el-button link type="primary" @click="viewCleanLog" v-if="settingsForm.enable_auto_clean_images">
                    <el-icon><Document /></el-icon>查看清理日志
                  </el-button>
                </div>
                <div class="form-tip">开启后，系统将自动扫描并删除未被任何用户任务（包括重新生成、重编辑等）所引用的多余/历史图片，释放服务器存储空间。</div>
              </el-form-item>

              <el-form-item label="清理策略（时长）" prop="auto_clean_images_time" v-if="settingsForm.enable_auto_clean_images">
                <div style="display: flex; gap: 10px; align-items: center; width: 100%;">
                  <el-input-number 
                    v-model="settingsForm.auto_clean_images_time" 
                    :min="0"
                    :step="1"
                    placeholder="例如: 7" 
                    style="flex: 1;"
                  >
                  </el-input-number>
                  <el-select v-model="settingsForm.auto_clean_images_unit" style="width: 100px;">
                    <el-option label="天" value="days"></el-option>
                    <el-option label="小时" value="hours"></el-option>
                    <el-option label="分钟" value="minutes"></el-option>
                  </el-select>
                </div>
                <div class="form-tip">设置清理多久之前的未被引用的图片。设为 0 表示立即清理（不建议，可能影响正在生成的任务，建议设置一定缓冲时间）。</div>
              </el-form-item>

              <el-form-item>
                <el-button type="primary" class="save-btn" @click="handleSaveImage" :loading="saving">
                  保存图片设置
                </el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-dialog
      v-model="logDialogVisible"
      title="自动清理历史图片日志"
      width="600px"
      class="custom-dialog"
    >
      <div v-loading="logLoading" class="log-content-container">
        <template v-if="cleanLogs.length > 0">
          <div v-for="(log, index) in cleanLogs" :key="index" class="log-item">
            {{ log }}
          </div>
        </template>
        <el-empty v-else description="暂无清理日志" />
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" plain @click="refreshCleanLog" :loading="refreshingLog">
            <el-icon><Refresh /></el-icon>刷新日志
          </el-button>
          <el-button type="danger" @click="clearCleanLog" :loading="clearingLog">清空日志</el-button>
          <el-button @click="logDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Setting, Monitor, Picture, Link, Document, Refresh } from '@element-plus/icons-vue'

const activeTab = ref('basic_settings')
const loading = ref(false)
const saving = ref(false)
const basicFormRef = ref()
const imageFormRef = ref()
const operationFormRef = ref()

const settingsForm = ref({
  site_name: '',
  site_logo: '',
  site_description: '',
  site_copyright: '',
  image_domain: '',
  enable_reference_image: true,
  enable_auto_clean_images: false,
  auto_clean_images_time: 7,
  auto_clean_images_unit: 'days',
  default_register_points: 5,
  enable_recharge: false,
  recharge_link: ''
})

const rules = {
  site_name: [{ required: true, message: '请输入网站名称', trigger: 'blur' }]
}

const logDialogVisible = ref(false)
const logLoading = ref(false)
const cleanLogs = ref<string[]>([])
const clearingLog = ref(false)
const refreshingLog = ref(false)

const refreshCleanLog = async () => {
  refreshingLog.value = true
  try {
    const res = await fetch(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/settings/get?key_name=auto_clean_images_log`, {
      headers: { ...getHeaders(), 'Content-Type': 'application/json' }
    })
    if (res.ok) {
      const json = await res.json()
      if (json.data && json.data.key_value) {
        cleanLogs.value = json.data.key_value.split('\n').filter((l: string) => l.trim() !== '')
      } else {
        cleanLogs.value = []
      }
      ElMessage.success('日志已刷新')
    }
  } catch (error) {
    ElMessage.error('获取日志失败')
  } finally {
    refreshingLog.value = false
  }
}

const clearCleanLog = async () => {
  clearingLog.value = true
  try {
    const res = await fetch(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/settings/save`, {
      method: 'POST',
      headers: { ...getHeaders(), 'Content-Type': 'application/json' },
      body: JSON.stringify({
        key_name: 'auto_clean_images_log', 
        key_value: '', 
        description: '自动清理历史图片日志'
      })
    })
    if (res.ok) {
      ElMessage.success('日志已清空')
      cleanLogs.value = []
    } else {
      ElMessage.error('清空日志失败')
    }
  } catch (error) {
    ElMessage.error('清空日志失败')
  } finally {
    clearingLog.value = false
  }
}

const viewCleanLog = async () => {
  logDialogVisible.value = true
  logLoading.value = true
  cleanLogs.value = []
  try {
    const res = await fetch(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/settings/get?key_name=auto_clean_images_log`, {
      headers: { ...getHeaders(), 'Content-Type': 'application/json' }
    })
    if (res.ok) {
      const json = await res.json()
      if (json.data && json.data.key_value) {
        cleanLogs.value = json.data.key_value.split('\n').filter((l: string) => l.trim() !== '')
      }
    }
  } catch (error) {
    ElMessage.error('获取日志失败')
  } finally {
    logLoading.value = false
  }
}

const getHeaders = () => {
  const token = localStorage.getItem('token')
  return {
    'Authorization': `Bearer ${token}`
  }
}

const formatLogoUrl = (url: string) => {
  if (!url) return ''
  if (url.startsWith('http') || url.startsWith('data:') || url.startsWith('blob:')) return url
  return `${url.startsWith('/') ? url : '/' + url}`
}

const selectedLogoFile = ref<File | null>(null)
const previewLogoUrl = ref('')

const handleLogoChange = (file: any) => {
  const isImage = file.raw.type.startsWith('image/')
  const isLt2M = file.raw.size / 1024 / 1024 < 2

  if (!isImage) {
    ElMessage.error('上传头像图片只能是图片格式!')
    return
  }
  if (!isLt2M) {
    ElMessage.error('上传头像图片大小不能超过 2MB!')
    return
  }
  
  selectedLogoFile.value = file.raw
  previewLogoUrl.value = URL.createObjectURL(file.raw)
}

// 获取指定配置
const fetchSetting = async (keyName: string) => {
  try {
    const res = await fetch(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/settings/get?key_name=${keyName}`, {
      headers: { ...getHeaders(), 'Content-Type': 'application/json' }
    })
    if (res.ok) {
      const json = await res.json()
      return json.data?.key_value || ''
    }
  } catch (error) {
    console.error(`Failed to fetch setting ${keyName}`, error)
  }
  return ''
}

// 批量保存配置
const saveSetting = async (keyName: string, keyValue: string, description: string) => {
  try {
    const res = await fetch(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/settings/save`, {
      method: 'POST',
      headers: { ...getHeaders(), 'Content-Type': 'application/json' },
      body: JSON.stringify({
        key_name: keyName,
        key_value: keyValue,
        description: description
      })
    })
    if (!res.ok) {
      throw new Error(`Failed to save ${keyName}`)
    }
  } catch (error) {
    throw error
  }
}

const loadSettings = async () => {
  loading.value = true
  try {
    settingsForm.value.site_name = await fetchSetting('site_name') || 'Miren Al'
    settingsForm.value.site_logo = await fetchSetting('site_logo')
    settingsForm.value.site_description = await fetchSetting('site_description')
    settingsForm.value.site_copyright = await fetchSetting('site_copyright') || '© 2026 Miren Al All Rights Reserved.'
    settingsForm.value.image_domain = await fetchSetting('image_domain')
    const enableRefImg = await fetchSetting('enable_reference_image')
    settingsForm.value.enable_reference_image = enableRefImg !== 'false'
    const enableAutoClean = await fetchSetting('enable_auto_clean_images')
    settingsForm.value.enable_auto_clean_images = enableAutoClean === 'true'
    const autoCleanTime = await fetchSetting('auto_clean_images_time')
    settingsForm.value.auto_clean_images_time = autoCleanTime ? parseInt(autoCleanTime, 10) : 7
    const autoCleanUnit = await fetchSetting('auto_clean_images_unit')
    settingsForm.value.auto_clean_images_unit = autoCleanUnit || 'days'
    const defaultPoints = await fetchSetting('default_register_points')
    settingsForm.value.default_register_points = defaultPoints ? parseInt(defaultPoints, 10) : 5
    
    const enableRecharge = await fetchSetting('enable_recharge')
    settingsForm.value.enable_recharge = enableRecharge === 'true'
    settingsForm.value.recharge_link = await fetchSetting('recharge_link') || ''
  } catch (error) {
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

const handleSaveBasic = async () => {
  if (!basicFormRef.value) return
  await basicFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      saving.value = true
      try {
        // Upload the selected logo if any
        if (selectedLogoFile.value) {
          const formData = new FormData()
          formData.append('file', selectedLogoFile.value)
          const res = await fetch(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/system/upload`, {
            method: 'POST',
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: formData
          })
          const json = await res.json()
          if (res.ok && json.code === 200 && json.data?.url) {
            settingsForm.value.site_logo = json.data.url
            selectedLogoFile.value = null
            previewLogoUrl.value = ''
          } else {
            ElMessage.error(json.message || 'Logo上传失败')
            saving.value = false
            return
          }
        }

        await saveSetting('site_name', settingsForm.value.site_name, '网站名称')
        await saveSetting('site_logo', settingsForm.value.site_logo, '网站 Logo')
        await saveSetting('site_description', settingsForm.value.site_description, '网站描述')
        await saveSetting('site_copyright', settingsForm.value.site_copyright, '版权信息')
        
        // Clear cached local storage to force frontend refresh
        localStorage.removeItem('site_logo')
        localStorage.removeItem('admin_site_logo')
        localStorage.removeItem('admin_site_name')
        
        ElMessage.success('基础设置保存成功')
      } catch (error) {
        ElMessage.error('基础设置保存失败')
      } finally {
        saving.value = false
      }
    }
  })
}

const handleSaveImage = async () => {
  if (!imageFormRef.value) return
  await imageFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      saving.value = true
      try {
        await saveSetting('image_domain', settingsForm.value.image_domain, '图片域名')
        await saveSetting('enable_reference_image', settingsForm.value.enable_reference_image ? 'true' : 'false', '是否开启参考图')
        await saveSetting('enable_auto_clean_images', settingsForm.value.enable_auto_clean_images ? 'true' : 'false', '是否开启自动清理历史图片')
        await saveSetting('auto_clean_images_time', String(settingsForm.value.auto_clean_images_time), '清理历史图片策略时长')
        await saveSetting('auto_clean_images_unit', settingsForm.value.auto_clean_images_unit, '清理历史图片策略单位')
        ElMessage.success('图片设置保存成功')
      } catch (error) {
        ElMessage.error('图片设置保存失败')
      } finally {
        saving.value = false
      }
    }
  })
}

const handleSaveOperation = async () => {
  if (!operationFormRef.value) return
  await operationFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      saving.value = true
      try {
        await saveSetting('default_register_points', String(settingsForm.value.default_register_points), '新用户默认积分')
        await saveSetting('enable_recharge', settingsForm.value.enable_recharge ? 'true' : 'false', '是否开启充值功能')
        await saveSetting('recharge_link', settingsForm.value.recharge_link, '充值跳转链接')
        ElMessage.success('运营设置保存成功')
      } catch (error) {
        ElMessage.error('运营设置保存失败')
      } finally {
        saving.value = false
      }
    }
  })
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.fade-in {
  animation: fadeIn 0.4s ease-in-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.settings-container {
  padding: 24px;
  background-color: #f1f5f9;
  min-height: calc(100vh - 60px);
}

.log-content-container {
  min-height: 200px;
  max-height: 400px;
  overflow-y: auto;
  background-color: #f8fafc;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
}

.log-item {
  font-family: monospace;
  font-size: 13px;
  color: #334155;
  padding: 6px 0;
  border-bottom: 1px dashed #e2e8f0;
  line-height: 1.5;
}

.log-item:last-child {
  border-bottom: none;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  background: #fff;
  padding: 20px 24px;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.header-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.icon-wrapper {
  width: 48px;
  height: 48px;
  background: #e6f2ff;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-title h2 {
  margin: 0 0 4px 0;
  font-size: 20px;
  color: #1e293b;
  font-weight: 600;
}

.subtitle {
  font-size: 13px;
  color: #64748b;
}

.settings-card {
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  background: #fff;
  min-height: 500px;
}

.tab-content {
  padding: 12px 24px;
}

.section-title {
  margin: 0 0 8px 0;
  font-size: 18px;
  color: #334155;
}

.section-desc {
  margin: 0 0 24px 0;
  color: #64748b;
  font-size: 14px;
}

.settings-form {
  max-width: 600px;
}

.form-tip {
  font-size: 12px;
  color: #94a3b8;
  margin-top: 4px;
  line-height: 1.4;
}

.save-btn {
  padding: 10px 32px;
  border-radius: 8px;
  font-weight: 500;
  margin-top: 12px;
}

.logo-upload-wrapper {
  display: flex;
  gap: 12px;
  width: 100%;
}

.logo-uploader {
  flex-shrink: 0;
}

.logo-preview {
  margin-top: 12px;
  width: 80px;
  height: 80px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f8fafc;
  overflow: hidden;
}

.logo-preview img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

:deep(.el-tabs__nav-wrap::after) {
  height: 1px;
  background-color: #e2e8f0;
}

:deep(.el-tabs__item) {
  font-size: 15px;
  padding: 0 24px;
  height: 48px;
  line-height: 48px;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 0 0 1px #e2e8f0 inset;
  padding: 4px 12px;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #cbd5e1 inset;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #409eff inset !important;
}
</style>
