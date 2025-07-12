<template>
	<view class="container">
		<view class="header">
			<image src="/static/logo.png" class="logo" mode="aspectFit"></image>
			<text class="title">欢迎使用优行出行</text>
			<text class="subtitle">请使用手机号登录</text>
		</view>
		
		<view class="form">
			<view class="input-group">
				<text class="label">手机号</text>
				<input 
					class="input"
					type="number"
					placeholder="请输入手机号"
					v-model="mobile"
					maxlength="11"
				/>
			</view>
			
			<view class="input-group">
				<text class="label">验证码</text>
				<view class="code-input">
					<input 
						class="input"
						type="number"
						placeholder="请输入验证码"
						v-model="smsCode"
						maxlength="6"
					/>
					<view 
						class="send-code-btn"
						:class="{ disabled: !canSendCode }"
						@tap="sendSmsCode"
					>
						{{ codeButtonText }}
					</view>
				</view>
			</view>
			
			<view class="login-btn" @tap="login" :class="{ disabled: !canLogin }">
				登录
			</view>
			
			<view class="register-link">
				<text>还没有账号？</text>
				<text class="link" @tap="goRegister">立即注册</text>
			</view>
		</view>
		
		<view class="wechat-login" @tap="wechatLogin">
			<image src="/static/wechat-icon.png" class="wechat-icon" mode="aspectFit"></image>
			<text>微信登录</text>
		</view>
	</view>
</template>

<script>
import api from '../../utils/api.js'

export default {
	data() {
		return {
			mobile: '',
			smsCode: '',
			countdown: 0,
			timer: null
		}
	},
	
	computed: {
		canSendCode() {
			return this.mobile.length === 11 && this.countdown === 0;
		},
		
		canLogin() {
			return this.mobile.length === 11 && this.smsCode.length === 6;
		},
		
		codeButtonText() {
			return this.countdown > 0 ? `${this.countdown}s` : '发送验证码';
		}
	},
	
	onUnload() {
		if (this.timer) {
			clearInterval(this.timer);
		}
	},
	
	methods: {
		// 发送短信验证码
		async sendSmsCode() {
			if (!this.canSendCode) return;
			
			if (!this.validateMobile(this.mobile)) {
				uni.showToast({
					title: '请输入正确的手机号',
					icon: 'none'
				});
				return;
			}
			
			try {
				await api.sendSms(this.mobile, 'login');
				uni.showToast({
					title: '验证码已发送',
					icon: 'success'
				});
				
				// 开始倒计时
				this.startCountdown();
			} catch (error) {
				console.error('发送验证码失败:', error);
			}
		},
		
		// 开始倒计时
		startCountdown() {
			this.countdown = 60;
			this.timer = setInterval(() => {
				this.countdown--;
				if (this.countdown <= 0) {
					clearInterval(this.timer);
					this.timer = null;
				}
			}, 1000);
		},
		
		// 登录
		async login() {
			if (!this.canLogin) return;
			
			if (!this.validateMobile(this.mobile)) {
				uni.showToast({
					title: '请输入正确的手机号',
					icon: 'none'
				});
				return;
			}
			
			if (!this.smsCode) {
				uni.showToast({
					title: '请输入验证码',
					icon: 'none'
				});
				return;
			}
			
			try {
				const response = await api.loginPassenger(this.mobile, this.smsCode);
				
				// 保存用户信息和token
				const userInfo = {
					id: response.data?.passengerId,
					mobile: this.mobile
				};
				
				uni.setStorageSync('userInfo', userInfo);
				uni.setStorageSync('token', response.data?.token || 'mock_token');
				
				uni.showToast({
					title: '登录成功',
					icon: 'success'
				});
				
				// 返回上一页或跳转到首页
				setTimeout(() => {
					const pages = getCurrentPages();
					if (pages.length > 1) {
						uni.navigateBack();
					} else {
						uni.switchTab({
							url: '/pages/index/index'
						});
					}
				}, 1500);
				
			} catch (error) {
				console.error('登录失败:', error);
			}
		},
		
		// 微信登录
		wechatLogin() {
			uni.showToast({
				title: '微信登录功能开发中',
				icon: 'none'
			});
		},
		
		// 前往注册
		goRegister() {
			uni.navigateTo({
				url: '/pages/register/register'
			});
		},
		
		// 验证手机号
		validateMobile(mobile) {
			const reg = /^1[3-9]\d{9}$/;
			return reg.test(mobile);
		}
	}
}
</script>

<style>
.container {
	background: linear-gradient(135deg, #4CAF50, #45a049);
	min-height: 100vh;
	padding: 80rpx 40rpx;
}

.header {
	text-align: center;
	margin-bottom: 80rpx;
}

.logo {
	width: 120rpx;
	height: 120rpx;
	margin-bottom: 40rpx;
}

.title {
	display: block;
	color: white;
	font-size: 48rpx;
	font-weight: bold;
	margin-bottom: 16rpx;
}

.subtitle {
	color: rgba(255, 255, 255, 0.8);
	font-size: 28rpx;
}

.form {
	background: white;
	border-radius: 20rpx;
	padding: 60rpx 40rpx;
	margin-bottom: 40rpx;
}

.input-group {
	margin-bottom: 40rpx;
}

.label {
	display: block;
	font-size: 28rpx;
	color: #333;
	margin-bottom: 16rpx;
}

.input {
	width: 100%;
	height: 88rpx;
	border: 2rpx solid #e0e0e0;
	border-radius: 12rpx;
	padding: 0 24rpx;
	font-size: 32rpx;
	box-sizing: border-box;
}

.input:focus {
	border-color: #4CAF50;
}

.code-input {
	display: flex;
	align-items: center;
	gap: 16rpx;
}

.code-input .input {
	flex: 1;
}

.send-code-btn {
	background: #4CAF50;
	color: white;
	font-size: 24rpx;
	padding: 20rpx 32rpx;
	border-radius: 12rpx;
	white-space: nowrap;
	text-align: center;
	min-width: 140rpx;
}

.send-code-btn.disabled {
	background: #ccc;
	color: #999;
}

.login-btn {
	background: #4CAF50;
	color: white;
	font-size: 36rpx;
	font-weight: bold;
	text-align: center;
	padding: 32rpx;
	border-radius: 48rpx;
	margin: 60rpx 0 40rpx;
}

.login-btn.disabled {
	background: #ccc;
}

.register-link {
	text-align: center;
	font-size: 28rpx;
	color: #666;
}

.link {
	color: #4CAF50;
	margin-left: 8rpx;
}

.wechat-login {
	background: rgba(255, 255, 255, 0.1);
	color: white;
	display: flex;
	align-items: center;
	justify-content: center;
	padding: 32rpx;
	border-radius: 16rpx;
	border: 2rpx solid rgba(255, 255, 255, 0.3);
}

.wechat-icon {
	width: 48rpx;
	height: 48rpx;
	margin-right: 16rpx;
}
</style> 