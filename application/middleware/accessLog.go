package middleware

import (
	"bytes"
	"fmt"
	"gin-api/pkg/app"
	"gin-api/pkg/helpers"
	"gin-api/pkg/logger"
	"github.com/gin-gonic/gin"
	"regexp"
	"time"
)

type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w AccessLogWriter) Write(p []byte) (int, error) {
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

//AccessLog 记录访问日志
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 静态资源不记录访问日志
		if helpers.IsMatchSubs(c.Request.URL.Path, app.LogExceptUri()...) {
			c.Next()
		} else {
			c.Request.ParseForm()

			//参数
			var param interface{}
			if c.GetHeader("Content-Type") == "application/json" {
				re, _     := regexp.Compile(`\n|\r|\s`)
				param      = re.ReplaceAllString(string(app.GetRequestBody()),"")
			} else {
				param = helpers.JsonEncode(c.Request.Form)
			}

			//重载 ResponseWriter 以便获取响应结果
			bodyWriter := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer    = bodyWriter

			beginTime  := time.Now().UnixNano() / 1e6

			c.Next()

			endTime := time.Now().UnixNano() / 1e6

			format := `
			begin_time : %v
			end_time   : %v
			cost_time  : %v ms
			header     : %v
			method     : %v
			url		   : %v
			params     : %v
			response   : %v
		`
			log := fmt.Sprintf(
				format,
				beginTime,
				endTime,
				endTime-beginTime,
				helpers.JsonEncode(c.Request.Header),
				c.Request.Method,
				app.GetFullUrl(),
				param,
				bodyWriter.body.String(),
			)
			logger.AccessLog(log, "")
		}
	}
}

