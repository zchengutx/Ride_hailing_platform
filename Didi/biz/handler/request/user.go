package request

type SendSmsReq struct {
	Mobile      string `json:"mobile" form:"mobile" binding:"required"`
	SendSmsCode string `json:"send_sms_code" form:"send_sms_code" binding:"required"`
}
type LoginReq struct {
	Mobile      string `json:"mobile" form:"mobile" binding:"required"`
	SendSmsCode string `json:"send_sms_code" form:"send_sms_code" binding:"required"`
}
type RealNameReq struct {
	UserName string `json:"user_name" form:"user_name" binding:"required"`
	Sex      string `json:"sex" form:"sex" binding:"required"`
	Age      int64  `json:"age" form:"age" binding:"required"`
	IdCard   string `json:"id_card" form:"id_card" binding:"required"`
}
type TakeACarReq struct {
	StartLocation string `json:"start_location" form:"start_location" binding:"required"`
	EndLocation   string `json:"end_location" form:"end_location" binding:"required"`
	CartType      string `json:"cart_type" form:"cart_type" binding:"required"`
}
