package farms

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type FarmHandler struct {
	service *FarmService
}

func NewFarmHandler(service *FarmService) *FarmHandler {
	return &FarmHandler{service: service}
}

func (h *FarmHandler) Create(c fiber.Ctx) error {
	var req CreateFarmRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	farm, err := h.service.Create(c.Context(), req)
	if err != nil {
		return farmError(c, err)
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Farm registered successfully",
		"data":    farm,
	})
}

func (h *FarmHandler) GetByID(c fiber.Ctx) error {
	farm, err := h.service.GetByID(
		c.Context(),
		c.Params("id"),
	)

	if err != nil {
		return farmError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    farm,
	})
}

func (h *FarmHandler) List(c fiber.Ctx) error {
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
		c.Query("user_id"),
		page,
		limit,
	)

	if err != nil {
		return farmError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (h *FarmHandler) Update(c fiber.Ctx) error {
	var req UpdateFarmRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Temporary development approach.
	// Replace this with the authenticated farmer ID from
	// JWT middleware before public deployment.
	userID := c.Query("user_id")

	if userID == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "user_id is required",
		})
	}

	farm, err := h.service.Update(
		c.Context(),
		c.Params("id"),
		userID,
		req,
	)

	if err != nil {
		return farmError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Farm updated successfully",
		"data":    farm,
	})
}

func (h *FarmHandler) Delete(c fiber.Ctx) error {
	// Temporary development approach.
	// Replace with authenticated farmer identity.
	userID := c.Query("user_id")

	if userID == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "user_id is required",
		})
	}

	err := h.service.Delete(
		c.Context(),
		c.Params("id"),
		userID,
	)

	if err != nil {
		return farmError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Farm deleted successfully",
	})
}

func farmError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrFarmNotFound):
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Farm not found",
		})

	case errors.Is(err, ErrFarmerMissing):
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Farmer does not exist",
		})

	case errors.Is(err, ErrInvalidFarmInput):
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
