package httpdelivery

import (
	"errors"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/modules/setting/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Handler struct{ service *usecase.Service }

func NewHandler(service *usecase.Service) *Handler { return &Handler{service: service} }

// @Summary Get setting by key
// @Tags setting
// @Security BearerAuth
// @Produce json
// @Param key path string true "key"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/settings/{key} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	item, err := h.service.Get(c.Context(), c.Params("key"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, item)
}

// @Summary Set setting
// @Tags setting
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]string true "key value"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/settings [put]
func (h *Handler) Set(c *fiber.Ctx) error {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	item, err := h.service.Set(c.Context(), req.Key, req.Value, middleware.UserIDFromContext(c))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, item)
}

// @Summary List settings
// @Tags setting
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/settings [get]
func (h *Handler) List(c *fiber.Ctx) error {
	items, err := h.service.List(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/settings", middleware.Auth(jwt))
	g.Get("/", middleware.RequireAnyPermission("setting.read"), h.List)
	g.Get("/:key", middleware.RequireAnyPermission("setting.read"), h.Get)
	g.Put("/", middleware.RequireAnyPermission("setting.write"), h.Set)
}

func mapErr(c *fiber.Ctx, err error) error {
	if errors.Is(err, usecase.ErrNotFound) {
		return response.NotFound(c, err.Error())
	}
	if errors.Is(err, usecase.ErrInvalidInput) {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.Internal(c, err.Error())
}
