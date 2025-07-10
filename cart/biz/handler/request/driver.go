package request

// 司机注册申请请求
type DriverPetitionReq struct {
	Name                 string `json:"name" form:"name" binding:"required"`                                       // 司机姓名
	Mobile               string `json:"mobile" form:"mobile" binding:"required"`                                   // 联系电话
	NickName             string `json:"nick_name" form:"nick_name" binding:"required"`                             // 昵称
	CarAge               int64  `json:"car_age" form:"car_age" binding:"required"`                                 // 车龄
	IdCardFileId         string `json:"id_card_file_id" form:"id_card_file_id" binding:"required"`                 // 身份证文件ID
	DriverLicenseFileId  string `json:"driver_license_file_id" form:"driver_license_file_id" binding:"required"`   // 驾照文件ID
	DrivingLicenseFileId string `json:"driving_license_file_id" form:"driving_license_file_id" binding:"required"` // 行驶证文件ID
	AvatarFileId         string `json:"avatar_file_id" form:"avatar_file_id" binding:"required"`                   // 司机头像文件ID
}

// 查询审核状态请求
type CheckStatusReq struct {
	DriverId int64 `json:"driver_id" form:"driver_id" binding:"required"` // 司机ID
}

// 司机登录请求
type DriverLoginReq struct {
	Mobile  string `json:"mobile" form:"mobile" binding:"required"`     // 手机号
	SmsCode string `json:"sms_code" form:"sms_code" binding:"required"` // 短信验证码
}

// 获取司机信息请求
type GetDriverInfoReq struct {
	DriverId int64 `json:"driver_id" form:"driver_id" binding:"required"` // 司机ID
}

// 更新司机信息请求
type UpdateDriverInfoReq struct {
	DriverId int64  `json:"driver_id" form:"driver_id" binding:"required"` // 司机ID
	NickName string `json:"nick_name" form:"nick_name"`                    // 昵称（可选）
	FileId   int64  `json:"file_id" form:"file_id"`                        // 头像文件ID（可选）
}

// 司机上线/下线请求
type ChangeStatusReq struct {
	DriverId  int64  `json:"driver_id" form:"driver_id" binding:"required"` // 司机ID
	Status    string `json:"status" form:"status" binding:"required"`       // 状态：online-上线, offline-下线
	Longitude string `json:"longitude" form:"longitude"`                    // 经度（上线时必填）
	Latitude  string `json:"latitude" form:"latitude"`                      // 纬度（上线时必填）
}

// 获取待接订单请求
type GetPendingOrdersReq struct {
	DriverId  int64  `json:"driver_id" form:"driver_id" binding:"required"` // 司机ID
	Longitude string `json:"longitude" form:"longitude" binding:"required"` // 司机当前经度
	Latitude  string `json:"latitude" form:"latitude" binding:"required"`   // 司机当前纬度
	Radius    int32  `json:"radius" form:"radius"`                          // 搜索半径（公里），默认5公里
}

// 接单请求
type AcceptOrderReq struct {
	DriverId int64 `json:"driver_id" form:"driver_id" binding:"required"` // 司机ID
	OrderId  int64 `json:"order_id" form:"order_id" binding:"required"`   // 订单ID
}

// 开始行程请求
type StartTripReq struct {
	DriverId  int64  `json:"driver_id" form:"driver_id" binding:"required"` // 司机ID
	OrderId   int64  `json:"order_id" form:"order_id" binding:"required"`   // 订单ID
	Longitude string `json:"longitude" form:"longitude" binding:"required"` // 当前经度
	Latitude  string `json:"latitude" form:"latitude" binding:"required"`   // 当前纬度
}

// 完成订单请求
type CompleteOrderReq struct {
	DriverId     int64   `json:"driver_id" form:"driver_id" binding:"required"`         // 司机ID
	OrderId      int64   `json:"order_id" form:"order_id" binding:"required"`           // 订单ID
	Longitude    string  `json:"longitude" form:"longitude" binding:"required"`         // 当前经度
	Latitude     string  `json:"latitude" form:"latitude" binding:"required"`           // 当前纬度
	ActualAmount float64 `json:"actual_amount" form:"actual_amount" binding:"required"` // 实际费用
}

// 取消订单请求
type CancelOrderReq struct {
	DriverId     int64  `json:"driver_id" form:"driver_id" binding:"required"`         // 司机ID
	OrderId      int64  `json:"order_id" form:"order_id" binding:"required"`           // 订单ID
	CancelReason string `json:"cancel_reason" form:"cancel_reason" binding:"required"` // 取消原因
	CancelRemark string `json:"cancel_remark" form:"cancel_remark"`                    // 取消备注
}

// 查询收益请求
type GetIncomeReq struct {
	DriverId  int64  `json:"driver_id" form:"driver_id" binding:"required"`   // 司机ID
	StartDate string `json:"start_date" form:"start_date" binding:"required"` // 开始日期 YYYY-MM-DD
	EndDate   string `json:"end_date" form:"end_date" binding:"required"`     // 结束日期 YYYY-MM-DD
}

// 更新位置请求
type UpdateLocationReq struct {
	DriverId  int64   `json:"driver_id" form:"driver_id" binding:"required"` // 司机ID
	Longitude string  `json:"longitude" form:"longitude" binding:"required"` // 经度
	Latitude  string  `json:"latitude" form:"latitude" binding:"required"`   // 纬度
	Speed     float64 `json:"speed" form:"speed"`                            // 行驶速度（可选）
	Direction float64 `json:"direction" form:"direction"`                    // 行驶方向（可选）
}
