package rotation

import (
	"errors"
	"strconv"

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

// Generate creates alternative crop rotation plans.
func (h *Handler) Generate(c fiber.Ctx) error {
	var req CreateRotationRequest

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

	plan, results, err := h.service.GenerateRotation(
		c.Context(),
		req,
	)

	if err != nil {
		status := fiber.StatusInternalServerError

		switch err.Error() {
		case "farm not found":
			status = fiber.StatusNotFound
		case "no crops available for rotation planning":
			status = fiber.StatusUnprocessableEntity
		}
		if errors.Is(err, errInsufficientCrops) {
			status = fiber.StatusUnprocessableEntity
		}

		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Rotation plans generated successfully",
		"data": fiber.Map{
			"plan":    plan,
			"results": results,
		},
	})
}

// GetPlan retrieves a rotation plan.
func (h *Handler) GetPlan(c fiber.Ctx) error {
	planID := c.Params("id")

	plan, err := h.service.GetPlan(
		c.Context(),
		planID,
	)

	if err != nil {
		status := fiber.StatusInternalServerError

		if err.Error() == "invalid rotation plan ID" {
			status = fiber.StatusBadRequest
		}

		if err.Error() == "rotation plan not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    plan,
	})
}

// GetResults retrieves generated alternatives for a plan.
func (h *Handler) GetResults(c fiber.Ctx) error {
	planID := c.Params("id")

	results, err := h.service.GetPlanResults(
		c.Context(),
		planID,
	)

	if err != nil {
		status := fiber.StatusInternalServerError

		if err.Error() == "invalid rotation plan ID" {
			status = fiber.StatusBadRequest
		}

		if err.Error() == "rotation plan not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    results,
	})
}

// ListFarmPlans retrieves rotation plans for a farm.
func (h *Handler) ListFarmPlans(c fiber.Ctx) error {
	farmID := c.Params("id")

	limit := 20
	offset := 0

	if value := c.Query("limit"); value != "" {
		parsed, err := strconv.Atoi(value)

		if err != nil || parsed < 1 || parsed > 100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "limit must be between 1 and 100",
			})
		}

		limit = parsed
	}

	if value := c.Query("offset"); value != "" {
		parsed, err := strconv.Atoi(value)

		if err != nil || parsed < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "offset must be zero or greater",
			})
		}

		offset = parsed
	}

	plans, err := h.service.ListFarmPlans(
		c.Context(),
		farmID,
		limit,
		offset,
	)

	if err != nil {
		status := fiber.StatusInternalServerError

		if err.Error() == "invalid farm ID" {
			status = fiber.StatusBadRequest
		}

		if err.Error() == "farm not found" {
			status = fiber.StatusNotFound
		}

		return c.Status(status).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    plans,
	})
}
