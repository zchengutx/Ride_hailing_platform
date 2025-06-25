package Ride_hailing_platform

import (
	"Ride_hailing_platform/basic/global"
	"Ride_hailing_platform/handler/model"
	driver0 "Ride_hailing_platform/kitex_gen/driver"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

// StartTrip implements the DriverServerImpl interface.
func (s *DriverServerImpl) StartTrip(ctx context.Context, req *driver0.StartTripReq) (resp *driver0.StartTripResp, err error) {
	orderID, err := strconv.ParseInt(req.OrderID, 10, 32)
	if err != nil {
		return &driver0.StartTripResp{
			Message: "订单ID格式错误",
			Code:    http.StatusBadRequest,
		}, nil
	}

	driverID, err := strconv.ParseInt(req.DriverID, 10, 32)
	if err != nil {
		return &driver0.StartTripResp{
			Message: "司机ID格式错误",
			Code:    http.StatusBadRequest,
		}, nil
	}

	// 检查订单是否存在且可以开始行程
	var order model.LxhOrder
	err = global.DB.Where("id = ? AND driver = ? AND order_status IN (?)",
		orderID, driverID, []string{"accepted", "arrived"}).First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &driver0.StartTripResp{
				Message: "订单不存在或状态不正确",
				Code:    http.StatusNotFound,
			}, nil
		}
		return &driver0.StartTripResp{
			Message: "查询订单失败",
			Code:    http.StatusInternalServerError,
		}, err
	}

	// 更新订单状态和开始时间
	now := time.Now()
	err = global.DB.Model(&order).Updates(map[string]interface{}{
		"order_status": "started",
		"start_time":   now,
	}).Error
	if err != nil {
		return &driver0.StartTripResp{
			Message: "开始行程失败",
			Code:    http.StatusInternalServerError,
		}, err
	}

	return &driver0.StartTripResp{
		Message: "行程已开始",
		Code:    http.StatusOK,
	}, nil
}

// EndTrip implements the DriverServerImpl interface.
func (s *DriverServerImpl) EndTrip(ctx context.Context, req *driver0.EndTripReq) (resp *driver0.EndTripResp, err error) {
	orderID, err := strconv.ParseInt(req.OrderID, 10, 32)
	if err != nil {
		return &driver0.EndTripResp{
			Message: "订单ID格式错误",
			Code:    http.StatusBadRequest,
		}, nil
	}

	driverID, err := strconv.ParseInt(req.DriverID, 10, 32)
	if err != nil {
		return &driver0.EndTripResp{
			Message: "司机ID格式错误",
			Code:    http.StatusBadRequest,
		}, nil
	}

	// 检查订单是否存在且正在进行中
	var order model.LxhOrder
	err = global.DB.Where("id = ? AND driver = ? AND order_status = ?",
		orderID, driverID, "started").First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &driver0.EndTripResp{
				Message: "订单不存在或状态不正确",
				Code:    http.StatusNotFound,
			}, nil
		}
		return &driver0.EndTripResp{
			Message: "查询订单失败",
			Code:    http.StatusInternalServerError,
		}, err
	}

	// 开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新订单状态和结束时间
	now := time.Now()
	err = tx.Model(&order).Updates(map[string]interface{}{
		"order_status": "completed",
		"end_time":     now,
		"pay_status":   "pending",
	}).Error
	if err != nil {
		tx.Rollback()
		return &driver0.EndTripResp{
			Message: "结束行程失败",
			Code:    http.StatusInternalServerError,
		}, err
	}

	// 更新司机状态为在线，准备接新订单
	err = tx.Model(&model.LxhDriver{}).Where("id = ?", driverID).Update("status", "online").Error
	if err != nil {
		tx.Rollback()
		return &driver0.EndTripResp{
			Message: "更新司机状态失败",
			Code:    http.StatusInternalServerError,
		}, err
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return &driver0.EndTripResp{
			Message: "结束行程失败",
			Code:    http.StatusInternalServerError,
		}, err
	}

	return &driver0.EndTripResp{
		Message: "行程已结束",
		Code:    http.StatusOK,
	}, nil
}
