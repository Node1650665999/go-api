package middleware

import (
	"fmt"
	"gin-api/application/errcode"
	"gin-api/pkg/logger"
	"github.com/gin-gonic/gin"
	//"gin-api/pkg/email"
	"gin-api/pkg/response"
)

//Catch 捕获路由异常
func Catch() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//处理路由不存在的情况
				code := errcode.Fatal
				v, isInt := err.(int)
				if isInt && v == 404 {
					code = errcode.No
				} else {
					str := fmt.Sprintf("\nException: %+v", err)
					logger.RuntimeLog(str)

					//邮件报警
					//todo
					/*e := email.SendMail("系统异常报警", str, []string{})
					if e != nil {
						logger.RuntimeLog(fmt.Sprintf("邮件发送失败: %s", e.Error()))
					}*/
				}

				//客户端响应
				response.Json(
					c,
					code,
					"",
					nil,
				)
				c.Abort()
			}
		}()
		c.Next()
	}
}
