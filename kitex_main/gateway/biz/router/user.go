package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"kitex_main/gateway/biz/handler"
	"kitex_main/pkg"
)

func SetUserRouter(r *server.Hertz) {
	r.POST("api/sendSms", handler.SendSms)                 //获取手机验证码
	r.POST("api/login", handler.Login)                     //手机号登录接口
	r.GET("api/wechat", handler.Wechat)                    //连接第三方微信
	r.GET("api/wechat/Login", handler.WechatLogin)         //微信登录接口
	r.GET("api/callback", handler.CallBack)                //第三方微信回调接口
	r.POST("api/Upload", Middleware(), pkg.Upload)         //修改用户头像
	r.POST("api/infoUser", Middleware(), handler.InfoUser) //用户详情接口
}
