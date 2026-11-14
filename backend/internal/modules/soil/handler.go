package soil

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type SoilHandler struct {
	service *SoilService
}

func NewSoilHandler(service *SoilService) *SoilHandler {
	return &SoilHandler{service: service}
}

// POST /api/v1/soil
func (h *SoilHandler) Create(c fiber.Ctx) error {
	var req CreateSoilRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	soil, err := h.service.Create(c.Context(), req)
	if err != nil {
		return soilError(c, err)
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Soil record created successfully",
		"data":    soil,
	})
}

// GET /api/v1/soil/:id
func (h *SoilHandler) GetByID(c fiber.Ctx) error {
	soil, err := h.service.GetByID(
		c.Context(),
		c.Params("id"),
	)

	if err != nil {
		return soilError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    soil,
	})
}

// GET /api/v1/soil
func (h *SoilHandler) List(c fiber.Ctx) error {
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

	result, err := h.service.List(
		c.Context(),
		page,
		limit,
	)

	if err != nil {
		return soilError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// GET /api/v1/farms/:id/soil
func (h *SoilHandler) ListByFarm(c fiber.Ctx) error {
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

	result, err := h.service.ListByFarm(
		c.Context(),
		c.Params("id"),
		page,
		limit,
	)

	if err != nil {
		return soilError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// PUT /api/v1/soil/:id
func (h *SoilHandler) Update(c fiber.Ctx) error {
	var req UpdateSoilRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	soil, err := h.service.Update(
		c.Context(),
		c.Params("id"),
		req,
	)

	if err != nil {
		return soilError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Soil record updated successfully",
		"data":    soil,
	})
}

// DELETE /api/v1/soil/:id
func (h *SoilHandler) Delete(c fiber.Ctx) error {
	err := h.service.Delete(
		c.Context(),
		c.Params("id"),
	)

	if err != nil {
		return soilError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Soil record deleted successfully",
	})
}

// Centralized error handling.
func soilError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrSoilNotFound):
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Soil record not found",
		})

	case errors.Is(err, ErrFarmNotFound):
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Farm not found",
		})

	case errors.Is(err, ErrInvalidSoilInput):
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})

	case err.Error() == "invalid soil record ID",
		err.Error() == "invalid farm ID":
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
