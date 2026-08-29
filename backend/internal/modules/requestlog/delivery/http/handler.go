package httpdelivery

import (
	"strconv"
	"time"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/modules/requestlog/domain"
	"crm-starterkit/internal/modules/requestlog/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Handler struct{ service *usecase.Service }

func NewHandler(service *usecase.Service) *Handler { return &Handler{service: service} }

// @Summary List request logs
// @Tags requestlog
// @Security BearerAuth
// @Produce json
// @Param ip query string false "ip"
// @Param path query string false "path"
// @Param user_id query string false "user"
// @Param status query int false "status"
// @Param from query string false "RFC3339"
// @Param to query string false "RFC3339"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/request-logs [get]
func (h *Handler) List(c *fiber.Ctx) error {
	f := domain.ListFilter{
		IP: c.Query("ip"), Path: c.Query("path"), UserID: c.Query("user_id"),
	}
	if v := c.Query("status"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			f.Status = &n
		}
	}
	if v := c.Query("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.From = &t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.To = &t
		}
	}
	f.Limit, _ = strconv.Atoi(c.Query("limit", "100"))
	f.Offset, _ = strconv.Atoi(c.Query("offset", "0"))
	items, err := h.service.List(c.Context(), f)
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary Request log stats
// @Tags requestlog
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/request-logs/stats [get]
func (h *Handler) Stats(c *fiber.Ctx) error {
	f := domain.ListFilter{IP: c.Query("ip"), Path: c.Query("path"), UserID: c.Query("user_id")}
	if v := c.Query("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.From = &t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.To = &t
		}
	}
	stats, err := h.service.Stats(c.Context(), f)
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, domain.StatsResponse{
		Total: stats.Total, AvgDurationMs: stats.AvgDurationMs,
		ErrorCount: stats.ErrorCount, SuspiciousCount: stats.Suspicious,
	})
}

// @Summary Check suspicious IP
// @Tags requestlog
// @Security BearerAuth
// @Produce json
// @Param ip path string true "ip"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/request-logs/suspicious/{ip} [get]
func (h *Handler) Suspicious(c *fiber.Ctx) error {
	ip := c.Params("ip")
	return response.OK(c, fiber.Map{"ip": ip, "suspicious": h.service.IsSuspiciousIP(ip)})
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/request-logs", middleware.Auth(jwt))
	g.Get("/", middleware.RequireAnyPermission("requestlog.read"), h.List)
	g.Get("/stats", middleware.RequireAnyPermission("requestlog.read"), h.Stats)
	g.Get("/suspicious/:ip", middleware.RequireAnyPermission("requestlog.read"), h.Suspicious)
}
