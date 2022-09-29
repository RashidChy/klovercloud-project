package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"project1/config"
	"project1/router"
)

func main() {
	e := echo.New()

	fmt.Println("Starting Application")

	err := config.InitEnVars()
	if err != nil {
		fmt.Println("[ERROR]: ", err.Error())
	}

	fmt.Println("main")

	router.Routes(e)
	e.Logger.Fatal(e.Start(":" + config.ServerPortStr))
}
