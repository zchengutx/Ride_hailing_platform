// inits 包用于初始化项目运行所需的全局配置等资源
package inits

import (
	"cart/biz/dal/global"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// init 函数是Go的特殊初始化函数，包被导入时自动执行
// 这里用于自动初始化Viper配置
func init() {
	InitViper()
}

// InitViper 初始化Viper配置，读取并解析配置文件到全局配置变量
// 1. 设置配置文件路径
// 2. 读取配置文件内容
// 3. 解析内容到 global.AppConf 结构体
func InitViper() {
	viper.SetConfigFile("biz/dal/dev.yaml") // 指定配置文件路径
	err := viper.ReadInConfig()             // 读取配置文件
	if err != nil {
		panic(fmt.Errorf("viper读取失败: %s \n", err))
	}
	log.Println("viper读取成功")
	err = viper.Unmarshal(&global.AppConf) // 解析到全局配置结构体
	if err != nil {
		panic(fmt.Errorf("viper解析失败: %s \n", err))
	}
	log.Println("viper解析成功", global.AppConf)
}

// InitBizMongoDB 初始化BFF层的MongoDB连接
func InitBizMongoDB() {
	// 读取配置文件
	InitBizViper()

	var err error
	mongoConf := &global.AppConf.MongoDB

	// 添加调试信息
	log.Printf("BFF层MongoDB配置信息 - User: %s, Host: %s, Port: %d, Database: %s",
		mongoConf.User, mongoConf.Host, mongoConf.Port, mongoConf.Database)

	// 验证配置参数
	if mongoConf.Host == "" {
		log.Printf("BFF层MongoDB Host为空，跳过MongoDB初始化")
		return
	}
	if mongoConf.Port <= 0 || mongoConf.Port > 65535 {
		log.Printf("BFF层MongoDB Port无效: %d，跳过MongoDB初始化", mongoConf.Port)
		return
	}
	if mongoConf.Database == "" {
		log.Printf("BFF层MongoDB Database为空，跳过MongoDB初始化")
		return
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

		log.Printf("BFF层尝试连接方式1 (authSource=admin): %s", uri)

		clientOptions = options.Client().ApplyURI(uri).
			SetMaxPoolSize(50).
			SetMinPoolSize(5).
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
				log.Printf("BFF层MongoDB连接成功 (authSource=admin) - Database: %s", mongoConf.Database)
				global.MongoDBName = mongoConf.Database
				return
			}
		}

		log.Printf("BFF层方式1连接失败: %v, 尝试方式2", err)

		// 方式2: 尝试使用目标数据库作为认证数据库
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=%s",
			mongoConf.User,
			mongoConf.Password,
			mongoConf.Host,
			mongoConf.Port,
			mongoConf.Database,
			mongoConf.Database)

		log.Printf("BFF层尝试连接方式2 (authSource=%s): %s", mongoConf.Database, uri)

		clientOptions = options.Client().ApplyURI(uri).
			SetMaxPoolSize(50).
			SetMinPoolSize(5).
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
				log.Printf("BFF层MongoDB连接成功 (authSource=%s) - Database: %s", mongoConf.Database, mongoConf.Database)
				global.MongoDBName = mongoConf.Database
				return
			}
		}

		log.Printf("BFF层方式2连接失败: %v, 尝试方式3", err)

		// 方式3: 尝试不指定authSource
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
			mongoConf.User,
			mongoConf.Password,
			mongoConf.Host,
			mongoConf.Port,
			mongoConf.Database)

		log.Printf("BFF层尝试连接方式3 (无authSource): %s", uri)

		clientOptions = options.Client().ApplyURI(uri).
			SetMaxPoolSize(50).
			SetMinPoolSize(5).
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
				log.Printf("BFF层MongoDB连接成功 (无authSource) - Database: %s", mongoConf.Database)
				global.MongoDBName = mongoConf.Database
				return
			}
		}

		// 所有方式都失败了
		log.Printf("BFF层MongoDB认证失败，请检查用户名密码和权限。最后尝试的错误: %s", err)

	} else {
		// 无用户名密码的情况
		uri = fmt.Sprintf("mongodb://%s:%d/%s",
			mongoConf.Host,
			mongoConf.Port,
			mongoConf.Database)

		log.Printf("BFF层尝试无认证连接: %s", uri)

		clientOptions = options.Client().ApplyURI(uri).
			SetMaxPoolSize(50).
			SetMinPoolSize(5).
			SetMaxConnIdleTime(30 * time.Second).
			SetServerSelectionTimeout(5 * time.Second).
			SetConnectTimeout(10 * time.Second).
			SetSocketTimeout(30 * time.Second)

		// 连接到MongoDB
		global.MongoDB, err = mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			log.Printf("BFF层MongoDB连接失败: %s", err)
			return
		}

		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := global.MongoDB.Ping(ctx, nil); err != nil {
			log.Printf("BFF层MongoDB连接测试失败: %s", err)
			return
		}

		log.Printf("BFF层MongoDB连接成功 (无认证) - Database: %s", mongoConf.Database)
		global.MongoDBName = mongoConf.Database
	}
}

// InitBizViper 初始化BFF层配置
func InitBizViper() {
	viper.SetConfigFile("rpc/basic/dev.yaml") // 使用同一个配置文件
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("BFF层viper读取失败: %s", err)
		return
	}
	log.Println("BFF层viper读取成功")
	err = viper.Unmarshal(&global.AppConf)
	if err != nil {
		log.Printf("BFF层viper解析失败: %s", err)
		return
	}
	log.Println("BFF层viper解析成功")
}
