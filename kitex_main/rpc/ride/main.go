package main

import (
	"github.com/cloudwego/kitex/server"
	_ "kitex_main/init"
	ride "kitex_main/kitex_gen/Ride_hailing_platform/ride/trips"
	"log"
	"net"
)

func main() {

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:50053")

	svr := ride.NewServer(new(TripsImpl), server.WithServiceAddr(addr))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
