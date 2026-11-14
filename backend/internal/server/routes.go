package server

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/soilcycle-earth2farm/backend/internal/middleware"
	"github.com/soilcycle-earth2farm/backend/internal/modules/crops"
	"github.com/soilcycle-earth2farm/backend/internal/modules/farmers"
	"github.com/soilcycle-earth2farm/backend/internal/modules/farms"
)

// RegisterRoutes attaches health, root, and versioned API routes to app.
func RegisterRoutes(app *fiber.App, db *pgxpool.Pool) {
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Welcome to SoilCycle by Earth2Farm",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	app.Get("/health", func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"success":  false,
				"status":   "unhealthy",
				"database": "disconnected",
			})
		}
		return c.JSON(fiber.Map{
			"success":  true,
			"status":   "healthy",
			"database": "connected",
		})
	})

	api := app.Group("/api/v1")
	// Register middleware before module routes so it runs before handlers.
	api.Use(middleware.RequireAuth())

	crops.RegisterRoutes(api, db)
	farmers.RegisterRoutes(api, db)
	farms.RegisterRoutes(api, db)
}
