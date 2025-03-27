package routes

import (
	"log-flow/internal/api/handler"
	"log-flow/internal/domain/response"

	"github.com/gofiber/fiber/v2"
)

func MountRoutes(
	app *fiber.App,
	httpHandler *handler.HttpHandler,
	websocketManager *handler.WebSocketManager,
) {
	// health check
	app.Get("/health", httpHandler.HealthCheck)

	//http routes
	mountAuthRoutes(app, httpHandler)
	mountLogRoutes(app, httpHandler)

	//websocket routes
	mountWebSocketRoutes(app, websocketManager)

}

func responseWrapper(handlerFunc func(*fiber.Ctx) response.HandledResponse) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return handlerFunc(c).WriteToJSON(c)
	}
}
