package global

import (
	"cart/rpc/basic/config"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	AppConf     config.Config
	DB          *gorm.DB
	Rdb         *redis.Client
	MongoDB     *mongo.Client
	MongoDBName string
)

// 导入dal包会造成循环导入，所以这里用接口类型
type MongoService interface {
	MapApiLog() interface{}
	RouteRecord() interface{}
	GeocodingCache() interface{}
	ReverseGeocodingCache() interface{}
	IPLocationCache() interface{}
	DistanceCache() interface{}
	CreateIndexes(ctx interface{}) error
	HealthCheck(ctx interface{}) error
	Close(ctx interface{}) error
}

var MongoSvc MongoService
