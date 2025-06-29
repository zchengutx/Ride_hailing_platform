package router

import (
	"log"
	"lxh_cx/biz/handler"
	"lxh_cx/utils"

	"github.com/cloudwego/hertz/pkg/route"
)

func UserGroup(a *route.RouterGroup) {
	// 创建 JWT 中间件
	authMiddleware, err := utils.NewJWTMiddleware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	user := a.Group("/user")
	{

		user.POST("sendSms", handler.SendSms)
		user.POST("register", handler.Register)
		user.GET("weChat", handler.WeChatLogin)
		user.GET("checkSignature", handler.CheckSignature)
		user.GET("callback", handler.Callback)

		user.POST("login", authMiddleware.LoginHandler)
		auth := user.Group("")
		auth.Use(authMiddleware.MiddlewareFunc())
		{
			auth.GET("userInfoList", handler.UserInfoList)
			auth.POST("bindMobile", handler.BindMobile)
			auth.POST("updateCancelOrder", handler.UpdateCancelOrder)
			auth.POST("passengerAddOrder", handler.PassengerAddOrder)
		}
	}
}
