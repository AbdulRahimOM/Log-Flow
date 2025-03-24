package routes

import (
	"log-flow/internal/api/handler"
	"log-flow/internal/infrastructure/config"
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/infrastructure/storage"

	"github.com/gofiber/fiber/v2"
)

func MountRoutes(
	app *fiber.App,
	rabbitMQLogQueue queue.LogQueueSender,
	liveProgressMessenger queue.ProgressMessenger,
) {
	supaBase := storage.NewSupabaseStorage(config.Env.SupaBaseURL, config.Env.SupaBaseKey, config.Env.SupaBaseBucket)
	httpHandler := handler.NewHttpHandler(rabbitMQLogQueue, supaBase)
	websocketManager := handler.NewWebSocketManager(liveProgressMessenger)

	//http routes
	mountAuthRoutes(app, httpHandler)
	mountLogRoutes(app, httpHandler)

	//websocket routes
	mountWebSocketRoutes(app, websocketManager)

}
