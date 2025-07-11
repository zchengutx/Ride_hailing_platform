package router

import "github.com/cloudwego/hertz/pkg/app/server"

func Router(h *server.Hertz) {
	routerRouter := h.Group("/v1")
	{
		apiGroup := routerRouter.Group("/api")
		{
			PassengerModel(apiGroup)
			DriverModel(apiGroup)
			WeChatModel(apiGroup)
			MapModel(apiGroup)
			UploadFileModel(apiGroup)
			MongoModel(apiGroup)
		}
	}
}
