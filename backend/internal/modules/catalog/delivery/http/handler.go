package httpdelivery

import (
	"errors"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/internal/modules/catalog/usecase"
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

// @Summary Create category
// @Tags catalog
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "name parent_id"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/catalog/categories [post]
func (h *Handler) CreateCategory(c *fiber.Ctx) error {
	var req struct {
		Name     string `json:"name" validate:"required"`
		ParentID string `json:"parent_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	item, err := h.service.CreateCategory(c.Context(), req.Name, req.ParentID)
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, item)
}

// @Summary List categories
// @Tags catalog
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/catalog/categories [get]
func (h *Handler) ListCategories(c *fiber.Ctx) error {
	items, err := h.service.ListCategories(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary Create catalog item
// @Tags catalog
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "item payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/catalog/items [post]
func (h *Handler) CreateItem(c *fiber.Ctx) error {
	var req struct {
		Type        string `json:"type" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		CategoryID  string `json:"category_id"`
		Currency    string `json:"currency"`
		BasePrice   int64  `json:"base_price"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	item, err := h.service.CreateItem(c.Context(), req.Type, req.Name, req.Description, req.CategoryID, req.Currency, req.BasePrice)
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, item)
}

// @Summary List catalog items
// @Tags catalog
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/catalog/items [get]
func (h *Handler) ListItems(c *fiber.Ctx) error {
	items, err := h.service.ListItems(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary Get catalog item
// @Tags catalog
// @Security BearerAuth
// @Produce json
// @Param id path string true "item id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/catalog/items/{id} [get]
func (h *Handler) GetItem(c *fiber.Ctx) error {
	item, attrs, variants, err := h.service.GetItem(c.Context(), c.Params("id"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, fiber.Map{"item": item, "attributes": attrs, "variants": variants})
}

// @Summary Set item attribute
// @Tags catalog
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "item id"
// @Param body body map[string]interface{} true "key value"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/catalog/items/{id}/attributes [post]
func (h *Handler) SetAttribute(c *fiber.Ctx) error {
	var req struct {
		Key   string `json:"key" validate:"required"`
		Value string `json:"value"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	a, err := h.service.SetAttribute(c.Context(), c.Params("id"), req.Key, req.Value)
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, a)
}

// @Summary Create item variant
// @Tags catalog
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "item id"
// @Param body body map[string]interface{} true "variant payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/catalog/items/{id}/variants [post]
func (h *Handler) CreateVariant(c *fiber.Ctx) error {
	var req struct {
		SKU       string `json:"sku" validate:"required"`
		Name      string `json:"name" validate:"required"`
		Price     int64  `json:"price"`
		AttrsJSON string `json:"attrs_json"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	v, err := h.service.CreateVariant(c.Context(), c.Params("id"), req.SKU, req.Name, req.AttrsJSON, req.Price)
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, v)
}

func (h *Handler) UpdateItem(c *fiber.Ctx) error {
	var req struct {
		Type        string `json:"type" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		CategoryID  string `json:"category_id"`
		Currency    string `json:"currency"`
		BasePrice   int64  `json:"base_price"`
		IsActive    *bool  `json:"is_active"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	item, err := h.service.UpdateItem(c.Context(), c.Params("id"), req.Type, req.Name, req.Description, req.CategoryID, req.Currency, req.BasePrice, active, middleware.UserIDFromContext(c))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, item)
}

func (h *Handler) DeleteItem(c *fiber.Ctx) error {
	if err := h.service.DeleteItem(c.Context(), c.Params("id")); err != nil {
		return mapErr(c, err)
	}
	return response.Message(c, fiber.StatusOK, "item deleted", nil)
}

func (h *Handler) UpdateCategory(c *fiber.Ctx) error {
	var req struct {
		Name     string `json:"name" validate:"required"`
		ParentID string `json:"parent_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	item, err := h.service.UpdateCategory(c.Context(), c.Params("id"), req.Name, req.ParentID)
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, item)
}

func (h *Handler) DeleteCategory(c *fiber.Ctx) error {
	if err := h.service.DeleteCategory(c.Context(), c.Params("id")); err != nil {
		return mapErr(c, err)
	}
	return response.Message(c, fiber.StatusOK, "category deleted", nil)
}

// @Summary Set inventory
// @Tags catalog
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param variant_id path string true "variant id"
// @Param body body map[string]interface{} true "quantity"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/catalog/inventory/{variant_id} [put]
func (h *Handler) SetInventory(c *fiber.Ctx) error {
	var req struct {
		Quantity int64 `json:"quantity"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	inv, err := h.service.SetInventory(c.Context(), c.Params("variant_id"), req.Quantity)
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, inv)
}

// @Summary Get inventory
// @Tags catalog
// @Security BearerAuth
// @Produce json
// @Param variant_id path string true "variant id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/catalog/inventory/{variant_id} [get]
func (h *Handler) GetInventory(c *fiber.Ctx) error {
	inv, err := h.service.GetInventory(c.Context(), c.Params("variant_id"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, inv)
}

func (h *Handler) ListPublicItems(c *fiber.Ctx) error {
	items, err := h.service.ListPublicItems(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

func (h *Handler) GetPublicItem(c *fiber.Ctx) error {
	item, attrs, variants, err := h.service.GetItemBySlug(c.Context(), c.Params("slug"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, fiber.Map{"item": item, "attributes": attrs, "variants": variants})
}

func (h *Handler) ListPublicCategories(c *fiber.Ctx) error {
	items, err := h.service.ListCategories(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

func (h *Handler) GetPublicCategory(c *fiber.Ctx) error {
	cat, items, err := h.service.GetCategoryBySlug(c.Context(), c.Params("slug"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, fiber.Map{"category": cat, "items": items})
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	pub := router.Group("/catalog/public")
	pub.Get("/items", h.ListPublicItems)
	pub.Get("/items/:slug", h.GetPublicItem)
	pub.Get("/categories", h.ListPublicCategories)
	pub.Get("/categories/:slug", h.GetPublicCategory)

	g := router.Group("/catalog", middleware.Auth(jwt))
	g.Post("/categories", middleware.RequireAnyPermission("catalog.category.manage"), h.CreateCategory)
	g.Get("/categories", middleware.RequireAnyPermission("catalog.read"), h.ListCategories)
	g.Post("/items", middleware.RequireAnyPermission("catalog.item.manage"), h.CreateItem)
	g.Get("/items", middleware.RequireAnyPermission("catalog.read"), h.ListItems)
	g.Get("/items/:id", middleware.RequireAnyPermission("catalog.read"), h.GetItem)
	g.Put("/items/:id", middleware.RequireAnyPermission("catalog.item.manage"), h.UpdateItem)
	g.Delete("/items/:id", middleware.RequireAnyPermission("catalog.item.manage"), h.DeleteItem)
	g.Post("/items/:id/attributes", middleware.RequireAnyPermission("catalog.item.manage"), h.SetAttribute)
	g.Post("/items/:id/variants", middleware.RequireAnyPermission("catalog.item.manage"), h.CreateVariant)
	g.Put("/categories/:id", middleware.RequireAnyPermission("catalog.category.manage"), h.UpdateCategory)
	g.Delete("/categories/:id", middleware.RequireAnyPermission("catalog.category.manage"), h.DeleteCategory)
	g.Put("/inventory/:variant_id", middleware.RequireAnyPermission("catalog.inventory.manage"), h.SetInventory)
	g.Get("/inventory/:variant_id", middleware.RequireAnyPermission("catalog.read"), h.GetInventory)
}

func mapErr(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		return response.NotFound(c, err.Error())
	case errors.Is(err, usecase.ErrInvalidInput):
		return response.BadRequest(c, err.Error(), nil)
	case errors.Is(err, usecase.ErrInUse):
		return response.Conflict(c, err.Error())
	default:
		return response.Internal(c, err.Error())
	}
}
