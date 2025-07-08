<template>
	<view class="login-container">
		<!-- 顶部logo区域 -->
		<view class="logo-section">
			<view class="logo">
				<text class="logo-text">滴滴出行</text>
			</view>
			<view class="subtitle">
				<text class="subtitle-text">一键登录，开启出行之旅</text>
			</view>
		</view>

		<!-- 登录表单 -->
		<view class="form-section">
			<!-- 手机号输入 -->
			<view class="input-group">
				<view class="input-wrapper">
					<uni-icons type="phone" size="20" color="#999"></uni-icons>
					<input 
						class="input" 
						type="number" 
						placeholder="请输入手机号" 
						v-model="mobile"
						maxlength="11"
						@input="onMobileInput"
					/>
				</view>
			</view>

			<!-- 验证码输入 -->
			<view class="input-group">
				<view class="input-wrapper">
					<uni-icons type="compose" size="20" color="#999"></uni-icons>
					<input 
						class="input" 
						type="number" 
						placeholder="请输入验证码" 
						v-model="smsCode"
						maxlength="6"
					/>
					<view class="sms-btn" @click="sendSmsCode" :class="{ disabled: !canSendSms }">
						<text class="sms-btn-text">{{ smsButtonText }}</text>
					</view>
				</view>
			</view>

			<!-- 登录按钮 -->
			<view class="login-btn" @click="doLogin" :class="{ disabled: !canLogin }">
				<text class="login-btn-text">登录</text>
			</view>

			<!-- 协议同意 -->
			<view class="agreement-section">
				<view class="agreement-item" @click="toggleAgreement">
					<uni-icons 
						:type="agreedToTerms ? 'checkbox-filled' : 'checkbox'" 
						size="16" 
						:color="agreedToTerms ? '#FF6600' : '#999'"
					></uni-icons>
					<text class="agreement-text">
						我已阅读并同意
						<text class="link-text" @click.stop="showTerms">《用户协议》</text>
						和
						<text class="link-text" @click.stop="showPrivacy">《隐私政策》</text>
					</text>
				</view>
			</view>
		</view>

		<!-- 快捷登录 -->
		<view class="quick-login-section">
			<view class="quick-login-title">
				<text class="title-text">其他登录方式</text>
			</view>
			<view class="quick-login-methods">
				<view class="login-method" @click="wechatLogin">
					<uni-icons type="weixin" size="30" color="#4CAF50"></uni-icons>
					<text class="method-text">微信</text>
				</view>
				<view class="login-method" @click="alipayLogin">
					<uni-icons type="wallet" size="30" color="#2196F3"></uni-icons>
					<text class="method-text">支付宝</text>
				</view>
			</view>
		</view>
	</view>
</template>

<script>
import { sendSms, login } from '@/api/user.js'
import { setToken, setUserInfo } from '@/utils/auth.js'

export default {
	data() {
		return {
			mobile: '',
			smsCode: '',
			agreedToTerms: false,
			smsCountdown: 0,
			countdownTimer: null
		}
	},
	
	computed: {
		// 是否可以发送验证码
		canSendSms() {
			return this.mobile.length === 11 && this.smsCountdown === 0 && this.agreedToTerms
		},
		
		// 是否可以登录
		canLogin() {
			return this.mobile.length === 11 && this.smsCode.length === 6 && this.agreedToTerms
		},
		
		// 验证码按钮文本
		smsButtonText() {
			if (this.smsCountdown > 0) {
				return `${this.smsCountdown}s`
			}
			return '获取验证码'
		}
	},
	
	onUnload() {
		// 清除定时器
		if (this.countdownTimer) {
			clearInterval(this.countdownTimer)
		}
	},
	
	methods: {
		// 手机号输入处理
		onMobileInput() {
			// 限制只能输入数字
			this.mobile = this.mobile.replace(/[^\d]/g, '')
		},
		
		// 发送验证码
		async sendSmsCode() {
			if (!this.canSendSms) {
				return
			}
			
			// 验证手机号格式
			if (!this.validateMobile(this.mobile)) {
				uni.showToast({
					title: '请输入正确的手机号',
					icon: 'none'
				})
				return
			}
			
			try {
				await sendSms(this.mobile)
				
				uni.showToast({
					title: '验证码已发送',
					icon: 'success'
				})
				
				// 开始倒计时
				this.startCountdown()
				
			} catch (error) {
				console.error('发送验证码失败：', error)
			}
		},
		
		// 开始倒计时
		startCountdown() {
			this.smsCountdown = 60
			this.countdownTimer = setInterval(() => {
				this.smsCountdown--
				if (this.smsCountdown <= 0) {
					clearInterval(this.countdownTimer)
					this.countdownTimer = null
				}
			}, 1000)
		},
		
		// 登录
		async doLogin() {
			if (!this.canLogin) {
				return
			}
			
			// 验证手机号格式
			if (!this.validateMobile(this.mobile)) {
				uni.showToast({
					title: '请输入正确的手机号',
					icon: 'none'
				})
				return
			}
			
			// 验证验证码格式
			if (this.smsCode.length !== 6) {
				uni.showToast({
					title: '请输入6位验证码',
					icon: 'none'
				})
				return
			}
			
			try {
				const result = await login(this.mobile, this.smsCode)
				
				// 保存token和用户信息
				setToken(result.data.token)
				setUserInfo({
					userId: result.data.user_id,
					mobile: this.mobile,
					expiresAt: result.data.expires_at
				})
				
				uni.showToast({
					title: '登录成功',
					icon: 'success'
				})
				
				// 延迟跳转到首页
				setTimeout(() => {
					uni.reLaunch({
						url: '/pages/index/index'
					})
				}, 1500)
				
			} catch (error) {
				console.error('登录失败：', error)
			}
		},
		
		// 验证手机号格式
		validateMobile(mobile) {
			const mobileReg = /^1[3-9]\d{9}$/
			return mobileReg.test(mobile)
		},
		
		// 切换协议同意状态
		toggleAgreement() {
			this.agreedToTerms = !this.agreedToTerms
		},
		
		// 显示用户协议
		showTerms() {
			uni.showModal({
				title: '用户协议',
				content: '这里是用户协议内容...',
				showCancel: false
			})
		},
		
		// 显示隐私政策
		showPrivacy() {
			uni.showModal({
				title: '隐私政策',
				content: '这里是隐私政策内容...',
				showCancel: false
			})
		},
		
		// 微信登录
		wechatLogin() {
			uni.showToast({
				title: '微信登录功能开发中',
				icon: 'none'
			})
		},
		
		// 支付宝登录
		alipayLogin() {
			uni.showToast({
				title: '支付宝登录功能开发中',
				icon: 'none'
			})
		}
	}
}
</script>

<style lang="scss" scoped>
.login-container {
	min-height: 100vh;
	background: linear-gradient(135deg, #FF6600, #FF8533);
	padding: 0 30px;
	display: flex;
	flex-direction: column;
}

.logo-section {
	flex: 1;
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	padding-top: 100px;
	
	.logo {
		margin-bottom: 20px;
		
		.logo-text {
			font-size: 36px;
			font-weight: bold;
			color: white;
		}
	}
	
	.subtitle {
		.subtitle-text {
			font-size: 16px;
			color: rgba(255, 255, 255, 0.8);
		}
	}
}

.form-section {
	flex: 2;
	padding: 40px 0;
	
	.input-group {
		margin-bottom: 20px;
		
		.input-wrapper {
			background-color: white;
			border-radius: 25px;
			padding: 0 20px;
			display: flex;
			align-items: center;
			height: 50px;
			box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
			
			.input {
				flex: 1;
				font-size: 16px;
				margin-left: 15px;
				color: #333;
				
				&::placeholder {
					color: #999;
				}
			}
			
			.sms-btn {
				padding: 8px 15px;
				background-color: #FF6600;
				border-radius: 15px;
				
				&.disabled {
					background-color: #ccc;
				}
				
				.sms-btn-text {
					color: white;
					font-size: 14px;
				}
			}
		}
	}
	
	.login-btn {
		width: 100%;
		height: 50px;
		background-color: white;
		border-radius: 25px;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-top: 30px;
		box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
		
		&.disabled {
			background-color: rgba(255, 255, 255, 0.5);
			
			.login-btn-text {
				color: #ccc;
			}
		}
		
		.login-btn-text {
			font-size: 18px;
			font-weight: bold;
			color: #FF6600;
		}
	}
	
	.agreement-section {
		margin-top: 20px;
		
		.agreement-item {
			display: flex;
			align-items: center;
			
			.agreement-text {
				font-size: 12px;
				color: rgba(255, 255, 255, 0.8);
				margin-left: 8px;
				line-height: 1.5;
				
				.link-text {
					color: white;
					text-decoration: underline;
				}
			}
		}
	}
}

.quick-login-section {
	flex: 1;
	padding-bottom: 50px;
	
	.quick-login-title {
		text-align: center;
		margin-bottom: 20px;
		
		.title-text {
			font-size: 14px;
			color: rgba(255, 255, 255, 0.8);
		}
	}
	
	.quick-login-methods {
		display: flex;
		justify-content: center;
		gap: 50px;
		
		.login-method {
			display: flex;
			flex-direction: column;
			align-items: center;
			
			.method-text {
				font-size: 12px;
				color: white;
				margin-top: 8px;
			}
		}
	}
}
</style> 