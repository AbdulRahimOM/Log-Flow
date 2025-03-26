package helper

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Check if the file is a .log file
func IsValidLogFile(filename string) bool {
	return strings.HasSuffix(filename, ".log")
}

// Extract log level and IP from log entry
func ExtractLogDetails(logEntry string) (string, string, error) {
	parts := strings.SplitN(logEntry, " ", 3)
	if len(parts) < 3 {
		return "UNKNOWN", "", fmt.Errorf("invalid log format: %s", logEntry)
	}

	level := parts[1]
	ip := ""

	// Extract IP if available in JSON payload
	if strings.Contains(logEntry, "{") {
		jsonPart := logEntry[strings.Index(logEntry, "{"):]
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(jsonPart), &data); err == nil {
			if ipVal, ok := data["ip"].(string); ok {
				ip = ipVal
			}
		}
	}

	return level, ip, nil
}
