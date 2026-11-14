package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/soilcycle-earth2farm/backend/internal/config"
	"github.com/soilcycle-earth2farm/backend/internal/database"
	"github.com/soilcycle-earth2farm/backend/internal/modules/crops"
	"github.com/soilcycle-earth2farm/backend/internal/modules/earthdata"
	"github.com/soilcycle-earth2farm/backend/internal/modules/farmers"
	"github.com/soilcycle-earth2farm/backend/internal/modules/farms"
	"github.com/soilcycle-earth2farm/backend/internal/modules/rotation"
	"github.com/soilcycle-earth2farm/backend/internal/modules/soil"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	defer db.Close()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	// API versioning
	api := app.Group("/api/v1")

	// Register Crop API

	crops.RegisterRoutes(api, db)
	farmers.RegisterRoutes(api, db)
	farms.RegisterRoutes(api, db)

	crops.RegisterRoutes(api, db)
	crops.RegisterRoutes(api, db)
	farmers.RegisterRoutes(api, db)
	farms.RegisterRoutes(api, db)
	soil.RegisterRoutes(api, db)
	rotation.RegisterRoutes(api, db)
	earthdata.RegisterRoutes(api)

	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to SoilCycle by Earth2Farm",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "healthy",
			"database": "connected",
		})
	})

	log.Fatal(app.Listen(":" + cfg.AppPort))
}
