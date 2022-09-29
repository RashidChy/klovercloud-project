package router

import (
	"fmt"
	"github.com/labstack/echo/v4"
	v1 "project1/api/v1"
)

func v1Apis(group *echo.Group) {

	fmt.Println("router")

	p1 := group.Group("/p1")
	p1.POST("/register", v1.Controller().CreateUser)

}
