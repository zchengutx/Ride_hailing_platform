package request

type WeChatServerValidationReq struct {
	Signature string `json:"signature" form:"signature" vd:"required"`
	Timestamp string `json:"timestamp" form:"timestamp" vd:"required"`
	Nonce     string `json:"nonce" form:"nonce" vd:"required"`
	Echostr   string `json:"echostr" form:"echostr" vd:"required"`
}

type SignReq struct {
	Signature string `json:"signature" form:"signature" vd:"required"`
	Timestamp string `json:"timestamp" form:"timestamp" vd:"required"`
	Nonce     string `json:"nonce" form:"nonce" vd:"required"`
	Echostr   string `json:"echostr" form:"echostr" vd:"required"`
}

type CalblackReq struct {
	Code string `form:"code"`
}
