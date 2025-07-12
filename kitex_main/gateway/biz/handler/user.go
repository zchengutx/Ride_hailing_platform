package handler

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	"github.com/skip2/go-qrcode"
	"kitex_main/config"
	"kitex_main/gateway/biz/handler/request"
	"kitex_main/gateway/biz/handler/response"
	users "kitex_main/kitex_gen/Ride_hailing_platform/user"
	"kitex_main/kitex_gen/Ride_hailing_platform/user/userserver"
	"kitex_main/pkg"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

// 短信验证码
func SendSms(ctx context.Context, c *app.RequestContext) {
	cil, err := userserver.NewClient("Ride_hailing_platform.user", client.WithHostPorts("0.0.0.0:50051"))
	if err != nil {
		log.Fatal(err)
	}

	var req request.SendSms
	err = c.BindAndValidate(&req)
	if err != nil {
		res := response.Response{
			Code:    400,
			Message: err.Error(),
			Data:    nil,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if !pkg.CheckMobile(req.Mobile) {
		res := response.Response{
			Code:    400,
			Message: "mobile not match",
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	sms, _ := cil.SendSms(ctx, &users.SendSmsReq{
		Mobile: req.Mobile,
		Source: req.Source,
	})

	if sms.Code != 200 {
		res := response.Response{
			Code:    int(sms.Code),
			Message: sms.Msg,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	res := response.Response{
		Code:    200,
		Message: "send sms success",
	}
	c.JSON(http.StatusOK, res)

}

// 手机号登录接口
func Login(ctx context.Context, c *app.RequestContext) {
	cil, err := userserver.NewClient("Ride_hailing_platform.user", client.WithHostPorts("0.0.0.0:50051"))
	if err != nil {
		log.Fatal(err)
	}

	var req request.Login
	err = c.BindAndValidate(&req)
	if err != nil {
		res := response.Response{
			Code:    400,
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if !pkg.CheckMobile(req.Mobile) {
		res := response.Response{
			Code:    400,
			Message: "mobile not match",
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	login, _ := cil.Login(ctx, &users.LoginReq{
		Mobile:      req.Mobile,
		SendSmsCode: req.SendSmsCode,
	})
	if login.Code != 200 {
		res := response.Response{
			Code:    400,
			Message: login.Msg,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	handler, err := pkg.TokenHandler(strconv.FormatInt(login.Token, 10))
	if err != nil {
		res := response.Response{
			Code:    400,
			Message: "create token error",
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	res := response.Response{
		Code:    200,
		Message: "login success",
		Data:    handler,
	}
	c.JSON(http.StatusOK, res)
}

func Wechat(ctx context.Context, c *app.RequestContext) {

	// 获取查询参数中的签名、时间戳和随机数
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")
	// 创建包含令牌、时间戳和随机数的字符串切片
	tmpArr := []string{config.WECHAT_TOKEN, timestamp, nonce}
	// 对切片进行字典排序
	sort.Strings(tmpArr)
	// 将排序后的元素拼接成单个字符串
	tmpStr := ""
	for _, v := range tmpArr {
		tmpStr += v
	}
	// 对字符串进行SHA-1哈希计算
	tmpHash := sha1.New()
	tmpHash.Write([]byte(tmpStr))
	tmpStr = fmt.Sprintf("%x", tmpHash.Sum(nil))
	fmt.Println(tmpStr)
	fmt.Println(signature)
	// 将计算得到的签名与请求中提供的签名进行比较，并根据结果发送相应的响应
	if tmpStr == signature {
		c.String(200, echostr)
		return
	} else {
		res := response.Response{
			Code:    403,
			Message: "签名验证失败 " + timestamp,
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}
}

// Redirect 微信扫码登录
// @Summary 用户登录接口3
// @Description 通过微信扫码登录，手机进行登录验证
// @Tags 公开
// @Accept json
// @Produce application/json
// @Param Url query string true "内网穿透地址"
// @Router /api/v1/wechat/login [get]
func WechatLogin(ctx context.Context, c *app.RequestContext) {

	state := pkg.RandomString(5)                                                           //防止跨站请求伪造攻击 增加安全性
	redirectURL := url.QueryEscape("http://" + "4efa1496.r8.cpolar.top" + "/api/callback") //userinfo,
	wechatLoginURL := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&state=%s&scope=snsapi_userinfo#wechat_redirect", "wx41cc84964c03b1a5", redirectURL, state)
	wechatLoginURL, _ = url.QueryUnescape(wechatLoginURL)
	// 生成二维码
	qrCode, err := qrcode.Encode(wechatLoginURL, qrcode.Medium, 256)
	if err != nil {
		// 错误处理
		res := response.Response{
			Code:    400,
			Message: "Error generating QR code",
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	// 将二维码图片作为响应返回给用户
	c.Header("Content-Type", "image/png")
	r := response.Response{
		Code:    http.StatusOK,
		Message: "wechat_redirect success",
		Data:    qrCode,
	}
	c.JSON(http.StatusOK, r)

}

func CallBack(ctx context.Context, c *app.RequestContext) {
	// 获取微信返回的授权码
	code := c.Query("code")
	// 向微信服务器发送请求，获取access_token和openid
	tokenResp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", "wx41cc84964c03b1a5", "d5162bc4a9f53cb277461b5f1bc06b50", code))
	if err != nil {
		fmt.Println(err)
		resp := &response.Response{
			Data:    nil,
			Message: "error,获取token失败",
			Code:    400,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	// 解析响应中的access_token和openid
	var tokenData struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenID       string `json:"openid"`
		Scope        string `json:"scope"`
	}
	if err1 := json.NewDecoder(tokenResp.Body).Decode(&tokenData); err1 != nil {
		resp := &response.Response{
			Data:    nil,
			Message: "error,获取token失败",
			Code:    400,
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	userInfoURL := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", tokenData.AccessToken, tokenData.OpenID)
	userInfoResp, err := http.Get(userInfoURL)
	if err != nil {
		// 错误处理
		log.Fatal("获取失败")
		return
	}
	defer userInfoResp.Body.Close()

	//------------------------------------
	var userData struct {
		OpenID   string `json:"openid"`
		Nickname string `json:"nickname"`
	}
	if err1 := json.NewDecoder(userInfoResp.Body).Decode(&userData); err1 != nil {
		// 错误处理
		log.Fatal("获取用户信息失败")
		return
	}

	cil, err := userserver.NewClient("Ride_hailing_platform.user", client.WithHostPorts("0.0.0.0:50051"))
	if err != nil {
		log.Fatal(err)
	}

	back, _ := cil.CallBack(ctx, &users.CallBackReq{
		OpenId:   userData.OpenID,
		NickName: userData.Nickname,
	})
	if back.Code != 200 {
		resp := &response.Response{
			Code:    400,
			Message: "callBack error",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	handler, err := pkg.TokenHandler(strconv.FormatInt(back.Token, 10))
	if err != nil {
		resp := &response.Response{
			Code:    400,
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	redirectUrl := "http://67594adb.r19.cpolar.top/login-callback?token=" + handler
	c.Redirect(consts.StatusFound, []byte(redirectUrl))
}

func InfoUser(ctx context.Context, c *app.RequestContext) {
	var userId int
	if value, exists := c.Get("userID"); exists {
		userId = value.(int)
	} else {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "userid is error",
		})
		return
	}

	cil, err := userserver.NewClient("Ride_hailing_platform.user", client.WithHostPorts("0.0.0.0:50051"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "init user client error",
		})
		return
	}

	user, _ := cil.InfoUser(ctx, &users.InfoUserReq{UserId: int64(userId)})
	if user.Code != 200 {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "user info error",
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "info success",
		Data:    user.Info,
	})

}
