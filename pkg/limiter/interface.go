package limiter

import (
	"errors"
	"github.com/spf13/cast"
	"strings"
	"time"
)

type LimiterIfac interface {
	Check(key string, format string) error
}

var timeRule = map[string]time.Duration{
	"S": time.Second,
	"M": time.Minute,
	"H": time.Hour,
	"D": 24 * time.Hour,
}

//ParseFormat 解析 format, 获取单位时间(every)的限制数量(limit)
func ParseFormat(format string) (limit int, everyDuration time.Duration, err error) {
	sp := strings.Split(format, "-")
	if len(sp) != 2 {
		return 0, 0, errors.New("限流格式有误")
	}
	limit = cast.ToInt(sp[0])
	everyDuration = timeRule[sp[1]]
	return limit, everyDuration, nil
}