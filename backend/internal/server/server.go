package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/soilcycle-earth2farm/backend/internal/config"
)

// New creates the configured Fiber application and registers all API routes.
func New(db *pgxpool.Pool, cfg *config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "SoilCycle by Earth2Farm",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
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

	RegisterRoutes(app, db, cfg)
	return app
}

// Listen serves requests until the context is cancelled or the server fails.
// Cancellation triggers a bounded graceful shutdown.
func Listen(ctx context.Context, app *fiber.App, address string) error {
	listenErr := make(chan error, 1)

	go func() {
		listenErr <- app.Listen(address)
	}()

	select {
	case err := <-listenErr:
		if err != nil {
			return fmt.Errorf("listen on %s: %w", address, err)
		}
		return nil

	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.ShutdownWithContext(shutdownCtx); err != nil {
			return fmt.Errorf("graceful server shutdown: %w", err)
		}

		select {
		case err := <-listenErr:
			if err != nil && !errors.Is(err, context.Canceled) {
				return fmt.Errorf("server stopped with error: %w", err)
			}
		default:
		}

		return nil
	}
}
