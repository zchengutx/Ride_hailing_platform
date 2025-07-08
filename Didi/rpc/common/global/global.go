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
	// RabbitMQ连接在utils包中管理，这里不需要全局变量
)
