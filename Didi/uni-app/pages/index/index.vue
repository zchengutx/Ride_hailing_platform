<template>
	<view class="container">
		<!-- 自定义导航栏 -->
		<view class="custom-navbar" :style="{ paddingTop: statusBarHeight + 'px' }">
			<view class="navbar-content">
				<!-- 左侧城市选择 -->
				<view class="city-selector" @click="selectCity">
					<text class="city-name">{{ currentCity }}</text>
					<uni-icons type="arrowdown" size="12" color="#333"></uni-icons>
				</view>
				
				<!-- 右侧功能按钮 -->
				<view class="navbar-actions">
					<view class="action-item" @click="scanCode">
						<uni-icons type="scan" size="20" color="#333"></uni-icons>
						<text class="action-text">乘车码</text>
					</view>
					<view class="action-item" @click="scanQR">
						<uni-icons type="camera" size="20" color="#333"></uni-icons>
						<text class="action-text">扫一扫</text>
					</view>
					<view class="action-item" @click="navigate">
						<uni-icons type="location" size="20" color="#333"></uni-icons>
						<text class="action-text">导航</text>
					</view>
					<view class="action-item" @click="locateMe">
						<uni-icons type="refreshempty" size="20" color="#333"></uni-icons>
						<text class="action-text">定位</text>
					</view>
				</view>
			</view>
		</view>

		<!-- 地图容器 -->
		<view class="map-container">
			<map 
				:latitude="latitude" 
				:longitude="longitude" 
				:scale="16"
				:markers="markers"
				:show-location="true"
				style="width: 100%; height: 100%;"
				@tap="onMapTap"
				@regionchange="onRegionChange"
			></map>
			
			<!-- 地图控制按钮 -->
			<view class="map-controls">
				<view class="control-btn" @click="zoomIn">
					<uni-icons type="plus" size="16" color="#333"></uni-icons>
				</view>
				<view class="control-btn" @click="zoomOut">
					<uni-icons type="minus" size="16" color="#333"></uni-icons>
				</view>
			</view>
		</view>

		<!-- 叫车面板 -->
		<view class="call-car-panel">
			<!-- 地址输入区域 -->
			<view class="address-section">
				<view class="address-item start-address" @click="selectStartAddress">
					<view class="address-dot start-dot"></view>
					<view class="address-content">
						<text class="address-text">{{ startAddress || '从东楚素质教育产业园上车' }}</text>
					</view>
				</view>
				
				<view class="address-item end-address" @click="selectEndAddress">
					<view class="address-dot end-dot"></view>
					<view class="address-content">
						<text class="address-text" :class="{ placeholder: !endAddress }">
							{{ endAddress || '您想去哪儿？' }}
						</text>
					</view>
				</view>
			</view>

			<!-- 快捷功能 -->
			<view class="quick-actions">
				<view class="quick-item" @click="bookCar">
					<uni-icons type="calendar" size="20" color="#666"></uni-icons>
					<text class="quick-text">预约</text>
				</view>
				<view class="quick-item" @click="helpOthers">
					<uni-icons type="person-add" size="20" color="#666"></uni-icons>
					<text class="quick-text">帮人叫车</text>
				</view>
				<view class="quick-item" @click="airportTransfer">
					<uni-icons type="paperplane" size="20" color="#666"></uni-icons>
					<text class="quick-text">接送机</text>
				</view>
			</view>

			<!-- 服务类型 -->
			<view class="service-types">
				<view class="service-row">
					<view class="service-item" @click="selectService('打车')">
						<uni-icons type="car" size="24" color="#FF6600"></uni-icons>
						<text class="service-text">打车</text>
					</view>
					<view class="service-item" @click="selectService('顺风车')">
						<uni-icons type="gift" size="24" color="#00C851"></uni-icons>
						<text class="service-text">顺风车</text>
					</view>
					<view class="service-item" @click="selectService('代驾')">
						<uni-icons type="gear" size="24" color="#FF6600"></uni-icons>
						<text class="service-text">代驾</text>
					</view>
					<view class="service-item" @click="selectService('城际拼车')">
						<uni-icons type="loop" size="24" color="#FF6600"></uni-icons>
						<text class="service-text">城际拼车</text>
					</view>
					<view class="service-item" @click="selectService('公交地铁')">
						<uni-icons type="navigate" size="24" color="#007AFF"></uni-icons>
						<text class="service-text">公交地铁</text>
					</view>
				</view>
				
				<view class="service-row">
					<view class="service-item" @click="selectService('特价拼车')">
						<uni-icons type="star" size="24" color="#FF6600"></uni-icons>
						<text class="service-text">特价拼车</text>
					</view>
					<view class="service-item" @click="selectService('快送跑腿')">
						<uni-icons type="cart" size="24" color="#FF6600"></uni-icons>
						<text class="service-text">快送跑腿</text>
					</view>
					<view class="service-item" @click="selectService('借钱')">
						<uni-icons type="wallet" size="24" color="#FFB74D"></uni-icons>
						<text class="service-text">借钱</text>
					</view>
					<view class="service-item" @click="selectService('优惠加油')">
						<uni-icons type="fire" size="24" color="#FF6600"></uni-icons>
						<text class="service-text">优惠加油</text>
					</view>
					<view class="service-item" @click="selectService('出租车')">
						<uni-icons type="car-filled" size="24" color="#FFD54F"></uni-icons>
						<text class="service-text">出租车</text>
					</view>
				</view>
				
				<view class="service-row">
					<view class="service-item" @click="selectService('青桔骑行')">
						<uni-icons type="compose" size="24" color="#4CAF50"></uni-icons>
						<text class="service-text">青桔骑行</text>
					</view>
					<view class="service-item" @click="selectService('送货')">
						<uni-icons type="gift-filled" size="24" color="#FF6600"></uni-icons>
						<text class="service-text">送货</text>
					</view>
					<view class="service-item" @click="selectService('滴滴租车')">
						<uni-icons type="car" size="24" color="#2196F3"></uni-icons>
						<text class="service-text">滴滴租车</text>
					</view>
					<view class="service-item" @click="selectService('搬家')">
						<uni-icons type="home" size="24" color="#FF6600"></uni-icons>
						<text class="service-text">搬家</text>
					</view>
					<view class="service-item" @click="selectService('更多服务')">
						<uni-icons type="more-filled" size="24" color="#999"></uni-icons>
						<text class="service-text">更多服务</text>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 叫车按钮 -->
		<view class="call-car-btn" @click="callCar" v-if="endAddress">
			<text class="btn-text">立即叫车</text>
		</view>
	</view>
</template>

<script>
import { takeACar } from '@/api/user.js'
import { checkLogin } from '@/utils/auth.js'

export default {
	data() {
		return {
			statusBarHeight: 0,
			currentCity: '上海',
			latitude: 31.230416,
			longitude: 121.473701,
			scale: 16,
			markers: [],
			startAddress: '',
			endAddress: '',
			selectedService: '快车'
		}
	},
	
	onLoad() {
		this.getSystemInfo()
		this.getCurrentLocation()
	},
	
	methods: {
		// 获取系统信息
		getSystemInfo() {
			uni.getSystemInfo({
				success: (res) => {
					this.statusBarHeight = res.statusBarHeight
				}
			})
		},
		
		// 获取当前位置
		getCurrentLocation() {
			uni.getLocation({
				type: 'gcj02',
				success: (res) => {
					this.latitude = res.latitude
					this.longitude = res.longitude
					this.reverseGeocode(res.latitude, res.longitude)
				},
				fail: (err) => {
					console.error('获取位置失败：', err)
					uni.showToast({
						title: '定位失败，请检查定位权限',
						icon: 'none'
					})
				}
			})
		},
		
		// 逆地理编码
		reverseGeocode(lat, lng) {
			// 这里应该调用地图API进行逆地理编码
			// 暂时使用固定地址
			this.startAddress = '东楚素质教育产业园'
		},
		
		// 城市选择
		selectCity() {
			uni.showActionSheet({
				itemList: ['上海', '北京', '广州', '深圳', '杭州'],
				success: (res) => {
					const cities = ['上海', '北京', '广州', '深圳', '杭州']
					this.currentCity = cities[res.tapIndex]
				}
			})
		},
		
		// 扫码功能
		scanCode() {
			uni.showToast({
				title: '乘车码功能开发中',
				icon: 'none'
			})
		},
		
		// 扫一扫
		scanQR() {
			uni.scanCode({
				success: (res) => {
					console.log('扫码结果：', res)
				}
			})
		},
		
		// 导航
		navigate() {
			uni.showToast({
				title: '导航功能开发中',
				icon: 'none'
			})
		},
		
		// 定位
		locateMe() {
			this.getCurrentLocation()
		},
		
		// 地图点击
		onMapTap(e) {
			console.log('地图点击：', e)
		},
		
		// 地图区域变化
		onRegionChange(e) {
			console.log('地图区域变化：', e)
		},
		
		// 放大地图
		zoomIn() {
			this.scale = Math.min(this.scale + 1, 20)
		},
		
		// 缩小地图
		zoomOut() {
			this.scale = Math.max(this.scale - 1, 5)
		},
		
		// 选择起始地址
		selectStartAddress() {
			// 跳转到地址选择页面
			uni.navigateTo({
				url: `/pages/address/address?type=start&address=${this.startAddress}`
			})
		},
		
		// 选择目的地址
		selectEndAddress() {
			// 跳转到地址选择页面
			uni.navigateTo({
				url: `/pages/address/address?type=end&address=${this.endAddress}`
			})
		},
		
		// 预约叫车
		bookCar() {
			uni.showToast({
				title: '预约功能开发中',
				icon: 'none'
			})
		},
		
		// 帮人叫车
		helpOthers() {
			uni.showToast({
				title: '帮人叫车功能开发中',
				icon: 'none'
			})
		},
		
		// 接送机
		airportTransfer() {
			uni.showToast({
				title: '接送机功能开发中',
				icon: 'none'
			})
		},
		
		// 选择服务类型
		selectService(service) {
			this.selectedService = service
			uni.showToast({
				title: `已选择${service}`,
				icon: 'none'
			})
		},
		
		// 叫车
		async callCar() {
			// 检查登录状态
			if (!checkLogin()) {
				return
			}
			
			if (!this.endAddress) {
				uni.showToast({
					title: '请选择目的地',
					icon: 'none'
				})
				return
			}
			
			try {
				const orderData = {
					startLocation: `${this.startAddress},${this.latitude},${this.longitude}`,
					endLocation: this.endAddress,
					cartType: this.selectedService
				}
				
				const result = await takeACar(orderData)
				
				uni.showModal({
					title: '叫车成功',
					content: result.msg,
					showCancel: false
				})
				
			} catch (error) {
				console.error('叫车失败：', error)
			}
		}
	}
}
</script>

<style lang="scss" scoped>
.container {
	position: relative;
	height: 100vh;
	background-color: #f5f5f5;
}

.custom-navbar {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	z-index: 1000;
	background-color: rgba(255, 255, 255, 0.95);
	backdrop-filter: blur(10px);
	border-bottom: 1px solid #e0e0e0;
	
	.navbar-content {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 10px 15px;
		
		.city-selector {
			display: flex;
			align-items: center;
			
			.city-name {
				font-size: 16px;
				font-weight: bold;
				margin-right: 5px;
			}
		}
		
		.navbar-actions {
			display: flex;
			align-items: center;
			
			.action-item {
				display: flex;
				flex-direction: column;
				align-items: center;
				margin-left: 15px;
				
				.action-text {
					font-size: 10px;
					margin-top: 2px;
					color: #666;
				}
			}
		}
	}
}

.map-container {
	position: relative;
	height: 50vh;
	margin-top: 88px;
	
	.map-controls {
		position: absolute;
		right: 15px;
		top: 50%;
		transform: translateY(-50%);
		
		.control-btn {
			width: 40px;
			height: 40px;
			background-color: rgba(255, 255, 255, 0.9);
			border-radius: 6px;
			display: flex;
			align-items: center;
			justify-content: center;
			margin-bottom: 10px;
			box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
		}
	}
}

.call-car-panel {
	position: absolute;
	bottom: 100px;
	left: 0;
	right: 0;
	background-color: white;
	border-radius: 20px 20px 0 0;
	padding: 20px 15px;
	box-shadow: 0 -5px 20px rgba(0, 0, 0, 0.1);
	
	.address-section {
		margin-bottom: 20px;
		
		.address-item {
			display: flex;
			align-items: center;
			padding: 15px 0;
			
			.address-dot {
				width: 8px;
				height: 8px;
				border-radius: 50%;
				margin-right: 15px;
				
				&.start-dot {
					background-color: #4CAF50;
				}
				
				&.end-dot {
					background-color: #FF6600;
				}
			}
			
			.address-content {
				flex: 1;
				
				.address-text {
					font-size: 16px;
					color: #333;
					
					&.placeholder {
						color: #999;
					}
				}
			}
		}
	}
	
	.quick-actions {
		display: flex;
		justify-content: space-around;
		margin-bottom: 20px;
		padding: 15px 0;
		border-top: 1px solid #f0f0f0;
		
		.quick-item {
			display: flex;
			flex-direction: column;
			align-items: center;
			
			.quick-text {
				font-size: 12px;
				color: #666;
				margin-top: 5px;
			}
		}
	}
	
	.service-types {
		.service-row {
			display: flex;
			justify-content: space-between;
			margin-bottom: 15px;
			
			.service-item {
				flex: 1;
				display: flex;
				flex-direction: column;
				align-items: center;
				padding: 10px 5px;
				
				.service-text {
					font-size: 12px;
					color: #333;
					margin-top: 5px;
					text-align: center;
				}
			}
		}
	}
}

.call-car-btn {
	position: fixed;
	bottom: 30px;
	left: 15px;
	right: 15px;
	height: 50px;
	background: linear-gradient(135deg, #FF6600, #FF8533);
	border-radius: 25px;
	display: flex;
	align-items: center;
	justify-content: center;
	box-shadow: 0 5px 15px rgba(255, 102, 0, 0.3);
	
	.btn-text {
		color: white;
		font-size: 18px;
		font-weight: bold;
	}
}
</style> 