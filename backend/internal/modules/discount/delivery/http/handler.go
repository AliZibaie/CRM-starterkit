package httpdelivery

import (
	"errors"
	"time"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/internal/modules/discount/usecase"
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

// @Summary Create discount
// @Tags discount
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "discount payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/discounts [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req struct {
		Code           string     `json:"code" validate:"required"`
		Type           string     `json:"type" validate:"required"`
		Value          int64      `json:"value" validate:"required"`
		MinOrderAmount int64      `json:"min_order_amount"`
		MaxUses        int64      `json:"max_uses"`
		StartsAt       *time.Time `json:"starts_at"`
		EndsAt         *time.Time `json:"ends_at"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	d, err := h.service.Create(c.Context(), req.Code, req.Type, req.Value, req.MinOrderAmount, req.MaxUses, req.StartsAt, req.EndsAt)
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, d)
}

// @Summary List discounts
// @Tags discount
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/discounts [get]
func (h *Handler) List(c *fiber.Ctx) error {
	items, err := h.service.List(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary Preview discount
// @Tags discount
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "code order_amount"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/discounts/preview [post]
func (h *Handler) Preview(c *fiber.Ctx) error {
	var req struct {
		Code        string `json:"code" validate:"required"`
		OrderAmount int64  `json:"order_amount" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	res, err := h.service.Preview(c.Context(), req.Code, req.OrderAmount)
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, res)
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/discounts", middleware.Auth(jwt))
	g.Post("/", middleware.RequireAnyPermission("discount.manage"), h.Create)
	g.Get("/", middleware.RequireAnyPermission("discount.read"), h.List)
	g.Post("/preview", middleware.RequireAnyPermission("discount.read"), h.Preview)
}

func mapErr(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		return response.NotFound(c, err.Error())
	case errors.Is(err, usecase.ErrInvalidInput), errors.Is(err, usecase.ErrNotApplicable):
		return response.BadRequest(c, err.Error(), nil)
	default:
		return response.Internal(c, err.Error())
	}
}
