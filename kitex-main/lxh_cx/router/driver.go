package router

import (
	"github.com/cloudwego/hertz/pkg/route"
	"log"
	"lxh_cx/biz/handler"
	"lxh_cx/utils"
)

func DriverGroup(a *route.RouterGroup) {
	// 创建 JWT 中间件
	authMiddleware, err := utils.NewJWTMiddleware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	driver := a.Group("/driver")
	{
		driver.POST("driverRegister", handler.DriverRegister)
		driver.POST("login", authMiddleware.LoginHandler)
		auth := driver.Group("")
		auth.Use(authMiddleware.MiddlewareFunc())
		{
			auth.GET("driverInfoList", handler.DriverInfoList)
			auth.POST("driverAdd", handler.DriverAdd)
			auth.POST("driverOverOrder", handler.DriverOverOrder)
			auth.POST("driverCancelOrder", handler.DriverCancelOrder)
		}
	}
}
