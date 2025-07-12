package request

type Driving struct {
	Origins      string `json:"origins,required"`
	Destinations string `json:"destinations,required"`
	CarType      string `json:"car_type,required"`
}

type ComputePrice struct {
	Distance    string `json:"distance,required"`
	VehicleType int    `json:"vehicleType,required"`
}

type CreateHistoricalSearch struct {
	AddrName string `json:"addrName,required"`
	IdAddr   string `json:"idAddr,required"`
}

type HistoricalSearchList struct {
	IdAddr string `json:"idAddr,required"`
}
