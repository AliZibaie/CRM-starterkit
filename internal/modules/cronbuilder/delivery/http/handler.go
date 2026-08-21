package httpdelivery

import (
	"errors"
	"strconv"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/internal/modules/cronbuilder/domain"
	"crm-starterkit/internal/modules/cronbuilder/usecase"
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

type createJobRequest struct {
	Name      string                `json:"name" validate:"required"`
	TableName string                `json:"table_name" validate:"required"`
	Condition domain.ConditionNode  `json:"condition"`
	Action    domain.JobAction      `json:"action" validate:"required"`
	CronExpr  string                `json:"cron_expr" validate:"required"`
	IsActive  *bool                 `json:"is_active"`
}

type buildConditionRequest struct {
	TableName string               `json:"table_name" validate:"required"`
	Condition domain.ConditionNode `json:"condition"`
}

// @Summary List tables
// @Tags cronbuilder
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/cronbuilder/tables [get]
func (h *Handler) ListTables(c *fiber.Ctx) error {
	return response.OK(c, h.service.ListTables())
}

// @Summary List table columns
// @Tags cronbuilder
// @Security BearerAuth
// @Produce json
// @Param table path string true "table name"
// @Success 200 {array} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/cronbuilder/tables/{table}/columns [get]
func (h *Handler) ListColumns(c *fiber.Ctx) error {
	cols, err := h.service.ListColumns(c.Params("table"))
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.OK(c, cols)
}

// @Summary List actions
// @Tags cronbuilder
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/cronbuilder/actions [get]
func (h *Handler) ListActions(c *fiber.Ctx) error {
	return response.OK(c, h.service.ListActions())
}

// @Summary Build condition
// @Tags cronbuilder
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "table_name condition"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/cronbuilder/conditions/build [post]
func (h *Handler) BuildCondition(c *fiber.Ctx) error {
	var req buildConditionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	built, err := h.service.BuildCondition(req.TableName, req.Condition)
	if err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	return response.OK(c, built)
}

// @Summary Create cron job
// @Tags cronbuilder
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "job payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/cronbuilder/jobs [post]
func (h *Handler) CreateJob(c *fiber.Ctx) error {
	var req createJobRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	job, err := h.service.CreateCronJob(c.Context(), domain.CreateJobInput{
		Name: req.Name, TableName: req.TableName, Condition: req.Condition,
		Action: req.Action, CronExpr: req.CronExpr, IsActive: active,
	})
	if err != nil {
		return mapCronError(c, err)
	}
	return response.Created(c, job)
}

// @Summary List cron jobs
// @Tags cronbuilder
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/cronbuilder/jobs [get]
func (h *Handler) ListJobs(c *fiber.Ctx) error {
	jobs, err := h.service.ListJobs(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, jobs)
}

// @Summary Run cron job now
// @Tags cronbuilder
// @Security BearerAuth
// @Produce json
// @Param id path string true "job id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/cronbuilder/jobs/{id}/run [post]
func (h *Handler) RunNow(c *fiber.Ctx) error {
	log, err := h.service.RunJobNow(c.Context(), c.Params("id"))
	if err != nil {
		return mapCronError(c, err)
	}
	return response.OK(c, log)
}

// @Summary Job execution logs
// @Tags cronbuilder
// @Security BearerAuth
// @Produce json
// @Param id path string true "job id"
// @Param limit query int false "limit"
// @Success 200 {array} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/cronbuilder/jobs/{id}/logs [get]
func (h *Handler) Logs(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	logs, err := h.service.JobExecutionLogs(c.Context(), c.Params("id"), limit)
	if err != nil {
		return mapCronError(c, err)
	}
	return response.OK(c, logs)
}

func RegisterRoutes(router fiber.Router, h *Handler, jwt *jwtutil.Manager) {
	g := router.Group("/cronbuilder", middleware.Auth(jwt))
	g.Get("/tables", middleware.RequireAnyPermission("cronbuilder.read"), h.ListTables)
	g.Get("/tables/:table/columns", middleware.RequireAnyPermission("cronbuilder.read"), h.ListColumns)
	g.Get("/actions", middleware.RequireAnyPermission("cronbuilder.read"), h.ListActions)
	g.Post("/conditions/build", middleware.RequireAnyPermission("cronbuilder.create"), h.BuildCondition)
	g.Post("/jobs", middleware.RequireAnyPermission("cronbuilder.create"), h.CreateJob)
	g.Get("/jobs", middleware.RequireAnyPermission("cronbuilder.read"), h.ListJobs)
	g.Post("/jobs/:id/run", middleware.RequireAnyPermission("cronbuilder.run"), h.RunNow)
	g.Get("/jobs/:id/logs", middleware.RequireAnyPermission("cronbuilder.read"), h.Logs)
}

func mapCronError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, usecase.ErrJobNotFound):
		return response.NotFound(c, err.Error())
	case errors.Is(err, usecase.ErrInvalidInput), errors.Is(err, usecase.ErrInvalidCron),
		errors.Is(err, usecase.ErrInvalidAction), errors.Is(err, usecase.ErrInvalidTable):
		return response.BadRequest(c, err.Error(), nil)
	default:
		return response.Internal(c, err.Error())
	}
}
