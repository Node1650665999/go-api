package helpers

import (
	"github.com/spf13/cast"
)

// String 转转 s 为 string
func String(s interface{}) string {
	return cast.ToString(s)
}

// Int 转换 s 为 int
func Int(s interface{}) int {
	return cast.ToInt(s)
}

// Float64 转换为 s 为 float64
func Float64(s interface{}) float64 {
	return cast.ToFloat64(s)
}

// Int64 转换 s 为 Int64
func Int64(s interface{}) int64 {
	return cast.ToInt64(s)
}

// Uint 转换 s 为 Uint
func Uint(s interface{}) uint {
	return cast.ToUint(s)
}

// Bool 转换 s 为 Bool
func Bool(s interface{}) bool {
	return cast.ToBool(s)
}

//StringMapString 转换 s 为 map[string]string
func StringMapString(s interface{}) map[string]string  {
	return cast.ToStringMapString(s)
}




