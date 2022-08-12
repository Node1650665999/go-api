package math

import (
	"github.com/shopspring/decimal"
	"strings"
)

//DecimalFormat 自定义保留places位小数
func DecimalFormat(a, b, places int64) string {
	if b == 0 {
		return "0." + strings.Repeat("0", int(places))
	}
	div := decimal.NewFromInt(b)
	dms := decimal.NewFromInt(a).Div(div).StringFixed(int32(places))
	return dms
}
