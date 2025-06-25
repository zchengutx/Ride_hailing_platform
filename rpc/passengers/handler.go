package main

import (
	"Ride_hailing_platform/basic/global"
	"Ride_hailing_platform/handler/model"
	"os"

	passengers "Ride_hailing_platform/kitex_gen/passengers"
	"Ride_hailing_platform/utils"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"time"
)

// PassengersServerImpl implements the last service interface defined in the IDL.
type PassengersServerImpl struct{}

// 发送短信验证码
func (s *PassengersServerImpl) SendSms(ctx context.Context, req *passengers.SendSmsReq) (resp *passengers.SendSmsResp, err error) {

	smsCode := rand.Intn(9000) + 1000

	err = global.Redis.Set(ctx, req.Mobile, smsCode, 5*time.Minute).Err()
	if err != nil {
		return &passengers.SendSmsResp{
			Message: "发送验证码失败",
			Code:    500,
		}, nil
	}

	return &passengers.SendSmsResp{
		Message: "验证码发送成功",
		Code:    200,
	}, nil
}

// 注册功能
func (s *PassengersServerImpl) Register(ctx context.Context, req *passengers.RegisterReq) (resp *passengers.RegisterResp, err error) {
	// TODO: 注册功能
	get := global.Redis.Get(context.Background(), "sendSms"+req.Mobile+"register")
	if get.Val() != req.SmsCode {
		resp = &passengers.RegisterResp{
			Message: "短信验证码错误，请查证后输入",
			Code:    http.StatusBadRequest,
		}
	}
	newUser := &model.LxhPassenger{
		Mobile:   req.Mobile,
		NickName: "滴滴用户" + req.Mobile[7:],
	}
	if err := global.DB.Create(newUser); err != nil {
		resp = &passengers.RegisterResp{
			Message: "用户创建失败",
			Code:    http.StatusBadRequest,
		}
	}
	return &passengers.RegisterResp{
		Message: "ok",
		Code:    http.StatusOK,
	}, nil
}

// 登录功能
func (s *PassengersServerImpl) Login(ctx context.Context, req *passengers.LoginReq) (resp *passengers.LoginResp, err error) {
	result, err := global.Redis.Get(ctx, req.Mobile).Result()
	if err != nil {
		if err == redis.Nil {
			return &passengers.LoginResp{
				Message: "验证码已过期或不存在",
				Code:    400,
			}, nil
		}
		return &passengers.LoginResp{
			Message: "验证码验证失败",
			Code:    500,
		}, nil
	}

	if result != req.SmsCode {
		return &passengers.LoginResp{
			Message: "验证码错误",
			Code:    400,
		}, nil
	}

	// 查找用户
	var userModel model.LxhPassenger
	err = global.DB.Where("mobile = ?", req.Mobile).First(&userModel).Error
	if err != nil {
		return &passengers.LoginResp{
			Message: "登录失败",
			Code:    500,
		}, nil
	}

	// 直接生成登录Token
	token := fmt.Sprintf("token_%s_%d", userModel.Id, time.Now().Unix())

	utils.InitDefaultJWTManager(&utils.JWTConfig{
		SecretKey:        "",
		ExpiresIn:        24 * time.Hour,
		RefreshExpiresIn: 7 * 24 * time.Hour,
	})

	generateToken, _ := utils.GenerateToken(1234, req.Mobile, token)

	// 删除已使用的验证码
	global.Redis.Del(ctx, req.Mobile)

	return &passengers.LoginResp{
		Message: "登录成功",
		Code:    200,
		Token:   generateToken,
	}, nil
}

// 实时定位
func (s *PassengersServerImpl) GetLocation(ctx context.Context, req *passengers.LocationReq) (resp *passengers.LocationResp, err error) {
	key := fmt.Sprintf("location:%v", req.PassengerID)
	latStr, err1 := global.Redis.Get(ctx, key+":lat").Result()
	lngStr, err2 := global.Redis.Get(ctx, key+":lng").Result()
	if err1 != nil || err2 != nil {
		return &passengers.LocationResp{
			Message: "未获取到定位信息",
			Code:    404,
		}, nil
	}
	return &passengers.LocationResp{
		Message: "获取定位成功",
		Code:    200,
		Lat:     utils.ParseFloat(latStr),
		Lng:     utils.ParseFloat(lngStr),
	}, nil
}

// 行程
func (s *PassengersServerImpl) GetTripInfo(ctx context.Context, req *passengers.TripReq) (resp *passengers.TripResp, err error) {
	var order []model.LxhOrder
	err = global.DB.Where("passenger_id = ?", req.PassengerID).Order("start_time desc").Limit(10).Find(&order).Error
	if err != nil {
		return &passengers.TripResp{
			Message: "历史记录加载失败",
			Code:    http.StatusInternalServerError,
		}, err
	}
	var history []string
	for _, v := range order {
		history = append(history, v.StartAddr+v.EndEnd)
	}
	dests := []string{}
	for _, o := range order {
		dests = append(dests, o.EndEnd)
		if len(dests) >= 10 {
			break
		}
	}

	return &passengers.TripResp{
		Message:               "获取行程成功",
		Code:                  http.StatusOK,
		HistoryTrips:          history,
		RecommendDestinations: dests,
	}, nil
}

// 计价
func (s *PassengersServerImpl) GetPricing(ctx context.Context, req *passengers.PricingReq) (resp *passengers.PricingResp, err error) {

	return
}

// 下单
func (s *PassengersServerImpl) PlaceOrder(ctx context.Context, req *passengers.OrderReq) (resp *passengers.OrderResp, err error) {

	return
}

// 路线规划与费用
func (s *PassengersServerImpl) RoutePlan(ctx context.Context, req *passengers.RoutePlanReq) (resp *passengers.RoutePlanResp, err error) {
	ak := os.Getenv("BAIDU_MAP_AK")
	if ak == "" {
		ak = "k8Q4DliDfEbE1dOs3TgacVcHW4AIvdup"
	}
	distance, duration, err := utils.GetRouteFromBaidu(
		req.StartLat, req.StartLng, req.EndLat, req.EndLng, ak,
	)
	if err != nil {
		return &passengers.RoutePlanResp{
			Message: "路径规划失败",
			Code:    500,
		}, nil
	}
	price := 10.0 + float64(distance)/1000*2.0
	if req.CarType == "快车" {
		price *= 1.2
	}
	return &passengers.RoutePlanResp{
		Message:  "路径规划成功",
		Code:     200,
		Distance: int32(distance),
		Duration: int32(duration),
		Price:    price,
	}, nil
}

// 微信登录
func (s *PassengersServerImpl) WechatLogin(ctx context.Context, req *passengers.WechatLoginReq) (resp *passengers.WechatLoginResp, err error) {
	// 根据OpenID查找用户
	var userModel model.LxhPassenger
	err = global.DB.Where("wechat_openid = ?", req.OpenID).First(&userModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//不存在则创建
			userID := uuid.New().ID()
			User := model.LxhPassenger{
				Id:             int32(userID),
				Name:           req.Nickname,
				WechatOpenID:   string(req.OpenID),
				WechatUnionID:  string(req.UnionID),
				WechatNickname: req.Nickname,
				WechatAvatar:   req.Avatar,
			}
			err = global.DB.Create(&User).Error
			if err != nil {
				return &passengers.WechatLoginResp{
					Message: "微信登录失败",
					Code:    500,
				}, nil
			}
			userModel = User
		} else {
			return &passengers.WechatLoginResp{
				Message: "微信登录失败",
				Code:    500,
			}, nil
		}
	}

	// 直接生成登录Token
	token := fmt.Sprintf("token_%s_%d", userModel.Id, time.Now().Unix())

	utils.InitDefaultJWTManager(&utils.JWTConfig{
		SecretKey:        "",
		ExpiresIn:        24 * time.Hour,
		RefreshExpiresIn: 7 * 24 * time.Hour,
	})

	generateToken, _ := utils.GenerateToken(1234, string(req.OpenID), token)

	return &passengers.WechatLoginResp{
		Message: "微信登录成功",
		Code:    200,
		Token:   generateToken,
	}, nil
}

// 绑定微信
func (s *PassengersServerImpl) BindWechat(ctx context.Context, req *passengers.BindWechatReq) (resp *passengers.BindWechatResp, err error) {
	var userModel model.LxhPassenger
	err = global.DB.Where("id = ?", req.UserID).First(&userModel).Error
	if err != nil {
		return &passengers.BindWechatResp{
			Message: "绑定失败",
			Code:    500,
		}, nil
	}

	now := time.Now()
	updates := map[string]interface{}{
		"wechat_openid":   req.OpenID,
		"wechat_unionid":  req.UnionID,
		"wechat_nickname": req.Nickname,
		"wechat_avatar":   req.Avatar,
		"wechat_bound_at": &now,
		"update_time":     now,
	}

	err = global.DB.Model(&userModel).Updates(updates).Error
	if err != nil {
		return &passengers.BindWechatResp{
			Message: "绑定失败",
			Code:    500,
		}, nil
	}

	// 重新查询用户信息
	err = global.DB.Where("id = ?", req.UserID).First(&userModel).Error
	if err != nil {
		return &passengers.BindWechatResp{
			Message: "绑定失败",
			Code:    500,
		}, nil
	}

	return &passengers.BindWechatResp{
		Message: "绑定成功",
		Code:    200,
	}, nil
}

// UnbindWechat implements the PassengersServerImpl interface.
func (s *PassengersServerImpl) UnbindWechat(ctx context.Context, req *passengers.UnbindWechatReq) (resp *passengers.UnbindWechatResp, err error) {
	// TODO: Your code here...
	return
}

// GetUserInfo implements the PassengersServerImpl interface.
func (s *PassengersServerImpl) GetUserInfo(ctx context.Context, req *passengers.GetUserInfoReq) (resp *passengers.GetUserInfoResp, err error) {
	// 查找用户
	var userModel model.LxhPassenger
	err = global.DB.Where("id = ?", req.UserID).First(&userModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &passengers.GetUserInfoResp{
				Message: "用户不存在",
				Code:    404,
			}, nil
		}
		return &passengers.GetUserInfoResp{
			Message: "获取用户信息失败",
			Code:    500,
		}, nil
	}

	return &passengers.GetUserInfoResp{
		Message: "获取成功",
		Code:    200,
	}, nil
}

// 更新用户信息
func (s *PassengersServerImpl) UpdateUserInfo(ctx context.Context, req *passengers.UpdateUserInfoReq) (resp *passengers.UpdateUserInfoResp, err error) {
	// 查找用户
	var userModel model.LxhPassenger
	err = global.DB.Where("id = ?", req.UserID).First(&userModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &passengers.UpdateUserInfoResp{
				Message: "用户不存在",
				Code:    404,
			}, nil
		}
		return &passengers.UpdateUserInfoResp{
			Message: "更新用户信息失败",
			Code:    500,
		}, nil
	}

	now := time.Now()
	updates := map[string]interface{}{
		"update_time": now,
	}

	err = global.DB.Model(&userModel).Updates(updates).Error
	if err != nil {
		return &passengers.UpdateUserInfoResp{
			Message: "更新用户信息失败",
			Code:    500,
		}, nil
	}

	err = global.DB.Where("id = ?", req.UserID).First(&userModel).Error
	if err != nil {
		return &passengers.UpdateUserInfoResp{
			Message: "更新用户信息失败",
			Code:    500,
		}, nil
	}

	return &passengers.UpdateUserInfoResp{
		Message: "更新成功",
		Code:    200,
	}, nil
}
