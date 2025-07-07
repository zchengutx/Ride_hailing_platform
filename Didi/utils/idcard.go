package utils

import (
	"regexp"
)

func verifyCheckCode(idCard string) bool {
	// 身份证号码前17位的权重因子
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

	// 校验码对应表
	checkCodeMap := map[int]rune{
		0: '1', 1: '0', 2: 'X', 3: '9', 4: '8',
		5: '7', 6: '6', 7: '5', 8: '4', 9: '3', 10: '2',
	}

	// 计算前17位的加权和
	sum := 0
	for i := 0; i < 17; i++ {
		num := int(idCard[i] - '0')
		sum += num * weights[i]
	}

	// 计算校验码
	mod := sum % 11
	expectedCode := checkCodeMap[mod]

	// 比较校验码
	actualCode := rune(idCard[17])
	if actualCode >= 'a' && actualCode <= 'z' {
		actualCode -= 32 // 转换为大写
	}

	return expectedCode == actualCode
}
func IdCardVerification(idCard string) bool {
	// 18位身份证正则表达式
	regex18 := `^[1-9]\d{5}(?:18|19|20)\d{2}(?:0\d|10|11|12)(?:0[1-9]|[1-2]\d|30|31)\d{3}[\dXx]$`
	// 验证18位身份证
	match18, _ := regexp.MatchString(regex18, idCard)
	if match18 {
		// 验证校验码
		return verifyCheckCode(idCard)
	}

	return false
}
