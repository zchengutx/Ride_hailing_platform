package middleware

import (
	"cart/biz/utils"
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/sirupsen/logrus"
)

// generateRequestID 生成请求ID
func generateRequestID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// LoggingMiddleware 请求日志中间件
func LoggingMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		requestID := generateRequestID()

		// 设置请求ID到上下文
		c.Set("request_id", requestID)

		// 创建请求日志记录器
		requestLogger := utils.NewRequestLogger(
			requestID,
			string(c.Method()),
			string(c.Path()),
			c.ClientIP(),
		)

		// 记录请求开始
		requestLogger.WithFields(logrus.Fields{
			"user_agent": string(c.UserAgent()),
			"referer":    string(c.Request.Header.Peek("Referer")),
		}).Info("Request started")

		// 继续处理请求
		c.Next(ctx)

		// 计算响应时间
		latency := time.Since(start)

		// 获取响应信息
		statusCode := c.Response.StatusCode()
		responseSize := int64(c.Response.Header.ContentLength())

		// 记录请求完成
		requestLogger.LogRequest(statusCode, latency, responseSize)

		// 记录慢请求
		if latency > 1*time.Second {
			utils.Warn("Slow request detected", logrus.Fields{
				"request_id": requestID,
				"path":       string(c.Path()),
				"latency":    latency.String(),
			})
		}

		// 记录错误响应的详细信息
		if statusCode >= 400 {
			errorDetails := map[string]interface{}{
				"request_id":  requestID,
				"method":      string(c.Method()),
				"path":        string(c.Path()),
				"status_code": statusCode,
				"client_ip":   c.ClientIP(),
				"user_agent":  string(c.UserAgent()),
				"latency":     latency.String(),
			}

			if statusCode >= 500 {
				utils.Error("Server error response", fmt.Errorf("HTTP %d", statusCode), errorDetails)
			} else {
				utils.Warn("Client error response", logrus.Fields(errorDetails))
			}
		}
	}
}

// RecoveryMiddleware 异常恢复中间件
func RecoveryMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get("request_id")

				errorDetails := map[string]interface{}{
					"request_id": requestID,
					"method":     string(c.Method()),
					"path":       string(c.Path()),
					"client_ip":  c.ClientIP(),
					"panic":      err,
				}

				utils.Error("Panic recovered", fmt.Errorf("panic: %v", err), errorDetails)

				// 返回500错误
				c.JSON(500, map[string]interface{}{
					"code": 500,
					"msg":  "服务器内部错误",
					"data": nil,
				})
				c.Abort()
			}
		}()

		c.Next(ctx)
	}
}
