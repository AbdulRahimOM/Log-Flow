package workers

import (
	"bufio"
	"fmt"
	"io"
	"log-flow/internal/domain/models"
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/infrastructure/storage"
	"log-flow/internal/utils/helper"
	"math"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LogProcessor struct {
	liveStatusQueue queue.LiveStatusQueueSession
	storage         storage.Storage
	db              *gorm.DB
	metrics         *LogMetrics
	stopChan        chan struct{}
	mutex           sync.Mutex
	jobID           string
	totalSize       int64
	status          string
}

func NewLogProcessor(
	progressMessenger queue.LiveStatusQueue,
	progressQueue queue.LiveStatusQueue,
	storage storage.Storage,
	db *gorm.DB,
	jobID string,
) (*LogProcessor, error) {

	queueSession, err := progressMessenger.StartQueue(jobID)
	if err != nil {
		return nil, fmt.Errorf("Error starting live stats queue: %v", err)
	}
	return &LogProcessor{
		liveStatusQueue: *queueSession,
		storage:         storage,
		metrics:         &LogMetrics{UniqueIPs: make(map[string]struct{})},
		stopChan:        make(chan struct{}),
		jobID:           jobID,
		db:              db,
	}, nil
}

func (lp *LogProcessor) ProcessLogFile(logMessage queue.LogMessage) {
	fileURL := logMessage.FileURL

	logStream, err := lp.storage.StreamLogs(fileURL)
	if err != nil {
		log.Errorf("Failed to stream logs: %v", err)
		//need to implement retry logic
		return
	}
	defer logStream.Close()

	lp.totalSize, _ = lp.storage.GetFileSize(fileURL)
	go lp.sendLiveUpdates()

	lp.processLogs(logStream)

	close(lp.stopChan)

	err = lp.SaveFinalMetrics()
	if err != nil {
		log.Errorf("Error saving final metrics: %v", err)
	}

	lp.liveStatusQueue.Delete()
}

func (lp *LogProcessor) processLogs(logStream io.ReadCloser) {
	scanner := bufio.NewScanner(logStream)
	for scanner.Scan() {
		logEntry := scanner.Text()

		lp.mutex.Lock()
		logLevel, ip, parseErr := helper.ExtractLogDetails(logEntry)
		if parseErr != nil {
			log.Debug("Parsing Error: %v", parseErr)
			lp.metrics.InvalidLogs++
		}

		lp.metrics.LogsProcessed++
		lp.metrics.ProcessedSize += int64(len(logEntry) + 1)

		switch logLevel {
		case "ERROR":
			lp.metrics.ErrorCount++
		case "WARN":
			lp.metrics.WarnCount++
		case "INFO":
			lp.metrics.InfoCount++
		}

		if ip != "" {
			if _, ok := lp.metrics.UniqueIPs[ip]; !ok {
				lp.metrics.UniqueIPs[ip] = struct{}{}
			}
		}

		lp.mutex.Unlock()
		// time.Sleep(1 * time.Second)
		time.Sleep(500 * time.Millisecond)
	}

	if err := scanner.Err(); err != nil {
		log.Errorf("Error reading log stream: %v", err)
	}
}

func (lp *LogProcessor) sendLiveUpdates() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if _, ok := ActiveJobs.Load(lp.jobID); !ok {
				log.Trace("No websockets listening for job: %v", lp.jobID)
				continue
			}
			lp.mutex.Lock()
			progress := lp.calculateProgress()
			stats := LogLiveStats{
				JobID:              lp.jobID,
				Progress:           progress,
				UniqueIPs:          len(lp.metrics.UniqueIPs),
				InvalidLogs:        lp.metrics.InvalidLogs,
				TotalLogsProcessed: lp.metrics.LogsProcessed,
				LogLevelCounts: map[string]int{
					"error": lp.metrics.ErrorCount,
					"warn":  lp.metrics.WarnCount,
					"info":  lp.metrics.InfoCount,
				},
				Status: "In Progress",
			}
			lp.mutex.Unlock()

			strMessage, err := stats.GetMessage()
			if err != nil {
				log.Errorf("Error marshalling live stats: %v", err)
				continue
			}

			lp.liveStatusQueue.SendIntermediateResult(strMessage)

		case <-lp.stopChan:
			if _, ok := ActiveJobs.Load(lp.jobID); !ok {
				log.Trace("No websockets listening for job: %v", lp.jobID)
				return
			}
			lp.mutex.Lock()
			stats := LogLiveStats{
				JobID:              lp.jobID,
				Progress:           100,
				UniqueIPs:          len(lp.metrics.UniqueIPs),
				InvalidLogs:        lp.metrics.InvalidLogs,
				TotalLogsProcessed: lp.metrics.LogsProcessed,
				LogLevelCounts: map[string]int{
					"error": lp.metrics.ErrorCount,
					"warn":  lp.metrics.WarnCount,
					"info":  lp.metrics.InfoCount,
				},
				Status: "Completed",
			}
			lp.mutex.Unlock()

			strMessage, err := stats.GetMessage()
			if err != nil {
				log.Errorf("Error marshalling live stats: %v", err)
				continue
			}

			lp.liveStatusQueue.SendIntermediateResult(strMessage)
			return
		}
	}
}

func (lp *LogProcessor) calculateProgress() float64 {
	return math.Round(min((float64(lp.metrics.ProcessedSize)/float64(lp.totalSize))*100, 100))
}

func (lp *LogProcessor) SaveFinalMetrics() error {
	lp.mutex.Lock()
	defer lp.mutex.Unlock()
	logReport := models.LogReport{
		JobID:       uuid.MustParse(lp.jobID),
		TotalLogs:   lp.metrics.LogsProcessed,
		ErrorCount:  lp.metrics.ErrorCount,
		WarnCount:   lp.metrics.WarnCount,
		InfoCount:   lp.metrics.InfoCount,
		UniqueIPs:   len(lp.metrics.UniqueIPs),
		InvalidLogs: lp.metrics.InvalidLogs,
		CreatedAt:   time.Now(),
	}

	err := logReport.Create(lp.db)
	if err != nil {
		return fmt.Errorf("Error saving log report: %v", err)
	}

	return nil
}
