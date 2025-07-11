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
		passengerModel.POST("/sendSms", api.SendSms)                     // 发送短信验证码
		passengerModel.POST("/registerPassenger", api.RegisterPassenger) // 乘客注册
		passengerModel.POST("/loginPassenger", api.LoginPassenger)       // 乘客登录
	}
	//jwt中间件
	passengerAuth := passengerModel.Group("")
	passengerAuth.Use(middleware.JWTAuth(global.JWT_SELECT_KEY))
	{
		// 个人信息管理
		passengerAuth.GET("/info", api.GetPassengerInfo)    // 获取乘客信息
		passengerAuth.PUT("/info", api.UpdatePassengerInfo) // 更新乘客信息

		// 订单管理
		passengerAuth.POST("/createOrder", api.CreateOrder)          // 创建订单
		passengerAuth.GET("/orders", api.GetPassengerOrders)         // 获取订单列表
		passengerAuth.GET("/order/:orderId", api.GetOrderDetail)     // 获取订单详情
		passengerAuth.POST("/cancelOrder", api.PassengerCancelOrder) // 取消订单
		passengerAuth.POST("/evaluateOrder", api.EvaluateOrder)      // 评价订单

		// 收藏地址管理
		passengerAuth.GET("/favoriteLocations", api.GetFavoriteLocations)                 // 获取收藏地址
		passengerAuth.POST("/favoriteLocation", api.AddFavoriteLocation)                  // 添加收藏地址
		passengerAuth.DELETE("/favoriteLocation/:locationId", api.DeleteFavoriteLocation) // 删除收藏地址

		// 热门地点推荐
		passengerAuth.GET("/hotLocations", api.GetHotLocations) // 获取热门地点

		// 微信绑定管理
		passengerAuth.POST("/bindWechat", api.BindWechat)               // 绑定微信账号
		passengerAuth.POST("/unbindWechat", api.UnbindWechat)           // 解绑微信账号
		passengerAuth.GET("/wechatBindStatus", api.GetWechatBindStatus) // 获取微信绑定状态

		// 路线记录查询
		passengerAuth.GET("/routeRecords", api.GetRouteRecords) // 获取路线记录

		// 地址搜索建议
		passengerAuth.GET("/searchAddress", api.SearchAddress) // 地址搜索建议

		// 主页服务（需要认证以获取个性化内容）
		passengerAuth.GET("/homePage", api.HomePage)  // 主页服务
		passengerAuth.POST("/callACar", api.CallACar) // 叫车服务
	}
}
