package handler

import (
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/infrastructure/storage"
)

type HttpHandler struct {
	fileStorage storage.Storage
	logQueue    queue.LogQueueSender
}

func NewHttpHandler(logQueue queue.LogQueueSender, storage storage.Storage) *HttpHandler {
	return &HttpHandler{
		fileStorage: storage,
		logQueue:    logQueue,
	}
}
