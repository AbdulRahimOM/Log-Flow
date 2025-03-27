package handler

import (
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/infrastructure/storage"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HttpHandler struct {
	fileStorage storage.Storage
	logQueue    queue.LogQueueSender
	db *gorm.DB
}

func NewHttpHandler(logQueue queue.LogQueueSender, storage storage.Storage, db *gorm.DB) *HttpHandler {
	return &HttpHandler{
		fileStorage: storage,
		logQueue:    logQueue,
		db: db,
	}
}

func (h *HttpHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"msg": "ok",
	})
}
