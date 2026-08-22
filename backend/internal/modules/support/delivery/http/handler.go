package httpdelivery

import (
	"errors"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/internal/modules/support/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service  *usecase.Service
	validate *validator.Validator
}

func NewHandler(service *usecase.Service, validate *validator.Validator) *Handler {
	return &Handler{service: service, validate: validate}
}

type createDeptRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type createTicketRequest struct {
	Subject      string `json:"subject" validate:"required"`
	Priority     string `json:"priority"`
	DepartmentID string `json:"department_id"`
}

type replyRequest struct {
	Body       string `json:"body" validate:"required"`
	IsInternal bool   `json:"is_internal"`
}

type assignRequest struct {
	AgentID string `json:"agent_id" validate:"required"`
}

type statusRequest struct {
	Status string `json:"status" validate:"required"`
}

// @Summary Create department
// @Tags support
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "name description"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/support/departments [post]
func (h *Handler) CreateDepartment(c *fiber.Ctx) error {
	var req createDeptRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	d, err := h.service.CreateDepartment(c.Context(), req.Name, req.Description)
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, d)
}

// @Summary List departments
// @Tags support
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/support/departments [get]
func (h *Handler) ListDepartments(c *fiber.Ctx) error {
	items, err := h.service.ListDepartments(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary Create ticket
// @Tags support
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "ticket payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/support/tickets [post]
func (h *Handler) CreateTicket(c *fiber.Ctx) error {
	var req createTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	t, err := h.service.CreateTicket(c.Context(), req.Subject, req.Priority, req.DepartmentID, middleware.UserIDFromContext(c))
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, t)
}

// @Summary List tickets
// @Tags support
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/support/tickets [get]
func (h *Handler) ListTickets(c *fiber.Ctx) error {
	items, err := h.service.ListTickets(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, items)
}

// @Summary Get ticket
// @Tags support
// @Security BearerAuth
// @Produce json
// @Param id path string true "ticket id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/support/tickets/{id} [get]
func (h *Handler) GetTicket(c *fiber.Ctx) error {
	t, msgs, err := h.service.GetTicket(c.Context(), c.Params("id"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, fiber.Map{"ticket": t, "messages": msgs})
}

// @Summary Reply to ticket
// @Tags support
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ticket id"
// @Param body body map[string]interface{} true "body is_internal"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/support/tickets/{id}/reply [post]
func (h *Handler) Reply(c *fiber.Ctx) error {
	var req replyRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	msg, err := h.service.Reply(c.Context(), c.Params("id"), middleware.UserIDFromContext(c), req.Body, req.IsInternal)
	if err != nil {
		return mapErr(c, err)
	}
	return response.Created(c, msg)
}

// @Summary Assign ticket agent
// @Tags support
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ticket id"
// @Param body body map[string]interface{} true "agent_id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/support/tickets/{id}/assign [post]
func (h *Handler) Assign(c *fiber.Ctx) error {
	var req assignRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	t, err := h.service.AssignAgent(c.Context(), c.Params("id"), req.AgentID)
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, t)
}

// @Summary Change ticket status
// @Tags support
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ticket id"
// @Param body body map[string]interface{} true "status"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/support/tickets/{id}/status [put]
func (h *Handler) ChangeStatus(c *fiber.Ctx) error {
	var req statusRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	t, err := h.service.ChangeStatus(c.Context(), c.Params("id"), req.Status)
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, t)
}

// @Summary Close ticket
// @Tags support
// @Security BearerAuth
// @Produce json
// @Param id path string true "ticket id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/support/tickets/{id}/close [post]
func (h *Handler) Close(c *fiber.Ctx) error {
	t, err := h.service.CloseTicket(c.Context(), c.Params("id"))
	if err != nil {
		return mapErr(c, err)
	}
	return response.OK(c, t)
}

// @Summary WebSocket upgrade
// @Tags support
// @Param id path string true "ticket id"
// @Success 101 {string} string "Switching Protocols"
// @Failure 426 {object} map[string]interface{}
// @Router /api/v1/support/ws/{id} [get]
func (h *Handler) WSUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (h *Handler) WSHandler() fiber.Handler {
	hub := h.service.Hub()
	return websocket.New(func(conn *websocket.Conn) {
		ticketID := conn.Params("id")
		hub.Join(ticketID, conn)
		defer hub.Leave(ticketID, conn)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	})
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/support", middleware.Auth(jwt))
	g.Post("/departments", middleware.RequireAnyPermission("support.department.manage"), h.CreateDepartment)
	g.Get("/departments", middleware.RequireAnyPermission("support.read"), h.ListDepartments)
	g.Post("/tickets", middleware.RequireAnyPermission("support.ticket.create"), h.CreateTicket)
	g.Get("/tickets", middleware.RequireAnyPermission("support.read"), h.ListTickets)
	g.Get("/tickets/:id", middleware.RequireAnyPermission("support.read"), h.GetTicket)
	g.Post("/tickets/:id/reply", middleware.RequireAnyPermission("support.ticket.reply"), h.Reply)
	g.Post("/tickets/:id/assign", middleware.RequireAnyPermission("support.ticket.assign"), h.Assign)
	g.Put("/tickets/:id/status", middleware.RequireAnyPermission("support.ticket.status"), h.ChangeStatus)
	g.Post("/tickets/:id/close", middleware.RequireAnyPermission("support.ticket.close"), h.Close)

	router.Use("/support/ws/:id", h.WSUpgrade)
	router.Get("/support/ws/:id", h.WSHandler())
}

func mapErr(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		return response.NotFound(c, err.Error())
	case errors.Is(err, usecase.ErrInvalidInput), errors.Is(err, usecase.ErrClosedTicket):
		return response.BadRequest(c, err.Error(), nil)
	default:
		return response.Internal(c, err.Error())
	}
}
