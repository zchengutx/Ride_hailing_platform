package inits

import (
	"cart/rpc/basic/global"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	InitViper()
	InitMysql()
	InitRedis()
	InitMongoDB()
}

func InitViper() {
	viper.SetConfigFile("rpc/basic/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("viper读取失败: %s \n", err))
	}
	log.Println("viper读取成功")
	err = viper.Unmarshal(&global.AppConf)
	if err != nil {
		panic(fmt.Errorf("viper解析失败: %s \n", err))
	}
	log.Println("viper解析成功", global.AppConf)
}

var (
	MysqlConf = &global.AppConf.Mysql
	RedisConf = &global.AppConf.Redis
	PoolConf  = &global.AppConf.DatabasePool
	MongoConf = &global.AppConf.MongoDB
)

func InitMysql() {
	var err error

	// 构建DSN
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		MysqlConf.User, MysqlConf.Password, MysqlConf.Host, MysqlConf.Port, MysqlConf.Database)

	// 配置GORM日志
	var gormLogger logger.Interface
	if PoolConf.EnableSlowLog {
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Duration(PoolConf.SlowThreshold) * time.Millisecond,
				LogLevel:                  logger.Warn, // 记录慢查询和错误
				IgnoreRecordNotFoundError: true,        // 忽略记录不存在错误
				Colorful:                  true,        // 彩色输出
			},
		)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// 打开数据库连接
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
	})

	if err != nil {
		panic(fmt.Errorf("mysql连接失败: %s \n", err))
	}

	// 获取底层的sql.DB对象来配置连接池
	sqlDB, err := global.DB.DB()
	if err != nil {
		panic(fmt.Errorf("获取数据库实例失败: %s \n", err))
	}

	// 配置连接池
	sqlDB.SetMaxOpenConns(PoolConf.MaxOpenConns)                                    // 最大打开连接数
	sqlDB.SetMaxIdleConns(PoolConf.MaxIdleConns)                                    // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Duration(PoolConf.ConnMaxLifetime) * time.Minute) // 连接最大生存时间
	sqlDB.SetConnMaxIdleTime(time.Duration(PoolConf.ConnMaxIdleTime) * time.Minute) // 连接最大空闲时间

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		panic(fmt.Errorf("数据库连接测试失败: %s \n", err))
	}

	log.Printf("数据库连接成功 - DSN: %s, 连接池配置: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%dm, MaxIdleTime=%dm",
		dsn, PoolConf.MaxOpenConns, PoolConf.MaxIdleConns, PoolConf.ConnMaxLifetime, PoolConf.ConnMaxIdleTime)
}

func InitRedis() {
	ctx := context.Background()
	global.Rdb = redis.NewClient(&redis.Options{
		Addr:         RedisConf.Addr,
		Password:     RedisConf.Password,
		DB:           RedisConf.Db,
		PoolSize:     10,              // 连接池大小
		MinIdleConns: 5,               // 最小空闲连接数
		DialTimeout:  5 * time.Second, // 连接超时
		ReadTimeout:  3 * time.Second, // 读取超时
		WriteTimeout: 3 * time.Second, // 写入超时
		PoolTimeout:  4 * time.Second, // 连接池超时
		IdleTimeout:  5 * time.Minute, // 空闲连接超时
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := global.Rdb.Ping(ctx).Err()
	if err != nil {
		panic(fmt.Errorf("redis连接失败: %s \n", err))
	}
	log.Printf("Redis连接成功 - Addr: %s, DB: %d", RedisConf.Addr, RedisConf.Db)
}

func InitMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构建MongoDB连接URI
	var uri string
	if MongoConf.User != "" && MongoConf.Password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin",
			MongoConf.User, MongoConf.Password, MongoConf.Host, MongoConf.Port, MongoConf.Database)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%d/%s",
			MongoConf.Host, MongoConf.Port, MongoConf.Database)
	}

	// 配置MongoDB客户端选项
	clientOptions := options.Client().ApplyURI(uri)

	// 设置连接池配置
	clientOptions.SetMaxPoolSize(100)                        // 最大连接池大小
	clientOptions.SetMinPoolSize(5)                          // 最小连接池大小
	clientOptions.SetMaxConnIdleTime(30 * time.Minute)       // 连接最大空闲时间
	clientOptions.SetConnectTimeout(10 * time.Second)        // 连接超时
	clientOptions.SetServerSelectionTimeout(5 * time.Second) // 服务器选择超时

	// 创建MongoDB客户端
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(fmt.Errorf("MongoDB连接失败: %s \n", err))
	}

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(fmt.Errorf("MongoDB连接测试失败: %s \n", err))
	}

	global.MongoDB = client
	log.Printf("MongoDB连接成功 - Host: %s:%d, Database: %s",
		MongoConf.Host, MongoConf.Port, MongoConf.Database)
}
