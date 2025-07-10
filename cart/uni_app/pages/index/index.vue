<template>
	<view class="container">
		<!-- 顶部状态栏 -->
		<view class="status-bar">
			<view class="location-info" @tap="selectCity">
				<text class="city-name">{{ currentCity }}</text>
				<text class="arrow-down">▼</text>
			</view>
			<view class="top-actions">
				<view class="qr-code">
					<text class="icon">🎫</text>
					<text>乘车码</text>
				</view>
				<view class="bike-sharing">
					<text class="icon">🔍</text>
					<text>扫一扫</text>
				</view>
			</view>
		</view>

		<!-- 地图区域 -->
		<view class="map-container">
			<!-- #ifdef H5 -->
			<view id="baiduMap" class="map"></view>
			<!-- #endif -->
			
			<!-- #ifndef H5 -->
			<map 
				id="map" 
				:longitude="longitude" 
				:latitude="latitude" 
				:scale="16"
				:markers="markers"
				:show-location="true"
				:enable-3D="false"
				:enable-traffic="false"
				class="map"
				@tap="onMapTap"
				@markertap="onMarkerTap"
			></map>
			<!-- #endif -->
			
			<!-- 最近上车位置提示 -->
			<view class="recent-pickup" v-if="recentLocation">
				<view class="pickup-tag">最近上车</view>
				<text class="pickup-location">{{ recentLocation }}</text>
				<text class="arrow-right">›</text>
			</view>

			<!-- 定位按钮 -->
			<view class="location-btn" @tap="getCurrentLocation">
				<text class="location-icon">📍</text>
			</view>
		</view>

		<!-- 促销横幅 -->
		<view class="promotion-banner" v-if="promotionInfo">
			<view class="promotion-content">
				<text class="promotion-title">{{ promotionInfo.title }}</text>
				<text class="promotion-subtitle">{{ promotionInfo.subtitle }}</text>
			</view>
			<view class="promotion-coupon">
				<text class="coupon-icon">🎫</text>
			</view>
		</view>

		<!-- 行程规划区域 -->
		<view class="trip-planner">
			<view class="location-input from-location" @tap="selectStartLocation">
				<view class="location-dot green"></view>
				<view class="location-content">
					<text class="location-label">从</text>
					<text class="location-text">{{ startLocation || '东禄素质教育产业园' }}</text>
				</view>
			</view>
			
			<view class="location-input to-location" @tap="selectEndLocation">
				<view class="location-dot orange"></view>
				<view class="location-content">
					<text class="location-label">您想去哪儿？</text>
					<text class="location-text">{{ endLocation || '去：上海市浦东医院' }}</text>
				</view>
			</view>
		</view>

		<!-- 功能按钮区域 -->
		<view class="action-buttons">
			<view class="action-btn" @tap="scheduleRide">
				<text class="action-icon">⏰</text>
				<text>预约</text>
			</view>
			<view class="action-btn" @tap="callForOthers">
				<text class="action-icon">👥</text>
				<text>帮人叫车</text>
			</view>
			<view class="action-btn" @tap="pickupService">
				<text class="action-icon">✈️</text>
				<text>接送机</text>
			</view>
		</view>

		<!-- 服务网格 -->
		<view class="services-grid">
			<!-- 第一行 -->
			<view class="service-row">
				<view class="service-item" @tap="selectService('taxi')" :class="{ active: selectedService === 'taxi' }">
					<text class="service-icon">🚕</text>
					<text>打车</text>
				</view>
				<view class="service-item" @tap="selectService('carpool')" :class="{ active: selectedService === 'carpool' }">
					<text class="service-icon">🚗</text>
					<text>顺风车</text>
				</view>
				<view class="service-item" @tap="selectService('driver')" :class="{ active: selectedService === 'driver' }">
					<text class="service-icon">🚘</text>
					<text>代驾</text>
				</view>
				<view class="service-item" @tap="selectService('pet')" :class="{ active: selectedService === 'pet' }">
					<text class="service-icon">🐕</text>
					<text>宠物出行</text>
				</view>
				<view class="service-item" @tap="selectService('intercity')" :class="{ active: selectedService === 'intercity' }">
					<text class="service-icon">🛣️</text>
					<text>城际拼车</text>
				</view>
			</view>

			<!-- 第二行 -->
			<view class="service-row">
				<view class="service-item" @tap="selectService('rental')" :class="{ active: selectedService === 'rental' }">
					<text class="service-icon">🚖</text>
					<text>出租车</text>
				</view>
				<view class="service-item" @tap="selectService('delivery')" :class="{ active: selectedService === 'delivery' }">
					<text class="service-icon">📦</text>
					<text>快送跑腿</text>
					<view class="service-badge">货因</view>
				</view>
				<view class="service-item" @tap="selectService('money')" :class="{ active: selectedService === 'money' }">
					<text class="service-icon">💰</text>
					<text>借钱</text>
				</view>
				<view class="service-item" @tap="selectService('gas')" :class="{ active: selectedService === 'gas' }">
					<text class="service-icon">⛽</text>
					<text>加油/充电</text>
				</view>
				<view class="service-item" @tap="selectService('transit')" :class="{ active: selectedService === 'transit' }">
					<text class="service-icon">🚇</text>
					<text>公交地铁</text>
				</view>
			</view>

			<!-- 第三行 -->
			<view class="service-row">
				<view class="service-item" @tap="selectService('bike')" :class="{ active: selectedService === 'bike' }">
					<text class="service-icon">🚲</text>
					<text>青桔骑行</text>
				</view>
				<view class="service-item" @tap="selectService('freight')" :class="{ active: selectedService === 'freight' }">
					<text class="service-icon">🚚</text>
					<text>送货</text>
				</view>
				<view class="service-item" @tap="selectService('car-rental')" :class="{ active: selectedService === 'car-rental' }">
					<text class="service-icon">🚙</text>
					<text>滴滴租车</text>
					<view class="service-badge free">免费领</view>
				</view>
				<view class="service-item" @tap="selectService('moving')" :class="{ active: selectedService === 'moving' }">
					<text class="service-icon">🏠</text>
					<text>搬家</text>
					<view class="service-badge discount">特惠搬</view>
				</view>
				<view class="service-item" @tap="selectService('more')" :class="{ active: selectedService === 'more' }">
					<text class="service-icon">➕</text>
					<text>更多服务</text>
				</view>
			</view>
		</view>

		<!-- 优惠券区域 -->
		<view class="coupon-section">
			<view class="coupon-header">
				<text class="coupon-title">App专属优惠券</text>
				<text class="coupon-subtitle">本单可减3元</text>
				<text class="more-link">更多</text>
			</view>
			<view class="coupon-list">
				<view class="coupon-item" v-for="(coupon, index) in coupons" :key="index">
					<text class="coupon-icon">🎟️</text>
					<text>{{ coupon.title }}</text>
				</view>
			</view>
		</view>

	</view>
</template>

<script>
import api from '../../utils/api.js'

export default {
	data() {
		return {
			// 基础位置信息
			currentCity: '上海',
			longitude: 121.473701,
			latitude: 31.230416,
			recentLocation: '东禄素质教育产业园',
			startLocation: '',
			endLocation: '',
			startLng: 0,
			startLat: 0,
			endLng: 0,
			endLat: 0,
			
			// 地图标记
			markers: [
				{
					id: 1,
					latitude: 31.230416,
					longitude: 121.473701,
					width: 30,
					height: 30,
					anchor: {x: 0.5, y: 1}
				}
			],
			
			// 选中的服务
			selectedService: 'taxi',
			
			// 促销信息
			promotionInfo: {
				title: '滴滴App专享',
				subtitle: '本单用App打车减3元',
				bgImage: '/static/promotion-bg.png'
			},
			
			// 优惠券列表
			coupons: [
				{ icon: '/static/coupon-taxi.png', title: '打车立减券' },
				{ icon: '/static/coupon-taxi.png', title: '打车立减券' },
				{ icon: '/static/coupon-taxi.png', title: '打车立减券' }
			],
			
			// 用户信息
			userInfo: null,
			token: '',
			
			// 百度地图实例
			baiduMap: null,
			mapContext: null
		}
	},
	
	onLoad() {
		this.checkLoginStatus();
		// 延迟获取位置，确保页面渲染完成
		setTimeout(() => {
			this.getCurrentLocation();
		}, 100);
	},
	
	onReady() {
		// #ifndef H5
		// 获取地图上下文
		this.mapContext = uni.createMapContext('map');
		// #endif
		
		// #ifdef H5
		// H5端初始化百度地图
		this.initBaiduMap();
		// #endif
	},
	
	methods: {
		// H5端初始化百度地图
		initBaiduMap() {
			// #ifdef H5
			// 动态加载百度地图JS API
			if (!window.BMap) {
				const script = document.createElement('script');
				script.type = 'text/javascript';
				script.src = 'https://api.map.baidu.com/api?v=3.0&ak=mGQEH6rmm11Em096OUIIY4SUfUI5AD3u&callback=initBMap';
				document.head.appendChild(script);
				
				// 定义全局回调函数
				window.initBMap = () => {
					this.createBaiduMap();
				};
			} else {
				this.createBaiduMap();
			}
			// #endif
		},
		
		// 创建百度地图实例
		createBaiduMap() {
			// #ifdef H5
			try {
				// 创建地图实例
				this.baiduMap = new window.BMap.Map('baiduMap');
				
				// 创建点坐标
				const point = new window.BMap.Point(this.longitude, this.latitude);
				
				// 初始化地图，设置中心点坐标和地图级别
				this.baiduMap.centerAndZoom(point, 16);
				
				// 开启鼠标滚轮缩放
				this.baiduMap.enableScrollWheelZoom(true);
				
				// 添加当前位置标记
				const marker = new window.BMap.Marker(point);
				this.baiduMap.addOverlay(marker);
				
				// 添加定位控件
				const geolocationControl = new window.BMap.GeolocationControl();
				this.baiduMap.addControl(geolocationControl);
				
			} catch (error) {
				console.error('创建百度地图失败:', error);
			}
			// #endif
		},
		
		// 检查登录状态
		checkLoginStatus() {
			this.token = uni.getStorageSync('token');
			this.userInfo = uni.getStorageSync('userInfo');
			
			if (!this.token) {
				// 未登录，跳转到登录页面
				// uni.navigateTo({ url: '/pages/login/login' });
			}
		},
		
		// 获取当前位置
		getCurrentLocation() {
			// #ifdef H5
			uni.getLocation({
				type: 'wgs84',
				success: (res) => {
					this.latitude = res.latitude;
					this.longitude = res.longitude;
					this.updateMarkers();
					this.reverseGeocode(res.latitude, res.longitude);
				},
				fail: (err) => {
					console.error('获取位置失败:', err);
					this.updateMarkers();
					uni.showModal({
						title: '提示',
						content: '获取位置失败，请检查是否开启定位权限',
						showCancel: false
					});
				}
			});
			// #endif

			// #ifndef H5
			uni.getLocation({
				type: 'gcj02',
				success: (res) => {
					this.latitude = res.latitude;
					this.longitude = res.longitude;
					this.updateMarkers();
					this.reverseGeocode(res.latitude, res.longitude);
				},
				fail: (err) => {
					console.error('获取位置失败:', err);
					this.updateMarkers();
					uni.showModal({
						title: '提示',
						content: '获取位置失败，请检查是否开启定位权限',
						showCancel: false
					});
				}
			});
			// #endif
		},
		
		// 更新地图标记
		updateMarkers() {
			this.markers = [
				{
					id: 1,
					latitude: this.latitude,
					longitude: this.longitude,
					width: 30,
					height: 30,
					anchor: {x: 0.5, y: 1}
				}
			];
		},
		
		// 逆地理编码获取地址
		async reverseGeocode(lat, lng) {
			try {
				const response = await api.reverseGeocoding(lat, lng);
				if (response.data) {
					this.recentLocation = response.data.formatted_address;
					// 从地址中提取城市信息
					if (response.data.addressComponent) {
						this.currentCity = response.data.addressComponent.city || this.currentCity;
					}
				}
			} catch (error) {
				console.error('逆地理编码失败:', error);
			}
		},
		
		// 加载首页数据
		async loadHomePageData() {
			if (!this.token) return;
			
			try {
				const response = await api.getHomePage(this.userInfo?.id, this.recentLocation);
				const data = response.data;
				this.recentLocation = data.current_location || this.recentLocation;
				// 可以根据返回的数据更新其他信息
				if (data.services) {
					// 更新服务信息
				}
				if (data.weather_info) {
					// 更新天气信息
				}
			} catch (error) {
				console.error('加载首页数据失败:', error);
			}
		},
		
		// 选择起始位置
		selectStartLocation() {
			uni.navigateTo({
				url: '/pages/location/select?type=start'
			});
		},
		
		// 选择目的地
		selectEndLocation() {
			uni.navigateTo({
				url: '/pages/location/select?type=end'
			});
		},
		
		// 选择服务类型
		selectService(serviceType) {
			this.selectedService = serviceType;
			
			// 根据不同服务类型跳转到对应页面
			switch(serviceType) {
				case 'taxi':
					// 打车服务保持在当前页面
					break;
				case 'carpool':
					uni.navigateTo({ url: '/pages/carpool/index' });
					break;
				case 'driver':
					uni.navigateTo({ url: '/pages/driver-service/index' });
					break;
				// 其他服务类型...
				default:
					uni.showToast({
						title: '服务即将上线',
						icon: 'none'
					});
			}
		},
		
		// 创建订单
		async createOrder(orderType) {
			if (!this.endLocation) {
				uni.showToast({
					title: '请选择目的地',
					icon: 'none'
				});
				return;
			}
			
			try {
				const orderData = {
					passengerId: this.userInfo?.id,
					startAddr: this.startLocation || this.recentLocation,
					startLng: this.startLng || this.longitude,
					startLat: this.startLat || this.latitude,
					endAddr: this.endLocation,
					endLng: this.endLng,
					endLat: this.endLat,
					orderType: orderType
				};
				
				const response = await api.createOrder(orderData);
				uni.navigateTo({
					url: `/pages/order/detail?orderCode=${response.data.orderCode}`
				});
			} catch (error) {
				console.error('创建订单失败:', error);
			}
		},
		
		// 预约用车
		scheduleRide() {
			uni.navigateTo({
				url: '/pages/schedule/index'
			});
		},
		
		// 帮人叫车
		callForOthers() {
			uni.navigateTo({
				url: '/pages/call-others/index'
			});
		},
		
		// 接送机
		pickupService() {
			uni.navigateTo({
				url: '/pages/pickup/index'
			});
		},
		
		// 地图点击事件
		onMapTap(e) {
			console.log('地图点击:', e);
		},
		
		// 标记点击事件
		onMarkerTap(e) {
			console.log('标记点击:', e);
		},
		
		// 选择城市
		selectCity() {
			uni.showActionSheet({
				itemList: ['上海', '北京', '广州', '深圳', '杭州'],
				success: (res) => {
					this.currentCity = ['上海', '北京', '广州', '深圳', '杭州'][res.tapIndex];
				}
			});
		}
		
	}
}
</script>

<style>
.container {
	min-height: 100vh;
	background-color: #f8f9fa;
}

/* 顶部状态栏 */
.status-bar {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 20rpx 32rpx;
	background-color: #fff;
	box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.1);
}

.location-info {
	display: flex;
	align-items: center;
}

.city-name {
	font-size: 32rpx;
	font-weight: bold;
	color: #333;
	margin-right: 8rpx;
}

.arrow-down {
	font-size: 16rpx;
	color: #666;
	margin-left: 4rpx;
}

.top-actions {
	display: flex;
	gap: 32rpx;
}

.qr-code, .bike-sharing {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.qr-code .icon, .bike-sharing .icon {
	font-size: 48rpx;
	margin-bottom: 8rpx;
}

.qr-code text:last-child, .bike-sharing text:last-child {
	font-size: 24rpx;
	color: #666;
}

/* 地图容器 */
.map-container {
	position: relative;
	height: 400rpx;
	margin: 16rpx;
	border-radius: 16rpx;
	overflow: hidden;
}

.map {
	width: 100%;
	height: 100%;
}

/* #ifdef H5 */
#baiduMap {
	width: 100%;
	height: 100%;
}
/* #endif */

.recent-pickup {
	position: absolute;
	top: 24rpx;
	left: 24rpx;
	right: 24rpx;
	background: rgba(255,255,255,0.95);
	border-radius: 24rpx;
	padding: 16rpx 24rpx;
	display: flex;
	align-items: center;
	backdrop-filter: blur(10rpx);
}

.pickup-tag {
	background: #4CAF50;
	color: white;
	font-size: 24rpx;
	padding: 8rpx 16rpx;
	border-radius: 12rpx;
	margin-right: 16rpx;
}

.pickup-location {
	flex: 1;
	font-size: 28rpx;
	color: #333;
}

.arrow-right {
	font-size: 32rpx;
	color: #999;
}

.location-btn {
	position: absolute;
	bottom: 24rpx;
	right: 24rpx;
	width: 80rpx;
	height: 80rpx;
	background: white;
	border-radius: 40rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	box-shadow: 0 4rpx 12rpx rgba(0,0,0,0.15);
}

.location-icon {
	font-size: 40rpx;
}

/* 促销横幅 */
.promotion-banner {
	position: relative;
	margin: 16rpx;
	height: 120rpx;
	border-radius: 16rpx;
	overflow: hidden;
	display: flex;
	align-items: center;
	background: linear-gradient(135deg, #FF6B6B 0%, #FF8E53 100%);
}

.promotion-content {
	flex: 1;
	padding: 0 32rpx;
}

.promotion-title {
	display: block;
	color: white;
	font-size: 28rpx;
	font-weight: bold;
}

.promotion-subtitle {
	color: rgba(255,255,255,0.9);
	font-size: 24rpx;
}

.promotion-coupon {
	padding-right: 32rpx;
}

.coupon-icon {
	font-size: 60rpx;
}

/* 行程规划 */
.trip-planner {
	background: white;
	margin: 16rpx;
	border-radius: 16rpx;
	padding: 32rpx;
}

.location-input {
	display: flex;
	align-items: center;
	padding: 24rpx 0;
}

.location-input:first-child {
	border-bottom: 2rpx solid #f0f0f0;
}

.location-dot {
	width: 16rpx;
	height: 16rpx;
	border-radius: 50%;
	margin-right: 24rpx;
}

.location-dot.green {
	background: #4CAF50;
}

.location-dot.orange {
	background: #FF9800;
}

.location-content {
	flex: 1;
}

.location-label {
	display: block;
	font-size: 24rpx;
	color: #999;
	margin-bottom: 8rpx;
}

.location-text {
	font-size: 32rpx;
	color: #333;
}

/* 功能按钮 */
.action-buttons {
	display: flex;
	background: white;
	margin: 0 16rpx 16rpx;
	border-radius: 16rpx;
	padding: 24rpx;
}

.action-btn {
	flex: 1;
	display: flex;
	flex-direction: column;
	align-items: center;
}

.action-icon {
	font-size: 48rpx;
	margin-bottom: 8rpx;
}

.action-btn text:last-child {
	font-size: 24rpx;
	color: #666;
}

/* 服务网格 */
.services-grid {
	background: white;
	margin: 0 16rpx 16rpx;
	border-radius: 16rpx;
	padding: 32rpx;
}

.service-row {
	display: flex;
	justify-content: space-between;
	margin-bottom: 48rpx;
}

.service-row:last-child {
	margin-bottom: 0;
}

.service-item {
	position: relative;
	display: flex;
	flex-direction: column;
	align-items: center;
	width: 120rpx;
}

.service-icon {
	font-size: 60rpx;
	margin-bottom: 12rpx;
}

.service-item text:last-child {
	font-size: 24rpx;
	color: #333;
	text-align: center;
}

.service-item.active {
	transform: scale(1.05);
}

.service-badge {
	position: absolute;
	top: -8rpx;
	right: -8rpx;
	background: #FF4444;
	color: white;
	font-size: 18rpx;
	padding: 4rpx 8rpx;
	border-radius: 8rpx;
	transform: scale(0.8);
}

.service-badge.free {
	background: #4CAF50;
}

.service-badge.discount {
	background: #FF9800;
}

/* 优惠券区域 */
.coupon-section {
	background: white;
	margin: 0 16rpx 16rpx;
	border-radius: 16rpx;
	padding: 32rpx;
}

.coupon-header {
	display: flex;
	align-items: center;
	margin-bottom: 24rpx;
}

.coupon-title {
	font-size: 32rpx;
	font-weight: bold;
	color: #333;
}

.coupon-subtitle {
	font-size: 24rpx;
	color: #FF4444;
	margin-left: 16rpx;
}

.more-link {
	margin-left: auto;
	font-size: 24rpx;
	color: #666;
}

.coupon-list {
	display: flex;
	gap: 24rpx;
}

.coupon-item {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.coupon-item .coupon-icon {
	font-size: 48rpx;
	margin-bottom: 8rpx;
}

.coupon-item text:last-child {
	font-size: 20rpx;
	color: #666;
}

</style>
