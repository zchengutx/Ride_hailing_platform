// api 包下的 health.go 实现健康检查相关接口
package api

import (
	"cart/biz/utils"
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/sirupsen/logrus"
)

// HealthCheck 健康检查接口，返回服务及依赖组件的健康状态
func HealthCheck(ctx context.Context, c *app.RequestContext) {
	startTime := time.Now()

	healthData := map[string]interface{}{
		"status":    "healthy",                    // 总体健康状态
		"timestamp": time.Now().Unix(),            // 当前时间戳
		"version":   "1.0.0",                      // 服务版本
		"services":  make(map[string]interface{}), // 各依赖服务健康信息
	}

	services := healthData["services"].(map[string]interface{})

	// 检查数据库连接
	if err := checkDatabaseHealth(); err != nil {
		services["database"] = map[string]interface{}{
			"status": "unhealthy",
			"error":  err.Error(),
		}
		healthData["status"] = "degraded"
		utils.Warn("Database health check failed", logrus.Fields{"error": err.Error()})
	} else {
		services["database"] = map[string]interface{}{
			"status": "healthy",
		}
	}

	// 检查Redis连接
	if err := checkRedisHealth(); err != nil {
		services["redis"] = map[string]interface{}{
			"status": "unhealthy",
			"error":  err.Error(),
		}
		healthData["status"] = "degraded"
		utils.Warn("Redis health check failed", logrus.Fields{"error": err.Error()})
	} else {
		services["redis"] = map[string]interface{}{
			"status": "healthy",
		}
	}

	// 添加响应时间
	healthData["response_time_ms"] = time.Since(startTime).Milliseconds()

	// 根据状态设置HTTP状态码
	if healthData["status"] == "healthy" {
		utils.Success(c, healthData)
	} else {
		c.JSON(503, map[string]interface{}{
			"code":    503,
			"message": "服务状态异常",
			"data":    healthData,
		})
	}
}

// checkDatabaseHealth 检查数据库健康状态
// 实际项目中应实现数据库连接检测，这里仅返回nil作为示例
func checkDatabaseHealth() error {
	// 这里需要导入global包来访问DB
	// 为了避免循环导入，这里先返回nil
	// 实际使用时需要调整包结构
	return nil
}

// checkRedisHealth 检查Redis健康状态
// 实际项目中应实现Redis连接检测，这里仅返回nil作为示例
func checkRedisHealth() error {
	// 这里需要导入global包来访问Rdb
	// 为了避免循环导入，这里先返回nil
	// 实际使用时需要调整包结构
	return nil
}
