package middleware

import (
	"fmt"
	"log-flow/internal/domain/response"
	"log-flow/internal/utils/helper"
	"log-flow/internal/utils/locals"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// To be called only after the user is authenticated by the AuthMiddleware.
// This checks if the user is the author of the job.
func JobAuthorCheck(c *fiber.Ctx) error {
	userID := locals.GetUserID(c)
	log.Debug("user id: ", userID)
	jobID := c.Params("jobID")
	if jobID == "" {
		log.Debug("Job ID is missing")
		return response.Response{
			HttpStatusCode: fiber.StatusBadRequest,
			Status:         false,
			ResponseCode:   "JOB_ID_MISSING",
			Error:          fmt.Errorf("Job ID is missing"),
		}.WriteToJSON(c)
	}
	if helper.IsUserIDMatchedWithJobID(userID, jobID) {
		return c.Next()
	} else {
		log.Debug("User is not authorized to access this job")
		return invalidAuthResponse(c, fmt.Errorf("User is not authorized to access this job"))
	}
}
