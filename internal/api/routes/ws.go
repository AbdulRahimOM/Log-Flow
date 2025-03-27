package routes

import (
	"log-flow/internal/api/handler"
	"log-flow/internal/api/middleware"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// func mountWebSocketRoutes(app, workers, liveProgressChannel)
func mountWebSocketRoutes(app *fiber.App, websocketManager *handler.WebSocketManager) {

	// WebSocket route
	app.Use(middleware.JobAuthorCheck)
	app.Get("/api/live-stats/:jobID", websocket.New(websocketManager.LiveProgressLogs))
}
