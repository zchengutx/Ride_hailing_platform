package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoMapApiLog MongoDB版本的地图API调用日志
type MongoMapApiLog struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	ApiType       string                 `bson:"api_type" json:"api_type"`                               // API类型
	RequestParams map[string]interface{} `bson:"request_params" json:"request_params"`                   // 请求参数
	ResponseData  map[string]interface{} `bson:"response_data" json:"response_data"`                     // 响应数据
	ResponseCode  int                    `bson:"response_code" json:"response_code"`                     // 响应状态码
	ResponseTime  int                    `bson:"response_time" json:"response_time"`                     // 响应时间(毫秒)
	IpAddress     string                 `bson:"ip_address" json:"ip_address"`                           // 客户端IP
	UserAgent     string                 `bson:"user_agent" json:"user_agent"`                           // 用户代理
	UserId        int64                  `bson:"user_id" json:"user_id"`                                 // 用户ID
	ErrorMessage  string                 `bson:"error_message,omitempty" json:"error_message,omitempty"` // 错误信息
	IsSuccess     bool                   `bson:"is_success" json:"is_success"`                           // 是否成功
	CreatedAt     time.Time              `bson:"created_at" json:"created_at"`                           // 创建时间
}

// GeoPoint 地理位置点
type GeoPoint struct {
	Type        string    `bson:"type" json:"type"`               // GeoJSON类型，固定为"Point"
	Coordinates []float64 `bson:"coordinates" json:"coordinates"` // [经度, 纬度]
}

// MongoRouteRecord MongoDB版本的路线记录
type MongoRouteRecord struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderId       int64              `bson:"order_id" json:"order_id"`                     // 关联订单ID
	PassengerId   int64              `bson:"passenger_id" json:"passenger_id"`             // 乘客ID
	DriverId      int64              `bson:"driver_id" json:"driver_id"`                   // 司机ID
	StartAddress  string             `bson:"start_address" json:"start_address"`           // 起点地址
	StartLocation GeoPoint           `bson:"start_location" json:"start_location"`         // 起点坐标(支持地理空间索引)
	EndAddress    string             `bson:"end_address" json:"end_address"`               // 终点地址
	EndLocation   GeoPoint           `bson:"end_location" json:"end_location"`             // 终点坐标(支持地理空间索引)
	Distance      float64            `bson:"distance" json:"distance"`                     // 距离(公里)
	EstimatedTime int                `bson:"estimated_time" json:"estimated_time"`         // 预估时间(分钟)
	ActualTime    int                `bson:"actual_time" json:"actual_time"`               // 实际用时(分钟)
	RoutePoints   []GeoPoint         `bson:"route_points" json:"route_points"`             // 路径坐标点数组
	RouteStatus   string             `bson:"route_status" json:"route_status"`             // 路径状态
	StartTime     time.Time          `bson:"start_time" json:"start_time"`                 // 开始时间
	EndTime       *time.Time         `bson:"end_time,omitempty" json:"end_time,omitempty"` // 结束时间
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`                 // 创建时间
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`                 // 更新时间
}

// MongoGeocodingCache MongoDB版本的地理编码缓存
type MongoGeocodingCache struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Address   string             `bson:"address" json:"address"`       // 详细地址
	Location  GeoPoint           `bson:"location" json:"location"`     // 地理位置(支持地理空间索引)
	Precise   int                `bson:"precise" json:"precise"`       // 精确度
	Level     string             `bson:"level" json:"level"`           // 地址级别
	Status    int                `bson:"status" json:"status"`         // API返回状态
	Source    string             `bson:"source" json:"source"`         // 数据来源
	HitCount  int                `bson:"hit_count" json:"hit_count"`   // 命中次数
	CreatedAt time.Time          `bson:"created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"` // 更新时间
}

// MongoReverseGeocodingCache MongoDB版本的逆地理编码缓存
type MongoReverseGeocodingCache struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Location  GeoPoint           `bson:"location" json:"location"`                   // 地理位置(支持地理空间索引)
	Address   string             `bson:"address" json:"address"`                     // 详细地址
	Formatted string             `bson:"formatted_address" json:"formatted_address"` // 格式化地址
	Status    int                `bson:"status" json:"status"`                       // API返回状态
	Source    string             `bson:"source" json:"source"`                       // 数据来源
	HitCount  int                `bson:"hit_count" json:"hit_count"`                 // 命中次数
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`               // 创建时间
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`               // 更新时间
}

// MongoIPLocationCache MongoDB版本的IP位置缓存
type MongoIPLocationCache struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	IpAddress     string             `bson:"ip_address" json:"ip_address"`         // IP地址
	Location      GeoPoint           `bson:"location" json:"location"`             // 地理位置
	Address       string             `bson:"address" json:"address"`               // 详细地址
	AddressDetail string             `bson:"address_detail" json:"address_detail"` // 地址详情
	ISP           string             `bson:"isp" json:"isp"`                       // 运营商
	Source        string             `bson:"source" json:"source"`                 // 数据来源
	HitCount      int                `bson:"hit_count" json:"hit_count"`           // 命中次数
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`         // 创建时间
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`         // 更新时间
}

// MongoDistanceCache MongoDB版本的距离计算缓存
type MongoDistanceCache struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StartLocation GeoPoint           `bson:"start_location" json:"start_location"` // 起点坐标
	EndLocation   GeoPoint           `bson:"end_location" json:"end_location"`     // 终点坐标
	Distance      float64            `bson:"distance" json:"distance"`             // 距离(公里)
	Duration      int                `bson:"duration" json:"duration"`             // 时长(分钟)
	Source        string             `bson:"source" json:"source"`                 // 数据来源
	HitCount      int                `bson:"hit_count" json:"hit_count"`           // 命中次数
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`         // 创建时间
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`         // 更新时间
}

// 集合名称常量
const (
	CollectionMapApiLog             = "map_api_logs"
	CollectionRouteRecord           = "route_records"
	CollectionGeocodingCache        = "geocoding_cache"
	CollectionReverseGeocodingCache = "reverse_geocoding_cache"
	CollectionIPLocationCache       = "ip_location_cache"
	CollectionDistanceCache         = "distance_cache"
)

// NewGeoPoint 创建GeoJSON格式的地理位置点
func NewGeoPoint(lng, lat float64) GeoPoint {
	return GeoPoint{
		Type:        "Point",
		Coordinates: []float64{lng, lat},
	}
}

// GetLng 获取经度
func (gp GeoPoint) GetLng() float64 {
	if len(gp.Coordinates) >= 2 {
		return gp.Coordinates[0]
	}
	return 0
}

// GetLat 获取纬度
func (gp GeoPoint) GetLat() float64 {
	if len(gp.Coordinates) >= 2 {
		return gp.Coordinates[1]
	}
	return 0
}
