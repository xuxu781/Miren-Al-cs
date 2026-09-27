import axios from 'axios'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/store/auth'

let apiUrl = ''
if ((window as any).APP_CONFIG?.API_BASE_URL) {
  apiUrl = (window as any).APP_CONFIG.API_BASE_URL
} else {
  if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
    apiUrl = window.location.protocol + "//" + window.location.hostname + ":8088"
  } else {
    apiUrl = window.location.origin
  }
}

const request = axios.create({
  baseURL: apiUrl,
  timeout: 120000
})

// 请求拦截器：自动携带 Token
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 用于防止短时间内重复显示相同的错误提示
let lastErrorMsg = ''
let lastErrorTime = 0
const showErrorMsg = (msg: string, type: 'error' | 'warning' = 'error') => {
  const now = Date.now()
  if (msg === lastErrorMsg && now - lastErrorTime < 3000) {
    return // 3秒内不重复提示相同的错误
  }
  lastErrorMsg = msg
  lastErrorTime = now
  if (type === 'warning') {
    ElMessage.warning(msg)
  } else {
    ElMessage.error(msg)
  }
}

// 响应拦截器：处理全局错误（包括 401 和 403 自动退出）
request.interceptors.response.use(
  (response) => {
    // 拦截因配置错误导致的 SPA HTML 兜底响应（如 API_BASE_URL 错误时请求到了前端本身）
    if (typeof response.data === 'string' && response.data.trim().toLowerCase().startsWith('<!doctype html>')) {
      showErrorMsg('服务器响应异常，请稍后重试')
      return Promise.reject(new Error('Invalid response format: Received HTML instead of JSON'))
    }
    return response
  },
  (error) => {
    if (error.response) {
      const status = error.response.status
      if (status === 401 || status === 403) {
        // 清除本地信息 (保存下token用于判断是否需要显示提示)
        const oldToken = localStorage.getItem('token')
        localStorage.removeItem('token')
        localStorage.removeItem('userInfo')
        
        // 提示用户
        const errMsg = error.response.data?.error || error.response.data?.message || ''
        if (status === 403) {
          showErrorMsg(errMsg || '您的账号已被封禁/禁用，请联系管理员')
        } else {
          // 如果本来就没有token，或者是单纯的“未提供授权Token”，则静默处理（只打开登录弹窗）
          // 由于后端返回可能是 {"error": "未提供授权Token"} 或者是其它字符串，这里统一对 "未提供授权Token" 不做全局错误弹窗
          if (oldToken && !errMsg.includes('未提供授权Token') && !errMsg.includes('未提供授权 Token')) {
             showErrorMsg(errMsg || '登录已过期或无权限，请重新登录')
          }
        }
        
        // 如果本来就没有 token 或者后端提示是未提供授权，不再强制弹出登录框
        // 除非用户触发了必须登录的操作（交由业务层去主动调用 openLogin）
        // 并且如果本来有 token 但现在失效了，可以打开登录框
        if (oldToken) {
          const authStore = useAuthStore()
          authStore.openLogin()
        } else if (status === 401 && (!errMsg || errMsg.includes('未提供授权Token') || errMsg.includes('未提供授权 Token'))) {
           // 满足用户期望: 刷新（或者首次访问）未登录时，主动弹出登录框
           const authStore = useAuthStore()
           authStore.openLogin()
        } else if (!oldToken && status === 401) {
           // 为了兜底，如果没有 token 时后端返回了 "无效的Token" / "Token格式错误" / 或者其他明确的 401，我们弹出登录框
           if (errMsg && !errMsg.includes('未提供授权Token') && !errMsg.includes('未提供授权 Token')) {
             const authStore = useAuthStore()
             authStore.openLogin()
           }
        }
      } else if (status === 404) {
        showErrorMsg('请求的服务不存在，请稍后重试')
      } else if (status >= 500) {
        showErrorMsg('服务器繁忙，请稍后重试')
      } else {
        // 其他错误全局提示
        const errMsg = error.response.data?.error || error.response.data?.message;
        if (errMsg === '积分不足') {
          showErrorMsg('积分不足，请前往充值页面进行充值', 'warning')
        } else {
          showErrorMsg(errMsg || '请求失败，请稍后重试')
        }
      }
    } else {
      showErrorMsg('无法连接到服务器，请检查网络连接或稍后重试')
    }
    return Promise.reject(error)
  }
)

export default request
