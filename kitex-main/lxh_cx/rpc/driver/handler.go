package main

import (
	"context"
	"encoding/json"
	"lxh_cx/config"
	"lxh_cx/kitex_gen/lxh_cx/base"
	driver0 "lxh_cx/kitex_gen/lxh_cx/driver"
	"lxh_cx/model"
	"net/http"
	"time"
)

// DriverServiceImpl implements the last service interface defined in the IDL.
type DriverServiceImpl struct{}

// DriverRegister implements the DriverServiceImpl interface.
func (s *DriverServiceImpl) DriverRegister(ctx context.Context, req *driver0.DriverRegisterReq) (resp *driver0.DriverRegisterResp, err error) {
	//司机注册发送审核
	//查看是否已经发送过审核了
	var drivers model.LxhDriverCheck
	err = config.Db.Where("driver_id = ?", req.DriverId).Find(&drivers).Limit(1).Error
	if err != nil {
		return &driver0.DriverRegisterResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "find driver failed",
			},
		}, err
	}

	if drivers.Id != 0 {
		return &driver0.DriverRegisterResp{BaseResp: &base.BaseResp{
			Code: http.StatusBadRequest,
			Msg:  "driver already exists",
		}}, nil
	}

	drivers = model.LxhDriverCheck{
		IdCardFileId:         req.IdCardFileId,
		DriverLicenseFileId:  req.DriverLicenseFileId,
		DrivingLicenseFileId: req.DrivingLicenseFileId,
		AvatorFileId:         req.AvatorFileId,
		CheckStatus:          "1",
		DriverId:             req.DriverId,
	}

	err = config.Db.Create(&drivers).Error
	if err != nil {
		return &driver0.DriverRegisterResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "driver register failed",
			},
		}, err
	}

	return &driver0.DriverRegisterResp{BaseResp: &base.BaseResp{
		Code: http.StatusOK,
		Msg:  "driver register success",
	}}, nil

}

// DriverAdd implements the DriverServiceImpl interface.
func (s *DriverServiceImpl) DriverAdd(ctx context.Context, req *driver0.DriverAddReq) (resp *driver0.DriverAddResp, err error) {
	//司机用户的添加
	result, _ := config.Rdb.Get(context.Background(), "passenger"+string(req.PassengerId)+":"+"order").Result()

	var orderMap map[string]interface{}
	json.Unmarshal([]byte(result), &orderMap)

	orders := model.LxhOrder{
		OrderCode:   orderMap["orderCode"].(string),
		Amount:      0,
		OrderStatus: orderMap["orderStatus"].(string),
		PassengerId: orderMap["passengerId"].(int32),
		StartAddr:   orderMap["startAddr"].(string),
		EndEnd:      orderMap["endEnd"].(string),
		Driver:      req.DriverId,
		StartTime:   time.Now(),
		PayStatus:   orderMap["payStatus"].(string),
	}

	err = config.Db.Create(&orders).Error
	if err != nil {
		return &driver0.DriverAddResp{BaseResp: &base.BaseResp{
			Code: http.StatusBadRequest,
			Msg:  "driver add failed",
		}}, nil
	}

	return &driver0.DriverAddResp{BaseResp: &base.BaseResp{
		Code: http.StatusOK,
		Msg:  "driver add success",
	}}, nil

}

// DriverOverOrder implements the DriverServiceImpl interface.
func (s *DriverServiceImpl) DriverOverOrder(ctx context.Context, req *driver0.DriverOverOrderReq) (resp *driver0.DriverOverOrderResp, err error) {
	//司机结束订单
	var orders model.LxhOrder
	err = config.Db.Where("id = ?", req.OrderId).Find(&orders).Limit(1).Error
	if err != nil {
		return &driver0.DriverOverOrderResp{BaseResp: &base.BaseResp{
			Code: http.StatusBadRequest,
			Msg:  "order not found",
		}}, nil
	}

	orders = model.LxhOrder{
		OrderStatus: "2",
		EndTime:     time.Now(),
		PayStatus:   "1",
	}

	err = config.Db.Where("id = ?", req.OrderId).Updates(&orders).Error
	if err != nil {
		return &driver0.DriverOverOrderResp{BaseResp: &base.BaseResp{
			Code: http.StatusBadRequest,
			Msg:  "over order failed",
		}}, nil
	}

	return &driver0.DriverOverOrderResp{
		BaseResp: &base.BaseResp{
			Code: http.StatusOK,
			Msg:  "order over success",
		},
	}, nil

}

// DriverCancelOrder implements the DriverServiceImpl interface.
func (s *DriverServiceImpl) DriverCancelOrder(ctx context.Context, req *driver0.DriverCancelOrderReq) (resp *driver0.DriverCancelOrderResp, err error) {
	//司机取消订单
	var orders model.LxhOrder
	err = config.Db.Where("id = ?", req.OrderId).Find(&orders).Limit(1).Error
	if err != nil {
		return &driver0.DriverCancelOrderResp{BaseResp: &base.BaseResp{
			Code: http.StatusBadRequest,
			Msg:  "find order failed",
		}}, err
	}

	orders = model.LxhOrder{
		OrderStatus:   "3",
		Driver:        req.DriverId,
		ConfirmBy:     "2",
		ConfirmPerson: req.ConfirmPerson,
		ConfirmReason: req.ConfirmReason,
		ConfirmRemark: req.ConfirmRemark,
	}

	err = config.Db.Where("id = ?", req.OrderId).Updates(&orders).Error
	if err != nil {
		return &driver0.DriverCancelOrderResp{BaseResp: &base.BaseResp{
			Code: http.StatusBadRequest,
			Msg:  "order cancel failed",
		}}, err
	}

	return &driver0.DriverCancelOrderResp{BaseResp: &base.BaseResp{
		Code: http.StatusOK,
		Msg:  "order cancel success",
	}}, nil

}

// DriverInfoList implements the DriverServiceImpl interface.
func (s *DriverServiceImpl) DriverInfoList(ctx context.Context, req *driver0.DriverInfoListReq) (resp *driver0.DriverInfoListResp, err error) {
	//司机首页展示的内容
	var drivers model.LxhDriver
	err = config.Db.Where("id = ?", req.DriverId).Find(&drivers).Limit(1).Error
	if err != nil {
		return &driver0.DriverInfoListResp{BaseResp: &base.BaseResp{
			Code: http.StatusBadRequest,
			Msg:  "driver find failed",
		}}, nil
	}

	if drivers.Id == 0 {
		return &driver0.DriverInfoListResp{BaseResp: &base.BaseResp{
			Code: http.StatusBadRequest,
			Msg:  "driver not exist",
		}}, nil
	}

	info := driver0.DriverInfo{
		Name:        drivers.Name,
		NickName:    drivers.NickName,
		AllAcount:   drivers.AllAcount,
		CarAge:      string(drivers.CarAge),
		HandlerFile: string(drivers.FileId),
	}

	return &driver0.DriverInfoListResp{
		BaseResp: &base.BaseResp{
			Code: http.StatusOK,
			Msg:  "driver info list success",
		},
		DriverInfo: &info,
	}, nil

}
