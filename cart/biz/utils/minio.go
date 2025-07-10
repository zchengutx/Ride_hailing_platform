package utils

import (
	"cart/biz/dal/global"
	"context"
	"log"
	"mime"
	"mime/multipart"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	MinioConf = &global.AppConf.Minio
)

func UploadMinio(name string, open multipart.File) bool {
	// MinIO 服务器配置

	endpoint := MinioConf.Endpoint
	accessKeyID := MinioConf.AccessKeyId
	secretAccessKey := MinioConf.AccessKeySecret

	// 文件配置
	fileName := name                   // 要上传的文件名
	bucketName := MinioConf.BucketName // 目标桶名
	objectName := fileName             // 上传后对象的名称（可包含路径前缀）

	// 根据文件扩展名自动检测文件类型
	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	if contentType == "" {
		contentType = "application/octet-stream" // 默认为二进制文件类型，支持所有文件
	}

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
	location := MinioConf.BucketName
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

	log.Printf("Successfully uploaded %s%s", MinioConf.BucketUrl, fileName)
	return true
}
