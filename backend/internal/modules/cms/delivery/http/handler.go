package httpdelivery

import (
	"errors"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/internal/modules/cms/domain"
	"crm-starterkit/internal/modules/cms/usecase"
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

// @Summary Create CMS page
// @Tags cms
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "title slug"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/cms/pages [post]
func (h *Handler) CreatePage(c *fiber.Ctx) error {
	var req struct {
		Title string `json:"title" validate:"required"`
		Slug  string `json:"slug"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	p, err := h.service.CreatePage(c.Context(), req.Title, req.Slug)
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, p)
}

// @Summary List CMS pages
// @Tags cms
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/cms/pages [get]
func (h *Handler) ListPages(c *fiber.Ctx) error {
	items, err := h.service.ListPages(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary Get CMS page by slug
// @Tags cms
// @Produce json
// @Param slug path string true "slug"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/cms/pages/{slug} [get]
func (h *Handler) GetPage(c *fiber.Ctx) error {
	p, blocks, err := h.service.GetPage(c.Context(), c.Params("slug"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, fiber.Map{"page": p, "blocks": blocks})
}

// @Summary Set page blocks
// @Tags cms
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "page id"
// @Param body body map[string]interface{} true "blocks"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/cms/pages/{id}/blocks [put]
func (h *Handler) SetBlocks(c *fiber.Ctx) error {
	var req struct {
		Blocks []domain.Block `json:"blocks"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.service.SetBlocks(c.Context(), c.Params("id"), req.Blocks); err != nil {
		return response.Internal(c, err.Error())
	}
	return response.Message(c, fiber.StatusOK, "blocks saved", nil)
}

// @Summary Upsert menu
// @Tags cms
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "menu payload"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/cms/menus [post]
func (h *Handler) UpsertMenu(c *fiber.Ctx) error {
	var req struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Location  string `json:"location"`
		ItemsJSON string `json:"items_json"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	m, err := h.service.UpsertMenu(c.Context(), req.ID, req.Name, req.Location, req.ItemsJSON)
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, m)
}

// @Summary List menus
// @Tags cms
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/cms/menus [get]
func (h *Handler) ListMenus(c *fiber.Ctx) error {
	items, err := h.service.ListMenus(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary Upsert SEO meta
// @Tags cms
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "seo payload"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/cms/seo [post]
func (h *Handler) UpsertSEO(c *fiber.Ctx) error {
	var req struct {
		EntityType  string `json:"entity_type"`
		EntityID    string `json:"entity_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Keywords    string `json:"keywords"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	meta, err := h.service.UpsertSEO(c.Context(), req.EntityType, req.EntityID, req.Title, req.Description, req.Keywords)
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, meta)
}

func (h *Handler) GetPageByID(c *fiber.Ctx) error {
	p, blocks, err := h.service.GetPageByID(c.Context(), c.Params("id"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, fiber.Map{"page": p, "blocks": blocks})
}

func (h *Handler) UpdatePage(c *fiber.Ctx) error {
	var req struct {
		Title  string `json:"title"`
		Slug   string `json:"slug"`
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	p, err := h.service.UpdatePage(c.Context(), c.Params("id"), req.Title, req.Slug, req.Status)
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, p)
}

// @Summary Get SEO meta
// @Tags cms
// @Produce json
// @Param entity_type query string false "entity type"
// @Param entity_id query string false "entity id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/cms/seo [get]
func (h *Handler) GetSEO(c *fiber.Ctx) error {
	meta, err := h.service.GetSEO(c.Context(), c.Query("entity_type"), c.Query("entity_id"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, meta)
}

// @Summary Upsert FAQ
// @Tags cms
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "faq payload"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/cms/faqs [post]
func (h *Handler) UpsertFAQ(c *fiber.Ctx) error {
	var req struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
		Sort     int    `json:"sort_order"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	f, err := h.service.UpsertFAQ(c.Context(), req.Question, req.Answer, req.Sort)
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, f)
}

// @Summary List FAQs
// @Tags cms
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/cms/faqs [get]
func (h *Handler) ListFAQs(c *fiber.Ctx) error {
	items, err := h.service.ListFAQs(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	pub := router.Group("/cms")
	pub.Get("/pages/:slug", h.GetPage)
	pub.Get("/faqs", h.ListFAQs)
	pub.Get("/seo", h.GetSEO)
	pub.Get("/menus", h.ListMenus)

	g := router.Group("/cms", middleware.Auth(jwt))
	g.Post("/pages", middleware.RequireAnyPermission("cms.manage"), h.CreatePage)
	g.Get("/pages", middleware.RequireAnyPermission("cms.manage"), h.ListPages)
	g.Get("/pages/id/:id", middleware.RequireAnyPermission("cms.manage"), h.GetPageByID)
	g.Put("/pages/:id", middleware.RequireAnyPermission("cms.manage"), h.UpdatePage)
	g.Put("/pages/:id/blocks", middleware.RequireAnyPermission("cms.manage"), h.SetBlocks)
	g.Post("/menus", middleware.RequireAnyPermission("cms.manage"), h.UpsertMenu)
	g.Get("/menus/all", middleware.RequireAnyPermission("cms.manage"), h.ListMenus)
	g.Post("/seo", middleware.RequireAnyPermission("cms.manage"), h.UpsertSEO)
	g.Post("/faqs", middleware.RequireAnyPermission("cms.manage"), h.UpsertFAQ)
}

func mapErr(c *fiber.Ctx, err error) error {
	if errors.Is(err, usecase.ErrNotFound) {
		return response.NotFound(c, err.Error())
	}
	if errors.Is(err, usecase.ErrInvalidInput) {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.Internal(c, err.Error())
}
