package httpdelivery

import (
	"errors"
	"time"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/modules/crm/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Handler struct{ service *usecase.Service }

func NewHandler(service *usecase.Service) *Handler { return &Handler{service: service} }

// @Summary Create lead
// @Tags crm
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "lead payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/crm/leads [post]
func (h *Handler) CreateLead(c *fiber.Ctx) error {
	var req struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Phone  string `json:"phone"`
		Source string `json:"source"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	lead, err := h.service.CreateLead(c.Context(), req.Name, req.Email, req.Phone, middleware.UserIDFromContext(c), req.Source)
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, lead)
}

// @Summary List leads
// @Tags crm
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/crm/leads [get]
func (h *Handler) ListLeads(c *fiber.Ctx) error {
	items, err := h.service.ListLeads(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary Convert lead to customer
// @Tags crm
// @Security BearerAuth
// @Produce json
// @Param id path string true "lead id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/crm/leads/{id}/convert [post]
func (h *Handler) ConvertLead(c *fiber.Ctx) error {
	contact, err := h.service.ConvertLeadToCustomer(c.Context(), c.Params("id"), middleware.UserIDFromContext(c))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, contact)
}

// @Summary List contacts
// @Tags crm
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/crm/contacts [get]
func (h *Handler) ListContacts(c *fiber.Ctx) error {
	items, err := h.service.ListContacts(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary Create deal
// @Tags crm
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "deal payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/crm/deals [post]
func (h *Handler) CreateDeal(c *fiber.Ctx) error {
	var req struct {
		Title      string `json:"title"`
		ContactID  string `json:"contact_id"`
		PipelineID string `json:"pipeline_id"`
		Stage      string `json:"stage"`
		Amount     int64  `json:"amount"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	deal, err := h.service.CreateDeal(c.Context(), req.Title, req.ContactID, req.PipelineID, req.Stage, middleware.UserIDFromContext(c), req.Amount)
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, deal)
}

// @Summary Move deal stage
// @Tags crm
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "deal id"
// @Param body body map[string]interface{} true "stage"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/crm/deals/{id}/stage [post]
func (h *Handler) MoveDealStage(c *fiber.Ctx) error {
	var req struct {
		Stage string `json:"stage"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	deal, err := h.service.MoveDealStage(c.Context(), c.Params("id"), req.Stage, middleware.UserIDFromContext(c))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, deal)
}

// @Summary Assign owner
// @Tags crm
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "entity_type entity_id owner_id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/crm/assign-owner [post]
func (h *Handler) AssignOwner(c *fiber.Ctx) error {
	var req struct {
		EntityType string `json:"entity_type"`
		EntityID   string `json:"entity_id"`
		OwnerID    string `json:"owner_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.service.AssignOwner(c.Context(), req.EntityType, req.EntityID, req.OwnerID, middleware.UserIDFromContext(c)); err != nil {
		return mapErr(c, err)
	}
	return response.Message(c, fiber.StatusOK, "assigned", nil)
}

// @Summary Log activity
// @Tags crm
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "activity payload"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/crm/activities [post]
func (h *Handler) LogActivity(c *fiber.Ctx) error {
	var req struct {
		EntityType string `json:"entity_type"`
		EntityID   string `json:"entity_id"`
		Type       string `json:"type"`
		Body       string `json:"body"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.service.LogActivity(c.Context(), req.EntityType, req.EntityID, req.Type, req.Body, middleware.UserIDFromContext(c)); err != nil {
		return response.Internal(c, err.Error())
	}
	return response.Message(c, fiber.StatusOK, "logged", nil)
}

// @Summary Create task reminder
// @Tags crm
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "task reminder payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/crm/tasks/reminders [post]
func (h *Handler) CreateTaskReminder(c *fiber.Ctx) error {
	var req struct {
		Title      string     `json:"title"`
		EntityType string     `json:"entity_type"`
		EntityID   string     `json:"entity_id"`
		DueAt      *time.Time `json:"due_at"`
		ReminderAt *time.Time `json:"reminder_at"`
		Notes      string     `json:"notes"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	reminderAt := req.ReminderAt
	if reminderAt == nil {
		reminderAt = req.DueAt
	}
	actorID := middleware.UserIDFromContext(c)
	task, err := h.service.CreateTaskReminder(c.Context(), req.Title, req.EntityType, req.EntityID, actorID, req.DueAt, reminderAt)
	if err != nil {
		return mapErr(c, err)
	}
	if req.Notes != "" && req.EntityType != "" && req.EntityID != "" {
		_ = h.service.LogActivity(c.Context(), req.EntityType, req.EntityID, "note", req.Notes, actorID)
	}
	return response.Created(c, task)
}

// @Summary Customer 360 view
// @Tags crm
// @Security BearerAuth
// @Produce json
// @Param contact_id path string true "contact id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/crm/customers/{contact_id}/360 [get]
func (h *Handler) Customer360(c *fiber.Ctx) error {
	view, err := h.service.Customer360(c.Context(), c.Params("contact_id"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, view)
}

// @Summary Upsert pipeline
// @Tags crm
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "name stages_json"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/crm/pipelines [post]
func (h *Handler) UpsertPipeline(c *fiber.Ctx) error {
	var req struct {
		Name       string `json:"name"`
		StagesJSON string `json:"stages_json"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	p, err := h.service.UpsertPipeline(c.Context(), req.Name, req.StagesJSON)
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, p)
}

// @Summary List pipelines
// @Tags crm
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/crm/pipelines [get]
func (h *Handler) ListPipelines(c *fiber.Ctx) error {
	items, err := h.service.ListPipelines(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary Upsert tag
// @Tags crm
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "name"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/crm/tags [post]
func (h *Handler) UpsertTag(c *fiber.Ctx) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	t, err := h.service.UpsertTag(c.Context(), req.Name)
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, t)
}

// @Summary Attach tag
// @Tags crm
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "entity_type entity_id tag_id"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/crm/tags/attach [post]
func (h *Handler) AttachTag(c *fiber.Ctx) error {
	var req struct {
		EntityType string `json:"entity_type"`
		EntityID   string `json:"entity_id"`
		TagID      string `json:"tag_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.service.AttachTag(c.Context(), req.EntityType, req.EntityID, req.TagID); err != nil {
		return response.Internal(c, err.Error())
	}
	return response.Message(c, fiber.StatusOK, "tagged", nil)
}

// @Summary Upsert segment
// @Tags crm
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "name rules_json"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/crm/segments [post]
func (h *Handler) UpsertSegment(c *fiber.Ctx) error {
	var req struct {
		Name      string `json:"name"`
		RulesJSON string `json:"rules_json"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	seg, err := h.service.UpsertSegment(c.Context(), req.Name, req.RulesJSON)
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, seg)
}

// @Summary List segments
// @Tags crm
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/crm/segments [get]
func (h *Handler) ListSegments(c *fiber.Ctx) error {
	items, err := h.service.ListSegments(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

func (h *Handler) ListDeals(c *fiber.Ctx) error {
	items, err := h.service.ListDealsByPipeline(c.Context(), c.Query("pipeline_id"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, items)
}

func (h *Handler) CreateContact(c *fiber.Ctx) error {
	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Company string `json:"company"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	contact, err := h.service.CreateContact(c.Context(), req.Name, req.Email, req.Phone, req.Company, middleware.UserIDFromContext(c))
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, contact)
}

func (h *Handler) UpdateContact(c *fiber.Ctx) error {
	var req struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Phone   string `json:"phone"`
		Company string `json:"company"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	contact, err := h.service.UpdateContact(c.Context(), c.Params("id"), req.Name, req.Email, req.Phone, req.Company)
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, contact)
}

func (h *Handler) DeleteContact(c *fiber.Ctx) error {
	if err := h.service.DeleteContact(c.Context(), c.Params("id")); err != nil {
		return mapErr(c, err)
	}
	return response.Message(c, fiber.StatusOK, "contact deleted", nil)
}

func (h *Handler) ListActivities(c *fiber.Ctx) error {
	items, err := h.service.ListActivities(c.Context(), c.Query("entity_type"), c.Query("entity_id"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, items)
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/crm", middleware.Auth(jwt))
	g.Post("/leads", middleware.RequireAnyPermission("crm.lead"), h.CreateLead)
	g.Get("/leads", middleware.RequireAnyPermission("crm.lead"), h.ListLeads)
	g.Post("/leads/:id/convert", middleware.RequireAnyPermission("crm.lead"), h.ConvertLead)
	g.Get("/contacts", middleware.RequireAnyPermission("crm.contact"), h.ListContacts)
	g.Post("/contacts", middleware.RequireAnyPermission("crm.contact"), h.CreateContact)
	g.Put("/contacts/:id", middleware.RequireAnyPermission("crm.contact"), h.UpdateContact)
	g.Delete("/contacts/:id", middleware.RequireAnyPermission("crm.contact"), h.DeleteContact)
	g.Get("/deals", middleware.RequireAnyPermission("crm.deal"), h.ListDeals)
	g.Post("/deals", middleware.RequireAnyPermission("crm.deal"), h.CreateDeal)
	g.Post("/deals/:id/stage", middleware.RequireAnyPermission("crm.deal"), h.MoveDealStage)
	g.Get("/activities", middleware.RequireAnyPermission("crm.manage"), h.ListActivities)
	g.Post("/assign-owner", middleware.RequireAnyPermission("crm.manage"), h.AssignOwner)
	g.Post("/activities", middleware.RequireAnyPermission("crm.manage"), h.LogActivity)
	g.Post("/tasks/reminders", middleware.RequireAnyPermission("crm.task"), h.CreateTaskReminder)
	g.Get("/customers/:contact_id/360", middleware.RequireAnyPermission("crm.view360"), h.Customer360)
	g.Post("/pipelines", middleware.RequireAnyPermission("crm.manage"), h.UpsertPipeline)
	g.Get("/pipelines", middleware.RequireAnyPermission("crm.manage"), h.ListPipelines)
	g.Post("/tags", middleware.RequireAnyPermission("crm.manage"), h.UpsertTag)
	g.Post("/tags/attach", middleware.RequireAnyPermission("crm.manage"), h.AttachTag)
	g.Post("/segments", middleware.RequireAnyPermission("crm.manage"), h.UpsertSegment)
	g.Get("/segments", middleware.RequireAnyPermission("crm.manage"), h.ListSegments)
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
