package main

import (
	pb "cart/kitex_gen/cart/wechat"
	"cart/rpc/basic/global"
	"cart/rpc/basic/model"
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

// Sign 微信签名验证接口，验证微信服务器推送的消息签名
func (s *WeChatServiceImpl) Sign(ctx context.Context, req *pb.SignReq) (resp *pb.SignResp, err error) {
	ctx = context.Background()

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
			Message: "签名验证成功",
			Echostr: &req.Echostr,
		}, nil
	}

	return &pb.SignResp{
		Code:    400,
		Message: "签名验证失败",
	}, nil
}

// GetQRCode 获取授权二维码接口，生成微信扫码登录二维码
func (s *WeChatServiceImpl) GetQRCode(ctx context.Context, req *pb.GetQRCodeReq) (resp *pb.GetQRCodeResp, err error) {
	ctx = context.Background()

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
		Message: "获取授权二维码成功",
		AuthUrl: &authURL,
	}, nil
}

// WeChatUserInfoAPI 微信用户信息结构（用于API调用）
type WeChatUserInfoAPI struct {
	OpenID     string   `json:"openid"`
	NickName   string   `json:"nickname"`
	HeadImgURL string   `json:"headimgurl"`
	Sex        int      `json:"sex"`
	Country    string   `json:"country"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Language   string   `json:"language"`
	Privilege  []string `json:"privilege"`
	UnionID    string   `json:"unionid"`
}

// WeChatAccessTokenResp 微信AccessToken响应
type WeChatAccessTokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}

// saveOrUpdateWeChatUser 保存或更新微信用户信息
func (s *WeChatServiceImpl) saveOrUpdateWeChatUser(userInfoData *WeChatUserInfoAPI) error {
	// 将特权信息转换为JSON字符串
	privilegeJSON, _ := json.Marshal(userInfoData.Privilege)

	// 查询是否已存在该用户
	var existingUser model.LxhWechatUser
	result := global.DB.Debug().Where("openid = ?", userInfoData.OpenID).First(&existingUser)

	if result.Error != nil {
		// 用户不存在，创建新用户
		newUser := model.LxhWechatUser{
			Openid:     userInfoData.OpenID,
			Nickname:   userInfoData.NickName,
			Headimgurl: userInfoData.HeadImgURL,
			Sex:        userInfoData.Sex,
			Country:    userInfoData.Country,
			Province:   userInfoData.Province,
			City:       userInfoData.City,
			Language:   userInfoData.Language,
			Privilege:  string(privilegeJSON),
			UnionId:    userInfoData.UnionID,
		}
		return global.DB.Create(&newUser).Error
	} else {
		// 用户已存在，更新用户信息
		updates := map[string]interface{}{
			"nickname":   userInfoData.NickName,
			"headimgurl": userInfoData.HeadImgURL,
			"sex":        userInfoData.Sex,
			"country":    userInfoData.Country,
			"province":   userInfoData.Province,
			"city":       userInfoData.City,
			"language":   userInfoData.Language,
			"privilege":  string(privilegeJSON),
		}
		if userInfoData.UnionID != "" {
			updates["unionid"] = userInfoData.UnionID
		}
		return global.DB.Debug().Model(&existingUser).Updates(updates).Error
	}
}

// saveOrUpdateWeChatToken 保存或更新微信令牌信息
func (s *WeChatServiceImpl) saveOrUpdateWeChatToken(tokenData *WeChatAccessTokenResp) error {
	// 计算令牌过期时间
	expiresAt := time.Now().Add(time.Duration(tokenData.ExpiresIn) * time.Second)

	// 查询是否已存在该用户的令牌
	var existingToken model.LxhWechatToken
	result := global.DB.Debug().Where("openid = ?", tokenData.OpenID).First(&existingToken)

	if result.Error != nil {
		// 令牌不存在，创建新令牌记录
		newToken := model.LxhWechatToken{
			Openid:       tokenData.OpenID,
			AccessToken:  tokenData.AccessToken,
			RefreshToken: tokenData.RefreshToken,
			ExpiresIn:    tokenData.ExpiresIn,
			Scope:        tokenData.Scope,
			TokenType:    "Bearer",
			ExpiresAt:    expiresAt,
		}
		return global.DB.Create(&newToken).Error
	} else {
		// 令牌已存在，更新令牌信息
		updates := map[string]interface{}{
			"access_token":  tokenData.AccessToken,
			"refresh_token": tokenData.RefreshToken,
			"expires_in":    tokenData.ExpiresIn,
			"scope":         tokenData.Scope,
			"expires_at":    expiresAt,
		}
		return global.DB.Debug().Model(&existingToken).Updates(updates).Error
	}
}

// Callback 处理微信授权回调接口，处理微信授权成功后的回调
func (s *WeChatServiceImpl) Callback(ctx context.Context, req *pb.CallbackReq) (resp *pb.CallbackResp, err error) {
	ctx = context.Background()

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

	var tokenData WeChatAccessTokenResp
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

	var userInfoData WeChatUserInfoAPI
	if err := json.Unmarshal(userInfoBody, &userInfoData); err != nil {
		return &pb.CallbackResp{
			Code:    500,
			Message: "解析用户信息响应失败",
		}, nil
	}

	// 3. 保存微信用户信息到数据库
	if err := s.saveOrUpdateWeChatUser(&userInfoData); err != nil {
		return &pb.CallbackResp{
			Code:    500,
			Message: "保存用户信息失败",
		}, nil
	}

	// 4. 保存微信令牌信息到数据库
	if err := s.saveOrUpdateWeChatToken(&tokenData); err != nil {
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
		Message:     "微信授权成功",
		UserInfo:    userInfo,
		AccessToken: &tokenData.AccessToken,
	}, nil
}
