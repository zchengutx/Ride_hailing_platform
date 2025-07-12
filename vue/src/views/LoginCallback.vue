<template>
  <div class="login-callback">
    <div class="callback-container">
      <div class="loading-content">
        <div class="loading-spinner">
          <div class="spinner"></div>
        </div>
        <h2>{{ statusText }}</h2>
        <p>{{ detailText }}</p>
        
        <!-- 进度条 -->
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: progress + '%' }"></div>
        </div>
        
        <!-- 成功/失败图标 -->
        <div v-if="status === 'success'" class="status-icon success">
          <el-icon size="48"><SuccessFilled /></el-icon>
        </div>
        
        <div v-if="status === 'error'" class="status-icon error">
          <el-icon size="48"><CircleCloseFilled /></el-icon>
          <el-button type="primary" @click="redirectToLogin" class="retry-btn">
            重新登录
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { SuccessFilled, CircleCloseFilled } from '@element-plus/icons-vue'
import Cookies from 'js-cookie'

const router = useRouter()
const route = useRoute()

const status = ref<'loading' | 'success' | 'error'>('loading')
const statusText = ref('正在处理登录信息...')
const detailText = ref('请稍候，即将为您跳转到打车页面')
const progress = ref(0)

// 处理登录回调
const handleLoginCallback = async () => {
  try {
    // 从URL参数中获取token
    const token = route.query.token as string
    
    if (!token) {
      throw new Error('未找到登录凭证')
    }
    
    // 验证token格式
    if (token.length < 10) {
      throw new Error('登录凭证格式不正确')
    }
    
    // 模拟处理进度
    statusText.value = '验证登录凭证...'
    detailText.value = '正在验证您的登录信息'
    progress.value = 30
    
    await new Promise(resolve => setTimeout(resolve, 800))
    
    statusText.value = '保存登录状态...'
    detailText.value = '正在保存您的登录信息'
    progress.value = 60
    
    // 将token存入cookie和localStorage
    // Cookie设置7天过期
    Cookies.set('authToken', token, { 
      expires: 7,
      secure: false,  // 开发环境设为false，生产环境应设为true
      sameSite: 'lax'
    })
    
    // 同时存入localStorage作为备用
    localStorage.setItem('token', token)
    localStorage.setItem('loginTime', Date.now().toString())
    
    await new Promise(resolve => setTimeout(resolve, 800))
    
    statusText.value = '登录成功！'
    detailText.value = '即将跳转到打车页面'
    progress.value = 100
    status.value = 'success'
    
    // 显示成功消息
    ElMessage.success('登录成功，欢迎使用打车服务！')
    
    // 等待一秒后跳转
    setTimeout(() => {
      router.replace('/')
    }, 1500)
    
  } catch (error) {
    console.error('登录回调处理失败:', error)
    
    status.value = 'error'
    statusText.value = '登录失败'
    detailText.value = error instanceof Error ? error.message : '处理登录信息时发生错误，请重试'
    progress.value = 0
    
    ElMessage.error('登录处理失败，请重新登录')
  }
}

// 跳转到登录页面
const redirectToLogin = () => {
  router.replace('/login')
}

// 组件挂载时执行
onMounted(() => {
  // 检查是否有token参数
  if (!route.query.token) {
    status.value = 'error'
    statusText.value = '缺少登录参数'
    detailText.value = '未找到有效的登录凭证，请重新登录'
    return
  }
  
  // 开始处理登录回调
  handleLoginCallback()
})
</script>

<style scoped>
.login-callback {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.callback-container {
  background: white;
  border-radius: 16px;
  padding: 40px;
  max-width: 480px;
  width: 100%;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
}

.loading-content h2 {
  color: #333;
  font-size: 24px;
  margin: 20px 0 12px 0;
  font-weight: 600;
}

.loading-content p {
  color: #666;
  font-size: 16px;
  margin: 0 0 30px 0;
  line-height: 1.5;
}

.loading-spinner {
  margin-bottom: 20px;
}

.spinner {
  width: 48px;
  height: 48px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #ff7e00;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.progress-bar {
  width: 100%;
  height: 6px;
  background: #f0f0f0;
  border-radius: 3px;
  overflow: hidden;
  margin: 20px 0;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #ff7e00 0%, #ff6600 100%);
  border-radius: 3px;
  transition: width 0.3s ease;
}

.status-icon {
  margin: 20px 0;
}

.status-icon.success {
  color: #67c23a;
}

.status-icon.error {
  color: #f56c6c;
}

.retry-btn {
  margin-top: 16px;
  background: linear-gradient(135deg, #ff7e00 0%, #ff6600 100%);
  border: none;
  border-radius: 8px;
  padding: 12px 24px;
  font-size: 16px;
  font-weight: 500;
}

.retry-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255, 126, 0, 0.3);
}

/* 响应式设计 */
@media (max-width: 480px) {
  .callback-container {
    padding: 24px;
    margin: 20px;
  }
  
  .loading-content h2 {
    font-size: 20px;
  }
  
  .loading-content p {
    font-size: 14px;
  }
}
</style> 