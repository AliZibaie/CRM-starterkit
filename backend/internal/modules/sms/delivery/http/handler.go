package httpdelivery

import (
	"errors"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/internal/modules/sms/usecase"
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

type sendRequest struct {
	To   string `json:"to" validate:"required"`
	Text string `json:"text" validate:"required"`
}

// @Summary Send SMS
// @Tags sms
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "to text"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/sms/send [post]
func (h *Handler) Send(c *fiber.Ctx) error {
	var req sendRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	if err := h.service.Send(c.Context(), req.To, req.Text); err != nil {
		if errors.Is(err, usecase.ErrInvalidInput) {
			return response.BadRequest(c, err.Error(), nil)
		}
		return response.Internal(c, err.Error())
	}
	return response.Message(c, fiber.StatusOK, "sms queued", nil)
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/sms", middleware.Auth(jwt))
	g.Post("/send", middleware.RequireAnyPermission("sms.send"), h.Send)
}
