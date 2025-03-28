package handler

import (
	"context"
	"fmt"
	"log-flow/internal/domain/models"
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/workers"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type WebSocketManager struct {
	ProgressMessenger queue.LiveStatusQueue
	db                *gorm.DB
}

func NewWebSocketManager(liveProgressMessenger queue.LiveStatusQueue, db *gorm.DB) *WebSocketManager {
	return &WebSocketManager{
		ProgressMessenger: liveProgressMessenger,
		db:                db,
	}
}

func (wsm *WebSocketManager) LiveProgressLogs(c *websocket.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute) // 5 minutes
	defer cancel()

	jobID := c.Params("jobID") // Extract job ID from URL

	logReport, err := models.GetLogReportByJobID(wsm.db, jobID)
	if err != nil && err.Error() != "record not found" {
		c.WriteMessage(websocket.TextMessage, []byte("Database error occured while fetching job details"))
		log.Error("Error while fetching job details:", err)
		return
	}
	if logReport != nil {
		//send log report to client
		logReport := workers.LogLiveStats{
			JobID:              jobID,
			Progress:           100,
			Status:             "Completed",
			UniqueIPs:          logReport.UniqueIPs,
			InvalidLogs:        logReport.InvalidLogs,
			TotalLogsProcessed: logReport.TotalLogs,
			LogLevelCounts: map[string]int{
				"ERROR": logReport.ErrorCount,
				"WARN":  logReport.WarnCount,
				"INFO":  logReport.InfoCount,
			},
		}
		message, err := logReport.GetMessage()
		if err != nil {
			c.WriteMessage(websocket.TextMessage, []byte("Some error occured while parsing log report"))
		} else {
			c.WriteMessage(websocket.TextMessage, []byte(message))
		}
		return
	}

	//if log report not found, then check if job is registered
	job, err := models.GetJobByID(wsm.db, jobID)
	if err != nil {
		c.WriteMessage(websocket.TextMessage, []byte("Database error occured while fetching job details"))
		log.Error("Error while fetching job details:", err)
		return
	}
	if job == nil {
		c.WriteMessage(websocket.TextMessage, []byte("Job not registered(Invalid Job ID)"))
		return
	}
	if job.Attempts >= 3 && job.Succeeded == false {
		c.WriteMessage(websocket.TextMessage, []byte("Job had been attempted 3 times, but failed"))
		return
	}

	//Job registered, but log report not found. So, listen for progress messages

	workers.ActiveJobs.Store(jobID, true)
	defer workers.ActiveJobs.Delete(jobID) // Remove when client disconnects

	msgs, err := wsm.ProgressMessenger.WaitAndRecieveProgressMsgsQueue(ctx, jobID)
	if err != nil {
		if err == queue.ErrCtxCancelled {
			fmt.Println("Context cancelled: stopping queue wait for ", jobID)
			c.WriteMessage(websocket.TextMessage, []byte("Error: Context timeout"))
		} else {
			c.WriteMessage(websocket.TextMessage, []byte("Error: "+err.Error()))
			log.Warn("error while waiting for progress messages:", err)
		}
		return
	}

	for msg := range msgs {
		if err := c.WriteMessage(websocket.TextMessage, msg.Body); err != nil {
			c.WriteMessage(websocket.TextMessage, []byte("Error: "+err.Error()))
			log.Error("WebSocket send error:", err)
			break
		}
	}

	c.WriteMessage(websocket.TextMessage, []byte("Completed"))
}
