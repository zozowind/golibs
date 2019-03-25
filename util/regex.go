package util

import (
	"regexp"
)

var chinaIDCard = regexp.MustCompile(`^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`)
var chinaMilitaryID = regexp.MustCompile(`^[a-zA-Z0-9]{7,21}$`)
var chinaHMTID = regexp.MustCompile(`^[a-zA-Z0-9]{5,21}$`)
var chinaPassport = regexp.MustCompile(`^[A-Za-z0-9]{5,20}$`)
var chinaMobileRegex = regexp.MustCompile(`^1[1-9][0-9]{9}$`)

// var phoneRegex = regexp.MustCompile(`^(([0\+]\d{2,3}-)?(0\d{2,3})-)?(\d{7,8})(-(\d{3,}))?$`)
var chinaPhoneWithAreaRegex = regexp.MustCompile(`^(0\d{2,3})-(\d{7,8})(-\d{1,})?$`)

// IsChinaCardID 身份证
func IsChinaCardID(s string) bool {
	return chinaIDCard.MatchString(s)
}

// IsChinaMilitaryID 军官证或士兵证
func IsChinaMilitaryID(s string) bool {
	return chinaMilitaryID.MatchString(s)
}

// IsChinaHMTID 港澳通信证，台胞证
func IsChinaHMTID(s string) bool {
	return chinaHMTID.MatchString(s)
}

// IsChinaPassport 护照
func IsChinaPassport(s string) bool {
	return chinaPassport.MatchString(s)
}

// IsChinaMobile 大陆手机
func IsChinaMobile(s string) bool {
	return chinaMobileRegex.MatchString(s)
}

// IsChinaPhone 大陆电话（包括手机）
func IsChinaPhone(s string) bool {
	ok := IsChinaMobile(s)
	if ok {
		return ok
	}
	return chinaPhoneWithAreaRegex.MatchString(s)
}
