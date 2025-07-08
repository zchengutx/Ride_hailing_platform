import request from '../utils/request.js'
import config from '../utils/config.js'

/**
 * 司机相关API
 */

/**
 * 司机申请
 * @param {Object} driverData 司机申请数据
 */
export function callACar(driverData) {
	return request.post(config.api.callACar, {
		address: driverData.address,
		driving_license_number: driverData.drivingLicenseNumber,
		quasi_driving_type: driverData.quasiDrivingType,
		driving_age: driverData.drivingAge,
		car_num: driverData.carNum,
		car_type: driverData.carType,
		vehicle_mileage: driverData.vehicleMileage,
		serving_the_city: driverData.servingTheCity
	})
}

/**
 * 司机审核
 * @param {Number} driverId 司机ID
 * @param {String} auditStatus 审核状态
 */
export function driverAudit(driverId, auditStatus) {
	return request.post(config.api.driverAudit, {
		driver_id: driverId,
		audit_status: auditStatus
	})
}

/**
 * 司机添加
 * @param {Number} driverId 司机ID
 */
export function addDriver(driverId) {
	return request.post(config.api.addDriver, {
		driver_id: driverId
	})
}

/**
 * 司机状态管理
 * @param {Number} driverId 司机ID
 * @param {String} driverStatus 司机状态
 */
export function driverOnline(driverId, driverStatus) {
	return request.post(config.api.driverOnline, {
		driver_id: driverId,
		driver_status: driverStatus
	})
}

/**
 * 司机接单
 * @param {Object} orderData 接单数据
 */
export function receivingOrder(orderData) {
	return request.post(config.api.receivingOrder, {
		order_id: orderData.orderId,
		current_latitude: orderData.currentLatitude,
		current_longitude: orderData.currentLongitude
	})
} 