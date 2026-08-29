package httpdelivery

import (
	"errors"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/internal/modules/payment/usecase"
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

type initRequest struct {
	OrderID     string `json:"order_id" validate:"required"`
	Gateway     string `json:"gateway" validate:"required"`
	Description string `json:"description"`
}

type verifyRequest struct {
	Authority string `json:"authority"`
	OrderID   string `json:"order_id"`
	Gateway   string `json:"gateway" validate:"required"`
}

// @Summary List payment gateways
// @Tags payment
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/payments/gateways [get]
func (h *Handler) ListGateways(c *fiber.Ctx) error {
	return response.OK(c, h.service.ListGatewayStatus())
}

func (h *Handler) ListFailed(c *fiber.Ctx) error {
	items, err := h.service.ListFailedPayments(c.Context(), 50)
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary Init payment
// @Tags payment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "order_id gateway description"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/payments/init [post]
func (h *Handler) Init(c *fiber.Ctx) error {
	var req initRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	p, result, err := h.service.InitPayment(c.Context(), req.OrderID, req.Gateway, req.Description)
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, fiber.Map{
		"payment":     p,
		"payment_url": result.PaymentURL,
		"authority":   result.Authority,
	})
}

// @Summary Verify payment
// @Tags payment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "authority order_id gateway"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/payments/verify [post]
func (h *Handler) Verify(c *fiber.Ctx) error {
	var req verifyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	p, result, err := h.service.VerifyPayment(c.Context(), req.Gateway, req.Authority, req.OrderID)
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, fiber.Map{"payment": p, "verify": result})
}

// @Summary Payment gateway callback
// @Tags payment
// @Produce json
// @Param gateway path string true "gateway"
// @Param Authority query string false "authority"
// @Param session_id query string false "session id"
// @Param order_id query string false "order id"
// @Param Status query string false "status"
// @Param canceled query string false "canceled"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/payments/callback/{gateway} [get]
func (h *Handler) Callback(c *fiber.Ctx) error {
	gateway := c.Params("gateway")
	authority := c.Query("Authority")
	if authority == "" {
		authority = c.Query("session_id")
	}
	orderID := c.Query("order_id")
	status := c.Query("Status")
	if status == "" {
		status = c.Query("canceled")
	}
	p, result, err := h.service.HandleCallback(c.Context(), gateway, authority, orderID, status)
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, fiber.Map{"payment": p, "verify": result})
}

// @Summary Refund payment
// @Tags payment
// @Security BearerAuth
// @Produce json
// @Param order_id path string true "order id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/payments/refund/{order_id} [post]
func (h *Handler) Refund(c *fiber.Ctx) error {
	p, err := h.service.Refund(c.Context(), c.Params("order_id"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, p)
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/payments")
	g.Get("/gateways", middleware.Auth(jwt), middleware.RequireAnyPermission("payment.init"), h.ListGateways)
	g.Get("/failed", middleware.Auth(jwt), middleware.RequireAnyPermission("payment.verify"), h.ListFailed)
	g.Post("/init", middleware.Auth(jwt), middleware.RequireAnyPermission("payment.init"), h.Init)
	g.Post("/verify", middleware.Auth(jwt), middleware.RequireAnyPermission("payment.verify"), h.Verify)
	g.Get("/callback/:gateway", h.Callback)
	g.Post("/refund/:order_id", middleware.Auth(jwt), middleware.RequireAnyPermission("payment.refund"), h.Refund)
}

func mapErr(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		return response.NotFound(c, err.Error())
	case errors.Is(err, usecase.ErrInvalidInput), errors.Is(err, usecase.ErrGatewayMissing), errors.Is(err, usecase.ErrBadOrder):
		return response.BadRequest(c, err.Error(), nil)
	default:
		return response.Internal(c, err.Error())
	}
}
