package rotation

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router fiber.Router, db *pgxpool.Pool) {
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	rotationGroup := router.Group("/rotation")

	rotationGroup.Post("/", handler.Generate)
	rotationGroup.Get("/:id", handler.GetPlan)
	rotationGroup.Get("/:id/results", handler.GetResults)

	router.Get("/farms/:id/rotations", handler.ListFarmPlans)
}
