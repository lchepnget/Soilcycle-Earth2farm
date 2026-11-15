package auth

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Invalid request body",
		})
	}
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false, "message": "Email or phone and password are required",
		})
	}

	result, err := h.service.Login(c.Context(), req)
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false, "message": "Invalid email, phone or password",
		})
	case errors.Is(err, ErrAuthNotConfigured):
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"success": false, "message": "Authentication is not configured",
		})
	case err != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false, "message": "Login failed",
		})
	}

	return c.JSON(fiber.Map{
		"success": true, "message": "Login successful", "data": result,
	})
}
