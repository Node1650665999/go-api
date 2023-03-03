package math

import (
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"strings"
)

//DecimalFormat 自定义保留places位小数(四舍五入)
func DecimalFormat(a, b, places int64) string {
	if b == 0 {
		return "0." + strings.Repeat("0", int(places))
	}
	div := decimal.NewFromInt(b)
	dms := decimal.NewFromInt(a).Div(div).StringFixed(int32(places))
	return dms
}

//FloatFormat 自定义保留places位小数(不四舍五入)
func FloatFormat(a, b, places int64) string {
	if b == 0 {
		return "0." + strings.Repeat("0", int(places))
	}
	div := decimal.NewFromFloat(float64(b))
	dms, _ := decimal.NewFromFloat(float64(a)).Div(div).RoundFloor(2).Float64()
	return cast.ToString(dms)
}
