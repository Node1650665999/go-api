package helpers

import (
	"bytes"
	"strings"
)

//StrJoin 用来拼接字符串
func StrJoin(args ...string) string {
	var buf bytes.Buffer
	for _, arg := range args {
		buf.WriteString(arg)
	}
	return buf.String()
}

//IsMatchSubs 检查是否包含子串中的任意一个
func IsMatchSubs(str string, subs... string) bool {
	isMatch := false
	for _, sub := range subs {
		if strings.Contains(str, sub) {
			isMatch = true
		}
	}

	return isMatch
}


