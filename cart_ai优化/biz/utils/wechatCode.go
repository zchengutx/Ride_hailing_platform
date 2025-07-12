package utils

import (
	"cart/biz/dal/global"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
