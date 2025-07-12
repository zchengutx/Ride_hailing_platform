package main

import (
	"github.com/cloudwego/kitex/server"
	_ "kitex_main/init"
	home "kitex_main/kitex_gen/Ride_hailing_platform/home/home"
	"log"
	"net"
)

func main() {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:50052")

	svr := home.NewServer(new(HomeImpl), server.WithServiceAddr(addr))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
