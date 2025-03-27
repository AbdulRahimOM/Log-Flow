package server

import (
	"log-flow/internal/api/handler"
	"log-flow/internal/api/routes"
	"log-flow/internal/infrastructure/config"
	"log-flow/internal/infrastructure/db"
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/infrastructure/storage"
	"log-flow/internal/workers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/supabase-community/gotrue-go"
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
	supabaseAuth := gotrue.New("mlvrrjjrhybrovqoijna",
		config.Env.SupaBaseKey)

	//workers
	workers := workers.NewWorkers(database, fileStore, logFileQueue, liveProgressMessenger, config.Env.LogConfig.Keywords)
	workers.StartMany(numOfWorkers)

	//handlers
	httpHandler := handler.NewHttpHandler(logFileQueue, fileStore, database, supabaseAuth)
	websocketManager := handler.NewWebSocketManager(liveProgressMessenger, database)

	//initialize routes
	routes.MountRoutes(app, httpHandler, websocketManager)

	return app
}
