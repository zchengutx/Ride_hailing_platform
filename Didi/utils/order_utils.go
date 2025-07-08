package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateOrderId 生成唯一订单ID
// 格式：年月日时分秒 + 4位随机数
// 例如：20231215143052001
func GenerateOrderId() string {
	now := time.Now()
	// 格式化时间：年月日时分秒
	timeStr := now.Format("20060102150405")
	// 生成4位随机数
	randomNum := rand.Intn(9999) + 1
	return fmt.Sprintf("%s%04d", timeStr, randomNum)
}

// GenerateOrderIdWithPrefix 生成带前缀的订单ID
// 例如：ORDER20231215143052001
func GenerateOrderIdWithPrefix(prefix string) string {
	return prefix + GenerateOrderId()
}
