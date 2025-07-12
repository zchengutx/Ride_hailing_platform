<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-header">
        <h1>欢迎使用打车平台</h1>
        <p>安全便捷的出行选择</p>
      </div>
      
      <div class="login-content">
        <LoginForm @login-success="handleLoginSuccess" />
      </div>
      
      <div class="back-home">
        <el-button text @click="$router.push('/')">
          <el-icon><ArrowLeft /></el-icon>
          返回首页
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import LoginForm from '@/components/LoginForm.vue'

const router = useRouter()

const handleLoginSuccess = () => {
  router.push('/')
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-container {
  background: white;
  border-radius: 16px;
  padding: 40px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 40px rgba(0,0,0,0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-header h1 {
  font-size: 24px;
  color: #333;
  margin-bottom: 8px;
}

.login-header p {
  color: #666;
  font-size: 14px;
  margin: 0;
}

.login-content {
  margin-bottom: 24px;
}

.back-home {
  text-align: center;
}

.back-home .el-button {
  color: #666;
  font-size: 14px;
}

@media (max-width: 480px) {
  .login-container {
    padding: 24px;
    margin: 0 16px;
  }
}
</style> 