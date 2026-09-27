import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import axios from 'axios'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/install',
    name: 'Install',
    component: () => import('../views/Install.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('../layout/index.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue')
      },
      {
        path: 'admin',
        name: 'Admin',
        component: () => import('../views/Admin.vue')
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('../views/Users.vue')
      },
      {
        path: 'tasks',
        name: 'Tasks',
        component: () => import('../views/Tasks.vue')
      },
      {
        path: 'models',
        name: 'Models',
        component: () => import('../views/Models.vue')
      },
      {
        path: 'media',
        name: 'Media',
        component: () => import('../views/Media.vue')
      },
      {
        path: 'business/cdkeys',
        name: 'Cdkeys',
        component: () => import('../views/Cdkeys.vue')
      },
      {
        path: 'content/inspirations',
        name: 'Inspirations',
        component: () => import('../views/Inspirations.vue')
      },
      {
        path: 'content/inspiration-categories',
        name: 'InspirationCategories',
        component: () => import('../views/InspirationCategories.vue')
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../views/Settings.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

let isCheckedInstall = false
let isInstalled = false

router.beforeEach(async (to, _from, next) => {
  if (!isCheckedInstall) {
    try {
      const res = await axios.get(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/check-install`)
      
      // 确保后端返回了正确的 JSON 格式，而不是前端开发服务器返回的 index.html (SPA Fallback)
      if (res.data && typeof res.data.installed === 'boolean') {
        isCheckedInstall = true
        isInstalled = res.data.installed
      } else {
        throw new Error('Invalid response format')
      }
    } catch (error) {
      console.error('Failed to check installation status', error)
      // 如果后端没起或网络错误，为了避免无限跳转安装页，这里标记检查完毕并假定已安装
      // 这样用户会留在当前页，页面内的具体接口请求会正常报网络错误
      isCheckedInstall = true
      isInstalled = true 
    }
  }

  // 如果已经安装，且访问的是安装页，跳转到首页
  if (isInstalled && to.path === '/install') {
    next('/')
    return
  }

  // 如果未安装，且访问的不是安装页，跳转到安装页
  if (isCheckedInstall && !isInstalled && to.path !== '/install') {
    next('/install')
    return
  }
  
  next()
})

export default router
