<template>
  <div class="install-container">
    <div class="install-box">
      <div class="install-left">
        <div class="install-info">
          <h2>Welcome to</h2>
          <h1>Miren Al System</h1>
          <p>只需几步，即可完成系统初始化配置</p>
          <ul class="features">
            <li><el-icon><Check /></el-icon> 自动配置数据库</li>
            <li><el-icon><Check /></el-icon> 初始化管理员账号</li>
            <li><el-icon><Check /></el-icon> 导入基础数据结构</li>
          </ul>
        </div>
        <div class="decoration-circle top-circle"></div>
        <div class="decoration-circle bottom-circle"></div>
      </div>
      <div class="install-right">
        <el-card class="install-card" shadow="never">
          <div class="card-header">
            <h3>系统安装向导</h3>
            <p>{{ activeStep === 0 ? '请配置 MySQL 数据库连接信息' : '请设置系统初始管理员账号' }}</p>
          </div>
          
          <el-steps :active="activeStep" finish-status="success" align-center style="margin-bottom: 30px">
            <el-step title="数据库配置" />
            <el-step title="管理员设置" />
          </el-steps>

          <el-form :model="form" :rules="rules" ref="formRef" label-position="top" size="large">
            <div v-show="activeStep === 0">
              <el-row :gutter="20">
                <el-col :span="16">
                  <el-form-item label="数据库地址" prop="db_host">
                    <el-input v-model="form.db_host" placeholder="例如: 127.0.0.1">
                      <template #prefix><el-icon><Platform /></el-icon></template>
                    </el-input>
                  </el-form-item>
                </el-col>
                <el-col :span="8">
                  <el-form-item label="端口" prop="db_port">
                    <el-input v-model="form.db_port" placeholder="3306" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-form-item label="数据库用户" prop="db_user">
                <el-input v-model="form.db_user" placeholder="例如: root">
                  <template #prefix><el-icon><User /></el-icon></template>
                </el-input>
              </el-form-item>
              <el-form-item label="数据库密码" prop="db_pass">
                <el-input v-model="form.db_pass" type="password" show-password placeholder="请输入密码">
                  <template #prefix><el-icon><Lock /></el-icon></template>
                </el-input>
              </el-form-item>
              <el-form-item label="数据库名称" prop="db_name">
                <el-input v-model="form.db_name" placeholder="例如: lty_db">
                  <template #prefix><el-icon><Coin /></el-icon></template>
                </el-input>
              </el-form-item>
              <el-form-item class="mt-4">
                <el-button type="primary" @click="nextStep" :loading="loading" class="install-btn">
                  下一步
                </el-button>
              </el-form-item>
            </div>

            <div v-show="activeStep === 1">
              <el-form-item label="管理员账号" prop="admin_user">
                <el-input v-model="form.admin_user" placeholder="例如: admin">
                  <template #prefix><el-icon><User /></el-icon></template>
                </el-input>
              </el-form-item>
              <el-form-item label="管理员密码" prop="admin_pass">
                <el-input v-model="form.admin_pass" type="password" show-password placeholder="请输入管理员密码">
                  <template #prefix><el-icon><Lock /></el-icon></template>
                </el-input>
              </el-form-item>
              <el-form-item class="mt-4" style="display: flex; gap: 15px;">
                <el-button @click="prevStep" style="flex: 1; height: 44px; border-radius: 8px; margin: 0;">
                  上一步
                </el-button>
                <el-button type="primary" @click="handleInstall" :loading="loading" style="flex: 1; height: 44px; border-radius: 8px; margin: 0;">
                  开始安装
                </el-button>
              </el-form-item>
            </div>
          </el-form>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Lock, Platform, Coin, Check } from '@element-plus/icons-vue'
import axios from 'axios'

const formRef = ref<FormInstance>()
const loading = ref(false)
const activeStep = ref(0)

const form = reactive({
  db_host: '127.0.0.1',
  db_port: '3306',
  db_user: 'root',
  db_pass: '123456',
  db_name: 'lty_db',
  admin_user: 'admin',
  admin_pass: '123456'
})

const rules = reactive<FormRules>({
  db_host: [{ required: true, message: '请输入数据库地址', trigger: 'blur' }],
  db_port: [{ required: true, message: '请输入数据库端口', trigger: 'blur' }],
  db_user: [{ required: true, message: '请输入数据库用户', trigger: 'blur' }],
  db_name: [{ required: true, message: '请输入数据库名称', trigger: 'blur' }],
  admin_user: [{ required: true, message: '请输入管理员账号', trigger: 'blur' }],
  admin_pass: [{ required: true, message: '请输入管理员密码', trigger: 'blur' }]
})

const nextStep = async () => {
  if (!formRef.value) return
  await formRef.value.validateField(['db_host', 'db_port', 'db_user', 'db_name'], async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res = await axios.post(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/check-db`, {
          db_host: form.db_host,
          db_port: form.db_port,
          db_user: form.db_user,
          db_pass: form.db_pass,
          db_name: form.db_name
        })
        if (res.status === 200) {
          ElMessage.success('数据库连接成功')
          activeStep.value = 1
        }
      } catch (error: any) {
        ElMessage.error(error.response?.data?.error || '数据库连接失败，请检查配置')
      } finally {
        loading.value = false
      }
    }
  })
}

const prevStep = () => {
  activeStep.value = 0
}

const handleInstall = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res = await axios.post(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/install`, form)
        if (res.status === 200) {
          ElMessage.success('系统安装成功！')
          window.location.href = '/login'
        }
      } catch (error: any) {
        ElMessage.error(error.response?.data?.error || '安装失败，请检查数据库配置是否正确')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.install-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f0f2f5;
  background-image: url('data:image/svg+xml,%3Csvg width="100%25" height="100%25" xmlns="http://www.w3.org/2000/svg"%3E%3Cdefs%3E%3Cpattern id="grid" width="40" height="40" patternUnits="userSpaceOnUse"%3E%3Cpath d="M 40 0 L 0 0 0 40" fill="none" stroke="%23e8e8e8" stroke-width="1"/%3E%3C/pattern%3E%3C/defs%3E%3Crect width="100%25" height="100%25" fill="url(%23grid)"/%3E%3C/svg%3E');
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
}

.install-box {
  display: flex;
  width: 1000px;
  height: 600px;
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 12px 32px 4px rgba(0, 0, 0, 0.04), 0 8px 20px rgba(0, 0, 0, 0.08);
  overflow: hidden;
}

.install-left {
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

.install-info {
  position: relative;
  z-index: 1;
}

.install-info h2 {
  font-size: 24px;
  font-weight: 400;
  margin: 0 0 10px 0;
  color: rgba(255, 255, 255, 0.8);
}

.install-info h1 {
  font-size: 36px;
  font-weight: 700;
  margin: 0 0 20px 0;
  color: #ffffff;
}

.install-info p {
  font-size: 15px;
  color: rgba(255, 255, 255, 0.7);
  margin: 0 0 30px 0;
}

.features {
  list-style: none;
  padding: 0;
  margin: 0;
}

.features li {
  display: flex;
  align-items: center;
  font-size: 15px;
  margin-bottom: 16px;
  color: rgba(255, 255, 255, 0.9);
}

.features li .el-icon {
  margin-right: 12px;
  font-size: 18px;
  background: rgba(255, 255, 255, 0.1);
  padding: 6px;
  border-radius: 50%;
  color: #409EFF;
}

.install-right {
  flex: 7;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px;
  background: #ffffff;
}

.install-card {
  width: 100%;
  max-width: 440px;
  border: none;
  background: transparent;
}

:deep(.el-card__body) {
  padding: 0;
}

.card-header {
  margin-bottom: 30px;
  text-align: center;
}

.card-header h3 {
  font-size: 24px;
  color: #303133;
  margin: 0 0 8px 0;
  font-weight: 600;
}

.card-header p {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.install-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  border-radius: 8px;
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

.mt-4 {
  margin-top: 16px;
}
</style>
