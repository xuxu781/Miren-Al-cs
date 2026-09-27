<template>
  <div class="media-container">
    <el-card class="box-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>媒体管理</span>
          <div class="header-actions">
            <el-radio-group v-model="activeCategory" size="small" style="margin-right: 16px;">
              <el-radio-button label="全部">全部</el-radio-button>
              <el-radio-button label="用户图片">用户图片</el-radio-button>
              <el-radio-button label="系统图片">系统图片</el-radio-button>
            </el-radio-group>
            <el-button type="primary" :icon="Refresh" @click="fetchMedia" :loading="loading">刷新</el-button>
          </div>
        </div>
      </template>

      <div class="media-grid" v-loading="loading">
        <div v-for="item in filteredMediaList" :key="item.path" class="media-item">
          <div class="image-wrapper">
            <el-image 
              :src="item.url" 
              :preview-src-list="[item.url]"
              fit="cover"
              class="media-image"
              lazy
            >
              <template #error>
                <div class="image-slot">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <div class="image-overlay">
              <el-button type="danger" size="small" circle :icon="Delete" @click="handleDelete(item)" />
            </div>
          </div>
          <div class="media-info">
            <div class="media-name" :title="item.name">{{ item.name }}</div>
            <div class="media-meta">
              <span>{{ formatSize(item.size) }}</span>
              <el-tag size="small" :type="item.category === '系统图片' ? 'warning' : 'info'">{{ item.category }}</el-tag>
            </div>
            <div class="media-meta" style="margin-top: 4px;">
              <span>{{ formatDate(item.created_at) }}</span>
            </div>
          </div>
        </div>
        
        <el-empty class="empty-state" v-if="!loading && filteredMediaList.length === 0" description="暂无媒体文件" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Picture, Delete, Refresh } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import axios from 'axios'

const mediaList = ref<any[]>([])
const loading = ref(false)
const activeCategory = ref('全部')

const filteredMediaList = computed(() => {
  if (activeCategory.value === '全部') {
    return mediaList.value
  }
  return mediaList.value.filter(item => item.category === activeCategory.value)
})

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (dateStr: string) => {
  // Check if it's already a formatted string like "2026-09-05 15:04:05"
  if (dateStr && dateStr.includes(' ')) {
    return dateStr
  }
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const fetchMedia = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await axios.get(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/media/list`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (res.data && res.data.code === 200) {
      mediaList.value = res.data.data.list || []
    }
  } catch (error) {
    ElMessage.error('获取媒体列表失败')
  } finally {
    loading.value = false
  }
}

const handleDelete = (item: any) => {
  ElMessageBox.confirm(
    `确定要删除该图片吗？此操作不可恢复。`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      const token = localStorage.getItem('token')
      const res = await axios.post(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/media/delete`, { path: item.path }, {
        headers: { Authorization: `Bearer ${token}` }
      })
      if (res.data && res.data.code === 200) {
        ElMessage.success('删除成功')
        fetchMedia()
      } else {
        ElMessage.error(res.data.message || '删除失败')
      }
    } catch (error) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {})
}

onMounted(() => {
  fetchMedia()
})
</script>

<style scoped>
.media-container {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  min-height: 200px;
}

.media-item {
  border: 1px solid #ebeef5;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s;
}

.media-item:hover {
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
  transform: translateY(-2px);
}

.image-wrapper {
  position: relative;
  width: 100%;
  height: 200px;
  background-color: #f5f7fa;
}

.media-image {
  width: 100%;
  height: 100%;
}

.image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  color: #909399;
  font-size: 30px;
}

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: 0;
  transition: opacity 0.3s;
}

.image-wrapper:hover .image-overlay {
  opacity: 1;
}

.media-info {
  padding: 12px;
  background-color: #fff;
}

.media-name {
  font-size: 14px;
  color: #303133;
  margin-bottom: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.media-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #909399;
}

.empty-state {
  grid-column: 1 / -1;
}
</style>