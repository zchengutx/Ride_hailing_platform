package config

type Viper struct {
	MysqlConfig struct {
		User string
		Pass string
		Host string
		Port int
		Data string
	}
	RedisConfig struct {
		Addr string
		Pass string
		DB   int
	}
	WechatConfig struct {
		AppID     string // 微信应用ID
		AppSecret string // 微信应用密钥
		RedirectURL string // 授权回调地址
	}
}
