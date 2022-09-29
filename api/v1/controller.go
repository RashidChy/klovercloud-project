package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"project1/helper"
)

type ControllerInf interface {
	CreateUser(e echo.Context) error
	//GetUserList(e echo.Context) error
	//DeleteUser(e echo.Context) error
	//UpdateUser(e echo.Context) error
}
type ControllerInstanceStruct struct {
}

func Controller() ControllerInf {
	fmt.Println("controller")
	return new(ControllerInstanceStruct)
}

func (s ControllerInstanceStruct) CreateUser(e echo.Context) error {
	return helper.CreateUser().Execute(e)
}
