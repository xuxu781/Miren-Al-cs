import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import ElementPlus, { ElMessage } from 'element-plus'
import 'element-plus/dist/index.css'
import './style.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import axios from 'axios'

// 拦截全局 axios 响应
axios.interceptors.response.use(
  response => response,
  error => {
    const requestUrl = error.config?.url || ''
    if (error.response && error.response.status === 401 && !requestUrl.includes('/login')) {
      if (!(window as any).__has_shown_401_error) {
        ;(window as any).__has_shown_401_error = true
        ElMessage.error('登录已过期，请重新登录')
        localStorage.removeItem('token')
        router.push('/login').finally(() => {
          setTimeout(() => {
            ;(window as any).__has_shown_401_error = false
          }, 1000)
        })
      }
    }
    return Promise.reject(error)
  }
)

// 拦截全局 fetch 响应
const originalFetch = window.fetch
window.fetch = async (...args) => {
  const response = await originalFetch(...args)
  const req = args[0]
  const url = typeof req === 'string' ? req : ('url' in req ? req.url : String(req))
  if (response.status === 401 && !url.includes('/login')) {
    if (!(window as any).__has_shown_401_error) {
      ;(window as any).__has_shown_401_error = true
      ElMessage.error('登录已过期，请重新登录')
      localStorage.removeItem('token')
      router.push('/login').finally(() => {
        setTimeout(() => {
          ;(window as any).__has_shown_401_error = false
        }, 1000)
      })
    }
  }
  return response
}

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(ElementPlus, {
  locale: zhCn,
})

app.mount('#app')
