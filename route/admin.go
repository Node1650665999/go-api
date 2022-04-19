package route

import (
	"github.com/gin-gonic/gin"
)

func RegisterAdminRouter(r *gin.Engine) *gin.Engine {
	admin := r.Group("/admin")
	{
		admin.Any("/foo", func(c *gin.Context) {
			c.String(200, "bar")
		})
	}

	return r
}

