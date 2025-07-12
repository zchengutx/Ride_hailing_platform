package main

import (
	driver "cart/kitex_gen/cart/driver/driverservice"
	"cart/rpc/basic/global"
	_ "cart/rpc/basic/inits"
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/kitex/server"
)

func main() {
	driverConf := &global.AppConf.DriverService
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", driverConf.Host, driverConf.Port))
	svr := driver.NewServer(new(DriverServiceImpl), server.WithServiceAddr(addr))

	fmt.Printf("司机服务启动在端口 %d\n", driverConf.Port)
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
