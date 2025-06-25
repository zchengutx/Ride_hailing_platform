package api_server

import (
	"Ride_hailing_platform/basic/config"
	"Ride_hailing_platform/handler/model"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// 微信授权相关常量
const (
	WechatAuthURL     = "https://open.weixin.qq.com/connect/qrconnect"
	WechatTokenURL    = "https://api.weixin.qq.com/sns/oauth2/access_token"
	WechatUserInfoURL = "https://api.weixin.qq.com/sns/userinfo"
	StateExpireTime   = 10 * time.Minute
)

// 微信授权handler
type WechatAuthHandler struct {
	db     *gorm.DB
	redis  *redis.Client
	config *config.Viper
}

func NewWechatAuthHandler(db *gorm.DB, redis *redis.Client, config *config.Viper) *WechatAuthHandler {
	return &WechatAuthHandler{
		db:     db,
		redis:  redis,
		config: config,
	}
}

// 生成随机状态码
func generateState() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// 获取微信授权URL
func (h *WechatAuthHandler) GetWechatAuthURL(ctx context.Context, c *app.RequestContext) {
	// 获取用户ID（从token或session中获取）
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	// 生成state并存储到Redis
	state := generateState()
	stateKey := fmt.Sprintf("wechat_auth_state:%s", state)

	err := h.redis.Set(ctx, stateKey, userID, StateExpireTime).Err()
	if err != nil {
		hlog.Errorf("存储微信授权状态失败: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "生成授权链接失败",
		})
		return
	}

	// 构建微信授权URL
	authURL := fmt.Sprintf("%s?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect",
		WechatAuthURL,
		h.config.WechatConfig.AppID,
		url.QueryEscape(h.config.WechatConfig.RedirectURL),
		state,
	)

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"message":  "获取授权链接成功",
		"auth_url": authURL,
		"state":    state,
	})
}

// 微信授权回调处理
func (h *WechatAuthHandler) WechatAuthCallback(ctx context.Context, c *app.RequestContext) {
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
	tokenResp, err := h.getWechatAccessToken(code)
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
	userInfo, err := h.getWechatUserInfo(tokenResp.AccessToken, tokenResp.OpenID)
	if err != nil {
		hlog.Errorf("获取微信用户信息失败: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取微信用户信息失败",
		})
		return
	}

	// 绑定微信账号
	err = h.bindWechatAccount(userID, userInfo)
	if err != nil {
		hlog.Errorf("绑定微信账号失败: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "绑定微信账号失败",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "微信绑定成功",
		"data": map[string]interface{}{
			"openid":   userInfo.OpenID,
			"unionid":  userInfo.UnionID,
			"nickname": userInfo.Nickname,
			"avatar":   userInfo.HeadImgURL,
		},
	})
}

// 获取微信访问令牌
func (h *WechatAuthHandler) getWechatAccessToken(code string) (*model.WechatAccessTokenResponse, error) {
	tokenURL := fmt.Sprintf("%s?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		WechatTokenURL,
		h.config.WechatConfig.AppID,
		h.config.WechatConfig.AppSecret,
		code,
	)

	resp, err := http.Get(tokenURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResp model.WechatAccessTokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// 获取微信用户信息
func (h *WechatAuthHandler) getWechatUserInfo(accessToken, openID string) (*model.WechatUserInfo, error) {
	userInfoURL := fmt.Sprintf("%s?access_token=%s&openid=%s&lang=zh_CN",
		WechatUserInfoURL,
		accessToken,
		openID,
	)

	resp, err := http.Get(userInfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo model.WechatUserInfo
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// 绑定微信账号
func (h *WechatAuthHandler) bindWechatAccount(userID string, wechatInfo *model.WechatUserInfo) error {
	now := time.Now()

	// 更新用户微信信息
	updates := map[string]interface{}{
		"wechat_openid":   wechatInfo.OpenID,
		"wechat_unionid":  wechatInfo.UnionID,
		"wechat_nickname": wechatInfo.Nickname,
		"wechat_avatar":   wechatInfo.HeadImgURL,
		"wechat_bound_at": &now,
		"update_time":     now,
	}

	err := h.db.Model(&model.LxhPassenger{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return fmt.Errorf("更新用户微信信息失败: %v", err)
	}

	return nil
}

// 解绑微信账号
func (h *WechatAuthHandler) UnbindWechatAccount(ctx context.Context, c *app.RequestContext) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	// 清除微信绑定信息
	updates := map[string]interface{}{
		"wechat_openid":   "",
		"wechat_unionid":  "",
		"wechat_nickname": "",
		"wechat_avatar":   "",
		"wechat_bound_at": nil,
		"update_time":     time.Now(),
	}

	err := h.db.Model(&model.LxhPassenger{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		hlog.Errorf("解绑微信账号失败: %v", err)
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "解绑微信账号失败",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "微信解绑成功",
	})
}

// 获取微信绑定状态
func (h *WechatAuthHandler) GetWechatBindStatus(ctx context.Context, c *app.RequestContext) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "用户ID不能为空",
		})
		return
	}

	var user model.LxhPassenger
	err := h.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"code":    404,
				"message": "用户不存在",
			})
		} else {
			hlog.Errorf("查询用户信息失败: %v", err)
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code":    500,
				"message": "查询用户信息失败",
			})
		}
		return
	}

	// 判断是否已绑定微信
	isBound := user.WechatOpenID != ""

	data := map[string]interface{}{
		"is_bound": isBound,
	}

	if isBound {
		data["wechat_info"] = map[string]interface{}{
			"openid":   user.WechatOpenID,
			"unionid":  user.WechatUnionID,
			"nickname": user.WechatNickname,
			"avatar":   user.WechatAvatar,
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "查询成功",
		"data":    data,
	})
}
