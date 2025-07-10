<template>
	<view class="container">
		<!-- 顶部搜索栏 -->
		<view class="search-header">
			<view class="search-box">
				<image src="/static/search-icon.png" class="search-icon" mode="aspectFit"></image>
				<input 
					class="search-input" 
					:placeholder="placeholder"
					v-model="searchKeyword"
					@input="onSearchInput"
					@confirm="onSearchConfirm"
					:focus="true"
				/>
				<view class="cancel-btn" @tap="goBack">取消</view>
			</view>
		</view>

		<!-- 当前位置 -->
		<view class="current-location" @tap="selectCurrentLocation" v-if="currentLocation">
			<view class="location-item">
				<image src="/static/location-current.png" class="location-icon" mode="aspectFit"></image>
				<view class="location-info">
					<text class="location-name">当前位置</text>
					<text class="location-address">{{ currentLocation }}</text>
				</view>
				<image src="/static/arrow-right.png" class="arrow-icon" mode="aspectFit"></image>
			</view>
		</view>

		<!-- 搜索结果 -->
		<view class="search-results" v-if="searchResults.length > 0">
			<view class="section-title">搜索结果</view>
			<view class="location-list">
				<view 
					class="location-item"
					v-for="(item, index) in searchResults"
					:key="index"
					@tap="selectLocation(item)"
				>
					<image src="/static/location-search.png" class="location-icon" mode="aspectFit"></image>
					<view class="location-info">
						<text class="location-name">{{ item.name }}</text>
						<text class="location-address">{{ item.address }}</text>
					</view>
					<view class="location-distance" v-if="item.distance">
						<text>{{ formatDistance(item.distance) }}</text>
					</view>
				</view>
			</view>
		</view>

		<!-- 收藏地址 -->
		<view class="favorite-locations" v-if="favoriteLocations.length > 0 && !searchKeyword">
			<view class="section-title">收藏地址</view>
			<view class="location-list">
				<view 
					class="location-item"
					v-for="(item, index) in favoriteLocations"
					:key="index"
					@tap="selectFavoriteLocation(item)"
				>
					<image :src="getFavoriteIcon(item.locationType)" class="location-icon" mode="aspectFit"></image>
					<view class="location-info">
						<text class="location-name">{{ item.locationName }}</text>
						<text class="location-address">{{ item.address }}</text>
					</view>
					<view class="favorite-actions">
						<image 
							src="/static/delete.png" 
							class="action-icon" 
							mode="aspectFit"
							@tap.stop="deleteFavorite(item)"
						></image>
					</view>
				</view>
			</view>
		</view>

		<!-- 热门地点 -->
		<view class="hot-locations" v-if="hotLocations.length > 0 && !searchKeyword">
			<view class="section-title">热门地点</view>
			<view class="location-list">
				<view 
					class="location-item"
					v-for="(item, index) in hotLocations"
					:key="index"
					@tap="selectHotLocation(item)"
				>
					<image src="/static/location-hot.png" class="location-icon" mode="aspectFit"></image>
					<view class="location-info">
						<text class="location-name">{{ item.locationName }}</text>
						<text class="location-address">{{ item.address }}</text>
					</view>
					<view class="hot-score">
						<text>热度 {{ item.hotScore }}</text>
					</view>
				</view>
			</view>
		</view>

		<!-- 空状态 -->
		<view class="empty-state" v-if="showEmptyState">
			<image src="/static/empty-location.png" class="empty-icon" mode="aspectFit"></image>
			<text class="empty-text">未找到相关位置</text>
			<text class="empty-hint">请尝试输入更详细的地址信息</text>
		</view>
	</view>
</template>

<script>
import api from '../../utils/api.js'

export default {
	data() {
		return {
			// 页面类型：start(起点) 或 end(终点)
			pageType: 'start',
			placeholder: '你要去哪儿?',
			
			// 搜索相关
			searchKeyword: '',
			searchResults: [],
			searchTimer: null,
			
			// 位置信息
			currentLocation: '',
			currentLng: 0,
			currentLat: 0,
			
			// 收藏地址
			favoriteLocations: [],
			
			// 热门地点
			hotLocations: [],
			
			// 用户信息
			userInfo: null,
			
			// 当前城市
			currentCity: '上海'
		}
	},
	
	computed: {
		showEmptyState() {
			return this.searchKeyword && 
				   this.searchResults.length === 0 && 
				   this.favoriteLocations.length === 0 && 
				   this.hotLocations.length === 0;
		}
	},
	
	onLoad(options) {
		// 获取页面类型
		this.pageType = options.type || 'end';
		this.placeholder = this.pageType === 'start' ? '从哪里出发?' : '你要去哪儿?';
		
		// 获取用户信息
		this.userInfo = uni.getStorageSync('userInfo');
		
		// 获取当前位置
		this.getCurrentLocation();
		
		// 加载收藏地址
		this.loadFavoriteLocations();
		
		// 加载热门地点
		this.loadHotLocations();
	},
	
	methods: {
		// 获取当前位置
		getCurrentLocation() {
			uni.getLocation({
				type: 'gcj02',
				success: (res) => {
					this.currentLat = res.latitude;
					this.currentLng = res.longitude;
					this.reverseGeocode(res.latitude, res.longitude);
				},
				fail: (err) => {
					console.error('获取位置失败:', err);
				}
			});
		},
		
		// 逆地理编码
		async reverseGeocode(lat, lng) {
			try {
				const response = await api.reverseGeocoding(lat, lng);
				this.currentLocation = response.data.formatted_address;
			} catch (error) {
				console.error('逆地理编码失败:', error);
			}
		},
		
		// 搜索输入
		onSearchInput(e) {
			const value = e.detail.value;
			this.searchKeyword = value;
			
			// 防抖搜索
			if (this.searchTimer) {
				clearTimeout(this.searchTimer);
			}
			
			if (value.trim()) {
				this.searchTimer = setTimeout(() => {
					this.searchLocation(value);
				}, 500);
			} else {
				this.searchResults = [];
			}
		},
		
		// 搜索确认
		onSearchConfirm(e) {
			const value = e.detail.value;
			if (value.trim()) {
				this.searchLocation(value);
			}
		},
		
		// 搜索位置
		async searchLocation(keyword) {
			try {
				const response = await api.searchAddress(
					keyword,
					this.currentCity,
					this.currentLng,
					this.currentLat,
					20
				);
				this.searchResults = response.data || [];
			} catch (error) {
				console.error('搜索地址失败:', error);
				this.searchResults = [];
			}
		},
		
		// 加载收藏地址
		async loadFavoriteLocations() {
			if (!this.userInfo?.id) return;
			
			try {
				const response = await api.getFavoriteLocations(this.userInfo.id);
				this.favoriteLocations = response.data || [];
			} catch (error) {
				console.error('加载收藏地址失败:', error);
			}
		},
		
		// 加载热门地点
		async loadHotLocations() {
			try {
				const response = await api.getHotLocations(this.currentCity, '', 10);
				this.hotLocations = response.data || [];
			} catch (error) {
				console.error('加载热门地点失败:', error);
			}
		},
		
		// 选择当前位置
		selectCurrentLocation() {
			const locationData = {
				name: '当前位置',
				address: this.currentLocation,
				lng: this.currentLng,
				lat: this.currentLat
			};
			this.returnLocation(locationData);
		},
		
		// 选择搜索结果位置
		selectLocation(item) {
			const locationData = {
				name: item.name,
				address: item.address,
				lng: item.lng,
				lat: item.lat
			};
			this.returnLocation(locationData);
		},
		
		// 选择收藏地址
		selectFavoriteLocation(item) {
			const locationData = {
				name: item.locationName,
				address: item.address,
				lng: item.lng,
				lat: item.lat
			};
			this.returnLocation(locationData);
		},
		
		// 选择热门地点
		selectHotLocation(item) {
			const locationData = {
				name: item.locationName,
				address: item.address,
				lng: item.lng,
				lat: item.lat
			};
			this.returnLocation(locationData);
		},
		
		// 返回选中的位置
		returnLocation(locationData) {
			// 使用uni.$emit发送事件，或者通过页面参数传递
			const pages = getCurrentPages();
			const prevPage = pages[pages.length - 2];
			
			if (prevPage) {
				// 调用上一页的方法
				if (this.pageType === 'start') {
					prevPage.$vm.startLocation = locationData.address;
					prevPage.$vm.startLng = locationData.lng;
					prevPage.$vm.startLat = locationData.lat;
				} else {
					prevPage.$vm.endLocation = locationData.address;
					prevPage.$vm.endLng = locationData.lng;
					prevPage.$vm.endLat = locationData.lat;
				}
			}
			
			// 返回上一页
			uni.navigateBack();
		},
		
		// 删除收藏地址
		async deleteFavorite(item) {
			try {
				await api.deleteFavoriteLocation(this.userInfo.id, item.id);
				uni.showToast({
					title: '删除成功',
					icon: 'success'
				});
				// 重新加载收藏地址
				this.loadFavoriteLocations();
			} catch (error) {
				console.error('删除收藏地址失败:', error);
			}
		},
		
		// 获取收藏地址图标
		getFavoriteIcon(type) {
			const iconMap = {
				'home': '/static/location-home.png',
				'company': '/static/location-company.png',
				'custom': '/static/location-favorite.png'
			};
			return iconMap[type] || '/static/location-favorite.png';
		},
		
		// 格式化距离
		formatDistance(distance) {
			if (distance < 1000) {
				return `${Math.round(distance)}m`;
			} else {
				return `${(distance / 1000).toFixed(1)}km`;
			}
		},
		
		// 返回上一页
		goBack() {
			uni.navigateBack();
		}
	}
}
</script>

<style>
.container {
	background-color: #f8f9fa;
	min-height: 100vh;
}

/* 顶部搜索栏 */
.search-header {
	background: white;
	padding: 20rpx;
	box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.1);
}

.search-box {
	display: flex;
	align-items: center;
	background: #f5f5f5;
	border-radius: 24rpx;
	padding: 16rpx 24rpx;
}

.search-icon {
	width: 32rpx;
	height: 32rpx;
	margin-right: 16rpx;
}

.search-input {
	flex: 1;
	font-size: 28rpx;
	color: #333;
	background: transparent;
}

.cancel-btn {
	color: #4CAF50;
	font-size: 28rpx;
	margin-left: 16rpx;
}

/* 当前位置 */
.current-location {
	background: white;
	margin-top: 16rpx;
}

/* 区域标题 */
.section-title {
	font-size: 24rpx;
	color: #999;
	padding: 24rpx 32rpx 16rpx;
	background: #f8f9fa;
}

/* 位置列表 */
.location-list {
	background: white;
}

.location-item {
	display: flex;
	align-items: center;
	padding: 24rpx 32rpx;
	border-bottom: 1rpx solid #f0f0f0;
}

.location-item:last-child {
	border-bottom: none;
}

.location-icon {
	width: 40rpx;
	height: 40rpx;
	margin-right: 24rpx;
}

.location-info {
	flex: 1;
}

.location-name {
	display: block;
	font-size: 32rpx;
	color: #333;
	margin-bottom: 8rpx;
}

.location-address {
	font-size: 24rpx;
	color: #999;
}

.arrow-icon {
	width: 20rpx;
	height: 20rpx;
}

.location-distance {
	font-size: 24rpx;
	color: #666;
	margin-right: 16rpx;
}

.hot-score {
	font-size: 20rpx;
	color: #FF4444;
	margin-right: 16rpx;
}

/* 收藏地址操作 */
.favorite-actions {
	display: flex;
	align-items: center;
}

.action-icon {
	width: 32rpx;
	height: 32rpx;
	margin-left: 16rpx;
}

/* 空状态 */
.empty-state {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 120rpx 60rpx;
	text-align: center;
}

.empty-icon {
	width: 120rpx;
	height: 120rpx;
	margin-bottom: 32rpx;
	opacity: 0.3;
}

.empty-text {
	font-size: 32rpx;
	color: #666;
	margin-bottom: 16rpx;
}

.empty-hint {
	font-size: 24rpx;
	color: #999;
}

/* 搜索结果区域 */
.search-results {
	margin-top: 16rpx;
}

.favorite-locations {
	margin-top: 16rpx;
}

.hot-locations {
	margin-top: 16rpx;
}
</style> 