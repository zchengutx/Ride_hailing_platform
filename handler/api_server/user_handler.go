package api_server

import (
	"Ride_hailing_platform/basic/config"
	"Ride_hailing_platform/kitex_gen/passengers"
	userserver "Ride_hailing_platform/kitex_gen/passengers/passengersserver"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// 用户Handler
type UserHandler struct {
	db            *gorm.DB
	redis         *redis.Client
	config        *config.Viper
	rpcClient     userserver.Client
	wechatHandler *WechatAuthHandler
}

func NewUserHandler(db *gorm.DB, redis *redis.Client, config *config.Viper, rpcClient userserver.Client) *UserHandler {
	return &UserHandler{
		db:            db,
		redis:         redis,
		config:        config,
		rpcClient:     rpcClient,
		wechatHandler: NewWechatAuthHandler(db, redis, config),
	}
}

// 发送短信验证码
func (h *UserHandler) SendSms(ctx context.Context, c *app.RequestContext) {
	// 获取手机号
	mobile := c.PostForm("mobile")
	source := c.PostForm("source")

	if mobile == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "手机号不能为空",
		})
		return
	}

	// 调用RPC服务
	req := &passengers.SendSmsReq{
		Mobile: mobile,
		Source: source,
	}

	resp, err := h.rpcClient.SendSms(ctx, req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		hlog.Errorf("发送短信验证码失败: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "发送短信验证码失败",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    resp.Code,
		"message": resp.Message,
	})
}

// 用户注册
func (h *UserHandler) Register(ctx context.Context, c *app.RequestContext) {
	// 获取注册参数
	mobile := c.PostForm("mobile")
	smsCode := c.PostForm("sms_code")

	if mobile == "" || smsCode == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "手机号和验证码不能为空",
		})
		return
	}

	// 构造RPC请求
	req := &passengers.RegisterReq{
		Mobile:  mobile,
		SmsCode: smsCode,
	}

	// 调用RPC服务
	resp, err := h.rpcClient.Register(ctx, req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		hlog.Errorf("用户注册失败: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "用户注册失败",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    resp.Code,
		"message": resp.Message,
	})
}

// 用户登录
func (h *UserHandler) Login(ctx context.Context, c *app.RequestContext) {
	// 获取登录参数
	mobile := c.PostForm("mobile")
	smsCode := c.PostForm("sms_code")

	if mobile == "" || smsCode == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "手机号和验证码不能为空",
		})
		return
	}

	// 构造RPC请求
	req := &passengers.LoginReq{
		Mobile:  mobile,
		SmsCode: smsCode,
	}

	// 调用RPC服务
	resp, err := h.rpcClient.Login(ctx, req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		hlog.Errorf("用户登录失败: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "用户登录失败",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    resp.Code,
		"message": resp.Message,
		"data": map[string]interface{}{
			"token": resp.Token,
		},
	})
}

// 微信登录
func (h *UserHandler) WechatLogin(ctx context.Context, c *app.RequestContext) {
	// 获取微信用户信息
	openID := c.PostForm("openid")
	nickname := c.PostForm("nickname")
	avatar := c.PostForm("avatar")

	if openID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "微信OpenID不能为空",
		})
		return
	}

	// 构造RPC请求
	req := &passengers.WechatLoginReq{
		Nickname: nickname,
		Avatar:   avatar,
	}

	// 调用RPC服务
	resp, err := h.rpcClient.WechatLogin(ctx, req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		hlog.Errorf("微信登录失败: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "微信登录失败",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    resp.Code,
		"message": resp.Message,
		"data": map[string]interface{}{
			"token": resp.Token,
		},
	})
}

// 绑定微信账号
func (h *UserHandler) BindWechat(ctx context.Context, c *app.RequestContext) {
	// 获取绑定参数
	userID := c.PostForm("user_id")
	openID := c.PostForm("openid")
	nickname := c.PostForm("nickname")
	avatar := c.PostForm("avatar")

	if userID == "" || openID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID和微信OpenID不能为空",
		})
		return
	}

	// 构造RPC请求
	req := &passengers.BindWechatReq{
		Nickname: nickname,
		Avatar:   avatar,
	}

	// 调用RPC服务
	resp, err := h.rpcClient.BindWechat(ctx, req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		hlog.Errorf("绑定微信账号失败: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "绑定微信账号失败",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    resp.Code,
		"message": resp.Message,
	})
}

// 解绑微信账号
func (h *UserHandler) UnbindWechat(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID
	userID := c.PostForm("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	// 构造RPC请求
	req := &passengers.UnbindWechatReq{}

	// 调用RPC服务
	resp, err := h.rpcClient.UnbindWechat(ctx, req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		hlog.Errorf("解绑微信账号失败: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "解绑微信账号失败",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    resp.Code,
		"message": resp.Message,
	})
}

// 获取用户信息
func (h *UserHandler) GetUserInfo(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID
	userID := c.Query("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	// 构造RPC请求
	req := &passengers.GetUserInfoReq{}

	// 调用RPC服务
	resp, err := h.rpcClient.GetUserInfo(ctx, req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		hlog.Errorf("获取用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取用户信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    resp.Code,
		"message": resp.Message,
	})
}

// 更新用户信息
func (h *UserHandler) UpdateUserInfo(ctx context.Context, c *app.RequestContext) {
	// 获取更新参数
	userID := c.PostForm("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	// 构造RPC请求
	req := &passengers.UpdateUserInfoReq{}

	// 调用RPC服务
	resp, err := h.rpcClient.UpdateUserInfo(ctx, req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		hlog.Errorf("更新用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "更新用户信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    resp.Code,
		"message": resp.Message,
	})
}

// 微信授权回调处理（集成到用户系统）
func (h *UserHandler) WechatAuthCallbackWithLogin(ctx context.Context, c *app.RequestContext) {
	// 获取回调参数
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "授权参数错误",
		})
		return
	}

	// 验证state
	stateKey := fmt.Sprintf("wechat_auth_state:%s", state)
	userID, err := h.redis.Get(ctx, stateKey).Result()
	if err != nil {
		hlog.Errorf("验证微信授权状态失败: %v", err)
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "授权状态无效或已过期",
		})
		return
	}

	// 删除已使用的state
	h.redis.Del(ctx, stateKey)

	// 获取access_token
	tokenResp, err := h.wechatHandler.getWechatAccessToken(code)
	if err != nil {
		hlog.Errorf("获取微信访问令牌失败: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取微信访问令牌失败",
		})
		return
	}

	if tokenResp.ErrorCode != 0 {
		hlog.Errorf("微信访问令牌错误: %s", tokenResp.ErrorMsg)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": fmt.Sprintf("微信授权失败: %s", tokenResp.ErrorMsg),
		})
		return
	}

	// 获取用户信息
	userInfo, err := h.wechatHandler.getWechatUserInfo(tokenResp.AccessToken, tokenResp.OpenID)
	if err != nil {
		hlog.Errorf("获取微信用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取微信用户信息失败",
		})
		return
	}

	var result map[string]interface{}

	// 如果有用户ID，则执行绑定操作
	if userID != "" {
		// 绑定微信账号
		bindReq := &passengers.BindWechatReq{
			Nickname: userInfo.Nickname,
			Avatar:   userInfo.HeadImgURL,
		}

		bindResp, err := h.rpcClient.BindWechat(ctx, bindReq, callopt.WithRPCTimeout(3*time.Second))
		if err != nil {
			hlog.Errorf("绑定微信账号失败: %v", err)
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "绑定微信账号失败",
			})
			return
		}

		result = map[string]interface{}{
			"code":    bindResp.Code,
			"message": "微信绑定成功",
			"data": map[string]interface{}{
				"type": "bind",
			},
		}
	} else {
		// 微信登录
		loginReq := &passengers.WechatLoginReq{
			Nickname: userInfo.Nickname,
			Avatar:   userInfo.HeadImgURL,
		}

		loginResp, err := h.rpcClient.WechatLogin(ctx, loginReq, callopt.WithRPCTimeout(3*time.Second))
		if err != nil {
			hlog.Errorf("微信登录失败: %v", err)
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "微信登录失败",
			})
			return
		}

		result = map[string]interface{}{
			"code":    loginResp.Code,
			"message": "微信登录成功",
			"data": map[string]interface{}{
				"token": loginResp.Token,
				"type":  "login",
			},
		}
	}

	c.JSON(http.StatusOK, result)
}
