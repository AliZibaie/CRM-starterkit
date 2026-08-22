package httpdelivery

import (
	"errors"
	"strconv"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/modules/media/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Handler struct{ service *usecase.Service }

func NewHandler(service *usecase.Service) *Handler { return &Handler{service: service} }

// @Summary Upload media
// @Tags media
// @Security BearerAuth
// @Accept mpfd
// @Produce json
// @Param file formData file true "file"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/media/upload [post]
func (h *Handler) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return response.BadRequest(c, "file required", nil)
	}
	f, err := file.Open()
	if err != nil {
		return response.Internal(c, err.Error())
	}
	defer f.Close()
	ct := file.Header.Get("Content-Type")
	media, err := h.service.Upload(c.Context(), middleware.UserIDFromContext(c), file.Filename, ct, f, file.Size)
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, media)
}

// @Summary List my media
// @Tags media
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/media [get]
func (h *Handler) ListMine(c *fiber.Ctx) error {
	items, err := h.service.ListMine(c.Context(), middleware.UserIDFromContext(c))
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary Delete media
// @Tags media
// @Security BearerAuth
// @Produce json
// @Param id path string true "media id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/media/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.service.Delete(c.Context(), c.Params("id"), middleware.UserIDFromContext(c)); err != nil {
		return mapErr(c, err)
	}
	return response.Message(c, fiber.StatusOK, "deleted", nil)
}

// @Summary Generate thumbnail
// @Tags media
// @Security BearerAuth
// @Produce json
// @Param id path string true "media id"
// @Param force query bool false "force"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/media/{id}/thumbnail [post]
func (h *Handler) Thumbnail(c *fiber.Ctx) error {
	force, _ := strconv.ParseBool(c.Query("force", "false"))
	_ = force
	media, err := h.service.ListMine(c.Context(), middleware.UserIDFromContext(c))
	if err != nil {
		return response.Internal(c, err.Error())
	}
	for i := range media {
		if media[i].ID == c.Params("id") {
			out, err := h.service.GenerateThumbnail(c.Context(), &media[i])
			if err != nil {
				return mapErr(c, err)
			}
			return response.OK(c, out)
		}
	}
	return response.NotFound(c, "media not found")
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/media", middleware.Auth(jwt))
	g.Post("/upload", middleware.RequireAnyPermission("media.upload"), h.Upload)
	g.Get("/", middleware.RequireAnyPermission("media.read"), h.ListMine)
	g.Post("/:id/thumbnail", middleware.RequireAnyPermission("media.upload"), h.Thumbnail)
	g.Delete("/:id", middleware.RequireAnyPermission("media.delete"), h.Delete)
}

func mapErr(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		return response.NotFound(c, err.Error())
	case errors.Is(err, usecase.ErrInvalidInput), errors.Is(err, usecase.ErrNoStorage):
		return response.BadRequest(c, err.Error(), nil)
	default:
		return response.Internal(c, err.Error())
	}
}
