package time

import (
	"fmt"
	"time"
)

//CurrentDate 返回当前日期
func CurrentDate() string {
	time := time.Now() //time.Time
	return TimeFormat(time)
}

//CurrentTimestamp 返回当前时间戳
func CurrentTimestamp() int64 {
	return time.Now().Unix() //获取当前时间
}

//TimestampInZone 基于时区获取 unix 时间戳
// eg. timezone = Asia/Shanghai
func TimestampInZone(timezone string) int64 {
	l, _ := time.LoadLocation(timezone)
	return time.Now().In(l).Unix() //获取当前时间
}

//Timestamp2Date 将时间戳转换为日期
func Timestamp2Date(timestamp int64) string {
	time := time.Unix(timestamp, 0) //time.Time
	return TimeFormat(time)
}

//Date2Timestamp 将日期转换时间戳
func Date2Timestamp(date string) int64 {
	format := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")                    //重要：获取时区
	timeObj, err := time.ParseInLocation(format, date, loc) //指定日期转当地日期对象 类型为 time.Time
	if err != nil {
		panic(err)
	}
	return timeObj.Unix()
}

// FormatDate 将time对象格式化为日期
func FormatDate(time time.Time) string {
	year := time.Year()     //年
	month := time.Month()   //月
	day := time.Day()       //日
	hour := time.Hour()     //小时
	minute := time.Minute() //分钟
	second := time.Second() //秒
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d\n", year, month, day, hour, minute, second)
}

// TimeFormat 将time对象格式化为日期
func TimeFormat(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

// SetInterval 定时执行 fn
func SetInterval(d time.Duration, fn func(args ...interface{}), args ...interface{}) chan bool {
	ticker := time.NewTicker(d)
	stopChan := make(chan bool)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fn(args...)
			case <-stopChan:
				break
			}
		}
	}()

	return stopChan
}

// SetTimeout 超时执行 fn
func SetTimeout(d time.Duration, fn func(args ...interface{}), args ...interface{}) {
	stopChan := make(chan bool)
	timer := time.NewTimer(d)
	go func() {
		select {
		case <-timer.C:
			fn(args)
		case <-stopChan:
			break
		}
		timer.Stop()
	}()
}
