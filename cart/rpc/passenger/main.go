package main

import (
	passenger "cart/kitex_gen/cart/passenger/passengerservice"
	"cart/rpc/basic/global"
	_ "cart/rpc/basic/inits"
	"fmt"
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
)

func main() {
	passengerConf := &global.AppConf.PassengerService
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", passengerConf.Host, passengerConf.Port))
	svr := passenger.NewServer(new(PassengerServiceImpl), server.WithServiceAddr(addr))
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
