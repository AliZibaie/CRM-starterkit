package middleware

import (
	"crm-starterkit/internal/core/http/response"

	"github.com/gofiber/fiber/v2"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

func Recover(log *zap.Logger) fiber.Handler {
	return fiberrecover.New(fiberrecover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			log.Error("panic recovered",
				zap.Any("error", e),
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
			)
		},
	})
}

func RecoverJSON(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic recovered",
					zap.Any("error", r),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
				)
				_ = response.Internal(c, "internal server error")
			}
		}()
		return c.Next()
	}
}
