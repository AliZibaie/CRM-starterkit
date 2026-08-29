package httpdelivery

import (
	"errors"
	"time"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/internal/modules/banner/domain"
	"crm-starterkit/internal/modules/banner/usecase"
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
	Title    string `json:"title" validate:"required"`
	ImageURL string `json:"image_url" validate:"required"`
	LinkURL  string `json:"link_url"`
	Position string `json:"position" validate:"required"`
}

type scheduleRequest struct {
	StartsAt *time.Time `json:"starts_at"`
	EndsAt   *time.Time `json:"ends_at"`
	IsActive *bool      `json:"is_active"`
}

type positionRequest struct {
	Position string `json:"position" validate:"required"`
}

// @Summary Create banner
// @Tags banner
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "banner payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/banners [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	b, err := h.service.Create(c.Context(), domain.CreateInput{
		Title: req.Title, ImageURL: req.ImageURL, LinkURL: req.LinkURL, Position: req.Position,
	})
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, b)
}

// @Summary Update banner
// @Tags banner
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "banner id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/banners/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	b, err := h.service.Update(c.Context(), c.Params("id"), domain.CreateInput{
		Title: req.Title, ImageURL: req.ImageURL, LinkURL: req.LinkURL, Position: req.Position,
	})
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, b)
}

// @Summary Delete banner
// @Tags banner
// @Security BearerAuth
// @Produce json
// @Param id path string true "banner id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/banners/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.service.Delete(c.Context(), c.Params("id")); err != nil {
		return mapErr(c, err)
	}
	return response.Message(c, fiber.StatusOK, "banner deleted", nil)
}

// @Summary Schedule banner
// @Tags banner
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "banner id"
// @Param body body map[string]interface{} true "schedule payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/banners/{id}/schedule [put]
func (h *Handler) Schedule(c *fiber.Ctx) error {
	var req scheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	b, err := h.service.Schedule(c.Context(), c.Params("id"), domain.ScheduleInput{
		StartsAt: req.StartsAt, EndsAt: req.EndsAt, IsActive: req.IsActive,
	})
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, b)
}

// @Summary Set banner position
// @Tags banner
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "banner id"
// @Param body body map[string]interface{} true "position"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/banners/{id}/position [put]
func (h *Handler) SetPosition(c *fiber.Ctx) error {
	var req positionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	b, err := h.service.SetPosition(c.Context(), c.Params("id"), req.Position)
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, b)
}

// @Summary List banners
// @Tags banner
// @Security BearerAuth
// @Produce json
// @Param position query string false "position"
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/banners [get]
func (h *Handler) List(c *fiber.Ctx) error {
	if pos := c.Query("position"); pos != "" {
		items, err := h.service.ListByPosition(c.Context(), pos)
		if err != nil {
			return response.Internal(c, err.Error())
		}
		return response.OK(c, items)
	}
	items, err := h.service.List(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary List banner positions
// @Tags banner
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/banners/positions [get]
func (h *Handler) Positions(c *fiber.Ctx) error {
	return response.OK(c, h.service.Positions())
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/banners", middleware.Auth(jwt))
	g.Post("/", middleware.RequireAnyPermission("banner.create"), h.Create)
	g.Get("/", middleware.RequireAnyPermission("banner.read"), h.List)
	g.Get("/positions", middleware.RequireAnyPermission("banner.read"), h.Positions)
	g.Put("/:id", middleware.RequireAnyPermission("banner.create"), h.Update)
	g.Delete("/:id", middleware.RequireAnyPermission("banner.create"), h.Delete)
	g.Put("/:id/schedule", middleware.RequireAnyPermission("banner.schedule"), h.Schedule)
	g.Put("/:id/position", middleware.RequireAnyPermission("banner.position"), h.SetPosition)
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
