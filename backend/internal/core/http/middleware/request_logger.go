package middleware

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type RequestLogEntry struct {
	RequestID string            `json:"request_id"`
	Method    string            `json:"method"`
	Path      string            `json:"path"`
	IP        string            `json:"ip"`
	UserAgent string            `json:"user_agent"`
	Status    int               `json:"status"`
	Duration  time.Duration     `json:"duration"`
	UserID    string            `json:"user_id,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
	Body      string            `json:"body,omitempty"`
}

type RequestRecorder func(entry RequestLogEntry)

var sensitiveFields = []string{
	"password", "token", "access_token", "refresh_token", "authorization", "secret", "credit_card",
}

func RequestLogger(log *zap.Logger, recorder RequestRecorder) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		headers := map[string]string{}
		c.Request().Header.VisitAll(func(key, value []byte) {
			k := string(key)
			if strings.EqualFold(k, "Authorization") || strings.EqualFold(k, "Cookie") {
				headers[k] = "***"
				return
			}
			headers[k] = string(value)
		})

		body := maskSensitiveBody(c.Body())
		userID := UserIDFromContext(c)
		entry := RequestLogEntry{
			RequestID: c.GetRespHeader("X-Request-ID"),
			Method:    c.Method(),
			Path:      c.Path(),
			IP:        c.IP(),
			UserAgent: c.Get("User-Agent"),
			Status:    c.Response().StatusCode(),
			Duration:  duration,
			UserID:    userID,
			Headers:   headers,
			Body:      body,
		}

		log.Info("http request",
			zap.String("method", entry.Method),
			zap.String("path", entry.Path),
			zap.String("ip", entry.IP),
			zap.Int("status", entry.Status),
			zap.Duration("duration", entry.Duration),
			zap.String("user_id", entry.UserID),
		)

		if recorder != nil {
			recorder(entry)
		}

		return err
	}
}

func maskSensitiveBody(raw []byte) string {
	if len(raw) == 0 {
		return ""
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return string(raw)
	}
	for key := range payload {
		lower := strings.ToLower(key)
		for _, sensitive := range sensitiveFields {
			if strings.Contains(lower, sensitive) {
				payload[key] = "***"
				break
			}
		}
	}
	masked, err := json.Marshal(payload)
	if err != nil {
		return string(raw)
	}
	return string(masked)
}
