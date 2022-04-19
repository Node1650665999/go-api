package response

import (
	"github.com/gin-gonic/gin"
	"gin-api/application/errcode"
	"gin-api/pkg/app"
)

var ctx *gin.Context

//Json response json
func Json(code int, msg string, data interface{}) {
	if msg == "" {
		msg = errcode.CodeText(code)
	}
	ctx = app.GetInstance().C
	ctx.PureJSON(errcode.HttpCode(code), gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
	return
}

//JsonAbort response json and abort current request
func JsonAbort(code int, msg string, data interface{})  {
	Json(code, msg, data)
	ctx.Abort()
	return
}

