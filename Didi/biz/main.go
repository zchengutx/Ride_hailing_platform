package main

import (
	_ "Didi/biz/dal/inits"
	"Didi/biz/router"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {

	h := server.Default(server.WithHostPorts("127.0.0.1:8080"))
	router.Register(h)
	h.Spin()
}
