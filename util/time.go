package util

import (
	"fmt"
	"time"
)

const (
	//TimeFormatYMDHIS 日期时间
	TimeFormatYMDHIS = "2006-01-02 15:04:05"
	//TimeFormatYMD 日期
	TimeFormatYMD = "2006-01-02"
	//TimeFormatyyyyMMddHHmmss 日期格式
	TimeFormatyyyyMMddHHmmss = "20060102150405"
)

//StringToTimestamp 字符串转时间戳
func StringToTimestamp(format string, dt string, loc *time.Location) (int64, error) {
	tm, err := time.ParseInLocation(format, dt, loc)
	if err != nil {
		return 0, err
	}
	return tm.Unix(), nil
}

//TimestampToString 时间戳转字符串
func TimestampToString(timestamp int64, format string) string {
	return time.Unix(timestamp, 0).Format(format)
}

//StringToDateTime 字符串转时间
func StringToDateTime(format, dt string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(format, dt, loc)
}

//TimeDurationToAboutStr 时间间隔
func TimeDurationToAboutStr(t time.Duration) string {
	if t.Hours() > float64(30*24) {
		return fmt.Sprintf("%d个月", int(t.Hours()/float64(30*24)))
	} else if t.Hours() > float64(7*24) {
		return fmt.Sprintf("%d个星期", int(t.Hours()/float64(7*24)))
	} else if t.Hours() > float64(24) {
		return fmt.Sprintf("%d天", int(t.Hours()/float64(24)))
	} else if t.Minutes() > float64(60) {
		return fmt.Sprintf("%d小时", int(t.Minutes()/float64(60)))
	} else if t.Seconds() > float64(60) {
		return fmt.Sprintf("%d分钟", int(t.Seconds()/float64(60)))
	}
	return fmt.Sprintf("%d秒", int(t.Seconds()))
}
