package workers

import (
	"encoding/json"
	"fmt"
	"log-flow/internal/domain/models"
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/infrastructure/storage"
	"sync"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

var ActiveJobs sync.Map

type Worker struct {
	db              *gorm.DB
	resultQueue     queue.LiveStatusQueue
	logQueue        queue.LogQueueReceiver
	storage         storage.Storage
	keyWordsToTrack []string
}

func NewWorkers(db *gorm.DB, storage storage.Storage, logQueue queue.LogQueueReceiver, progressQueue queue.LiveStatusQueue, keyWordsToTrack []string) *Worker {
	return &Worker{
		db:              db,
		resultQueue:     progressQueue,
		logQueue:        logQueue,
		storage:         storage,
		keyWordsToTrack: keyWordsToTrack,
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
			log.Errorf("❌ Failed to unmarshal message: %v", err)
			//marshalling errors are not supposed to be happen, and not meaningful to retry. Hence, directly sending to failed queue (for manual inspection, if required)
			w.logQueue.SendToFailedQueue(msg) 
			continue
		}

		err:=models.AddFailAttemptForJob(w.db, logMsg.JobID)
		if err != nil {
			log.Errorf("❌ Failed to add attempt for job in database: %v", err)
			w.logQueue.SentForRetry(msg)
			continue
		}

		logProcessor, err := NewLogProcessor(w.resultQueue, w.resultQueue, w.storage, w.db, w.keyWordsToTrack, logMsg.JobID)
		if err != nil {
			log.Errorf("❌ Failed to create log processor: %v", err)
			w.logQueue.SentForRetry(msg)
			continue
		}

		err = logProcessor.ProcessLogFile(logMsg)
		if err != nil {
			log.Errorf("❌ Failed to process log file: %v", err)
			w.logQueue.SentForRetry(msg)
			continue
		}

		log.Trace("✅ @Received message from RabbitMQ by worker(%d)..: %s\n", workerID, logMsg.FileURL)
	}
}
