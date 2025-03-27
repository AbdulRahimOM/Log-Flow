package routes

import (
	"log-flow/internal/api/handler"

	"github.com/gofiber/fiber/v2"
)

func mountAuthRoutes(app *fiber.App, handler *handler.HttpHandler) {
	auth := app.Group("/auth")
	{
		auth.Post("/login", responseWrapper(handler.Login))
		auth.Post("/register", responseWrapper(handler.Register))
	}
}
