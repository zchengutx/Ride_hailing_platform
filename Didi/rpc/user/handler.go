package main

import (
	pb "Didi/kitex_gen/Didi/user"
	"Didi/rpc/common/global"
	"Didi/rpc/common/model"
	"Didi/utils"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// UserServerImpl implements the last service interface defined in the IDL.
type UserServerImpl struct{}

// redisHashManager Redis哈希存储管理器实例
var redisHashManager = utils.NewRedisHashManager()

// SendSms implements the UserServerImpl interface.
func (s *UserServerImpl) SendSms(ctx context.Context, req *pb.SendSmsReq) (*pb.SendSmsResp, error) {
	ctx = context.Background()

	// 生成6位数验证码
	code := rand.Intn(900000) + 100000

	// 使用哈希存储保存验证码，包含更多信息
	if err := redisHashManager.SetSmsCode(req.Mobile, code, time.Minute*5); err != nil {
		return &pb.SendSmsResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}

	return &pb.SendSmsResp{
		Code:    200,
		Message: "短信发送成功",
	}, nil
}

// LoginUser implements the UserServerImpl interface.
func (s *UserServerImpl) LoginUser(ctx context.Context, req *pb.LoginUserReq) (*pb.LoginUserResp, error) {
	ctx = context.Background()

	// 从哈希存储中获取验证码
	storedCode, err := redisHashManager.GetSmsCode(req.Mobile)
	if err != nil {
		return &pb.LoginUserResp{
			Code:    301,
			Message: "验证码已过期或不存在",
		}, nil
	}

	// 增加尝试次数（用于防刷）
	attempts, _ := redisHashManager.IncrementSmsAttempts(req.Mobile)

	// 检查尝试次数是否超限（防止暴力破解）
	if attempts > 5 {
		// 删除验证码，防止继续尝试
		redisHashManager.DeleteSmsCode(req.Mobile)
		return &pb.LoginUserResp{
			Code:    301,
			Message: "验证码尝试次数过多，请重新获取",
		}, nil
	}

	// 验证验证码
	if storedCode != req.SendSmsCode {
		return &pb.LoginUserResp{
			Code:    301,
			Message: "验证码错误，剩余尝试次数：" + strconv.FormatInt(5-attempts, 10),
		}, nil
	}

	// 验证成功后删除验证码（避免重复使用）
	redisHashManager.DeleteSmsCode(req.Mobile)

	// 查询或创建用户
	var user model.User
	if err := global.DB.Debug().Where("mobile = ?", req.Mobile).Limit(1).Find(&user).Error; err != nil {
		return &pb.LoginUserResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}

	nikeName := "Didi" + req.Mobile
	if user.Id == 0 {
		user.Mobile = req.Mobile
		user.NikeName = nikeName
		global.DB.Create(&user)
	}

	return &pb.LoginUserResp{
		Code:    200,
		Message: "登录成功",
		UId:     int64(user.Id),
	}, nil
}

// RealName implements the UserServerImpl interface.
func (s *UserServerImpl) RealName(ctx context.Context, req *pb.RealNameReq) (*pb.RealNameResp, error) {
	// TODO: Your code here...
	ctx = context.Background()
	var user model.User
	if err := global.DB.Debug().Where("id = ?", req.UId).Limit(1).Find(&user).Error; err != nil {
		return &pb.RealNameResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}
	if user.Id == 0 {
		return &pb.RealNameResp{
			Code:    404,
			Message: "用户未登录",
		}, nil
	}
	m := model.User{
		Id:       uint64(req.UId),
		UserName: req.UserName,
		Sex:      req.Sex,
		Age:      uint64(req.Age),
		IdCard:   req.IdCard,
	}
	if err := global.DB.Debug().Create(&m).Error; err != nil {
		return &pb.RealNameResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}
	return &pb.RealNameResp{
		Code:    200,
		Message: "实名信息提交成功",
	}, nil
}

// CallACar implements the UserServerImpl interface.
func (s *UserServerImpl) TakeACar(ctx context.Context, req *pb.TakeACarReq) (*pb.TakeACarResp, error) {
	ctx = context.Background()

	// 参数验证
	if req.UId <= 0 {
		return &pb.TakeACarResp{
			Code:    400,
			Message: "用户ID无效",
		}, nil
	}
	if req.StartLocation == "" || req.EndLocation == "" {
		return &pb.TakeACarResp{
			Code:    400,
			Message: "起始地址和目的地址不能为空",
		}, nil
	}

	// 验证用户是否存在并已实名认证
	var user model.User
	if err := global.DB.Debug().Where("id = ?", req.UId).Limit(1).Find(&user).Error; err != nil {
		fmt.Printf("查询用户信息失败: %v\n", err)
		return &pb.TakeACarResp{
			Code:    503,
			Message: "服务器异常",
		}, nil
	}
	if user.Id == 0 {
		return &pb.TakeACarResp{
			Code:    404,
			Message: "用户未登录",
		}, nil
	}

	// 检查用户是否已实名认证（可选验证）
	if user.UserName == "" || user.IdCard == "" {
		return &pb.TakeACarResp{
			Code:    403,
			Message: "请先完成实名认证",
		}, nil
	}

	// 生成唯一订单ID
	orderId := utils.GenerateOrderId()
	fmt.Printf("开始处理用户 %d 的打车订单，订单ID: %s\n", req.UId, orderId)

	// 解析起始位置坐标
	// 支持格式: "地址名称,纬度,经度" 如: "北京市朝阳区,39.918058,116.397026"
	var latitude, longitude float64
	var addressName string

	locationParts := strings.Split(req.StartLocation, ",")
	if len(locationParts) >= 3 {
		addressName = strings.TrimSpace(locationParts[0])
		if lat, err := strconv.ParseFloat(strings.TrimSpace(locationParts[1]), 64); err == nil {
			latitude = lat
		}
		if lng, err := strconv.ParseFloat(strings.TrimSpace(locationParts[2]), 64); err == nil {
			longitude = lng
		}
	}

	// 坐标有效性验证
	if latitude == 0 || longitude == 0 || latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180 {
		// 坐标无效，返回错误
		fmt.Printf("错误: 起始位置坐标无效，纬度=%f, 经度=%f\n", latitude, longitude)
		return &pb.TakeACarResp{
			Code:    400,
			Message: "起始位置坐标无效，请确保定位功能已开启并重新获取位置",
		}, nil
	}

	fmt.Printf("成功解析起始位置 %s: 纬度=%f, 经度=%f\n", addressName, latitude, longitude)

	// 准备订单数据
	orderData := map[string]interface{}{
		"order_id":       orderId,
		"user_id":        fmt.Sprintf("%d", req.UId),
		"user_name":      user.UserName,
		"user_mobile":    user.Mobile,
		"start_location": addressName,
		"end_location":   req.EndLocation,
		"cart_type":      req.CartType,
		"latitude":       latitude,
		"longitude":      longitude,
		"status":         "waiting", // 等待接单
		"created_at":     time.Now().Unix(),
	}

	// 将订单信息存储到Redis哈希存储（24小时过期）
	if err := redisHashManager.SetOrderInfo(orderId, orderData, time.Hour*24); err != nil {
		fmt.Printf("订单存储失败: %v\n", err)
		return &pb.TakeACarResp{
			Code:    503,
			Message: "订单创建失败，服务器异常",
		}, nil
	}
	fmt.Printf("订单 %s 已成功存储到Redis\n", orderId)

	// 查找周边3公里的在线司机
	geoManager := utils.NewRedisGeoManager()
	nearbyDrivers, err := geoManager.GetNearbyDrivers(latitude, longitude, 3.0, req.CartType)
	if err != nil {
		fmt.Printf("查找周边司机失败: %v\n", err)
		// 订单已创建，但查找司机失败
		return &pb.TakeACarResp{
			Code:    201,
			Message: fmt.Sprintf("订单创建成功（订单号：%s），但暂时无法查找周边司机，请稍后重试", orderId),
		}, nil
	}

	fmt.Printf("在3公里范围内找到 %d 位符合条件的在线司机\n", len(nearbyDrivers))

	// 如果找到了周边在线司机，通过RabbitMQ推送订单信息
	if len(nearbyDrivers) > 0 {
		// 提取司机ID列表并记录司机信息
		var driverIDs []string
		for i, driver := range nearbyDrivers {
			driverIDs = append(driverIDs, driver.DriverID)
			fmt.Printf("司机 %d: ID=%s, 距离=%.0f米, 坐标=(%.6f,%.6f)\n",
				i+1, driver.DriverID, driver.Distance, driver.Latitude, driver.Longitude)
		}

		// 构建订单消息
		orderMessage := utils.OrderMessage{
			OrderID:       orderId,
			UserID:        fmt.Sprintf("%d", req.UId),
			StartLocation: addressName,
			EndLocation:   req.EndLocation,
			CartType:      req.CartType,
			CreatedAt:     time.Now().Unix(),
			Latitude:      latitude,
			Longitude:     longitude,
			DriverIDs:     driverIDs,
		}

		// 获取RabbitMQ管理器并推送消息
		rabbitMQManager := utils.GetRabbitMQManager()
		if rabbitMQManager != nil {
			if err := rabbitMQManager.PublishOrderToNearbyDrivers(orderMessage); err != nil {
				fmt.Printf("推送订单消息失败: %v\n", err)
				// 推送失败但订单已创建，返回部分成功状态
				return &pb.TakeACarResp{
					Code:    202,
					Message: fmt.Sprintf("订单创建成功（订单号：%s），找到%d位司机，但消息推送失败，系统将重试", orderId, len(nearbyDrivers)),
				}, nil
			} else {
				fmt.Printf("成功推送订单 %s 给 %d 个周边在线司机: %v\n", orderId, len(driverIDs), driverIDs)
			}
		} else {
			fmt.Println("警告: RabbitMQ未初始化，无法推送消息给司机")
			return &pb.TakeACarResp{
				Code:    202,
				Message: fmt.Sprintf("订单创建成功（订单号：%s），找到%d位司机，但消息系统暂时不可用", orderId, len(nearbyDrivers)),
			}, nil
		}

		// 更新订单状态为已推送
		if err := redisHashManager.UpdateOrderStatus(orderId, "pushed"); err != nil {
			fmt.Printf("更新订单状态失败: %v\n", err)
		}

		return &pb.TakeACarResp{
			Code:    200,
			Message: fmt.Sprintf("订单创建成功！订单号：%s，已通知%d位司机，请等待司机接单", orderId, len(nearbyDrivers)),
		}, nil
	} else {
		fmt.Printf("订单 %s 在3公里范围内暂无符合条件的在线司机（车型：%s）\n", orderId, req.CartType)

		// 更新订单状态为无司机
		if err := redisHashManager.UpdateOrderStatus(orderId, "no_driver"); err != nil {
			fmt.Printf("更新订单状态失败: %v\n", err)
		}

		return &pb.TakeACarResp{
			Code:    200,
			Message: fmt.Sprintf("订单创建成功（订单号：%s），但附近3公里内暂无可用的%s司机，系统将持续为您寻找", orderId, req.CartType),
		}, nil
	}
}
