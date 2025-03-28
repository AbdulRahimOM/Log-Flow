package helper

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

// Check if the file is a .log file
func IsValidLogFile(filename string) bool {
	return strings.HasSuffix(filename, ".log")
}

// Extract log level, log payload and IP from log entry
func ExtractLogDetails(logEntry string) (string, string, string, error) {
	pattern := `^\[(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)\]\s+(INFO|DEBUG|WARN|ERROR)\s+(.+)$`
	re := regexp.MustCompile(pattern)

	matches := re.FindStringSubmatch(logEntry)
	if len(matches) < 4 {
		return "", "", "", fmt.Errorf("invalid log format")
	}

	timestampStr := matches[1]
	logLevel := matches[2]
	message := matches[3]

	log.Trace("logLevel:", logLevel, "message:", message, " ")

	// Validate timestamp format
	if _, err := time.Parse(time.RFC3339, timestampStr); err != nil {
		return "", "", "", fmt.Errorf("invalid timestamp: %s", timestampStr)
	}

	ip := ""

	// Extract IP if JSON exists
	if strings.Contains(message, "{") {
		jsonPart := message[strings.Index(message, "{"):]
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(jsonPart), &data); err == nil {
			if ipVal, ok := data["ip"].(string); ok {
				ip = ipVal
			}
		}
	}

	log.Trace("ip:", ip)
	return logLevel, message, ip, nil
}
