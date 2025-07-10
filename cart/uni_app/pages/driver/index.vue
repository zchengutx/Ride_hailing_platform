<template>
	<view class="container">
		<view class="header">
			<text class="title">车主服务</text>
			<text class="subtitle">成为司机，开启收益之旅</text>
		</view>
		
		<view class="services">
			<view class="service-item" @tap="driverRegister">
				<image src="/static/driver-register.png" mode="aspectFit"></image>
				<text class="service-title">司机注册</text>
				<text class="service-desc">注册成为司机，获得更多收益</text>
			</view>
			
			<view class="service-item" @tap="driverLogin">
				<image src="/static/driver-login.png" mode="aspectFit"></image>
				<text class="service-title">司机登录</text>
				<text class="service-desc">已注册司机请登录管理订单</text>
			</view>
		</view>
	</view>
</template>

<script>
export default {
	data() {
		return {
			
		}
	},
	methods: {
		driverRegister() {
			uni.showToast({
				title: '司机注册功能开发中',
				icon: 'none'
			});
		},
		
		driverLogin() {
			uni.showToast({
				title: '司机登录功能开发中',
				icon: 'none'
			});
		}
	}
}
</script>

<style>
.container {
	background-color: #f8f9fa;
	min-height: 100vh;
	padding: 40rpx;
}

.header {
	text-align: center;
	margin-bottom: 80rpx;
}

.title {
	display: block;
	font-size: 48rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 16rpx;
}

.subtitle {
	font-size: 28rpx;
	color: #666;
}

.services {
	display: flex;
	flex-direction: column;
	gap: 32rpx;
}

.service-item {
	background: white;
	border-radius: 16rpx;
	padding: 40rpx;
	text-align: center;
	box-shadow: 0 4rpx 12rpx rgba(0,0,0,0.1);
}

.service-item image {
	width: 120rpx;
	height: 120rpx;
	margin-bottom: 24rpx;
}

.service-title {
	display: block;
	font-size: 36rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 16rpx;
}

.service-desc {
	font-size: 28rpx;
	color: #666;
}
</style> 