package util

import (
	"fmt"
	"strconv"
	"strings"
)

//IsStringInSlice 检查是否在slice中
func IsStringInSlice(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

//StringToIntArray 字符串转整形数组
func StringToIntArray(str string, sep string) ([]int, error) {
	if len(str) == 0 {
		return make([]int, 0), nil
	}
	strs := strings.Split(str, sep)
	var err error
	var result = make([]int, len(strs))
	for i := 0; i < len(strs); i++ {
		if len(strs[i]) == 0 {
			continue
		}
		result[i], err = strconv.Atoi(strs[i])
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

//StringToIntArray 字符串转整形数组
func StringToStringArray(str string, sep string) ([]string, error) {
	if len(str) == 0 {
		return make([]string, 0), nil
	}
	strs := strings.Split(str, sep)
	var err error
	var result = make([]string, len(strs))
	for i := 0; i < len(strs); i++ {
		if len(strs[i]) == 0 {
			continue
		}
		result[i] = strs[i]
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

//IntArrayToString 将整型数组转成字符串
func IntArrayToString(arr []int, sep string) string {
	var arrStr = make([]string, len(arr))
	for i, v := range arr {
		arrStr[i] = strconv.Itoa(v)
	}
	return strings.Join(arrStr, sep)
}

//StringArrayToString 将字符串数组转成字符串
func StringArrayToString(arr []string, sep string) string {
	var arrStr = make([]string, len(arr))
	for i, v := range arr {
		arrStr[i] = v
	}
	return strings.Join(arrStr, sep)
}

//StringArrayToString 将字符串数组转成字符串
func StringArrayToStringIn(arr []string) string {
	s := strings.Replace(fmt.Sprint(arr), " ", "','", -1)
	s = strings.Replace(s, "[", "('", -1)
	s = strings.Replace(s, "]", "')", -1)
	return s
}
