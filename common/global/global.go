package global

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"kitexTwo/common/config"
)

var (
	AppConf config.Config
	DB      *gorm.DB
	Rdb     *redis.Client
)
