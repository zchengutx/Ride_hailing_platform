package inits

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"kitexTwo/common/global"
	"log"
)

func init() {
	InitViper()
	InitMysql()
	InitRedis()
}
func InitViper() {
	viper.SetConfigFile("../kitexTwo/common/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("viper读取失败: %s \n", err))
	}
	log.Println("Viper<读取成功>")
	err = viper.Unmarshal(&global.AppConf)
	if err != nil {
		panic(fmt.Errorf("viper<解析失败>: %s \n", err))
	}
	log.Println("Viper<解析成功>", global.AppConf)
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
		panic(fmt.Errorf("mysql<连接失败>: %s \n", err))
	}
	log.Println("Mysql<连接成功>", dsn, global.DB)

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
		panic(fmt.Errorf("redis <连接失败>: %s \n", err))
	}
	log.Println("Redis<连接成功>", global.Rdb)
}
