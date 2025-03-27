package middleware

import (
	"fmt"
	"log-flow/internal/domain/response"
	jwttoken "log-flow/internal/utils/jwt"
	"log-flow/internal/utils/locals"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return invalidAuthResponse(c, fmt.Errorf("Missing token"))
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := jwttoken.ValidateTokenAndGetUserID(tokenStr)
	if err != nil {
		return invalidAuthResponse(c, err)
	}

	//set user id in context
	locals.SetUserID(c, userID)

	return c.Next()
}

func invalidAuthResponse(c *fiber.Ctx, err error) error {
	return response.UnauthorizedResponse(err).WriteToJSON(c)
}
