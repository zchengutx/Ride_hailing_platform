package utils

import (
	"regexp"
	"strconv"
)

// 验证日期是否有效
func isValidDate(year, month, day int) bool {
	// 月份范围检查
	if month < 1 || month > 12 {
		return false
	}

	// 每月天数
	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	// 闰年判断
	if isLeapYear(year) {
		daysInMonth[1] = 29
	}

	// 日期范围检查
	if day < 1 || day > daysInMonth[month-1] {
		return false
	}

	return true
}

// 判断是否为闰年
func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || year%400 == 0
}

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
	// 改进的18位身份证正则表达式
	// 支持更广泛的年份范围(1800-2099)，修正月份验证
	regex18 := `^[1-9]\d{5}(18|19|20|21)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dXx]$`

	// 验证基本格式
	match18, _ := regexp.MatchString(regex18, idCard)
	if !match18 {
		return false
	}

	// 提取出生日期进行精确验证
	yearStr := idCard[6:10]
	monthStr := idCard[10:12]
	dayStr := idCard[12:14]

	year, err1 := strconv.Atoi(yearStr)
	month, err2 := strconv.Atoi(monthStr)
	day, err3 := strconv.Atoi(dayStr)

	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}

	// 验证日期有效性
	if !isValidDate(year, month, day) {
		return false
	}

	// 验证校验码
	return verifyCheckCode(idCard)
}
