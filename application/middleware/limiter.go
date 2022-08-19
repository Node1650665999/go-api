package middleware

import (
	"gin-api/application/errcode"
	"gin-api/pkg/limiter"
	"gin-api/pkg/response"
	"github.com/gin-gonic/gin"
	"strings"
)

//LimitIp 全局限流器(限制单个用户访问能访问系统几次)
func LimitIp(format string, driver ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP() + ":" + format
		if err := Limiter(driver...).Check(key, format); err != nil {
			response.JsonAbort(c, errcode.TooManyRequests, err.Error(),nil)
			return
		}
		c.Next()
	}
}

//LimitRoute 针对某个接口限流(即所有人限制访问该接口总共几次)
func LimitRoute(format string, driver ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		key    := c.FullPath()
		if err := Limiter(driver...).Check(key, format); err != nil {
			response.JsonAbort(c, errcode.TooManyRequests, err.Error(), nil)
			return
		}
		c.Next()
	}
}

//LimitRouteAndIp 对某个ip访问某接口进行限流(即每人限制访问该接口几次)
func LimitRouteAndIp(format string, driver ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		key    := routeToKeyString(c.FullPath() + c.ClientIP())
		if err := Limiter(driver...).Check(key, format); err != nil {
			response.JsonAbort(c, errcode.TooManyRequests, err.Error(), nil)
			return
		}
		c.Next()
	}
}

// routeToKeyString 辅助方法，将 URL 中的 / 格式为 -
func routeToKeyString(routeName string) string {
	routeName = strings.TrimLeft(routeName, "/")
	routeName = strings.ReplaceAll(routeName, "/", "-")
	routeName = strings.ReplaceAll(routeName, ":", "_")
	return routeName
}

//Limiter 根据driver类型来实例化对应的限流器,
//同时对分布式限流器进行健康检查，一旦有问题则切换到单机版。
// driver=1 单机版限流
// driver=2 分布式限流
func Limiter(driver ...int) limiter.LimiterIfac  {
	if len(driver) == 0 || driver[0] == 1 {
		return  limiter.NewAloneBucket()
	}

	distribute := limiter.NewDistributeBucket()
	if distribute.HealthCheck() == false {
		return  limiter.NewAloneBucket()
	}

	return distribute
}

