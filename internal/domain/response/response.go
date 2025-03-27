package response

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Response struct {
	HttpStatusCode int         `json:"-"`
	Status         bool        `json:"status"`
	ResponseCode   string      `json:"resp_code"`
	Error          error       `json:"-"` //will be marshalled to string when WriteToJSON is called
	Data           interface{} `json:"data,omitempty"`
}

type custError struct {
	Response
	Error string `json:"error"`
}

type ValidationErrorResponse struct {
	Status       bool           `json:"status"`
	ResponseCode string         `json:"resp_code"`
	Errors       []InvalidField `json:"errors"`
}

type InvalidField struct {
	FailedField string      `json:"field"`
	Tag         string      `json:"tag"`
	Value       interface{} `json:"value"`
}

func ErrorResponse(statusCode int, respcode string, err error) *Response {
	return &Response{
		HttpStatusCode: statusCode,
		Status:         false,
		ResponseCode:   respcode,
		Error:          err,
	}
}

func SuccessResponse(statusCode int, respcode string, data interface{}) *Response {
	return &Response{
		HttpStatusCode: statusCode,
		Status:         true,
		ResponseCode:   respcode,
		Data:           data,
	}
}

func CreatedResponse[T int | uuid.UUID](id T) *Response {
	return SuccessResponse(http.StatusCreated, Created, map[string]interface{}{
		"id": id,
	})
}

func InternalServerErrorResponse(err error) *Response {
	return ErrorResponse(http.StatusInternalServerError, InternalServerError, fmt.Errorf("internal server error: %w", err))
}

func DBErrorResponse(err error) *Response {
	return &Response{
		HttpStatusCode: 500,
		Status:         false,
		ResponseCode:   DbError,
		Error:          err,
	}
}

func InvalidURLParamResponse(param string, err error) *Response {
	return ErrorResponse(http.StatusBadRequest, InvalidURLParam, fmt.Errorf("error parsing %v from url: %w", param, err))
}

func BugResponse(err error) *Response { //Development purpose only
	return ErrorResponse(http.StatusInternalServerError, "BUG", fmt.Errorf("bug found, notify BE: %w", err))
}

func UnauthorizedResponse(err error) *Response {
	return ErrorResponse(http.StatusUnauthorized, Unauthorized, fmt.Errorf("unauthorized: %w", err))
}

func NotFoundResponse(itemName string) *Response {
	return ErrorResponse(http.StatusNotFound, NotFound, fmt.Errorf("not found: %v", itemName))
}

func (resp *Response) WriteToJSON(c *fiber.Ctx) error {

	if resp.Error == nil {
		return c.Status(resp.HttpStatusCode).JSON(resp)
	}

	newCustError := custError{
		Response: *resp,
	}
	if resp.Error != nil {
		newCustError.Error = resp.Error.Error()
	}

	return c.Status(resp.HttpStatusCode).JSON(newCustError)
}
