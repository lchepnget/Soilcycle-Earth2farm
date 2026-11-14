package crops

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type CropHandler struct {
	service *CropService
}

func NewCropHandler(service *CropService) *CropHandler {
	return &CropHandler{service: service}
}

func (h *CropHandler) Create(c fiber.Ctx) error {
	var req CreateCropRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	crop, err := h.service.Create(c.Context(), req)
	if err != nil {
		return respondError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Crop created successfully",
		"data":    crop,
	})
}

func (h *CropHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")

	crop, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return respondError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    crop,
	})
}

func (h *CropHandler) List(c fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid page parameter",
		})
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Limit must be between 1 and 100",
		})
	}

	result, err := h.service.List(
		c.Context(),
		c.Query("search"),
		c.Query("category"),
		page,
		limit,
	)

	if err != nil {
		return respondError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (h *CropHandler) Update(c fiber.Ctx) error {
	var req UpdateCropRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	crop, err := h.service.Update(
		c.Context(),
		c.Params("id"),
		req,
	)

	if err != nil {
		return respondError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Crop updated successfully",
		"data":    crop,
	})
}

func (h *CropHandler) Delete(c fiber.Ctx) error {
	err := h.service.Delete(c.Context(), c.Params("id"))
	if err != nil {
		return respondError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Crop deleted successfully",
	})
}

func (h *CropHandler) GetCharacteristics(c fiber.Ctx) error {
	result, err := h.service.GetCharacteristics(
		c.Context(),
		c.Params("id"),
	)

	if err != nil {
		return respondError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

func (h *CropHandler) UpsertCharacteristics(c fiber.Ctx) error {
	var req CropCharacteristicRequest

	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	result, err := h.service.UpsertCharacteristics(
		c.Context(),
		c.Params("id"),
		req,
	)

	if err != nil {
		return respondError(c, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Crop characteristics saved successfully",
		"data":    result,
	})
}

func respondError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrCropNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Crop not found",
		})

	case errors.Is(err, ErrDuplicateCrop):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "A crop with this name already exists",
		})

	default:
		// Log the actual error server-side in production.
		var validationError interface{ Error() string }

		if errors.As(err, &validationError) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "An internal server error occurred",
		})
	}
}
