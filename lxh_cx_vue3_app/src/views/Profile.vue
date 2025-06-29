<template>
  <div class="profile">
    <!-- 顶部导航 -->
    <header class="header">
      <div class="header-content">
        <button @click="goBack" class="back-btn">
          <span class="back-icon">←</span>
        </button>
        <h1 class="page-title">个人中心</h1>
        <div class="header-placeholder"></div>
      </div>
    </header>

    <div class="profile-container">
      <!-- 用户信息卡片 -->
      <div class="user-card">
        <div class="user-avatar">
          <span class="avatar-icon">🐱</span>
        </div>
        <div class="user-info">
          <h2 class="user-name">{{ userInfo.mobile || '罗小黑用户' }}</h2>
          <p class="user-desc">专享会员 · 已认证</p>
        </div>
        <div class="user-actions">
          <button @click="editProfile" class="edit-btn">编辑</button>
        </div>
      </div>

      <!-- 统计信息 -->
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-number">{{ stats.rides }}</div>
          <div class="stat-label">总行程</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ stats.distance }}</div>
          <div class="stat-label">总里程(km)</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ stats.points }}</div>
          <div class="stat-label">积分</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ stats.coupons }}</div>
          <div class="stat-label">优惠券</div>
        </div>
      </div>

      <!-- 快捷服务 -->
      <div class="services-section">
        <h3 class="section-title">快捷服务</h3>
        <div class="services-grid">
          <div class="service-item" @click="goToOrders">
            <div class="service-icon">📋</div>
            <div class="service-text">我的订单</div>
          </div>
          <div class="service-item" @click="goToWallet">
            <div class="service-icon">💰</div>
            <div class="service-text">我的钱包</div>
          </div>
          <div class="service-item" @click="goToAddress">
            <div class="service-icon">📍</div>
            <div class="service-text">常用地址</div>
          </div>
          <div class="service-item" @click="goToInvite">
            <div class="service-icon">🎁</div>
            <div class="service-text">邀请好友</div>
          </div>
        </div>
      </div>

      <!-- 菜单列表 -->
      <div class="menu-section">
        <div class="menu-item" @click="showBindMobile = true" v-if="!userInfo.mobile">
          <div class="menu-icon">📱</div>
          <div class="menu-text">绑定手机号</div>
          <div class="menu-arrow">›</div>
        </div>
        <div class="menu-item" @click="goToSecurity">
          <div class="menu-icon">🔒</div>
          <div class="menu-text">账号与安全</div>
          <div class="menu-arrow">›</div>
        </div>
        <div class="menu-item" @click="goToSettings">
          <div class="menu-icon">⚙️</div>
          <div class="menu-text">设置</div>
          <div class="menu-arrow">›</div>
        </div>
        <div class="menu-item" @click="goToHelp">
          <div class="menu-icon">❓</div>
          <div class="menu-text">帮助与反馈</div>
          <div class="menu-arrow">›</div>
        </div>
        <div class="menu-item" @click="goToAbout">
          <div class="menu-icon">ℹ️</div>
          <div class="menu-text">关于我们</div>
          <div class="menu-arrow">›</div>
        </div>
      </div>

      <!-- 退出登录 -->
      <div class="logout-section">
        <button @click="confirmLogout" class="logout-btn">
          退出登录
        </button>
      </div>
    </div>

    <!-- 绑定手机号弹窗 -->
    <div v-if="showBindMobile" class="modal-overlay" @click="closeBindModal">
      <div class="bind-modal" @click.stop>
        <div class="modal-header">
          <h3>绑定手机号</h3>
          <button @click="closeBindModal" class="close-btn">×</button>
        </div>
        <div class="modal-body">
          <div class="input-group">
            <label>手机号</label>
            <input
              v-model="bindForm.mobile"
              type="tel"
              placeholder="请输入手机号"
              maxlength="11"
              :disabled="bindForm.codeSent"
            />
          </div>
          
          <div v-if="!bindForm.codeSent" class="button-group">
            <button 
              @click="sendBindCode" 
              :disabled="!isValidBindMobile || bindForm.sending"
              class="send-btn"
            >
              {{ bindForm.sending ? '发送中...' : '获取验证码' }}
            </button>
          </div>
          
          <div v-if="bindForm.codeSent" class="input-group">
            <label>验证码</label>
            <div class="code-input-wrapper">
              <input
                v-model="bindForm.code"
                type="text"
                placeholder="请输入6位验证码"
                maxlength="6"
              />
              <button 
                @click="resendBindCode" 
                :disabled="bindForm.countdown > 0"
                class="resend-btn"
              >
                {{ bindForm.countdown > 0 ? `${bindForm.countdown}s后重发` : '重新发送' }}
              </button>
            </div>
          </div>
          
          <div v-if="bindForm.codeSent" class="button-group">
            <button 
              @click="handleBindMobile" 
              :disabled="!isValidBindCode || bindForm.binding"
              class="bind-btn"
            >
              {{ bindForm.binding ? '绑定中...' : '确认绑定' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { userAPI } from '@/api/user'

const router = useRouter()

// 响应式数据
const userInfo = ref({
  mobile: '',
  nickname: '',
  username: '',
  avatar: '',
  user_type: ''
})

const stats = ref({
  rides: 25,
  distance: 368,
  points: 1280,
  coupons: 3
})

// 绑定手机号相关
const showBindMobile = ref(false)
const bindForm = ref({
  mobile: '',
  code: '',
  codeSent: false,
  sending: false,
  binding: false,
  countdown: 0
})

let bindCountdownTimer = null

// 计算属性
const isValidBindMobile = computed(() => {
  return /^1[3-9]\d{9}$/.test(bindForm.value.mobile)
})

const isValidBindCode = computed(() => {
  return bindForm.value.code.length === 6
})

// 方法
const loadUserInfo = async () => {
  try {
    // 先从本地存储获取基本信息
    const savedUserInfo = localStorage.getItem('userInfo')
    if (savedUserInfo) {
      const parsed = JSON.parse(savedUserInfo)
      userInfo.value = { ...userInfo.value, ...parsed }
    }

    // 调用API获取最新的用户信息
    const response = await userAPI.auth.getUserInfoList()
    
    console.log('获取用户信息响应:', response.data)
    
    // 修复响应格式处理
    if (response.data.code === 200 || response.data.code === 0) {
      const userData = response.data.data?.user || response.data.data
      if (userData) {
        userInfo.value = { ...userInfo.value, ...userData }
        // 更新本地存储
        localStorage.setItem('userInfo', JSON.stringify(userInfo.value))
      }
    }
  } catch (error) {
    console.error('获取用户信息失败:', error)
    // 如果API调用失败，使用本地存储的信息
    if (error.response?.status === 401) {
      // Token失效，清除本地信息
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      router.push('/login')
    }
  }
}

const goBack = () => {
  router.go(-1)
}

const editProfile = () => {
  alert('编辑资料功能开发中...')
}

const goToOrders = () => {
  alert('我的订单功能开发中...')
}

const goToWallet = () => {
  alert('我的钱包功能开发中...')
}

const goToAddress = () => {
  alert('常用地址功能开发中...')
}

const goToInvite = () => {
  alert('邀请好友功能开发中...')
}

const goToSecurity = () => {
  alert('账号与安全功能开发中...')
}

const goToSettings = () => {
  alert('设置功能开发中...')
}

const goToHelp = () => {
  alert('帮助与反馈功能开发中...')
}

const goToAbout = () => {
  alert('关于我们功能开发中...')
}

const confirmLogout = () => {
  if (confirm('确定要退出登录吗？')) {
    logout()
  }
}

const logout = () => {
  // 清除本地存储
  localStorage.removeItem('token')
  localStorage.removeItem('userInfo')
  
  alert('已退出登录')
  
  // 跳转到首页
  router.push('/')
}

// 绑定手机号相关方法
const closeBindModal = () => {
  showBindMobile.value = false
  bindForm.value = {
    mobile: '',
    code: '',
    codeSent: false,
    sending: false,
    binding: false,
    countdown: 0
  }
  if (bindCountdownTimer) {
    clearInterval(bindCountdownTimer)
    bindCountdownTimer = null
  }
}

const sendBindCode = async () => {
  if (!isValidBindMobile.value) {
    alert('请输入正确的手机号')
    return
  }

  bindForm.value.sending = true
  
  try {
    const response = await userAPI.sendSms(bindForm.value.mobile, 'bind')
    console.log('绑定验证码响应:', response.data)
    
    if (response.data.code === 200 || response.data.code === 0) {
      bindForm.value.codeSent = true
      startBindCountdown()
    } else {
      alert(response.data.message || response.data.msg || '发送验证码失败')
    }
  } catch (error) {
    console.error('发送验证码失败:', error)
    if (error.response?.data) {
      const errorData = error.response.data
      const errorMsg = errorData.message || errorData.msg || errorData.status_msg || '发送验证码失败'
      alert(`发送失败: ${errorMsg}`)
    } else {
      alert('网络连接失败，请检查后端服务状态')
    }
  } finally {
    bindForm.value.sending = false
  }
}

const startBindCountdown = () => {
  bindForm.value.countdown = 60
  bindCountdownTimer = setInterval(() => {
    bindForm.value.countdown--
    if (bindForm.value.countdown <= 0) {
      clearInterval(bindCountdownTimer)
      bindCountdownTimer = null
    }
  }, 1000)
}

const resendBindCode = () => {
  bindForm.value.code = ''
  sendBindCode()
}

const handleBindMobile = async () => {
  if (!isValidBindMobile.value) {
    alert('请输入正确的手机号')
    return
  }

  if (!isValidBindCode.value) {
    alert('请输入6位验证码')
    return
  }

  bindForm.value.binding = true

  try {
    const response = await userAPI.auth.bindMobile(bindForm.value.mobile, bindForm.value.code)
    
    console.log('绑定手机号响应:', response.data)
    
    // 修复响应格式处理
    if (response.data.code === 200 || response.data.code === 0) {
      alert('绑定成功！')
      userInfo.value.mobile = bindForm.value.mobile
      // 更新本地存储
      localStorage.setItem('userInfo', JSON.stringify(userInfo.value))
      closeBindModal()
      // 重新加载用户信息
      loadUserInfo()
    } else {
      const errorMsg = response.data.message || response.data.msg || '绑定失败'
      alert(errorMsg)
    }
  } catch (error) {
    console.error('绑定失败:', error)
    if (error.response?.data) {
      const errorData = error.response.data
      const errorMsg = errorData.message || errorData.msg || errorData.status_msg || '绑定失败'
      alert(`绑定失败: ${errorMsg}`)
    } else {
      alert('网络连接失败，请检查后端服务状态')
    }
  } finally {
    bindForm.value.binding = false
  }
}

// 生命周期
onMounted(() => {
  loadUserInfo()
})
</script>

<style scoped>
.profile {
  min-height: 100vh;
  background: #f5f5f5;
}

.header {
  background: #ff6600;
  color: white;
  padding: 1rem 0;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1rem;
}

.back-btn {
  background: none;
  border: none;
  color: white;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 50%;
  transition: all 0.3s ease;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.back-icon {
  font-size: 1.5rem;
  font-weight: bold;
}

.page-title {
  font-size: 1.2rem;
  font-weight: bold;
}

.header-placeholder {
  width: 40px;
}

.profile-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 1rem;
  gap: 1rem;
  display: flex;
  flex-direction: column;
}

.user-card {
  background: white;
  border-radius: 15px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.user-avatar {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #ff6600, #ff8533);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-icon {
  font-size: 2rem;
}

.user-info {
  flex: 1;
}

.user-name {
  color: #333;
  font-size: 1.2rem;
  font-weight: bold;
  margin-bottom: 0.25rem;
}

.user-desc {
  color: #999;
  font-size: 0.9rem;
}

.edit-btn {
  background: #f5f5f5;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.edit-btn:hover {
  background: #eee;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
  background: white;
  padding: 1.5rem;
  border-radius: 15px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.stat-item {
  text-align: center;
}

.stat-number {
  font-size: 1.5rem;
  font-weight: bold;
  color: #ff6600;
  margin-bottom: 0.25rem;
}

.stat-label {
  color: #666;
  font-size: 0.9rem;
}

.services-section {
  background: white;
  border-radius: 15px;
  padding: 1.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.section-title {
  color: #333;
  font-size: 1.1rem;
  font-weight: bold;
  margin-bottom: 1rem;
}

.services-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
}

.service-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 1rem;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.service-item:hover {
  background: #f8f8f8;
}

.service-icon {
  font-size: 2rem;
  margin-bottom: 0.5rem;
}

.service-text {
  color: #666;
  font-size: 0.9rem;
}

.menu-section {
  background: white;
  border-radius: 15px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 1rem 1.5rem;
  cursor: pointer;
  transition: all 0.3s ease;
  border-bottom: 1px solid #f5f5f5;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-item:hover {
  background: #f8f8f8;
}

.menu-icon {
  font-size: 1.5rem;
  margin-right: 1rem;
}

.menu-text {
  flex: 1;
  color: #333;
  font-weight: 500;
}

.menu-arrow {
  color: #ccc;
  font-size: 1.2rem;
}

.logout-section {
  margin-top: 1rem;
}

.logout-btn {
  width: 100%;
  background: white;
  color: #ff4444;
  border: none;
  padding: 1rem;
  border-radius: 15px;
  font-size: 1rem;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.logout-btn:hover {
  background: #fff5f5;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .services-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .profile-container {
    padding: 1rem 0.5rem;
  }
  
  .user-card {
    padding: 1rem;
  }
  
  .stats-grid, .services-section, .menu-section {
    padding: 1rem;
  }
}

@media (max-width: 480px) {
  .header-content {
    padding: 0 0.5rem;
  }
  
  .user-name {
    font-size: 1rem;
  }
  
  .stat-number {
    font-size: 1.2rem;
  }
  
  .service-icon {
    font-size: 1.5rem;
  }
  
  .menu-icon {
    font-size: 1.2rem;
  }
}

/* 绑定手机号弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.bind-modal {
  background: white;
  border-radius: 15px;
  width: 100%;
  max-width: 400px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.5rem;
  border-bottom: 1px solid #f0f0f0;
}

.modal-header h3 {
  margin: 0;
  color: #333;
  font-size: 1.2rem;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #999;
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 50%;
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: #f0f0f0;
  color: #666;
}

.modal-body {
  padding: 1.5rem;
}

.input-group {
  margin-bottom: 1.5rem;
}

.input-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #333;
  font-weight: 600;
  font-size: 0.9rem;
}

.input-group input {
  width: 100%;
  padding: 1rem;
  border: 2px solid #eee;
  border-radius: 10px;
  font-size: 1rem;
  transition: border-color 0.3s ease;
  box-sizing: border-box;
}

.input-group input:focus {
  outline: none;
  border-color: #ff6600;
}

.input-group input:disabled {
  background: #f5f5f5;
  color: #999;
}

.code-input-wrapper {
  display: flex;
  gap: 0.5rem;
}

.code-input-wrapper input {
  flex: 1;
}

.resend-btn {
  background: #f5f5f5;
  color: #666;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s ease;
  white-space: nowrap;
}

.resend-btn:hover:not(:disabled) {
  background: #eee;
}

.resend-btn:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.button-group {
  margin-top: 1rem;
}

.send-btn, .bind-btn {
  width: 100%;
  background: #ff6600;
  color: white;
  border: none;
  padding: 1rem;
  border-radius: 10px;
  font-size: 1rem;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s ease;
}

.send-btn:hover:not(:disabled), .bind-btn:hover:not(:disabled) {
  background: #e55a00;
  transform: translateY(-2px);
}

.send-btn:disabled, .bind-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
}
</style> 