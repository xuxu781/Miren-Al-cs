<template>
  <el-container class="dy-layout-container">
    <!-- 侧边栏：后台管理系统风格 (纯白底色、无边框、胶囊菜单) -->
    <el-aside :width="isCollapse ? '68px' : '240px'" class="dy-sidebar" :class="{ 'is-mini': isCollapse }">
      <!-- 品牌 Logo -->
      <div class="dy-brand">
        <div class="dy-logo-wrapper">
          <!-- 品牌 Logo 元素的配色点缀 -->
          <div class="dy-logo-icon">
            <span class="cyan-part"></span>
            <span class="red-part"></span>
            <el-icon class="icon-front"><VideoCameraFilled /></el-icon>
          </div>
        </div>
        <div class="dy-brand-text" :class="{ 'is-hidden': isCollapse }">Miren Al后台管理</div>
      </div>

      <el-scrollbar class="dy-sidebar-scroll">
        <div class="menu-group" v-show="!isCollapse">创作与数据</div>
        <el-menu
          :default-active="$route.path"
          class="dy-menu"
          :collapse="isCollapse"
          :collapse-transition="false"
          router
        >
          <el-menu-item index="/dashboard">
            <el-icon><DataLine /></el-icon>
            <template #title>数据概览</template>
          </el-menu-item>
          
          <el-sub-menu index="/user-operations">
            <template #title>
              <el-icon><User /></el-icon>
              <span>用户运营</span>
            </template>
            <el-menu-item index="/users">
              <template #title>用户管理</template>
            </el-menu-item>
            <el-menu-item index="/tasks">
              <template #title>任务列表</template>
            </el-menu-item>
          </el-sub-menu>

          <el-menu-item index="/admin">
            <el-icon><Avatar /></el-icon>
            <template #title>管理员管理</template>
          </el-menu-item>
          <el-menu-item index="/media">
            <el-icon><Picture /></el-icon>
            <template #title>媒体管理</template>
          </el-menu-item>

          <div class="menu-divider" v-show="!isCollapse"></div>
          <div class="menu-group" v-show="!isCollapse">内容管理</div>
          
          <el-sub-menu index="/content">
            <template #title>
              <el-icon><Document /></el-icon>
              <span>内容管理</span>
            </template>
            <el-menu-item index="/content/inspirations">
              <template #title>探索管理</template>
            </el-menu-item>
            <el-menu-item index="/content/inspiration-categories">
              <template #title>探索分类设置</template>
            </el-menu-item>
          </el-sub-menu>

          <div class="menu-divider" v-show="!isCollapse"></div>
          <div class="menu-group" v-show="!isCollapse">商业与营销</div>
          
          <el-sub-menu index="/business">
            <template #title>
              <el-icon><Goods /></el-icon>
              <span>商业中心</span>
            </template>
            <el-menu-item index="/business/cdkeys">
              <template #title>卡密管理</template>
            </el-menu-item>
          </el-sub-menu>
          
          <div class="menu-divider" v-show="!isCollapse"></div>
          <div class="menu-group" v-show="!isCollapse">系统与设置</div>
          
          <el-sub-menu index="/system">
            <template #title>
              <el-icon><Setting /></el-icon>
              <span>系统配置</span>
            </template>
            <el-menu-item index="/settings">
              <template #title>系统设置</template>
            </el-menu-item>
            <el-menu-item index="/models">
              <template #title>模型上游</template>
            </el-menu-item>
          </el-sub-menu>
        </el-menu>
      </el-scrollbar>

      <!-- 侧边栏底部折叠控制 -->
      <div class="sidebar-footer" @click="toggleCollapse">
        <div class="collapse-btn" :class="{ 'is-collapsed': isCollapse }">
          <el-icon><Fold v-if="!isCollapse" /><Expand v-else /></el-icon>
          <span class="collapse-text" :class="{ 'is-hidden': isCollapse }">收起导航</span>
        </div>
      </div>
    </el-aside>

    <el-container class="dy-main-container">
      <!-- 顶部导航：后台 Web 端经典样式 -->
      <el-header class="dy-header">
        <div class="header-left">
          <h2 class="page-title">{{ currentRouteName }}</h2>
        </div>

        <!-- 居中搜索框 (胶囊形) -->
        <div class="header-center" v-if="!isMobile">
          <div class="dy-search-bar">
            <input type="text" placeholder="搜索功能、数据或探索内容" class="dy-search-input" />
            <div class="dy-search-btn">
              <el-icon><Search /></el-icon>
              <span>搜索</span>
            </div>
          </div>
        </div>

        <div class="header-right">
          <div class="icon-action">
            <el-badge is-dot class="dy-badge">
              <el-icon><Bell /></el-icon>
            </el-badge>
          </div>
          
          <el-dropdown trigger="hover" class="dy-avatar-dropdown">
            <div class="dy-avatar-wrapper">
              <img src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png" class="dy-avatar" />
            </div>
            <template #dropdown>
              <el-dropdown-menu class="dy-dropdown-menu">
                <el-dropdown-item>
                  <el-icon><User /></el-icon> 创作者主页
                </el-dropdown-item>
                <el-dropdown-item>
                  <el-icon><Setting /></el-icon> 设置
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout" class="dy-logout">
                  <el-icon><SwitchButton /></el-icon> 退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 主体内容区 -->
      <el-main class="dy-content">
        <div class="dy-page-wrapper">
          <router-view v-slot="{ Component }">
            <transition name="fade-transform" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { 
  DataLine, User, Setting, Picture, Document,
  Search, Bell, SwitchButton, VideoCameraFilled, Fold, Expand, Avatar, Goods
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const isCollapse = ref(false)
const isMobile = ref(false)

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const currentRouteName = computed(() => {
  const map: Record<string, string> = {
    '/dashboard': '数据概览',
    '/admin': '管理员管理',
    '/users': '用户管理',
    '/tasks': '任务列表',
    '/media': '媒体管理',
    '/settings': '系统设置',
    '/models': '模型上游',
    '/business/cdkeys': '卡密管理',
    '/content/inspirations': '探索管理',
    '/content/inspiration-categories': '探索分类设置'
  }
  return map[route.path] || '页面'
})

const handleLogout = () => {
  router.push('/login')
}
</script>

<style scoped>
/* 全局变量 */
:root {
  --dy-red: #FE2C55; /* 主题红 */
  --dy-cyan: #25F4EE; /* 主题青 */
  --dy-bg: #F6F7F9; /* 背景灰 */
  --dy-text-main: #161823; /* 主文本 */
  --dy-text-sub: #808394; /* 副文本 */
  --dy-hover: #F2F4F7; /* hover 浅灰 */
}

.dy-layout-container {
  height: 100vh;
  background-color: var(--dy-bg, #F6F7F9);
  font-family: "PingFang SC", "Microsoft YaHei", -apple-system, BlinkMacSystemFont, sans-serif;
  color: var(--dy-text-main, #161823);
}

/* ================= 侧边栏 ================= */
.dy-sidebar {
  background-color: #ffffff;
  transition: width 0.3s cubic-bezier(0.25, 0.1, 0.25, 1), box-shadow 0.3s ease;
  display: flex;
  flex-direction: column;
  z-index: 20;
  border-right: 1px solid rgba(22, 24, 35, 0.05);
}

.dy-sidebar.is-mini {
  box-shadow: 4px 0 16px rgba(0, 0, 0, 0.06);
  border-right-color: transparent;
}

.dy-brand {
  height: 68px;
  display: flex;
  align-items: center;
  padding: 0 20px;
  gap: 12px;
  overflow: hidden;
  white-space: nowrap;
  flex-shrink: 0;
}

.dy-logo-wrapper {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  flex-shrink: 0;
}

.dy-logo-icon {
  position: relative;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cyan-part {
  position: absolute;
  width: 20px;
  height: 20px;
  border-radius: 4px;
  background: #25F4EE;
  top: 0;
  left: 0;
  mix-blend-mode: multiply;
}

.red-part {
  position: absolute;
  width: 20px;
  height: 20px;
  border-radius: 4px;
  background: #FE2C55;
  bottom: 0;
  right: 0;
  mix-blend-mode: multiply;
}

.icon-front {
  position: relative;
  z-index: 2;
  color: #161823;
  font-size: 18px;
}

.dy-brand-text {
  font-size: 18px;
  font-weight: 700;
  color: #161823;
  letter-spacing: 0.5px;
  transition: opacity 0.2s, width 0.2s;
  opacity: 1;
}

.dy-brand-text.is-hidden {
  opacity: 0;
  width: 0;
  pointer-events: none;
}

/* 菜单区域 */
.dy-sidebar-scroll {
  flex: 1;
}

.menu-group {
  padding: 16px 20px 8px;
  font-size: 12px;
  color: #808394;
  font-weight: 500;
}

.menu-divider {
  height: 1px;
  background-color: rgba(22, 24, 35, 0.05);
  margin: 8px 20px;
}

.dy-menu {
  border-right: none;
  padding: 0 12px;
}

.dy-menu:not(.el-menu--collapse) {
  width: 100%;
}

:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  height: 44px;
  line-height: 44px;
  border-radius: 8px; /* 胶囊状 */
  margin-bottom: 4px;
  color: #161823;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s ease;
}

:deep(.el-menu-item:hover),
:deep(.el-sub-menu__title:hover) {
  background-color: var(--dy-hover, #F2F4F7) !important;
}

:deep(.el-menu-item.is-active) {
  background-color: rgba(254, 44, 85, 0.08) !important;
  color: #FE2C55 !important;
  font-weight: 600;
}

:deep(.el-menu-item .el-icon) {
  font-size: 18px;
  margin-right: 12px;
}

/* 侧边栏底部 */
.sidebar-footer {
  padding: 12px 16px;
  border-top: 1px solid rgba(22, 24, 35, 0.05);
  flex-shrink: 0;
}

.collapse-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 40px;
  padding: 0 12px;
  border-radius: 8px;
  cursor: pointer;
  color: #808394;
  font-size: 14px;
  transition: all 0.2s;
  overflow: hidden;
  white-space: nowrap;
}

.collapse-btn.is-collapsed {
  justify-content: center;
  padding: 0;
}

.collapse-text {
  transition: opacity 0.2s, width 0.2s;
  opacity: 1;
}

.collapse-text.is-hidden {
  opacity: 0;
  width: 0;
  pointer-events: none;
}

.collapse-btn:hover {
  background-color: var(--dy-hover, #F2F4F7);
  color: #161823;
}

/* ================= 头部 ================= */
.dy-main-container {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.dy-header {
  height: 68px;
  background-color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  z-index: 10;
  border-bottom: 1px solid rgba(22, 24, 35, 0.05);
}

.header-left {
  display: flex;
  align-items: center;
  min-width: 150px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  color: #161823;
}

/* 居中搜索框 */
.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
  padding: 0 40px;
}

.dy-search-bar {
  display: flex;
  align-items: center;
  width: 100%;
  max-width: 480px;
  height: 40px;
  background: var(--dy-hover, #F2F4F7);
  border-radius: 20px;
  border: 1px solid transparent;
  transition: all 0.3s;
  overflow: hidden;
}

.dy-search-bar:focus-within {
  border-color: rgba(22, 24, 35, 0.1);
  background: #ffffff;
  box-shadow: 0 4px 16px rgba(0,0,0,0.04);
}

.dy-search-input {
  flex: 1;
  height: 100%;
  border: none;
  background: transparent;
  padding: 0 16px;
  font-size: 14px;
  color: #161823;
  outline: none;
}

.dy-search-input::placeholder {
  color: #808394;
}

.dy-search-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 0 16px;
  height: 100%;
  cursor: pointer;
  color: #161823;
  font-size: 14px;
  font-weight: 500;
  border-left: 1px solid rgba(22, 24, 35, 0.08);
}

.dy-search-btn:hover {
  color: #FE2C55;
}

/* 头部右侧 */
.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
  min-width: 120px;
  justify-content: flex-end;
}

.icon-action {
  font-size: 22px;
  color: #161823;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-action:hover {
  color: #FE2C55;
}

.dy-badge :deep(.el-badge__content.is-fixed.is-dot) {
  right: 2px;
  top: 4px;
  background-color: #FE2C55;
}

.dy-avatar-wrapper {
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dy-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: 1px solid rgba(22, 24, 35, 0.08);
  transition: transform 0.2s;
  object-fit: cover;
}

.dy-avatar-wrapper:hover .dy-avatar {
  transform: scale(1.05);
}

/* ================= 主体内容 ================= */
.dy-content {
  padding: 24px;
  overflow-y: auto;
}

.dy-page-wrapper {
  max-width: 1200px;
  margin: 0 auto;
}

/* 页面切换动画 */
.fade-transform-leave-active,
.fade-transform-enter-active {
  transition: all 0.3s cubic-bezier(0.25, 0.1, 0.25, 1);
}
.fade-transform-enter-from {
  opacity: 0;
  transform: translateY(10px) scale(0.99);
}
.fade-transform-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.99);
}

/* 下拉菜单定制 */
.dy-dropdown-menu .el-dropdown-menu__item {
  padding: 10px 24px;
  font-size: 14px;
  color: #161823;
  border-radius: 6px;
  margin: 0 8px;
}
.dy-dropdown-menu .el-dropdown-menu__item:hover {
  background-color: #F2F4F7;
  color: #161823;
}
.dy-dropdown-menu .dy-logout {
  color: #161823;
}
.dy-dropdown-menu .dy-logout:hover {
  background-color: rgba(254, 44, 85, 0.08);
  color: #FE2C55;
}
</style>
