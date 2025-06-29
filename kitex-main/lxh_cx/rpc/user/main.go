package main

import (
	"github.com/cloudwego/kitex/server"
	"log"
	_ "lxh_cx/init"
	user "lxh_cx/kitex_gen/lxh_cx/user/userservice"
	"net"
)

func main() {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:5001")
	svr := user.NewServer(new(UserServiceImpl), server.WithServiceAddr(addr))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
