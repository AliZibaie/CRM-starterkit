package httpdelivery

import (
	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/internal/modules/notification/domain"
	"crm-starterkit/internal/modules/notification/usecase"
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

type dispatchRequest struct {
	UserID  string         `json:"user_id"`
	To      string         `json:"to"`
	Title   string         `json:"title" validate:"required"`
	Body    string         `json:"body" validate:"required"`
	Channel string         `json:"channel"`
	Data    map[string]any `json:"data"`
}

// @Summary Dispatch notification
// @Tags notification
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "notification payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/notifications/dispatch [post]
func (h *Handler) Dispatch(c *fiber.Ctx) error {
	var req dispatchRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	if err := h.service.Dispatch(c.Context(), domain.Payload{
		UserID: req.UserID, To: req.To, Title: req.Title, Body: req.Body, Channel: req.Channel, Data: req.Data,
	}); err != nil {
		return response.Internal(c, err.Error())
	}
	return response.Message(c, fiber.StatusOK, "dispatched", nil)
}

// @Summary List my notifications
// @Tags notification
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/notifications/me [get]
func (h *Handler) ListMine(c *fiber.Ctx) error {
	items, err := h.service.ListMine(c.Context(), middleware.UserIDFromContext(c))
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary Mark notification read
// @Tags notification
// @Security BearerAuth
// @Produce json
// @Param id path string true "notification id"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/notifications/{id}/read [post]
func (h *Handler) MarkRead(c *fiber.Ctx) error {
	if err := h.service.MarkRead(c.Context(), c.Params("id"), middleware.UserIDFromContext(c)); err != nil {
		return response.Internal(c, err.Error())
	}
	return response.Message(c, fiber.StatusOK, "marked read", nil)
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/notifications", middleware.Auth(jwt))
	g.Post("/dispatch", middleware.RequireAnyPermission("notification.dispatch"), h.Dispatch)
	g.Get("/me", middleware.RequireAnyPermission("notification.read"), h.ListMine)
	g.Post("/:id/read", middleware.RequireAnyPermission("notification.read"), h.MarkRead)
}
