package init

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"lxh_cx/config"
	"lxh_cx/kitex_gen/lxh_cx/driver/driverservice"
	"lxh_cx/kitex_gen/lxh_cx/order/orderservice"
	"lxh_cx/kitex_gen/lxh_cx/user/userservice"
	"path/filepath"
	"runtime"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	GetViper()
	MysqlInit()
	RedisInit()
	UserClient()
	MongoDBInit()
	DriverClient()
	OrderClient()
}

func OrderClient() {
	c, err := orderservice.NewClient("lxh_cx.order", client.WithHostPorts("127.0.0.1:5002"))
	if err != nil {
		log.Fatal(err)
	}
	config.OrderCli = c
}

func DriverClient() {
	c, err := driverservice.NewClient("lxh_cx.driver", client.WithHostPorts("127.0.0.1:5003"))
	if err != nil {
		log.Fatal(err)
	}
	config.DriverCli = c
}

func MongoDBInit() {
	var err error

	credential := options.Credential{
		AuthSource: "lxh",
		Username:   "xiaoyu",
		Password:   "123456",
	}

	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://14.103.133.61/:27017/?authSource=lxh").SetAuth(credential)

	// 连接到MongoDB
	config.Mdb, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("MongoDB connect error, %s", err)
	}
	// 检查连接
	err = config.Mdb.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("MongoDB connect error, %s", err)
	}

	log.Println("MongoDB connect success")
}

func UserClient() {
	c, err := userservice.NewClient("lxh_cx.user", client.WithHostPorts("127.0.0.1:5001"))
	if err != nil {
		log.Fatal(err)
	}
	config.UserCli = c
}

func RedisInit() {
	config.Rdb = redis.NewClient(&redis.Options{
		Addr:     config.DataConfig.Redis.Addr,
		Password: config.DataConfig.Redis.Password, // no password set
		DB:       config.DataConfig.Redis.Db,       // use default DB
	})

	err := config.Rdb.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
}

func MysqlInit() {
	var err error
	m := config.DataConfig.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", m.User, m.Password, m.Host, m.Port, m.Database)
	config.Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}
	log.Println("mysql init success")
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := config.Db.DB()

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func GetViper() {
	// 获取当前文件所在目录
	_, filename, _, _ := runtime.Caller(0)
	configDir := filepath.Dir(filename)
	configPath := filepath.Join(configDir, "..", "kitex_info.yaml")

	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = viper.Unmarshal(&config.DataConfig)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	log.Printf("%+v", config.DataConfig)
}
