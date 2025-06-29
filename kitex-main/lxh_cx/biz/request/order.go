package request

type PathPlanningReq struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
}

type PassengerAddOrderReq struct {
	StartAddr string `json:"startAddr"`
	EndEnd    string `json:"endEnd"`
}
