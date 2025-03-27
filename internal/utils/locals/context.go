package locals

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	UserIdKey = "userID"
)

func GetUserID(c *fiber.Ctx) uuid.UUID {
	userID, _ := uuid.Parse(c.Locals(UserIdKey).(string))
	return userID
}

func SetUserID(c *fiber.Ctx, userID string) {
	c.Locals(UserIdKey, userID)
}
