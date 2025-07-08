import config from './config.js'
import { getToken, removeToken } from './auth.js'

/**
 * 网络请求工具类
 */
class Request {
	constructor() {
		this.baseURL = config.baseURL
		this.timeout = config.timeout
		this.header = config.headers
	}

	/**
	 * 发起请求
	 * @param {Object} options 请求配置
	 */
	request(options = {}) {
		// 显示加载提示
		if (options.loading !== false) {
			uni.showLoading({
				title: '加载中...',
				mask: true
			})
		}

		return new Promise((resolve, reject) => {
			// 构建请求头
			let header = { ...this.header }
			
			// 添加token
			const token = getToken()
			if (token) {
				header.Authorization = `Bearer ${token}`
			}

			// 构建请求参数
			const requestOptions = {
				url: this.baseURL + options.url,
				method: options.method || 'GET',
				data: options.data || {},
				header: { ...header, ...options.header },
				timeout: options.timeout || this.timeout,
				success: (res) => {
					this.handleSuccess(res, resolve, reject)
				},
				fail: (err) => {
					this.handleError(err, reject)
				},
				complete: () => {
					if (options.loading !== false) {
						uni.hideLoading()
					}
				}
			}

			// 发起请求
			uni.request(requestOptions)
		})
	}

	/**
	 * 处理请求成功
	 */
	handleSuccess(res, resolve, reject) {
		const { statusCode, data } = res

		if (statusCode === 200) {
			// 根据业务状态码处理
			if (data.code === 200) {
				resolve(data)
			} else if (data.code === 401) {
				// token过期，跳转登录
				this.handleTokenExpired()
				reject(data)
			} else {
				// 业务错误
				this.showError(data.msg || '请求失败')
				reject(data)
			}
		} else {
			// HTTP状态码错误
			this.showError(`请求失败，状态码：${statusCode}`)
			reject({ code: statusCode, msg: '网络错误' })
		}
	}

	/**
	 * 处理请求失败
	 */
	handleError(err, reject) {
		console.error('请求失败：', err)
		let message = '网络异常，请检查网络连接'
		
		if (err.errMsg) {
			if (err.errMsg.includes('timeout')) {
				message = '请求超时，请重试'
			} else if (err.errMsg.includes('fail')) {
				message = '网络连接失败'
			}
		}
		
		this.showError(message)
		reject({ code: -1, msg: message })
	}

	/**
	 * 显示错误信息
	 */
	showError(message) {
		uni.showToast({
			title: message,
			icon: 'none',
			duration: 2000
		})
	}

	/**
	 * 处理token过期
	 */
	handleTokenExpired() {
		removeToken()
		uni.showModal({
			title: '提示',
			content: '登录已过期，请重新登录',
			showCancel: false,
			success: () => {
				uni.reLaunch({
					url: '/pages/login/login'
				})
			}
		})
	}

	/**
	 * GET请求
	 */
	get(url, data = {}, options = {}) {
		return this.request({
			url,
			method: 'GET',
			data,
			...options
		})
	}

	/**
	 * POST请求
	 */
	post(url, data = {}, options = {}) {
		return this.request({
			url,
			method: 'POST',
			data,
			...options
		})
	}

	/**
	 * PUT请求
	 */
	put(url, data = {}, options = {}) {
		return this.request({
			url,
			method: 'PUT',
			data,
			...options
		})
	}

	/**
	 * DELETE请求
	 */
	delete(url, data = {}, options = {}) {
		return this.request({
			url,
			method: 'DELETE',
			data,
			...options
		})
	}
}

// 创建请求实例
const request = new Request()

export default request 