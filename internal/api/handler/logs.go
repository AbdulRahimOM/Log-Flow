package handler

import (
	"fmt"
	"log-flow/internal/domain/response"
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/utils/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

const (
	uploadsDir = "./uploads"
)

func init() {
	err := helper.EnsureUploadsDir(uploadsDir)
	if err != nil {
		log.Fatal("Error ensuring uploads directory. ", err)
	}
}

func (h *HttpHandler) UploadLogs(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return response.ErrorResponse(fiber.StatusBadRequest, "INVALID_FILE", fmt.Errorf("Invalid file. %v", err)).WriteToJSON(c)
	}

	if !helper.IsValidLogFile(file.Filename) {
		return response.ErrorResponse(fiber.StatusBadRequest, "NOT_SUPPORTED_FILE", fmt.Errorf("File type not supported. %v", err)).WriteToJSON(c)
	}

	url, err := h.fileStorage.UploadFile(file)
	if err != nil {
		return response.ErrorResponse(fiber.StatusInternalServerError, "UPLOAD_FAILED", fmt.Errorf("Failed to upload file. %v", err)).WriteToJSON(c)
	}

	logMsg := queue.LogMessage{
		JobID:   uuid.New().String(),
		FileURL: url,
	}

	log.Debug("Log message: ", logMsg)
	err = h.logQueue.SendToQueue(logMsg)
	if err != nil {
		return response.ErrorResponse(fiber.StatusInternalServerError, "QUEUE_ERROR", fmt.Errorf("Failed to send to queue. %v", err)).WriteToJSON(c)
	}

	return response.SuccessResponse(200, response.Success, "File uploaded successfully").WriteToJSON(c)
}
