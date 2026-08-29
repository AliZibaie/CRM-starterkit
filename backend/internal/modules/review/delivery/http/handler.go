package httpdelivery

import (
	"errors"
	"strconv"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/modules/review/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Handler struct{ service *usecase.Service }

func NewHandler(service *usecase.Service) *Handler { return &Handler{service: service} }

// @Summary Create review
// @Tags review
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "item_id rating comment"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/reviews [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req struct {
		ItemID  string `json:"item_id"`
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	r, err := h.service.Create(c.Context(), req.ItemID, middleware.UserIDFromContext(c), req.Comment, req.Rating)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidInput) {
			return response.BadRequest(c, err.Error(), nil)
		}
		return response.Internal(c, err.Error())
	}
	return response.Created(c, r)
}

// @Summary List reviews by item
// @Tags review
// @Produce json
// @Param item_id path string true "item id"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/reviews/items/{item_id} [get]
func (h *Handler) ListByItem(c *fiber.Ctx) error {
	items, avg, err := h.service.ListByItem(c.Context(), c.Params("item_id"))
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, fiber.Map{"reviews": items, "average": avg})
}

// @Summary Delete review
// @Tags review
// @Security BearerAuth
// @Produce json
// @Param id path string true "review id"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/reviews/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.service.Delete(c.Context(), c.Params("id"), middleware.UserIDFromContext(c)); err != nil {
		return response.Internal(c, err.Error())
	}
	return response.Message(c, fiber.StatusOK, "deleted", nil)
}

func (h *Handler) ListAll(c *fiber.Ctx) error {
	var rating *int
	if v := c.Query("rating"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			rating = &n
		}
	}
	items, err := h.service.ListAll(c.Context(), c.Query("item_id"), rating, 100)
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

func (h *Handler) SetVisibility(c *fiber.Ctx) error {
	var req struct {
		IsVisible bool `json:"is_visible"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.service.SetVisibility(c.Context(), c.Params("id"), req.IsVisible); err != nil {
		return response.Internal(c, err.Error())
	}
	return response.Message(c, fiber.StatusOK, "updated", nil)
}

func (h *Handler) AdminDelete(c *fiber.Ctx) error {
	if err := h.service.AdminDelete(c.Context(), c.Params("id")); err != nil {
		return response.Internal(c, err.Error())
	}
	return response.Message(c, fiber.StatusOK, "deleted", nil)
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	router.Get("/reviews/items/:item_id", h.ListByItem)
	g := router.Group("/reviews", middleware.Auth(jwt))
	g.Get("/admin", middleware.RequireAnyPermission("review.moderate"), h.ListAll)
	g.Post("/", middleware.RequireAnyPermission("review.create"), h.Create)
	g.Put("/:id/visibility", middleware.RequireAnyPermission("review.moderate"), h.SetVisibility)
	g.Delete("/:id/admin", middleware.RequireAnyPermission("review.moderate"), h.AdminDelete)
	g.Delete("/:id", middleware.RequireAnyPermission("review.create"), h.Delete)
}
