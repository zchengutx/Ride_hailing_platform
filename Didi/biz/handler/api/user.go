package api

import (
	"Didi/biz/dal/global"
	"Didi/biz/handler/request"
	"Didi/biz/middleware"
	"Didi/kitex_gen/Didi/user"
	pb "Didi/kitex_gen/Didi/user/userserver"
	"Didi/utils"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func SendSms(ctx context.Context, c *app.RequestContext) {
	var req request.SendSmsReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "查询失败",
			"data": err.Error(),
		})
		return
	}
	if !utils.ValidateMobile(req.Mobile) {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "手机号码格式不正确",
			"data": nil,
		})
		return
	}
	newClient, _ := pb.NewClient("user", client.WithHostPorts("127.0.0.1:50051"))
	sms, _ := newClient.SendSms(ctx, &user.SendSmsReq{
		Mobile:      req.Mobile,
		SendSmsCode: req.SendSmsCode,
	})
	c.JSON(200, sms)
}
func Login(ctx context.Context, c *app.RequestContext) {
	var req request.LoginReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "查询失败",
			"data": err.Error(),
		})
		return
	}
	if !utils.ValidateMobile(req.Mobile) {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "手机号码格式不正确",
			"data": nil,
		})
		return
	}
	newClient, _ := pb.NewClient("user", client.WithHostPorts("127.0.0.1:50051"))
	loginUser, _ := newClient.LoginUser(ctx, &user.LoginUserReq{
		Mobile:      req.Mobile,
		SendSmsCode: req.SendSmsCode,
	})
	if loginUser.Code != 200 {
		c.JSON(200, loginUser)
	}
	token, _ := middleware.NewJWT(global.JWT_SELECT_KEY).CreateToken(middleware.CustomClaims{
		ID: int(loginUser.UId),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 86400,
		},
	})
	c.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "登录成功",
		"data": map[string]string{
			"token": token,
		},
	})
}

func RealName(ctx context.Context, c *app.RequestContext) {
	var req request.RealNameReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "查询失败",
			"data": err.Error(),
		})
		return
	}
	if !utils.IdCardVerification(req.IdCard) {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "身份证格式不正确",
			"data": nil,
		})
		return
	}
	newClient, _ := pb.NewClient("user", client.WithHostPorts("127.0.0.1:50051"))
	name, _ := newClient.RealName(ctx, &user.RealNameReq{
		UId:      int64(c.GetInt("userId")),
		UserName: req.UserName,
		Sex:      req.Sex,
		Age:      req.Age,
		IdCard:   req.IdCard,
	})
	c.JSON(200, name)
}
