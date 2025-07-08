package utils

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"mime/multipart"
)

/*
file, _ := c.FormFile("file")

	if file.Size >= 500*1024*1024 {
		c.JSON(200, gin.H{
			"code": 10001,
			"msg":  "The file is too large, and only files within 500MB are allowed to be uploaded.",
			"data": nil,
		})
		return
	}

ext := filepath.Ext(file.Filename)

	if ext != ".mp4" {
		c.JSON(200, gin.H{
			"code": 10002,
			"msg":  "The file only supports mp4 format.",
			"data": nil,
		})
		return
	}

fileName := time.Now().Format("20060102150405") + ext
open, err := file.Open()

	if err != nil {
		log.Println("file open error")
		return
	}

defer open.Close()

	if !utils.MiNio(fileName, open) {
		c.JSON(200, gin.H{
			"code": 10003,
			"msg":  "File upload failed.",
			"data": nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "File upload successfully.",
		"data": nil,
	})
*/
func UploadMinio(name string, open multipart.File) bool {
	// MinIO 服务器配置

	endpoint := ""
	accessKeyID := ""
	secretAccessKey := ""

	// 文件配置
	fileName := name           // 要上传的文件名
	bucketName := ""           // 目标桶名
	objectName := fileName     // 上传后对象的名称（可包含路径前缀）
	contentType := "video/mp4" // 文件的 MIME 类型

	// 初始化 MinIO 客户端
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// 创建上下文
	ctx := context.Background()

	// 创建桶（如果不存在）
	location := "test"
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists != nil {
			log.Fatalln(errBucketExists)
		}
		if exists {
			log.Printf("Bucket %s already exists", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created bucket %s", bucketName)
	}

	// 上传文件
	//filePath := filepath.Join("./tmp/", fileName) // 文件路径
	//fmt.Println(filePath)
	_, err = minioClient.PutObject(ctx, bucketName, objectName, open, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Printf("Failed to upload %s: %v", objectName, err)
		return false
	}

	log.Printf("Successfully uploaded %s%s", "url", fileName)
	return true
}
