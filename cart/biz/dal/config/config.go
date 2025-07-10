// 配置包，定义了服务的各种配置结构体
package config

// Config 是整个服务的总配置结构体，包含了系统配置、各客户端配置等
// 用于集中管理和加载所有配置项
// 通过嵌套结构体实现各模块配置的统一管理
// 例如：Config.System 表示系统配置，Config.PassengerClient 表示乘客服务客户端配置
type Config struct {
	System          // 系统基础配置
	PassengerClient // 乘客服务客户端配置
	DriverClient    // 司机服务客户端配置
	MapClient       // 地图服务客户端配置
	WeChatClient    // 微信服务客户端配置
	Minio           // Minio 对象存储配置
}

// System 结构体定义了系统的基础配置
// Host: 服务监听的主机地址
// Port: 服务监听的端口
// Name: 服务名称
type System struct {
	Host string // 服务主机地址
	Port int    // 服务端口
	Name string // 服务名称
}

// PassengerClient 结构体定义了乘客服务客户端的连接配置
// Host: 乘客服务主机地址
// Port: 乘客服务端口
type PassengerClient struct {
	Host string // 乘客服务主机地址
	Port string // 乘客服务端口
}

// DriverClient 结构体定义了司机服务客户端的连接配置
// Host: 司机服务主机地址
// Port: 司机服务端口
type DriverClient struct {
	Host string // 司机服务主机地址
	Port string // 司机服务端口
}

// MapClient 结构体定义了地图服务客户端的连接配置
// Host: 地图服务主机地址
// Port: 地图服务端口
type MapClient struct {
	Host string // 地图服务主机地址
	Port string // 地图服务端口
}

// WeChatClient 结构体定义了微信服务客户端的连接配置
// Host: 微信服务主机地址
// Port: 微信服务端口
type WeChatClient struct {
	Host string // 微信服务主机地址
	Port string // 微信服务端口
}

// Minio 结构体定义了 Minio 对象存储的相关配置
// Endpoint: Minio 服务地址
// AccessKeyId: 访问密钥 ID
// AccessKeySecret: 访问密钥 Secret
// BucketName: 存储桶名称
// UseSsl: 是否使用 SSL
// BasePath: 基础路径
// BucketUrl: 存储桶访问 URL
type Minio struct {
	Endpoint        string // Minio 服务地址
	AccessKeyId     string // 访问密钥 ID
	AccessKeySecret string // 访问密钥 Secret
	BucketName      string // 存储桶名称
	UseSsl          string // 是否使用 SSL
	BasePath        string // 基础路径
	BucketUrl       string // 存储桶访问 URL
}
