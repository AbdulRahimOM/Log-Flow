package server

import (
	"log-flow/internal/api/handler"
	"log-flow/internal/api/routes"
	"log-flow/internal/infrastructure/config"
	"log-flow/internal/infrastructure/db"
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/infrastructure/storage"
<<<<<<< HEAD
=======
	"log-flow/internal/repo"
>>>>>>> 96eb9961e0f7697b47c6ea0b2bdd61f4581f4779
	"log-flow/internal/workers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const (
	gigaByte     = 1024 * 1024 * 1024
	bodyLimit    = 2 * gigaByte
	numOfWorkers = 4
	appName      = "Log-Flow"
)

func InitializeServer() *fiber.App {

	app := fiber.New(fiber.Config{
		AppName:       appName,
		StrictRouting: true,
		BodyLimit:     bodyLimit,
	})
	app.Use(logger.New())

	//dependencies
	database := db.GetDB()
	fileStore := storage.NewSupabaseStorage(config.Env.SupaBaseURL, config.Env.SupaBaseKey, config.Env.SupaBaseBucket)
	logFileQueue := queue.InitLogQueue()
	liveProgressMessenger := queue.InitLiveStatusQueue()

	//workers
	workers := workers.NewWorkers(database, fileStore, logFileQueue, liveProgressMessenger)
	workers.StartMany(numOfWorkers)

	//handlers
	httpHandler := handler.NewHttpHandler(logFileQueue, fileStore, database)
	websocketManager := handler.NewWebSocketManager(liveProgressMessenger, database)

	//initialize routes
	routes.MountRoutes(app, httpHandler, websocketManager)

	return app
}
