package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterWebRouter(r *gin.Engine) *gin.Engine {
	web := r.Group("/web/")
	{
		web.GET("/index", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"title": "hello world !!!",
			})
		})

		web.GET("/welcome", func(c *gin.Context) {
			//注意，所有的模板子目录都相对于 `application/http/view`
			c.HTML(http.StatusOK, "home/welcome.html", gin.H{
				"name": "tony",
			})
		})
	}

	return r
}
