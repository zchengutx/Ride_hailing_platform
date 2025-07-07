package main

import (
	pb "Didi/kitex_gen/Didi/driver"
	"Didi/rpc/common/global"
	"Didi/rpc/common/model"
	"context"
)

// DriverServerImpl implements the last service interface defined in the IDL.
type DriverServerImpl struct{}

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
