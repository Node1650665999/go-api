package middleware

import (
	"errors"
	"gin-api/application/errcode"
	"gin-api/pkg/jwt"
	"gin-api/pkg/response"
	"github.com/gin-gonic/gin"
	"strings"
)

var (
	TokenCanNotEmpty   = errors.New("缺失token")
	TokenFormatError   = errors.New("请求头中 Authorization 格式有误")
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//分别从表单和header中查找token
		token := c.Request.FormValue("token")
		if len(token) == 0 {
			headerToken, err := getTokenFromHeader(c)
			if err != nil {
				response.JsonAbort(c, errcode.Unauthorized, err.Error(),nil)
				return
			}
			token = headerToken
		}

		customClaims, err := jwt.VerifyToken(token)
		if err != nil {
			response.JsonAbort(c, errcode.Unauthorized, err.Error(),nil)
			return
		}
		c.Set("custom_claims", customClaims)

		c.Next()
	}
}

//getTokenFromHeader 获取token, token 既支持通过参数传递,也支持通过 header 传递.
//header中格式为： Authorization:Bearer xxxxx
func getTokenFromHeader(c *gin.Context) (string,error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", TokenCanNotEmpty
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && strings.TrimSpace(parts[0]) == "Bearer") {
		return "", TokenFormatError
	}
	return strings.TrimSpace(parts[1]), nil
}
