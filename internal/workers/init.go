package workers

import (
	"encoding/json"
	"fmt"
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/infrastructure/storage"
	"sync"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

var ActiveJobs sync.Map

type Worker struct {
	db          *gorm.DB
	resultQueue queue.LiveStatusQueue
	logQueue    queue.LogQueueReceiver
	storage     storage.Storage
}

func NewWorkers(db *gorm.DB, storage storage.Storage, logQueue queue.LogQueueReceiver, progressQueue queue.LiveStatusQueue) *Worker {
	return &Worker{
		db:          db,
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

		logProcessor, err := NewLogProcessor(w.resultQueue, w.resultQueue, w.storage, w.db, logMsg.JobID)
		if err != nil {
			log.Error("❌ Failed to create log processor: %v", err) //need to implement retry logic
			continue
		}

		logProcessor.ProcessLogFile(logMsg)

		log.Trace("✅ @Received message from RabbitMQ by worker(%d)..: %s\n", workerID, logMsg.FileURL)
	}
}
