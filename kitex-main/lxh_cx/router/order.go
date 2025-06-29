package router

import (
	"github.com/cloudwego/hertz/pkg/route"
	"lxh_cx/biz/handler"
)

func OrderGroup(a *route.RouterGroup) {
	order := a.Group("/order")
	{
		order.GET("pathPlanning", handler.PathPlanning)
	}
}
