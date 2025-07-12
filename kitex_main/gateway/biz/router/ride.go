package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"kitex_main/gateway/biz/handler"
)

func SetRideRouter(r *server.Hertz) {
	r.POST("/api/CreateTripe", Middleware(), handler.CreateOrder) //创建订单接口
}
