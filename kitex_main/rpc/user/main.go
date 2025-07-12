package main

import (
	"github.com/cloudwego/kitex/server"
	_ "kitex_main/init"
	user "kitex_main/kitex_gen/Ride_hailing_platform/user/userserver"
	"log"
	"net"
)

func main() {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:50051")
	svr := user.NewServer(new(UserServerImpl), server.WithServiceAddr(addr))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
