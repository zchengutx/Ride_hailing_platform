package dal

import (
	"cart/rpc/basic/model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 路线记录服务实现
type mongoRouteRecordServiceImpl struct {
	collection *mongo.Collection
}

// Create 创建路线记录
func (s *mongoRouteRecordServiceImpl) Create(ctx context.Context, record *model.MongoRouteRecord) error {
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Now()
	}
	if record.UpdatedAt.IsZero() {
		record.UpdatedAt = time.Now()
	}
	_, err := s.collection.InsertOne(ctx, record)
	return err
}

// Update 更新路线记录
func (s *mongoRouteRecordServiceImpl) Update(ctx context.Context, id primitive.ObjectID, record *model.MongoRouteRecord) error {
	record.UpdatedAt = time.Now()
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: record}}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

// FindByOrderId 根据订单ID查询路线记录
func (s *mongoRouteRecordServiceImpl) FindByOrderId(ctx context.Context, orderId int64) (*model.MongoRouteRecord, error) {
	filter := bson.D{{Key: "order_id", Value: orderId}}
	var record model.MongoRouteRecord
	err := s.collection.FindOne(ctx, filter).Decode(&record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// FindByDriverId 根据司机ID查询路线记录
func (s *mongoRouteRecordServiceImpl) FindByDriverId(ctx context.Context, driverId int64, limit int) ([]model.MongoRouteRecord, error) {
	filter := bson.D{{Key: "driver_id", Value: driverId}}
	opts := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []model.MongoRouteRecord
	if err = cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}

// FindByPassengerId 根据乘客ID查询路线记录
func (s *mongoRouteRecordServiceImpl) FindByPassengerId(ctx context.Context, passengerId int64, limit int) ([]model.MongoRouteRecord, error) {
	filter := bson.D{{Key: "passenger_id", Value: passengerId}}
	opts := options.Find().SetLimit(int64(limit)).SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []model.MongoRouteRecord
	if err = cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}

// FindNearbyRoutes 查询附近的路线记录（基于地理位置）
func (s *mongoRouteRecordServiceImpl) FindNearbyRoutes(ctx context.Context, lng, lat float64, maxDistance int, limit int) ([]model.MongoRouteRecord, error) {
	filter := bson.M{
		"start_location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{lng, lat},
				},
				"$maxDistance": maxDistance,
			},
		},
	}
	opts := options.Find().SetLimit(int64(limit))

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var records []model.MongoRouteRecord
	if err = cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}

// UpdateStatus 更新路线状态
func (s *mongoRouteRecordServiceImpl) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "route_status", Value: status},
			{Key: "updated_at", Value: time.Now()},
		}},
	}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

// 地理编码缓存服务实现
type mongoGeocodingCacheServiceImpl struct {
	collection *mongo.Collection
}

// Create 创建缓存记录
func (s *mongoGeocodingCacheServiceImpl) Create(ctx context.Context, cache *model.MongoGeocodingCache) error {
	if cache.CreatedAt.IsZero() {
		cache.CreatedAt = time.Now()
	}
	if cache.UpdatedAt.IsZero() {
		cache.UpdatedAt = time.Now()
	}
	if cache.HitCount == 0 {
		cache.HitCount = 1
	}
	_, err := s.collection.InsertOne(ctx, cache)
	return err
}

// FindByAddress 根据地址查询缓存
func (s *mongoGeocodingCacheServiceImpl) FindByAddress(ctx context.Context, address string) (*model.MongoGeocodingCache, error) {
	filter := bson.D{{Key: "address", Value: address}}
	var cache model.MongoGeocodingCache
	err := s.collection.FindOne(ctx, filter).Decode(&cache)
	if err != nil {
		return nil, err
	}
	return &cache, nil
}

// UpdateHitCount 更新命中次数
func (s *mongoGeocodingCacheServiceImpl) UpdateHitCount(ctx context.Context, address string) error {
	filter := bson.D{{Key: "address", Value: address}}
	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "hit_count", Value: 1}}},
		{Key: "$set", Value: bson.D{{Key: "updated_at", Value: time.Now()}}},
	}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

// FindNearbyCache 查询附近的地理编码缓存
func (s *mongoGeocodingCacheServiceImpl) FindNearbyCache(ctx context.Context, lng, lat float64, maxDistance int, limit int) ([]model.MongoGeocodingCache, error) {
	filter := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{lng, lat},
				},
				"$maxDistance": maxDistance,
			},
		},
	}
	opts := options.Find().SetLimit(int64(limit))

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var caches []model.MongoGeocodingCache
	if err = cursor.All(ctx, &caches); err != nil {
		return nil, err
	}
	return caches, nil
}

// CleanExpiredCache 清理过期缓存
func (s *mongoGeocodingCacheServiceImpl) CleanExpiredCache(ctx context.Context, expireDays int) (int64, error) {
	expireTime := time.Now().AddDate(0, 0, -expireDays)
	filter := bson.D{{Key: "created_at", Value: bson.D{{Key: "$lt", Value: expireTime}}}}
	result, err := s.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// 逆地理编码缓存服务实现
type mongoReverseGeocodingCacheServiceImpl struct {
	collection *mongo.Collection
}

// Create 创建缓存记录
func (s *mongoReverseGeocodingCacheServiceImpl) Create(ctx context.Context, cache *model.MongoReverseGeocodingCache) error {
	if cache.CreatedAt.IsZero() {
		cache.CreatedAt = time.Now()
	}
	if cache.UpdatedAt.IsZero() {
		cache.UpdatedAt = time.Now()
	}
	if cache.HitCount == 0 {
		cache.HitCount = 1
	}
	_, err := s.collection.InsertOne(ctx, cache)
	return err
}

// FindByLocation 根据位置查询缓存
func (s *mongoReverseGeocodingCacheServiceImpl) FindByLocation(ctx context.Context, lng, lat float64, precision float64) (*model.MongoReverseGeocodingCache, error) {
	// 使用地理空间查询查找精度范围内的缓存
	filter := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{lng, lat},
				},
				"$maxDistance": precision,
			},
		},
	}

	var cache model.MongoReverseGeocodingCache
	err := s.collection.FindOne(ctx, filter).Decode(&cache)
	if err != nil {
		return nil, err
	}
	return &cache, nil
}

// UpdateHitCount 更新命中次数
func (s *mongoReverseGeocodingCacheServiceImpl) UpdateHitCount(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "hit_count", Value: 1}}},
		{Key: "$set", Value: bson.D{{Key: "updated_at", Value: time.Now()}}},
	}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

// CleanExpiredCache 清理过期缓存
func (s *mongoReverseGeocodingCacheServiceImpl) CleanExpiredCache(ctx context.Context, expireDays int) (int64, error) {
	expireTime := time.Now().AddDate(0, 0, -expireDays)
	filter := bson.D{{Key: "created_at", Value: bson.D{{Key: "$lt", Value: expireTime}}}}
	result, err := s.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// IP位置缓存服务实现
type mongoIPLocationCacheServiceImpl struct {
	collection *mongo.Collection
}

// Create 创建缓存记录
func (s *mongoIPLocationCacheServiceImpl) Create(ctx context.Context, cache *model.MongoIPLocationCache) error {
	if cache.CreatedAt.IsZero() {
		cache.CreatedAt = time.Now()
	}
	if cache.UpdatedAt.IsZero() {
		cache.UpdatedAt = time.Now()
	}
	if cache.HitCount == 0 {
		cache.HitCount = 1
	}
	_, err := s.collection.InsertOne(ctx, cache)
	return err
}

// FindByIP 根据IP地址查询缓存
func (s *mongoIPLocationCacheServiceImpl) FindByIP(ctx context.Context, ipAddress string) (*model.MongoIPLocationCache, error) {
	filter := bson.D{{Key: "ip_address", Value: ipAddress}}
	var cache model.MongoIPLocationCache
	err := s.collection.FindOne(ctx, filter).Decode(&cache)
	if err != nil {
		return nil, err
	}
	return &cache, nil
}

// UpdateHitCount 更新命中次数
func (s *mongoIPLocationCacheServiceImpl) UpdateHitCount(ctx context.Context, ipAddress string) error {
	filter := bson.D{{Key: "ip_address", Value: ipAddress}}
	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "hit_count", Value: 1}}},
		{Key: "$set", Value: bson.D{{Key: "updated_at", Value: time.Now()}}},
	}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

// CleanExpiredCache 清理过期缓存
func (s *mongoIPLocationCacheServiceImpl) CleanExpiredCache(ctx context.Context, expireDays int) (int64, error) {
	expireTime := time.Now().AddDate(0, 0, -expireDays)
	filter := bson.D{{Key: "created_at", Value: bson.D{{Key: "$lt", Value: expireTime}}}}
	result, err := s.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// 距离计算缓存服务实现
type mongoDistanceCacheServiceImpl struct {
	collection *mongo.Collection
}

// Create 创建缓存记录
func (s *mongoDistanceCacheServiceImpl) Create(ctx context.Context, cache *model.MongoDistanceCache) error {
	if cache.CreatedAt.IsZero() {
		cache.CreatedAt = time.Now()
	}
	if cache.UpdatedAt.IsZero() {
		cache.UpdatedAt = time.Now()
	}
	if cache.HitCount == 0 {
		cache.HitCount = 1
	}
	_, err := s.collection.InsertOne(ctx, cache)
	return err
}

// FindByLocations 根据起终点查询缓存
func (s *mongoDistanceCacheServiceImpl) FindByLocations(ctx context.Context, startLng, startLat, endLng, endLat float64, precision float64) (*model.MongoDistanceCache, error) {
	// 使用复合查询查找精度范围内的起终点缓存
	filter := bson.M{
		"$and": []bson.M{
			{
				"start_location": bson.M{
					"$near": bson.M{
						"$geometry": bson.M{
							"type":        "Point",
							"coordinates": []float64{startLng, startLat},
						},
						"$maxDistance": precision,
					},
				},
			},
			{
				"end_location": bson.M{
					"$near": bson.M{
						"$geometry": bson.M{
							"type":        "Point",
							"coordinates": []float64{endLng, endLat},
						},
						"$maxDistance": precision,
					},
				},
			},
		},
	}

	var cache model.MongoDistanceCache
	err := s.collection.FindOne(ctx, filter).Decode(&cache)
	if err != nil {
		return nil, err
	}
	return &cache, nil
}

// UpdateHitCount 更新命中次数
func (s *mongoDistanceCacheServiceImpl) UpdateHitCount(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "hit_count", Value: 1}}},
		{Key: "$set", Value: bson.D{{Key: "updated_at", Value: time.Now()}}},
	}
	_, err := s.collection.UpdateOne(ctx, filter, update)
	return err
}

// CleanExpiredCache 清理过期缓存
func (s *mongoDistanceCacheServiceImpl) CleanExpiredCache(ctx context.Context, expireDays int) (int64, error) {
	expireTime := time.Now().AddDate(0, 0, -expireDays)
	filter := bson.D{{Key: "created_at", Value: bson.D{{Key: "$lt", Value: expireTime}}}}
	result, err := s.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
