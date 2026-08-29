package httpdelivery

import (
	"errors"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/internal/modules/announcement/domain"
	"crm-starterkit/internal/modules/announcement/usecase"
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

type createRequest struct {
	Title       string `json:"title" validate:"required"`
	Body        string `json:"body" validate:"required"`
	Audience    string `json:"audience"`
	AudienceRef string `json:"audience_ref"`
}

// @Summary Create announcement
// @Tags announcement
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "announcement payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/announcements [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	a, err := h.service.Create(c.Context(), domain.CreateInput{
		Title: req.Title, Body: req.Body, Audience: req.Audience,
		AudienceRef: req.AudienceRef, CreatedBy: middleware.UserIDFromContext(c),
	})
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, a)
}

// @Summary List announcements
// @Tags announcement
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/announcements [get]
func (h *Handler) List(c *fiber.Ctx) error {
	items, err := h.service.List(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary Get announcement
// @Tags announcement
// @Security BearerAuth
// @Produce json
// @Param id path string true "announcement id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/announcements/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	a, err := h.service.Get(c.Context(), c.Params("id"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, fiber.Map{
		"announcement": a,
		"audience":     h.service.ResolveAudience(a),
	})
}

type updateRequest struct {
	Title       string `json:"title" validate:"required"`
	Body        string `json:"body" validate:"required"`
	Audience    string `json:"audience"`
	AudienceRef string `json:"audience_ref"`
}

// @Summary Update announcement
// @Tags announcement
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "announcement id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/announcements/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	var req updateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	a, err := h.service.Update(c.Context(), c.Params("id"), domain.CreateInput{
		Title: req.Title, Body: req.Body, Audience: req.Audience, AudienceRef: req.AudienceRef,
	})
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, a)
}

// @Summary Delete announcement
// @Tags announcement
// @Security BearerAuth
// @Produce json
// @Param id path string true "announcement id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/announcements/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.service.Delete(c.Context(), c.Params("id")); err != nil {
		return mapErr(c, err)
	}
	return response.Message(c, fiber.StatusOK, "announcement deleted", nil)
}

// @Summary Broadcast announcement
// @Tags announcement
// @Security BearerAuth
// @Produce json
// @Param id path string true "announcement id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/announcements/{id}/broadcast [post]
func (h *Handler) Broadcast(c *fiber.Ctx) error {
	a, err := h.service.Broadcast(c.Context(), c.Params("id"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, a)
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/announcements", middleware.Auth(jwt))
	g.Post("/", middleware.RequireAnyPermission("announcement.create"), h.Create)
	g.Get("/", middleware.RequireAnyPermission("announcement.read"), h.List)
	g.Get("/:id", middleware.RequireAnyPermission("announcement.read"), h.Get)
	g.Put("/:id", middleware.RequireAnyPermission("announcement.create"), h.Update)
	g.Delete("/:id", middleware.RequireAnyPermission("announcement.create"), h.Delete)
	g.Post("/:id/broadcast", middleware.RequireAnyPermission("announcement.broadcast"), h.Broadcast)
}

func mapErr(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		return response.NotFound(c, err.Error())
	case errors.Is(err, usecase.ErrInvalidInput):
		return response.BadRequest(c, err.Error(), nil)
	default:
		return response.Internal(c, err.Error())
	}
}
