package main

import (
	user "Didi/kitex_gen/Didi/user/userserver"
	"Didi/rpc/common/global"
	_ "Didi/rpc/common/inits"
	"fmt"
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
)

func main() {
	UserConf := &global.AppConf.UserClient
	userAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", UserConf.Host, UserConf.Port))
	svr := user.NewServer(new(UserServerImpl), server.WithServiceAddr(userAddr))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
