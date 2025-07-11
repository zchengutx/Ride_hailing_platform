// api 包下的 mongo_logs.go 实现MongoDB日志查询相关接口
package api

import (
	"cart/biz/dal/global"
	"cart/biz/dal/inits"
	"cart/biz/utils"
	"cart/rpc/basic/dal"
	"context"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

var mongoService dal.MongoService

// 初始化MongoDB服务
func init() {
	// 延迟初始化，确保MongoDB连接已建立
	go func() {
		time.Sleep(3 * time.Second)

		// 初始化BFF层的MongoDB连接
		inits.InitBizMongoDB()

		// 创建MongoDB服务，使用BFF层的MongoDB连接
		if global.MongoDB != nil && global.MongoDBName != "" {
			mongoService = dal.NewMongoServiceWithClient(global.MongoDB, global.MongoDBName)
		}
	}()
}

// GetMapApiLogs 获取地图API调用日志
// 支持按API类型、用户ID、时间范围查询
func GetMapApiLogs(ctx context.Context, c *app.RequestContext) {
	if mongoService == nil {
		utils.InternalError(c, nil)
		return
	}

	// 解析查询参数
	apiType := string(c.Query("api_type"))
	userIdStr := string(c.Query("user_id"))
	startTimeStr := string(c.Query("start_time"))
	endTimeStr := string(c.Query("end_time"))
	limitStr := string(c.Query("limit"))

	// 设置默认限制
	limit := 50
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 1000 {
			limit = parsedLimit
		}
	}

	var logs interface{}
	var err error

	// 根据查询条件获取日志
	if apiType != "" {
		// 按API类型查询
		logs, err = mongoService.MapApiLog().FindByApiType(ctx, apiType, limit)
	} else if userIdStr != "" {
		// 按用户ID查询
		if userId, parseErr := strconv.ParseInt(userIdStr, 10, 64); parseErr == nil {
			logs, err = mongoService.MapApiLog().FindByUserId(ctx, userId, limit)
		} else {
			utils.BadRequest(c, "用户ID格式不正确", parseErr)
			return
		}
	} else if startTimeStr != "" && endTimeStr != "" {
		// 按时间范围查询
		startTime, err1 := time.Parse("2006-01-02 15:04:05", startTimeStr)
		endTime, err2 := time.Parse("2006-01-02 15:04:05", endTimeStr)

		if err1 != nil || err2 != nil {
			utils.BadRequest(c, "时间格式不正确，请使用 YYYY-MM-DD HH:MM:SS 格式", nil)
			return
		}

		logs, err = mongoService.MapApiLog().FindByTimeRange(ctx, startTime, endTime, limit)
	} else {
		// 默认获取最近的日志
		endTime := time.Now()
		startTime := endTime.Add(-24 * time.Hour) // 最近24小时
		logs, err = mongoService.MapApiLog().FindByTimeRange(ctx, startTime, endTime, limit)
	}

	if err != nil {
		utils.DatabaseError(c, err)
		return
	}

	// 获取记录数量
	var count int
	if logsSlice, ok := logs.([]interface{}); ok {
		count = len(logsSlice)
	} else {
		count = 0
	}

	utils.Success(c, map[string]interface{}{
		"logs":  logs,
		"count": count,
		"limit": limit,
	})
}

// GetMapApiStatistics 获取地图API调用统计信息
func GetMapApiStatistics(ctx context.Context, c *app.RequestContext) {
	if mongoService == nil {
		utils.InternalError(c, nil)
		return
	}

	// 解析查询参数
	startTimeStr := string(c.Query("start_time"))
	endTimeStr := string(c.Query("end_time"))

	// 默认统计最近7天的数据
	endTime := time.Now()
	startTime := endTime.Add(-7 * 24 * time.Hour)

	if startTimeStr != "" && endTimeStr != "" {
		var err1, err2 error
		startTime, err1 = time.Parse("2006-01-02 15:04:05", startTimeStr)
		endTime, err2 = time.Parse("2006-01-02 15:04:05", endTimeStr)

		if err1 != nil || err2 != nil {
			utils.BadRequest(c, "时间格式不正确，请使用 YYYY-MM-DD HH:MM:SS 格式", nil)
			return
		}
	}

	statistics, err := mongoService.MapApiLog().GetStatistics(ctx, startTime, endTime)
	if err != nil {
		utils.DatabaseError(c, err)
		return
	}

	utils.Success(c, map[string]interface{}{
		"statistics": statistics,
		"time_range": map[string]interface{}{
			"start_time": startTime.Format("2006-01-02 15:04:05"),
			"end_time":   endTime.Format("2006-01-02 15:04:05"),
		},
	})
}

// GetMongoRouteRecords 获取MongoDB中的路线记录
func GetMongoRouteRecords(ctx context.Context, c *app.RequestContext) {
	if mongoService == nil {
		utils.InternalError(c, nil)
		return
	}

	// 解析查询参数
	driverIdStr := string(c.Query("driver_id"))
	passengerIdStr := string(c.Query("passenger_id"))
	orderIdStr := string(c.Query("order_id"))
	limitStr := string(c.Query("limit"))

	// 设置默认限制
	limit := 20
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	var records interface{}
	var err error

	if orderIdStr != "" {
		// 按订单ID查询单条记录
		if orderId, parseErr := strconv.ParseInt(orderIdStr, 10, 64); parseErr == nil {
			records, err = mongoService.RouteRecord().FindByOrderId(ctx, orderId)
		} else {
			utils.BadRequest(c, "订单ID格式不正确", parseErr)
			return
		}
	} else if driverIdStr != "" {
		// 按司机ID查询
		if driverId, parseErr := strconv.ParseInt(driverIdStr, 10, 64); parseErr == nil {
			records, err = mongoService.RouteRecord().FindByDriverId(ctx, driverId, limit)
		} else {
			utils.BadRequest(c, "司机ID格式不正确", parseErr)
			return
		}
	} else if passengerIdStr != "" {
		// 按乘客ID查询
		if passengerId, parseErr := strconv.ParseInt(passengerIdStr, 10, 64); parseErr == nil {
			records, err = mongoService.RouteRecord().FindByPassengerId(ctx, passengerId, limit)
		} else {
			utils.BadRequest(c, "乘客ID格式不正确", parseErr)
			return
		}
	} else {
		utils.BadRequest(c, "请提供查询参数：order_id、driver_id 或 passenger_id", nil)
		return
	}

	if err != nil {
		utils.DatabaseError(c, err)
		return
	}

	utils.Success(c, map[string]interface{}{
		"records": records,
		"limit":   limit,
	})
}

// GetGeocodingCache 获取地理编码缓存
func GetGeocodingCache(ctx context.Context, c *app.RequestContext) {
	if mongoService == nil {
		utils.InternalError(c, nil)
		return
	}

	address := string(c.Query("address"))
	if address == "" {
		utils.BadRequest(c, "请提供地址参数", nil)
		return
	}

	cache, err := mongoService.GeocodingCache().FindByAddress(ctx, address)
	if err != nil {
		utils.NotFound(c, "未找到该地址的缓存记录")
		return
	}

	utils.Success(c, cache)
}

// GetNearbyRoutes 获取附近的路线记录
func GetNearbyRoutes(ctx context.Context, c *app.RequestContext) {
	if mongoService == nil {
		utils.InternalError(c, nil)
		return
	}

	lngStr := string(c.Query("lng"))
	latStr := string(c.Query("lat"))
	maxDistanceStr := string(c.Query("max_distance"))
	limitStr := string(c.Query("limit"))

	if lngStr == "" || latStr == "" {
		utils.BadRequest(c, "请提供经纬度参数", nil)
		return
	}

	lng, err1 := strconv.ParseFloat(lngStr, 64)
	lat, err2 := strconv.ParseFloat(latStr, 64)
	if err1 != nil || err2 != nil {
		utils.BadRequest(c, "经纬度格式不正确", nil)
		return
	}

	maxDistance := 1000 // 默认1000米
	if maxDistanceStr != "" {
		if parsed, err := strconv.Atoi(maxDistanceStr); err == nil && parsed > 0 {
			maxDistance = parsed
		}
	}

	limit := 10
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	routes, err := mongoService.RouteRecord().FindNearbyRoutes(ctx, lng, lat, maxDistance, limit)
	if err != nil {
		utils.DatabaseError(c, err)
		return
	}

	utils.Success(c, map[string]interface{}{
		"routes":       routes,
		"count":        len(routes),
		"center_point": map[string]float64{"lng": lng, "lat": lat},
		"max_distance": maxDistance,
		"limit":        limit,
	})
}
