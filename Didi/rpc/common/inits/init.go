package inits

import (
	"Didi/rpc/common/global"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	InitViper()
	InitMysql()
	InitRedis()
}

// findConfigFile 查找配置文件的路径
func findConfigFile() string {
	// 尝试的路径列表
	possiblePaths := []string{
		"Didi/rpc/common/dev.yaml",          // 从项目根目录
		"rpc/common/dev.yaml",               // 从 Didi 目录
		"common/dev.yaml",                   // 从 rpc 目录
		"dev.yaml",                          // 当前目录
		"../dev.yaml",                       // 上一级目录
		"../../dev.yaml",                    // 上两级目录
		"../../../Didi/rpc/common/dev.yaml", // 从深层目录往上找
	}

	// 获取当前工作目录
	wd, err := os.Getwd()
	if err == nil {
		log.Printf("当前工作目录: %s", wd)
	}

	// 尝试每个可能的路径
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			log.Printf("找到配置文件: %s", path)
			return path
		}

		// 也尝试绝对路径
		if wd != "" {
			absPath := filepath.Join(wd, path)
			if _, err := os.Stat(absPath); err == nil {
				log.Printf("找到配置文件: %s", absPath)
				return absPath
			}
		}
	}

	// 如果都找不到，返回默认路径
	log.Println("警告: 找不到配置文件，使用默认路径")
	return "Didi/rpc/common/dev.yaml"
}

func InitViper() {
	configPath := findConfigFile()
	viper.SetConfigFile(configPath)
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
)

func InitMysql() {
	var err error
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", MysqlConf.User, MysqlConf.Password, MysqlConf.Host, MysqlConf.Port, MysqlConf.Database)
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("mysql连接失败: %s \n", err))
	}
	log.Println("mysql连接成功", dsn, global.DB)
}
func InitRedis() {
	ctx := context.Background()
	global.Rdb = redis.NewClient(&redis.Options{
		Addr:     RedisConf.Addr,     // use default Addr
		Password: RedisConf.Password, // no password set
		DB:       RedisConf.Db,       // use default DB
	})
	err := global.Rdb.Ping(ctx).Err()
	if err != nil {
		panic(fmt.Errorf("redis连接失败: %s \n", err))
	}
	log.Println("redis连接成功", global.Rdb)
}
