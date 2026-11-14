package earthdata

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetClimateData handles POST /earthdata/climate.
func (h *Handler) GetClimateData(c fiber.Ctx) error {
	var req EarthdataRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	ctx := c.Context()

	response, err := h.service.GetClimateData(
		ctx,
		req,
	)

	if err != nil {
		var validationError interface{ Unwrap() error }

		if errors.As(err, &validationError) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"success": false,
			"message": "Unable to retrieve climate data from NASA POWER",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Climate data retrieved successfully",
		"data":    response,
	})
}

// Keep context imported for compile-time compatibility with context.Context.
var _ context.Context
