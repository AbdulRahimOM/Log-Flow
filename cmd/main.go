package main

import (
	"fmt"
	"log-flow/internal/api/routes"
	"log-flow/internal/infrastructure/config"
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/workers"
	_ "log-flow/internal/workers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const (
	gigaByte  = 1024 * 1024 * 1024
	bodyLimit = 2 * gigaByte
)

func healthCheck(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"msg": "ok",
	})
}

func main() {
	app := fiber.New(fiber.Config{
		AppName:       "Log-Flow",
		StrictRouting: true,
		BodyLimit:     2 * 1024 * 1024 * 1024,
	})
	app.Use(logger.New())

	// health check
	app.Get("/health", healthCheck)

	log.SetLevel(log.LevelInfo)

	rabbitMqConfig := queue.RabbitMQConfig{
		Host:     "localhost",
		Port:     "5672",
		User:     "guest",
		Password: "guest",
	}

	logFileQueue, err := queue.NewRabbitMQLogQueue(rabbitMqConfig, "log_files_exchange", "log_processing_queue")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	liveProgressMessenger, err := queue.NewResultChan(rabbitMqConfig)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ live progress channel: %v", err)
	}

	workers.NewWorkers(logFileQueue, liveProgressMessenger).StartMany(4)

	routes.MountRoutes(app, logFileQueue, liveProgressMessenger)

	err = app.Listen(fmt.Sprintf(":%s", config.Env.AppSettings.Port))
	if err != nil {
		log.Fatal("Couldn't start the server. Error: " + err.Error())
	}
}
