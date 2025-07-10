<template>
	<view class="container">
		<view class="header">
			<text class="title">我的订单</text>
		</view>
		
		<view class="filter-tabs">
			<view 
				class="tab-item"
				:class="{ active: activeTab === 'all' }"
				@tap="switchTab('all')"
			>
				全部
			</view>
			<view 
				class="tab-item"
				:class="{ active: activeTab === 'pending' }"
				@tap="switchTab('pending')"
			>
				进行中
			</view>
			<view 
				class="tab-item"
				:class="{ active: activeTab === 'completed' }"
				@tap="switchTab('completed')"
			>
				已完成
			</view>
			<view 
				class="tab-item"
				:class="{ active: activeTab === 'cancelled' }"
				@tap="switchTab('cancelled')"
			>
				已取消
			</view>
		</view>
		
		<view class="order-list">
			<view 
				class="order-item"
				v-for="(order, index) in orderList"
				:key="index"
				@tap="goToDetail(order)"
			>
				<view class="order-header">
					<text class="order-time">{{ order.startTime }}</text>
					<text class="order-status" :class="getStatusClass(order.orderStatus)">
						{{ order.orderStatus }}
					</text>
				</view>
				
				<view class="route-info">
					<view class="location">
						<view class="dot start-dot"></view>
						<text class="address">{{ order.startAddr }}</text>
					</view>
					<view class="location">
						<view class="dot end-dot"></view>
						<text class="address">{{ order.endEnd }}</text>
					</view>
				</view>
				
				<view class="order-footer">
					<text class="order-code">{{ order.orderCode }}</text>
					<text class="order-amount">¥{{ order.amount }}</text>
				</view>
			</view>
		</view>
		
		<view class="empty-state" v-if="orderList.length === 0">
			<image src="/static/empty-order.png" class="empty-icon" mode="aspectFit"></image>
			<text class="empty-text">暂无订单记录</text>
		</view>
	</view>
</template>

<script>
import api from '../../utils/api.js'

export default {
	data() {
		return {
			activeTab: 'all',
			orderList: [],
			userInfo: null
		}
	},
	
	onShow() {
		this.userInfo = uni.getStorageSync('userInfo');
		this.loadOrderList();
	},
	
	methods: {
		// 切换标签
		switchTab(tab) {
			this.activeTab = tab;
			this.loadOrderList();
		},
		
		// 加载订单列表
		async loadOrderList() {
			if (!this.userInfo?.id) {
				uni.showToast({
					title: '请先登录',
					icon: 'none'
				});
				return;
			}
			
			try {
				const status = this.activeTab === 'all' ? '' : this.activeTab;
				const response = await api.getOrderList(this.userInfo.id, 1, 20, status);
				this.orderList = response.data || [];
			} catch (error) {
				console.error('加载订单列表失败:', error);
				this.orderList = [
					// 模拟数据
					{
						id: 1,
						orderCode: 'DL202312120001',
						startTime: '2023-12-12 14:30',
						orderStatus: '已完成',
						startAddr: '东禄素质教育产业园',
						endEnd: '上海市浦东医院',
						amount: '25.80'
					},
					{
						id: 2,
						orderCode: 'DL202312120002',
						startTime: '2023-12-11 09:15',
						orderStatus: '已取消',
						startAddr: '人民广场',
						endEnd: '虹桥机场T2',
						amount: '68.50'
					}
				];
			}
		},
		
		// 获取状态样式
		getStatusClass(status) {
			const statusMap = {
				'已完成': 'completed',
				'进行中': 'pending',
				'已取消': 'cancelled',
				'等待接单': 'waiting'
			};
			return statusMap[status] || 'default';
		},
		
		// 前往订单详情
		goToDetail(order) {
			uni.navigateTo({
				url: `/pages/order/detail?orderCode=${order.orderCode}`
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

.header {
	background: white;
	text-align: center;
	padding: 40rpx 0;
	box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.1);
}

.title {
	font-size: 36rpx;
	font-weight: bold;
	color: #333;
}

.filter-tabs {
	background: white;
	display: flex;
	padding: 0 20rpx;
}

.tab-item {
	flex: 1;
	text-align: center;
	padding: 24rpx 0;
	font-size: 28rpx;
	color: #666;
	border-bottom: 4rpx solid transparent;
}

.tab-item.active {
	color: #4CAF50;
	border-bottom-color: #4CAF50;
}

.order-list {
	padding: 20rpx;
}

.order-item {
	background: white;
	border-radius: 16rpx;
	padding: 32rpx;
	margin-bottom: 16rpx;
	box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.1);
}

.order-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 24rpx;
}

.order-time {
	font-size: 28rpx;
	color: #666;
}

.order-status {
	font-size: 28rpx;
	font-weight: bold;
}

.order-status.completed {
	color: #4CAF50;
}

.order-status.pending {
	color: #FF9800;
}

.order-status.cancelled {
	color: #999;
}

.order-status.waiting {
	color: #2196F3;
}

.route-info {
	margin-bottom: 24rpx;
}

.location {
	display: flex;
	align-items: center;
	margin-bottom: 16rpx;
}

.dot {
	width: 12rpx;
	height: 12rpx;
	border-radius: 50%;
	margin-right: 16rpx;
}

.start-dot {
	background: #4CAF50;
}

.end-dot {
	background: #FF9800;
}

.address {
	font-size: 28rpx;
	color: #333;
}

.order-footer {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding-top: 16rpx;
	border-top: 1rpx solid #f0f0f0;
}

.order-code {
	font-size: 24rpx;
	color: #999;
}

.order-amount {
	font-size: 32rpx;
	font-weight: bold;
	color: #4CAF50;
}

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
}
</style> 