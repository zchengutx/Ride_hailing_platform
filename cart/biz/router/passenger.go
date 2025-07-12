package router

import (
	"cart/biz/dal/global"
	"cart/biz/handler/api"
	"cart/biz/middleware"
	"log"

	"github.com/cloudwego/hertz/pkg/route"
)

func PassengerModel(r *route.RouterGroup) {
	log.Println("乘客服务")
	passengerModel := r.Group("/passenger")
	{
		log.Println("短信验证码:::::/v1/passenger/sendSms")
		passengerModel.POST("/sendSms", api.SendSms) // 发送短信验证码
		log.Println("乘客注册:::::/v1/passenger/registerPassenger")
		passengerModel.POST("/registerPassenger", api.RegisterPassenger) // 乘客注册
		log.Println("乘客登录:::::/v1/passenger/loginPassenger")
		passengerModel.POST("/loginPassenger", api.LoginPassenger) // 乘客登录
	}

	// jwt
	passengerAuth := passengerModel.Group("")
	passengerAuth.Use(middleware.JWTAuth(global.JWT_SELECT_KEY))
	{
		passengerAuth.GET("/info", api.GetPassengerInfo)    // 获取乘客信息
		passengerAuth.PUT("/info", api.UpdatePassengerInfo) // 更新乘客信息

		// 订单管理
		log.Println("创建订单:::::/v1/passenger/createOrder")
		passengerAuth.POST("/createOrder", api.CreateOrder) // 创建订单
		log.Println("获取订单列表:::::/v1/passenger/orders")
		passengerAuth.GET("/orders", api.GetPassengerOrders) // 获取订单列表
		log.Println("获取订单详情:::::/v1/passenger/order/:orderId")
		passengerAuth.GET("/order/:orderId", api.GetOrderDetail) // 获取订单详情
		log.Println("取消订单:::::/v1/passenger/cancelOrder")
		passengerAuth.POST("/cancelOrder", api.PassengerCancelOrder) // 取消订单
		log.Println("评价订单:::::/v1/passenger/evaluateOrder")
		passengerAuth.POST("/evaluateOrder", api.EvaluateOrder) // 评价订单
		// 地址管理
		log.Println("获取收藏地址:::::/v1/passenger/favoriteLocations")
		passengerAuth.GET("/favoriteLocations", api.GetFavoriteLocations) // 获取收藏地址
		log.Println("添加收藏地址:::::/v1/passenger/favoriteLocation")
		passengerAuth.POST("/favoriteLocation", api.AddFavoriteLocation) // 添加收藏地址
		log.Println("删除收藏地址:::::/v1/passenger/favoriteLocation/:locationId")
		passengerAuth.DELETE("/favoriteLocation/:locationId", api.DeleteFavoriteLocation) // 删除收藏地址
		log.Println("获取热门地点:::::/v1/passenger/hotLocations")
		passengerAuth.GET("/hotLocations", api.GetHotLocations) // 获取热门地点
		// 微信绑定管理
		log.Println("绑定微信账号:::::/v1/passenger/bindWechat")
		passengerAuth.POST("/bindWechat", api.BindWechat) // 绑定微信账号
		log.Println("解绑微信账号:::::/v1/passenger/unbindWechat")
		passengerAuth.POST("/unbindWechat", api.UnbindWechat) // 解绑微信账号
		log.Println("获取微信绑定状态:::::/v1/passenger/wechatBindStatus")
		passengerAuth.GET("/wechatBindStatus", api.GetWechatBindStatus) // 获取微信绑定状态
		// 路线记录查询
		log.Println("获取路线记录:::::/v1/passenger/routeRecords")
		passengerAuth.GET("/routeRecords", api.GetRouteRecords) // 获取路线记录
		log.Println("历史地址:::::/v1/passenger/searchAddress")
		passengerAuth.GET("/searchAddress", api.SearchAddress) // 历史地址
		// 主页服务
		log.Println("主页服务:::::/v1/passenger/homePage")
		passengerAuth.GET("/homePage", api.HomePage) // 主页服务
		log.Println("叫车服务:::::/v1/passenger/callACar")
		passengerAuth.POST("/callACar", api.CallACar) // 叫车服务
	}
}
