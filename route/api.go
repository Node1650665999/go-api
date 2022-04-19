package route

import (
	"gin-api/application/errcode"
	"gin-api/application/http/controller"
	"gin-api/application/middleware"
	"gin-api/pkg/jwt"
	"gin-api/pkg/response"
	"github.com/gin-gonic/gin"
	"time"
)

func RegisterApiRouter(r *gin.Engine) *gin.Engine {
	api := r.Group("/api")

	api.Use(middleware.LimitRouteAndIp("500-H"))
	{
		api.Any("/foo", func(ctx *gin.Context) {
			response.Json(errcode.Success, "", nil)
			return
		})

		//panic test
		api.POST("/error", func(c *gin.Context) {
			//panic
			var slice = []int{1, 2, 3, 4, 5}
			slice[6] = 6
		})

		//token 相关
		tokenGroup := api.Group("/token").Use(middleware.LimitRouteAndIp("2-M"))
		{
			//生成token
			tokenGroup.Any("/get", func(ctx *gin.Context) {
				user := gin.H{"name": "tcl", "age": 30}
				token := jwt.GenerateToken(user, time.Hour)
				response.Json(errcode.Success, "success", token)
				return
			})

			//刷新token
			tokenGroup.Any("/fresh", middleware.JwtAuth(), func(ctx *gin.Context) {
				token, err := jwt.RefreshToken(ctx.Request.FormValue("token"), time.Hour)
				if err != nil {
					response.Json(errcode.Fail, err.Error(), nil)
				} else {
					response.Json(errcode.Success, "success", token)
				}
				return
			})
		}

		//带有版本号的接口
		v1 := api.Group("/v1").Use(middleware.LimitRoute("2-M"))
		{
			v1.Any("/info", middleware.JwtAuth(), controller.NewUserController().Info)
		}
	}

	return r
}
