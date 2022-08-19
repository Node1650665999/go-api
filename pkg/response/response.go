package response

import (
	"gin-api/application/errcode"
	"github.com/gin-gonic/gin"
)


//Json response json
func Json(ctx *gin.Context, code int, msg string, data interface{}) {
	if msg == "" {
		msg = errcode.CodeText(code)
	}
	ctx.PureJSON(errcode.HttpCode(code), gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
	return
}

//JsonAbort response json and abort current request
func JsonAbort(ctx *gin.Context, code int, msg string, data interface{})  {
	Json(ctx, code, msg, data)
	ctx.Abort()
	return
}

