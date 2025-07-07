package global

import (
	"Didi/rpc/common/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	AppConf config.Config
	DB      *gorm.DB
	Rdb     *redis.Client
)
