package validation

import (
	"log-flow/internal/domain/response"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

const (
	bindingErrCode      = "BINDING_ERROR"
	validationErrCode   = "VALIDATION_ERROR"
	queryBindingErrCode = "URL_QUERY_BINDING_ERROR"
)

func bindErrResponse(err error) *response.Response {
	log.Debug("error parsing request:", err)
	return &response.Response{
		HttpStatusCode: http.StatusBadRequest,
		Status:         false,
		ResponseCode:   bindingErrCode,
		Error:          err,
	}
}

func validationErrResponse(err []response.InvalidField) *response.ValidationErrorResponse {
	log.Debug("error validating request:", err)
	return &response.ValidationErrorResponse{
		Status:       false,
		ResponseCode: validationErrCode,
		Errors:       err,
	}
}

// BindAndValidateRequest binds and validates the request.
// Req should be a pointer to the request struct.
func BindAndValidateJSONRequest(c *fiber.Ctx, req interface{}) response.HandledResponse {
	if err := c.BodyParser(req); err != nil {
		return bindErrResponse(err)
	}
	if err := validateJSONRequestDetailed(req); err != nil {
		return validationErrResponse(err)
	}

	return nil
}
