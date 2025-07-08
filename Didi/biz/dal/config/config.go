package config

type Config struct {
	System
	UserClient
	DriverClient
	MiNio
}
type System struct {
	Host string
	Port int
	Name string
}
type MiNio struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
	UseSsl          string
	BasePath        string
	BucketUrl       string
}
type UserClient struct {
	Host string
	Port string
}
type DriverClient struct {
	Host string
	Port string
}
