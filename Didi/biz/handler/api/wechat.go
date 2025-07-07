package api

import (
	"Didi/biz/dal/global"
	"Didi/biz/handler/request"
	"Didi/utils"
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/skip2/go-qrcode"
)

func Sign(ctx context.Context, c *app.RequestContext) {
	var req request.SignReq
	if err := c.Bind(&req); err != nil {
		c.JSON(200, map[string]interface{}{
			"code": 400,
			"msg":  "查询失败",
			"data": err.Error(),
		})
		return
	}
	token := "555" //与后台token相同
	// 1. 将 token、timestamp、nonce 放入切片
	tmpArr := []string{token, req.Timestamp, req.Nonce}
	// 2. 对切片进行字典序排序
	sort.Strings(tmpArr)
	// 3. 拼接字符串
	tmpStr := strings.Join(tmpArr, "")
	// 4. 计算 SHA1 哈希值
	hash := sha1.New()
	hash.Write([]byte(tmpStr))
	computedSignature := fmt.Sprintf("%x", hash.Sum(nil))
	// 5. 比较签名
	if computedSignature == req.Signature {
		c.String(200, req.Echostr)
		return
	}
	return
}

func One(ctx context.Context, c *app.RequestContext) {
	redirect_uri := "http://11f623fe.r19.vip.cpolar.cn/v1/api/wechat/calBlack"                                                                                                                               //回调域名 (记得修改自己的)
	encodedRedirectURI := url.QueryEscape(redirect_uri)                                                                                                                                                      //对域名进行编码(防止特殊字符出现意外)
	urls := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=123#wechat_redirect", global.AppId, encodedRedirectURI) //生成url
	// 生成二维码
	qr, err := qrcode.New(urls, qrcode.Medium)
	fmt.Println(urls)
	if err != nil {
		return
	}
	// 获取PNG格式的二维码数据
	qrPNG, err := qr.PNG(256) // 256x256 像素
	if err != nil {
		return
	}
	//设置相应head头
	c.Header("Content-Type", "image/png")
	//输出二维码
	c.Data(200, "image/png", qrPNG)
}

func CalBlack(ctx context.Context, c *app.RequestContext) {
	//获取code
	var req request.CalblackReq
	if err := c.Bind(&req); err != nil {
		c.JSON(200, map[string]interface{}{
			"code": 400,
			"msg":  "查询失败",
			"data": err.Error(),
		})
		return
	}
	//获取token
	token := utils.GetAccessTokenByCode(req.Code)
	//获取用户信息
	resData, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%v&lang=zh_CN", token, "OPENID"))
	if err != nil {
		c.JSON(200, map[string]interface{}{
			"code": 400,
			"msg":  "查询失败",
			"data": err.Error(),
		})
		return
	}
	defer resData.Body.Close()
	//读取内容
	resDataBytes, err := ioutil.ReadAll(resData.Body)

	resDataMap := make(map[string]interface{})
	//解码
	err = json.Unmarshal(resDataBytes, &resDataMap)
	fmt.Println(resDataMap)
	//打印用户信息
	fmt.Println("openid:", resDataMap["openid"])
	fmt.Println("nickname:", resDataMap["nickname"])
}
