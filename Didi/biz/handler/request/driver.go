package request

type CallACarReq struct {
	Address              string `json:"address" form:"address" binding:"required"`                               // 家庭地址
	DrivingLicenseNumber string `json:"driving_license_number" form:"driving_license_number" binding:"required"` // 驾驶证号
	QuasiDrivingType     string `json:"quasi_driving_type" form:"quasi_driving_type" binding:"required"`         // 准驾车型
	DrivingAge           uint64 `json:"driving_age" form:"driving_age" binding:"required"`                       // 驾龄
	CarNum               string `json:"car_num" form:"car_num" binding:"required"`                               // 车牌号
	CarType              string `json:"car_type" form:"car_type" binding:"required"`                             // 车辆品牌类型
	VehicleMileage       string `json:"vehicle_mileage" form:"vehicle_mileage" binding:"required"`               // 车辆行驶里程
	ServingTheCity       string `json:"serving_the_city" form:"serving_the_city" binding:"required"`
}
type DriverAuditReq struct {
	DriverId    uint64 `json:"driver_id" form:"driver_id" binding:"required"`
	AuditStatus string `json:"audit_status" form:"audit_status" binding:"required"`
}
type AddDriverReq struct {
	DriverId uint64 `json:"driver_id" form:"driver_id" binding:"required"`
}
type DriverOnlineReq struct {
	DriverId     uint64 `json:"driver_id" form:"driver_id" binding:"required"`
	DriverStatus string `json:"driver_status" form:"driver_status" binding:"required"`
}

type ReceivingOrderReq struct {
	OrderId          string  `json:"order_id" form:"order_id" binding:"required"`                   // 订单ID
	CurrentLatitude  float64 `json:"current_latitude" form:"current_latitude" binding:"required"`   // 司机当前纬度
	CurrentLongitude float64 `json:"current_longitude" form:"current_longitude" binding:"required"` // 司机当前经度
}
