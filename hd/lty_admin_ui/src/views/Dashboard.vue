<template>
  <div class="dashboard-container">
    <!-- 数据概览 -->
    <el-row :gutter="20" class="mb-4">
      <el-col :span="6" v-for="(stat, index) in statistics" :key="index">
        <el-card shadow="hover" class="stat-card" :body-style="{ padding: '20px' }">
          <div class="stat-content">
            <div class="stat-icon" :style="{ backgroundColor: stat.bgColor }">
              <el-icon :style="{ color: stat.color }"><component :is="stat.icon" /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-title">{{ stat.title }}</div>
              <div class="stat-value">{{ stat.value }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <!-- 快捷操作 -->
      <el-col :span="16">
        <el-card shadow="hover" class="box-card">
          <template #header>
            <div class="card-header">
              <span>系统概览</span>
              <el-tag type="success" size="small" effect="light">运行正常</el-tag>
            </div>
          </template>
          <div class="welcome-section">
            <div class="welcome-text">
              <h3>欢迎使用 Miren Al 后台管理系统</h3>
              <p>这是一个基于 Vue 3 + Element Plus 构建的现代化后台管理模板，提供了开箱即用的功能，帮助您快速搭建企业级中后台产品。</p>
              <div class="action-buttons">
                <el-button type="primary">系统设置</el-button>
                <el-button plain>查看文档</el-button>
              </div>
            </div>
            <img src="../assets/vue.svg" class="welcome-img" alt="vue" />
          </div>
        </el-card>
      </el-col>

      <!-- 最近动态 -->
      <el-col :span="8">
        <el-card shadow="hover" class="box-card">
          <template #header>
            <div class="card-header">
              <span>最近动态</span>
              <el-button link type="primary">查看更多</el-button>
            </div>
          </template>
          <el-timeline class="custom-timeline">
            <el-timeline-item
              v-for="(activity, index) in activities"
              :key="index"
              :type="activity.type"
              :color="activity.color"
              :size="activity.size"
              :timestamp="activity.timestamp"
            >
              {{ activity.content }}
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { User, DataLine, ShoppingCart, ChatDotRound } from '@element-plus/icons-vue'
import { markRaw } from 'vue'

const statistics = [
  {
    title: '总访问量',
    value: '1,234,567',
    icon: markRaw(DataLine),
    color: '#40c9c6',
    bgColor: '#e6f7f7'
  },
  {
    title: '新增用户',
    value: '8,432',
    icon: markRaw(User),
    color: '#36a3f7',
    bgColor: '#e6f3fe'
  },
  {
    title: '本月订单',
    value: '6,231',
    icon: markRaw(ShoppingCart),
    color: '#f4516c',
    bgColor: '#fde8eb'
  },
  {
    title: '未读消息',
    value: '95',
    icon: markRaw(ChatDotRound),
    color: '#34bfa3',
    bgColor: '#e6f8f5'
  }
]

const activities = [
  {
    content: '管理员 admin 登录系统',
    timestamp: '2026-08-17 10:30',
    type: 'primary',
    size: 'large'
  },
  {
    content: '用户管理模块更新完成',
    timestamp: '2026-08-16 15:20',
    color: '#0bbd87'
  },
  {
    content: '系统初始化配置完成',
    timestamp: '2026-08-15 09:00',
    type: 'success'
  },
  {
    content: 'Miren Al 项目立项',
    timestamp: '2026-08-10 14:00',
    type: 'info'
  }
]
</script>

<style scoped>
.dashboard-container {
  animation: fadeIn 0.4s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.mb-4 {
  margin-bottom: 20px;
}

.stat-card {
  border: none;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 20px -10px rgba(0,0,0,0.1) !important;
}

.stat-content {
  display: flex;
  align-items: center;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-right: 16px;
  transition: transform 0.3s;
}

.stat-card:hover .stat-icon {
  transform: scale(1.05);
}

.stat-icon .el-icon {
  font-size: 28px;
}

.stat-info {
  flex: 1;
}

.stat-title {
  font-size: 14px;
  color: #8c939d;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 26px;
  font-weight: bold;
  color: #303133;
  line-height: 1;
}

.box-card {
  border: none;
  border-radius: 8px;
  min-height: 400px;
}

:deep(.el-card__header) {
  border-bottom: 1px solid #f0f2f5;
  padding: 16px 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 500;
  color: #303133;
  font-size: 16px;
}

.welcome-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 10px;
}

.welcome-text h3 {
  margin: 0 0 16px 0;
  font-size: 24px;
  color: #303133;
  font-weight: 600;
}

.welcome-text p {
  color: #606266;
  line-height: 1.8;
  margin-bottom: 24px;
  font-size: 14px;
  max-width: 500px;
}

.welcome-img {
  width: 140px;
  height: 140px;
  animation: float 6s ease-in-out infinite;
}

@keyframes float {
  0% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
  100% { transform: translateY(0px); }
}

.action-buttons .el-button {
  padding: 10px 24px;
}

.custom-timeline {
  padding-top: 10px;
  padding-left: 5px;
}

:deep(.el-timeline-item__content) {
  color: #303133;
  font-size: 14px;
}

:deep(.el-timeline-item__timestamp) {
  color: #909399;
  font-size: 13px;
}
</style>
