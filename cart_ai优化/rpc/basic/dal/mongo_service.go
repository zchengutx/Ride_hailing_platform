package dal

import (
	"cart/rpc/basic/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoMapApiLogService 地图API日志服务接口
type MongoMapApiLogService interface {
	// Create 创建API日志记录
	Create(ctx context.Context, log *model.MongoMapApiLog) error

	// FindByApiType 根据API类型查询日志
	FindByApiType(ctx context.Context, apiType string, limit int) ([]model.MongoMapApiLog, error)

	// FindByUserId 根据用户ID查询日志
	FindByUserId(ctx context.Context, userId int64, limit int) ([]model.MongoMapApiLog, error)

	// FindByTimeRange 根据时间范围查询日志
	FindByTimeRange(ctx context.Context, start, end time.Time, limit int) ([]model.MongoMapApiLog, error)

	// GetStatistics 获取API调用统计信息
	GetStatistics(ctx context.Context, start, end time.Time) (map[string]interface{}, error)
}

// MongoRouteRecordService 路线记录服务接口
type MongoRouteRecordService interface {
	// Create 创建路线记录
	Create(ctx context.Context, record *model.MongoRouteRecord) error

	// Update 更新路线记录
	Update(ctx context.Context, id primitive.ObjectID, record *model.MongoRouteRecord) error

	// FindByOrderId 根据订单ID查询路线记录
	FindByOrderId(ctx context.Context, orderId int64) (*model.MongoRouteRecord, error)

	// FindByDriverId 根据司机ID查询路线记录
	FindByDriverId(ctx context.Context, driverId int64, limit int) ([]model.MongoRouteRecord, error)

	// FindByPassengerId 根据乘客ID查询路线记录
	FindByPassengerId(ctx context.Context, passengerId int64, limit int) ([]model.MongoRouteRecord, error)

	// FindNearbyRoutes 查询附近的路线记录（基于地理位置）
	FindNearbyRoutes(ctx context.Context, lng, lat float64, maxDistance int, limit int) ([]model.MongoRouteRecord, error)

	// UpdateStatus 更新路线状态
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error
}

// MongoGeocodingCacheService 地理编码缓存服务接口
type MongoGeocodingCacheService interface {
	// Create 创建缓存记录
	Create(ctx context.Context, cache *model.MongoGeocodingCache) error

	// FindByAddress 根据地址查询缓存
	FindByAddress(ctx context.Context, address string) (*model.MongoGeocodingCache, error)

	// UpdateHitCount 更新命中次数
	UpdateHitCount(ctx context.Context, address string) error

	// FindNearbyCache 查询附近的地理编码缓存
	FindNearbyCache(ctx context.Context, lng, lat float64, maxDistance int, limit int) ([]model.MongoGeocodingCache, error)

	// CleanExpiredCache 清理过期缓存
	CleanExpiredCache(ctx context.Context, expireDays int) (int64, error)
}

// MongoReverseGeocodingCacheService 逆地理编码缓存服务接口
type MongoReverseGeocodingCacheService interface {
	// Create 创建缓存记录
	Create(ctx context.Context, cache *model.MongoReverseGeocodingCache) error

	// FindByLocation 根据位置查询缓存
	FindByLocation(ctx context.Context, lng, lat float64, precision float64) (*model.MongoReverseGeocodingCache, error)

	// UpdateHitCount 更新命中次数
	UpdateHitCount(ctx context.Context, id primitive.ObjectID) error

	// CleanExpiredCache 清理过期缓存
	CleanExpiredCache(ctx context.Context, expireDays int) (int64, error)
}

// MongoIPLocationCacheService IP位置缓存服务接口
type MongoIPLocationCacheService interface {
	// Create 创建缓存记录
	Create(ctx context.Context, cache *model.MongoIPLocationCache) error

	// FindByIP 根据IP地址查询缓存
	FindByIP(ctx context.Context, ipAddress string) (*model.MongoIPLocationCache, error)

	// UpdateHitCount 更新命中次数
	UpdateHitCount(ctx context.Context, ipAddress string) error

	// CleanExpiredCache 清理过期缓存
	CleanExpiredCache(ctx context.Context, expireDays int) (int64, error)
}

// MongoDistanceCacheService 距离计算缓存服务接口
type MongoDistanceCacheService interface {
	// Create 创建缓存记录
	Create(ctx context.Context, cache *model.MongoDistanceCache) error

	// FindByLocations 根据起终点查询缓存
	FindByLocations(ctx context.Context, startLng, startLat, endLng, endLat float64, precision float64) (*model.MongoDistanceCache, error)

	// UpdateHitCount 更新命中次数
	UpdateHitCount(ctx context.Context, id primitive.ObjectID) error

	// CleanExpiredCache 清理过期缓存
	CleanExpiredCache(ctx context.Context, expireDays int) (int64, error)
}

// MongoService MongoDB综合服务接口
type MongoService interface {
	// 子服务
	MapApiLog() MongoMapApiLogService
	RouteRecord() MongoRouteRecordService
	GeocodingCache() MongoGeocodingCacheService
	ReverseGeocodingCache() MongoReverseGeocodingCacheService
	IPLocationCache() MongoIPLocationCacheService
	DistanceCache() MongoDistanceCacheService

	// 通用方法
	CreateIndexes(ctx context.Context) error // 创建索引
	HealthCheck(ctx context.Context) error   // 健康检查
	Close(ctx context.Context) error         // 关闭连接
}
