<template>
  <div class="user-profile">
    <!-- 顶部用户信息卡片 -->
    <div class="user-header">
      <div class="user-info">
        <div class="avatar">
          <img :src="defaultAvatar" alt="头像" />
          <div class="edit-badge">编辑</div>
        </div>
        <div class="user-details">
          <h2 class="username">{{ userInfo.nickname || '未设置昵称' }}</h2>
          <p class="phone">{{ formatPhone(userInfo.phone) }}</p>
        </div>
      </div>
      <div class="user-stats">
        <div class="stat-item">
          <span class="number">{{ userInfo.tripCount || 0 }}</span>
          <span class="label">行程</span>
        </div>
        <div class="stat-item">
          <span class="number">{{ userInfo.favoriteCount || 0 }}</span>
          <span class="label">收藏</span>
        </div>
      </div>
    </div>

    <!-- 钱包卡片 -->
    <div class="wallet-card">
      <div class="balance-section">
        <div class="balance-info">
          <span class="label">钱包余额</span>
          <span class="amount">¥{{ userInfo.balance || '0.00' }}</span>
        </div>
        <button class="recharge-btn">充值</button>
      </div>
      <div class="wallet-items">
        <div class="wallet-item">
          <span class="icon">🎫</span>
          <span class="text">优惠券</span>
          <span class="count">{{ userInfo.couponCount || 0 }}张</span>
        </div>
        <div class="wallet-item">
          <span class="icon">🎁</span>
          <span class="text">礼品卡</span>
          <span class="count">{{ userInfo.giftCardCount || 0 }}张</span>
        </div>
      </div>
    </div>

    <!-- 功能列表 -->
    <div class="feature-list">
      <div class="feature-group">
        <div class="feature-item" @click="navigateTo('emergency-contact')">
          <span class="icon">👥</span>
          <span class="text">紧急联系人</span>
          <span class="arrow">›</span>
        </div>
        <div class="feature-item" @click="navigateTo('favorite-addresses')">
          <span class="icon">📍</span>
          <span class="text">常用地址</span>
          <span class="arrow">›</span>
        </div>
      </div>

      <div class="feature-group">
        <div class="feature-item" @click="navigateTo('payment-methods')">
          <span class="icon">💳</span>
          <span class="text">支付管理</span>
          <span class="arrow">›</span>
        </div>
        <div class="feature-item" @click="navigateTo('invoice')">
          <span class="icon">📃</span>
          <span class="text">发票管理</span>
          <span class="arrow">›</span>
        </div>
      </div>

      <div class="feature-group">
        <div class="feature-item" @click="navigateTo('settings')">
          <span class="icon">⚙️</span>
          <span class="text">设置</span>
          <span class="arrow">›</span>
        </div>
        <div class="feature-item" @click="navigateTo('help')">
          <span class="icon">❓</span>
          <span class="text">帮助与客服</span>
          <span class="arrow">›</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const defaultAvatar = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMTAwIiBoZWlnaHQ9IjEwMCIgdmlld0JveD0iMCAwIDEwMCAxMDAiIGZpbGw9Im5vbmUiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CiAgPGNpcmNsZSBjeD0iNTAiIGN5PSI1MCIgcj0iNTAiIGZpbGw9IiNFMEUwRTAiLz4KICA8Y2lyY2xlIGN4PSI1MCIgY3k9IjQwIiByPSIyMCIgZmlsbD0iI0JEQkRCRCIvPgogIDxwYXRoIGQ9Ik01MCA2NUMyOCA2NSAxMCA4NSAxMCA4NVYxMDBIOTBWODVDOTAgODUgNzIgNjUgNTAgNjVaIiBmaWxsPSIjQkRCREJEIi8+Cjwvc3ZnPgo='

const userInfo = ref({
  nickname: '罗小黑',
  phone: '13800138000',
  avatar: '',
  tripCount: 12,
  favoriteCount: 3,
  balance: '288.00',
  couponCount: 5,
  giftCardCount: 2
})

// 格式化手机号
const formatPhone = (phone) => {
  if (!phone) return ''
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

// 页面导航
const navigateTo = (route) => {
  router.push(`/user/${route}`)
}

onMounted(() => {
  // 这里可以添加获取用户信息的API调用
  console.log('个人页面加载完成')
})
</script>

<style scoped>
.user-profile {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 20px;
}

.user-header {
  background: linear-gradient(135deg, #00C896 0%, #00B085 100%);
  padding: 20px;
  color: white;
  border-radius: 0 0 20px 20px;
}

.user-info {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.avatar {
  position: relative;
  width: 80px;
  height: 80px;
  margin-right: 15px;
}

.avatar img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 2px solid rgba(255, 255, 255, 0.3);
}

.edit-badge {
  position: absolute;
  bottom: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  font-size: 12px;
  padding: 2px 6px;
  border-radius: 10px;
}

.user-details {
  flex: 1;
}

.username {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

.phone {
  margin: 5px 0 0;
  font-size: 14px;
  opacity: 0.8;
}

.user-stats {
  display: flex;
  justify-content: space-around;
  text-align: center;
}

.stat-item {
  display: flex;
  flex-direction: column;
}

.stat-item .number {
  font-size: 20px;
  font-weight: 600;
}

.stat-item .label {
  font-size: 14px;
  opacity: 0.8;
}

.wallet-card {
  margin: 20px;
  background: white;
  border-radius: 15px;
  padding: 20px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.balance-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.balance-info {
  display: flex;
  flex-direction: column;
}

.balance-info .label {
  font-size: 14px;
  color: #666;
}

.balance-info .amount {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin-top: 5px;
}

.recharge-btn {
  background: #00C896;
  color: white;
  border: none;
  padding: 8px 20px;
  border-radius: 20px;
  font-size: 14px;
  cursor: pointer;
}

.wallet-items {
  display: flex;
  justify-content: space-between;
}

.wallet-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.wallet-item .icon {
  font-size: 20px;
}

.wallet-item .text {
  font-size: 14px;
  color: #666;
}

.wallet-item .count {
  font-size: 14px;
  color: #00C896;
}

.feature-list {
  margin: 20px;
}

.feature-group {
  background: white;
  border-radius: 15px;
  margin-bottom: 15px;
  overflow: hidden;
}

.feature-item {
  display: flex;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #f5f5f5;
  cursor: pointer;
}

.feature-item:last-child {
  border-bottom: none;
}

.feature-item .icon {
  font-size: 20px;
  margin-right: 12px;
}

.feature-item .text {
  flex: 1;
  font-size: 16px;
  color: #333;
}

.feature-item .arrow {
  color: #999;
  font-size: 18px;
}

@media (max-width: 768px) {
  .wallet-card,
  .feature-list {
    margin: 15px;
  }

  .username {
    font-size: 20px;
  }

  .balance-info .amount {
    font-size: 20px;
  }
}
</style> 