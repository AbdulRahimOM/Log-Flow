package jwttoken

import (
	"fmt"
	"log-flow/internal/infrastructure/config"

	"github.com/golang-jwt/jwt/v5"
)

// ValidateTokenAndGetUserID validates the JWT token and returns the User ID
func ValidateTokenAndGetUserID(tokenStr string) (string, error) {

	// Parse the JWT token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Env.SupaBaseJwtSecret), nil
	})
	if err != nil {
		return "", fmt.Errorf("Error parsing token: %v", err)
	}
	if !token.Valid {
		return "", fmt.Errorf("Invalid token")
	}

	// Extract claims from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Invalid token claims")
	}

	// Extract `sub` (Considering it as User ID, as it is unique)
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("Subject (sub) claim missing or invalid")
	}

	return userID, nil
}
