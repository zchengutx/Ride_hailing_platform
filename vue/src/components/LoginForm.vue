<template>
  <div class="login-form">
    <div class="login-tabs">
      <div 
        class="tab"
        :class="{ active: loginType === 'sms' }"
        @click="loginType = 'sms'"
      >
        验证码登录
      </div>
      <div 
        class="tab"
        :class="{ active: loginType === 'wechat' }"
        @click="loginType = 'wechat'"
      >
        微信登录
      </div>
    </div>

    <!-- 短信登录 -->
    <div v-if="loginType === 'sms'" class="sms-login">
      <el-form :model="smsForm" :rules="smsRules" ref="smsFormRef">
        <el-form-item prop="phone">
          <el-input
            v-model="smsForm.phone"
            placeholder="请输入手机号"
            size="large"
            :prefix-icon="Phone"
          />
        </el-form-item>
        
        <el-form-item prop="code">
          <div class="code-input">
            <el-input
              v-model="smsForm.code"
              placeholder="请输入验证码"
              size="large"
              :prefix-icon="Message"
            />
            <el-button
              class="send-code-btn"
              :disabled="codeSending || countdown > 0"
              @click="sendCode"
            >
              {{ countdown > 0 ? `${countdown}s` : '获取验证码' }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>

      <el-button
        type="primary"
        size="large"
        class="login-btn"
        :loading="logging"
        @click="handleSmsLogin"
      >
        登录
      </el-button>
    </div>

    <!-- 微信登录 -->
    <div v-if="loginType === 'wechat'" class="wechat-login">
      <div class="wechat-qr">
        <div v-if="!qrCodeImage" class="qr-placeholder">
          <el-icon size="80" color="#07c160"><MessageBox /></el-icon>
          <p>请使用微信扫码登录</p>
        </div>
        <div v-else class="qr-code">
          <div class="qr-code-wrapper" :class="{ expired: qrCodeStatus === 'expired' }">
            <img :src="qrCodeImage" alt="微信登录二维码" />
            
            <!-- 扫码状态遮罩 -->
            <div v-if="qrCodeStatus === 'scanning'" class="qr-status-overlay scanning">
              <el-icon size="40" color="#07c160"><MessageBox /></el-icon>
              <p>检测到扫码</p>
              <p>请在手机上确认登录</p>
            </div>
            
            <div v-if="qrCodeStatus === 'success'" class="qr-status-overlay success">
              <el-icon size="40" color="#67c23a"><SuccessFilled /></el-icon>
              <p>登录成功！</p>
            </div>
            
            <div v-if="qrCodeStatus === 'expired'" class="qr-status-overlay expired">
              <el-icon size="40" color="#f56c6c"><CircleClose /></el-icon>
              <p>二维码已过期</p>
              <p>请点击刷新</p>
            </div>
            
            <div v-if="qrCodeStatus === 'cancelled'" class="qr-status-overlay cancelled">
              <el-icon size="40" color="#909399"><Close /></el-icon>
              <p>用户取消登录</p>
            </div>
          </div>
          
          <!-- 状态提示文字 -->
          <p class="qr-status-text" :class="qrCodeStatus">
            <template v-if="qrCodeStatus === 'waiting'">
              请使用微信扫描上方二维码
            </template>
            <template v-else-if="qrCodeStatus === 'scanning'">
              已扫码，请在手机上确认登录
            </template>
            <template v-else-if="qrCodeStatus === 'success'">
              登录成功，正在跳转...
            </template>
            <template v-else-if="qrCodeStatus === 'expired'">
              二维码已过期，请重新获取
            </template>
            <template v-else-if="qrCodeStatus === 'cancelled'">
              登录已取消
            </template>
          </p>
          
          <!-- 开发提示 -->
          <div class="dev-tip">
            <p>✅ 演示版本提示：</p>
            <p>已禁用状态轮询，避免404错误</p>
            <p>请点击下方"模拟登录成功"按钮测试完整流程</p>
          </div>
        </div>
      </div>
      
      <el-button
        type="success"
        size="large"
        class="wechat-btn"
        :loading="wechatLogging"
        :disabled="qrCodeStatus === 'success'"
        @click="handleWechatLogin"
      >
        <el-icon><MessageBox /></el-icon>
        {{ getWechatButtonText() }}
      </el-button>
      
      <!-- 开发测试按钮 -->
      <el-button
        v-if="qrCodeImage"
        type="primary"
        size="large"
        class="mock-success-btn"
        @click="mockLoginSuccess"
      >
        🧪 模拟登录成功（仅测试）
      </el-button>
    </div>

    <div class="login-footer">
      <p>登录即表示同意 <a href="#">《用户协议》</a> 和 <a href="#">《隐私政策》</a></p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onUnmounted, watch } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import { Phone, Message, MessageBox, SuccessFilled, CircleClose, Close } from '@element-plus/icons-vue'
import { sendSms, login, wechatLogin, checkWechatLoginStatus } from '@/api/index'

const emit = defineEmits(['login-success'])

const loginType = ref<'sms' | 'wechat'>('sms')
const logging = ref(false)
const wechatLogging = ref(false)
const codeSending = ref(false)
const countdown = ref(0)
const qrCodeImage = ref('')

// 微信扫码状态
const qrCodeStatus = ref<'waiting' | 'scanning' | 'success' | 'expired' | 'cancelled'>('waiting')
const qrCodeKey = ref('') // 二维码的唯一标识
const statusCheckTimer = ref<number | null>(null)
const qrExpireTimer = ref<number | null>(null) // 二维码过期定时器

const smsFormRef = ref<FormInstance>()
const smsForm = reactive({
  phone: '',
  code: ''
})

const smsRules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码为6位数字', trigger: 'blur' }
  ]
}

let countdownTimer: number | null = null

const sendCode = async () => {
  if (!smsForm.phone) {
    ElMessage.warning('请先输入手机号')
    return
  }

  if (!/^1[3-9]\d{9}$/.test(smsForm.phone)) {
    ElMessage.warning('请输入正确的手机号')
    return
  }

  try {
    codeSending.value = true
    await sendSms(smsForm.phone)
    ElMessage.success('验证码已发送')
    
    // 开始倒计时
    countdown.value = 60
    countdownTimer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(countdownTimer!)
        countdownTimer = null
      }
    }, 1000)
  } catch (error) {
    ElMessage.error('发送验证码失败')
  } finally {
    codeSending.value = false
  }
}

const handleSmsLogin = async () => {
  if (!smsFormRef.value) return

  try {
    await smsFormRef.value.validate()
    logging.value = true
    
    const result = await login({
      phone: smsForm.phone,
      code: smsForm.code
    })
    
    if (result.code === 200) {
      // 保存token - 优先保存到localStorage，同时保存到cookie
      const token = result.data || result.token || ''
      if (token) {
        localStorage.setItem('token', token)
        // 同步保存到cookie（7天过期）
        document.cookie = `x-token=${token}; expires=${new Date(Date.now() + 7*24*60*60*1000).toUTCString()}; path=/`
        console.log('💾 登录成功，token已保存:', token.substring(0, 20) + '...')
      }
      
      ElMessage.success('登录成功，正在跳转...')
      
      // 发送登录成功事件给父组件
      emit('login-success')
      
      // 延迟一下再跳转，让用户看到成功消息
      setTimeout(() => {
        // 跳转到首页而不是重新加载页面
        window.location.href = '/'
      }, 500)
    } else {
      ElMessage.error(result.message || '登录失败')
    }
  } catch (error) {
    ElMessage.error('登录失败')
  } finally {
    logging.value = false
  }
}

const handleWechatLogin = async () => {
  try {
    wechatLogging.value = true
    
    // 清除之前的轮询和过期定时器
    clearAllTimers()
    
    const result = await wechatLogin()
    
    if (result.code === 200) {
      // 处理返回的二维码数据
      if (result.data) {
        // 假设后端返回格式：{ qrCode: "base64图片", qrKey: "唯一标识" }
        if (typeof result.data === 'object' && result.data.qrCode && result.data.qrKey) {
          qrCodeImage.value = `data:image/png;base64,${result.data.qrCode}`
          qrCodeKey.value = result.data.qrKey
          qrCodeStatus.value = 'waiting'
          
          ElMessage.success('二维码获取成功，请使用微信扫码')
          
          // 暂时禁用轮询检查，避免404错误（等待后端接口开发完成）
          // startStatusCheck()
        } else if (typeof result.data === 'string') {
          // 兼容旧版本：如果只返回图片base64字符串
          qrCodeImage.value = `data:image/png;base64,${result.data}`
          qrCodeKey.value = Date.now().toString() // 临时使用时间戳作为key
          qrCodeStatus.value = 'waiting'
          
          ElMessage.success('二维码获取成功，请使用微信扫码')
          
          // 暂时禁用轮询检查，避免404错误（等待后端接口开发完成）
          // startStatusCheck()
        } else {
          ElMessage.error('二维码数据格式不正确')
        }
      } else {
        ElMessage.error('未获取到二维码数据')
      }
    } else {
      ElMessage.error(result.message || '获取微信二维码失败')
    }
  } catch (error) {
    console.error('微信登录错误:', error)
    ElMessage.error('获取微信二维码失败')
  } finally {
    wechatLogging.value = false
  }
}

// 开始轮询检查扫码状态
const startStatusCheck = () => {
  if (!qrCodeKey.value) return
  
  // 设置二维码过期定时器（3分钟）
  qrExpireTimer.value = setTimeout(() => {
    if (qrCodeStatus.value === 'waiting' || qrCodeStatus.value === 'scanning') {
      qrCodeStatus.value = 'expired'
      
      // 停止轮询
      if (statusCheckTimer.value) {
        clearInterval(statusCheckTimer.value)
        statusCheckTimer.value = null
      }
      
      ElMessage.warning('二维码已过期，请重新获取')
    }
  }, 3 * 60 * 1000) // 3分钟过期
  
  statusCheckTimer.value = setInterval(async () => {
    try {
      const result = await checkWechatLoginStatus(qrCodeKey.value)
      
      if (result.code === 200) {
        const status = result.data?.status
        
        switch (status) {
          case 'waiting':
            qrCodeStatus.value = 'waiting'
            break
            
          case 'scanning':
            qrCodeStatus.value = 'scanning'
            ElMessage.info('检测到微信扫码，请在手机上确认登录')
            break
            
          case 'success':
            qrCodeStatus.value = 'success'
            
            // 停止轮询和过期定时器
            if (statusCheckTimer.value) {
              clearInterval(statusCheckTimer.value)
              statusCheckTimer.value = null
            }
            if (qrExpireTimer.value) {
              clearTimeout(qrExpireTimer.value)
              qrExpireTimer.value = null
            }
            
            // 保存token
            if (result.data?.token) {
              const token = result.data.token
              localStorage.setItem('token', token)
              // 同步保存到cookie
              document.cookie = `x-token=${token}; expires=${new Date(Date.now() + 7*24*60*60*1000).toUTCString()}; path=/`
              console.log('💾 微信登录成功，token已保存:', token.substring(0, 20) + '...')
              
              ElMessage.success('微信登录成功，正在跳转...')
              emit('login-success')
              
              // 延迟跳转到首页
              setTimeout(() => {
                window.location.href = '/'
              }, 500)
            } else {
              ElMessage.success('微信登录成功！')
              emit('login-success')
            }
            break
            
          case 'expired':
            qrCodeStatus.value = 'expired'
            
            // 停止轮询
            if (statusCheckTimer.value) {
              clearInterval(statusCheckTimer.value)
              statusCheckTimer.value = null
            }
            
            ElMessage.warning('二维码已过期，请重新获取')
            break
            
          case 'cancelled':
            qrCodeStatus.value = 'cancelled'
            
            // 停止轮询
            if (statusCheckTimer.value) {
              clearInterval(statusCheckTimer.value)
              statusCheckTimer.value = null
            }
            
            ElMessage.info('用户取消了登录')
            break
        }
      }
    } catch (error) {
      // 如果是404错误，说明后端接口还未实现，静默处理
      if (error.response?.status === 404) {
        console.log('扫码状态检查接口暂未实现，请手动刷新或等待后端接口开发完成')
        return
      }
      
      // 其他错误才打印到控制台
      console.warn('检查扫码状态失败:', error)
    }
  }, 3000) // 调整为每3秒检查一次，减少请求频率
}

// 获取微信按钮文字
const getWechatButtonText = () => {
  if (qrCodeStatus.value === 'success') {
    return '登录成功'
  } else if (qrCodeStatus.value === 'expired') {
    return '重新获取二维码'
  } else if (qrCodeImage.value) {
    return '刷新二维码'
  } else {
    return '获取微信二维码'
  }
}

// 清理所有定时器的统一函数
const clearAllTimers = () => {
  if (statusCheckTimer.value) {
    clearInterval(statusCheckTimer.value)
    statusCheckTimer.value = null
  }
  if (qrExpireTimer.value) {
    clearTimeout(qrExpireTimer.value)
    qrExpireTimer.value = null
  }
}

// 模拟登录成功（仅用于测试）
const mockLoginSuccess = () => {
  qrCodeStatus.value = 'success'
  
  // 停止轮询和过期定时器
  clearAllTimers()
  
  // 模拟保存token
  const mockToken = 'mock_token_' + Date.now()
  localStorage.setItem('token', mockToken)
  // 同步保存到cookie
  document.cookie = `x-token=${mockToken}; expires=${new Date(Date.now() + 7*24*60*60*1000).toUTCString()}; path=/`
  
  // 显示成功提示
  ElMessage.success('模拟登录成功，正在跳转...')
  
  // 通知父组件登录成功
  emit('login-success')
  
  // 延迟跳转到首页
  setTimeout(() => {
    window.location.href = '/'
  }, 500)
}

// 监听登录类型变化，切换到微信登录时自动获取二维码
watch(loginType, (newType) => {
  if (newType === 'wechat' && !qrCodeImage.value) {
    handleWechatLogin()
  }
  
  // 切换登录类型时清理可能存在的定时器
  clearAllTimers()
})

// 组件销毁时清理定时器
onUnmounted(() => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
  }
  clearAllTimers()
})
</script>

<style scoped>
.login-form {
  width: 100%;
}

.login-tabs {
  display: flex;
  margin-bottom: 24px;
  border-bottom: 1px solid #eee;
}

.tab {
  flex: 1;
  text-align: center;
  padding: 12px 0;
  cursor: pointer;
  color: #666;
  font-size: 14px;
  transition: all 0.3s;
}

.tab.active {
  color: #ff7e00;
  border-bottom: 2px solid #ff7e00;
}

.sms-login {
  margin-bottom: 20px;
}

.code-input {
  display: flex;
  gap: 8px;
}

.code-input .el-input {
  flex: 1;
}

.send-code-btn {
  width: 100px;
  font-size: 12px;
}

.login-btn, .wechat-btn {
  width: 100%;
  height: 44px;
  border-radius: 22px;
  font-size: 16px;
  font-weight: 500;
}

.login-btn {
  background: linear-gradient(135deg, #ff7e00 0%, #ff6600 100%);
  border: none;
}

.wechat-login {
  text-align: center;
  margin-bottom: 20px;
}

.wechat-qr {
  padding: 20px;
  margin-bottom: 20px;
}

.wechat-qr p {
  margin: 12px 0 0 0;
  color: #666;
  font-size: 14px;
}

.qr-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.qr-code {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.qr-code-wrapper {
  position: relative;
  display: inline-block;
}

.qr-code-wrapper.expired {
  opacity: 0.5;
}

.qr-code img {
  width: 200px;
  height: 200px;
  border: 1px solid #eee;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: opacity 0.3s;
}

.qr-status-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  animation: fadeIn 0.3s;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.qr-status-overlay p {
  margin: 4px 0;
  font-size: 14px;
  font-weight: 500;
}

.qr-status-overlay.scanning p {
  color: #07c160;
}

.qr-status-overlay.success p {
  color: #67c23a;
}

.qr-status-overlay.expired p {
  color: #f56c6c;
}

.qr-status-overlay.cancelled p {
  color: #909399;
}

.qr-status-text {
  margin: 12px 0 0 0;
  color: #666;
  font-size: 14px;
  transition: color 0.3s;
}

.qr-status-text.scanning {
  color: #07c160;
  animation: pulse 2s infinite;
}

.qr-status-text.success {
  color: #67c23a;
}

.qr-status-text.expired {
  color: #f56c6c;
}

.qr-status-text.cancelled {
  color: #909399;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}

.dev-tip {
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  border-radius: 6px;
  padding: 8px 12px;
  margin-top: 12px;
  font-size: 11px;
  line-height: 1.4;
}

.dev-tip p {
  margin: 2px 0;
  color: #166534;
}

.dev-tip p:first-child {
  font-weight: 600;
}

.mock-success-btn {
  width: 100%;
  height: 40px;
  margin-top: 12px;
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  border: none;
  border-radius: 20px;
  font-size: 14px;
  opacity: 0.8;
}

.mock-success-btn:hover {
  opacity: 1;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.wechat-btn {
  background: #07c160;
  border-color: #07c160;
}

.wechat-btn:hover {
  background: #06ad56;
  border-color: #06ad56;
}

.login-footer {
  margin-top: 20px;
  text-align: center;
}

.login-footer p {
  font-size: 12px;
  color: #999;
  margin: 0;
  line-height: 1.4;
}

.login-footer a {
  color: #ff7e00;
  text-decoration: none;
}

.login-footer a:hover {
  text-decoration: underline;
}
</style> 