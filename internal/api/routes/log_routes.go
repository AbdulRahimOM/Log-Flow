package routes

import (
	"log-flow/internal/api/handler"

	"github.com/gofiber/fiber/v2"
)

func mountLogRoutes(app *fiber.App, handler *handler.HttpHandler) {
	api := app.Group("/api")
	{
		api.Post("/upload-logs", handler.UploadLogs)
		// api.Get("/stats", handler.FetchStats)
		// api.Get("/stats/:jobId", handler.FetchStatsByJobId)
		// api.Get("/queue-status", handler.GetQueueStatus)
		// api.Get("/live-stats", handler.LiveStats)
	}
}
