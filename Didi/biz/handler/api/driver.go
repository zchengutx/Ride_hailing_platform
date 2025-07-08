package api

import (
	"Didi/biz/handler/request"
	"Didi/kitex_gen/Didi/driver"
	"Didi/utils"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

var (
	// 使用封装的客户端管理器获取司机服务客户端，避免重复创建连接
	DriverClient = utils.GetDefaultDriverClient()
)

func CallACar(ctx context.Context, c *app.RequestContext) {
	var req request.CallACarReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	// 动态验证必要参数
	if req.Address == "" || req.DrivingLicenseNumber == "" || req.CarNum == "" || req.CarType == "" {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "必要参数不能为空",
			"data": nil,
		})
		return
	}

	// 调用司机服务进行叫车请求处理
	car, err := DriverClient.CallACar(ctx, &driver.CallACarReq{
		UserId:               int64(c.GetInt("userId")),
		Address:              req.Address,
		DrivingLicenseNumber: req.DrivingLicenseNumber,
		QuasiDrivingType:     req.QuasiDrivingType,
		DrivingAge:           int64(req.DrivingAge),
		CarNum:               req.CarNum,
		CarType:              req.CarType,
		VehicleMileage:       req.VehicleMileage,
		ServingTheCity:       req.ServingTheCity,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	// 根据业务响应动态返回HTTP状态码
	httpCode := 200
	if car.Code != 200 {
		if car.Code >= 500 {
			httpCode = 500
		} else if car.Code >= 400 {
			httpCode = 400
		}
	}
	c.JSON(httpCode, car)
}

func DriverAudit(ctx context.Context, c *app.RequestContext) {
	var req request.DriverAuditReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	// 动态验证审核状态参数
	if req.AuditStatus == "" {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "审核状态不能为空",
			"data": nil,
		})
		return
	}

	audit, err := DriverClient.DriverAudit(ctx, &driver.DriverAuditReq{
		DriverId:    int64(c.GetInt("driverId")),
		AuditStatus: req.AuditStatus,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	// 根据业务响应动态返回HTTP状态码
	httpCode := 200
	if audit.Code != 200 {
		if audit.Code >= 500 {
			httpCode = 500
		} else if audit.Code >= 400 {
			httpCode = 400
		}
	}
	c.JSON(httpCode, audit)
}

func AddDriver(ctx context.Context, c *app.RequestContext) {
	var req request.AddDriverReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	addDriver, err := DriverClient.AddDriver(ctx, &driver.AddDriverReq{
		DriverId: int64(c.GetInt("driverId")),
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	// 根据业务响应动态返回HTTP状态码
	httpCode := 200
	if addDriver.Code != 200 {
		if addDriver.Code >= 500 {
			httpCode = 500
		} else if addDriver.Code >= 400 {
			httpCode = 400
		}
	}
	c.JSON(httpCode, addDriver)
}

func DriverOnline(ctx context.Context, c *app.RequestContext) {
	var req request.DriverOnlineReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	// 动态验证司机状态参数
	if req.DriverStatus == "" {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "司机状态不能为空",
			"data": nil,
		})
		return
	}

	online, err := DriverClient.DriverOnline(ctx, &driver.DriverOnlineReq{
		DriverId:     int64(c.GetInt("driverId")),
		DriverStatus: req.DriverStatus,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	// 根据业务响应动态返回HTTP状态码
	httpCode := 200
	if online.Code != 200 {
		if online.Code >= 500 {
			httpCode = 500
		} else if online.Code >= 400 {
			httpCode = 400
		}
	}
	c.JSON(httpCode, online)
}

func ReceivingOrder(ctx context.Context, c *app.RequestContext) {
	var req request.ReceivingOrderReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	// 动态参数验证
	if req.OrderId == "" {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "订单ID不能为空",
			"data": nil,
		})
		return
	}

	// 动态坐标有效性验证
	if req.CurrentLatitude == 0 || req.CurrentLongitude == 0 {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "司机位置坐标不能为空",
			"data": nil,
		})
		return
	}

	// 坐标范围验证
	if req.CurrentLatitude < -90 || req.CurrentLatitude > 90 || req.CurrentLongitude < -180 || req.CurrentLongitude > 180 {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "司机位置坐标超出有效范围",
			"data": nil,
		})
		return
	}

	// 调用司机服务进行接单处理
	receivingOrder, err := DriverClient.ReceivingOrder(ctx, &driver.ReceivingOrderReq{
		DriverId:         int64(c.GetInt("driverId")),
		OrderId:          req.OrderId,
		CurrentLatitude:  req.CurrentLatitude,
		CurrentLongitude: req.CurrentLongitude,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	// 根据业务响应动态返回HTTP状态码
	httpCode := 200
	if receivingOrder.Code != 200 {
		if receivingOrder.Code >= 500 {
			httpCode = 500
		} else if receivingOrder.Code >= 400 {
			httpCode = 400
		}
	}
	c.JSON(httpCode, receivingOrder)
}
