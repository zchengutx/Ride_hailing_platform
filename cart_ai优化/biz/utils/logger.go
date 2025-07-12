package utils

import (
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Logger
)

// LogConfig 日志配置
type LogConfig struct {
	Level    string `json:"level"`
	Format   string `json:"format"` // json, text
	Output   string `json:"output"` // file, console, both
	FilePath string `json:"file_path"`
}

// InitLogger 初始化日志系统
func InitLogger(config LogConfig) error {
	Logger = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// 设置日志格式
	if config.Format == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "function",
			},
		})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// 设置日志输出
	switch config.Output {
	case "file":
		if err := setupFileOutput(config.FilePath); err != nil {
			return err
		}
	case "both":
		if err := setupFileOutput(config.FilePath); err != nil {
			return err
		}
		// 保持控制台输出
	default: // console
		Logger.SetOutput(os.Stdout)
	}

	return nil
}

// setupFileOutput 设置文件输出
func setupFileOutput(filePath string) error {
	if filePath == "" {
		filePath = "logs/app.log"
	}

	// 确保日志目录存在
	logDir := filepath.Dir(filePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 打开日志文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	Logger.SetOutput(file)
	return nil
}

// RequestLogger 请求日志记录器
type RequestLogger struct {
	*logrus.Entry
}

// NewRequestLogger 创建请求日志记录器
func NewRequestLogger(requestID, method, path, clientIP string) *RequestLogger {
	entry := Logger.WithFields(logrus.Fields{
		"request_id": requestID,
		"method":     method,
		"path":       path,
		"client_ip":  clientIP,
		"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
	})
	return &RequestLogger{Entry: entry}
}

// LogRequest 记录请求信息
func (rl *RequestLogger) LogRequest(statusCode int, latency time.Duration, responseSize int64) {
	entry := rl.WithFields(logrus.Fields{
		"status_code":   statusCode,
		"latency_ms":    latency.Milliseconds(),
		"response_size": responseSize,
	})

	switch {
	case statusCode >= 500:
		entry.Error("Server Error")
	case statusCode >= 400:
		entry.Warn("Client Error")
	case statusCode >= 300:
		entry.Info("Redirect")
	default:
		entry.Info("Success")
	}
}

// DatabaseLogger 数据库日志记录器
func LogDatabase(operation, table string, duration time.Duration, err error) {
	entry := Logger.WithFields(logrus.Fields{
		"component": "database",
		"operation": operation,
		"table":     table,
		"duration":  duration.Milliseconds(),
	})

	if err != nil {
		entry.WithError(err).Error("Database operation failed")
	} else if duration > 100*time.Millisecond {
		entry.Warn("Slow database operation")
	} else {
		entry.Debug("Database operation completed")
	}
}

// RedisLogger Redis日志记录器
func LogRedis(operation, key string, duration time.Duration, err error) {
	entry := Logger.WithFields(logrus.Fields{
		"component": "redis",
		"operation": operation,
		"key":       key,
		"duration":  duration.Milliseconds(),
	})

	if err != nil {
		entry.WithError(err).Error("Redis operation failed")
	} else {
		entry.Debug("Redis operation completed")
	}
}

// BusinessLogger 业务日志记录器
func LogBusiness(module, action string, userID int, details map[string]interface{}) {
	fields := logrus.Fields{
		"component": "business",
		"module":    module,
		"action":    action,
		"user_id":   userID,
	}

	// 合并详细信息
	for k, v := range details {
		fields[k] = v
	}

	Logger.WithFields(fields).Info("Business operation")
}

// SecurityLogger 安全日志记录器
func LogSecurity(event, source, userID string, success bool, details map[string]interface{}) {
	fields := logrus.Fields{
		"component": "security",
		"event":     event,
		"source":    source,
		"user_id":   userID,
		"success":   success,
	}

	// 合并详细信息
	for k, v := range details {
		fields[k] = v
	}

	entry := Logger.WithFields(fields)
	if success {
		entry.Info("Security event")
	} else {
		entry.Warn("Security event failed")
	}
}

// ErrorLogger 错误日志记录器
func LogError(component, operation string, err error, context map[string]interface{}) {
	fields := logrus.Fields{
		"component": component,
		"operation": operation,
		"error":     err.Error(),
	}

	// 合并上下文信息
	for k, v := range context {
		fields[k] = v
	}

	Logger.WithFields(fields).Error("Operation failed")
}

// 便捷方法
func Info(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Info(msg)
	} else {
		Logger.Info(msg)
	}
}

func Error(msg string, err error, fields ...logrus.Fields) {
	entry := Logger.WithError(err)
	if len(fields) > 0 {
		entry.WithFields(fields[0]).Error(msg)
	} else {
		entry.Error(msg)
	}
}

func Warn(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Warn(msg)
	} else {
		Logger.Warn(msg)
	}
}

func Debug(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Debug(msg)
	} else {
		Logger.Debug(msg)
	}
}
