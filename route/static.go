package route

import (
	"github.com/gin-gonic/gin"
	"gin-api/pkg/app"
)

func RegisterStaticRouter(r *gin.Engine) *gin.Engine {

	r.StaticFS("/asset/public/", app.EmbedAsset())

	return r
}
