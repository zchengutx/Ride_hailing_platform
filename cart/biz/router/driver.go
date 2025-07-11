package router

import (
	"cart/biz/dal/global"
	"cart/biz/handler/api"
	"cart/biz/middleware"
	"log"

	"github.com/cloudwego/hertz/pkg/route"
)

func DriverModel(r *route.RouterGroup) {
	log.Println("司机服务")
	driverModel := r.Group("/driver")

	{
		driverModel.POST("/petition", api.DriverPetition) // 司机注册申请
		driverModel.POST("/login", api.DriverLogin)       // 司机登录
	}

	// jwt中间件
	driverAuth := driverModel.Group("")
	driverAuth.Use(middleware.JWTAuth(global.JWT_SELECT_KEY))
	{
		// 司机基本信息管理
		driverAuth.GET("/checkStatus", api.CheckStatus) // 查询审核状态
		driverAuth.GET("/info", api.GetDriverInfo)      // 获取司机信息
		driverAuth.PUT("/info", api.UpdateDriverInfo)   // 更新司机信息

		// 司机状态管理
		driverAuth.POST("/changeStatus", api.ChangeStatus) // 上线/下线

		// 订单管理
		driverAuth.GET("/pendingOrders", api.GetPendingOrders) // 获取待接订单
		driverAuth.POST("/acceptOrder", api.AcceptOrder)       // 接单
		driverAuth.POST("/startTrip", api.StartTrip)           // 开始行程
		driverAuth.POST("/completeOrder", api.CompleteOrder)   // 完成订单
		driverAuth.POST("/cancelOrder", api.CancelOrder)       // 取消订单

		// 收益管理
		driverAuth.GET("/income", api.GetIncome) // 查询收益

		// 位置服务
		driverAuth.POST("/updateLocation", api.UpdateLocation) // 更新位置
	}
}
