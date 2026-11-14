package earthdata

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router) {
	client := NewPOWERClient()
	service := NewService(client)
	handler := NewHandler(service)

	earthdataGroup := router.Group("/earthdata")

	earthdataGroup.Post("/climate", handler.GetClimateData)
}
