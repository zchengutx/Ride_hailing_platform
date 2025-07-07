package router

import "github.com/cloudwego/hertz/pkg/app/server"

func Register(h *server.Hertz) {
	routerGroup := h.Group("/v1")
	{
		apiGroup := routerGroup.Group("/api")
		{
			UserModel(apiGroup)
			WechatModel(apiGroup)
			DriverModel(apiGroup)
		}
	}
}
