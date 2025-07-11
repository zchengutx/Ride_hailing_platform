// api 包下的 driver.go 主要实现司机相关的接口处理逻辑
package api

import (
	"cart/biz/dal/global"
	"cart/biz/handler/request"
	"cart/biz/middleware"
	"cart/biz/utils"
	pb "cart/kitex_gen/cart/driver"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// DriverClient 是与司机服务交互的RPC客户端
var (
	DriverClient = utils.GetDefaultDriverClient()
)

// DriverPetition 司机注册申请接口，处理司机注册申请
func DriverPetition(ctx context.Context, c *app.RequestContext) {
	var req request.DriverPetitionReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	// 验证手机号格式
	if !utils.ValidateMobile(req.Mobile) {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "手机号码格式不正确",
			"data": nil,
		})
		return
	}

	// 调用RPC服务
	response, _ := DriverClient.DriverPetition(ctx, &pb.DriverPetitionReq{
		Name:                 req.Name,
		Mobile:               req.Mobile,
		NickName:             req.NickName,
		CarAge:               req.CarAge,
		IdCardFileId:         req.IdCardFileId,
		DriverLicenseFileId:  req.DriverLicenseFileId,
		DrivingLicenseFileId: req.DrivingLicenseFileId,
		AvatarFileId:         req.AvatarFileId,
	})

	c.JSON(200, response)
}

// CheckStatus 查询司机审核状态接口，查询司机审核状态
func CheckStatus(ctx context.Context, c *app.RequestContext) {
	var req request.CheckStatusReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	response, _ := DriverClient.CheckStatus(ctx, &pb.CheckStatusReq{
		DriverId: req.DriverId,
	})

	c.JSON(200, response)
}

// DriverLogin 司机登录接口，司机登录功能
func DriverLogin(ctx context.Context, c *app.RequestContext) {
	var req request.DriverLoginReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	// 验证手机号格式
	if !utils.ValidateMobile(req.Mobile) {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "手机号码格式不正确",
			"data": nil,
		})
		return
	}

	driver, _ := DriverClient.DriverLogin(ctx, &pb.DriverLoginReq{
		Mobile:  req.Mobile,
		SmsCode: req.SmsCode,
	})

	if driver.Code != 200 {
		c.JSON(200, driver)
		return
	}

	// 生成JWT token
	token, _ := middleware.NewJWT(global.JWT_SELECT_KEY).CreateToken(middleware.CustomClaims{
		ID: int(*driver.DriverId),
	})

	c.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "登录成功",
		"data": map[string]interface{}{
			"token":       token,
			"driver_info": driver.DriverInfo,
		},
	})
}

// GetDriverInfo 获取司机信息接口，获取司机详细信息
func GetDriverInfo(ctx context.Context, c *app.RequestContext) {
	var req request.GetDriverInfoReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	response, _ := DriverClient.GetDriverInfo(ctx, &pb.GetDriverInfoReq{
		DriverId: req.DriverId,
	})

	c.JSON(200, response)
}

// UpdateDriverInfo 更新司机信息接口，更新司机个人信息
func UpdateDriverInfo(ctx context.Context, c *app.RequestContext) {
	var req request.UpdateDriverInfoReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	// 构建更新请求
	updateReq := &pb.UpdateDriverInfoReq{
		DriverId: req.DriverId,
	}

	// 只传递非空值
	if req.NickName != "" {
		updateReq.NickName = &req.NickName
	}
	if req.FileId != 0 {
		updateReq.FileId = &req.FileId
	}

	response, _ := DriverClient.UpdateDriverInfo(ctx, updateReq)
	c.JSON(200, response)
}

// ChangeStatus 司机状态切换接口，司机上线下线功能
func ChangeStatus(ctx context.Context, c *app.RequestContext) {
	var req request.ChangeStatusReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	// 验证状态值
	if req.Status != "online" && req.Status != "offline" {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "状态值无效，只能是online或offline",
			"data": nil,
		})
		return
	}

	// 如果是上线，必须提供位置信息
	if req.Status == "online" && (req.Longitude == "" || req.Latitude == "") {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "上线时必须提供位置信息",
			"data": nil,
		})
		return
	}

	changeReq := &pb.ChangeStatusReq{
		DriverId: req.DriverId,
		Status:   req.Status,
	}

	// 添加位置信息（如果有）
	if req.Longitude != "" {
		changeReq.Longitude = &req.Longitude
	}
	if req.Latitude != "" {
		changeReq.Latitude = &req.Latitude
	}

	response, _ := DriverClient.ChangeStatus(ctx, changeReq)
	c.JSON(200, response)
}

// GetPendingOrders 获取待接订单接口，获取司机可接取的订单列表
func GetPendingOrders(ctx context.Context, c *app.RequestContext) {
	var req request.GetPendingOrdersReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	getReq := &pb.GetPendingOrdersReq{
		DriverId:  req.DriverId,
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
	}

	// 设置搜索半径，默认5公里
	if req.Radius > 0 {
		getReq.Radius = &req.Radius
	} else {
		defaultRadius := int32(5)
		getReq.Radius = &defaultRadius
	}

	response, _ := DriverClient.GetPendingOrders(ctx, getReq)
	c.JSON(200, response)
}

// AcceptOrder 司机接单接口，司机接受订单
func AcceptOrder(ctx context.Context, c *app.RequestContext) {
	var req request.AcceptOrderReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	response, _ := DriverClient.AcceptOrder(ctx, &pb.AcceptOrderReq{
		DriverId: req.DriverId,
		OrderId:  req.OrderId,
	})

	c.JSON(200, response)
}

// StartTrip 开始行程接口，司机开始行程
func StartTrip(ctx context.Context, c *app.RequestContext) {
	var req request.StartTripReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	response, _ := DriverClient.StartTrip(ctx, &pb.StartTripReq{
		DriverId:  req.DriverId,
		OrderId:   req.OrderId,
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
	})

	c.JSON(200, response)
}

// CompleteOrder 完成订单接口，司机完成订单
func CompleteOrder(ctx context.Context, c *app.RequestContext) {
	var req request.CompleteOrderReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	// 验证金额
	if req.ActualAmount <= 0 {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "实际金额必须大于0",
			"data": nil,
		})
		return
	}

	response, _ := DriverClient.CompleteOrder(ctx, &pb.CompleteOrderReq{
		DriverId:     req.DriverId,
		OrderId:      req.OrderId,
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
		ActualAmount: req.ActualAmount,
	})

	c.JSON(200, response)
}

// CancelOrder 取消订单接口，司机取消订单
func CancelOrder(ctx context.Context, c *app.RequestContext) {
	var req request.CancelOrderReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	cancelReq := &pb.CancelOrderReq{
		DriverId:     req.DriverId,
		OrderId:      req.OrderId,
		CancelReason: req.CancelReason,
	}

	// 添加取消备注（如果有）
	if req.CancelRemark != "" {
		cancelReq.CancelRemark = &req.CancelRemark
	}

	response, _ := DriverClient.CancelOrder(ctx, cancelReq)
	c.JSON(200, response)
}

// GetIncome 查询司机收益接口，查询司机收益统计
func GetIncome(ctx context.Context, c *app.RequestContext) {
	var req request.GetIncomeReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	// 简单的日期格式验证
	if len(req.StartDate) != 10 || len(req.EndDate) != 10 {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "日期格式错误，请使用YYYY-MM-DD格式",
			"data": nil,
		})
		return
	}

	response, _ := DriverClient.GetIncome(ctx, &pb.GetIncomeReq{
		DriverId:  req.DriverId,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	})

	c.JSON(200, response)
}

// UpdateLocation 更新司机位置接口，更新司机实时位置
func UpdateLocation(ctx context.Context, c *app.RequestContext) {
	var req request.UpdateLocationReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	updateReq := &pb.UpdateLocationReq{
		DriverId:  req.DriverId,
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
	}

	// 添加可选参数
	if req.Speed > 0 {
		updateReq.Speed = &req.Speed
	}
	if req.Direction >= 0 {
		updateReq.Direction = &req.Direction
	}

	response, _ := DriverClient.UpdateLocation(ctx, updateReq)
	c.JSON(200, response)
}
