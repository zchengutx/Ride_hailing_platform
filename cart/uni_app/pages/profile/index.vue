<template>
	<view class="container">
		<!-- 用户信息卡片 -->
		<view class="user-card" v-if="userInfo">
			<view class="avatar-section">
				<image :src="userInfo.avatar || '/static/default-avatar.png'" class="avatar" mode="aspectFill"></image>
				<view class="user-info">
					<text class="username">{{ userInfo.nickName || userInfo.name || '用户' }}</text>
					<text class="phone">{{ userInfo.mobile || '未绑定手机' }}</text>
				</view>
			</view>
			<view class="edit-btn" @tap="editProfile">
				<image src="/static/edit.png" mode="aspectFit"></image>
			</view>
		</view>
		
		<!-- 登录提示 -->
		<view class="login-card" v-else>
			<text class="login-text">登录后享受更多服务</text>
			<view class="login-btn" @tap="goLogin">立即登录</view>
		</view>
		
		<!-- 功能菜单 -->
		<view class="menu-section">
			<view class="menu-group">
				<view class="menu-item" @tap="goToOrders">
					<view class="menu-left">
						<image src="/static/order-icon.png" class="menu-icon" mode="aspectFit"></image>
						<text class="menu-title">我的订单</text>
					</view>
					<image src="/static/arrow-right.png" class="arrow" mode="aspectFit"></image>
				</view>
				
				<view class="menu-item" @tap="goToFavorites">
					<view class="menu-left">
						<image src="/static/favorite-icon.png" class="menu-icon" mode="aspectFit"></image>
						<text class="menu-title">收藏地址</text>
					</view>
					<image src="/static/arrow-right.png" class="arrow" mode="aspectFit"></image>
				</view>
				
				<view class="menu-item" @tap="goToWallet">
					<view class="menu-left">
						<image src="/static/wallet-icon.png" class="menu-icon" mode="aspectFit"></image>
						<text class="menu-title">我的钱包</text>
					</view>
					<image src="/static/arrow-right.png" class="arrow" mode="aspectFit"></image>
				</view>
			</view>
			
			<view class="menu-group">
				<view class="menu-item" @tap="goToSettings">
					<view class="menu-left">
						<image src="/static/setting-icon.png" class="menu-icon" mode="aspectFit"></image>
						<text class="menu-title">设置</text>
					</view>
					<image src="/static/arrow-right.png" class="arrow" mode="aspectFit"></image>
				</view>
				
				<view class="menu-item" @tap="goToHelp">
					<view class="menu-left">
						<image src="/static/help-icon.png" class="menu-icon" mode="aspectFit"></image>
						<text class="menu-title">帮助中心</text>
					</view>
					<image src="/static/arrow-right.png" class="arrow" mode="aspectFit"></image>
				</view>
				
				<view class="menu-item" @tap="goToAbout">
					<view class="menu-left">
						<image src="/static/about-icon.png" class="menu-icon" mode="aspectFit"></image>
						<text class="menu-title">关于我们</text>
					</view>
					<image src="/static/arrow-right.png" class="arrow" mode="aspectFit"></image>
				</view>
			</view>
		</view>
		
		<!-- 退出登录 -->
		<view class="logout-section" v-if="userInfo">
			<view class="logout-btn" @tap="logout">退出登录</view>
		</view>
	</view>
</template>

<script>
import api from '../../utils/api.js'

export default {
	data() {
		return {
			userInfo: null
		}
	},
	
	onShow() {
		this.loadUserInfo();
	},
	
	methods: {
		// 加载用户信息
		loadUserInfo() {
			this.userInfo = uni.getStorageSync('userInfo');
		},
		
		// 编辑个人资料
		editProfile() {
			uni.navigateTo({
				url: '/pages/profile/edit'
			});
		},
		
		// 前往登录
		goLogin() {
			uni.navigateTo({
				url: '/pages/login/login'
			});
		},
		
		// 我的订单
		goToOrders() {
			if (!this.checkLogin()) return;
			uni.switchTab({
				url: '/pages/order/list'
			});
		},
		
		// 收藏地址
		goToFavorites() {
			if (!this.checkLogin()) return;
			uni.navigateTo({
				url: '/pages/profile/favorites'
			});
		},
		
		// 我的钱包
		goToWallet() {
			if (!this.checkLogin()) return;
			uni.navigateTo({
				url: '/pages/profile/wallet'
			});
		},
		
		// 设置
		goToSettings() {
			uni.navigateTo({
				url: '/pages/profile/settings'
			});
		},
		
		// 帮助中心
		goToHelp() {
			uni.navigateTo({
				url: '/pages/profile/help'
			});
		},
		
		// 关于我们
		goToAbout() {
			uni.navigateTo({
				url: '/pages/profile/about'
			});
		},
		
		// 检查登录状态
		checkLogin() {
			if (!this.userInfo) {
				uni.showModal({
					title: '提示',
					content: '请先登录',
					success: (res) => {
						if (res.confirm) {
							this.goLogin();
						}
					}
				});
				return false;
			}
			return true;
		},
		
		// 退出登录
		logout() {
			uni.showModal({
				title: '提示',
				content: '确定要退出登录吗？',
				success: (res) => {
					if (res.confirm) {
						// 清除本地存储
						uni.removeStorageSync('token');
						uni.removeStorageSync('userInfo');
						
						// 更新用户信息
						this.userInfo = null;
						
						uni.showToast({
							title: '已退出登录',
							icon: 'success'
						});
					}
				}
			});
		}
	}
}
</script>

<style>
.container {
	background-color: #f8f9fa;
	min-height: 100vh;
}

/* 用户信息卡片 */
.user-card {
	background: white;
	margin: 20rpx;
	border-radius: 16rpx;
	padding: 40rpx;
	display: flex;
	align-items: center;
	justify-content: space-between;
}

.avatar-section {
	display: flex;
	align-items: center;
}

.avatar {
	width: 120rpx;
	height: 120rpx;
	border-radius: 60rpx;
	margin-right: 32rpx;
}

.user-info {
	display: flex;
	flex-direction: column;
}

.username {
	font-size: 36rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 8rpx;
}

.phone {
	font-size: 28rpx;
	color: #666;
}

.edit-btn {
	width: 60rpx;
	height: 60rpx;
	display: flex;
	align-items: center;
	justify-content: center;
}

.edit-btn image {
	width: 40rpx;
	height: 40rpx;
}

/* 登录卡片 */
.login-card {
	background: white;
	margin: 20rpx;
	border-radius: 16rpx;
	padding: 60rpx 40rpx;
	text-align: center;
}

.login-text {
	display: block;
	font-size: 32rpx;
	color: #666;
	margin-bottom: 40rpx;
}

.login-btn {
	background: #4CAF50;
	color: white;
	font-size: 32rpx;
	padding: 24rpx 60rpx;
	border-radius: 48rpx;
	display: inline-block;
}

/* 菜单区域 */
.menu-section {
	margin: 20rpx;
}

.menu-group {
	background: white;
	border-radius: 16rpx;
	margin-bottom: 20rpx;
	overflow: hidden;
}

.menu-item {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 32rpx 40rpx;
	border-bottom: 1rpx solid #f0f0f0;
}

.menu-item:last-child {
	border-bottom: none;
}

.menu-left {
	display: flex;
	align-items: center;
}

.menu-icon {
	width: 48rpx;
	height: 48rpx;
	margin-right: 24rpx;
}

.menu-title {
	font-size: 32rpx;
	color: #333;
}

.arrow {
	width: 24rpx;
	height: 24rpx;
}

/* 退出登录 */
.logout-section {
	margin: 40rpx 20rpx;
}

.logout-btn {
	background: white;
	color: #FF4444;
	font-size: 32rpx;
	text-align: center;
	padding: 32rpx;
	border-radius: 16rpx;
	border: 2rpx solid #FF4444;
}
</style> 