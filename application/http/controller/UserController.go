package controller

import (
	"gin-api/application/errcode"
	"gin-api/application/http/logic"
	"gin-api/application/http/validate"
	"gin-api/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

type Person struct {
	Name string `form:"name" json:"name" binding:"required,oneof=red green"`
	//Email   string `form:"email" json:"email" binding:"required,email"`
	Address string `form:"address" json:"address" binding:"required"`
	Age     int    `form:"age" json:"age" binding:"required,gt=0,lt=120"`
}

func (u *UserController) Info(c *gin.Context) {
	var person Person
	valid, err := validate.BindAndValid(c, &person)
	if ! valid {
		response.Json(errcode.Fail, err.First(), nil)
		return
	}

	info := logic.NewUserLogic().Info()

	response.Json(
		errcode.Success,
		"success",
		gin.H{
			"p":     "",
			"m":     c.Request.Method,
			"num":   "123456",
			"phone": info.Phone,
			"html":  "<h1>world</h1>",
		},
	)
	return
}
