package workers

import (
	"encoding/json"
	"fmt"
)

type LogLiveStats struct {
	JobID             string         `json:"job_id"`
	Progress          float64        `json:"progressInPercentage"`
	UniqueIPs         int            `json:"unique_ips"`
	ParsingErrorCount int            `json:"parsingErrorCount"`
	Status            string         `json:"status"` //Started, InProgress, Completed
	LogLevelCounts    map[string]int `json:"logLevelCounts"`
}

func (lls *LogLiveStats) GetMessage() (string, error) {
	message, err := json.Marshal(lls)
	if err != nil {
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
	UniqueIPs     map[string]struct{}
}
