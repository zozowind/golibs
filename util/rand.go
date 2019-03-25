package util

import (
	"math/rand"
	"time"
)

// RandString 生成特定长度的随机字符串
func RandString(length int) string {
	rand.Seed(time.Now().Unix())
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	nonce := make([]byte, length)
	for i := 0; i < length; i++ {
		nonce[i] = byte(chars[rand.Intn(len(chars))])
	}
	return string(nonce)
}

// RandNumString 生成特定长度的数字字符串
func RandNumString(length int) string {
	rand.Seed(time.Now().Unix())
	chars := "0123456789"
	nonce := make([]byte, length)
	for i := 0; i < length; i++ {
		nonce[i] = byte(chars[rand.Intn(len(chars))])
	}
	return string(nonce)
}

// RandNum RandNum
func RandNum(min, max int) int {
	return min + rand.Intn(max-min+1)
}
