package time

import (
	"fmt"
	"github.com/spf13/cast"
	"strings"
	"time"
)

var (
	LayoutSecond         = "2006-01-02 15:04:05"     //横杠区分 精度到秒
	LayoutSecondZone     = "2006-01-02 15:04:05 MST" //横杠区分 进度到秒 加时区
	LayoutDay            = "2006-01-02"              //横杠区分 精度到天
	LayoutMonth          = "2006-01"                 //横杠区分 精度到月
	LayoutBackslashDay   = "2006/01/02"              //反斜杠区分 精度到天
	LayoutBackslashMonth = "2006/01"                 //反斜杠区分 精度到月
	LayoutNumSecond      = "20060102150405"          //数字格式 精度到秒
	LayoutNumDay         = "20060102"                //数字格式 精度到天
	LayoutNumMonth       = "200601"                  //数字格式 精度到月
	LayoutYear           = "2006"                    //数字格式 精度到年
	MonthBackslashDay    = "01/02"                   //反斜杠区分 精度到天
)

//CurrentDate 返回当前日期
func CurrentDate() string {
	return time.Now().In(time.Local).Format(LayoutSecond)
}

//CurrentTimestamp 返回当前时间戳
func CurrentTimestamp() int64 {
	return time.Now().In(time.Local).Unix()
}

//TimestampToDate 时间戳转日期
func TimestampToDate(timestamp int64, layout string) string {
	if layout == "" {
		layout = LayoutSecond
	}
	return time.Unix(timestamp, 0).In(time.Local).Format(layout)
}

//MsToDate 毫秒时间戳转日期
func MsToDate(ms int64, layout string) string {
	if layout == "" {
		layout = LayoutSecond
	}
	return time.Unix(0, ms*int64(time.Millisecond)).Format(layout)
}

//DateToTimestamp 日期转时间戳
func DateToTimestamp(date string, layout string) int64 {
	timeObj, _ := time.ParseInLocation(layout, date, time.Local)
	return timeObj.Unix()
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

//FormatNumDayToLayoutDay 格式化时间：20060102 -> 2006-01-02
func FormatNumDayToLayoutDay(value string) string {
	tm, _ := time.ParseInLocation(LayoutNumDay, value, time.Local)
	return tm.Format(LayoutDay)
}

//FormatLayoutDayToNumDay 格式化时间：2006-01-02 -> 20060102
func FormatLayoutDayToNumDay(value string) string {
	tm, _ := time.ParseInLocation(LayoutDay, value, time.Local)
	return tm.Format(LayoutNumDay)
}

//FormatNumDayToBackslashDay 格式化时间：20060102 -> 2006/01/02
func FormatNumDayToBackslashDay(value string) string {
	tm, _ := time.ParseInLocation(LayoutNumDay, value, time.Local)
	return tm.Format(LayoutBackslashDay)
}

//FormatNumMonthToLayoutMonth 格式化时间：200601 -> 2006-01
func FormatNumMonthToLayoutMonth(value string) string {
	tm, _ := time.ParseInLocation(LayoutNumMonth, value, time.Local)
	return tm.Format(LayoutMonth)
}

//FormatLayoutMonthToNumMonth  格式化时间：2006-01 -> 200601
func FormatLayoutMonthToNumMonth(value string) string {
	tm, _ := time.ParseInLocation(LayoutMonth, value, time.Local)
	return tm.Format(LayoutNumMonth)
}

//FormatNumMonthToBackslashMonth 格式化时间：200601 -> 2006/01
func FormatNumMonthToBackslashMonth(value string) string {
	tm, _ := time.ParseInLocation(LayoutNumMonth, value, time.Local)
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
func IsToday(timeObj time.Time) bool {
	now := time.Now().In(time.Local)
	return timeObj.Year() == now.Year() &&
		timeObj.Month() == now.Month() &&
		timeObj.Day() == now.Day()
}

//GetTimeRangeOfDay 获取某日0点时间和24点时间
func GetTimeRangeOfDay(tmObj time.Time) (time.Time, time.Time) {
	firstSecondObj := time.Date(tmObj.Year(), tmObj.Month(), tmObj.Day(), 0, 0, 0, 0, time.Local)
	lastSecondObj := time.Date(tmObj.Year(), tmObj.Month(), tmObj.Day(), 23, 59, 59, 0, time.Local)
	return firstSecondObj, lastSecondObj
}

//GetTimeRangeOfMonth 获取某月1日0点时间和最后一日24点时间
func GetTimeRangeOfMonth(tmObj time.Time) (time.Time, time.Time) {
	firstDay := tmObj.AddDate(0, 0, -tmObj.Day()+1)

	firstSecondObj := time.Date(firstDay.Year(), firstDay.Month(), firstDay.Day(), 0, 0, 0, 0, time.Local)

	lastDay := firstSecondObj.AddDate(0, 1, -1)
	lastSecondObj := time.Date(lastDay.Year(), lastDay.Month(), lastDay.Day(), 23, 59, 59, 0, time.Local)
	return firstSecondObj, lastSecondObj
}

//GetTimeRangeOfYear 获取某年第一天0点时间和最后一天24点时间
func GetTimeRangeOfYear(tmObj time.Time) (time.Time, time.Time) {
	firstDay := tmObj.AddDate(0, -int(tmObj.Month())+1, -tmObj.Day()+1)

	firstSecondObj := time.Date(firstDay.Year(), firstDay.Month(), firstDay.Day(), 0, 0, 0, 0, time.Local)

	lastDay := firstSecondObj.AddDate(1, 0, -1)
	lastSecondObj := time.Date(lastDay.Year(), lastDay.Month(), lastDay.Day(), 23, 59, 59, 0, time.Local)
	return firstSecondObj, lastSecondObj
}

//GetTimeRangeOfWeek 返回某日所在的周一和周日的时间
func GetTimeRangeOfWeek(timeObj time.Time) (int64, int64) {
	offset := int(time.Monday - timeObj.Weekday())
	if offset > 0 {
		offset = -6
	}
	timeStart := time.Date(timeObj.Year(), timeObj.Month(), timeObj.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	timeEnd := timeStart.AddDate(0, 0, 6)
	timeEnd = time.Date(timeEnd.Year(), timeEnd.Month(), timeEnd.Day(), 23, 59, 59, 0, time.Local)

	ts := timeStart.Format("20060102")
	te := timeEnd.Format("20060102")

	return cast.ToInt64(ts), cast.ToInt64(te)
}

//GetWeekStartAndWeekEnd 传入一年中的第几周,返回这周的开始时间与结束时间
//refer:https://stackoverflow.com/questions/52300644/date-range-by-week-number-golang
func GetWeekStartAndWeekEnd(year, week int) (int64, int64) {
	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.Local)

	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)
	return GetTimeRangeOfWeek(t)
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

//GetDayOfYear 返回当前时间为一年中的第几天
func GetDayOfYear(tmObj time.Time) int {
	return tmObj.YearDay()
}

//GetWeekOfYear 返回当前时间为一年中的第几周
func GetWeekOfYear(tmObj time.Time) (year, week int) {
	return tmObj.ISOWeek()
}

//GetWeekday 返回当前时间为一周的第几天(周几)
func GetWeekday(tmObj time.Time) int {
	weekDay := int(tmObj.Weekday())
	if weekDay == 0 {
		weekDay = 7
	}
	return weekDay
}

//SecondFormat 将秒格式化为几分几秒(xx:xx)
func SecondFormat(second int64) string {
	m := second / 60
	s := second % 60
	return fmt.Sprintf("%02d:%02d", m, s)
}

//ConvertSecond 将几分几秒(xx:xx) 转换为秒
func ConvertSecond(str string) int64 {
	arr := strings.Split(str, ":")
	if len(arr) != 2 {
		return 0
	}
	return cast.ToInt64(arr[0])*60 + cast.ToInt64(arr[1])
}

//SecondFormatDate 转换秒为时分秒
func SecondFormatDate(second int64) string {
	d := second / 60 / 60 / 24
	h := (second / 60 / 60) % 24
	m := (second / 60) % 60
	s := second % 60
	return fmt.Sprintf("%02d天%02d时%02d分%02d秒", d, h, m, s)
}

//FillDate 返回两个日期间的日期列表
func FillDate(ts, te string) []string {
	dateList := make([]string, 0)
	layout := "20060102"

	timeStar, _ := time.ParseInLocation(layout, ts, time.Local)
	timeEnd, _ := time.ParseInLocation(layout, te, time.Local)

	for timeEnd.After(timeStar) {
		dateList = append(dateList, timeStar.Format(layout))
		timeStar = timeStar.AddDate(0, 0, 1)
	}
	dateList = append(dateList, timeEnd.Format(layout))
	return dateList
}

//GetMaxPersistDays 统计一组日期中的最大连续天数
func GetMaxPersistDays(days []string) int {
	length := len(days)
	if length <= 0 {
		return 0
	}

	maxPersist := 1
	temp := 1

	for i := 1; i < length; i++ {
		today, _ := time.ParseInLocation(LayoutNumDay, days[i], time.Local)
		prevDay, _ := time.ParseInLocation(LayoutNumDay, days[i-1], time.Local)
		subTime := today.Sub(prevDay)
		if subTime.Hours()/24 == 1 {
			temp++
		} else {
			if temp > maxPersist {
				maxPersist = temp
			}
			temp = 1
		}
	}

	if maxPersist > temp {
		return maxPersist
	}

	return temp
}

// DateChunk 将日期范围切分为一组日期数组
func DateChunk(start, end time.Time, chunkSize int) [][]time.Time {
	chunks := make([][]time.Time, 0)
	chunk := make([]time.Time, 0)

	// 循环遍历开始日期到结束日期之间的每一天
	for start.Before(end) || start.Equal(end) {
		chunk = append(chunk, start)

		// 如果切块大小已满或已达到结束日期，则将当前切块添加到切块数组中，并创建新的切块
		if len(chunk) == chunkSize || start.Equal(end) {
			chunks = append(chunks, chunk)
			chunk = make([]time.Time, 0)
		}

		start = start.AddDate(0, 0, 1)
	}

	return chunks
}
