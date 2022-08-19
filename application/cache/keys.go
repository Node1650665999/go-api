package cache

import "fmt"

var (
	Prefix       = "xxx"             //业务前缀
	UserInfo     = "user:info:%d"    //用户数据     user:info:{用户ID}
	TokenInfo    = "user:token:%s"   //token数据   user:token:{token值}
)

//FormatKey 格式化key，拼接业务前缀以及的参数
func FormatKey(key string, params ...interface{}) string {
	if len(params) <= 0 {
		return key
	}

	if len(Prefix) > 0 {
		key = Prefix + ":" + key
	}
	return fmt.Sprintf(key, params)
}
