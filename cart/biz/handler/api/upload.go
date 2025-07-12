// api 包下的 upload.go 实现文件上传相关接口
package api

import (
	"cart/biz/utils"
	"context"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

// UploadFile 文件上传接口，支持视频、音频、图片格式，最大500MB
// 1. 校验文件大小
// 2. 校验文件格式
// 3. 上传到 Minio 对象存储
func UploadFile(ctx context.Context, c *app.RequestContext) {
	file, _ := c.FormFile("file")

	// 校验文件大小，最大500MB
	if file.Size >= 500*1024*1024 {
		c.JSON(200, map[string]interface{}{
			"code": 10001,
			"msg":  "The file is too large, and only files within 500MB are allowed to be uploaded.",
			"data": nil,
		})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	// 支持的文件格式
	supportedFormats := []string{
		// 视频格式
		".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv", ".webm", ".m4v",
		// 音频格式
		".mp3", ".wav", ".aac", ".flac", ".m4a", ".ogg", ".wma",
		// 图片格式
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff",
	}

	// 检查文件格式是否支持
	supported := false
	for _, format := range supportedFormats {
		if ext == format {
			supported = true
			break
		}
	}

	if !supported {
		c.JSON(200, map[string]interface{}{
			"code": 10002,
			"msg":  "The file format is not supported. Only video, audio and image files are allowed.",
			"data": nil,
		})
		return
	}

	// 生成唯一文件名
	fileName := time.Now().Format("20060102150405") + ext
	open, err := file.Open()

	if err != nil {
		log.Println("file open error")
		return
	}

	defer open.Close()

	// 上传到 Minio
	if !utils.UploadMinio(fileName, open) {
		c.JSON(200, map[string]interface{}{
			"code": 10003,
			"msg":  "File upload failed.",
			"data": nil,
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"code": 200,
		"msg":  "File upload successfully.",
		"data": nil,
	})

}
