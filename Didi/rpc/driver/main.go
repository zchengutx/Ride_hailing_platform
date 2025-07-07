package main

import (
	driver "Didi/kitex_gen/Didi/driver/driverserver"
	"Didi/rpc/common/global"
	_ "Didi/rpc/common/inits"
	"fmt"
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
)

func main() {
	driverConf := &global.AppConf.DriverClient
	driverAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", driverConf.Host, driverConf.Port))
	svr := driver.NewServer(new(DriverServerImpl), server.WithServiceAddr(driverAddr))
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
