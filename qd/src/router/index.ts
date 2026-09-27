import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/Generate.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/explore',
      name: 'explore',
      component: () => import('../views/Generate.vue'),
      meta: { requiresAuth: true }
    }
  ]
})

// 路由守卫：检查是否登录
router.beforeEach((_to, _from) => {
  // 不再强制跳转或弹出登录框，直接放行
  return true
})

export default router
