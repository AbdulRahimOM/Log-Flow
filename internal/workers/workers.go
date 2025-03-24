package workers

import (
	"encoding/json"
	"fmt"
	"log-flow/internal/infrastructure/queue"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

var ActiveJobs sync.Map

type Worker struct {
	resultMessenger queue.ProgressMessenger
	logQueue        queue.LogQueueReceiver
}

func NewWorkers(logQueue queue.LogQueueReceiver, progressMessenger queue.ProgressMessenger) *Worker {
	return &Worker{
		resultMessenger: progressMessenger,
		logQueue:        logQueue,
	}
}

func (w *Worker) StartMany(count int) {
	fmt.Println("Starting all workers...")
	for i := 1; i <= count; i++ {
		go w.ProcessIncomingLogs(i)
	}
}

func (w *Worker) ProcessIncomingLogs(workerID int) {
	msgs, err := w.logQueue.RecieveLogFileDetails()
	if err != nil {
		log.Fatalf("Failed to consume messages from RabbitMQ: %v", err)
	}

	for msg := range msgs {

		log.Debug("log recieved by worker:", workerID)
		var logMsg queue.LogMessage
		if err := json.Unmarshal(msg.Body, &logMsg); err != nil {
			log.Error("❌ Failed to unmarshal message: %v", err)
			continue
		}

		w.processLogFile(logMsg)

		log.Debug("✅ @Received message from RabbitMQ by worker(%d)..: %s\n", workerID, logMsg.FileURL)
	}
}

func (w *Worker) processLogFile(msg queue.LogMessage) {
	queue, err := w.resultMessenger.StartQueue(msg.JobID)
	if err != nil {
		log.Error("Error starting queue: %v", err)
		return
	}

	for i := range 10 {
		time.Sleep(3 * time.Second) //to simulate processing

		// Checking if the any websocket client is listening (to this jobID)
		if _, exists := ActiveJobs.Load(msg.JobID); exists {
			queue.SendIntermediateResult(fmt.Sprintf("_Progress: %d0%%", i+1))
		} else {
			log.Debug("No WebSocket client is listening. Continuing processing...")
		}
	}

	queue.Delete()
	fmt.Println("✅ Done processing message:", msg)
}
