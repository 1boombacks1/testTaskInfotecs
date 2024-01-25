package v1

import (
	"github.com/1boombacks1/testTaskInfotecs/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewRouter(app *fiber.App, service *service.Service) {
	app.Use(recover.New())
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	v1 := app.Group("/api/v1")
	{
		newWalletRoutes(v1, *service)
	}
}
