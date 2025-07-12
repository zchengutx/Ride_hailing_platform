package pkg

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"kitex_main/gateway/biz/handler/response"
	_ "kitex_main/init"

	"mime"
	"net/http"
	"path/filepath"
	"time"
)

const (
	MinioEndpoint  = "14.103.134.228:9000"
	MinioAccessKey = "EsC0ClOoyswxNguwC6JP"
	MinioSecretKey = "6OOrbmDCIh1lZ0urTdLtkz4iWR8EPauw1Hyxtfzj"
	BucketName     = "chen"
)

func Upload(ctx context.Context, c *app.RequestContext) {

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "file not exists" + err.Error(),
		})
		return
	}

	client, err := minio.New(MinioEndpoint, &minio.Options{
		Creds: credentials.NewStaticV4(MinioAccessKey, MinioSecretKey, ""),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Code:    400,
			Message: "file not exists" + err.Error(),
		})
		return
	}

	open, err := file.Open()
	defer open.Close()

	location, _ := time.LoadLocation("Asia/Shanghai")
	Now := time.Now().In(location).Truncate(time.Second)
	format := Now.Format("20060102150405")

	ext := filepath.Ext(file.Filename)
	extension := mime.TypeByExtension(ext)

	sprintf := fmt.Sprintf("%s/%s", format, file.Filename)
	client.PutObject(ctx, BucketName, sprintf, open, file.Size, minio.PutObjectOptions{ContentType: mime.TypeByExtension(extension)})

	Url := fmt.Sprintf("http://%s/%s/%s", MinioEndpoint, BucketName, sprintf)

	c.JSON(http.StatusOK, response.Response{
		Code:    200,
		Message: "file uploaded",
		Data:    Url,
	})

}
