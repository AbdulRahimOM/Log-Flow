package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimit(rate int) func(ctx *fiber.Ctx) error {
	return limiter.New(limiter.Config{
		Max:        rate,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			realIP := c.Get("X-Real-IP") //real ip, set by nginx
			if realIP != "" {
				return realIP
			}
			return c.Context().RemoteIP().String() //remote ip
		},
	})
}
