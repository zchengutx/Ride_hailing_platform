package main

import (
	"Ride_hailing_platform/basic/global"
	"Ride_hailing_platform/handler/api_server"
	u "Ride_hailing_platform/kitex_gen/passengers"
	user "Ride_hailing_platform/kitex_gen/passengers/passengersserver"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/hertz-contrib/cors"
	etcd "github.com/kitex-contrib/registry-etcd"

	"context"
	"time"

	_ "Ride_hailing_platform/basic/init"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client/callopt"
)

var (
	userRPCClient user.Client
)

func main() {
	Http()
	resolver, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		hlog.Fatalf("创建etcd resolver失败: %v", err)
	}

	userRPCClient, err = user.NewClient("passengers", client.WithResolver(resolver))
	if err != nil {
		hlog.Fatalf("创建用户服务客户端失败: %v", err)
	}

	hz := server.New(server.WithHostPorts("localhost:8889"))

	userHandler := api_server.NewUserHandler(global.DB, global.Redis, global.Conf, userRPCClient)

	wechatAuthHandler := api_server.NewWechatAuthHandler(global.DB, global.Redis, global.Conf)

	userGroup := hz.Group("/api/user")
	{
		userGroup.POST("/sms/send", userHandler.SendSms)           // 发送短信验证码
		userGroup.POST("/register", userHandler.Register)          // 用户注册
		userGroup.POST("/login", userHandler.Login)                // 用户登录
		userGroup.POST("/wechat/login", userHandler.WechatLogin)   // 微信登录
		userGroup.POST("/wechat/bind", userHandler.BindWechat)     // 绑定微信账号
		userGroup.GET("/info", userHandler.GetUserInfo)            // 获取用户信息
		userGroup.POST("/info/update", userHandler.UpdateUserInfo) // 更新用户信息
	}

	wechatGroup := hz.Group("/api/wechat")
	{
		wechatGroup.GET("/auth/url", wechatAuthHandler.GetWechatAuthURL)           // 获取微信授权URL
		wechatGroup.GET("/auth/callback", userHandler.WechatAuthCallbackWithLogin) // 微信授权回调（集成登录绑定）
	}

	// 兼容原有的API路由
	hz.GET("/api/item", OldHandler)

	// 静态文件服务（用于测试页面）
	hz.Static("/static", "../")
	hz.StaticFile("/index.html", "../index.html")
	hz.StaticFile("/example.html", "../example.html")

	hlog.Info("HTTP服务器启动中...")
	hlog.Info("用户服务API地址: http://localhost:8889/api/user/")
	hlog.Info("微信授权API地址: http://localhost:8889/api/wechat/")
	hlog.Info("测试页面地址: http://localhost:8889/static/example.html")

	if err := hz.Run(); err != nil {
		hlog.Fatalf("HTTP服务器启动失败: %v", err)
	}
}

// 原有的测试Handler
func OldHandler(ctx context.Context, c *app.RequestContext) {
	req := u.NewRegisterReq()
	req.Mobile = "1024"
	resp, err := userRPCClient.Register(context.Background(), req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		hlog.Errorf("调用注册服务失败: %v", err)
		c.String(500, "调用注册服务失败")
		return
	}

	c.String(200, resp.String())
}

func Http() {
	h := server.Default()

	// 配置 CORS 中间件
	h.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 或指定域名 ["https://example.com"]
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// 注册路由...
	h.Spin()
}
