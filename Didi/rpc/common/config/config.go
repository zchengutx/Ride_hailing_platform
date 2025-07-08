package config

type Config struct {
	Mysql
	Redis
	RabbitMQ
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
type RabbitMQ struct {
	Host     string
	Port     int
	User     string
	Password string
	Vhost    string
}
type UserClient struct {
	Host string
	Port int
}
type DriverClient struct {
	Host string
	Port int
}
