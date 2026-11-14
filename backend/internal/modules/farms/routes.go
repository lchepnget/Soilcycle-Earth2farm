package farms

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router fiber.Router, db *pgxpool.Pool) {
	repository := NewFarmRepository(db)
	service := NewFarmService(repository)
	handler := NewFarmHandler(service)

	group := router.Group("/farms")

	group.Post("/", handler.Create)
	group.Get("/", handler.List)
	group.Get("/:id", handler.GetByID)
	group.Put("/:id", handler.Update)
	group.Delete("/:id", handler.Delete)
}
