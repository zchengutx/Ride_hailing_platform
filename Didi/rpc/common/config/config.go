package config

type Config struct {
	Mysql
	Redis
	UserClient
	DriverClient
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
type UserClient struct {
	Host string
	Port int
}
type DriverClient struct {
	Host string
	Port int
}
