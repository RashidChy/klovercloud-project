package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"project1/config"
	"project1/database"
	"project1/router"
)

func main() {

	e := echo.New()

	fmt.Println("Starting Application")

	err := config.InitEnVars()
	if err != nil {
		fmt.Println("[ERROR]: ", err.Error())
	}

	err = database.InitDBConnection()
	if err != nil {
		log.Fatal().Err(err).Msg("FATAL_EXIT")
	}

	fmt.Println("main")

	router.Routes(e)
	e.Logger.Fatal(e.Start(":" + config.ServerPortStr))
}
