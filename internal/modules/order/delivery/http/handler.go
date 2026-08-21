package httpdelivery

import (
	"errors"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/internal/modules/order/usecase"
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

// @Summary Create order from cart
// @Tags order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} false "discount_code"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/orders [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req struct {
		DiscountCode string `json:"discount_code"`
	}
	_ = c.BodyParser(&req)
	order, items, invoice, err := h.service.CreateFromCart(c.Context(), middleware.UserIDFromContext(c), req.DiscountCode)
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, fiber.Map{"order": order, "items": items, "invoice": invoice})
}

// @Summary Checkout order
// @Tags order
// @Security BearerAuth
// @Produce json
// @Param id path string true "order id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/orders/{id}/checkout [post]
func (h *Handler) Checkout(c *fiber.Ctx) error {
	order, err := h.service.Checkout(c.Context(), c.Params("id"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, order)
}

// @Summary Apply discount to order
// @Tags order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "order id"
// @Param body body map[string]interface{} true "code"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/orders/{id}/discount [post]
func (h *Handler) ApplyDiscount(c *fiber.Ctx) error {
	var req struct {
		Code string `json:"code" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	order, err := h.service.ApplyDiscount(c.Context(), c.Params("id"), req.Code)
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, order)
}

// @Summary Refund order
// @Tags order
// @Security BearerAuth
// @Produce json
// @Param id path string true "order id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/orders/{id}/refund [post]
func (h *Handler) Refund(c *fiber.Ctx) error {
	order, err := h.service.Refund(c.Context(), c.Params("id"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, order)
}

// @Summary Get order
// @Tags order
// @Security BearerAuth
// @Produce json
// @Param id path string true "order id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/orders/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	order, items, invoice, err := h.service.Get(c.Context(), c.Params("id"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, fiber.Map{"order": order, "items": items, "invoice": invoice})
}

// @Summary List my orders
// @Tags order
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/orders [get]
func (h *Handler) ListMine(c *fiber.Ctx) error {
	items, err := h.service.ListMine(c.Context(), middleware.UserIDFromContext(c))
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/orders", middleware.Auth(jwt))
	g.Post("/", middleware.RequireAnyPermission("order.create"), h.Create)
	g.Get("/", middleware.RequireAnyPermission("order.read"), h.ListMine)
	g.Get("/:id", middleware.RequireAnyPermission("order.read"), h.Get)
	g.Post("/:id/checkout", middleware.RequireAnyPermission("order.checkout"), h.Checkout)
	g.Post("/:id/discount", middleware.RequireAnyPermission("order.discount"), h.ApplyDiscount)
	g.Post("/:id/refund", middleware.RequireAnyPermission("order.refund"), h.Refund)
}

func mapErr(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		return response.NotFound(c, err.Error())
	case errors.Is(err, usecase.ErrInvalidInput), errors.Is(err, usecase.ErrEmptyCart), errors.Is(err, usecase.ErrBadStatus):
		return response.BadRequest(c, err.Error(), nil)
	default:
		return response.Internal(c, err.Error())
	}
}
