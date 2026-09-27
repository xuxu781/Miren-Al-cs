<script setup lang="ts">
import { onMounted } from 'vue'
import LoginModal from '@/components/LoginModal.vue'

onMounted(async () => {
  try {
    const apiUrl = (window as any).APP_CONFIG?.API_BASE_URL || ''
    const res = await fetch(`${apiUrl}/api/public/settings`)
    if (res.ok) {
      const json = await res.json()
      if (json.code === 200 && json.data) {
        localStorage.setItem('site_settings', JSON.stringify(json.data))
        // 触发自定义事件以便其他组件感知
        window.dispatchEvent(new StorageEvent('storage', {
          key: 'site_settings',
          newValue: JSON.stringify(json.data)
        }))
      }
    }
  } catch (error) {
    console.error('Failed to load public settings:', error)
  }
})
</script>

<template>
  <router-view />
  <LoginModal />
</template>

<style>
/* 去除默认边距 */
html, body {
  margin: 0 !important;
  padding: 0 !important;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  overflow: hidden !important;
  width: 100vw !important;
  height: 100vh !important;
}
</style>
