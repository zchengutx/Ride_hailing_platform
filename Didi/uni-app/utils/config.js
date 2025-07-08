// API配置
const config = {
	// 基础URL配置
	baseURL: 'https://api.didi.com', // 替换为实际的API地址
	
	// 接口路径配置
	api: {
		// 用户相关接口
		sendSms: '/user/sendSms',           // 发送验证码
		login: '/user/login',               // 用户登录
		realName: '/user/realName',         // 实名认证
		takeACar: '/user/takeACar',         // 用户叫车
		
		// 司机相关接口
		callACar: '/driver/callACar',       // 司机申请
		driverAudit: '/driver/driverAudit', // 司机审核
		addDriver: '/driver/addDriver',     // 司机添加
		driverOnline: '/driver/driverOnline', // 司机状态管理
		receivingOrder: '/driver/receivingOrder' // 司机接单
	},
	
	// 请求超时时间
	timeout: 10000,
	
	// 默认headers
	headers: {
		'Content-Type': 'application/json'
	}
}

export default config 