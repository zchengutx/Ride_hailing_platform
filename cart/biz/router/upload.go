package router

import (
	"cart/biz/dal/global"
	"cart/biz/handler/api"
	"cart/biz/middleware"
	"log"

	"github.com/cloudwego/hertz/pkg/route"
)

func UploadFileModel(r *route.RouterGroup) {
	log.Println("minio文件上传服务")
	uploadFileModel := r.Group("/upload")

	// 文件上传需要认证 - 防止匿名用户滥用存储资源
	uploadFileModel.Use(middleware.JWTAuth(global.JWT_SELECT_KEY))
	{
		log.Println("文件上传::::::/v1/upload/uploadFile")
		uploadFileModel.POST("/uploadFile", api.UploadFile)
	}
}
