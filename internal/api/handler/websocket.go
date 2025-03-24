package handler

import (
	"context"
	"fmt"
	"log"
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/workers"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketManager struct {
	ProgressMessenger queue.ProgressMessenger
}

func NewWebSocketManager(liveProgressMessenger queue.ProgressMessenger) *WebSocketManager {
	return &WebSocketManager{
		ProgressMessenger: liveProgressMessenger,
	}
}

func (wsm *WebSocketManager) LiveProgressLogs(c *websocket.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute) // 5 minutes
	defer cancel()

	jobID := c.Params("jobID") // Extract job ID from URL
	fmt.Println("Job ID:", jobID)

	//first check in database. -->to do later

	c.WriteMessage(websocket.TextMessage, []byte("Your job is in queue. Please wait..."))

	workers.ActiveJobs.Store(jobID, true)
	defer workers.ActiveJobs.Delete(jobID) // Remove when client disconnects

	msgs, err := wsm.ProgressMessenger.WaitAndRecieveProgressMsgsQueue(ctx, jobID)
	if err != nil {
		if err == queue.ErrCtxCancelled {
			c.WriteMessage(websocket.TextMessage, []byte("Error: Context timeout"))
		} else {
			c.WriteMessage(websocket.TextMessage, []byte("Error: "+err.Error()))
			log.Println("consumeResults error:", err)
		}
		return
	}

	for msg := range msgs {
		fmt.Println("Message received:", string(msg.Body))
		if err := c.WriteMessage(websocket.TextMessage, msg.Body); err != nil {
			c.WriteMessage(websocket.TextMessage, []byte("Error: "+err.Error()))
			log.Println("WebSocket send error:", err)
			break
		}
	}

	//get final data from db: -->to do later

	c.WriteMessage(websocket.TextMessage, []byte("Final result: bla bla bal"))
}
