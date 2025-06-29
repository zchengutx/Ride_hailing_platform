package request

type SendSmsReq struct {
	Mobile string `json:"mobile,required"`
	Source string `json:"source,required"`
}

type RegisterReq struct {
	Mobile   string `json:"mobile,required"`
	SendCode string `json:"sendCode,required"`
}

type LoginReq struct {
	Mobile   string `json:"mobile,required"`
	SendCode string `json:"sendCode,required"`
}

type BindMobileReq struct {
	Mobile string `json:"mobile,required"`
}

type UpdateCancelOrderReq struct {
	OrderId       int32  `json:"orderId"`
	ConfirmPerson int32  `json:"confirmPerson"`
	ConfirmReason string `json:"confirmReason"`
	Id            int32  `json:"id"`
	ConfirmRemark string `json:"confirmRemark"`
}
