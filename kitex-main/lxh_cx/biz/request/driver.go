package request

type DriverRegisterReq struct {
	IdCardFileId         int32 `json:"idCardFileId"`
	DriverLicenseFileId  int32 `json:"driverLicenseFileId"`
	DrivingLicenseFileId int32 `json:"drivingLicenseFileId"`
	AvatorFileId         int32 `json:"avatorFileId"`
	DriverId             int32 `json:"driverId"`
}
type DriverAddReq struct {
	IdCardFileId         int32 `json:"idCardFileId"`
	DriverLicenseFileId  int32 `json:"driverLicenseFileId"`
	DrivingLicenseFileId int32 `json:"drivingLicenseFileId"`
	AvatorFileId         int32 `json:"avatorFileId"`
	DriverId             int32 `json:"driverId"`
	PassengerId          int32 `json:"passengerId"`
}
type DriverOverOrderReq struct {
	OrderId int32 `json:"orderId"`
}
type DriverCancelOrderReq struct {
	ConfirmPerson int32  `json:"confirmPerson"`
	ConfirmReason int32  `json:"confirmReason"`
	ConfirmRemark string `json:"confirmRemark"`
	OrderId       int32  `json:"orderId"`
}
