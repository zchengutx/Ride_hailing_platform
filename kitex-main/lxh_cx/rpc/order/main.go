package main

import (
	"github.com/cloudwego/kitex/server"
	"log"
	_ "lxh_cx/init"
	order "lxh_cx/kitex_gen/lxh_cx/order/orderservice"
	"net"
)

func main() {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:5002")
	svr := order.NewServer(new(OrderServiceImpl), server.WithServiceAddr(addr))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
