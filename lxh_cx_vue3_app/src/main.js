import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'

// 配置axios基础URL - 使用相对路径，通过Vite代理
axios.defaults.baseURL = '/v1/api'
axios.defaults.timeout = 10000

const app = createApp(App)

// 将axios挂载到全局属性
app.config.globalProperties.$http = axios

app.use(router)
app.mount('#app') 