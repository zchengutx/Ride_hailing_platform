package utils

import (
	"cart/biz/dal/global"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gorm.io/gorm"

	"cart/rpc/basic/model"
)

func GetAccessTokenByCode(code string) string {
	var resD map[string]interface{}
	//accesstoken（获取用户信息）
	res, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", global.AppId, global.Secret, code))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()
	resData, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(resData, &resD)
	fmt.Println(resD)
	fmt.Println(resD["access_token"])
	return resD["access_token"].(string)
}

// WeChatUserInfoAPI 微信用户信息结构体（用于API调用）
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

// WeChatAccessTokenResp 微信AccessToken响应结构体
type WeChatAccessTokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}

// SaveOrUpdateWeChatUser 保存或更新微信用户信息
func SaveOrUpdateWeChatUser(userInfoData *WeChatUserInfoAPI, db *gorm.DB) error {
	privilegeJSON, _ := json.Marshal(userInfoData.Privilege)
	var existingUser model.LxhWechatUser
	result := db.Where("openid = ?", userInfoData.OpenID).First(&existingUser)
	if result.Error != nil {
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
		return db.Create(&newUser).Error
	} else {
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
		return db.Model(&existingUser).Updates(updates).Error
	}
}

// SaveOrUpdateWeChatToken 保存或更新微信令牌信息
func SaveOrUpdateWeChatToken(tokenData *WeChatAccessTokenResp, db *gorm.DB) error {
	expiresAt := time.Now().Add(time.Duration(tokenData.ExpiresIn) * time.Second)
	var existingToken model.LxhWechatToken
	result := db.Where("openid = ?", tokenData.OpenID).First(&existingToken)
	if result.Error != nil {
		newToken := model.LxhWechatToken{
			Openid:       tokenData.OpenID,
			AccessToken:  tokenData.AccessToken,
			RefreshToken: tokenData.RefreshToken,
			ExpiresIn:    tokenData.ExpiresIn,
			Scope:        tokenData.Scope,
			TokenType:    "Bearer",
			ExpiresAt:    expiresAt,
		}
		return db.Create(&newToken).Error
	} else {
		updates := map[string]interface{}{
			"access_token":  tokenData.AccessToken,
			"refresh_token": tokenData.RefreshToken,
			"expires_in":    tokenData.ExpiresIn,
			"scope":         tokenData.Scope,
			"expires_at":    expiresAt,
		}
		return db.Model(&existingToken).Updates(updates).Error
	}
}
