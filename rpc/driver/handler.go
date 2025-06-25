package main

import (
	"Ride_hailing_platform/basic/global"
	"Ride_hailing_platform/handler/model"
	driver0 "Ride_hailing_platform/kitex_gen/driver"
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// DriverServerImpl implements the last service interface defined in the IDL.
type DriverServerImpl struct{}

// 短信验证码
func (s *DriverServerImpl) SendSms(ctx context.Context, req *driver0.SendSmsDriverReq) (resp *driver0.SendSmsDriverResp, err error) {
	// TODO: Your code here...
	code := rand.Intn(9000) + 1000
	err = global.Redis.Set(ctx, "sendSms"+req.Mobile+req.Source, code, time.Minute*5).Err()
	if err != nil {
		return &driver0.SendSmsDriverResp{
			Message: "短信验证码发送失败",
			Code:    http.StatusBadRequest,
		}, err
	}
	return &driver0.SendSmsDriverResp{
		Message: "ok",
		Code:    http.StatusOK,
	}, nil
}

// 司机注册功能
func (s *DriverServerImpl) Register(ctx context.Context, req *driver0.DriverRegisterReq) (resp *driver0.DriverRegisterResp, err error) {
	// TODO: Your code here...
	get := global.Redis.Get(ctx, "sendSms"+req.Mobile+"register")
	if get.Val() != req.SmsCode {
		return &driver0.DriverRegisterResp{
			Message: "短信验证码错误，请查证后在输入",
			Code:    http.StatusBadRequest,
		}, err
	}

	var driver model.LxhDriver
	err = global.DB.Where("mobile = ?", req.Mobile).First(&driver).Error
	if err == nil {
		return &driver0.DriverRegisterResp{
			Message: "此司机已存在，请直接登录",
			Code:    http.StatusBadRequest,
		}, nil
	}
	driver = model.LxhDriver{
		Name:     req.Name,
		Mobile:   req.Mobile,
		NickName: req.Name,
	}
	err = global.DB.Create(&driver).Error
	if err != nil {
		return &driver0.DriverRegisterResp{
			Message: "司机注册失败",
			Code:    http.StatusInternalServerError,
		}, err
	}

	global.Redis.Del(ctx, "sendSms"+req.Mobile+"register")

	return &driver0.DriverRegisterResp{
		Message: "司机注册成功",
		Code:    http.StatusOK,
	}, nil
}

// 获取司机信息
func (s *DriverServerImpl) GetDriverInfo(ctx context.Context, req *driver0.GetDriverInfoReq) (resp *driver0.GetDriverInfoResp, err error) {
	var driver model.LxhDriver
	err = global.DB.Where("id = ?", req.DriverID).First(&driver).Error
	if err != nil {
		return &driver0.GetDriverInfoResp{
			Message: "司机不存在",
			Code:    http.StatusBadRequest,
		}, err
	}
	return &driver0.GetDriverInfoResp{
		Message: "获取司机信息成功",
		Code:    http.StatusOK,
	}, nil
}

// 更新司机状态
func (s *DriverServerImpl) UpdateDriverStatus(ctx context.Context, req *driver0.UpdateDriverStatusReq) (resp *driver0.UpdateDriverStatusResp, err error) {
	var status string
	if req.WorkStatus == 0 {
		status = "离线"
	} else if req.WorkStatus == 1 {
		status = "在线"
	} else if req.WorkStatus == 2 {
		status = "载客中"
	} else {
		status = "司机未出车"
	}

	update := global.DB.Model(&model.LxhDriver{}).Where("id = ?", req.DriverID).Update("status", status)
	if update.Error != nil {
		return &driver0.UpdateDriverStatusResp{
			Message: "更新司机状态失败",
			Code:    http.StatusBadRequest,
		}, err
	}

	return &driver0.UpdateDriverStatusResp{
		Message: "更新司机状态成功",
		Code:    http.StatusOK,
	}, nil
}

// 获取待接订单
func (s *DriverServerImpl) GetAwaitOrder(ctx context.Context, req *driver0.GetAwaitOrderReq) (resp *driver0.GetAwaitOrderResp, err error) {
	// 检查司机是否存在并且在线
	var driver model.LxhDriver
	err = global.DB.Where("id = ? AND status = ?", req.DriverID, "在线").First(&driver).Error
	if err != nil {
		return &driver0.GetAwaitOrderResp{
			Message: "查询司机信息失败",
			Code:    http.StatusInternalServerError,
		}, err
	}
	var orders []model.LxhOrder
	err = global.DB.Where("order_status = ? AND (driver IS NULL OR driver = 0)", "待接").Limit(10).Find(&orders).Error
	if err != nil {
		return &driver0.GetAwaitOrderResp{
			Message: "未查询到待接订单",
			Code:    http.StatusBadRequest,
		}, err
	}
	if len(orders) == 0 {
		return &driver0.GetAwaitOrderResp{
			Message: "暂时没有新的订单",
			Code:    http.StatusBadRequest,
		}, err
	}
	return &driver0.GetAwaitOrderResp{
		Message: fmt.Sprintf("现在有%d个待接订单", len(orders)),
		Code:    http.StatusOK,
	}, nil
}

// 接单功能
func (s *DriverServerImpl) AccessOrder(ctx context.Context, req *driver0.AccessOrderReq) (resp *driver0.AccessOrderResp, err error) {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var order model.LxhOrder
	err = tx.Where("id = ? AND order_status = ? AND (driver IS NULL OR driver = 0)", req.DriverID, "待接").First(&order).Error
	if err != nil {
		tx.Rollback()
		return &driver0.AccessOrderResp{
			Message: "查询订单失败",
			Code:    http.StatusBadRequest,
		}, err
	}

	var driver model.LxhDriver
	err = tx.Where("id = ? AND status = ?", req.DriverID, "在线").First(&driver).Error
	if err != nil {
		tx.Rollback()
		return &driver0.AccessOrderResp{
			Message: "司机不在线或是还未出车",
			Code:    http.StatusBadRequest,
		}, err
	}

	err = tx.Model(&driver).Update("status", "载客中").Error
	if err != nil {
		tx.Rollback()
		return &driver0.AccessOrderResp{
			Message: "更新司机状态失败",
			Code:    http.StatusBadRequest,
		}, err
	}

	err = tx.Commit().Error
	if err != nil {
		return &driver0.AccessOrderResp{
			Message: "接单失败",
			Code:    http.StatusBadRequest,
		}, err
	}
	return &driver0.AccessOrderResp{
		Message: "接单成功",
		Code:    http.StatusOK,
	}, err
}

// 拒接功能
func (s *DriverServerImpl) RejectionOrder(ctx context.Context, req *driver0.RejectionOrderReq) (resp *driver0.RejectionOrderResp, err error) {
	var order model.LxhOrder
	err = global.DB.Where("id = ? AND driver = ?", req.OrderID, req.DriverID).First(&order).Error
	if err != nil {
		return &driver0.RejectionOrderResp{
			Message: "订单不存在或司机未出车",
			Code:    http.StatusBadRequest,
		}, err
	}

	if order.OrderStatus == "开始" {
		return &driver0.RejectionOrderResp{
			Message: "订单已开始，无法拒接",
			Code:    http.StatusBadRequest,
		}, err
	}

	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = tx.Model(&model.LxhDriver{}).Where("id = ?", req.DriverID).Update("status", "在线").Error
	if err != nil {
		tx.Rollback()
		return &driver0.RejectionOrderResp{
			Message: "司机状态更新失败",
			Code:    http.StatusBadRequest,
		}, err
	}

	err = tx.Commit().Error
	if err != nil {
		return &driver0.RejectionOrderResp{
			Message: "拒接失败",
			Code:    http.StatusBadRequest,
		}, err
	}
	return &driver0.RejectionOrderResp{
		Message: "司机成功拒接此单",
		Code:    http.StatusOK,
	}, err
}

// ArriveLocation implements the DriverServerImpl interface.
func (s *DriverServerImpl) ArriveLocation(ctx context.Context, req *driver0.ArriveLocationReq) (resp *driver0.ArriveLocationResp, err error) {
	var order model.LxhOrder
	err = global.DB.Where("id = ? AND driver = ? AND order_status", req.OrderID, req.DriverID, "载客中").First(&order).Error
	if err != nil {
		return &driver0.ArriveLocationResp{
			Message: "订单不存在或司机状态不正确",
			Code:    http.StatusBadRequest,
		}, err
	}

	err = global.DB.Model(&order).Update("order_status", "完成").Error
	if err != nil {
		return &driver0.ArriveLocationResp{
			Message: "订单状态更新失败",
			Code:    http.StatusBadRequest,
		}, err
	}

	return &driver0.ArriveLocationResp{
		Message: "订单状态更新成功",
		Code:    http.StatusOK,
	}, err
}

// 开启行程
func (s *DriverServerImpl) StartTrip(ctx context.Context, req *driver0.StartTripReq) (resp *driver0.StartTripResp, err error) {
	// TODO: Your code here...
	return
}

// 结束行程
func (s *DriverServerImpl) EndTrip(ctx context.Context, req *driver0.EndTripReq) (resp *driver0.EndTripResp, err error) {
	// TODO: Your code here...
	return
}
