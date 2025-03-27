package handler

import (
	"fmt"
	"log-flow/internal/domain/response"
	"log-flow/internal/utils/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/supabase-community/gotrue-go/types"
)

func (h *HttpHandler) Login(c *fiber.Ctx) response.HandledResponse {
	req := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	if errResponse := validation.BindAndValidateJSONRequest(c, req); errResponse != nil {
		return errResponse
	}

	user, err := h.supabaseAuth.SignInWithEmailPassword(
		req.Email,
		req.Password,
	)
	if err != nil {
		return response.ErrorResponse(fiber.StatusInternalServerError, "LOGIN_FAILED", fmt.Errorf("Failed to login. %v", err))
	}

	return response.SuccessResponse(fiber.StatusOK, "LOGIN_SUCCESS", user)
}

func (h *HttpHandler) Register(c *fiber.Ctx) response.HandledResponse {
	req := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	if errResponse := validation.BindAndValidateJSONRequest(c, req); errResponse != nil {
		return errResponse
	}

	user, err := h.supabaseAuth.Signup(types.SignupRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return response.ErrorResponse(fiber.StatusInternalServerError, "SIGNUP_FAILED", fmt.Errorf("Failed to signup. %v", err))
	}

	return response.SuccessResponse(fiber.StatusOK, "SIGNUP_SUCCESS", user)
}
