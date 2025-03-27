package helper

import (
	"fmt"

	"github.com/google/uuid"
)

// GenerateNewUUIDWithPrefixFromString takes a user ID as a string (UUID format)
// and generates a new UUID that has the same first four hex digits.
func GenerateUserIDPrefixedUUID(userID uuid.UUID) (uuid.UUID, error) {
	// Combine the first 8 hex characters of the user ID with the rest of the new UUID
	modifiedUUIDStr := userID.String()[:8] + uuid.New().String()[8:]

	finalUUID, err := uuid.Parse(modifiedUUIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID: %v", err)
	}

	return finalUUID, nil
}

// Compare the first 8 hex characters of the user ID with the first 8 hex characters of the job ID
// and return true if they match.
func IsUserIDMatchedWithJobID(userID uuid.UUID, jobID string) bool {
	if len(jobID) < 8 {
		return false
	}
	return userID.String()[:8] == jobID[:8]
}
