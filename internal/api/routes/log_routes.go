package routes

import (
	"log-flow/internal/api/handler"
	"log-flow/internal/api/middleware"
	"log-flow/internal/infrastructure/config"

	"github.com/gofiber/fiber/v2"
)

func mountLogRoutes(app *fiber.App, handler *handler.HttpHandler) {
	api := app.Group("/api")
	app.Use(middleware.RateLimit(config.Env.GeneralRateLimit))
	api.Use(middleware.AuthMiddleware)
	{
		api.Post("/upload-logs", responseWrapper(handler.UploadLogs))
		api.Get("/stats", responseWrapper(handler.FetchStats))
		api.Get("/stats/:jobId", middleware.JobAuthorCheck, responseWrapper(handler.FetchStatsByJobId))
		// api.Get("/queue-status", handler.GetQueueStatus)
	}
}
