package utils

import (
	"Didi/rpc/common/global"
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

// RedisGeoManager Redis地理位置管理器
type RedisGeoManager struct {
	ctx context.Context
	rdb *redis.Client
}

// NewRedisGeoManager 创建Redis地理位置管理器实例
func NewRedisGeoManager() *RedisGeoManager {
	return &RedisGeoManager{
		ctx: context.Background(),
		rdb: global.Rdb,
	}
}

// DriverLocation 司机位置信息
type DriverLocation struct {
	DriverID  string  `json:"driver_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    string  `json:"status"`   // online, offline, busy
	CarType   string  `json:"car_type"` // 车型
	UpdatedAt int64   `json:"updated_at"`
}

// NearbyDriver 附近司机信息
type NearbyDriver struct {
	DriverID  string  `json:"driver_id"`
	Distance  float64 `json:"distance"` // 距离（米）
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

const (
	GEO_KEY_DRIVERS = "drivers:location" // 司机位置的GEO key
	STATUS_ONLINE   = "online"           // 在线状态
	STATUS_OFFLINE  = "offline"          // 离线状态
	STATUS_BUSY     = "busy"             // 忙碌状态
)

// UpdateDriverLocation 更新司机位置
func (g *RedisGeoManager) UpdateDriverLocation(location DriverLocation) error {
	// 添加司机到地理位置集合
	err := g.rdb.GeoAdd(g.ctx, GEO_KEY_DRIVERS, &redis.GeoLocation{
		Name:      location.DriverID,
		Longitude: location.Longitude,
		Latitude:  location.Latitude,
	}).Err()

	if err != nil {
		return fmt.Errorf("更新司机地理位置失败: %v", err)
	}

	// 同时在哈希表中存储司机的详细信息
	driverInfoKey := "driver:info:" + location.DriverID
	driverInfo := map[string]interface{}{
		"status":     location.Status,
		"car_type":   location.CarType,
		"latitude":   location.Latitude,
		"longitude":  location.Longitude,
		"updated_at": location.UpdatedAt,
	}

	err = g.rdb.HMSet(g.ctx, driverInfoKey, driverInfo).Err()
	if err != nil {
		return fmt.Errorf("更新司机信息失败: %v", err)
	}

	return nil
}

// GetNearbyDrivers 获取指定位置周边的司机
func (g *RedisGeoManager) GetNearbyDrivers(latitude, longitude float64, radiusKm float64, carType string) ([]NearbyDriver, error) {
	// 使用GEORADIUS查找周边司机
	results, err := g.rdb.GeoRadius(g.ctx, GEO_KEY_DRIVERS, longitude, latitude, &redis.GeoRadiusQuery{
		Radius:    radiusKm,
		Unit:      "km",
		WithCoord: true,  // 返回坐标
		WithDist:  true,  // 返回距离
		Count:     50,    // 最多返回50个司机
		Sort:      "ASC", // 按距离升序排列
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("查找周边司机失败: %v", err)
	}

	var nearbyDrivers []NearbyDriver

	// 遍历结果，筛选在线且车型匹配的司机
	for _, result := range results {
		driverID := result.Name

		// 获取司机详细信息
		driverInfo, err := g.GetDriverInfo(driverID)
		if err != nil {
			continue // 跳过获取信息失败的司机
		}

		// 筛选条件：在线状态 && 车型匹配（如果指定了车型）
		if driverInfo["status"] == STATUS_ONLINE {
			if carType == "" || driverInfo["car_type"] == carType {
				// 转换距离为米
				distanceKm := result.Dist
				distanceM := distanceKm * 1000

				nearbyDriver := NearbyDriver{
					DriverID:  driverID,
					Distance:  distanceM,
					Latitude:  result.Latitude,
					Longitude: result.Longitude,
				}
				nearbyDrivers = append(nearbyDrivers, nearbyDriver)
			}
		}
	}

	return nearbyDrivers, nil
}

// GetDriverInfo 获取司机详细信息
func (g *RedisGeoManager) GetDriverInfo(driverID string) (map[string]string, error) {
	driverInfoKey := "driver:info:" + driverID
	return g.rdb.HGetAll(g.ctx, driverInfoKey).Result()
}

// UpdateDriverStatus 更新司机状态
func (g *RedisGeoManager) UpdateDriverStatus(driverID, status string) error {
	driverInfoKey := "driver:info:" + driverID
	return g.rdb.HSet(g.ctx, driverInfoKey, "status", status).Err()
}

// RemoveDriver 移除司机（下线时调用）
func (g *RedisGeoManager) RemoveDriver(driverID string) error {
	// 从地理位置集合中移除
	err := g.rdb.ZRem(g.ctx, GEO_KEY_DRIVERS, driverID).Err()
	if err != nil {
		return fmt.Errorf("从地理位置集合移除司机失败: %v", err)
	}

	// 删除司机信息
	driverInfoKey := "driver:info:" + driverID
	err = g.rdb.Del(g.ctx, driverInfoKey).Err()
	if err != nil {
		return fmt.Errorf("删除司机信息失败: %v", err)
	}

	return nil
}

// GetDriverLocation 获取司机当前位置
func (g *RedisGeoManager) GetDriverLocation(driverID string) (*DriverLocation, error) {
	// 获取地理位置坐标
	positions, err := g.rdb.GeoPos(g.ctx, GEO_KEY_DRIVERS, driverID).Result()
	if err != nil {
		return nil, fmt.Errorf("获取司机位置失败: %v", err)
	}

	if len(positions) == 0 || positions[0] == nil {
		return nil, fmt.Errorf("司机位置不存在")
	}

	// 获取司机详细信息
	driverInfo, err := g.GetDriverInfo(driverID)
	if err != nil {
		return nil, fmt.Errorf("获取司机信息失败: %v", err)
	}

	// 转换时间戳
	updatedAt, _ := strconv.ParseInt(driverInfo["updated_at"], 10, 64)

	location := &DriverLocation{
		DriverID:  driverID,
		Latitude:  positions[0].Latitude,
		Longitude: positions[0].Longitude,
		Status:    driverInfo["status"],
		CarType:   driverInfo["car_type"],
		UpdatedAt: updatedAt,
	}

	return location, nil
}

// GetOnlineDriverCount 获取在线司机数量
func (g *RedisGeoManager) GetOnlineDriverCount() (int64, error) {
	// 使用Scan遍历所有司机信息key
	var cursor uint64
	var count int64

	for {
		keys, nextCursor, err := g.rdb.Scan(g.ctx, cursor, "driver:info:*", 100).Result()
		if err != nil {
			return 0, fmt.Errorf("扫描司机信息失败: %v", err)
		}

		// 检查每个司机的状态
		for _, key := range keys {
			status, err := g.rdb.HGet(g.ctx, key, "status").Result()
			if err == nil && status == STATUS_ONLINE {
				count++
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return count, nil
}

// GetDriversInArea 获取指定区域内的所有司机
func (g *RedisGeoManager) GetDriversInArea(latitude, longitude, radiusKm float64) ([]string, error) {
	results, err := g.rdb.GeoRadius(g.ctx, GEO_KEY_DRIVERS, longitude, latitude, &redis.GeoRadiusQuery{
		Radius: radiusKm,
		Unit:   "km",
		Count:  1000, // 获取更多司机
		Sort:   "ASC",
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("查找区域内司机失败: %v", err)
	}

	var driverIDs []string
	for _, result := range results {
		driverIDs = append(driverIDs, result.Name)
	}

	return driverIDs, nil
}

// CalculateDistance 计算两点之间的距离
func (g *RedisGeoManager) CalculateDistance(lat1, lon1, lat2, lon2 float64) (float64, error) {
	// 临时添加两个点来计算距离
	tempKey := "temp:distance"

	// 添加两个临时位置点
	err := g.rdb.GeoAdd(g.ctx, tempKey,
		&redis.GeoLocation{Name: "point1", Longitude: lon1, Latitude: lat1},
		&redis.GeoLocation{Name: "point2", Longitude: lon2, Latitude: lat2},
	).Err()
	if err != nil {
		return 0, fmt.Errorf("添加临时位置点失败: %v", err)
	}

	// 计算距离
	distance, err := g.rdb.GeoDist(g.ctx, tempKey, "point1", "point2", "m").Result()
	if err != nil {
		return 0, fmt.Errorf("计算距离失败: %v", err)
	}

	// 清理临时数据
	g.rdb.Del(g.ctx, tempKey)

	return distance, nil
}
