package main

import (
	"context"
	"kitex_main/config"
	user "kitex_main/kitex_gen/Ride_hailing_platform/user"
	"kitex_main/model"
	"math/rand"
	"time"
)

// UserServerImpl implements the last service interface defined in the IDL.
type UserServerImpl struct{}

func (s *UserServerImpl) CallBack(ctx context.Context, req *user.CallBackReq) (r *user.CallBackResp, err error) {
	//TODO implement me
	var auth model.Auth
	var users model.LxhPassenger

	err = auth.FindAuthOpenId(req.OpenId)
	if err != nil {
		auth = model.Auth{
			AppName:    "微信",
			AppUnionid: req.OpenId,
		}
		users = model.LxhPassenger{
			NickName: req.NickName,
		}
		err = auth.CreateAuth()
		err = users.CreateUser()

		if err != nil {
			return &user.CallBackResp{
				Code: 400,
				Msg:  "create auth failed",
			}, nil
		}

	}

	return &user.CallBackResp{
		Code:  200,
		Msg:   "login success",
		Token: int64(users.Id),
	}, nil
}

// SendSms implements the UserServerImpl interface.
func (s *UserServerImpl) SendSms(ctx context.Context, req *user.SendSmsReq) (resp *user.SendSmsResp, err error) {
	// TODO: Your code here...
	code := rand.Intn(900000) + 100000

	//sms, err := pkg.SendSms(req.Mobile, strconv.Itoa(code))
	//if err != nil {
	//	return &user.SendSmsResp{
	//		Code: 400,
	//		Msg:  "send sms failed",
	//	}, nil
	//}
	//
	//if *sms.Body.Code != "OK" {
	//	return &user.SendSmsResp{
	//		Code: 400,
	//		Msg:  *sms.Body.Message,
	//	}, nil
	//}

	config.RDB.Set(config.Ctx, "sendSms"+req.Mobile+req.Source, code, 10*time.Minute)

	return &user.SendSmsResp{
		Code: 200,
	}, nil
}

// Login implements the UserServerImpl interface.
func (s *UserServerImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	// TODO: Your code here...

	get := config.RDB.Get(config.Ctx, "sendSms"+req.Mobile+"login")

	if get.Val() != req.SendSmsCode {
		return &user.LoginResp{
			Code: 400,
			Msg:  "send sms failed",
		}, err
	}

	var User model.LxhPassenger
	err = User.FindUserMobile(req.Mobile)
	if err != nil {
		User.Mobile = req.Mobile
		err = User.CreateUser()
		if err != nil {
			return &user.LoginResp{
				Code: 400,
				Msg:  "register failed",
			}, err
		}
	}

	return &user.LoginResp{
		Code:  200,
		Token: int64(User.Id),
	}, nil
}

// InfoUser implements the UserServerImpl interface.
func (s *UserServerImpl) InfoUser(ctx context.Context, req *user.InfoUserReq) (resp *user.InfoUserResp, err error) {
	// TODO: Your code here...

	var User model.LxhPassenger

	err = User.FindUserId(int(req.UserId))
	if err != nil {
		return &user.InfoUserResp{
			Code: 400,
		}, err
	}

	var InfoUser *user.InfoUser

	InfoUser = &user.InfoUser{
		UserName: User.Name,
		Mobile:   User.Mobile,
		NikeName: User.NickName,
		Sex:      User.Sex,
		Mileage:  User.Mileage,
		Age:      string(User.Age),
	}

	return &user.InfoUserResp{
		Code: 200,
		Info: InfoUser,
	}, nil
}
