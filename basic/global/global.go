package global

import (
	"Ride_hailing_platform/basic/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
	Conf  *config.Viper
)
