package main

import (
	"context"
	"kitexTwo/common/global"
	"kitexTwo/common/model"
	pb "kitexTwo/kitex_gen/kitexTwo/user"
)

// UserImpl implements the last service interface defined in the IDL.
type UserImpl struct {
}

// CreateUser implements the UserImpl interface.
func (s *UserImpl) CreateUser(_ context.Context, req *pb.CreateUserReq) (*pb.CreateUserResp, error) {
	// TODO: Your code here...
	code := "8769"
	var user model.User
	if err := global.DB.Debug().Where("mobile = ?", req.Title).Limit(1).Find(&user).Error; err != nil {
		return &pb.CreateUserResp{
			Message: "该用户已存在",
			Code:    603,
		}, nil
	}
	if code != req.SendSmsCode {
		return &pb.CreateUserResp{
			Message: "验证码不正确",
			Code:    603,
		}, nil
	}
	m := model.User{
		Mobile: req.Title,
	}
	if err := global.DB.Debug().Create(&m).Error; err != nil {
		return &pb.CreateUserResp{
			Message: "注册失败",
			Code:    403,
		}, nil
	}
	return &pb.CreateUserResp{
		Message: "注册成功",
		Code:    200,
	}, nil
}
func (s *UserImpl) LoginUser(_ context.Context, req *pb.LoginUserReq) (*pb.LoginUserResp, error) {
	code := "8769"
	var user model.User
	if err := global.DB.Debug().Where("mobile = ?", req.Title).Limit(1).Find(&user).Error; err != nil {
		return &pb.LoginUserResp{
			Message: "账号未注册",
			Code:    603,
		}, nil
	}
	if code != req.SendSmsCode {
		return &pb.LoginUserResp{
			Message: "验证码不存在",
			Code:    403,
		}, nil
	}
	return &pb.LoginUserResp{
		Message: "登录成功",
		Code:    200,
		Uid:     user.Id,
	}, nil
}
