package main

import (
	"cart/biz/utils"
	pb "cart/kitex_gen/cart/wechat"
	"cart/rpc/basic/global"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// WeChatServiceImpl implements the last service interface defined in the IDL.
type WeChatServiceImpl struct{}

// Sign 微信签名验证接口
// 验证微信服务器的签名，用于微信公众平台配置验证
func (s *WeChatServiceImpl) Sign(ctx context.Context, req *pb.SignReq) (resp *pb.SignResp, err error) {
	// 获取微信配置
	wechatConfig := &global.AppConf.WeChatConfig

	// 将token、timestamp、nonce三个参数进行字典序排序
	params := []string{wechatConfig.Token, req.Timestamp, req.Nonce}
	sort.Strings(params)

	// 将三个参数字符串拼接成一个字符串进行sha1加密
	str := strings.Join(params, "")
	hash := sha1.Sum([]byte(str))
	signature := hex.EncodeToString(hash[:])

	// 验证签名
	if signature == req.Signature {
		return &pb.SignResp{
			Code:    200,
			Message: "验证成功",
			Echostr: &req.Echostr,
		}, nil
	}

	return &pb.SignResp{
		Code:    400,
		Message: "签名验证失败",
	}, nil
}

// GetQRCode 获取授权二维码接口
// 生成微信授权二维码，用于用户扫码登录
func (s *WeChatServiceImpl) GetQRCode(ctx context.Context, req *pb.GetQRCodeReq) (resp *pb.GetQRCodeResp, err error) {
	// 获取微信配置
	wechatConfig := &global.AppConf.WeChatConfig

	// 构建微信授权URL
	params := url.Values{}
	params.Add("appid", wechatConfig.AppID)
	params.Add("redirect_uri", wechatConfig.RedirectURI)
	params.Add("response_type", "code")
	params.Add("scope", "snsapi_userinfo")
	params.Add("state", "cart_login")

	authURL := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?%s#wechat_redirect", params.Encode())

	// 这里简化处理，直接返回授权URL
	// 实际项目中可以生成二维码图片数据
	return &pb.GetQRCodeResp{
		Code:    200,
		Message: "获取成功",
		AuthUrl: &authURL,
	}, nil
}

// Callback 处理微信授权回调接口
// 处理微信授权回调，获取用户信息并保存到数据库
func (s *WeChatServiceImpl) Callback(ctx context.Context, req *pb.CallbackReq) (resp *pb.CallbackResp, err error) {
	// 获取微信配置
	wechatConfig := &global.AppConf.WeChatConfig

	// 1. 通过code获取access_token
	tokenURL := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		wechatConfig.AppID, wechatConfig.AppSecret, req.Code)

	// 创建HTTP客户端
	client := &http.Client{Timeout: 10 * time.Second}

	// 请求access_token
	tokenResp, err := client.Get(tokenURL)
	if err != nil {
		return &pb.CallbackResp{
			Code:    500,
			Message: "获取access_token失败",
		}, nil
	}
	defer tokenResp.Body.Close()

	tokenBody, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return &pb.CallbackResp{
			Code:    500,
			Message: "读取access_token响应失败",
		}, nil
	}

	var tokenData utils.WeChatAccessTokenResp
	if err := json.Unmarshal(tokenBody, &tokenData); err != nil {
		return &pb.CallbackResp{
			Code:    500,
			Message: "解析access_token响应失败",
		}, nil
	}

	// 检查是否有错误
	if tokenData.ErrCode != 0 {
		return &pb.CallbackResp{
			Code:    400,
			Message: fmt.Sprintf("获取access_token失败: %s", tokenData.ErrMsg),
		}, nil
	}

	// 2. 通过access_token获取用户信息
	userInfoURL := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		tokenData.AccessToken, tokenData.OpenID)

	userInfoResp, err := client.Get(userInfoURL)
	if err != nil {
		return &pb.CallbackResp{
			Code:    500,
			Message: "获取用户信息失败",
		}, nil
	}
	defer userInfoResp.Body.Close()

	userInfoBody, err := io.ReadAll(userInfoResp.Body)
	if err != nil {
		return &pb.CallbackResp{
			Code:    500,
			Message: "读取用户信息响应失败",
		}, nil
	}

	var userInfoData utils.WeChatUserInfoAPI
	if err := json.Unmarshal(userInfoBody, &userInfoData); err != nil {
		return &pb.CallbackResp{
			Code:    500,
			Message: "解析用户信息响应失败",
		}, nil
	}

	// 3. 保存微信用户信息到数据库
	if err := utils.SaveOrUpdateWeChatUser(&userInfoData, global.DB); err != nil {
		return &pb.CallbackResp{
			Code:    500,
			Message: "保存用户信息失败",
		}, nil
	}

	// 4. 保存微信令牌信息到数据库
	if err := utils.SaveOrUpdateWeChatToken(&tokenData, global.DB); err != nil {
		return &pb.CallbackResp{
			Code:    500,
			Message: "保存令牌信息失败",
		}, nil
	}

	// 5. 构建返回的用户信息
	userInfo := &pb.WeChatUserInfo{
		Openid:     userInfoData.OpenID,
		Nickname:   userInfoData.NickName,
		Headimgurl: userInfoData.HeadImgURL,
		Sex:        int16(userInfoData.Sex),
		Country:    userInfoData.Country,
		Province:   userInfoData.Province,
		City:       userInfoData.City,
		Language:   userInfoData.Language,
		Privilege:  userInfoData.Privilege,
	}

	// 返回成功响应
	return &pb.CallbackResp{
		Code:        200,
		Message:     "授权成功",
		UserInfo:    userInfo,
		AccessToken: &tokenData.AccessToken,
	}, nil
}
