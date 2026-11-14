package farmers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type FarmerHandler struct {
	service *FarmerService
}

func NewFarmerHandler(service *FarmerService) *FarmerHandler {
	return &FarmerHandler{service: service}
}

func (h *FarmerHandler) Create(c fiber.Ctx) error {
	var req CreateFarmerRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	farmer, err := h.service.Create(c.Context(), req)
	if err != nil {
		return farmerError(c, err)
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Farmer registered successfully",
		"data":    farmer,
	})
}

func (h *FarmerHandler) GetByID(c fiber.Ctx) error {
	farmer, err := h.service.GetByID(
		c.Context(),
		c.Params("id"),
	)

	if err != nil {
		return farmerError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    farmer,
	})
}

func (h *FarmerHandler) List(c fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid page",
		})
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Limit must be between 1 and 100",
		})
	}

	result, err := h.service.List(c.Context(), page, limit)
	if err != nil {
		return farmerError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (h *FarmerHandler) Update(c fiber.Ctx) error {
	var req UpdateFarmerRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	farmer, err := h.service.Update(
		c.Context(),
		c.Params("id"),
		req,
	)

	if err != nil {
		return farmerError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Farmer updated successfully",
		"data":    farmer,
	})
}

func farmerError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrFarmerNotFound):
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Farmer not found",
		})

	case errors.Is(err, ErrFarmerExists):
		return c.Status(409).JSON(fiber.Map{
			"success": false,
			"message": "Email or phone already registered",
		})

	case errors.Is(err, ErrInvalidFarmerInput):
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})

	default:
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Internal server error",
		})
	}
}
