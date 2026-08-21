package httpcore

import (
	"github.com/gofiber/fiber/v2"
)

func MountCoreRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":    "crm-starterkit",
			"message": "api is running",
		})
	})
}
