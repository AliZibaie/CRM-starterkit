package httpdelivery

import (
	"errors"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/modules/wishlist/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Handler struct{ service *usecase.Service }

func NewHandler(service *usecase.Service) *Handler { return &Handler{service: service} }

// @Summary Add wishlist item
// @Tags wishlist
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "item_id"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/wishlist [post]
func (h *Handler) Add(c *fiber.Ctx) error {
	var req struct {
		ItemID string `json:"item_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	item, err := h.service.Add(c.Context(), middleware.UserIDFromContext(c), req.ItemID)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidInput) {
			return response.BadRequest(c, err.Error(), nil)
		}
		return response.Internal(c, err.Error())
	}
	return response.Created(c, item)
}

// @Summary Remove wishlist item
// @Tags wishlist
// @Security BearerAuth
// @Produce json
// @Param item_id path string true "item id"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/wishlist/{item_id} [delete]
func (h *Handler) Remove(c *fiber.Ctx) error {
	if err := h.service.Remove(c.Context(), middleware.UserIDFromContext(c), c.Params("item_id")); err != nil {
		return response.Internal(c, err.Error())
	}
	return response.Message(c, fiber.StatusOK, "removed", nil)
}

// @Summary List wishlist
// @Tags wishlist
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/wishlist [get]
func (h *Handler) List(c *fiber.Ctx) error {
	items, err := h.service.List(c.Context(), middleware.UserIDFromContext(c))
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/wishlist", middleware.Auth(jwt))
	g.Get("/", middleware.RequireAnyPermission("wishlist.manage"), h.List)
	g.Post("/", middleware.RequireAnyPermission("wishlist.manage"), h.Add)
	g.Delete("/:item_id", middleware.RequireAnyPermission("wishlist.manage"), h.Remove)
}
