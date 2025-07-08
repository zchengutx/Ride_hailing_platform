package api

import (
	"Didi/biz/dal/global"
	"Didi/biz/handler/request"
	"Didi/biz/middleware"
	"Didi/kitex_gen/Didi/user"
	"Didi/utils"
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/dgrijalva/jwt-go"
)

var (
	// 使用封装的客户端管理器获取用户服务客户端，避免重复创建连接
	UserClient = utils.GetDefaultUserClient()
)

func SendSms(ctx context.Context, c *app.RequestContext) {
	var req request.SendSmsReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}
	if !utils.ValidateMobile(req.Mobile) {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "手机号码格式不正确",
			"data": nil,
		})
		return
	}

	// 调用用户服务发送短信验证码
	sms, err := UserClient.SendSms(ctx, &user.SendSmsReq{
		Mobile:      req.Mobile,
		SendSmsCode: req.SendSmsCode,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	// 根据业务响应动态返回HTTP状态码
	httpCode := 200
	if sms.Code != 200 {
		httpCode = 400
	}
	c.JSON(httpCode, sms)
}
func Login(ctx context.Context, c *app.RequestContext) {
	var req request.LoginReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}
	if !utils.ValidateMobile(req.Mobile) {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "手机号码格式不正确",
			"data": nil,
		})
		return
	}

	// 调用用户服务进行登录验证
	loginUser, err := UserClient.LoginUser(ctx, &user.LoginUserReq{
		Mobile:      req.Mobile,
		SendSmsCode: req.SendSmsCode,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	// 登录失败，直接返回业务错误
	if loginUser.Code != 200 {
		httpCode := 400
		if loginUser.Code >= 500 {
			httpCode = 500
		}
		c.JSON(httpCode, loginUser)
		return
	}

	// 动态获取JWT配置的过期时间 - 避免硬编码
	jwtConfig := global.JWT_SELECT_KEY
	defaultExpireHours := 24 // 默认24小时
	expirationTime := time.Now().Unix() + int64(defaultExpireHours*3600)

	token, err := middleware.NewJWT(jwtConfig).CreateToken(middleware.CustomClaims{
		ID: int(loginUser.UId),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "Token生成失败",
			"data": err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "登录成功",
		"data": map[string]interface{}{
			"token":      token,
			"user_id":    loginUser.UId,
			"expires_at": expirationTime,
		},
	})
}

func RealName(ctx context.Context, c *app.RequestContext) {
	var req request.RealNameReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}
	if !utils.IdCardVerification(req.IdCard) {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "身份证格式不正确",
			"data": nil,
		})
		return
	}

	// 调用用户服务进行实名认证
	name, err := UserClient.RealName(ctx, &user.RealNameReq{
		UId:      int64(c.GetInt("userId")),
		UserName: req.UserName,
		Sex:      req.Sex,
		Age:      req.Age,
		IdCard:   req.IdCard,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	// 根据业务响应动态返回HTTP状态码
	httpCode := 200
	if name.Code != 200 {
		httpCode = 400
		if name.Code >= 500 {
			httpCode = 500
		}
	}
	c.JSON(httpCode, name)
}

func TakeACar(ctx context.Context, c *app.RequestContext) {
	var req request.TakeACarReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	// 动态验证请求参数
	if req.StartLocation == "" || req.EndLocation == "" {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "起始位置和目的地不能为空",
			"data": nil,
		})
		return
	}

	car, err := UserClient.TakeACar(ctx, &user.TakeACarReq{
		UId:           int64(c.GetInt("userId")),
		StartLocation: req.StartLocation,
		EndLocation:   req.EndLocation,
		CartType:      req.CartType,
	})
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "服务调用失败",
			"data": err.Error(),
		})
		return
	}

	// 根据业务响应动态返回HTTP状态码
	httpCode := 200
	if car.Code != 200 {
		if car.Code >= 500 {
			httpCode = 500
		} else if car.Code >= 400 {
			httpCode = 400
		}
	}
	c.JSON(httpCode, car)
}
