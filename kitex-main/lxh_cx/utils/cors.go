package utils

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// CorsMiddleware 跨域中间件
func CorsMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 设置CORS头
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, x-token")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if string(c.Method()) == "OPTIONS" {
			c.Status(200)
			return
		}

		c.Next(ctx)
	}
}
