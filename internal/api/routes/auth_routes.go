package routes

import (
	"log-flow/internal/api/handler"
	"log-flow/internal/api/middleware"
	"log-flow/internal/infrastructure/config"

	"github.com/gofiber/fiber/v2"
)

func mountAuthRoutes(app *fiber.App, handler *handler.HttpHandler) {
	auth := app.Group("/auth")
	app.Use(middleware.RateLimit(config.Env.AuthEndpointsRateLimit))
	{
		auth.Post("/login", responseWrapper(handler.Login))
		auth.Post("/register", responseWrapper(handler.Register))
	}
}
