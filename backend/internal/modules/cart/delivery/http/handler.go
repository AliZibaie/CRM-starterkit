package httpdelivery

import (
	"errors"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/internal/modules/cart/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service  *usecase.Service
	validate *validator.Validator
}

func NewHandler(service *usecase.Service, validate *validator.Validator) *Handler {
	return &Handler{service: service, validate: validate}
}

// @Summary Get cart
// @Tags cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/cart [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	cart, items, total, err := h.service.GetCart(c.Context(), middleware.UserIDFromContext(c))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, fiber.Map{"cart": cart, "items": items, "total": total})
}

// @Summary Add cart item
// @Tags cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "item_id variant_id quantity"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/cart/items [post]
func (h *Handler) Add(c *fiber.Ctx) error {
	var req struct {
		ItemID    string `json:"item_id" validate:"required"`
		VariantID string `json:"variant_id"`
		Quantity  int64  `json:"quantity" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	item, err := h.service.AddItem(c.Context(), middleware.UserIDFromContext(c), req.ItemID, req.VariantID, req.Quantity)
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, item)
}

// @Summary Remove cart item
// @Tags cart
// @Security BearerAuth
// @Produce json
// @Param item_id path string true "item id"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/cart/items/{item_id} [delete]
func (h *Handler) Remove(c *fiber.Ctx) error {
	if err := h.service.RemoveItem(c.Context(), middleware.UserIDFromContext(c), c.Params("item_id")); err != nil {
		return mapErr(c, err)
	}
	return response.Message(c, fiber.StatusOK, "removed", nil)
}

// @Summary Clear cart
// @Tags cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/cart [delete]
func (h *Handler) Clear(c *fiber.Ctx) error {
	if err := h.service.Clear(c.Context(), middleware.UserIDFromContext(c)); err != nil {
		return mapErr(c, err)
	}
	return response.Message(c, fiber.StatusOK, "cleared", nil)
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/cart", middleware.Auth(jwt))
	g.Get("/", middleware.RequireAnyPermission("cart.read"), h.Get)
	g.Post("/items", middleware.RequireAnyPermission("cart.write"), h.Add)
	g.Delete("/items/:item_id", middleware.RequireAnyPermission("cart.write"), h.Remove)
	g.Delete("/", middleware.RequireAnyPermission("cart.write"), h.Clear)
}

func mapErr(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		return response.NotFound(c, err.Error())
	case errors.Is(err, usecase.ErrInvalidInput), errors.Is(err, usecase.ErrOutOfStock):
		return response.BadRequest(c, err.Error(), nil)
	default:
		return response.Internal(c, err.Error())
	}
}
