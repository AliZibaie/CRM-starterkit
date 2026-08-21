package response

import (
	"github.com/gofiber/fiber/v2"
)

type Envelope struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

func OK(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(Envelope{
		Success: true,
		Data:    data,
	})
}

func Created(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Envelope{
		Success: true,
		Data:    data,
	})
}

func Message(c *fiber.Ctx, status int, message string, data any) error {
	return c.Status(status).JSON(Envelope{
		Success: status >= 200 && status < 300,
		Message: message,
		Data:    data,
	})
}

func WithMeta(c *fiber.Ctx, data any, meta any) error {
	return c.Status(fiber.StatusOK).JSON(Envelope{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}
