package soil

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router fiber.Router, db *pgxpool.Pool) {
	repository := NewSoilRepository(db)
	service := NewSoilService(repository)
	handler := NewSoilHandler(service)

	group := router.Group("/soil")

	group.Post("/", handler.Create)
	group.Get("/", handler.List)
	group.Get("/:id", handler.GetByID)
	group.Put("/:id", handler.Update)
	group.Delete("/:id", handler.Delete)

	// List soil records belonging to a specific farm.
	router.Get("/farms/:id/soil", handler.ListByFarm)
}
