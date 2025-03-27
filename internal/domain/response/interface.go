package response

import "github.com/gofiber/fiber/v2"

type HandledResponse interface {
	WriteToJSON(c *fiber.Ctx) error
}
