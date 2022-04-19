package app

import (
	"bytes"
	"gin-api/pkg/config"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

//MountApp 挂载 gin.Context 到 app 上
func MountApp(c *gin.Context) {
	scheme := "http://"
	if c.Request.TLS != nil {
		scheme = "https://"
	}

	app.C 		    = c
	app.host        = scheme   + c.Request.Host
	app.fullUrl     = app.host + c.Request.RequestURI
	//由于 request body 不能读取两次, 为了后续能继续读取 body，因此将body数据回写至 Request.Body
	body,_          := c.GetRawData()
	app.requestBody = body
	c.Request.Body  = ioutil.NopCloser(bytes.NewBuffer(body))
}

//GetFullUrl 获取当前请求完整的url
func GetFullUrl() string {
	return app.fullUrl
}

//GetRequestBody 获取当前请求body
func GetRequestBody() []byte {
	return app.requestBody
}

//HttpPort 返回监听的http端口
func HttpPort() int {
	return config.GetInt("app.port")
}


