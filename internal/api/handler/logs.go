package handler

import (
	"fmt"
	"log-flow/internal/domain/models"
	"log-flow/internal/domain/response"
	"log-flow/internal/infrastructure/queue"
	"log-flow/internal/utils/helper"
	"log-flow/internal/utils/locals"
	"time"

	_ "log-flow/internal/infrastructure/db"

	_ "log-flow/internal/infrastructure/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

const (
	uploadsDir = "./uploads"
)

func (h *HttpHandler) UploadLogs(c *fiber.Ctx) response.HandledResponse {
	file, err := c.FormFile("file")
	if err != nil {
		return response.ErrorResponse(fiber.StatusBadRequest, "INVALID_FILE", fmt.Errorf("Invalid file. %v", err))
	}

	if !helper.IsValidLogFile(file.Filename) {
		return response.ErrorResponse(fiber.StatusBadRequest, "NOT_SUPPORTED_FILE", fmt.Errorf("File type not supported. %v", err))
	}

	userID := locals.GetUserID(c)

	url, err := h.fileStorage.UploadFile(file)
	if err != nil {
		return response.ErrorResponse(fiber.StatusInternalServerError, "UPLOAD_FAILED", fmt.Errorf("Failed to upload file. %v", err))
	}
	log.Debug("File uploaded. URL: ", url)

	jobID, err := helper.GenerateUserIDPrefixedUUID(userID)
	if err != nil {
		return response.ErrorResponse(fiber.StatusInternalServerError, "UUID_ERROR", fmt.Errorf("Failed to generate UUID. %v", err))
	}

	logMsg := queue.LogMessage{
		JobID:   jobID.String(),
		FileURL: url,
	}

	err = h.logQueue.SendToQueue(logMsg)
	if err != nil {
		return response.ErrorResponse(fiber.StatusInternalServerError, "QUEUE_ERROR", fmt.Errorf("Failed to send to queue. %v", err))
	}

	job := models.Job{
		ID:         jobID,
		UserID:     userID,
		FileURL:    url,
		UploadedAt: time.Now(),
	}
	if err = job.Create(h.db); err != nil {
		return response.ErrorResponse(fiber.StatusInternalServerError, "DB_ERROR", fmt.Errorf("Failed to save job to db. %v", err))
	}

	return response.SuccessResponse(200, response.Success,
		map[string]any{
			"jobID": logMsg.JobID,
		},
	)
}

func (h *HttpHandler) FetchStatsByJobId(c *fiber.Ctx) response.HandledResponse {
	jobID := c.Params("jobID")

	job, err := models.GetLogReportByJobID(h.db, jobID)
	if err != nil {
		return response.ErrorResponse(fiber.StatusInternalServerError, "DB_ERROR", fmt.Errorf("Failed to get job details. %v", err))
	}
	if job == nil {
		return response.ErrorResponse(fiber.StatusNotFound, "JOB_NOT_FOUND", fmt.Errorf("Job not found."))
	}

	return response.SuccessResponse(200, response.Success, job)
}

func (h *HttpHandler) FetchStats(c *fiber.Ctx) response.HandledResponse {
	userID := locals.GetUserID(c)
	results, err := models.GetWholeLogReportsAggregate(h.db, userID)
	if err != nil {
		return response.DBErrorResponse(fmt.Errorf("Failed to get whole log reports aggregate. %v", err))
	}

	return response.SuccessResponse(200, response.Success, results)
}
