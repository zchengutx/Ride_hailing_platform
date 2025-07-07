package request

type SignReq struct {
	Signature string `json:"signature" form:"signature" binding:"required"`
	Timestamp string `json:"timestamp" form:"timestamp" binding:"required"`
	Nonce     string `json:"nonce" form:"nonce" binding:"required"`
	Echostr   string `json:"echostr" form:"echostr" binding:"required"`
}

type CalblackReq struct {
	Code string `form:"code"`
}
