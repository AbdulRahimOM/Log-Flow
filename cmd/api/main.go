package main

import (
	"fmt"
	"log-flow/internal/infrastructure/config"
	"log-flow/internal/infrastructure/server"
	"log-flow/internal/utils/helper"

	"github.com/gofiber/fiber/v2/log"
)

func main() {
	helper.SetFiberLogLevel(config.Env.AppSettings.LogLevel)

	app := server.InitializeServer()
	err := app.Listen(fmt.Sprintf(":%s", config.Env.AppSettings.Port))
	if err != nil {
		log.Fatal("Couldn't start the server. Error: " + err.Error())
	}
}
