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
