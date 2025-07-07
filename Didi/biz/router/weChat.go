package router

import (
	"Didi/biz/handler/api"
	"github.com/cloudwego/hertz/pkg/route"
)

func WechatModel(r *route.RouterGroup) {
	weChatModel := r.Group("/wechat")
	{
		weChatModel.GET("/sign", api.Sign)
		weChatModel.GET("/one", api.One)
		weChatModel.GET("/calBlack", api.CalBlack)
	}
}
