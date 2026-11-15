package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router fiber.Router, db *pgxpool.Pool, jwtSecret string) {
	repository := NewAuthRepository(db)
	service := NewAuthService(repository, jwtSecret)
	handler := NewAuthHandler(service)

	group := router.Group("/auth")
	group.Post("/login", handler.Login)
}
