package config

type AppConfig struct {
	Mysql Mysql
	Redis Redis
}

type Mysql struct {
	User     string
	Password string
	Host     string
	Port     int32
	Database string
}

type Redis struct {
	Addr     string
	Password string
	DB       int32
}
