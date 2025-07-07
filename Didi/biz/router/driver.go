package router

import (
	"Didi/biz/dal/global"
	"Didi/biz/handler/api"
	"Didi/biz/middleware"
	"github.com/cloudwego/hertz/pkg/route"
)

func DriverModel(r *route.RouterGroup) {
	driverModel := r.Group("/driver")
	{
		driverModel.Use(middleware.JWTAuth(global.JWT_SELECT_KEY))
		driverModel.POST("/callACar", api.CallACar)
	}
}
