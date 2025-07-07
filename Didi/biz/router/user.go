package router

import (
	"Didi/biz/dal/global"
	"Didi/biz/handler/api"
	"Didi/biz/middleware"
	"github.com/cloudwego/hertz/pkg/route"
)

func UserModel(r *route.RouterGroup) {
	userModel := r.Group("/user")
	{
		userModel.POST("/sendSms", api.SendSms)
		userModel.POST("/login", api.Login)
		userModel.Use(middleware.JWTAuth(global.JWT_SELECT_KEY))
		userModel.POST("/realName", api.RealName)
	}
}
