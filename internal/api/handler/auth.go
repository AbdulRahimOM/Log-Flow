package handler

import "github.com/gofiber/fiber/v2"

func (h *HttpHandler) Login(c *fiber.Ctx) error {
	return c.SendString("Login")
}

func (h *HttpHandler) Register(c *fiber.Ctx) error {
	return c.SendString("Register")
}
