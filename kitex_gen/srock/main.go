package main

import (
	"github.com/cloudwego/kitex/server"
	_ "kitexTwo/common/inits"
	__ "kitexTwo/kitex_gen/kitexTwo/user/user"
	"log"
	"net"
)

func main() {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:50051")
	svr := __.NewServer(new(UserImpl), server.WithServiceAddr(addr))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
