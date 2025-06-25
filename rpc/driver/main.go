package main

import (
	driver "Ride_hailing_platform/kitex_gen/driver/driverserver"
	"log"
)

func main() {
	svr := driver.NewServer(new(DriverServerImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
