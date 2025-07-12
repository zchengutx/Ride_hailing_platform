package dal

import (
	"cart/rpc/basic/global"
	"cart/rpc/basic/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoServiceImpl MongoDB服务实现
type mongoServiceImpl struct {
	db *mongo.Database
}

// NewMongoService 创建MongoDB服务实例（使用RPC层global）
func NewMongoService() MongoService {
	return &mongoServiceImpl{
		db: global.MongoDB.Database(global.MongoDBName),
	}
}

// NewMongoServiceWithClient 创建MongoDB服务实例（使用自定义客户端）
func NewMongoServiceWithClient(client *mongo.Client, dbName string) MongoService {
	return &mongoServiceImpl{
		db: client.Database(dbName),
	}
}

// MapApiLog 获取地图API日志服务
func (s *mongoServiceImpl) MapApiLog() MongoMapApiLogService {
	return &mongoMapApiLogServiceImpl{
		collection: s.db.Collection(model.CollectionMapApiLog),
	}
}

// RouteRecord 获取路线记录服务
func (s *mongoServiceImpl) RouteRecord() MongoRouteRecordService {
	return &mongoRouteRecordServiceImpl{
		collection: s.db.Collection(model.CollectionRouteRecord),
	}
}

// GeocodingCache 获取地理编码缓存服务
func (s *mongoServiceImpl) GeocodingCache() MongoGeocodingCacheService {
	return &mongoGeocodingCacheServiceImpl{
		collection: s.db.Collection(model.CollectionGeocodingCache),
	}
}

// ReverseGeocodingCache 获取逆地理编码缓存服务
func (s *mongoServiceImpl) ReverseGeocodingCache() MongoReverseGeocodingCacheService {
	return &mongoReverseGeocodingCacheServiceImpl{
		collection: s.db.Collection(model.CollectionReverseGeocodingCache),
	}
}

// IPLocationCache 获取IP位置缓存服务
func (s *mongoServiceImpl) IPLocationCache() MongoIPLocationCacheService {
	return &mongoIPLocationCacheServiceImpl{
		collection: s.db.Collection(model.CollectionIPLocationCache),
	}
}

// DistanceCache 获取距离计算缓存服务
func (s *mongoServiceImpl) DistanceCache() MongoDistanceCacheService {
	return &mongoDistanceCacheServiceImpl{
		collection: s.db.Collection(model.CollectionDistanceCache),
	}
}

// CreateIndexes 创建所有集合的索引
func (s *mongoServiceImpl) CreateIndexes(ctx context.Context) error {
	// 地图API日志索引
	_, err := s.db.Collection(model.CollectionMapApiLog).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "api_type", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "is_success", Value: 1}}},
	})
	if err != nil {
		return err
	}

	// 路线记录索引
	_, err = s.db.Collection(model.CollectionRouteRecord).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "order_id", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "driver_id", Value: 1}}},
		{Keys: bson.D{{Key: "passenger_id", Value: 1}}},
		{Keys: bson.D{{Key: "start_location", Value: "2dsphere"}}}, // 地理空间索引
		{Keys: bson.D{{Key: "end_location", Value: "2dsphere"}}},   // 地理空间索引
		{Keys: bson.D{{Key: "route_status", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	})
	if err != nil {
		return err
	}

	// 地理编码缓存索引
	_, err = s.db.Collection(model.CollectionGeocodingCache).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "address", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "location", Value: "2dsphere"}}}, // 地理空间索引
		{Keys: bson.D{{Key: "hit_count", Value: -1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	})
	if err != nil {
		return err
	}

	// 逆地理编码缓存索引
	_, err = s.db.Collection(model.CollectionReverseGeocodingCache).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "location", Value: "2dsphere"}}}, // 地理空间索引
		{Keys: bson.D{{Key: "hit_count", Value: -1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	})
	if err != nil {
		return err
	}

	// IP位置缓存索引
	_, err = s.db.Collection(model.CollectionIPLocationCache).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "ip_address", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "location", Value: "2dsphere"}}}, // 地理空间索引
		{Keys: bson.D{{Key: "hit_count", Value: -1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	})
	if err != nil {
		return err
	}

	// 距离计算缓存索引
	_, err = s.db.Collection(model.CollectionDistanceCache).Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "start_location", Value: "2dsphere"}}}, // 地理空间索引
		{Keys: bson.D{{Key: "end_location", Value: "2dsphere"}}},   // 地理空间索引
		{Keys: bson.D{{Key: "hit_count", Value: -1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	})
	if err != nil {
		return err
	}

	return nil
}

// HealthCheck 健康检查
func (s *mongoServiceImpl) HealthCheck(ctx context.Context) error {
	return global.MongoDB.Ping(ctx, nil)
}

// Close 关闭连接
func (s *mongoServiceImpl) Close(ctx context.Context) error {
	return global.MongoDB.Disconnect(ctx)
}

// 地图API日志服务实现
type mongoMapApiLogServiceImpl struct {
	collection *mongo.Collection
}

// Create 创建API日志记录
func (s *mongoMapApiLogServiceImpl) Create(ctx context.Context, log *model.MongoMapApiLog) error {
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}
	_, err := s.collection.InsertOne(ctx, log)
	return err
}

// FindByApiType 根据API类型查询日志
func (s *mongoMapApiLogServiceImpl) FindByApiType(ctx context.Context, apiType string, limit int) ([]model.MongoMapApiLog, error) {
	filter := bson.D{{Key: "api_type", Value: apiType}}
	opts := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []model.MongoMapApiLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// FindByUserId 根据用户ID查询日志
func (s *mongoMapApiLogServiceImpl) FindByUserId(ctx context.Context, userId int64, limit int) ([]model.MongoMapApiLog, error) {
	filter := bson.D{{Key: "user_id", Value: userId}}
	opts := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []model.MongoMapApiLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// FindByTimeRange 根据时间范围查询日志
func (s *mongoMapApiLogServiceImpl) FindByTimeRange(ctx context.Context, start, end time.Time, limit int) ([]model.MongoMapApiLog, error) {
	filter := bson.D{
		{Key: "created_at", Value: bson.D{
			{Key: "$gte", Value: start},
			{Key: "$lte", Value: end},
		}},
	}
	opts := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []model.MongoMapApiLog
	if err = cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// GetStatistics 获取API调用统计信息
func (s *mongoMapApiLogServiceImpl) GetStatistics(ctx context.Context, start, end time.Time) (map[string]interface{}, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"created_at": bson.M{
					"$gte": start,
					"$lte": end,
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"api_type":   "$api_type",
					"is_success": "$is_success",
				},
				"count":               bson.M{"$sum": 1},
				"avg_response_time":   bson.M{"$avg": "$response_time"},
				"total_response_time": bson.M{"$sum": "$response_time"},
			},
		},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	statistics := map[string]interface{}{
		"total_calls":       0,
		"success_rate":      0.0,
		"api_types":         make(map[string]interface{}),
		"avg_response_time": 0.0,
	}

	totalCalls := 0
	successCalls := 0
	totalResponseTime := 0

	for _, result := range results {
		count := result["count"].(int32)
		totalCalls += int(count)

		if result["_id"].(bson.M)["is_success"].(bool) {
			successCalls += int(count)
		}

		totalResponseTime += int(result["total_response_time"].(int32))
	}

	if totalCalls > 0 {
		statistics["total_calls"] = totalCalls
		statistics["success_rate"] = float64(successCalls) / float64(totalCalls) * 100
		statistics["avg_response_time"] = float64(totalResponseTime) / float64(totalCalls)
	}

	return statistics, nil
}
