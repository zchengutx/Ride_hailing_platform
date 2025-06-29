package config

type Config struct {
	System
}
type System struct {
	Host string
	Port int
	Name string
}
