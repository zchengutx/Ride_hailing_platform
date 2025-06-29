package main

import (
	"context"
	"fmt"
	"lxh_cx/config"
	"lxh_cx/kitex_gen/lxh_cx/base"
	user "lxh_cx/kitex_gen/lxh_cx/user"
	"lxh_cx/model"
	"lxh_cx/utils"
	"net/http"
	"time"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// SendSms implements the UserServiceImpl interface.
func (s *UserServiceImpl) SendSms(ctx context.Context, req *user.SendSmsReq) (resp *user.SendSmsResp, err error) {
	// TODO: Your code here...
	//生成一个随机数，作为验证码
	code, err := utils.GenerateUniqueSixDigitString()
	if err != nil {
		return &user.SendSmsResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "send sms failed",
			},
		}, err
	}
	//sms, err := utils.SendSms(req.Mobile, code)
	//if err != nil {
	//	return &user.SendSmsResp{
	//		BaseResp: &base.BaseResp{
	//			Code: http.StatusBadRequest,
	//			Msg:  "send sms failed",
	//		},
	//	}, err
	//}
	//if *sms.Body.Code != "OK" {
	//	return &user.SendSmsResp{
	//		BaseResp: &base.BaseResp{
	//			Code: http.StatusBadRequest,
	//			Msg:  *sms.Body.Message,
	//		},
	//	}, err
	//}
	//fmt.Println(req)
	config.Rdb.Set(context.Background(), "sendSms"+req.Mobile+req.Source, code, time.Minute*5)

	return &user.SendSmsResp{
		BaseResp: &base.BaseResp{
			Code: http.StatusOK,
			Msg:  "send sms successfully",
		},
	}, err

}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	// 从Redis中获取验证码
	key := "sendSms" + req.Mobile + "register"
	code, err := config.Rdb.Get(ctx, key).Result()
	if err != nil {
		return &user.RegisterResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "verification code not found or expired",
			},
		}, nil
	}

	// 验证验证码
	if code != req.SendCode {
		return &user.RegisterResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "invalid verification code",
			},
		}, nil
	}

	digitString, _ := utils.GenerateUniqueSixDigitString()
	users := model.LxhPassenger{
		Mobile:   req.Mobile,
		NickName: "用户" + digitString,
	}

	err = config.Db.Create(&users).Error
	if err != nil {
		return &user.RegisterResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "user register failed",
			},
		}, nil
	}
	// 注册成功后删除验证码
	config.Rdb.Del(ctx, key)

	return &user.RegisterResp{
		BaseResp: &base.BaseResp{
			Code: http.StatusOK,
			Msg:  "register successfully",
		},
	}, nil

}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {

	key := "sendSms" + req.Mobile + "login"
	code, err := config.Rdb.Get(ctx, key).Result()
	if err != nil {
		return &user.LoginResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "verification code not found or expired",
			},
		}, nil
	}

	// 验证验证码
	if code != req.SendCode {
		return &user.LoginResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "invalid verification code",
			},
		}, nil
	}
	//在数据库查询是否存在用户
	var users model.LxhPassenger
	err = config.Db.Where("mobile = ?", req.Mobile).Find(&users).Limit(1).Error
	if err != nil {
		return &user.LoginResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "find users failed",
			},
		}, nil
	}
	//如果用户不存在时
	if users.Id == 0 {
		digitString, _ := utils.GenerateUniqueSixDigitString()
		users = model.LxhPassenger{
			Mobile:   req.Mobile,
			NickName: "用户" + digitString,
		}
		err = config.Db.Create(&users).Error
		if err != nil {
			return &user.LoginResp{
				BaseResp: &base.BaseResp{
					Code: http.StatusBadRequest,
					Msg:  "user register failed",
				},
			}, nil
		}
	}

	return &user.LoginResp{
		BaseResp: &base.BaseResp{
			Code: http.StatusOK,
			Msg:  "user login successfully",
		},
		Id: int64(users.Id),
	}, nil

}

// UserInfoList implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserInfoList(ctx context.Context, req *user.UserInfoListReq) (resp *user.UserInfoListResp, err error) {

	//查询用户详情
	var users model.LxhPassenger
	err = config.Db.Where("id = ?", req.Id).Find(&users).Error
	if err != nil {
		return &user.UserInfoListResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "find users failed",
			},
		}, nil
	}
	//查询用户是否存在
	if users.Id == 0 {
		return &user.UserInfoListResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "user not exist",
			},
		}, nil
	}

	userInfo := &user.UserInfoResp{
		NickName: users.NickName,
		Sex:      users.Sex,
		Mileage:  users.Mileage,
	}

	return &user.UserInfoListResp{
		BaseResp: &base.BaseResp{
			Code: http.StatusOK,
			Msg:  "user info list successfully",
		},
		UserInfo: userInfo,
	}, nil

}

// BindMobile implements the UserServiceImpl interface.
func (s *UserServiceImpl) BindMobile(ctx context.Context, req *user.BindMobileReq) (resp *user.BindMobileResp, err error) {

	//用户绑定第三方登录
	var users model.LxhPassenger
	err = config.Db.Where("mobile = ?", req.Mobile).Find(&users).Error
	if err != nil {
		return &user.BindMobileResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "user find failed",
			},
		}, nil
	}

	if users.Id == 0 {
		return &user.BindMobileResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "user not exist",
			},
		}, nil
	}
	//查询是否第三方已登录
	var auth model.LxhAuth
	err = config.Db.Where("id = ?", req.Id).Find(&auth).Limit(1).Error
	if err != nil {
		return &user.BindMobileResp{BaseResp: &base.BaseResp{
			Code: http.StatusBadRequest,
			Msg:  "auth find failed",
		}}, nil
	}

	if auth.Id == 0 {
		return &user.BindMobileResp{BaseResp: &base.BaseResp{
			Code: http.StatusBadRequest,
			Msg:  "auth not exist",
		}}, nil
	}
	//进行绑定第三方
	userAuth := model.LxhUserAuth{
		PassengerId: users.Id,
		AuthId:      int32(req.Id),
	}

	err = config.Db.Create(&userAuth).Error
	if err != nil {
		return &user.BindMobileResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "userAuth create failed",
			},
		}, nil
	}

	return &user.BindMobileResp{
		BaseResp: &base.BaseResp{
			Code: http.StatusOK,
			Msg:  "user info bind successfully",
		},
	}, nil

}

// UpdateUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateUserInfo(ctx context.Context, req *user.UpdateUserInfoReq) (resp *user.UpdateUserInfoResp, err error) {
	//修改用户的个人信息
	var users model.LxhPassenger
	err = config.Db.Where("id = ?", req.Id).Find(&users).Limit(1).Error
	if err != nil {
		return &user.UpdateUserInfoResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "find users failed",
			},
		}, nil
	}

	if users.Id == 0 {
		return &user.UpdateUserInfoResp{BaseResp: &base.BaseResp{
			Code: http.StatusBadRequest,
			Msg:  "user not exist",
		}}, nil
	}

	users = model.LxhPassenger{
		Mobile:   req.Mobile,
		NickName: req.NickName,
		Sex:      req.Sex,
	}

	err = config.Db.Where("id = ?", req.Id).Updates(&users).Error
	if err != nil {
		return &user.UpdateUserInfoResp{BaseResp: &base.BaseResp{
			Code: http.StatusBadRequest,
			Msg:  "update user failed",
		}}, nil
	}

	return &user.UpdateUserInfoResp{BaseResp: &base.BaseResp{
		Code: http.StatusOK,
		Msg:  "user info update successfully",
	}}, nil

}

// UpdateCancelOrder implements the UserServiceImpl interface.
func (s *UserServiceImpl) UpdateCancelOrder(ctx context.Context, req *user.UpdateCancelOrderReq) (resp *user.UpdateCancelOrderResp, err error) {
	//用户取消订单
	var orders model.LxhOrder
	err = config.Db.Where("id = ? and passenger_id = ?", req.OrderId, req.Id).Find(&orders).Limit(1).Error
	if err != nil {
		return &user.UpdateCancelOrderResp{BaseResp: &base.BaseResp{
			Code: http.StatusBadRequest,
			Msg:  "find order failed",
		}}, err
	}

	//取消订单时的订单状态为3
	orders = model.LxhOrder{
		OrderStatus:   "3",
		Driver:        req.Id,
		ConfirmBy:     "1",
		ConfirmPerson: req.ConfirmPerson,
		ConfirmReason: req.ConfirmReason,
		ConfirmRemark: req.ConfirmRemark,
	}

	err = config.Db.Where("id = ?", req.OrderId).Updates(&orders).Error
	if err != nil {
		return &user.UpdateCancelOrderResp{BaseResp: &base.BaseResp{
			Code: http.StatusBadRequest,
			Msg:  "order cancel failed",
		}}, err
	}

	return &user.UpdateCancelOrderResp{BaseResp: &base.BaseResp{
		Code: http.StatusOK,
		Msg:  "order cancel success",
	}}, nil

}

// FindRouteRecord implements the UserServiceImpl interface.
func (s *UserServiceImpl) FindRouteRecord(ctx context.Context, req *user.FindRouteRecordReq) (resp *user.FindRouteRecordResp, err error) {
	//取出用户的历史路线搜索前十
	key := fmt.Sprintf("passenger:%d", req.UserId)
	result, _ := config.Rdb.ZRevRange(ctx, key, 0, 9).Result()
	if result != nil {
		return &user.FindRouteRecordResp{
			BaseResp: &base.BaseResp{
				Code: http.StatusBadRequest,
				Msg:  "find route record failed",
			},
		}, nil
	}

	var route []*user.FindRouteRecord

	for _, rou := range result {
		route = append(route, &user.FindRouteRecord{
			Address: rou,
		})
	}

	return &user.FindRouteRecordResp{
		BaseResp: &base.BaseResp{
			Code: http.StatusOK,
			Msg:  "user find route record successfully",
		},
		Route: route,
	}, nil
}

// SetRouteUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) SetRouteUser(ctx context.Context, req *user.SetRouteUserReq) (resp *user.SetRouteUserResp, err error) {
	//存储用户已搜索的地址
	key := fmt.Sprintf("passenger:%d:%s", req.UserId, req.StartAddr)
	config.Rdb.ZIncrBy(ctx, key, 1, req.StartAddr)

	return &user.SetRouteUserResp{
		BaseResp: &base.BaseResp{
			Code: http.StatusOK,
			Msg:  "user set route successfully",
		},
	}, nil
}
