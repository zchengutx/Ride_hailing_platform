package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"kitex_main/gateway/biz/handler"
)

func SetHomeRouter(r *server.Hertz) {
	r.GET("/api/driving", handler.Driving)                                              //获取当前经纬度接口
	r.POST("/api/directionlite", Middleware(), handler.Directionlite)                   //规划路线接口
	r.POST("/api/CreateHistoricalSearch", Middleware(), handler.CreateHistoricalSearch) //保存历史记录接口
	r.POST("/api/HistoricalSearchList", Middleware(), handler.HistoricalSearchList)     //历史记录展示接口
}
