package main

import (
	pb "Didi/kitex_gen/Didi/user"
	"Didi/rpc/common/global"
	"Didi/rpc/common/model"
	"context"
	"math/rand"
	"time"
)

// UserServerImpl implements the last service interface defined in the IDL.
type UserServerImpl struct{}

// SendSms implements the UserServerImpl interface.
func (s *UserServerImpl) SendSms(ctx context.Context, req *pb.SendSmsReq) (*pb.SendSmsResp, error) {
	// TODO: Your code here...
	ctx = context.Background()
	code := rand.Intn(900000) + 100000
	if err := global.Rdb.Set(ctx, "sendSms"+req.Mobile, code, time.Minute*5).Err(); err != nil {
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
	// TODO: Your code here...
	ctx = context.Background()
	login := global.Rdb.Get(ctx, "sendSms"+req.Mobile).Val()
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
	if login != req.SendSmsCode {
		return &pb.LoginUserResp{
			Code:    301,
			Message: "验证码错误",
		}, nil
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
