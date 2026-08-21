package response

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type ErrorBody struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func Fail(c *fiber.Ctx, status int, code, message string, details any) error {
	return c.Status(status).JSON(ErrorBody{
		Success: false,
		Code:    code,
		Message: message,
		Details: details,
	})
}

func BadRequest(c *fiber.Ctx, message string, details any) error {
	return Fail(c, fiber.StatusBadRequest, "bad_request", message, details)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "unauthorized"
	}
	return Fail(c, fiber.StatusUnauthorized, "unauthorized", message, nil)
}

func Forbidden(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "forbidden"
	}
	return Fail(c, fiber.StatusForbidden, "forbidden", message, nil)
}

func NotFound(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "not found"
	}
	return Fail(c, fiber.StatusNotFound, "not_found", message, nil)
}

func Conflict(c *fiber.Ctx, message string) error {
	return Fail(c, fiber.StatusConflict, "conflict", message, nil)
}

func TooManyRequests(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "too many requests"
	}
	return Fail(c, fiber.StatusTooManyRequests, "rate_limited", message, nil)
}

func Internal(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "internal server error"
	}
	return Fail(c, fiber.StatusInternalServerError, "internal_error", message, nil)
}

func FromError(c *fiber.Ctx, err error) error {
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return Fail(c, fiberErr.Code, "fiber_error", fiberErr.Message, nil)
	}
	return Internal(c, err.Error())
}
