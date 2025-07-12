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
		log.Println("司机注册:::::/v1/driver/petition")
		driverModel.POST("/petition", api.DriverPetition) // 司机注册申请
		log.Println("司机登录:::::/v1/petition/login")
		driverModel.POST("/login", api.DriverLogin) // 司机登录
	}

	// jwt
	driverAuth := driverModel.Group("")
	driverAuth.Use(middleware.JWTAuth(global.JWT_SELECT_KEY))
	{
		driverAuth.GET("/info", api.GetDriverInfo)    // 获取司机信息
		driverAuth.PUT("/info", api.UpdateDriverInfo) // 更新司机信息

		// 司机状态管理
		log.Println("司机上线/下线:::::/v1/driver/changeStatus")
		driverAuth.POST("/changeStatus", api.ChangeStatus) // 上线/下线

		// 订单管理
		log.Println("获取待接订单:::::/v1/driver/pendingOrders")
		driverAuth.GET("/pendingOrders", api.GetPendingOrders) // 获取待接订单
		log.Println("接单:::::/v1/driver/acceptOrder")
		driverAuth.POST("/acceptOrder", api.AcceptOrder) // 接单
		log.Println("开始行程:::::/v1/driver/startTrip")
		driverAuth.POST("/startTrip", api.StartTrip) // 开始行程
		log.Println("完成订单:::::/v1/driver/completeOrder")
		driverAuth.POST("/completeOrder", api.CompleteOrder) // 完成订单
		log.Println("取消订单:::::/v1/driver/cancelOrder")
		driverAuth.POST("/cancelOrder", api.CancelOrder) // 取消订单
		log.Println("查询收益:::::/v1/driver/income")
		driverAuth.GET("/income", api.GetIncome) // 查询收益

		// 位置服务
		log.Println("更新位置:::::/v1/driver/updateLocation")
		driverAuth.POST("/updateLocation", api.UpdateLocation) // 更新位置
	}
}
