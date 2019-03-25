package util

import (
	"fmt"
	"math"
	"strings"
)

//Round 四舍五入
func Round(f float64, n int) float64 {
	p := math.Pow10(n)
	return math.Trunc(f*p+0.5) / p
}

const (
	hexStr      = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	tenThousand = 10000
	billion     = 100000000
)

// DecimalToAny 10进制到任意
func DecimalToAny(num, n int64) (string, error) {
	if n > 62 {
		return "", fmt.Errorf("over 62")
	}
	str := ""
	var r int64
	var rStr string
	for num != 0 {
		r = num % n
		rStr = string(hexStr[r])
		str = rStr + str
		num = num / n
	}
	return str, nil
}

func getHexStrIndex(str string) int64 {
	r := []rune(str)
	if r[0] > 96 && r[0] < 123 {
		return int64(r[0] - 87)
	}
	if r[0] > 64 && r[0] < 91 {
		return int64(r[0] - 29)
	}
	if r[0] > 47 && r[0] < 58 {
		return int64(r[0] - 48)
	}
	return -1
}

//AnyToDecimal 任意进制转10进制
func AnyToDecimal(str string, n int64) (int64, error) {
	var result int64
	arr := strings.Split(str, "")
	l := len(arr)
	for i, value := range arr {
		k := getHexStrIndex(value)
		if k < 0 {
			return 0, fmt.Errorf("over 62")
		}
		if k >= n {
			return 0, fmt.Errorf("over %d", n)
		}
		result += k * int64(math.Pow(float64(n), float64(l-i-1)))
	}
	return result, nil
}

//StatisticsInt 显示数字格式化以万、亿为单位转换
func StatisticsInt(num int64) *string {
	var str string
	if num >= billion {
		f := Round(float64(num)/float64(billion), 1)
		str = fmt.Sprintf("%.1f", f)
		strArr := strings.Split(str, ".")
		if len(strArr) < 2 {
			str = fmt.Sprintf("%s亿", strArr[0])
		} else {
			if strArr[1] == "0" {
				str = fmt.Sprintf("%s亿", strArr[0])
			} else {
				str = fmt.Sprintf("%.1f亿", f)
			}
		}
	} else if num >= tenThousand && num < billion {
		f := Round(float64(num)/float64(tenThousand), 1)
		str = fmt.Sprintf("%.1f", f)
		strArr := strings.Split(str, ".")
		if len(strArr) < 2 {
			str = fmt.Sprintf("%s万", strArr[0])
		} else {
			if strArr[1] == "0" {
				str = fmt.Sprintf("%s万", strArr[0])
			} else {
				str = fmt.Sprintf("%.1f万", f)
			}
		}
	} else if num > 0 && num < tenThousand {
		str = fmt.Sprintf("%d", num)
	} else if num <= 0 {
		str = fmt.Sprintf("%d", 0)
	}
	return &str
}

//Abs 绝对值
func Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

//StatisticsFloat 显示数字格式化以万、亿为单位转换
func StatisticsFloat(num float64) *string {
	var str string
	if num >= billion {
		f := Round(num/float64(billion), 1)
		str = fmt.Sprintf("%.1f", f)
		strArr := strings.Split(str, ".")
		if len(strArr) < 2 {
			str = fmt.Sprintf("%s亿", strArr[0])
		} else {
			if strArr[1] == "0" {
				str = fmt.Sprintf("%s亿", strArr[0])
			} else {
				str = fmt.Sprintf("%.1f亿", f)
			}
		}
	} else if num >= tenThousand && num < billion {
		f := Round(num/float64(tenThousand), 1)
		str = fmt.Sprintf("%.1f", f)
		strArr := strings.Split(str, ".")
		if len(strArr) < 2 {
			str = fmt.Sprintf("%s万", strArr[0])
		} else {
			if strArr[1] == "0" {
				str = fmt.Sprintf("%s万", strArr[0])
			} else {
				str = fmt.Sprintf("%.1f万", f)
			}
		}
	} else if num > 0 && num < tenThousand {
		f := Round(num, 1)
		str = fmt.Sprintf("%.1f", f)
		strArr := strings.Split(str, ".")
		if len(strArr) < 2 {
			str = fmt.Sprintf("%s", strArr[0])
		} else {
			if strArr[1] == "0" {
				str = fmt.Sprintf("%s", strArr[0])
			} else {
				str = fmt.Sprintf("%.1f", f)
			}
		}
	} else if num <= 0 {
		str = fmt.Sprintf("%d", 0)
	}
	return &str
}
