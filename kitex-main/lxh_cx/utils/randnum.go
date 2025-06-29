package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func RandString(i int) interface{} {
	// 验证长度参数有效性
	if i < 1 {
		return ""
	}
	if i > 18 {
		return ""
	}

	// 初始化随机种子
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// 计算生成范围
	min := pow10(i - 1) // 10^(length-1) 确保首位非0
	max := pow10(i) - 1 // 10^length - 1

	// 生成范围内的随机数
	num := min + rand.Int63n(max-min+1)
	return fmt.Sprintf("%d", num)
}

// 辅助函数：计算10的幂
func pow10(n int) int64 {
	result := int64(1)
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}
