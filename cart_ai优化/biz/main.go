package main

import (
	"cart/biz/dal/global"
	_ "cart/biz/dal/inits"
	"cart/biz/middleware" // 新增
	"cart/biz/router"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	HertzConf := &global.AppConf.System
	h := server.Default(server.WithHostPorts(fmt.Sprintf("%s:%d", HertzConf.Host, HertzConf.Port)))

	// 注册CORS中间件
	h.Use(middleware.CORS())

	router.Router(h)

	h.Spin()
}
