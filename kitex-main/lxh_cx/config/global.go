package config

import (
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"lxh_cx/kitex_gen/lxh_cx/driver/driverservice"
	"lxh_cx/kitex_gen/lxh_cx/order/orderservice"
	"lxh_cx/kitex_gen/lxh_cx/user/userservice"
)

var (
	DataConfig AppConfig
	Db         *gorm.DB
	Rdb        *redis.Client
	UserCli    userservice.Client
	Mdb        *mongo.Client
	OrderCli   orderservice.Client
	DriverCli  driverservice.Client
)
