package config

type Config struct {
	System
	UserClient
}
type System struct {
	Host string
	Port int
	Name string
}
type UserClient struct {
	Host string
	Port int
}
