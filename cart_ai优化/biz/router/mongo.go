package router

import (
	"cart/biz/dal/global"
	"cart/biz/handler/api"
	"cart/biz/middleware"
	"log"

	"github.com/cloudwego/hertz/pkg/route"
)

// MongoModel MongoDB数据查询模块路由
func MongoModel(r *route.RouterGroup) {
	log.Println("MongoDB数据查询服务")
	mongoModel := r.Group("/mongo")

	// 需要认证的MongoDB查询接口
	mongoAuth := mongoModel.Group("")
	mongoAuth.Use(middleware.JWTAuth(global.JWT_SELECT_KEY))
	{
		// 地图API日志查询
		mongoAuth.GET("/api-logs", api.GetMapApiLogs)             // 获取地图API调用日志
		mongoAuth.GET("/api-statistics", api.GetMapApiStatistics) // 获取API调用统计信息

		// 路线记录查询
		mongoAuth.GET("/route-records", api.GetMongoRouteRecords) // 获取路线记录
		mongoAuth.GET("/nearby-routes", api.GetNearbyRoutes)      // 获取附近的路线

		// 缓存数据查询
		mongoAuth.GET("/geocoding-cache", api.GetGeocodingCache) // 获取地理编码缓存
	}
}
