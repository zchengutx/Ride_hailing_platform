package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
)

// 用于存储已生成的字符串，确保不重复
var (
	generatedStrings = make(map[string]bool)
	mutex            sync.RWMutex
)

// GenerateUniqueSixDigitString 生成一个6位不重复的数字字符串
// 返回生成的字符串和可能的错误
func GenerateUniqueSixDigitString() (string, error) {
	maxAttempts := 1000 // 最大尝试次数，避免无限循环

	for i := 0; i < maxAttempts; i++ {
		// 生成6位随机数字
		str, err := generateRandomSixDigit()
		if err != nil {
			return "", fmt.Errorf("生成随机数字失败: %w", err)
		}

		// 检查是否已存在
		mutex.RLock()
		exists := generatedStrings[str]
		mutex.RUnlock()

		if !exists {
			// 添加到已生成集合中
			mutex.Lock()
			generatedStrings[str] = true
			mutex.Unlock()

			return str, nil
		}
	}

	return "", fmt.Errorf("无法生成不重复的6位数字字符串，已达到最大尝试次数")
}

// generateRandomSixDigit 生成一个6位随机数字字符串
func generateRandomSixDigit() (string, error) {
	// 生成0-999999之间的随机数
	max := big.NewInt(1000000) // 1000000 = 10^6
	randomNum, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// 格式化为6位数字字符串，不足6位前面补0
	return fmt.Sprintf("%06d", randomNum.Int64()), nil
}
