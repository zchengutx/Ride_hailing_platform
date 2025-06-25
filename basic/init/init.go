package init

import (
	"Ride_hailing_platform/basic/config"
	"Ride_hailing_platform/basic/global"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

func init() {
	ViperInit()
	MysqlInit()
	RedisInit()
}

func ViperInit() {
	viper.SetConfigFile("./config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	} else {
		log.Println("ViperInit success")
	}
	global.Conf = &config.Viper{}
	err = viper.Unmarshal(global.Conf)
	if err != nil {
		panic(err)
	} else {
		log.Println("ViperInit success")
	}
}

func MysqlInit() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	var err error
	c := global.Conf.MysqlConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.User, c.Pass, c.Host, c.Port, c.Data)
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	log.Println("数据库初始化完成")
}

func RedisInit() {
	global.Redis = redis.NewClient(&redis.Options{
		Addr:     global.Conf.RedisConfig.Addr, // use default Addr
		Password: global.Conf.RedisConfig.Pass, // no password set
		DB:       global.Conf.RedisConfig.DB,   // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := global.Redis.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	}
	log.Println("Redis初始化完成")
}
