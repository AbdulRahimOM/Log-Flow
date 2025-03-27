package workers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
)

type LogLiveStats struct {
	JobID              string         `json:"jobID"`
	Progress           float64        `json:"progressInPercentage"`
	UniqueIPs          int            `json:"uniqueIPs"`
	InvalidLogs        int            `json:"invalidLogs"`
	TotalLogsProcessed int            `json:"totalLogsProcessed"`
	Status             string         `json:"status"` //Started, In Progress, Completed
	LogLevelCounts     map[string]int `json:"logLevelCounts"`
	KeyWordCounts      map[string]int `json:"keyWordCounts"`
}

func (lls *LogLiveStats) GetMessage() (string, error) {
	message, err := json.Marshal(lls)
	if err != nil {
		log.Errorf("Error marshalling live stats: %v", err)
		return "", fmt.Errorf("Error marshalling live stats: %v", err)
	}
	return string(message), nil
}

type LogMetrics struct {
	ProcessedSize int64
	LogsProcessed int
	InvalidLogs   int
	ErrorCount    int
	WarnCount     int
	InfoCount     int
	UniqueIPs      map[string]struct{}
	KeyWordsCount map[string]int
}
