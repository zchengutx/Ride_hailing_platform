import request from '../utils/request.js'
import config from '../utils/config.js'

/**
 * 用户相关API
 */

/**
 * 发送短信验证码
 * @param {String} mobile 手机号
 * @param {String} sendSmsCode 验证码类型
 */
export function sendSms(mobile, sendSmsCode = 'register') {
	return request.post(config.api.sendSms, {
		mobile,
		send_sms_code: sendSmsCode
	})
}

/**
 * 用户登录
 * @param {String} mobile 手机号
 * @param {String} smsCode 验证码
 */
export function login(mobile, smsCode) {
	return request.post(config.api.login, {
		mobile,
		send_sms_code: smsCode
	})
}

/**
 * 实名认证
 * @param {Object} realNameData 实名认证数据
 */
export function realName(realNameData) {
	return request.post(config.api.realName, {
		user_name: realNameData.userName,
		sex: realNameData.sex,
		age: realNameData.age,
		id_card: realNameData.idCard
	})
}

/**
 * 用户叫车
 * @param {Object} orderData 订单数据
 */
export function takeACar(orderData) {
	return request.post(config.api.takeACar, {
		start_location: orderData.startLocation,
		end_location: orderData.endLocation,
		cart_type: orderData.cartType
	})
} 