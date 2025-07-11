package request

type SendSmsReq struct {
	Mobile      string `json:"mobile" form:"mobile" vd:"required"`
	SendSmsCode string `json:"send_sms_code" form:"send_sms_code" vd:"required"`
}
type RegisterPassengerReq struct {
	Mobile      string `json:"mobile" form:"mobile" vd:"required"`
	SendSmsCode string `json:"send_sms_code" form:"send_sms_code" vd:"required"`
}
type LoginPassengerReq struct {
	Mobile      string `json:"mobile" form:"mobile" vd:"required"`
	SendSmsCode string `json:"send_sms_code" form:"send_sms_code" vd:"required"`
}

// UpdatePassengerInfoReq 更新乘客信息请求
type UpdatePassengerInfoReq struct {
	Name     *string `json:"name" form:"name"`
	NickName *string `json:"nick_name" form:"nick_name"`
	FileId   *int64  `json:"file_id" form:"file_id"`
	Age      *int64  `json:"age" form:"age"`
	Sex      *string `json:"sex" form:"sex"`
}

// CreateOrderReq 创建订单请求
type CreateOrderReq struct {
	StartAddr   string  `json:"start_addr" form:"start_addr" vd:"required"`
	StartLng    float64 `json:"start_lng" form:"start_lng" vd:"required"`
	StartLat    float64 `json:"start_lat" form:"start_lat" vd:"required"`
	EndAddr     string  `json:"end_addr" form:"end_addr" vd:"required"`
	EndLng      float64 `json:"end_lng" form:"end_lng" vd:"required"`
	EndLat      float64 `json:"end_lat" form:"end_lat" vd:"required"`
	OrderType   string  `json:"order_type" form:"order_type" vd:"required"`
	AppointTime *string `json:"appoint_time" form:"appoint_time"`
}

// GetOrderListReq 获取订单列表请求
type GetOrderListReq struct {
	Page     *int32  `json:"page" form:"page"`
	PageSize *int32  `json:"page_size" form:"page_size"`
	Status   *string `json:"status" form:"status"`
}

// PassengerCancelOrderReq 取消订单请求
type PassengerCancelOrderReq struct {
	OrderId int64   `json:"order_id" form:"order_id" vd:"required"`
	Reason  string  `json:"reason" form:"reason" vd:"required"`
	Remark  *string `json:"remark" form:"remark"`
}

// EvaluateOrderReq 评价订单请求
type EvaluateOrderReq struct {
	OrderId int64    `json:"order_id" form:"order_id" vd:"required"`
	Rating  int32    `json:"rating" form:"rating" vd:"required,min=1,max=5"`
	Comment *string  `json:"comment" form:"comment"`
	Tags    []string `json:"tags" form:"tags"`
}

// GetFavoriteLocationsReq 获取收藏地址请求
type GetFavoriteLocationsReq struct {
	LocationType *string `json:"location_type" form:"location_type"`
}

// AddFavoriteLocationReq 添加收藏地址请求
type AddFavoriteLocationReq struct {
	LocationType string  `json:"location_type" form:"location_type" vd:"required"`
	LocationName string  `json:"location_name" form:"location_name" vd:"required"`
	Address      string  `json:"address" form:"address" vd:"required"`
	Lng          float64 `json:"lng" form:"lng" vd:"required"`
	Lat          float64 `json:"lat" form:"lat" vd:"required"`
	Province     string  `json:"province" form:"province" vd:"required"`
	City         string  `json:"city" form:"city" vd:"required"`
	District     string  `json:"district" form:"district" vd:"required"`
	IsDefault    *bool   `json:"is_default" form:"is_default"`
}

// GetHotLocationsReq 获取热门地点请求
type GetHotLocationsReq struct {
	City     string  `json:"city" form:"city" vd:"required"`
	Category *string `json:"category" form:"category"`
	Limit    *int32  `json:"limit" form:"limit"`
}

// BindWechatReq 绑定微信请求
type BindWechatReq struct {
	Code string `json:"code" form:"code" vd:"required"`
}

// GetRouteRecordsReq 获取路线记录请求
type GetRouteRecordsReq struct {
	Page     *int32 `json:"page" form:"page"`
	PageSize *int32 `json:"page_size" form:"page_size"`
}

// SearchAddressReq 地址搜索请求
type SearchAddressReq struct {
	Keyword string   `json:"keyword" form:"keyword" vd:"required"`
	City    *string  `json:"city" form:"city"`
	Lng     *float64 `json:"lng" form:"lng"`
	Lat     *float64 `json:"lat" form:"lat"`
	Limit   *int32   `json:"limit" form:"limit"`
}

// HomePageReq 主页请求
type HomePageReq struct {
	Location *string `json:"location" form:"location"`
}

// CallACarReq 叫车请求
type CallACarReq struct {
	StartingPlace string `json:"starting_place" form:"starting_place" vd:"required"`
	Destination   string `json:"destination" form:"destination" vd:"required"`
}
