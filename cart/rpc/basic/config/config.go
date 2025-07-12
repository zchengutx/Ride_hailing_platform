package config

type Config struct {
	Mysql
	Redis
	PassengerService
	DriverService
	MapService
	WeChatService
	AppConfig
	SMSConfig
	BaiduMapConfig
	WeatherConfig
	WeChatConfig
	DatabasePool
	MongoDB
}

type Mysql struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

// DatabasePool 数据库连接池配置
type DatabasePool struct {
	MaxOpenConns    int  // 最大打开连接数
	MaxIdleConns    int  // 最大空闲连接数
	ConnMaxLifetime int  // 连接最大生存时间（分钟）
	ConnMaxIdleTime int  // 连接最大空闲时间（分钟）
	EnableSlowLog   bool // 是否启用慢查询日志
	SlowThreshold   int  // 慢查询阈值（毫秒）
}

type Redis struct {
	Addr     string
	Password string
	Db       int
}

type PassengerService struct {
	Host string
	Port int
}

type DriverService struct {
	Host string
	Port int
}

type MapService struct {
	Host string
	Port int
}

type WeChatService struct {
	Host string
	Port int
}

// 应用配置
type AppConfig struct {
	AppName         string // 应用名称
	WelcomeMessage  string // 欢迎消息
	DefaultLocation string // 默认位置
	NicknamePrefix  string // 昵称前缀
}

// 短信配置
type SMSConfig struct {
	CodeLength     int // 验证码长度
	CodeExpireTime int // 验证码过期时间（分钟）
	MinCodeValue   int // 验证码最小值
	MaxCodeValue   int // 验证码最大值
}

// 百度地图配置
type BaiduMapConfig struct {
	APIKey  string // 百度地图API Key
	APIHost string // 百度地图API Host
}

// 天气配置
type WeatherConfig struct {
	DefaultWeather string // 默认天气信息
	APIKey         string // 天气API Key（预留）
	APIHost        string // 天气API Host（预留）
}

// 微信配置
type WeChatConfig struct {
	AppID          string // 微信公众号AppID
	AppSecret      string // 微信公众号AppSecret
	Token          string // 微信公众号Token
	EncodingAESKey string // 微信公众号EncodingAESKey
	RedirectURI    string // 微信授权回调地址
}

// MongoDB配置
type MongoDB struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}
