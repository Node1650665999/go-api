package middleware

import (
	"gin-api/pkg/app"
	"github.com/gin-gonic/gin"
)

//MountApp 上下文信息到App上
func MountApp() gin.HandlerFunc {
	return func(c *gin.Context) {
		app.MountApp(c)
		c.Next()
	}
}

