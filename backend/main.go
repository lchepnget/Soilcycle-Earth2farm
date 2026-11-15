package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/soilcycle-earth2farm/backend/internal/config"
	"github.com/soilcycle-earth2farm/backend/internal/database"
	"github.com/soilcycle-earth2farm/backend/internal/middleware"
	authmodule "github.com/soilcycle-earth2farm/backend/internal/modules/auth"
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
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowedOrigins,
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
	}))

	// API versioning and authentication.
	api := app.Group("/api/v1")
	api.Use(middleware.RequireAuth(cfg.JWTSecret))
	authmodule.RegisterRoutes(api, db, cfg.JWTSecret)

	// Register API modules after authentication middleware.
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
