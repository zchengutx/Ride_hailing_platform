package utils

import "strconv"

// ParseFloat 将字符串解析为float64，解析失败返回0
func ParseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}
