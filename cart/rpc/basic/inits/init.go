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
	log.Println("viper解析成功")

	// 打印配置信息用于调试
	log.Printf("MySQL配置: %+v", global.AppConf.Mysql)
	log.Printf("Redis配置: %+v", global.AppConf.Redis)
	log.Printf("MongoDB配置: %+v", global.AppConf.MongoDB)
}

func InitMysql() {
	var err error

	mysqlConf := &global.AppConf.Mysql
	poolConf := &global.AppConf.DatabasePool

	// 构建DSN
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConf.User, mysqlConf.Password, mysqlConf.Host, mysqlConf.Port, mysqlConf.Database)

	// 配置GORM日志
	var gormLogger logger.Interface
	if poolConf.EnableSlowLog {
		gormLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Duration(poolConf.SlowThreshold) * time.Millisecond,
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
	sqlDB.SetMaxOpenConns(poolConf.MaxOpenConns)                                    // 最大打开连接数
	sqlDB.SetMaxIdleConns(poolConf.MaxIdleConns)                                    // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Duration(poolConf.ConnMaxLifetime) * time.Minute) // 连接最大生存时间
	sqlDB.SetConnMaxIdleTime(time.Duration(poolConf.ConnMaxIdleTime) * time.Minute) // 连接最大空闲时间

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		panic(fmt.Errorf("数据库连接测试失败: %s \n", err))
	}

	log.Printf("数据库连接成功 - DSN: %s, 连接池配置: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%dm, MaxIdleTime=%dm",
		dsn, poolConf.MaxOpenConns, poolConf.MaxIdleConns, poolConf.ConnMaxLifetime, poolConf.ConnMaxIdleTime)
}

func InitRedis() {
	ctx := context.Background()
	redisConf := &global.AppConf.Redis

	global.Rdb = redis.NewClient(&redis.Options{
		Addr:         redisConf.Addr,
		Password:     redisConf.Password,
		DB:           redisConf.Db,
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
	log.Printf("Redis连接成功 - Addr: %s, DB: %d", redisConf.Addr, redisConf.Db)
}

func InitMongoDB() {
	var err error

	mongoConf := &global.AppConf.MongoDB

	// 添加调试信息
	log.Printf("MongoDB配置信息 - User: %s, Host: %s, Port: %d, Database: %s",
		mongoConf.User, mongoConf.Host, mongoConf.Port, mongoConf.Database)

	// 验证配置参数
	if mongoConf.Host == "" {
		panic(fmt.Errorf("MongoDB Host不能为空"))
	}
	if mongoConf.Port <= 0 || mongoConf.Port > 65535 {
		panic(fmt.Errorf("MongoDB Port无效: %d", mongoConf.Port))
	}
	if mongoConf.Database == "" {
		panic(fmt.Errorf("MongoDB Database不能为空"))
	}

	// 尝试多种连接方式
	var clientOptions *options.ClientOptions
	var uri string

	if mongoConf.User != "" && mongoConf.Password != "" {
		// 方式1: 尝试使用admin作为认证数据库
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin",
			mongoConf.User,
			mongoConf.Password,
			mongoConf.Host,
			mongoConf.Port,
			mongoConf.Database)

		log.Printf("尝试连接方式1 (authSource=admin): %s", uri)

		clientOptions = options.Client().ApplyURI(uri).
			SetMaxPoolSize(100).
			SetMinPoolSize(10).
			SetMaxConnIdleTime(30 * time.Second).
			SetServerSelectionTimeout(5 * time.Second).
			SetConnectTimeout(10 * time.Second).
			SetSocketTimeout(30 * time.Second)

		// 尝试连接
		global.MongoDB, err = mongo.Connect(context.Background(), clientOptions)
		if err == nil {
			// 测试连接
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = global.MongoDB.Ping(ctx, nil)
			cancel()

			if err == nil {
				log.Printf("MongoDB连接成功 (authSource=admin) - Database: %s", mongoConf.Database)
				global.MongoDBName = mongoConf.Database
				initMongoService()
				return
			}
		}

		log.Printf("方式1连接失败: %v, 尝试方式2", err)

		// 方式2: 尝试使用目标数据库作为认证数据库
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s",
			mongoConf.User,
			mongoConf.Password,
			mongoConf.Host,
			mongoConf.Port,
			mongoConf.Database,
			mongoConf.Database)

		log.Printf("尝试连接方式2 (authSource=%s): %s", mongoConf.Database, uri)

		clientOptions = options.Client().ApplyURI(uri).
			SetMaxPoolSize(100).
			SetMinPoolSize(10).
			SetMaxConnIdleTime(30 * time.Second).
			SetServerSelectionTimeout(5 * time.Second).
			SetConnectTimeout(10 * time.Second).
			SetSocketTimeout(30 * time.Second)

		// 尝试连接
		global.MongoDB, err = mongo.Connect(context.Background(), clientOptions)
		if err == nil {
			// 测试连接
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = global.MongoDB.Ping(ctx, nil)
			cancel()

			if err == nil {
				log.Printf("MongoDB连接成功 (authSource=%s) - Database: %s", mongoConf.Database, mongoConf.Database)
				global.MongoDBName = mongoConf.Database
				initMongoService()
				return
			}
		}

		log.Printf("方式2连接失败: %v, 尝试方式3", err)

		// 方式3: 尝试不指定authSource
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
			mongoConf.User,
			mongoConf.Password,
			mongoConf.Host,
			mongoConf.Port,
			mongoConf.Database)

		log.Printf("尝试连接方式3 (无authSource): %s", uri)

		clientOptions = options.Client().ApplyURI(uri).
			SetMaxPoolSize(100).
			SetMinPoolSize(10).
			SetMaxConnIdleTime(30 * time.Second).
			SetServerSelectionTimeout(5 * time.Second).
			SetConnectTimeout(10 * time.Second).
			SetSocketTimeout(30 * time.Second)

		// 尝试连接
		global.MongoDB, err = mongo.Connect(context.Background(), clientOptions)
		if err == nil {
			// 测试连接
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = global.MongoDB.Ping(ctx, nil)
			cancel()

			if err == nil {
				log.Printf("MongoDB连接成功 (无authSource) - Database: %s", mongoConf.Database)
				global.MongoDBName = mongoConf.Database
				initMongoService()
				return
			}
		}

		// 所有方式都失败了
		panic(fmt.Errorf("MongoDB认证失败，请检查用户名密码和权限。最后尝试的错误: %s", err))

	} else {
		// 无用户名密码的情况
		uri = fmt.Sprintf("mongodb://%s:%d/%s",
			mongoConf.Host,
			mongoConf.Port,
			mongoConf.Database)

		log.Printf("尝试无认证连接: %s", uri)

		clientOptions = options.Client().ApplyURI(uri).
			SetMaxPoolSize(100).
			SetMinPoolSize(10).
			SetMaxConnIdleTime(30 * time.Second).
			SetServerSelectionTimeout(5 * time.Second).
			SetConnectTimeout(10 * time.Second).
			SetSocketTimeout(30 * time.Second)

		// 连接到MongoDB
		global.MongoDB, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			panic(fmt.Errorf("MongoDB连接失败: %s", err))
		}

		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := global.MongoDB.Ping(ctx, nil); err != nil {
			panic(fmt.Errorf("MongoDB连接测试失败: %s", err))
		}

		log.Printf("MongoDB连接成功 (无认证) - Database: %s", mongoConf.Database)
		global.MongoDBName = mongoConf.Database
		initMongoService()
	}
}

// initMongoService 初始化MongoDB服务（延迟导入避免循环依赖）
func initMongoService() {
	// 这里不能直接导入dal包，需要在业务逻辑中手动初始化
	// 或者使用工厂模式来创建服务实例
}
