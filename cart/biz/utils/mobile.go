package utils

import (
	"regexp"
)

// ValidatePhone 验证手机号格式
func ValidatePhone(mobile string) bool {
	// 手机号正则表达式，支持中国大陆手机号码
	// 规则：1开头，第二位3-9，后面跟9位数字
	reg := `^1[3-9]\d{9}$`
	match, _ := regexp.MatchString(reg, mobile)
	return match
}

func ValidateMobile(mobile string) bool {
	if !ValidatePhone(mobile) {
		return false
	}
	return true
}
