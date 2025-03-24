package routes

import (
	"log-flow/internal/api/handler"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// func mountWebSocketRoutes(app, workers, liveProgressChannel)
func mountWebSocketRoutes(app *fiber.App, websocketManager *handler.WebSocketManager) {

	// WebSocket route
	app.Get("/ws/:jobID", websocket.New(websocketManager.LiveProgressLogs))
}
