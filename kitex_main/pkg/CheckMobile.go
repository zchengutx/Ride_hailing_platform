package pkg

import "regexp"

func CheckMobile(phone string) bool {
	// 匹配规则：^1[3456789]\d{9}$
	regRuler := "^1[3456789]\\d{9}$"
	reg := regexp.MustCompile(regRuler)
	return reg.MatchString(phone)
}
