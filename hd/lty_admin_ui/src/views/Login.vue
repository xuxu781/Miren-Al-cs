<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-left">
        <div class="login-info">
          <h2>Welcome to</h2>
          <h1>{{ siteName }}</h1>
          <p>高效、简洁、美观的现代化后台管理系统</p>
        </div>
        <div class="decoration-circle top-circle"></div>
        <div class="decoration-circle bottom-circle"></div>
      </div>
      <div class="login-right">
        <div class="login-form-container">
          <div class="form-header">
            <div class="logo-image" v-if="siteLogo">
              <img :src="siteLogo" alt="logo" />
            </div>
            <h3>账号登录</h3>
            <p>请输入您的用户名和密码</p>
          </div>
          <el-form :model="form" label-position="top" size="large" class="login-form">
            <el-form-item>
              <el-input 
                v-model="form.username" 
                placeholder="用户名" 
                :prefix-icon="User"
              />
            </el-form-item>
            <el-form-item>
              <el-input 
                v-model="form.password" 
                type="password" 
                placeholder="密码" 
                :prefix-icon="Lock"
                show-password 
                @keyup.enter="handleLogin"
              />
            </el-form-item>
            <div class="form-options">
              <el-checkbox v-model="form.remember">记住密码</el-checkbox>
              <el-link type="primary" :underline="false">忘记密码？</el-link>
            </div>
            <el-form-item>
              <el-button type="primary" class="login-btn" @click="handleLogin" :loading="loading">
                登录
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const loading = ref(false)
const siteLogo = ref(localStorage.getItem('admin_site_logo') || '')
const siteName = ref(localStorage.getItem('admin_site_name') || 'Miren Al')
const form = reactive({
  username: '',
  password: '',
  remember: false
})

const fetchSettings = async () => {
  try {
    const res = await fetch(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/public/settings`)
    const json = await res.json()
    if (res.ok && json.code === 200 && json.data) {
      if (json.data.site_logo !== undefined) {
        if (json.data.site_logo) {
          siteLogo.value = json.data.site_logo
          localStorage.setItem('admin_site_logo', json.data.site_logo)
        } else {
          siteLogo.value = ''
          localStorage.removeItem('admin_site_logo')
        }
      }
      if (json.data.site_name !== undefined) {
        if (json.data.site_name) {
          siteName.value = json.data.site_name
          localStorage.setItem('admin_site_name', json.data.site_name)
        } else {
          siteName.value = 'Miren Al'
          localStorage.removeItem('admin_site_name')
        }
      }
    }
  } catch (error) {
    console.error('Failed to fetch public settings', error)
  }
}

onMounted(() => {
  fetchSettings()
})

const handleLogin = async () => {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    const res = await fetch(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: form.username,
        password: form.password
      })
    })
    const data = await res.json()
    if (res.ok) {
      ElMessage.success('登录成功')
      localStorage.setItem('token', data.token)
      router.push('/dashboard')
    } else {
      ElMessage.error(data.error || '登录失败')
    }
  } catch (error) {
    ElMessage.error('网络错误，请稍后重试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f0f2f5;
  background-image: url('data:image/svg+xml,%3Csvg width="100%25" height="100%25" xmlns="http://www.w3.org/2000/svg"%3E%3Cdefs%3E%3Cpattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"%3E%3Cpath d="M 40 0 L 0 0 0 40" fill="none" stroke="%23e8e8e8" stroke-width="1"/%3E%3C/pattern%3E%3C/defs%3E%3Crect width="100%25" height="100%25" fill="url(%23grid)"/%3E%3C/svg%3E');
  position: relative;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
}

.login-box {
  display: flex;
  width: 900px;
  height: 520px;
  background-color: #ffffff;
  border-radius: 16px;
  box-shadow: 0 12px 32px 4px rgba(0, 0, 0, 0.04), 0 8px 20px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  z-index: 10;
}

.login-left {
  flex: 5;
  background: linear-gradient(135deg, #1e1e2d 0%, #2b2b40 100%);
  color: #ffffff;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 50px;
  position: relative;
  overflow: hidden;
}

.decoration-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.05);
}

.top-circle {
  width: 300px;
  height: 300px;
  top: -100px;
  right: -100px;
}

.bottom-circle {
  width: 200px;
  height: 200px;
  bottom: -50px;
  left: -50px;
}

.login-info {
  position: relative;
  z-index: 1;
}

.login-info h2 {
  font-size: 24px;
  font-weight: 400;
  margin: 0 0 10px 0;
  color: rgba(255, 255, 255, 0.8);
}

.login-info h1 {
  font-size: 42px;
  font-weight: 700;
  margin: 0 0 20px 0;
  color: #ffffff;
}

.login-info p {
  font-size: 15px;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
  line-height: 1.6;
}

.login-right {
  flex: 6;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px;
  background-color: #ffffff;
}

.login-form-container {
  width: 100%;
  max-width: 340px;
}

.form-header {
  margin-bottom: 30px;
  text-align: center;
}

.form-header h3 {
  font-size: 26px;
  color: #303133;
  margin: 0 0 8px 0;
  font-weight: 600;
}

.form-header p {
  font-size: 14px;
  color: #909399;
  margin: 0;
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
  border-radius: 8px;
}

.login-form {
  margin-top: 20px;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  margin-top: -10px;
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  border-radius: 8px;
  margin-top: 10px;
  font-weight: 500;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  padding: 4px 12px;
  box-shadow: 0 0 0 1px #dcdfe6 inset;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #409eff inset !important;
}

:deep(.el-input__inner) {
  height: 40px;
}

:deep(.el-input__prefix) {
  font-size: 18px;
  color: #909399;
}
</style>
