// api 包下的 passenger.go 实现乘客相关接口，包括注册、登录、信息管理、订单、收藏、微信绑定等
package api

import (
	"cart/biz/dal/global"
	"cart/biz/handler/request"
	"cart/biz/middleware"
	"cart/biz/utils"
	pb "cart/kitex_gen/cart/passenger"
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

// PassengerClient 是与乘客服务交互的RPC客户端
var (
	PassengerClient = utils.GetDefaultPassengerClient()
)

// SendSms 发送短信验证码接口
func SendSms(ctx context.Context, c *app.RequestContext) {
	var req request.SendSmsReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "查询失败",
			"data": err.Error(),
		})
		return
	}
	if !utils.ValidateMobile(req.Mobile) {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "手机号码格式不正确",
			"data": nil,
		})
		return
	}
	sms, _ := PassengerClient.SendSms(ctx, &pb.SendSmsReq{
		Mobile:      req.Mobile,
		SendSmsCode: req.SendSmsCode,
	})
	c.JSON(200, sms)
}

// RegisterPassenger 乘客注册接口
func RegisterPassenger(ctx context.Context, c *app.RequestContext) {
	var req request.RegisterPassengerReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "查询失败",
			"data": err.Error(),
		})
		return
	}
	if !utils.ValidateMobile(req.Mobile) {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "手机号码格式不正确",
			"data": nil,
		})
		return
	}
	passenger, _ := PassengerClient.RegisterPassenger(ctx, &pb.RegisterPassengerReq{
		Mobile:      req.Mobile,
		SendSmsCode: req.SendSmsCode,
	})
	c.JSON(200, passenger)
}

// LoginPassenger 乘客登录接口，支持手机号+验证码登录，登录成功后生成JWT token
func LoginPassenger(ctx context.Context, c *app.RequestContext) {
	var req request.LoginPassengerReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "查询失败",
			"data": err.Error(),
		})
		return
	}
	if !utils.ValidateMobile(req.Mobile) {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "手机号码格式不正确",
			"data": nil,
		})
		return
	}
	passenger, _ := PassengerClient.LoginPassenger(ctx, &pb.LoginPassengerReq{
		Mobile:      req.Mobile,
		SendSmsCode: req.SendSmsCode,
	})
	if passenger.Code != 200 {
		c.JSON(200, passenger)
	}
	token, _ := middleware.NewJWT(global.JWT_SELECT_KEY).CreateToken(middleware.CustomClaims{
		ID: int(passenger.PassengerId),
	})
	c.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "登录成功",
		"data": map[string]string{
			"token": token,
		},
	})
}

// GetPassengerInfo 获取乘客信息
func GetPassengerInfo(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	resp, _ := PassengerClient.GetPassengerInfo(ctx, &pb.GetPassengerInfoReq{
		PassengerId: int16(currentUser.ID),
	})
	c.JSON(200, resp)
}

// UpdatePassengerInfo 更新乘客信息
func UpdatePassengerInfo(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	var req request.UpdatePassengerInfoReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	updateReq := &pb.UpdatePassengerInfoReq{
		PassengerId: int16(currentUser.ID),
	}

	if req.Name != nil {
		updateReq.Name = req.Name
	}
	if req.NickName != nil {
		updateReq.NickName = req.NickName
	}
	if req.FileId != nil {
		updateReq.FileId = req.FileId
	}
	if req.Age != nil {
		updateReq.Age = req.Age
	}
	if req.Sex != nil {
		updateReq.Sex = req.Sex
	}

	resp, _ := PassengerClient.UpdatePassengerInfo(ctx, updateReq)
	c.JSON(200, resp)
}

// CreateOrder 创建订单
func CreateOrder(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	var req request.CreateOrderReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	resp, _ := PassengerClient.CreateOrder(ctx, &pb.CreateOrderReq{
		PassengerId: int16(currentUser.ID),
		StartAddr:   req.StartAddr,
		StartLng:    req.StartLng,
		StartLat:    req.StartLat,
		EndAddr:     req.EndAddr,
		EndLng:      req.EndLng,
		EndLat:      req.EndLat,
		OrderType:   req.OrderType,
		AppointTime: req.AppointTime,
	})
	c.JSON(200, resp)
}

// GetPassengerOrders 获取订单列表
func GetPassengerOrders(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	var req request.GetOrderListReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	resp, _ := PassengerClient.GetOrderList(ctx, &pb.GetOrderListReq{
		PassengerId: int16(currentUser.ID),
		Page:        req.Page,
		PageSize:    req.PageSize,
		Status:      req.Status,
	})
	c.JSON(200, resp)
}

// GetOrderDetail 获取订单详情
func GetOrderDetail(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	orderIdStr := c.Param("orderId")
	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "订单ID格式错误",
			"data": nil,
		})
		return
	}

	resp, _ := PassengerClient.GetOrderDetail(ctx, &pb.GetOrderDetailReq{
		PassengerId: int16(currentUser.ID),
		OrderId:     orderId,
	})
	c.JSON(200, resp)
}

// PassengerCancelOrder 乘客取消订单
func PassengerCancelOrder(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	var req request.PassengerCancelOrderReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	resp, _ := PassengerClient.CancelOrder(ctx, &pb.CancelOrderReq{
		PassengerId: int16(currentUser.ID),
		OrderId:     req.OrderId,
		Reason:      req.Reason,
		Remark:      req.Remark,
	})
	c.JSON(200, resp)
}

// EvaluateOrder 评价订单
func EvaluateOrder(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	var req request.EvaluateOrderReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	resp, _ := PassengerClient.EvaluateOrder(ctx, &pb.EvaluateOrderReq{
		PassengerId: int16(currentUser.ID),
		OrderId:     req.OrderId,
		Rating:      req.Rating,
		Comment:     req.Comment,
		Tags:        req.Tags,
	})
	c.JSON(200, resp)
}

// GetFavoriteLocations 获取收藏地址
func GetFavoriteLocations(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	var req request.GetFavoriteLocationsReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	resp, _ := PassengerClient.GetFavoriteLocations(ctx, &pb.GetFavoriteLocationsReq{
		PassengerId:  int16(currentUser.ID),
		LocationType: req.LocationType,
	})
	c.JSON(200, resp)
}

// AddFavoriteLocation 添加收藏地址
func AddFavoriteLocation(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	var req request.AddFavoriteLocationReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	resp, _ := PassengerClient.AddFavoriteLocation(ctx, &pb.AddFavoriteLocationReq{
		PassengerId:  int16(currentUser.ID),
		LocationType: req.LocationType,
		LocationName: req.LocationName,
		Address:      req.Address,
		Lng:          req.Lng,
		Lat:          req.Lat,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		IsDefault:    req.IsDefault,
	})
	c.JSON(200, resp)
}

// DeleteFavoriteLocation 删除收藏地址
func DeleteFavoriteLocation(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	locationIdStr := c.Param("locationId")
	locationId, err := strconv.ParseInt(locationIdStr, 10, 64)
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "地址ID格式错误",
			"data": nil,
		})
		return
	}

	resp, _ := PassengerClient.DeleteFavoriteLocation(ctx, &pb.DeleteFavoriteLocationReq{
		PassengerId: int16(currentUser.ID),
		LocationId:  locationId,
	})
	c.JSON(200, resp)
}

// GetHotLocations 获取热门地点
func GetHotLocations(ctx context.Context, c *app.RequestContext) {
	var req request.GetHotLocationsReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	resp, _ := PassengerClient.GetHotLocations(ctx, &pb.GetHotLocationsReq{
		City:     req.City,
		Category: req.Category,
		Limit:    req.Limit,
	})
	c.JSON(200, resp)
}

// BindWechat 绑定微信账号
func BindWechat(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	var req request.BindWechatReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	resp, _ := PassengerClient.BindWechat(ctx, &pb.BindWechatReq{
		PassengerId: int16(currentUser.ID),
		Code:        req.Code,
	})
	c.JSON(200, resp)
}

// UnbindWechat 解绑微信账号
func UnbindWechat(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	resp, _ := PassengerClient.UnbindWechat(ctx, &pb.UnbindWechatReq{
		PassengerId: int16(currentUser.ID),
	})
	c.JSON(200, resp)
}

// GetWechatBindStatus 获取微信绑定状态
func GetWechatBindStatus(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	resp, _ := PassengerClient.GetWechatBindStatus(ctx, &pb.GetWechatBindStatusReq{
		PassengerId: int16(currentUser.ID),
	})
	c.JSON(200, resp)
}

// GetRouteRecords 获取路线记录
func GetRouteRecords(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	var req request.GetRouteRecordsReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	resp, _ := PassengerClient.GetRouteRecords(ctx, &pb.GetRouteRecordsReq{
		PassengerId: int16(currentUser.ID),
		Page:        req.Page,
		PageSize:    req.PageSize,
	})
	c.JSON(200, resp)
}

// SearchAddress 地址搜索建议
func SearchAddress(ctx context.Context, c *app.RequestContext) {
	var req request.SearchAddressReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	resp, _ := PassengerClient.SearchAddress(ctx, &pb.SearchAddressReq{
		Keyword: req.Keyword,
		City:    req.City,
		Lng:     req.Lng,
		Lat:     req.Lat,
		Limit:   req.Limit,
	})
	c.JSON(200, resp)
}

// HomePage 主页服务
func HomePage(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	var req request.HomePageReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	resp, _ := PassengerClient.HomePage(ctx, &pb.HomePageReq{
		PassengerId: int16(currentUser.ID),
		Location:    req.Location,
	})
	c.JSON(200, resp)
}

// CallACar 叫车服务
func CallACar(ctx context.Context, c *app.RequestContext) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*middleware.CustomClaims)

	var req request.CallACarReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	resp, _ := PassengerClient.CallACar(ctx, &pb.CallACarReq{
		PassengerId:   int16(currentUser.ID),
		StartingPlace: req.StartingPlace,
		Destination:   req.Destination,
	})
	c.JSON(200, resp)
}
