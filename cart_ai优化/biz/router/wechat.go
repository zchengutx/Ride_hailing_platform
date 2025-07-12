package router

import (
	"cart/biz/handler/api"
	"github.com/cloudwego/hertz/pkg/route"
	"log"
)

func WeChatModel(r *route.RouterGroup) {
	log.Println("微信服务")
	weChatModel := r.Group("/wechat")
	{
		weChatModel.GET("/sign", api.Sign)
		weChatModel.GET("/one", api.GetQRCode)
		weChatModel.GET("/calBlack", api.Callback)
	}
}
