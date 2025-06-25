package main

import (
	passengers "Ride_hailing_platform/kitex_gen/passengers/passengersserver"
	"log"
)

func main() {
	svr := passengers.NewServer(new(PassengersServerImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
