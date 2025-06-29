package inits

import (
	"Didi/biz/dal/global"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func init() {
	InitViper()
}

// findConfigFile 查找配置文件的路径
func findConfigFile() string {
	// 尝试的路径列表
	possiblePaths := []string{
		"Didi/biz/dal/dev.yaml",          // 从项目根目录
		"biz/dal/dev.yaml",               // 从 Didi 目录
		"dal/dev.yaml",                   // 从 biz 目录
		"dev.yaml",                       // 当前目录
		"../dev.yaml",                    // 上一级目录
		"../../dev.yaml",                 // 上两级目录
		"../../../Didi/biz/dal/dev.yaml", // 从深层目录往上找
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
	return "Didi/biz/dal/dev.yaml"
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
