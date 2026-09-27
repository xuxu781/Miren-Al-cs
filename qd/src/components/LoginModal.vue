<template>
  <el-dialog
    v-model="authStore.showLoginModal"
    :show-close="false"
    :lock-scroll="false"
    width="400px"
    class="custom-login-dialog"
    destroy-on-close
  >
    <div class="login-box">
      <div class="close-btn" @click="authStore.closeLogin()">
        <el-icon><Close /></el-icon>
      </div>
      
      <div class="login-content">
        <div class="form-header">
          <div class="logo-placeholder" v-if="!siteLogo">
            <el-icon><MagicStick /></el-icon>
          </div>
          <div class="logo-image" v-else>
            <img :src="siteLogo" alt="logo" />
          </div>
          <h3>{{ isLogin ? '欢迎回来' : '创建账号' }}</h3>
          <p>{{ isLogin ? '登录以继续使用 AI 生图' : '注册以开启您的创作之旅' }}</p>
        </div>
        
        <el-form :model="form" label-position="top" size="large" class="login-form">
          <el-form-item>
            <el-input 
              v-model="form.email" 
              placeholder="请输入邮箱" 
              :prefix-icon="Message"
            />
          </el-form-item>
          <el-form-item>
            <el-input 
              v-model="form.password" 
              type="password" 
              placeholder="请输入密码" 
              :prefix-icon="Lock"
              show-password 
              @keyup.enter="handleSubmit"
            />
          </el-form-item>
          
          <div class="form-options" v-if="isLogin">
            <el-checkbox v-model="form.remember">自动登录</el-checkbox>
            <el-link type="primary" :underline="false">忘记密码</el-link>
          </div>
          
          <el-form-item>
            <el-button type="primary" class="login-btn" @click="handleSubmit" :loading="loading">
              {{ isLogin ? '登录' : '注册' }}
            </el-button>
          </el-form-item>
          
          <div class="toggle-mode">
            <span class="text-muted">{{ isLogin ? '还没有账号？' : '已有账号？' }}</span>
            <el-link type="primary" :underline="false" @click="toggleMode">
              {{ isLogin ? '立即注册' : '直接登录' }}
            </el-link>
          </div>
        </el-form>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Lock, Message, Close, MagicStick } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'
import { useAuthStore } from '@/store/auth'
import defaultLogo from '@/assets/logo/logo.png'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const isLogin = ref(true)
const cachedLogo = localStorage.getItem('site_logo')
const siteLogo = ref(cachedLogo || defaultLogo)

const fetchSettings = async () => {
  try {
    const res = await request.get('/api/public/settings')
    if (res.data && res.data.code === 200) {
      const settings = res.data.data || {}
      if (settings.site_logo !== undefined) {
        if (settings.site_logo) {
          siteLogo.value = settings.site_logo
          localStorage.setItem('site_logo', settings.site_logo)
        } else {
          siteLogo.value = defaultLogo
          localStorage.removeItem('site_logo')
        }
      }
    }
  } catch (error) {
    console.warn('Failed to fetch public settings', error)
  }
}

onMounted(() => {
  fetchSettings()
})

const form = reactive({
  email: '',
  password: '',
  remember: false
})

const toggleMode = () => {
  isLogin.value = !isLogin.value
  form.password = ''
  form.email = ''
}

const handleSubmit = async () => {
  if (!form.email || !form.password) {
    ElMessage.warning('请输入邮箱和密码')
    return
  }
  loading.value = true
  try {
    const endpoint = isLogin.value ? '/api/user/login' : '/api/user/register'
    
    const payload = { email: form.email, password: form.password }

    const res = await request.post(endpoint, payload)
    
    if (res.data.token) {
      // 登录成功
      ElMessage.success('登录成功')
      authStore.setToken(res.data.token)
      authStore.setUserInfo(res.data.user)
      
      authStore.closeLogin()
      
      // 如果当前不是首页，则跳转到首页；如果是首页，不需要刷新了，因为使用了 store 的响应式数据
      if (router.currentRoute.value.path !== '/') {
        router.push('/')
      }
    } else if (res.data.message === '注册成功') {
      // 注册成功
      ElMessage.success('注册成功，请登录')
      isLogin.value = true
      form.password = '' // 清空密码让用户重新输入
    }
  } catch (error: any) {
    console.error(error)
    ElMessage.error(error.response?.data?.error || (isLogin.value ? '登录失败' : '注册失败'))
  } finally {
    loading.value = false
  }
}
</script>

<style>
/* 覆盖 el-dialog 默认样式，使其透明无边框 */
.custom-login-dialog {
  background: transparent !important;
  box-shadow: none !important;
  padding: 0 !important;
}
.custom-login-dialog .el-dialog__header {
  display: none;
}
.custom-login-dialog .el-dialog__body {
  padding: 0 !important;
}

@media (max-width: 768px) {
  .custom-login-dialog {
    width: 85% !important;
    max-width: 320px !important;
  }
}
</style>

<style scoped>
.login-box {
  display: flex;
  flex-direction: column;
  width: 400px;
  max-width: 100%;
  background-color: #ffffff;
  border-radius: 16px;
  box-shadow: 0 12px 32px 4px rgba(0, 0, 0, 0.04), 0 8px 20px rgba(0, 0, 0, 0.08);
  margin: 0 auto;
  position: relative;
}

.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  z-index: 20;
  cursor: pointer;
  font-size: 20px;
  color: #909399;
  transition: color 0.3s;
  padding: 4px;
  border-radius: 4px;
}

.close-btn:hover {
  color: #303133;
  background-color: #f5f7fa;
}

.login-content {
  padding: 40px 32px 32px;
}

.form-header {
  margin-bottom: 32px;
  text-align: center;
}

.logo-placeholder {
  width: 48px;
  height: 48px;
  background-color: #f3f4f6;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
  font-size: 24px;
  color: #111827;
}

.logo-image {
  width: 48px;
  height: 48px;
  margin: 0 auto 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-image img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  border-radius: 12px;
}

.form-header h3 {
  font-size: 24px;
  color: #111827;
  margin: 0 0 8px 0;
  font-weight: 600;
}

.form-header p {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  margin-top: -12px;
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  border-radius: 8px;
  font-weight: 500;
  background-color: #111827;
  border-color: #111827;
}

.login-btn:hover {
  background-color: #374151;
  border-color: #374151;
}

.toggle-mode {
  text-align: center;
  margin-top: 16px;
  font-size: 14px;
}

.text-muted {
  color: #6b7280;
  margin-right: 8px;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  padding: 4px 12px;
  box-shadow: 0 0 0 1px #e5e7eb inset;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #111827 inset !important;
}

:deep(.el-input__inner) {
  height: 40px;
}

:deep(.el-input__prefix) {
  font-size: 18px;
  color: #9ca3af;
}

@media (max-width: 768px) {
  .login-box {
    width: 100%;
  }
  .login-content {
    padding: 24px 20px 20px;
  }
  .form-header {
    margin-bottom: 20px;
  }
  .form-header h3 {
    font-size: 20px;
  }
  .logo-placeholder, .logo-image {
    width: 40px;
    height: 40px;
    margin: 0 auto 12px;
  }
  .login-btn {
    height: 40px;
    font-size: 15px;
  }
  :deep(.el-input__inner) {
    height: 36px;
  }
  .form-options {
    margin-bottom: 16px;
  }
}
</style>
