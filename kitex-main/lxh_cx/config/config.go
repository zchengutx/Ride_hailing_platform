package config

type AppConfig struct {
	Mysql struct {
		User     string
		Password string
		Host     string
		Port     int
		Database string
	}
	Redis struct {
		Addr     string
		Password string
		Db       int
	}
	ALiYun struct {
		AccessKeyID     string
		AccessKeySecret string
	}
}
