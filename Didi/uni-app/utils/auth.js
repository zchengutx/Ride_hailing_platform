/**
 * 用户认证工具类
 */

const TOKEN_KEY = 'didi_token'
const USER_INFO_KEY = 'didi_user_info'

/**
 * 获取token
 */
export function getToken() {
	return uni.getStorageSync(TOKEN_KEY)
}

/**
 * 设置token
 * @param {String} token 
 */
export function setToken(token) {
	return uni.setStorageSync(TOKEN_KEY, token)
}

/**
 * 移除token
 */
export function removeToken() {
	return uni.removeStorageSync(TOKEN_KEY)
}

/**
 * 获取用户信息
 */
export function getUserInfo() {
	const userInfo = uni.getStorageSync(USER_INFO_KEY)
	return userInfo ? JSON.parse(userInfo) : null
}

/**
 * 设置用户信息
 * @param {Object} userInfo 
 */
export function setUserInfo(userInfo) {
	return uni.setStorageSync(USER_INFO_KEY, JSON.stringify(userInfo))
}

/**
 * 移除用户信息
 */
export function removeUserInfo() {
	return uni.removeStorageSync(USER_INFO_KEY)
}

/**
 * 检查是否已登录
 */
export function isLoggedIn() {
	return !!getToken()
}

/**
 * 清除所有认证信息
 */
export function clearAuth() {
	removeToken()
	removeUserInfo()
}

/**
 * 检查是否需要登录
 * @param {Function} callback 需要登录时的回调
 */
export function checkLogin(callback) {
	if (!isLoggedIn()) {
		uni.showModal({
			title: '提示',
			content: '请先登录',
			success: (res) => {
				if (res.confirm) {
					uni.navigateTo({
						url: '/pages/login/login'
					})
				}
			}
		})
		return false
	}
	
	if (callback && typeof callback === 'function') {
		callback()
	}
	return true
} 