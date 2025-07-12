package main

import (
	pb "cart/kitex_gen/cart/driver"
	"cart/rpc/basic/global"
	"cart/rpc/basic/model"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"gorm.io/gorm"
)

// DriverServiceImpl implements the last service interface defined in the IDL.
type DriverServiceImpl struct{}

// validateDriverInfo 验证司机基本信息
func (s *DriverServiceImpl) validateDriverInfo(req *pb.DriverPetitionReq) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("姓名不能为空")
	}
	if strings.TrimSpace(req.NickName) == "" {
		return fmt.Errorf("昵称不能为空")
	}
	if strings.TrimSpace(req.Mobile) == "" {
		return fmt.Errorf("手机号不能为空")
	}
	if len(req.Mobile) != 11 {
		return fmt.Errorf("手机号格式不正确")
	}
	if req.CarAge < 0 || req.CarAge > 20 {
		return fmt.Errorf("车龄应在0-20年之间")
	}
	if strings.TrimSpace(req.IdCardFileId) == "" {
		return fmt.Errorf("身份证照片不能为空")
	}
	if strings.TrimSpace(req.DriverLicenseFileId) == "" {
		return fmt.Errorf("驾驶证照片不能为空")
	}
	if strings.TrimSpace(req.DrivingLicenseFileId) == "" {
		return fmt.Errorf("行驶证照片不能为空")
	}
	if strings.TrimSpace(req.AvatarFileId) == "" {
		return fmt.Errorf("头像不能为空")
	}
	return nil
}

// generateOrderCode 生成订单编号
func (s *DriverServiceImpl) generateOrderCode() string {
	timestamp := time.Now().Unix()
	randomNum := rand.Intn(10000)
	return fmt.Sprintf("D%d%04d", timestamp, randomNum)
}

// convertOrderToInfo 将订单模型转换为返回信息
func (s *DriverServiceImpl) convertOrderToInfo(order model.LxhOrder) *pb.OrderInfo {
	// 获取乘客信息
	var passenger model.LxhPassenger
	passengerName := "未知乘客"
	passengerMobile := ""
	if err := global.DB.Where("id = ?", order.PassengerId).First(&passenger).Error; err == nil {
		if passenger.NickName != "" {
			passengerName = passenger.NickName
		} else if passenger.Name != "" {
			passengerName = passenger.Name
		}
		// 脱敏处理手机号
		if len(passenger.Mobile) == 11 {
			passengerMobile = passenger.Mobile[:3] + "****" + passenger.Mobile[7:]
		}
	}

	return &pb.OrderInfo{
		Id:              order.Id,
		OrderCode:       order.OrderCode,
		Amount:          order.Amount,
		OrderStatus:     order.OrderStatus,
		PassengerId:     order.PassengerId,
		PassengerName:   passengerName,
		PassengerMobile: passengerMobile,
		StartAddr:       order.StartAddr,
		EndAddr:         order.EndEnd,
		StartTime:       order.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:         order.EndTime.Format("2006-01-02 15:04:05"),
		OrderType:       order.OrderType,
	}
}

// DriverPetition 司机注册申请接口，实现司机注册申请功能
func (s *DriverServiceImpl) DriverPetition(ctx context.Context, req *pb.DriverPetitionReq) (*pb.DriverPetitionResp, error) {
	ctx = context.Background()

	// 参数验证
	if err := s.validateDriverInfo(req); err != nil {
		return &pb.DriverPetitionResp{
			Code:    400,
			Message: err.Error(),
		}, nil
	}

	// 检查手机号是否已注册
	var existingDriver model.LxhDriver
	if err := global.DB.Debug().Where("mobile = ?", req.Mobile).First(&existingDriver).Error; err == nil {
		return &pb.DriverPetitionResp{
			Code:    400,
			Message: "该手机号已注册司机账户",
		}, nil
	}

	// 检查是否已有审核中的申请
	var existingCheck model.LxhDriverCheck
	if err := global.DB.Debug().Joins("JOIN lxh_driver ON lxh_driver_check.id = lxh_driver.id").
		Where("lxh_driver.mobile = ? AND lxh_driver_check.check_status = ?", req.Mobile, "pending").
		First(&existingCheck).Error; err == nil {
		return &pb.DriverPetitionResp{
			Code:    400,
			Message: "您已有审核中的申请，请耐心等待",
		}, nil
	}

	// 使用事务确保数据一致性
	tx := global.DB.Begin()
	if tx.Error != nil {
		return &pb.DriverPetitionResp{
			Code:    503,
			Message: "服务器异常，请重试",
		}, nil
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建司机记录
	driver := model.LxhDriver{
		Name:      strings.TrimSpace(req.Name),
		NickName:  strings.TrimSpace(req.NickName),
		Mobile:    req.Mobile,
		CarAge:    req.CarAge,
		Status:    "offline",
		FileId:    0,
		AllAcount: 0,
	}

	if err := tx.Create(&driver).Error; err != nil {
		tx.Rollback()
		return &pb.DriverPetitionResp{
			Code:    503,
			Message: "申请提交失败，请重试",
		}, nil
	}

	// 创建审核记录
	driverCheck := model.LxhDriverCheck{
		Id:                   driver.Id,
		IdCardFileId:         req.IdCardFileId,
		DriverLicenseFileId:  req.DriverLicenseFileId,
		DrivingLicenseFileId: req.DrivingLicenseFileId,
		AvatorFileId:         req.AvatarFileId,
		CheckStatus:          "pending",
		Remark:               "申请已提交，正在审核中",
	}

	if err := tx.Create(&driverCheck).Error; err != nil {
		tx.Rollback()
		return &pb.DriverPetitionResp{
			Code:    503,
			Message: "申请提交失败，请重试",
		}, nil
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return &pb.DriverPetitionResp{
			Code:    503,
			Message: "申请提交失败，请重试",
		}, nil
	}

	return &pb.DriverPetitionResp{
		Code:          200,
		Message:       "申请提交成功，我们将在1-3个工作日内完成审核",
		ApplicationId: &driver.Id,
	}, nil
}

// CheckStatus 查询司机审核状态接口，查询司机审核状态
func (s *DriverServiceImpl) CheckStatus(ctx context.Context, req *pb.CheckStatusReq) (*pb.CheckStatusResp, error) {
	ctx = context.Background()

	if req.DriverId <= 0 {
		return &pb.CheckStatusResp{
			Code:    400,
			Message: "司机ID无效",
		}, nil
	}

	var driverCheck model.LxhDriverCheck
	if err := global.DB.Debug().Where("id = ?", req.DriverId).First(&driverCheck).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.CheckStatusResp{
				Code:    404,
				Message: "未找到审核记录",
			}, nil
		}
		return &pb.CheckStatusResp{
			Code:    503,
			Message: "查询失败，请重试",
		}, nil
	}

	// 根据审核状态返回友好信息
	var statusMessage string
	switch driverCheck.CheckStatus {
	case "pending":
		statusMessage = "审核中，请耐心等待"
	case "approved":
		statusMessage = "审核通过，可以开始接单"
	case "rejected":
		statusMessage = "审核未通过"
	default:
		statusMessage = "未知状态"
	}

	return &pb.CheckStatusResp{
		Code:        200,
		Message:     statusMessage,
		CheckStatus: &driverCheck.CheckStatus,
		Remark:      &driverCheck.Remark,
	}, nil
}

// DriverLogin 司机登录接口，实现司机登录功能
func (s *DriverServiceImpl) DriverLogin(ctx context.Context, req *pb.DriverLoginReq) (*pb.DriverLoginResp, error) {
	ctx = context.Background()

	// 参数验证
	if strings.TrimSpace(req.Mobile) == "" {
		return &pb.DriverLoginResp{
			Code:    400,
			Message: "手机号不能为空",
		}, nil
	}
	if strings.TrimSpace(req.SmsCode) == "" {
		return &pb.DriverLoginResp{
			Code:    400,
			Message: "验证码不能为空",
		}, nil
	}

	// 验证短信验证码
	result, err := global.Rdb.Get(ctx, "sendSms"+req.Mobile).Result()
	if err != nil {
		return &pb.DriverLoginResp{
			Code:    601,
			Message: "验证码已过期，请重新获取",
		}, nil
	}

	if result != req.SmsCode {
		return &pb.DriverLoginResp{
			Code:    400,
			Message: "验证码错误",
		}, nil
	}

	// 查找司机记录
	var driver model.LxhDriver
	if err := global.DB.Debug().Where("mobile = ?", req.Mobile).First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.DriverLoginResp{
				Code:    401,
				Message: "该手机号未注册司机账户",
			}, nil
		}
		return &pb.DriverLoginResp{
			Code:    503,
			Message: "登录失败，请重试",
		}, nil
	}

	// 检查审核状态
	var driverCheck model.LxhDriverCheck
	if err := global.DB.Debug().Where("id = ?", driver.Id).First(&driverCheck).Error; err != nil {
		return &pb.DriverLoginResp{
			Code:    403,
			Message: "账户信息异常，请联系客服",
		}, nil
	}

	if driverCheck.CheckStatus == "pending" {
		return &pb.DriverLoginResp{
			Code:    403,
			Message: "账户正在审核中，请耐心等待",
		}, nil
	}

	if driverCheck.CheckStatus == "rejected" {
		return &pb.DriverLoginResp{
			Code:    403,
			Message: fmt.Sprintf("账户审核未通过：%s", driverCheck.Remark),
		}, nil
	}

	if driverCheck.CheckStatus != "approved" {
		return &pb.DriverLoginResp{
			Code:    403,
			Message: "账户状态异常，请联系客服",
		}, nil
	}

	// 生成token（这里简化处理，实际应该用JWT）
	token := fmt.Sprintf("driver_token_%d_%d", driver.Id, time.Now().Unix())

	// 删除验证码（防止重复使用）
	global.Rdb.Del(ctx, "sendSms"+req.Mobile)

	// 构建司机信息
	driverInfo := &pb.DriverInfo{
		Id:         driver.Id,
		Name:       driver.Name,
		NickName:   driver.NickName,
		Mobile:     driver.Mobile,
		AllAccount: driver.AllAcount,
		CarAge:     driver.CarAge,
		Status:     driver.Status,
		FileId:     driver.FileId,
	}

	return &pb.DriverLoginResp{
		Code:       200,
		Message:    "登录成功",
		DriverId:   &driver.Id,
		Token:      &token,
		DriverInfo: driverInfo,
	}, nil
}

// GetDriverInfo 获取司机信息接口，获取司机详细信息
func (s *DriverServiceImpl) GetDriverInfo(ctx context.Context, req *pb.GetDriverInfoReq) (*pb.GetDriverInfoResp, error) {
	ctx = context.Background()

	if req.DriverId <= 0 {
		return &pb.GetDriverInfoResp{
			Code:    400,
			Message: "司机ID无效",
		}, nil
	}

	var driver model.LxhDriver
	if err := global.DB.Debug().Where("id = ?", req.DriverId).First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.GetDriverInfoResp{
				Code:    404,
				Message: "司机不存在",
			}, nil
		}
		return &pb.GetDriverInfoResp{
			Code:    503,
			Message: "查询失败，请重试",
		}, nil
	}

	driverInfo := &pb.DriverInfo{
		Id:         driver.Id,
		Name:       driver.Name,
		NickName:   driver.NickName,
		Mobile:     driver.Mobile,
		AllAccount: driver.AllAcount,
		CarAge:     driver.CarAge,
		Status:     driver.Status,
		FileId:     driver.FileId,
	}

	return &pb.GetDriverInfoResp{
		Code:       200,
		Message:    "获取成功",
		DriverInfo: driverInfo,
	}, nil
}

// UpdateDriverInfo 更新司机信息接口，更新司机个人信息
func (s *DriverServiceImpl) UpdateDriverInfo(ctx context.Context, req *pb.UpdateDriverInfoReq) (*pb.UpdateDriverInfoResp, error) {
	ctx = context.Background()

	if req.DriverId <= 0 {
		return &pb.UpdateDriverInfoResp{
			Code:    400,
			Message: "司机ID无效",
		}, nil
	}

	// 验证司机是否存在
	var driver model.LxhDriver
	if err := global.DB.Debug().Where("id = ?", req.DriverId).First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.UpdateDriverInfoResp{
				Code:    404,
				Message: "司机不存在",
			}, nil
		}
		return &pb.UpdateDriverInfoResp{
			Code:    503,
			Message: "更新失败，请重试",
		}, nil
	}

	// 构建更新数据
	updateData := make(map[string]interface{})
	if req.NickName != nil {
		nickName := strings.TrimSpace(*req.NickName)
		if nickName == "" {
			return &pb.UpdateDriverInfoResp{
				Code:    400,
				Message: "昵称不能为空",
			}, nil
		}
		updateData["nick_name"] = nickName
	}
	if req.FileId != nil {
		updateData["file_id"] = *req.FileId
	}

	if len(updateData) == 0 {
		return &pb.UpdateDriverInfoResp{
			Code:    400,
			Message: "没有需要更新的数据",
		}, nil
	}

	if err := global.DB.Debug().Model(&model.LxhDriver{}).Where("id = ?", req.DriverId).Updates(updateData).Error; err != nil {
		return &pb.UpdateDriverInfoResp{
			Code:    503,
			Message: "更新失败，请重试",
		}, nil
	}

	return &pb.UpdateDriverInfoResp{
		Code:    200,
		Message: "更新成功",
	}, nil
}

// ChangeStatus 司机状态切换接口，实现司机上线下线功能
func (s *DriverServiceImpl) ChangeStatus(ctx context.Context, req *pb.ChangeStatusReq) (*pb.ChangeStatusResp, error) {
	ctx = context.Background()

	if req.DriverId <= 0 {
		return &pb.ChangeStatusResp{
			Code:    400,
			Message: "司机ID无效",
		}, nil
	}

	// 验证状态值
	if req.Status != "online" && req.Status != "offline" {
		return &pb.ChangeStatusResp{
			Code:    400,
			Message: "状态值无效，只能是online或offline",
		}, nil
	}

	// 验证司机是否存在且已审核通过
	var driver model.LxhDriver
	if err := global.DB.Debug().Where("id = ?", req.DriverId).First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.ChangeStatusResp{
				Code:    404,
				Message: "司机不存在",
			}, nil
		}
		return &pb.ChangeStatusResp{
			Code:    503,
			Message: "状态更新失败，请重试",
		}, nil
	}

	// 检查审核状态
	var driverCheck model.LxhDriverCheck
	if err := global.DB.Debug().Where("id = ?", driver.Id).First(&driverCheck).Error; err != nil || driverCheck.CheckStatus != "approved" {
		return &pb.ChangeStatusResp{
			Code:    403,
			Message: "账户未审核通过，无法上线",
		}, nil
	}

	// 如果是上线，需要位置信息
	if req.Status == "online" && (req.Longitude == nil || req.Latitude == nil) {
		return &pb.ChangeStatusResp{
			Code:    400,
			Message: "上线时必须提供位置信息",
		}, nil
	}

	// 检查是否有未完成的订单
	if req.Status == "offline" {
		var activeOrderCount int64
		global.DB.Debug().Model(&model.LxhOrder{}).Where("driver = ? AND order_status IN ?", req.DriverId, []string{"已接单", "进行中"}).Count(&activeOrderCount)
		if activeOrderCount > 0 {
			return &pb.ChangeStatusResp{
				Code:    400,
				Message: "您还有未完成的订单，无法下线",
			}, nil
		}
	}

	// 更新司机状态
	if err := global.DB.Debug().Model(&model.LxhDriver{}).Where("id = ?", req.DriverId).Update("status", req.Status).Error; err != nil {
		return &pb.ChangeStatusResp{
			Code:    503,
			Message: "状态更新失败，请重试",
		}, nil
	}

	// 如果上线，记录位置信息
	if req.Status == "online" {
		locationKey := fmt.Sprintf("driver_location_%d", req.DriverId)
		locationData := fmt.Sprintf("%s,%s", *req.Longitude, *req.Latitude)
		if err := global.Rdb.Set(ctx, locationKey, locationData, 30*time.Minute).Err(); err != nil {
			fmt.Printf("位置信息存储失败: %v\n", err)
		}
	} else {
		// 下线时清除位置信息
		locationKey := fmt.Sprintf("driver_location_%d", req.DriverId)
		global.Rdb.Del(ctx, locationKey)
	}

	message := "下线成功"
	if req.Status == "online" {
		message = "上线成功，开始接收订单"
	}

	return &pb.ChangeStatusResp{
		Code:    200,
		Message: message,
	}, nil
}

// GetPendingOrders 获取待接订单接口，获取司机可接取的订单列表
func (s *DriverServiceImpl) GetPendingOrders(ctx context.Context, req *pb.GetPendingOrdersReq) (*pb.GetPendingOrdersResp, error) {
	ctx = context.Background()

	if req.DriverId <= 0 {
		return &pb.GetPendingOrdersResp{
			Code:    400,
			Message: "司机ID无效",
		}, nil
	}

	// 验证司机是否在线
	var driver model.LxhDriver
	if err := global.DB.Debug().Where("id = ? AND status = ?", req.DriverId, "online").First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.GetPendingOrdersResp{
				Code:    403,
				Message: "请先上线再查看订单",
			}, nil
		}
		return &pb.GetPendingOrdersResp{
			Code:    503,
			Message: "查询失败，请重试",
		}, nil
	}

	// 查询待接订单（未分配司机的订单）
	var orders []model.LxhOrder
	query := global.DB.Debug().Where("order_status = ? AND (driver = 0 OR driver IS NULL)", "待接单").
		Order("start_time ASC").
		Limit(20) // 默认最多返回20条

	if err := query.Find(&orders).Error; err != nil {
		return &pb.GetPendingOrdersResp{
			Code:    503,
			Message: "查询订单失败，请重试",
		}, nil
	}

	// 转换为返回格式
	var orderInfos []*pb.OrderInfo
	for _, order := range orders {
		orderInfo := s.convertOrderToInfo(order)
		orderInfos = append(orderInfos, orderInfo)
	}

	return &pb.GetPendingOrdersResp{
		Code:    200,
		Message: fmt.Sprintf("共找到%d个待接订单", len(orderInfos)),
		Orders:  orderInfos,
	}, nil
}

// AcceptOrder 司机接单接口，司机接受订单
func (s *DriverServiceImpl) AcceptOrder(ctx context.Context, req *pb.AcceptOrderReq) (*pb.AcceptOrderResp, error) {
	ctx = context.Background()

	if req.DriverId <= 0 {
		return &pb.AcceptOrderResp{
			Code:    400,
			Message: "司机ID无效",
		}, nil
	}
	if req.OrderId <= 0 {
		return &pb.AcceptOrderResp{
			Code:    400,
			Message: "订单ID无效",
		}, nil
	}

	// 验证司机是否在线
	var driver model.LxhDriver
	if err := global.DB.Debug().Where("id = ? AND status = ?", req.DriverId, "online").First(&driver).Error; err != nil {
		return &pb.AcceptOrderResp{
			Code:    403,
			Message: "请先上线再接单",
		}, nil
	}

	// 检查司机是否已有进行中的订单
	var activeOrderCount int64
	global.DB.Debug().Model(&model.LxhOrder{}).Where("driver = ? AND order_status IN ?", req.DriverId, []string{"已接单", "进行中"}).Count(&activeOrderCount)
	if activeOrderCount > 0 {
		return &pb.AcceptOrderResp{
			Code:    400,
			Message: "您还有未完成的订单，请先完成当前订单",
		}, nil
	}

	// 使用事务确保数据一致性
	tx := global.DB.Begin()
	if tx.Error != nil {
		return &pb.AcceptOrderResp{
			Code:    503,
			Message: "接单失败，请重试",
		}, nil
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查订单状态并锁定
	var order model.LxhOrder
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ? AND order_status = ? AND (driver = 0 OR driver IS NULL)", req.OrderId, "待接单").First(&order).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return &pb.AcceptOrderResp{
				Code:    404,
				Message: "订单不存在或已被其他司机接取",
			}, nil
		}
		return &pb.AcceptOrderResp{
			Code:    503,
			Message: "接单失败，请重试",
		}, nil
	}

	// 更新订单状态和司机ID
	if err := tx.Model(&order).Updates(map[string]interface{}{
		"driver":       req.DriverId,
		"order_status": "已接单",
	}).Error; err != nil {
		tx.Rollback()
		return &pb.AcceptOrderResp{
			Code:    503,
			Message: "接单失败，请重试",
		}, nil
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return &pb.AcceptOrderResp{
			Code:    503,
			Message: "接单失败，请重试",
		}, nil
	}

	// 重新查询订单信息
	if err := global.DB.Debug().Where("id = ?", req.OrderId).First(&order).Error; err != nil {
		return &pb.AcceptOrderResp{
			Code:    503,
			Message: "接单成功，但获取订单信息失败",
		}, nil
	}

	// 构建订单信息返回
	orderInfo := s.convertOrderToInfo(order)
	orderInfo.OrderStatus = "已接单"

	return &pb.AcceptOrderResp{
		Code:      200,
		Message:   "接单成功，请尽快前往乘客位置",
		OrderInfo: orderInfo,
	}, nil
}

// StartTrip 开始行程接口，司机开始行程
func (s *DriverServiceImpl) StartTrip(ctx context.Context, req *pb.StartTripReq) (*pb.StartTripResp, error) {
	ctx = context.Background()

	if req.DriverId <= 0 {
		return &pb.StartTripResp{
			Code:    400,
			Message: "司机ID无效",
		}, nil
	}
	if req.OrderId <= 0 {
		return &pb.StartTripResp{
			Code:    400,
			Message: "订单ID无效",
		}, nil
	}

	// 检查订单状态
	var order model.LxhOrder
	if err := global.DB.Debug().Where("id = ? AND driver = ? AND order_status = ?", req.OrderId, req.DriverId, "已接单").First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.StartTripResp{
				Code:    404,
				Message: "订单不存在或状态异常",
			}, nil
		}
		return &pb.StartTripResp{
			Code:    503,
			Message: "开始行程失败，请重试",
		}, nil
	}

	// 更新订单状态为进行中
	currentTime := time.Now()
	if err := global.DB.Debug().Model(&order).Updates(map[string]interface{}{
		"order_status": "进行中",
		"start_time":   currentTime,
	}).Error; err != nil {
		return &pb.StartTripResp{
			Code:    503,
			Message: "开始行程失败，请重试",
		}, nil
	}

	// 更新路线记录状态（如果存在）
	global.DB.Debug().Model(&model.LxhRouteRecord{}).Where("order_id = ?", req.OrderId).Updates(map[string]interface{}{
		"route_status": "running",
		"start_time":   currentTime,
	})

	return &pb.StartTripResp{
		Code:    200,
		Message: "行程已开始，祝您一路平安",
	}, nil
}

// CompleteOrder 完成订单接口，司机完成订单
func (s *DriverServiceImpl) CompleteOrder(ctx context.Context, req *pb.CompleteOrderReq) (*pb.CompleteOrderResp, error) {
	ctx = context.Background()

	if req.DriverId <= 0 {
		return &pb.CompleteOrderResp{
			Code:    400,
			Message: "司机ID无效",
		}, nil
	}
	if req.OrderId <= 0 {
		return &pb.CompleteOrderResp{
			Code:    400,
			Message: "订单ID无效",
		}, nil
	}
	if req.ActualAmount < 0 {
		return &pb.CompleteOrderResp{
			Code:    400,
			Message: "实际金额不能为负数",
		}, nil
	}

	// 使用事务确保数据一致性
	tx := global.DB.Begin()
	if tx.Error != nil {
		return &pb.CompleteOrderResp{
			Code:    503,
			Message: "完成订单失败，请重试",
		}, nil
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查订单状态
	var order model.LxhOrder
	if err := tx.Where("id = ? AND driver = ? AND order_status = ?", req.OrderId, req.DriverId, "进行中").First(&order).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return &pb.CompleteOrderResp{
				Code:    404,
				Message: "订单不存在或状态异常",
			}, nil
		}
		return &pb.CompleteOrderResp{
			Code:    503,
			Message: "完成订单失败，请重试",
		}, nil
	}

	// 计算司机收入（假设司机分成80%）
	driverIncome := req.ActualAmount * 0.8
	currentTime := time.Now()

	// 更新订单状态
	if err := tx.Model(&order).Updates(map[string]interface{}{
		"order_status": "已完成",
		"end_time":     currentTime,
		"amount":       req.ActualAmount,
		"pay_status":   "已支付",
	}).Error; err != nil {
		tx.Rollback()
		return &pb.CompleteOrderResp{
			Code:    503,
			Message: "完成订单失败，请重试",
		}, nil
	}

	// 更新司机总收益
	if err := tx.Model(&model.LxhDriver{}).Where("id = ?", req.DriverId).
		Update("all_acount", gorm.Expr("all_acount + ?", driverIncome)).Error; err != nil {
		tx.Rollback()
		return &pb.CompleteOrderResp{
			Code:    503,
			Message: "收益更新失败，请联系客服",
		}, nil
	}

	// 更新路线记录状态
	tx.Model(&model.LxhRouteRecord{}).Where("order_id = ?", req.OrderId).Updates(map[string]interface{}{
		"route_status": "completed",
		"end_time":     currentTime,
	})

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return &pb.CompleteOrderResp{
			Code:    503,
			Message: "完成订单失败，请重试",
		}, nil
	}

	return &pb.CompleteOrderResp{
		Code:         200,
		Message:      fmt.Sprintf("订单完成，您的收益为%.2f元", driverIncome),
		DriverIncome: &driverIncome,
	}, nil
}

// CancelOrder implements the DriverServiceImpl interface.
func (s *DriverServiceImpl) CancelOrder(ctx context.Context, req *pb.CancelOrderReq) (*pb.CancelOrderResp, error) {
	if req.DriverId <= 0 {
		return &pb.CancelOrderResp{
			Code:    400,
			Message: "司机ID无效",
		}, nil
	}
	if req.OrderId <= 0 {
		return &pb.CancelOrderResp{
			Code:    400,
			Message: "订单ID无效",
		}, nil
	}
	if strings.TrimSpace(req.CancelReason) == "" {
		return &pb.CancelOrderResp{
			Code:    400,
			Message: "取消原因不能为空",
		}, nil
	}

	// 检查订单状态
	var order model.LxhOrder
	if err := global.DB.Where("id = ? AND driver = ?", req.OrderId, req.DriverId).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.CancelOrderResp{
				Code:    404,
				Message: "订单不存在",
			}, nil
		}
		return &pb.CancelOrderResp{
			Code:    503,
			Message: "取消订单失败，请重试",
		}, nil
	}

	// 只能取消已接单或进行中的订单
	if order.OrderStatus != "已接单" && order.OrderStatus != "进行中" {
		return &pb.CancelOrderResp{
			Code:    400,
			Message: "当前订单状态不允许取消",
		}, nil
	}

	// 构建取消备注
	remark := strings.TrimSpace(req.CancelReason)
	if req.CancelRemark != nil {
		additionalRemark := strings.TrimSpace(*req.CancelRemark)
		if additionalRemark != "" {
			remark += " - " + additionalRemark
		}
	}

	// 更新订单状态
	if err := global.DB.Model(&order).Updates(map[string]interface{}{
		"order_status":   "已取消",
		"confirm_by":     "司机",
		"confirm_person": req.DriverId,
		"confirm_reason": req.CancelReason,
		"confirm_remark": remark,
		"driver":         0, // 释放司机绑定
		"end_time":       time.Now(),
	}).Error; err != nil {
		return &pb.CancelOrderResp{
			Code:    503,
			Message: "取消订单失败，请重试",
		}, nil
	}

	// 更新路线记录状态
	global.DB.Model(&model.LxhRouteRecord{}).Where("order_id = ?", req.OrderId).Updates(map[string]interface{}{
		"route_status": "cancelled",
		"end_time":     time.Now(),
	})

	return &pb.CancelOrderResp{
		Code:    200,
		Message: "订单已取消",
	}, nil
}

// GetIncome implements the DriverServiceImpl interface.
func (s *DriverServiceImpl) GetIncome(ctx context.Context, req *pb.GetIncomeReq) (*pb.GetIncomeResp, error) {
	if req.DriverId <= 0 {
		return &pb.GetIncomeResp{
			Code:    400,
			Message: "司机ID无效",
		}, nil
	}

	// 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return &pb.GetIncomeResp{
			Code:    400,
			Message: "开始日期格式错误，请使用YYYY-MM-DD格式",
		}, nil
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return &pb.GetIncomeResp{
			Code:    400,
			Message: "结束日期格式错误，请使用YYYY-MM-DD格式",
		}, nil
	}

	// 验证日期范围
	if endDate.Before(startDate) {
		return &pb.GetIncomeResp{
			Code:    400,
			Message: "结束日期不能早于开始日期",
		}, nil
	}

	// 限制查询范围不超过3个月
	if endDate.Sub(startDate) > 90*24*time.Hour {
		return &pb.GetIncomeResp{
			Code:    400,
			Message: "查询范围不能超过3个月",
		}, nil
	}

	// 验证司机是否存在
	var driver model.LxhDriver
	if err := global.DB.Where("id = ?", req.DriverId).First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.GetIncomeResp{
				Code:    404,
				Message: "司机不存在",
			}, nil
		}
		return &pb.GetIncomeResp{
			Code:    503,
			Message: "查询失败，请重试",
		}, nil
	}

	// 查询司机在指定时间段内的已完成订单
	var orders []model.LxhOrder
	if err := global.DB.Where("driver = ? AND order_status = ? AND DATE(end_time) BETWEEN ? AND ?",
		req.DriverId, "已完成", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Order("end_time DESC").Find(&orders).Error; err != nil {
		return &pb.GetIncomeResp{
			Code:    503,
			Message: "查询收益失败，请重试",
		}, nil
	}

	// 按日期分组统计
	incomeMap := make(map[string]*pb.IncomeDetail)
	totalIncome := 0.0
	totalOrders := 0

	for _, order := range orders {
		dateStr := order.EndTime.Format("2006-01-02")
		driverIncome := order.Amount * 0.8 // 司机分成80%
		platformFee := order.Amount * 0.2  // 平台费20%

		if detail, exists := incomeMap[dateStr]; exists {
			detail.OrderCount++
			detail.TotalIncome += order.Amount
			detail.PlatformFee += platformFee
			detail.NetIncome += driverIncome
		} else {
			incomeMap[dateStr] = &pb.IncomeDetail{
				Date:        dateStr,
				OrderCount:  1,
				TotalIncome: order.Amount,
				PlatformFee: platformFee,
				NetIncome:   driverIncome,
			}
		}
		totalIncome += driverIncome
		totalOrders++
	}

	// 转换为数组并按日期排序
	var incomeDetails []*pb.IncomeDetail
	for _, detail := range incomeMap {
		incomeDetails = append(incomeDetails, detail)
	}

	// 按日期排序（最新的在前）
	for i := 0; i < len(incomeDetails)-1; i++ {
		for j := i + 1; j < len(incomeDetails); j++ {
			if incomeDetails[i].Date < incomeDetails[j].Date {
				incomeDetails[i], incomeDetails[j] = incomeDetails[j], incomeDetails[i]
			}
		}
	}

	message := fmt.Sprintf("查询成功，共完成%d单，总收益%.2f元", totalOrders, totalIncome)

	return &pb.GetIncomeResp{
		Code:          200,
		Message:       message,
		IncomeDetails: incomeDetails,
		TotalIncome:   &totalIncome,
	}, nil
}

// UpdateLocation implements the DriverServiceImpl interface.
func (s *DriverServiceImpl) UpdateLocation(ctx context.Context, req *pb.UpdateLocationReq) (*pb.UpdateLocationResp, error) {
	if req.DriverId <= 0 {
		return &pb.UpdateLocationResp{
			Code:    400,
			Message: "司机ID无效",
		}, nil
	}

	if strings.TrimSpace(req.Longitude) == "" || strings.TrimSpace(req.Latitude) == "" {
		return &pb.UpdateLocationResp{
			Code:    400,
			Message: "经纬度不能为空",
		}, nil
	}

	// 验证司机是否存在且在线
	var driver model.LxhDriver
	if err := global.DB.Where("id = ? AND status = ?", req.DriverId, "online").First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &pb.UpdateLocationResp{
				Code:    403,
				Message: "请先上线再更新位置",
			}, nil
		}
		return &pb.UpdateLocationResp{
			Code:    503,
			Message: "位置更新失败，请重试",
		}, nil
	}

	// 构建位置数据
	locationData := fmt.Sprintf("%s,%s", req.Longitude, req.Latitude)
	if req.Speed != nil {
		locationData += fmt.Sprintf(",speed:%.2f", *req.Speed)
	}
	if req.Direction != nil {
		locationData += fmt.Sprintf(",direction:%.2f", *req.Direction)
	}
	locationData += fmt.Sprintf(",timestamp:%d", time.Now().Unix())

	// 存储到Redis中，30分钟过期
	locationKey := fmt.Sprintf("driver_location_%d", req.DriverId)
	if err := global.Rdb.Set(ctx, locationKey, locationData, 30*time.Minute).Err(); err != nil {
		return &pb.UpdateLocationResp{
			Code:    503,
			Message: "位置更新失败，请重试",
		}, nil
	}

	// 记录日志（生产环境可以考虑异步处理）
	fmt.Printf("[%s] 司机 %d 位置更新：%s\n", time.Now().Format("2006-01-02 15:04:05"), req.DriverId, locationData)

	return &pb.UpdateLocationResp{
		Code:    200,
		Message: "位置更新成功",
	}, nil
}
