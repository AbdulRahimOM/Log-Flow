package workers

import (
	"encoding/json"
	"fmt"
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/infrastructure/storage"
	"log-flow/internal/repo"
	"sync"

	"github.com/gofiber/fiber/v2/log"
)

var ActiveJobs sync.Map

type Worker struct {
	repo        repo.Repository
	resultQueue queue.LiveStatusQueue
	logQueue    queue.LogQueueReceiver
	storage     storage.Storage
}

func NewWorkers(repo repo.Repository, storage storage.Storage, logQueue queue.LogQueueReceiver, progressQueue queue.LiveStatusQueue) *Worker {
	return &Worker{
		repo:        repo,
		resultQueue: progressQueue,
		logQueue:    logQueue,
		storage:     storage,
	}
}

func (w *Worker) StartMany(count int) {
	fmt.Println("Starting all workers...")
	for i := 1; i <= count; i++ {
		go w.start(i)
	}
}

func (w *Worker) start(workerID int) {
	msgs, err := w.logQueue.RecieveLogFileDetails()
	if err != nil {
		log.Fatalf("Failed to consume messages from RabbitMQ: %v", err)
	}

	for msg := range msgs {

		log.Debug("✅log recieved by worker:", workerID)
		var logMsg queue.LogMessage
		if err := json.Unmarshal(msg.Body, &logMsg); err != nil {
			log.Error("❌ Failed to unmarshal message: %v", err)
			continue
		}

		logProcessor, err := NewLogProcessor(w.resultQueue, w.storage, w.resultQueue, logMsg.JobID)
		if err != nil {
			log.Error("❌ Failed to create log processor: %v", err) //need to implement retry logic
			continue
		}

		logProcessor.ProcessLogFile(logMsg)

		// log.Debugf("✅ @Received message from RabbitMQ by worker(%d)..: %s\n", workerID, logMsg.FileURL)
	}
}
