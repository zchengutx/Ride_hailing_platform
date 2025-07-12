package request

type SendSms struct {
	Mobile string `json:"mobile,required"`
	Source string `json:"source,required"`
}

type Login struct {
	Mobile      string `json:"mobile,required"`
	SendSmsCode string `json:"send_sms_code,required"`
}
