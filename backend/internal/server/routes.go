package server

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/soilcycle-earth2farm/backend/internal/config"
	"github.com/soilcycle-earth2farm/backend/internal/middleware"
	authmodule "github.com/soilcycle-earth2farm/backend/internal/modules/auth"
	"github.com/soilcycle-earth2farm/backend/internal/modules/crops"
	"github.com/soilcycle-earth2farm/backend/internal/modules/earthdata"
	"github.com/soilcycle-earth2farm/backend/internal/modules/farmers"
	"github.com/soilcycle-earth2farm/backend/internal/modules/farms"
	"github.com/soilcycle-earth2farm/backend/internal/modules/rotation"
	"github.com/soilcycle-earth2farm/backend/internal/modules/soil"
)

// RegisterRoutes attaches health, root, and versioned API routes to app.
func RegisterRoutes(app *fiber.App, db *pgxpool.Pool, cfg *config.Config) {
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
	api.Use(middleware.RequireAuth(cfg.JWTSecret))
	authmodule.RegisterRoutes(api, db, cfg.JWTSecret)

	crops.RegisterRoutes(api, db)
	farmers.RegisterRoutes(api, db)
	farms.RegisterRoutes(api, db)
	soil.RegisterRoutes(api, db)
	rotation.RegisterRoutes(api, db)
	earthdata.RegisterRoutes(api)
}
