package bootstrap

import (
	"gin-api/application/middleware"
	"gin-api/pkg/app"
	"gin-api/route"
	"github.com/gin-gonic/gin"
	"html/template"
)

func SetupRoute(router *gin.Engine) {
	registerMiddleware(router)
	registerRouter(router)
}

//registerMiddleware 注册中间件
func registerMiddleware(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(middleware.MountApp())
	router.Use(middleware.Catch())
	router.Use(middleware.Cors())
	router.Use(middleware.AccessLog())
	router.Use(middleware.Translations())
}

//registerRouter 注册路由
func registerRouter(router *gin.Engine) {
	//处理 favicon.ico 导致的两次请求问题
	router.Any("favicon.ico", func(ctx *gin.Context) {
		return
	})

	//处理路由 404
	router.NoRoute(func(c *gin.Context) {
		panic(404)
	})

	//解析嵌入的html模板
	templ := template.Must(template.New("").ParseFS(app.EmbedTemplate(), app.TmplPattern()...))
	router.SetHTMLTemplate(templ)

	//注册各类路由
	route.RegisterApiRouter(router)
	route.RegisterAdminRouter(router)
	route.RegisterWebRouter(router)
	route.RegisterStaticRouter(router)
}
