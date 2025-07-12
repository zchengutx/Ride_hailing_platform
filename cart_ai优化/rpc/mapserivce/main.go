package main

import (
	mapservice "cart/kitex_gen/cart/mapservice/mapservice"
	"cart/rpc/basic/global"
	_ "cart/rpc/basic/inits"
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
)

func main() {
	mapConf := &global.AppConf.MapService
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", mapConf.Host, mapConf.Port))
	svr := mapservice.NewServer(new(MapServiceImpl), server.WithServiceAddr(addr))
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
