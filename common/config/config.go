package config

type Config struct {
	Mysql
	Redis
}
type Mysql struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}
type Redis struct {
	Addr     string
	Password string
	Db       int
}
