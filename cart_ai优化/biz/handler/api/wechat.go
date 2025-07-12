// api 包下的 wechat.go 实现微信相关接口，如签名校验、二维码生成、回调处理等
package api

import (
	"cart/biz/handler/request"
	"cart/biz/utils"
	pb "cart/kitex_gen/cart/wechat"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/skip2/go-qrcode"
)

// WeChatClient 是与微信服务交互的RPC客户端
var (
	WeChatClient = utils.GetDefaultWeChatClient()
)

// Sign 微信签名校验接口，校验微信服务器推送的消息签名
func Sign(ctx context.Context, c *app.RequestContext) {
	var req request.SignReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	// 调用RPC服务验证签名
	response, _ := WeChatClient.Sign(ctx, &pb.SignReq{
		Signature: req.Signature,
		Timestamp: req.Timestamp,
		Nonce:     req.Nonce,
		Echostr:   req.Echostr,
	})

	if response.Code == 200 {
		// 签名验证成功，返回echostr
		if response.Echostr != nil {
			c.String(200, *response.Echostr)
		} else {
			c.String(200, "")
		}
		return
	}

	c.JSON(200, response)
}

// GetQRCode 生成微信授权二维码接口，生成微信扫码登录二维码
func GetQRCode(ctx context.Context, c *app.RequestContext) {
	// 调用RPC服务获取授权URL
	response, _ := WeChatClient.GetQRCode(ctx, &pb.GetQRCodeReq{})

	if response.Code != 200 {
		c.JSON(200, response)
		return
	}

	if response.AuthUrl == nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "获取授权URL失败",
			"data": nil,
		})
		return
	}

	// 生成二维码
	qr, err := qrcode.New(*response.AuthUrl, qrcode.Medium)
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "生成二维码失败",
			"data": err.Error(),
		})
		return
	}

	// 获取PNG格式的二维码数据
	qrPNG, err := qr.PNG(256) // 256x256 像素
	if err != nil {
		c.JSON(500, map[string]interface{}{
			"code": 500,
			"msg":  "生成二维码图片失败",
			"data": err.Error(),
		})
		return
	}

	// 设置响应头并输出二维码
	c.Header("Content-Type", "image/png")
	c.Data(200, "image/png", qrPNG)
}

// Callback 微信回调处理接口，处理微信授权成功后的回调
func Callback(ctx context.Context, c *app.RequestContext) {
	var req request.CalblackReq
	if err := c.Bind(&req); err != nil {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "参数绑定失败",
			"data": err.Error(),
		})
		return
	}

	// 验证code参数
	if req.Code == "" {
		c.JSON(400, map[string]interface{}{
			"code": 400,
			"msg":  "授权码不能为空",
			"data": nil,
		})
		return
	}

	// 调用RPC服务处理回调
	response, _ := WeChatClient.Callback(ctx, &pb.CallbackReq{
		Code: req.Code,
	})

	c.JSON(200, response)
}
