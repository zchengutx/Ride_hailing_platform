package main

import (
	"cart/biz/dal/global"
	_ "cart/biz/dal/inits"
	"cart/biz/middleware"
	"cart/biz/router"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	HertzConf := &global.AppConf.System
	h := server.Default(
		server.WithHostPorts(fmt.Sprintf("%s:%d", HertzConf.Host, HertzConf.Port)),
		server.WithDisablePrintRoute(true), // 关闭所有路由注册日志
	)

	// 注册CORS中间件
	h.Use(middleware.CORS())

	router.Router(h)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	cancel := func() {} // 只保留 cancel 以兼容后续清理逻辑
	go func() {
		<-quit
		cancel()
	}()

	h.Spin()
	// 可以在这里添加清理逻辑
	os.Exit(0)
}
