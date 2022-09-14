// Package middlewares Gin 中间件
package middleware

import (
	"errors"
	"fmt"
	"gin-api/application/errcode"
	"gin-api/pkg/hash"
	"gin-api/pkg/request"
	"gin-api/pkg/response"
	"github.com/gin-gonic/gin"
	"sort"
	"strings"
)


const APP_SECRET = "aPL3UW45zRcAA5shjNdFbBg44o0WRB"

func Sign() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := signCheck(c)
		if err != nil {
			response.JsonAbort(c, errcode.Sign, err.Error(),nil)
		}
		c.Next()
	}
}

func signCheck(c *gin.Context) error {
    //从请求头中获取token和sign
	token := c.Request.Header.Get("token")
	sign  := c.Request.Header.Get("sign")

	//获取全部参数
	any, err := request.Input(c, "")
	if err != nil {
		return err
	}

	allParams, ok := any.(map[string]interface{})

	if ok == false {
		return errors.New("params type not is map")
	}

	//校验参数
	if len(allParams) <= 0 {
		return errors.New("params is empty")
	}
	if len(sign) <= 0 {
		return errors.New("[sign] 缺失")
	}

	//对key进行 ascii 排序
	if len(token) > 0 {
		allParams["token"] = token
	}
	keys := []string{}
	for key := range allParams {
		//排除要签名的字段
		if key == "sign" {
			continue
		}

		keys = append(keys, key)
	}
	sort.Strings(keys)

	//使用 & 拼接字段
	str := ""
	for _, key := range keys {
		//非空参数才参与签名
		if val, _ := allParams[key]; val != nil {
			str += fmt.Sprintf("&%s=%v", key, allParams[key])
		}
	}
	str = strings.Trim(str, "&") + "&" + APP_SECRET

	if sign != strings.ToUpper(hash.HashBySha1(str)) {
		return errors.New("[sign] 失败")
	}

	return nil
}
