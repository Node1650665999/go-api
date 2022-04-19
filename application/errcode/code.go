package errcode

import "net/http"

const (
	Success         = 200
	Fail            = 400
	No              = 404
	Unauthorized    = 401
	TooManyRequests = 429
	Fatal           = 500
)

var textMap = map[int]string{
	Success:         "成功",
	Fail:            "失败",
	No:              "路由不存在",
	Unauthorized:    "认证失败",
	TooManyRequests: "请求太频繁",
	Fatal:           "系统异常",
}

var httpMap = map[int]int{
	Success:         http.StatusOK,
	Fail:            http.StatusOK,
	No:              http.StatusOK,
	Fatal:           http.StatusInternalServerError,
}

//CodeText 获取错误码描述
func CodeText(code int) string {
	return textMap[code]
}

//HttpCode 获取错误码映射的http状态码
func HttpCode(code int) int {
	v, ok := httpMap[code]
	if ok == false {
		return http.StatusOK
	}
	return v
}
