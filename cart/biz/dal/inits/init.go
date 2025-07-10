// inits 包用于初始化项目运行所需的全局配置等资源
package inits

import (
	"cart/biz/dal/global"
	"fmt"
	"log"

	"github.com/spf13/viper"
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
