package httpdelivery

import (
	"strconv"
	"time"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/modules/auditlog/domain"
	"crm-starterkit/internal/modules/auditlog/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Handler struct{ service *usecase.Service }

func NewHandler(service *usecase.Service) *Handler { return &Handler{service: service} }

// @Summary List audit logs
// @Tags auditlog
// @Security BearerAuth
// @Produce json
// @Param actor_id query string false "actor"
// @Param module query string false "module"
// @Param entity_type query string false "entity type"
// @Param entity_id query string false "entity id"
// @Param from query string false "RFC3339"
// @Param to query string false "RFC3339"
// @Param limit query int false "limit"
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/audit-logs [get]
func (h *Handler) List(c *fiber.Ctx) error {
	f := domain.Filter{
		ActorID: c.Query("actor_id"), Module: c.Query("module"),
		EntityType: c.Query("entity_type"), EntityID: c.Query("entity_id"),
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
	if v := c.Query("limit"); v != "" {
		f.Limit, _ = strconv.Atoi(v)
	}
	items, err := h.service.List(c.Context(), f)
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/audit-logs", middleware.Auth(jwt))
	g.Get("/", middleware.RequireAnyPermission("auditlog.read"), h.List)
}
