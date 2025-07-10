// 配置包，负责加载和管理环境变量相关的配置
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// EnvConfig 环境配置结构体，包含数据库、Redis、JWT、微信、百度地图、MinIO及应用本身的相关配置
// 通过环境变量进行初始化，便于不同环境下灵活配置
type EnvConfig struct {
	// 数据库配置
	DBUser     string // 数据库用户名
	DBPassword string // 数据库密码
	DBHost     string // 数据库主机地址
	DBPort     int    // 数据库端口
	DBName     string // 数据库名称

	// Redis配置
	RedisAddr     string // Redis 地址
	RedisPassword string // Redis 密码
	RedisDB       int    // Redis 数据库编号

	// JWT配置
	JWTSecretKey string // JWT 密钥

	// 微信配置
	WeChatAppID          string // 微信AppID
	WeChatAppSecret      string // 微信AppSecret
	WeChatToken          string // 微信Token
	WeChatEncodingAESKey string // 微信消息加解密Key

	// 百度地图配置
	BaiduMapAPIKey string // 百度地图API Key

	// MinIO配置
	MinIOEndpoint   string // MinIO服务地址
	MinIOAccessKey  string // MinIO访问Key
	MinIOSecretKey  string // MinIO密钥
	MinIOBucketName string // MinIO桶名称

	// 应用配置
	AppEnv   string // 应用环境（如development/production）
	LogLevel string // 日志级别
}

// LoadEnvConfig 加载环境变量配置，优先从.env文件读取，否则读取系统环境变量
// 返回 EnvConfig 结构体指针
func LoadEnvConfig() *EnvConfig {
	// 尝试加载.env文件（如果存在）
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &EnvConfig{
		// 数据库配置
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnvAsInt("DB_PORT", 3306),
		DBName:     getEnv("DB_NAME", "cart"),

		// Redis配置
		RedisAddr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		// JWT配置
		JWTSecretKey: getEnv("JWT_SECRET_KEY", "default_jwt_secret_key_change_in_production"),

		// 微信配置
		WeChatAppID:          getEnv("WECHAT_APP_ID", ""),
		WeChatAppSecret:      getEnv("WECHAT_APP_SECRET", ""),
		WeChatToken:          getEnv("WECHAT_TOKEN", ""),
		WeChatEncodingAESKey: getEnv("WECHAT_ENCODING_AES_KEY", ""),

		// 百度地图配置
		BaiduMapAPIKey: getEnv("BAIDU_MAP_API_KEY", ""),

		// MinIO配置
		MinIOEndpoint:   getEnv("MINIO_ENDPOINT", "127.0.0.1:9000"),
		MinIOAccessKey:  getEnv("MINIO_ACCESS_KEY", ""),
		MinIOSecretKey:  getEnv("MINIO_SECRET_KEY", ""),
		MinIOBucketName: getEnv("MINIO_BUCKET_NAME", "test"),

		// 应用配置
		AppEnv:   getEnv("APP_ENV", "development"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv 获取指定key的环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取指定key的环境变量并转换为int类型，如果不存在或转换失败则返回默认值
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: Invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}

// ValidateConfig 验证配置项是否合理，主要检查必填项是否为空
// 若有必填项未配置，则输出警告日志
func (c *EnvConfig) ValidateConfig() error {
	requiredFields := map[string]string{
		"数据库密码": c.DBPassword,
		"JWT密钥": c.JWTSecretKey,
	}

	for fieldName, value := range requiredFields {
		if value == "" {
			log.Printf("Warning: %s 未配置，请检查环境变量", fieldName)
		}
	}

	return nil
}
