package farmers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router fiber.Router, db *pgxpool.Pool) {
	repository := NewFarmerRepository(db)
	service := NewFarmerService(repository)
	handler := NewFarmerHandler(service)

	group := router.Group("/farmers")

	group.Post("/", handler.Create)
	group.Get("/", handler.List)
	group.Get("/:id", handler.GetByID)
	group.Put("/:id", handler.Update)
}
