package logic

import (
	"gin-api/application/http/model"
)

type UserLogic struct {}

func NewUserLogic() *UserLogic {
	return &UserLogic{}
}

func (ul *UserLogic) Info() *model.User {
	return  model.NewUser().Info()
}

