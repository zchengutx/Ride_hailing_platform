package inits

import (
	"fmt"
	"github.com/spf13/viper"
	"kitexTwo/client/basic/config"
	"log"
)

func init() {
	viper.SetConfigFile("../client/dev.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("viper读取失败: %s \n", err))
	}
	log.Println("viper读取成功")
	err = viper.Unmarshal(&config.AppConf)
	if err != nil {
		panic(fmt.Errorf("viper<解析失败>: %s \n", err))
	}
	log.Println("viper<解析成功>", config.AppConf)
}
