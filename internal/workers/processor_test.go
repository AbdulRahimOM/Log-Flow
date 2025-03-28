package workers

import (
	"bytes"
	"io"
	"log-flow/internal/infrastructure/config"
	"log-flow/internal/utils/helper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	helper.SetFiberLogLevel(config.Env.AppSettings.LogLevel)
}

func TestProcessLogs(t *testing.T) {
	type wantA struct {
		logsProcessed int
		errorCount    int
		warnCount     int
		infoCount     int
		invalidLogs   int
		keywordCounts map[string]int
		uniqueIPs     int
	}

	tests := []struct {
		name           string
		logs           []string
		keywords       []string
		mockProcessLag bool
		lagMs          int
		want           wantA
	}{
		{
			name: "Basic log processing",
			logs: []string{
				`[2025-02-20T10:03:50Z] WARN High memory usage detected (78%)`,
				`[2025-02-20T10:05:23Z] ERROR Database timeout {"userId": 123, "ip": "192.168.1.1"}`,
				`[2025-02-20T10:09:30Z] ERROR Payment gateway unreachable {"error": "Connection refused", "orderId": 56789}`,
				`[2025-02-20T10:17:00Z] ERROR Failed to send email {"recipient": "user@example.com", "error": "SMTP timeout"}`,
				`[2025-02-20T10:02:30Z] INFO User login successful {"userId": 101, "ip": "192.168.1.5"}`,
			},
			keywords: []string{"timeout", "unreachable", "memory"},
			want: wantA{
				logsProcessed: 5,
				errorCount:    3,
				warnCount:     1,
				infoCount:     1,
				invalidLogs:   0,
				keywordCounts: map[string]int{
					"timeout":     2,
					"unreachable": 1,
					"memory":      1,
				},
				uniqueIPs: 2, // 192.168.1.1 and 192.168.1.5
			},
		},
		{
			name: "Invalid log format",
			logs: []string{
				"Invalid log format",
				`[2025-02-20T10:00:00Z] INFO Server started`,
				"Another invalid log",
			},
			keywords: []string{"started", "invalid"},
			want: wantA{
				logsProcessed: 3,
				errorCount:    0,
				warnCount:     0,
				infoCount:     1,
				invalidLogs:   2,
				keywordCounts: map[string]int{
					"started": 1,
					// "invalid": 0, // 0 because those logs are invalid
				},
				uniqueIPs: 0,
			},
		},
		{
			name: "JSON logs with IPs",
			logs: []string{
				`[2025-02-20T10:02:30Z] INFO User login successful {"userId": 101, "ip": "192.168.1.5"}`,
				`[2025-02-20T10:05:23Z] ERROR Database timeout {"userId": 123, "ip": "192.168.1.1"}`,
				`[2025-02-20T10:07:10Z] INFO API request received {"method": "GET", "endpoint": "/orders", "userId": 234}`,
			},
			keywords: []string{"login", "timeout"},
			want: wantA{
				logsProcessed: 3,
				errorCount:    1,
				warnCount:     0,
				infoCount:     2,
				invalidLogs:   0,
				keywordCounts: map[string]int{
					"login":   1,
					"timeout": 1,
				},
				uniqueIPs: 2,
			},
		},
		{
			name:     "Empty logs",
			logs:     []string{},
			keywords: []string{"error"},
			want: wantA{
				logsProcessed: 0,
				errorCount:    0,
				warnCount:     0,
				infoCount:     0,
				invalidLogs:   0,
				keywordCounts: map[string]int{},
				uniqueIPs:     0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuffer bytes.Buffer
			for _, log := range tt.logs {
				logBuffer.WriteString(log + "\n")
			}

			processor := &LogProcessor{
				keyWordsToTrack: tt.keywords,
				mockProcessLag:  tt.mockProcessLag,
				metrics: &LogMetrics{
					KeyWordsCount: make(map[string]int),
					UniqueIPs:     make(map[string]struct{}),
				},
			}

			// Process logs
			processor.processLogs(io.NopCloser(&logBuffer))

			// Assert results
			assert.Equal(t, tt.want.logsProcessed, processor.metrics.LogsProcessed, "logs processed count mismatch")
			assert.Equal(t, tt.want.errorCount, processor.metrics.ErrorCount, "error count mismatch")
			assert.Equal(t, tt.want.warnCount, processor.metrics.WarnCount, "warn count mismatch")
			assert.Equal(t, tt.want.infoCount, processor.metrics.InfoCount, "info count mismatch")
			assert.Equal(t, tt.want.invalidLogs, processor.metrics.InvalidLogs, "invalid logs count mismatch")
			assert.Equal(t, tt.want.uniqueIPs, len(processor.metrics.UniqueIPs), "unique IPs count mismatch")
			assert.EqualValues(t, tt.want.keywordCounts, processor.metrics.KeyWordsCount, "Keyword count mismatch")

		})
	}
}

func TestProcessLogsWithLogStreamError(t *testing.T) { //simulate error in log stream
	// Create a mock reader that returns an error
	errReader := &errorReader{err: io.ErrUnexpectedEOF}

	processor := &LogProcessor{
		keyWordsToTrack: []string{"test"},
		metrics: &LogMetrics{
			KeyWordsCount: make(map[string]int),
			UniqueIPs:     make(map[string]struct{}),
		},
	}

	processor.processLogs(io.NopCloser(errReader))
	assert.Equal(t, 0, processor.metrics.LogsProcessed, "no logs should be processed when reader errors")
	assert.Equal(t, 0, processor.metrics.ErrorCount, "no errors should be counted when reader errors")
	assert.Equal(t, 0, processor.metrics.WarnCount, "no warnings should be counted when reader errors")
	assert.Equal(t, 0, processor.metrics.InfoCount, "no info logs should be counted when reader errors")
	assert.Equal(t, 0, processor.metrics.InvalidLogs, "no invalid logs should be counted when reader errors")
	assert.Equal(t, 0, len(processor.metrics.UniqueIPs), "no unique IPs should be counted when reader errors")
	assert.EqualValues(t, 0, len(processor.metrics.KeyWordsCount), "no keywords should be counted when reader errors")

}

// Mock reader that always returns an error
type errorReader struct {
	err error
}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}
