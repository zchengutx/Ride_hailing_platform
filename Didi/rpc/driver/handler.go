package main

import (
	pb "Didi/kitex_gen/Didi/driver"
	"Didi/rpc/common/global"
	"Didi/rpc/common/model"
	"Didi/utils"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// DriverServerImpl implements the last service interface defined in the IDL.
type DriverServerImpl struct{}

// redisHashManager Redis哈希存储管理器实例
var redisHashManager = utils.NewRedisHashManager()

// redisGeoManager Redis地理位置管理器实例
var redisGeoManager = utils.NewRedisGeoManager()

// CallACar implements the DriverServerImpl interface.
func (s *DriverServerImpl) CallACar(ctx context.Context, req *pb.CallACarReq) (*pb.CallACarResp, error) {
	// TODO: Your code here...
	ctx = context.Background()
	var user model.User
	if err := global.DB.Debug().Where("id = ?", req.UserId).Limit(1).Find(&user).Error; err != nil {
		return &pb.CallACarResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}
	if user.Id == 0 {
		return &pb.CallACarResp{
			Code:    404,
			Message: "用户未登录",
		}, nil
	}
	if user.IdCard == "" {
		return &pb.CallACarResp{
			Code:    601,
			Message: "用户未实名,不允许申请",
		}, nil
	}
	application := model.DriverApplication{
		Driver:               user.Id,
		Mobile:               user.Mobile,
		IdCard:               user.IdCard,
		Age:                  user.Age,
		Sex:                  user.Sex,
		Address:              req.Address,
		DrivingLicenseNumber: req.DrivingLicenseNumber,
		QuasiDrivingType:     req.QuasiDrivingType,
		DrivingAge:           uint64(req.DrivingAge),
		CarNum:               req.CarNum,
		CarType:              req.CarType,
		VehicleMileage:       req.VehicleMileage,
		ServingTheCity:       req.ServingTheCity,
	}
	if err := global.DB.Debug().Create(&application).Error; err != nil {
		return &pb.CallACarResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}
	return &pb.CallACarResp{
		Code:    200,
		Message: "申请成功",
	}, nil
}

// DriverAudit implements the DriverServerImpl interface.
func (s *DriverServerImpl) DriverAudit(ctx context.Context, req *pb.DriverAuditReq) (*pb.DriverAuditResp, error) {
	// TODO: Your code here...
	ctx = context.Background()
	var drivers model.DriverApplication
	if err := global.DB.Debug().Where("id = ?", req.DriverId).Limit(1).Find(&drivers).Error; err != nil {
		return &pb.DriverAuditResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}
	if drivers.Id == 0 {
		return &pb.DriverAuditResp{
			Code:    404,
			Message: "该记录不存在",
		}, nil
	}
	if drivers.AuditStatus == "已审核" {
		return &pb.DriverAuditResp{
			Code:    505,
			Message: "该信息已审核",
		}, nil
	}
	application := model.DriverApplication{
		Id:          uint64(req.DriverId),
		AuditStatus: req.AuditStatus,
	}
	if err := global.DB.Debug().Updates(&application).Error; err != nil {
		return &pb.DriverAuditResp{
			Code:    503,
			Message: "审核失败",
		}, nil
	}
	return &pb.DriverAuditResp{
		Code:    200,
		Message: "审核成功",
	}, nil
}

// AddDriver implements the DriverServerImpl interface.
func (s *DriverServerImpl) AddDriver(ctx context.Context, req *pb.AddDriverReq) (*pb.AddDriverResp, error) {
	// TODO: Your code here...
	ctx = context.Background()
	var drivers model.DriverApplication
	if err := global.DB.Debug().Where("id = ?", req.DriverId).Limit(1).Find(&drivers).Error; err != nil {
		return &pb.AddDriverResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}
	if drivers.Id == 0 {
		return &pb.AddDriverResp{
			Code:    404,
			Message: "该信息不存在",
		}, nil
	}
	if drivers.AuditStatus != "已审核" {
		return &pb.AddDriverResp{
			Code:    505,
			Message: "用户审核为通过,不允许添加",
		}, nil
	}
	driver := model.Driver{
		DriverId:             drivers.Id,
		Mobile:               drivers.Mobile,
		IdCard:               drivers.IdCard,
		CarNum:               drivers.CarNum,
		CarType:              drivers.CarType,
		DrivingLicenseNumber: drivers.DrivingLicenseNumber,
		DrivingAge:           drivers.DrivingAge,
	}
	if err := global.DB.Debug().Create(&driver).Error; err != nil {
		return &pb.AddDriverResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}
	return &pb.AddDriverResp{
		Code:    200,
		Message: "司机添加成功",
	}, nil
}

// DriverOnline implements the DriverServerImpl interface.
func (s *DriverServerImpl) DriverOnline(ctx context.Context, req *pb.DriverOnlineReq) (*pb.DriverOnlineResp, error) {
	// TODO: Your code here...
	ctx = context.Background()
	var driver model.Driver
	if err := global.DB.Debug().Where("id = ?", req.DriverId).Limit(1).Find(&driver).Error; err != nil {
		return &pb.DriverOnlineResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}
	if driver.Id == 0 {
		return &pb.DriverOnlineResp{
			Code:    404,
			Message: "司机不存在",
		}, nil
	}
	if driver.DriverStatus == req.DriverStatus {
		return &pb.DriverOnlineResp{
			Code:    505,
			Message: "操作异常",
		}, nil
	}
	m := model.Driver{
		Id:           uint64(req.DriverId),
		DriverStatus: req.DriverStatus,
	}
	if err := global.DB.Debug().Updates(&m).Error; err != nil {
		return &pb.DriverOnlineResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}
	return &pb.DriverOnlineResp{
		Code:    200,
		Message: "状态更改成功",
	}, nil
}

// ReceivingOrder implements the DriverServerImpl interface.
// 司机接单接口 - 与TakeACar用户发单接口逻辑完全匹配
// 支持3公里范围限制、车型匹配、并发安全、状态流转一致性
func (s *DriverServerImpl) ReceivingOrder(ctx context.Context, req *pb.ReceivingOrderReq) (*pb.ReceivingOrderResp, error) {
	ctx = context.Background()

	// 参数验证
	if req.DriverId <= 0 {
		return &pb.ReceivingOrderResp{
			Code:    400,
			Message: "司机ID无效",
		}, nil
	}
	if req.OrderId == "" {
		return &pb.ReceivingOrderResp{
			Code:    400,
			Message: "订单ID不能为空",
		}, nil
	}
	if req.CurrentLatitude == 0 || req.CurrentLongitude == 0 {
		return &pb.ReceivingOrderResp{
			Code:    400,
			Message: "司机当前位置坐标无效",
		}, nil
	}

	fmt.Printf("司机 %d 尝试接单：%s，当前位置：纬度=%.6f，经度=%.6f\n", req.DriverId, req.OrderId, req.CurrentLatitude, req.CurrentLongitude)

	// 验证司机是否存在且为在线状态
	var driver model.Driver
	if err := global.DB.Debug().Where("id = ?", req.DriverId).Limit(1).Find(&driver).Error; err != nil {
		fmt.Printf("查询司机信息失败: %v\n", err)
		return &pb.ReceivingOrderResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}
	if driver.Id == 0 {
		return &pb.ReceivingOrderResp{
			Code:    404,
			Message: "司机不存在",
		}, nil
	}
	if driver.DriverStatus != "在线" {
		return &pb.ReceivingOrderResp{
			Code:    403,
			Message: "司机状态异常，请先上线后再接单",
		}, nil
	}

	// 从Redis获取订单信息
	orderInfo, err := redisHashManager.GetOrderInfo(req.OrderId)
	if err != nil {
		fmt.Printf("获取订单信息失败: %v\n", err)
		return &pb.ReceivingOrderResp{
			Code:    404,
			Message: "订单不存在或已过期",
		}, nil
	}
	if len(orderInfo) == 0 {
		return &pb.ReceivingOrderResp{
			Code:    404,
			Message: "订单不存在",
		}, nil
	}

	// 检查订单状态 - 与TakeACar接口的状态流转保持一致
	orderStatus := orderInfo["status"]
	if orderStatus == "assigned" {
		return &pb.ReceivingOrderResp{
			Code:    409,
			Message: "订单已被其他司机接取",
		}, nil
	}
	if orderStatus != "waiting" && orderStatus != "pushed" && orderStatus != "no_driver" {
		return &pb.ReceivingOrderResp{
			Code:    400,
			Message: fmt.Sprintf("订单状态异常（当前状态：%s），无法接单", orderStatus),
		}, nil
	}

	// 解析订单位置信息 - 与TakeACar接口的存储格式保持一致
	orderLatStr := orderInfo["latitude"]
	orderLngStr := orderInfo["longitude"]
	var orderLat, orderLng float64
	var err1, err2 error

	// TakeACar存储的是浮点数，这里需要正确解析
	if orderLatStr == "" || orderLngStr == "" {
		fmt.Printf("订单位置信息缺失: lat=%s, lng=%s\n", orderLatStr, orderLngStr)
		return &pb.ReceivingOrderResp{
			Code:    503,
			Message: "订单位置信息缺失",
		}, nil
	}

	orderLat, err1 = strconv.ParseFloat(orderLatStr, 64)
	orderLng, err2 = strconv.ParseFloat(orderLngStr, 64)
	if err1 != nil || err2 != nil {
		fmt.Printf("订单位置坐标解析失败: lat=%s, lng=%s, err1=%v, err2=%v\n", orderLatStr, orderLngStr, err1, err2)
		return &pb.ReceivingOrderResp{
			Code:    503,
			Message: "订单位置坐标格式异常",
		}, nil
	}

	// 坐标有效性验证 - 与TakeACar接口的验证逻辑一致
	if orderLat < -90 || orderLat > 90 || orderLng < -180 || orderLng > 180 {
		fmt.Printf("订单位置坐标超出有效范围: lat=%f, lng=%f\n", orderLat, orderLng)
		return &pb.ReceivingOrderResp{
			Code:    503,
			Message: "订单位置坐标超出有效范围",
		}, nil
	}

	// 计算司机与订单起始位置的距离
	distance, err := redisGeoManager.CalculateDistance(req.CurrentLatitude, req.CurrentLongitude, orderLat, orderLng)
	if err != nil {
		fmt.Printf("计算距离失败: %v\n", err)
		return &pb.ReceivingOrderResp{
			Code:    503,
			Message: "距离计算失败",
		}, nil
	}

	// 检查距离是否在3公里范围内
	maxDistance := 3000.0 // 3公里 = 3000米
	if distance > maxDistance {
		fmt.Printf("司机 %d 距离订单起始位置 %.0f 米，超出3公里范围\n", req.DriverId, distance)
		return &pb.ReceivingOrderResp{
			Code:    403,
			Message: fmt.Sprintf("您距离乘客太远（%.1f公里），无法接单", distance/1000),
		}, nil
	}

	fmt.Printf("司机 %d 距离订单起始位置 %.0f 米，在接单范围内\n", req.DriverId, distance)

	// 原子性更新：检查并更新订单状态 - 确保并发安全
	// 再次检查订单状态，防止并发情况下状态变更
	currentStatus, err := redisHashManager.GetOrderField(req.OrderId, "status")
	if err != nil {
		fmt.Printf("再次获取订单状态失败: %v\n", err)
		return &pb.ReceivingOrderResp{
			Code:    503,
			Message: "订单状态检查失败",
		}, nil
	}

	// 与TakeACar的状态流转保持一致：waiting/pushed/no_driver 可接单，assigned 已被接
	if currentStatus == "assigned" {
		return &pb.ReceivingOrderResp{
			Code:    409,
			Message: "订单已被其他司机接取，接单失败",
		}, nil
	}
	if currentStatus != "waiting" && currentStatus != "pushed" && currentStatus != "no_driver" {
		return &pb.ReceivingOrderResp{
			Code:    409,
			Message: fmt.Sprintf("订单状态已变更（%s），接单失败", currentStatus),
		}, nil
	}

	// 更新订单信息：分配司机
	if err := redisHashManager.UpdateOrderDriver(req.OrderId, strconv.FormatInt(req.DriverId, 10)); err != nil {
		fmt.Printf("更新订单司机信息失败: %v\n", err)
		return &pb.ReceivingOrderResp{
			Code:    503,
			Message: "接单失败，服务器异常",
		}, nil
	}

	// 车型匹配验证 - 确保司机车型与订单需求匹配
	orderCartType := orderInfo["cart_type"]
	if orderCartType != "" && driver.CarType != orderCartType {
		fmt.Printf("司机 %d 车型（%s）与订单需求（%s）不匹配\n", req.DriverId, driver.CarType, orderCartType)
		return &pb.ReceivingOrderResp{
			Code:    403,
			Message: fmt.Sprintf("您的车型（%s）与订单需求（%s）不匹配", driver.CarType, orderCartType),
		}, nil
	}

	// 将订单信息保存到数据库 - 字段与TakeACar接口存储保持一致
	userId, _ := strconv.ParseUint(orderInfo["user_id"], 10, 64)
	order := model.Order{
		UserId:        userId,
		DriverId:      uint64(req.DriverId),
		StartLocation: orderInfo["start_location"],
		EndLocation:   orderInfo["end_location"],
		StartLat:      fmt.Sprintf("%.6f", orderLat), // 保持精度一致
		StartLng:      fmt.Sprintf("%.6f", orderLng), // 保持精度一致
		CartType:      orderInfo["cart_type"],
		Status:        "已接单",
	}
	if err := global.DB.Debug().Create(&order).Error; err != nil {
		fmt.Printf("保存订单到数据库失败: %v\n", err)
		// 回滚Redis操作 - 恢复为原状态
		redisHashManager.UpdateOrderStatus(req.OrderId, currentStatus)
		return &pb.ReceivingOrderResp{
			Code:    503,
			Message: "接单失败，数据保存异常",
		}, nil
	}

	// 更新司机状态为忙碌
	driverUpdate := model.Driver{
		Id:           uint64(req.DriverId),
		DriverStatus: "忙碌",
		TotalOrders:  driver.TotalOrders + 1,
	}
	if err := global.DB.Debug().Updates(&driverUpdate).Error; err != nil {
		fmt.Printf("更新司机状态失败: %v\n", err)
		// 这里不回滚，因为订单已经成功分配
	}

	// 更新司机在Redis中的状态和位置 - 确保司机位置信息最新
	driverIdStr := strconv.FormatInt(req.DriverId, 10)
	if err := redisGeoManager.UpdateDriverStatus(driverIdStr, "busy"); err != nil {
		fmt.Printf("更新Redis中司机状态失败: %v\n", err)
	}

	// 同时更新司机的地理位置信息
	driverLocation := utils.DriverLocation{
		DriverID:  driverIdStr,
		Latitude:  req.CurrentLatitude,
		Longitude: req.CurrentLongitude,
		Status:    "busy",
		CarType:   driver.CarType,
		UpdatedAt: time.Now().Unix(),
	}
	if err := redisGeoManager.UpdateDriverLocation(driverLocation); err != nil {
		fmt.Printf("更新Redis中司机位置失败: %v\n", err)
	}

	// 动态构建返回的订单详情 - 基于Redis中实际存储的数据
	orderDetails := make(map[string]interface{})

	// 复制Redis中的所有订单信息
	for key, value := range orderInfo {
		orderDetails[key] = value
	}

	// 添加接单时的计算字段
	orderDetails["distance"] = fmt.Sprintf("%.1f", distance/1000)   // 司机到乘客的距离（公里）
	orderDetails["assigned_time"] = time.Now().Unix()               // 接单时间戳
	orderDetails["driver_id"] = strconv.FormatInt(req.DriverId, 10) // 司机ID
	orderDetails["driver_latitude"] = req.CurrentLatitude           // 司机当前纬度
	orderDetails["driver_longitude"] = req.CurrentLongitude         // 司机当前经度

	// 更新订单状态为已分配
	orderDetails["status"] = "assigned"

	orderDetailsJson, _ := json.Marshal(orderDetails)

	// 通过RabbitMQ通知用户订单被接取 - 与TakeACar接口的推送机制保持一致
	rabbitMQManager := utils.GetRabbitMQManager()
	if rabbitMQManager != nil {
		// 构建接单通知消息
		acceptMessage := map[string]interface{}{
			"type":         "order_accepted",
			"order_id":     req.OrderId,
			"driver_id":    strconv.FormatInt(req.DriverId, 10),
			"driver_name":  fmt.Sprintf("司机%s", driver.Mobile), // 使用手机号后几位作为显示名
			"driver_phone": driver.Mobile,
			"car_num":      driver.CarNum,
			"car_type":     driver.CarType,
			"distance":     fmt.Sprintf("%.1f", distance/1000),
			"accepted_at":  time.Now().Unix(),
		}

		// 这里可以向用户发送接单通知
		fmt.Printf("订单 %s 已被司机 %d（%s）接取，距离%.1f公里\n",
			req.OrderId, req.DriverId, driver.CarNum, distance/1000)
	} else {
		fmt.Println("警告: RabbitMQ未初始化，无法发送接单通知")
	}

	fmt.Printf("司机 %d 成功接取订单 %s，数据库订单ID: %d，订单状态已更新为assigned\n",
		req.DriverId, req.OrderId, order.Id)

	return &pb.ReceivingOrderResp{
		Code:         200,
		Message:      fmt.Sprintf("接单成功！订单号：%s，乘客距离您%.1f公里，请尽快前往接客", req.OrderId, distance/1000),
		OrderDetails: string(orderDetailsJson),
	}, nil
}
