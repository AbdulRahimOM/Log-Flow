package handler

import (
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/infrastructure/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/supabase-community/gotrue-go"
	"gorm.io/gorm"
)

type HttpHandler struct {
	fileStorage  storage.Storage
	logQueue     queue.LogQueueSender
	db           *gorm.DB
	supabaseAuth gotrue.Client
}

func NewHttpHandler(
	logQueue queue.LogQueueSender,
	storage storage.Storage,
	db *gorm.DB,
	supabaseAuth gotrue.Client,
) *HttpHandler {
	return &HttpHandler{
		fileStorage:  storage,
		logQueue:     logQueue,
		db:           db,
		supabaseAuth: supabaseAuth,
	}
}

func (h *HttpHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"msg": "ok",
	})
}
