package time

import (
	"fmt"
	"github.com/spf13/cast"
	"math"
	"time"
)

var (
	Loc,_                      = time.LoadLocation("Asia/Shanghai")
	LayoutSecond               = "2006-01-02 15:04:05"     //横杠区分 精度到秒
	LayoutSecondZone           = "2006-01-02 15:04:05 MST" //横杠区分 进度到秒 加时区
	LayoutDay                  = "2006-01-02"              //横杠区分 精度到天
	LayoutMonth                = "2006-01"                 //横杠区分 精度到月
	LayoutBackslashDay         = "2006/01/02"              //反斜杠区分 精度到天
	LayoutBackslashMonth       = "2006/01"                 //反斜杠区分 精度到月
	LayoutNumSecond            = "20060102150405"          //数字格式 精度到秒
	LayoutNumDay               = "20060102"                //数字格式 精度到天
	LayoutNumMonth             = "200601"                  //数字格式 精度到月
	LayoutYear                 = "2006"                    //数字格式 精度到年
	MonthBackslashDay          = "01/02"                   //反斜杠区分 精度到天
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


//格式化时间：20060102 -> 2006-01-02
func FormatNumDayToLayoutDay(value string, loc *time.Location) string {
	tm, _ := time.ParseInLocation(LayoutNumDay, value, loc)
	return tm.Format(LayoutDay)
}

//格式化时间：2006-01-02 -> 20060102
func FormatLayoutDayToNumDay(value string, loc *time.Location) string {
	tm, _ := time.ParseInLocation(LayoutDay, value, loc)
	return tm.Format(LayoutNumDay)
}

//格式化时间：20060102 -> 2006/01/02
func FormatNumDayToBackslashDay(value string, loc *time.Location) string {
	tm, _ := time.ParseInLocation(LayoutNumDay, value, loc)
	return tm.Format(LayoutBackslashDay)
}

//格式化时间：200601 -> 2006-01
func FormatNumMonthToLayoutMonth(value string, loc *time.Location) string {
	tm, _ := time.ParseInLocation(LayoutNumMonth, value, loc)
	return tm.Format(LayoutMonth)
}

//格式化时间：2006-01 -> 200601
func FormatLayoutMonthToNumMonth(value string, loc *time.Location) string {
	tm, _ := time.ParseInLocation(LayoutMonth, value, loc)
	return tm.Format(LayoutNumMonth)
}

//格式化时间：200601 -> 2006/01
func FormatNumMonthToBackslashMonth(value string, loc *time.Location) string {
	tm, _ := time.ParseInLocation(LayoutNumMonth, value, loc)
	return tm.Format(LayoutBackslashMonth)
}

//GetFirstDayOfMonth 返回传入的时间所在月份的第一天
func GetFirstDayOfMonth(timeObj time.Time) time.Time {
	return timeObj.AddDate(0, 0, -timeObj.Day()+1)
}

//GetLastDayOfMonth 返回传入的时间所在月份的最后一天
func GetLastDayOfMonth(timeObj time.Time) time.Time {
	return timeObj.AddDate(0, 1, -timeObj.Day())
}

//IsFirstDayOfMonth 判断传入的时间是否为所在月份第一天时间
func IsFirstDayOfMonth(timeObj time.Time) bool {
	firstDayTime := GetFirstDayOfMonth(timeObj)
	if timeObj.Day() == firstDayTime.Day() {
		return true
	}
	return false
}

//IsLastDayOfMonth 判断传入的时间是否为所在月份最后一天时间
func IsLastDayOfMonth(timeObj time.Time) bool {
	lastDayTime := GetLastDayOfMonth(timeObj)
	if timeObj.Day() == lastDayTime.Day() {
		return true
	}
	return false
}

//IsToday 判断传入的时间是否为今天
func IsToday(timeObj time.Time, loc *time.Location) bool {
	now := time.Now().In(loc)
	return timeObj.Year() == now.Year() &&
		timeObj.Month() == now.Month() &&
		timeObj.Day() == now.Day()
}

//GetTimeRangeOfMonth 获取某月1日0点时间和最后一日24点时间
func GetTimeRangeOfMonth(tmObj time.Time) (time.Time, time.Time) {
	//某月1日0点时间和24点时间
	firstDay := tmObj.AddDate(0, 0, -tmObj.Day()+1)
	firstDayZeroTime := time.Date(firstDay.Year(), firstDay.Month(), firstDay.Day(), 0, 0, 0, 0, time.Local)
	lastDay := firstDayZeroTime.AddDate(0, 1, -1)
	lastDayLastTime := time.Date(lastDay.Year(), lastDay.Month(), lastDay.Day(), 23, 59, 59, 0, time.Local)

	return firstDayZeroTime, lastDayLastTime
}

//GetFirstDayTimeRange 获取某月1日0点时间和24点时间
func GetFirstDayTimeRange(tmObj time.Time) (time.Time, time.Time) {
	//某月1日0点时间和24点时间
	firstDay := tmObj.AddDate(0, 0, -tmObj.Day()+1)
	firstDayZeroTime := time.Date(firstDay.Year(), firstDay.Month(), firstDay.Day(), 0, 0, 0, 0, time.Local)
	firstDayLastTime := time.Date(firstDay.Year(), firstDay.Month(), firstDay.Day(), 23, 59, 59, 0, time.Local)

	return firstDayZeroTime, firstDayLastTime
}

//GetLastDayTimeRange 获取某月最后一日0点时间和24点时间
func GetLastDayTimeRange(tmObj time.Time) (time.Time, time.Time) {
	//某月最后一日0点时间和24点时间
	lastDay := tmObj.AddDate(0, 1, -tmObj.Day())
	lastDayZeroTime := time.Date(lastDay.Year(), lastDay.Month(), lastDay.Day(), 0, 0, 0, 0, time.Local)
	lastDayLastTime := time.Date(lastDay.Year(), lastDay.Month(), lastDay.Day(), 23, 59, 59, 0, time.Local)
	return lastDayZeroTime, lastDayLastTime
}

//一年中的第几天
func getDayOfYear(tmObj time.Time) int {
	return tmObj.YearDay()
}

//一年中的第几周
func getWeekOfYear(tmObj time.Time) int {
	weekFloat := math.Ceil(float64(tmObj.YearDay()) / float64(7))
	return cast.ToInt(weekFloat)
}

//一周的第几天(周几)
func getWeekday(tmObj time.Time) int {
	weekDay := int(tmObj.Weekday())
	if weekDay == 0 {
		weekDay = 7
	}
	return weekDay
}

//ParseTime 解析前端传入的时间区间
func ParseTime(timeStart, timeEnd string, timeUnit int32, loc *time.Location)  (tmStart time.Time, tmEnd time.Time, err error) {
	if timeStart == "" || timeEnd == "" {
		switch timeUnit {
		case 1:
			//默认最近七天
			tmEnd = time.Now().In(loc)
			tmStart = tmEnd.AddDate(0, 0, -6)
		case 2:
			//默认最近两个月
			tmEnd = time.Now().In(loc)
			tmStart = tmEnd.AddDate(0, -1, 0)
		}
	}

	//某日(timeStart=timeEnd)
	if timeUnit == 2 && timeStart != "" && timeStart == timeEnd {
		//某日0点时间
		tmStart, err = time.ParseInLocation(LayoutDay, timeStart, loc)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		tmStart = time.Date(tmStart.Year(), tmStart.Month(), tmStart.Day(), 0, 0, 0, 0, loc)

		//某日24点时间
		//tmEnd = tmStart.AddDate(0, 0, 1)
		tmEnd = time.Date(tmStart.Year(), tmStart.Month(), tmStart.Day(), 23, 59, 59, 0, loc)
	}

	//某月(timeStart=timeEnd)
	if timeUnit == 4 && timeStart != "" && timeStart == timeEnd {
		tmObj, err := time.ParseInLocation(LayoutDay, timeStart, loc)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		//某月1日0点时间
		tmStart = tmObj.AddDate(0, 0, -tmObj.Day()+1)
		tmStart = time.Date(tmStart.Year(), tmStart.Month(), tmStart.Day(), 0, 0, 0, 0, loc)
		//某月最后一日24点时间
		tmEnd   = tmStart.AddDate(0, 1, -1)
		tmEnd   = time.Date(tmEnd.Year(), tmEnd.Month(), tmEnd.Day(), 23, 59, 59, 0, loc)
	}

	if timeStart != timeEnd {
		tmStart, err = time.ParseInLocation(LayoutDay, timeStart, loc)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		tmEnd, err = time.ParseInLocation(LayoutDay, timeEnd, loc)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	return tmStart, tmEnd, nil
}


