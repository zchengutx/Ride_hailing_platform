package request

type CreateTripe struct {
	StartingPoint string `json:"startingPoint,required"`
	Terminal      string `json:"terminal,required"`
	OrderDistance string `json:"orderDistance,required"`
	VehicleType   int    `json:"vehicleType,required"`
}
