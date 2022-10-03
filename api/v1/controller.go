package v1

import (
	"github.com/labstack/echo/v4"
	"project1/helper"
)

type ControllerInf interface {
	CreateUser(e echo.Context) error
	LoginUser(e echo.Context) error
	//GetUserList(e echo.Context) error
	//DeleteUser(e echo.Context) error
	//UpdateUser(e echo.Context) error
}
type ControllerInstanceStruct struct {
}

func Controller() ControllerInf {
	return new(ControllerInstanceStruct)
}

func (s ControllerInstanceStruct) CreateUser(e echo.Context) error {
	return helper.CreateUser().Execute(e)
}

func (s ControllerInstanceStruct) LoginUser(e echo.Context) error {
	return helper.LoginUser().Execute(e)
}
