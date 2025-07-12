package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// CORS 跨域中间件
func CORS() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		method := string(c.Method())

		// 设置CORS头
		c.Response.Header.Set("Access-Control-Allow-Origin", "*")
		c.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-token")
		c.Response.Header.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		c.Response.Header.Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if method == "OPTIONS" {
			c.Response.SetStatusCode(200)
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}
