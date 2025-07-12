// global 包用于存放全局变量和常量，便于在项目各处统一访问
package global

import (
	"cart/biz/dal/config"

	"go.mongodb.org/mongo-driver/mongo"
)

// AppConf 全局配置变量，保存应用的所有配置信息
var (
	AppConf     config.Config
	MongoDB     *mongo.Client
	MongoDBName string
)

// 统一定义全局常量
const (
	JWT_SELECT_KEY = "badb6fe15b05b84863ff1e2b249deab2" // JWT密钥（示例值，生产环境应安全存储）
	AppId          = "wxd475c55bec3c3d14"               // 微信AppId
	Secret         = "8f635851ef0b4b934ca54e914c469b20" // 微信AppSecret
)
