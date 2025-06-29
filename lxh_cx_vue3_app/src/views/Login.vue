<template>
  <div class="login">
    <div class="login-container">
      <div class="login-card">
        <div class="logo">
          <div class="logo-icon">🐱</div>
          <h1>罗小黑出行</h1>
          <p>欢迎回来</p>
        </div>

        <!-- 登录方式切换 -->
        <div class="login-tabs">
          <div 
            :class="['tab-item', { active: loginType === 'phone' }]"
            @click="switchLoginType('phone')"
          >
            手机登录
          </div>
          <div 
            :class="['tab-item', { active: loginType === 'wechat' }]"
            @click="switchLoginType('wechat')"
          >
            微信登录
          </div>
        </div>

        <!-- 手机登录 -->
        <div v-if="loginType === 'phone'" class="phone-login">
          <form @submit.prevent="handlePhoneLogin" class="login-form">
            <div class="input-group">
              <label for="phone">手机号</label>
              <input
                id="phone"
                v-model="phoneForm.phone"
                type="tel"
                placeholder="请输入手机号"
                maxlength="11"
                class="phone-input"
                :disabled="showCodeInput"
              />
            </div>

            <div v-if="!showCodeInput" class="button-group">
              <button 
                type="button" 
                @click="sendCode" 
                :disabled="!isValidPhone || codeSending"
                class="send-code-btn"
              >
                {{ codeSending ? '发送中...' : '获取验证码' }}
              </button>
            </div>

            <div v-if="showCodeInput" class="input-group">
              <label for="code">请输入验证码</label>
              <p class="code-tip">验证码已发送至 {{ formatPhoneNumber(phoneForm.phone) }}</p>
              
              <!-- 6位验证码输入框 -->
              <div class="verification-code-input">
                <input
                  v-for="(digit, index) in codeDigits"
                  :key="index"
                  :ref="(el) => codeInputRefs[index] = el"
                  v-model="codeDigits[index]"
                  type="text"
                  maxlength="1"
                  class="code-digit"
                  :class="{ 'active': activeDigitIndex === index, 'filled': codeDigits[index] }"
                  @input="handleDigitInput(index, $event)"
                  @keydown="handleDigitKeydown(index, $event)"
                  @focus="activeDigitIndex = index"
                  @blur="activeDigitIndex = -1"
                />
              </div>
              
              <div class="resend-section">
                <span v-if="countdown > 0" class="countdown-text">
                  {{ countdown }}s后可重新发送
                </span>
                <button 
                  v-else
                  type="button" 
                  @click="resendCode" 
                  class="resend-link"
                >
                  重新发送验证码
                </button>
              </div>
            </div>

            <div v-if="showCodeInput" class="button-group">
              <button 
                type="submit" 
                :disabled="!isValidCode || loggingIn"
                class="login-btn"
              >
                {{ loggingIn ? '登录中...' : '登录' }}
              </button>
            </div>
          </form>
        </div>

        <!-- 微信登录 -->
        <div v-if="loginType === 'wechat'" class="wechat-login">
          <div class="qr-container">
            <!-- 加载中状态 -->
            <div v-if="wechatLoading" class="qr-loading">
              <div class="loading-spinner"></div>
              <p>正在获取微信二维码...</p>
            </div>
            
            <!-- 二维码显示 -->
            <div v-else-if="wechatQRCode && !wechatLoginSuccess" class="qr-code">
              <img :src="wechatQRCode" alt="微信登录二维码" />
              <div class="qr-refresh">
                <button @click="refreshWeChatQRCode" class="refresh-btn">
                  🔄 刷新二维码
                </button>
              </div>
            </div>
            
            <!-- 登录成功状态 -->
            <div v-else-if="wechatLoginSuccess" class="qr-success">
              <div class="success-icon">✅</div>
              <p>扫码成功，正在登录...</p>
            </div>
            
            <!-- 默认占位符 -->
            <div v-else class="qr-placeholder">
              <div class="qr-icon">📱</div>
              <p>请使用微信扫码登录</p>
              <button @click="loadWeChatQRCode" class="retry-btn">
                重新获取二维码
              </button>
            </div>
          </div>
          <div class="wechat-tips">
            <p>📱 打开微信，扫一扫二维码即可登录</p>
            <p class="tip-small">⏰ 二维码5分钟内有效</p>
            <p class="tip-warning">⚠️ 扫码后请手动刷新页面完成登录</p>
          </div>
        </div>

        <!-- 底部链接 -->
        <div class="bottom-links">
          <router-link to="/" class="back-home">返回首页</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onUnmounted, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { userApi } from '@/api/user'

const router = useRouter()
const route = useRoute()

// 响应式数据
const loginType = ref('phone')
const showCodeInput = ref(false)
const codeSending = ref(false)
const loggingIn = ref(false)
const countdown = ref(0)
let countdownTimer = null

const phoneForm = ref({
  phone: '',
  code: ''
})

// 验证码输入相关
const codeDigits = ref(['', '', '', '', '', ''])
const codeInputRefs = ref([])
const activeDigitIndex = ref(-1)

// 微信登录相关
const wechatQRCode = ref('')
const wechatLoading = ref(false)
const wechatPolling = ref(null)
const wechatLoginSuccess = ref(false)

// 计算属性
const isValidPhone = computed(() => {
  return /^1[3-9]\d{9}$/.test(phoneForm.value.phone)
})

const isValidCode = computed(() => {
  return codeDigits.value.every(digit => digit !== '') && codeDigits.value.join('').length === 6
})

// 方法
// 发送验证码
const sendCode = async () => {
  if (!isValidPhone.value) {
    alert('请输入正确的手机号')
    return
  }

  codeSending.value = true
  
  try {
    const response = await userApi.sendSms({
      mobile: phoneForm.value.phone,
      source: 'login'  // 使用正确的参数名source
    })

    console.log('发送验证码响应:', response)

    if (response.code === 200 || response.code === 0) {
      showCodeInput.value = true
      startCountdown()
      // 清空验证码输入框
      codeDigits.value = ['', '', '', '', '', '']
      // 聚焦到第一个输入框
      setTimeout(() => {
        if (codeInputRefs.value[0]) {
          codeInputRefs.value[0].focus()
        }
      }, 100)
    } else {
      alert(response.message || response.msg || '发送验证码失败')
    }
  } catch (error) {
    console.error('发送验证码失败:', error)
    alert('发送验证码失败，请重试')
  } finally {
    codeSending.value = false
  }
}

const resendCode = () => {
  codeDigits.value = ['', '', '', '', '', '']
  phoneForm.value.code = ''
  sendCode()
}

const startCountdown = () => {
  countdown.value = 60
  countdownTimer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      clearInterval(countdownTimer)
    }
  }, 1000)
}

// 处理单个验证码输入
const handleDigitInput = (index, event) => {
  const value = event.target.value
  
  // 只允许数字
  if (!/^\d$/.test(value)) {
    codeDigits.value[index] = ''
    return
  }
  
  codeDigits.value[index] = value
  
  // 自动跳转到下一个输入框
  if (value && index < 5) {
    const nextInput = codeInputRefs.value[index + 1]
    if (nextInput) {
      nextInput.focus()
    }
  }
  
  // 如果6位都填完了，自动登录
  if (isValidCode.value) {
    phoneForm.value.code = codeDigits.value.join('')
    setTimeout(() => {
      handlePhoneLogin()
    }, 300)
  }
}

// 处理键盘事件
const handleDigitKeydown = (index, event) => {
  // 退格键处理
  if (event.key === 'Backspace') {
    if (!codeDigits.value[index] && index > 0) {
      // 当前框为空，删除前一个框的内容并聚焦
      const prevInput = codeInputRefs.value[index - 1]
      if (prevInput) {
        codeDigits.value[index - 1] = ''
        prevInput.focus()
      }
    } else {
      // 当前框有内容，清空当前框
      codeDigits.value[index] = ''
    }
  }
  // 左右箭头键处理
  else if (event.key === 'ArrowLeft' && index > 0) {
    codeInputRefs.value[index - 1].focus()
  }
  else if (event.key === 'ArrowRight' && index < 5) {
    codeInputRefs.value[index + 1].focus()
  }
}

// 格式化手机号显示
const formatPhoneNumber = (phone) => {
  if (!phone || phone.length !== 11) return phone
  return phone.replace(/(\d{3})(\d{4})(\d{4})/, '$1****$3')
}

// 处理登录
const handlePhoneLogin = async () => {
  if (!isValidPhone.value) {
    alert('请输入正确的手机号')
    return
  }

  if (!isValidCode.value) {
    alert('请输入6位验证码')
    return
  }

  loggingIn.value = true

  try {
    const response = await userApi.login({
      mobile: phoneForm.value.phone,
      sendCode: codeDigits.value.join('')  // 使用正确的参数名sendCode
    })

    console.log('登录响应:', response)

    if (response.code === 200 || response.code === 0) {
      if (response.data?.token) {
        localStorage.setItem('token', response.data.token)
        alert('登录成功！')
        router.push('/')
      } else {
        alert('登录失败：未获取到有效token')
      }
    } else {
      alert(response.message || response.msg || '登录失败')
    }
  } catch (error) {
    console.error('登录失败:', error)
    alert('登录失败，请重试')
  } finally {
    loggingIn.value = false
  }
}

// 微信登录方法
const loadWeChatQRCode = async () => {
  wechatLoading.value = true
  try {
    const response = await userApi.wechat.getQRCode()
    
    if (response.code === 200 || response.code === 0) {
      wechatQRCode.value = response.data.qrcode
      startWeChatPolling()
    } else {
      alert('获取微信二维码失败，请稍后重试')
    }
  } catch (error) {
    console.error('获取微信二维码失败:', error)
    alert('获取微信二维码失败，请稍后重试')
  } finally {
    wechatLoading.value = false
  }
}

const startWeChatPolling = () => {
  // 清除之前的轮询
  if (wechatPolling.value) {
    clearInterval(wechatPolling.value)
  }
  
  // 暂时禁用自动轮询功能，避免403错误
  console.log('微信登录提示：扫码后请手动刷新页面完成登录')
  
  // 显示提示信息给用户
  setTimeout(() => {
    if (!wechatLoginSuccess.value) {
      console.log('微信登录提示：如果您已完成扫码，请刷新页面查看登录状态')
    }
  }, 10000) // 10秒后提示
  
  // 原来的轮询代码暂时注释
  /*
  wechatPolling.value = pollWeChatLoginStatus((loginData) => {
    wechatLoginSuccess.value = true
    
    // 保存登录信息
    localStorage.setItem('token', loginData.token)
    localStorage.setItem('userInfo', JSON.stringify(loginData))
    
    alert('微信登录成功！')
    router.push('/')
  })
  
  // 5分钟后停止轮询
  setTimeout(() => {
    if (wechatPolling.value) {
      clearInterval(wechatPolling.value)
      wechatPolling.value = null
    }
  }, 5 * 60 * 1000)
  */
}

const refreshWeChatQRCode = () => {
  // 清理之前的URL对象
  if (wechatQRCode.value && wechatQRCode.value.startsWith('blob:')) {
    URL.revokeObjectURL(wechatQRCode.value)
  }
  wechatQRCode.value = ''
  wechatLoginSuccess.value = false
  loadWeChatQRCode()
}

// 组件挂载时加载微信二维码
onMounted(() => {
  if (loginType.value === 'wechat') {
    loadWeChatQRCode()
  }
})

// 监听登录方式切换
const switchLoginType = (type) => {
  loginType.value = type
  if (type === 'wechat') {
    loadWeChatQRCode()
  } else {
    // 切换到手机登录时，停止微信轮询
    if (wechatPolling.value) {
      clearInterval(wechatPolling.value)
      wechatPolling.value = null
    }
  }
}

// 清理定时器
onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
  }
  if (wechatPolling.value) {
    clearInterval(wechatPolling.value)
  }
})
</script>

<style scoped>
.login {
  min-height: 100vh;
  background: linear-gradient(135deg, #ff6600 0%, #ff8533 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.login-container {
  width: 100%;
  max-width: 400px;
}

.login-card {
  background: white;
  border-radius: 20px;
  padding: 2.5rem;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
}

.logo {
  text-align: center;
  margin-bottom: 2rem;
}

.logo-icon {
  font-size: 48px;
  margin-bottom: 10px;
}

.logo h1 {
  font-size: 24px;
  margin: 10px 0;
  color: #333;
}

.logo p {
  color: #666;
  font-size: 14px;
}

.login-tabs {
  display: flex;
  background: #f5f5f5;
  border-radius: 10px;
  padding: 0.25rem;
  margin-bottom: 2rem;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 0.75rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  font-weight: 500;
}

.tab-item.active {
  background: white;
  color: #ff6600;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.phone-login {
  margin-bottom: 2rem;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.input-group label {
  color: #333;
  font-weight: 600;
  font-size: 0.9rem;
}

.phone-input, .code-input {
  padding: 1rem;
  border: 2px solid #eee;
  border-radius: 10px;
  font-size: 1rem;
  transition: border-color 0.3s ease;
}

.phone-input:focus, .code-input:focus {
  outline: none;
  border-color: #ff6600;
}

.phone-input:disabled {
  background: #f5f5f5;
  color: #999;
}

.code-input-wrapper {
  display: flex;
  gap: 0.5rem;
}

.code-input {
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
  margin-top: 0.5rem;
}

.send-code-btn, .login-btn {
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

.send-code-btn:hover:not(:disabled), .login-btn:hover:not(:disabled) {
  background: #e55a00;
  transform: translateY(-2px);
}

.send-code-btn:disabled, .login-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
}

.wechat-login {
  text-align: center;
  margin-bottom: 2rem;
}

.qr-container {
  margin-bottom: 1rem;
}

.qr-placeholder {
  width: 200px;
  height: 200px;
  border: 2px dashed #ddd;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  color: #999;
}

.qr-loading {
  width: 200px;
  height: 200px;
  border: 2px solid #f0f0f0;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  color: #666;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #ff6600;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.qr-code {
  text-align: center;
  margin: 0 auto;
}

.qr-code img {
  width: 200px;
  height: 200px;
  border: 2px solid #eee;
  border-radius: 10px;
  object-fit: contain;
}

.qr-refresh {
  margin-top: 1rem;
}

.refresh-btn {
  background: #f5f5f5;
  color: #666;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s ease;
}

.refresh-btn:hover {
  background: #eee;
  color: #333;
}

.qr-success {
  width: 200px;
  height: 200px;
  border: 2px solid #00C896;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  color: #00C896;
  background: #f0f9f7;
}

.success-icon {
  font-size: 3rem;
  margin-bottom: 0.5rem;
}

.retry-btn {
  background: #ff6600;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  margin-top: 1rem;
  transition: all 0.3s ease;
}

.retry-btn:hover {
  background: #e55a00;
}

.qr-icon {
  font-size: 3rem;
  margin-bottom: 0.5rem;
}

.wechat-tips {
  color: #666;
  font-size: 0.9rem;
  text-align: center;
}

.tip-small {
  font-size: 0.8rem;
  color: #999;
  margin-top: 0.5rem;
}

.tip-warning {
  font-size: 0.9rem;
  color: #ff6600;
  font-weight: 500;
  margin-top: 0.5rem;
}

/* 验证码输入相关样式 */
.code-tip {
  color: #666;
  font-size: 0.9rem;
  margin: 0.5rem 0 1.5rem 0;
  text-align: center;
}

.verification-code-input {
  display: flex;
  justify-content: space-between;
  gap: 0.75rem;
  margin: 1.5rem 0;
}

.code-digit {
  width: 3rem;
  height: 3rem;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  text-align: center;
  font-size: 1.5rem;
  font-weight: bold;
  color: #333;
  outline: none;
  transition: all 0.3s ease;
  background: #fafafa;
}

.code-digit:focus,
.code-digit.active {
  border-color: #ff6600;
  background: white;
  box-shadow: 0 0 0 3px rgba(255, 102, 0, 0.1);
  transform: scale(1.05);
}

.code-digit.filled {
  border-color: #ff6600;
  background: #fff5f0;
  color: #ff6600;
}

.code-digit:hover {
  border-color: #ff8533;
}

.resend-section {
  text-align: center;
  margin-top: 1.5rem;
}

.countdown-text {
  color: #999;
  font-size: 0.9rem;
}

.resend-link {
  background: none;
  border: none;
  color: #ff6600;
  font-size: 0.9rem;
  cursor: pointer;
  text-decoration: underline;
  transition: color 0.3s ease;
}

.resend-link:hover {
  color: #e55a00;
}

.bottom-links {
  text-align: center;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}

.back-home {
  color: #ff6600;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s ease;
}

.back-home:hover {
  text-decoration: underline;
}

@media (max-width: 480px) {
  .login {
    padding: 1rem;
  }
  
  .login-card {
    padding: 2rem;
  }
  
  .logo {
    font-size: 1.5rem;
  }
  
  .qr-placeholder {
    width: 150px;
    height: 150px;
  }
  
  /* 移动端验证码输入框调整 */
  .verification-code-input {
    gap: 0.5rem;
  }
  
  .code-digit {
    width: 2.5rem;
    height: 2.5rem;
    font-size: 1.2rem;
  }
  
  .code-tip {
    font-size: 0.8rem;
  }
}
</style>