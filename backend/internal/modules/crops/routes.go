package crops

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router fiber.Router, db *pgxpool.Pool) {
	repository := NewCropRepository(db)
	service := NewCropService(repository)
	handler := NewCropHandler(service)

	crops := router.Group("/crops")

	crops.Get("/", handler.List)
	crops.Get("/:id/characteristics", handler.GetCharacteristics)
	crops.Put("/:id/characteristics", handler.UpsertCharacteristics)

	crops.Get("/:id", handler.GetByID)
	crops.Post("/", handler.Create)
	crops.Put("/:id", handler.Update)
	crops.Delete("/:id", handler.Delete)
}
