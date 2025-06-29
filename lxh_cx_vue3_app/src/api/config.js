import axios from 'axios'

// 创建axios实例
const apiClient = axios.create({
  baseURL: 'http://localhost:8899/v1/api',  // 直接指向后端服务
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
apiClient.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
apiClient.interceptors.response.use(
  response => {
    // 直接返回响应数据
    return response.data
  },
  error => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          // token过期或未登录
          localStorage.removeItem('token')
          window.location.href = '/login'
          break
        case 403:
          // 权限不足
          console.error('没有权限访问该资源')
          break
        case 500:
          console.error('服务器内部错误:', error.response.data)
          break
        default:
          console.error('请求失败:', error.response.data)
      }
    } else if (error.request) {
      console.error('网络请求失败，请检查网络连接')
    }
    return Promise.reject(error)
  }
)

export default apiClient
 