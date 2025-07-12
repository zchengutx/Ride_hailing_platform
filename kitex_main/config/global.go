package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	Config AppConfig
	DB     *gorm.DB
	RDB    *redis.Client
	Ctx    = context.Background()
	Client *mongo.Client
)

const (
	WECHAT_TOKEN = "111"
	JWT_TOKEN    = "www.topgoer.com"
)
