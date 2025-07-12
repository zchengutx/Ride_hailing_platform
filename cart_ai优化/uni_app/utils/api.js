// API配置和服务文件
const config = {
	// 根据实际后端部署地址修改
	baseURL: 'http://localhost:8080/v1',
	timeout: 10000
}

/**
 * 通用请求方法
 * @param {string} url 请求地址
 * @param {object} data 请求数据
 * @param {string} method 请求方法
 * @param {object} customHeader 自定义请求头
 */
function request(url, data = {}, method = 'POST', customHeader = {}) {
	return new Promise((resolve, reject) => {
		// 获取存储的token
		const token = uni.getStorageSync('token');
		
		// 默认请求头
		const header = {
			'Content-Type': 'application/json',
			...customHeader
		};
		
		// 如果有token，添加到请求头
		if (token) {
			header['x-token'] = token;
		}
		
		// 发起请求
		uni.request({
			url: config.baseURL + url,
			method: method,
			data: data,
			header: header,
			timeout: config.timeout,
			success: (res) => {
				const statusCode = res.statusCode;
				const responseData = res.data;
				
				// 请求成功
				if (statusCode === 200) {
					// 检查业务状态码
					if (responseData.code === 200) {
						resolve(responseData);
					} else {
						// 业务错误
						uni.showToast({
							title: responseData.message || '请求失败',
							icon: 'none'
						});
						reject(responseData);
					}
				} else {
					// HTTP错误
					uni.showToast({
						title: `网络错误 ${statusCode}`,
						icon: 'none'
					});
					reject(res);
				}
			},
			fail: (err) => {
				console.error('请求失败:', err);
				uni.showToast({
					title: '网络连接失败',
					icon: 'none'
				});
				reject(err);
			}
		});
	});
}

// ============ 乘客服务API ============

/**
 * 发送短信验证码
 * @param {string} mobile 手机号
 * @param {string} sendSmsCode 验证码类型
 */
export function sendSms(mobile, sendSmsCode = 'login') {
	return request('/api/passenger/sendSms', {
		mobile,
		sendSmsCode
	});
}

/**
 * 乘客注册
 * @param {string} mobile 手机号
 * @param {string} sendSmsCode 验证码
 */
export function registerPassenger(mobile, sendSmsCode) {
	return request('/api/passenger/registerPassenger', {
		mobile,
		sendSmsCode
	});
}

/**
 * 乘客登录
 * @param {string} mobile 手机号
 * @param {string} sendSmsCode 验证码
 */
export function loginPassenger(mobile, sendSmsCode) {
	return request('/api/passenger/loginPassenger', {
		mobile,
		sendSmsCode
	});
}

/**
 * 获取首页数据
 * @param {number} passengerId 乘客ID
 * @param {string} location 当前位置
 */
export function getHomePage(passengerId, location) {
	return request('/api/passenger/homePage', {
		passengerId,
		location
	});
}

/**
 * 获取乘客信息
 * @param {number} passengerId 乘客ID
 */
export function getPassengerInfo(passengerId) {
	return request('/api/passenger/info', { passengerId }, 'GET');
}

/**
 * 更新乘客信息
 * @param {object} updateData 更新数据
 */
export function updatePassengerInfo(updateData) {
	return request('/api/passenger/info', updateData, 'PUT');
}

/**
 * 创建订单
 * @param {object} orderData 订单数据
 */
export function createOrder(orderData) {
	return request('/api/passenger/createOrder', orderData);
}

/**
 * 获取订单列表
 * @param {number} passengerId 乘客ID
 * @param {number} page 页码
 * @param {number} pageSize 每页数量
 * @param {string} status 订单状态
 */
export function getOrderList(passengerId, page = 1, pageSize = 10, status = '') {
	return request('/api/passenger/orders', {
		passengerId,
		page,
		pageSize,
		status
	}, 'GET');
}

/**
 * 获取订单详情
 * @param {number} passengerId 乘客ID
 * @param {number} orderId 订单ID
 */
export function getOrderDetail(passengerId, orderId) {
	return request(`/api/passenger/order/${orderId}`, { passengerId }, 'GET');
}

/**
 * 取消订单
 * @param {number} passengerId 乘客ID
 * @param {number} orderId 订单ID
 * @param {string} reason 取消原因
 * @param {string} remark 备注
 */
export function cancelOrder(passengerId, orderId, reason, remark = '') {
	return request('/api/passenger/cancelOrder', {
		passengerId,
		orderId,
		reason,
		remark
	});
}

/**
 * 获取收藏地址
 * @param {number} passengerId 乘客ID
 * @param {string} locationType 地址类型
 */
export function getFavoriteLocations(passengerId, locationType = '') {
	return request('/api/passenger/favoriteLocations', {
		passengerId,
		locationType
	}, 'GET');
}

/**
 * 添加收藏地址
 * @param {object} locationData 地址数据
 */
export function addFavoriteLocation(locationData) {
	return request('/api/passenger/favoriteLocation', locationData);
}

/**
 * 删除收藏地址
 * @param {number} passengerId 乘客ID
 * @param {number} locationId 地址ID
 */
export function deleteFavoriteLocation(passengerId, locationId) {
	return request(`/api/passenger/favoriteLocation/${locationId}`, {
		passengerId
	}, 'DELETE');
}

/**
 * 获取热门地点
 * @param {string} city 城市
 * @param {string} category 分类
 * @param {number} limit 数量限制
 */
export function getHotLocations(city, category = '', limit = 10) {
	return request('/api/passenger/hotLocations', {
		city,
		category,
		limit
	}, 'GET');
}

/**
 * 地址搜索建议
 * @param {string} keyword 搜索关键词
 * @param {string} city 城市
 * @param {number} lng 经度
 * @param {number} lat 纬度
 * @param {number} limit 数量限制
 */
export function searchAddress(keyword, city = '', lng = 0, lat = 0, limit = 10) {
	return request('/api/passenger/searchAddress', {
		keyword,
		city,
		lng,
		lat,
		limit
	});
}

// ============ 地图服务API ============

/**
 * 地理编码 - 根据地址获取坐标
 * @param {string} address 地址
 */
export function geocoding(address) {
	return request('/api/map/geocoding', { address });
}

/**
 * 逆地理编码 - 根据坐标获取地址
 * @param {number} lat 纬度
 * @param {number} lng 经度
 */
export function reverseGeocoding(lat, lng) {
	return request('/api/map/reverse-geocoding', { lat, lng });
}

/**
 * IP定位
 * @param {string} ip IP地址
 */
export function ipLocation(ip) {
	return request('/api/map/ip-location', { ip });
}

/**
 * 距离计算
 * @param {number} originLat 起点纬度
 * @param {number} originLng 起点经度
 * @param {number} destLat 终点纬度
 * @param {number} destLng 终点经度
 */
export function calculateDistance(originLat, originLng, destLat, destLng) {
	return request('/api/map/distance', {
		origin_lat: originLat,
		origin_lng: originLng,
		dest_lat: destLat,
		dest_lng: destLng
	});
}

/**
 * 获取当前位置
 * @param {string} clientIp 客户端IP
 */
export function getCurrentLocation(clientIp) {
	return request('/api/map/current-location', { client_ip: clientIp }, 'GET');
}

/**
 * 获取省份列表
 */
export function getProvinces() {
	return request('/api/map/provinces', {}, 'GET');
}

/**
 * 获取城市列表
 * @param {number} provinceCode 省份代码
 */
export function getCities(provinceCode) {
	return request('/api/map/cities', { province_code: provinceCode });
}

/**
 * 获取区县列表
 * @param {number} cityCode 城市代码
 */
export function getDistricts(cityCode) {
	return request('/api/map/districts', { city_code: cityCode });
}

// ============ 微信服务API ============

/**
 * 绑定微信账号
 * @param {number} passengerId 乘客ID
 * @param {string} code 微信授权码
 */
export function bindWechat(passengerId, code) {
	return request('/api/passenger/bindWechat', {
		passengerId,
		code
	});
}

/**
 * 解绑微信账号
 * @param {number} passengerId 乘客ID
 */
export function unbindWechat(passengerId) {
	return request('/api/passenger/unbindWechat', { passengerId });
}

/**
 * 获取微信绑定状态
 * @param {number} passengerId 乘客ID
 */
export function getWechatBindStatus(passengerId) {
	return request('/api/passenger/getWechatBindStatus', { passengerId }, 'GET');
}

// ============ 健康检查API ============

/**
 * 健康检查
 */
export function healthCheck() {
	return request('/api/health', {}, 'GET');
}

// ============ 工具方法 ============

/**
 * 设置API基础地址
 * @param {string} baseURL 基础地址
 */
export function setBaseURL(baseURL) {
	config.baseURL = baseURL;
}

/**
 * 设置请求超时时间
 * @param {number} timeout 超时时间(毫秒)
 */
export function setTimeout(timeout) {
	config.timeout = timeout;
}

/**
 * 获取当前配置
 */
export function getConfig() {
	return { ...config };
}

// 默认导出配置对象
export default {
	request,
	setBaseURL,
	setTimeout,
	getConfig,
	// 乘客服务
	sendSms,
	registerPassenger,
	loginPassenger,
	getHomePage,
	getPassengerInfo,
	updatePassengerInfo,
	createOrder,
	getOrderList,
	getOrderDetail,
	cancelOrder,
	getFavoriteLocations,
	addFavoriteLocation,
	deleteFavoriteLocation,
	getHotLocations,
	searchAddress,
	// 地图服务
	geocoding,
	reverseGeocoding,
	ipLocation,
	calculateDistance,
	getCurrentLocation,
	getProvinces,
	getCities,
	getDistricts,
	// 微信服务
	bindWechat,
	unbindWechat,
	getWechatBindStatus,
	// 健康检查
	healthCheck
}; 