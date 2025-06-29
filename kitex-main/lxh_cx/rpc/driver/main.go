package main

import (
	"github.com/cloudwego/kitex/server"
	"log"
	_ "lxh_cx/init"
	driver "lxh_cx/kitex_gen/lxh_cx/driver/driverservice"
	"net"
)

func main() {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:5003")
	svr := driver.NewServer(new(DriverServiceImpl), server.WithServiceAddr(addr))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
