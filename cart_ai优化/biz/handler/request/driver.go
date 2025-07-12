package request

// DriverPetitionReq 司机注册申请请求
type DriverPetitionReq struct {
	Name                 string `json:"name" form:"name" vd:"required"`
	Mobile               string `json:"mobile" form:"mobile" vd:"required"`
	NickName             string `json:"nick_name" form:"nick_name" vd:"required"`
	CarAge               int64  `json:"car_age" form:"car_age" vd:"required"`
	IdCardFileId         string `json:"id_card_file_id" form:"id_card_file_id" vd:"required"`
	DriverLicenseFileId  string `json:"driver_license_file_id" form:"driver_license_file_id" vd:"required"`
	DrivingLicenseFileId string `json:"driving_license_file_id" form:"driving_license_file_id" vd:"required"`
	AvatarFileId         string `json:"avatar_file_id" form:"avatar_file_id" vd:"required"`
}

// CheckStatusReq 查询审核状态请求
type CheckStatusReq struct {
	DriverId int64 `json:"driver_id" form:"driver_id" vd:"required"`
}

// DriverLoginReq 司机登录请求
type DriverLoginReq struct {
	Mobile  string `json:"mobile" form:"mobile" vd:"required"`
	SmsCode string `json:"sms_code" form:"sms_code" vd:"required"`
}

// GetDriverInfoReq 获取司机信息请求
type GetDriverInfoReq struct {
	DriverId int64 `json:"driver_id" form:"driver_id" vd:"required"`
}

// UpdateDriverInfoReq 更新司机信息请求
type UpdateDriverInfoReq struct {
	DriverId int64  `json:"driver_id" form:"driver_id" vd:"required"`
	Name     string `json:"name" form:"name" vd:"required"`
	NickName string `json:"nick_name" form:"nick_name" vd:"required"`
	Age      int64  `json:"age" form:"age"`
	Sex      string `json:"sex" form:"sex" vd:"required"`
	FileId   int64  `json:"file_id" form:"file_id"` // 非指针
}

// ChangeStatusReq 上线/下线请求
type ChangeStatusReq struct {
	DriverId  int64  `json:"driver_id" form:"driver_id" vd:"required"`
	Status    string `json:"status" form:"status" vd:"required"`
	Longitude string `json:"longitude" form:"longitude"`
	Latitude  string `json:"latitude" form:"latitude"`
}

// UpdateLocationReq 更新司机位置请求
type UpdateLocationReq struct {
	DriverId  int64   `json:"driver_id" form:"driver_id" vd:"required"`
	Longitude string  `json:"longitude" form:"longitude" vd:"required"`
	Latitude  string  `json:"latitude" form:"latitude" vd:"required"`
	Speed     float64 `json:"speed" form:"speed"`         // 非指针
	Direction float64 `json:"direction" form:"direction"` // 新增
}

// GetPendingOrdersReq 获取待接订单请求
type GetPendingOrdersReq struct {
	DriverId  int64  `json:"driver_id" form:"driver_id" vd:"required"`
	Longitude string `json:"longitude" form:"longitude" vd:"required"`
	Latitude  string `json:"latitude" form:"latitude" vd:"required"`
	Radius    int32  `json:"radius" form:"radius"` // 非指针
	Page      int    `json:"page" form:"page"`
	PageSize  int    `json:"page_size" form:"page_size"`
}

// AcceptOrderReq 接单请求
type AcceptOrderReq struct {
	DriverId int64 `json:"driver_id" form:"driver_id" vd:"required"`
	OrderId  int64 `json:"order_id" form:"order_id" vd:"required"`
}

// StartTripReq 开始行程请求
type StartTripReq struct {
	DriverId  int64  `json:"driver_id" form:"driver_id" vd:"required"`
	OrderId   int64  `json:"order_id" form:"order_id" vd:"required"`
	Longitude string `json:"longitude" form:"longitude" vd:"required"`
	Latitude  string `json:"latitude" form:"latitude" vd:"required"`
}

// CompleteOrderReq 完成订单请求
type CompleteOrderReq struct {
	DriverId     int64   `json:"driver_id" form:"driver_id" vd:"required"`
	OrderId      int64   `json:"order_id" form:"order_id" vd:"required"`
	Longitude    string  `json:"longitude" form:"longitude" vd:"required"`
	Latitude     string  `json:"latitude" form:"latitude" vd:"required"`
	ActualAmount float64 `json:"actual_amount" form:"actual_amount" vd:"required"`
}

// CancelOrderReq 司机取消订单请求
type CancelOrderReq struct {
	DriverId     int64  `json:"driver_id" form:"driver_id" vd:"required"`
	OrderId      int64  `json:"order_id" form:"order_id" vd:"required"`
	CancelReason string `json:"cancel_reason" form:"cancel_reason" vd:"required"`
	CancelRemark string `json:"cancel_remark" form:"cancel_remark"` // 非指针
}

// GetIncomeReq 查询收益请求
type GetIncomeReq struct {
	DriverId  int64  `json:"driver_id" form:"driver_id" vd:"required"`
	StartDate string `json:"start_date" form:"start_date" vd:"required"`
	EndDate   string `json:"end_date" form:"end_date" vd:"required"`
}

// GetNearbyDriversReq 查询附近司机请求
type GetNearbyDriversReq struct {
	DriverId  int64   `json:"driver_id" form:"driver_id" vd:"required"`
	Longitude string  `json:"longitude" form:"longitude" vd:"required"`
	Latitude  string  `json:"latitude" form:"latitude" vd:"required"`
	Radius    float64 `json:"radius" form:"radius"`
}
