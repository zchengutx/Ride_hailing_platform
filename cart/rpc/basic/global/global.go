package global

import (
	"cart/rpc/basic/config"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	AppConf config.Config
	DB      *gorm.DB
	Rdb     *redis.Client
	MongoDB *mongo.Client
)
