package util

import (
	"fmt"
	"strconv"
	"time"
)

//IndentityData 身份证数据
type IndentityData struct {
	ProviceCode string
	CityCode    string
	Birthday    time.Time
	Gender      int
}

//CheckIdentity 只支持二代身份证
func CheckIdentity(val []byte) (*IndentityData, error) {
	// 校验位对应的规则。
	gb11643Map := []byte{'1', '0', 'x', '9', '8', '7', '6', '5', '4', '3', '2'}

	// 前17位号码对应的权值，为一个固定数组。可由gb11643_test.getWeight()计算得到。
	gb11643Weight := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

	if len(val) != 18 {
		return nil, fmt.Errorf("invalid id length")
	}

	sum := 0
	for i := 0; i < 17; i++ {
		sum += (gb11643Weight[i] * int((val[i] - '0')))
	}
	if val[17] == 'X' {
		val[17] = 'x'
	}

	if gb11643Map[sum%11] != val[17] {
		return nil, fmt.Errorf("invalid id")
	}

	return identityParse(val)
}

//identityParse 身份证解析
func identityParse(val []byte) (*IndentityData, error) {
	i, err := strconv.ParseInt(string(val[14:17]), 10, 64)
	if nil != err {
		return nil, fmt.Errorf("invalid id gender %s", string(val[14:17]))
	}
	t, err := time.Parse("20060102", string(val[6:14]))
	if nil != err {
		return nil, fmt.Errorf("invalide id birthday")
	}
	data := &IndentityData{
		ProviceCode: string(val[0:3]),
		CityCode:    string(val[3:6]),
		Birthday:    t,
		Gender:      int(i % 2),
	}
	return data, nil
}
