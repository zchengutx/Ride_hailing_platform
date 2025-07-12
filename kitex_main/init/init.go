package init

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kitex_main/config"
	"log"
)

func init() {
	InitViper()
	InitMysql()
	InitRedis()
	InitMongoDB()
}

func InitViper() {
	viper.SetConfigFile("D:\\GoWork\\Ride_hailing_platform\\kitex_main\\config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&config.Config); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	log.Println("init viper success")
}

func InitMysql() {
	conf := config.Config.Mysql

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
	var err error
	config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("MySQL connect error, %s", err)
	}

	log.Println("MySQL connect success")
}

func InitRedis() {
	conf := config.Config.Redis

	config.RDB = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,     // use default Addr
		Password: conf.Password, // no password set

		DB: int(conf.DB), // use default DB
	})

	pong, err := config.RDB.Ping(config.Ctx).Result()
	if err != nil {
		log.Fatalf("Redis connect error, %s", err)
	}
	fmt.Println(pong, err)
}

func InitMongoDB() {

	var err error

	credential := options.Credential{
		AuthSource: "lxh",
		Username:   "chen2",
		Password:   "123456",
	}

	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://14.103.134.228:27017/?authSource=lxh").SetAuth(credential)

	// 连接到MongoDB
	config.Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("MongoDB connect error, %s", err)
	}
	// 检查连接
	err = config.Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("MongoDB connect error, %s", err)
	}

	log.Println("MongoDB connect success")

}
