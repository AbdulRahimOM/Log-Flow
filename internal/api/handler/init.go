package handler

import (
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/infrastructure/storage"
	"log-flow/internal/repo"

	"github.com/gofiber/fiber/v2"
)

type HttpHandler struct {
	fileStorage storage.Storage
	logQueue    queue.LogQueueSender
	repo        repo.Repository
}

func NewHttpHandler(logQueue queue.LogQueueSender, storage storage.Storage, repo repo.Repository) *HttpHandler {
	return &HttpHandler{
		fileStorage: storage,
		logQueue:    logQueue,
		repo:        repo,
	}
}

func (h *HttpHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"msg": "ok",
	})
}
