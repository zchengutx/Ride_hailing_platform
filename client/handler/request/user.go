package request

type CreateUserReq struct {
	Mobile      string `json:"mobile" form:"mobile" binding:"required"`
	SendSmsCode string `json:"send_sms_code" form:"send_sms_code" binding:"required"`
}
type LoginUserReq struct {
	Mobile      string `json:"mobile" form:"mobile" binding:"required"`
	SendSmsCode string `json:"send_sms_code" form:"send_sms_code" binding:"required"`
}
