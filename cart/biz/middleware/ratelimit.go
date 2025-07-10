package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"golang.org/x/time/rate"
)

// RateLimiter 请求限流器
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	r        rate.Limit
	b        int
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

// getLimiter 获取IP对应的限流器
func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.limiters[ip] = limiter
	}

	return limiter
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(ip string) bool {
	return rl.getLimiter(ip).Allow()
}

// 全局限流器实例
var globalRateLimiter = NewRateLimiter(10, 20) // 每秒10个请求，突发20个

// RateLimit 限流中间件
func RateLimit() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 获取客户端IP
		ip := c.ClientIP()

		// 检查是否允许请求
		if !globalRateLimiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, map[string]interface{}{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
				"data": nil,
			})
			c.Abort()
			return
		}

		c.Next(ctx)
	}
}

// CleanupRateLimiters 定期清理过期的限流器
func (rl *RateLimiter) CleanupRateLimiters() {
	for {
		time.Sleep(time.Hour) // 每小时清理一次
		rl.mu.Lock()
		for ip, limiter := range rl.limiters {
			// 如果限流器很久没有使用，删除它
			if limiter.Tokens() == float64(rl.b) {
				delete(rl.limiters, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// InitRateLimit 初始化限流器
func InitRateLimit() {
	go globalRateLimiter.CleanupRateLimiters()
}
