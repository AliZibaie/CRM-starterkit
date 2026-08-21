package httpdelivery

import (
	"errors"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/module"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/internal/modules/rbac/domain"
	"crm-starterkit/internal/modules/rbac/usecase"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service  *usecase.Service
	validate *validator.Validator
	registry *module.Registry
}

func NewHandler(service *usecase.Service, validate *validator.Validator, registry *module.Registry) *Handler {
	return &Handler{service: service, validate: validate, registry: registry}
}

type createRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type createPermissionRequest struct {
	Key         string `json:"key" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Module      string `json:"module"`
}

type editPermissionRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Module      *string `json:"module"`
}

type assignPermissionRequest struct {
	RoleID       string `json:"role_id" validate:"required"`
	PermissionID string `json:"permission_id" validate:"required"`
}

type assignRoleRequest struct {
	UserID string `json:"user_id" validate:"required"`
	RoleID string `json:"role_id" validate:"required"`
}

// @Summary Create role
// @Tags rbac
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "role payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /api/v1/rbac/roles [post]
func (h *Handler) CreateRole(c *fiber.Ctx) error {
	var req createRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	role, err := h.service.CreateRole(c.Context(), domain.CreateRoleInput{
		Name: req.Name, Description: req.Description,
	})
	if err != nil {
		return mapRBACError(c, err)
	}
	return response.Created(c, role)
}

// @Summary List roles
// @Tags rbac
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/rbac/roles [get]
func (h *Handler) ListRoles(c *fiber.Ctx) error {
	roles, err := h.service.ListRoles(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, roles)
}

// @Summary Get role
// @Tags rbac
// @Security BearerAuth
// @Produce json
// @Param id path string true "role id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/rbac/roles/{id} [get]
func (h *Handler) GetRole(c *fiber.Ctx) error {
	role, permIDs, err := h.service.GetRole(c.Context(), c.Params("id"))
	if err != nil {
		return mapRBACError(c, err)
	}
	return response.OK(c, fiber.Map{"role": role, "permission_ids": permIDs})
}

// @Summary Create permission
// @Tags rbac
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "permission payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /api/v1/rbac/permissions [post]
func (h *Handler) CreatePermission(c *fiber.Ctx) error {
	var req createPermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	perm, err := h.service.CreatePermission(c.Context(), domain.CreatePermissionInput{
		Key: req.Key, Name: req.Name, Description: req.Description, Module: req.Module,
	})
	if err != nil {
		return mapRBACError(c, err)
	}
	return response.Created(c, perm)
}

// @Summary Edit permission
// @Tags rbac
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "permission id"
// @Param body body map[string]interface{} true "permission payload"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/rbac/permissions/{id} [put]
func (h *Handler) EditPermission(c *fiber.Ctx) error {
	var req editPermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	perm, err := h.service.EditPermission(c.Context(), c.Params("id"), domain.EditPermissionInput{
		Name: req.Name, Description: req.Description, Module: req.Module,
	})
	if err != nil {
		return mapRBACError(c, err)
	}
	return response.OK(c, perm)
}

// @Summary List permissions
// @Tags rbac
// @Security BearerAuth
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/rbac/permissions [get]
func (h *Handler) ListPermissions(c *fiber.Ctx) error {
	perms, err := h.service.ListPermissions(c.Context())
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, perms)
}

// @Summary Assign permission to role
// @Tags rbac
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "role_id permission_id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/rbac/assign-permission [post]
func (h *Handler) AssignPermissionToRole(c *fiber.Ctx) error {
	var req assignPermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	if err := h.service.AssignPermissionToRole(c.Context(), req.RoleID, req.PermissionID); err != nil {
		return mapRBACError(c, err)
	}
	return response.Message(c, fiber.StatusOK, "permission assigned to role", nil)
}

// @Summary Assign role to user
// @Tags rbac
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "user_id role_id"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/rbac/assign-role [post]
func (h *Handler) AssignRoleToUser(c *fiber.Ctx) error {
	var req assignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	if err := h.service.AssignRoleToUser(c.Context(), req.UserID, req.RoleID); err != nil {
		return mapRBACError(c, err)
	}
	return response.Message(c, fiber.StatusOK, "role assigned to user", nil)
}

// @Summary Sync permissions from registry
// @Tags rbac
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rbac/permissions/sync [post]
func (h *Handler) SyncPermissions(c *fiber.Ctx) error {
	var perms []module.Permission
	if h.registry != nil {
		perms = h.registry.Permissions()
	}
	count, err := h.service.SyncPermissionsFromRegistry(c.Context(), perms)
	if err != nil {
		return response.Internal(c, err.Error())
	}
	return response.OK(c, fiber.Map{"synced": count})
}

func RegisterRoutes(router fiber.Router, h *Handler, jwtManager *jwtutil.Manager) {
	g := router.Group("/rbac", middleware.Auth(jwtManager))
	g.Post("/roles", middleware.RequireAnyPermission("rbac.role.create"), h.CreateRole)
	g.Get("/roles", middleware.RequireAnyPermission("rbac.role.read"), h.ListRoles)
	g.Get("/roles/:id", middleware.RequireAnyPermission("rbac.role.read"), h.GetRole)
	g.Post("/permissions", middleware.RequireAnyPermission("rbac.permission.create"), h.CreatePermission)
	g.Get("/permissions", middleware.RequireAnyPermission("rbac.permission.read"), h.ListPermissions)
	g.Put("/permissions/:id", middleware.RequireAnyPermission("rbac.permission.edit"), h.EditPermission)
	g.Post("/assign-permission", middleware.RequireAnyPermission("rbac.assign.permission"), h.AssignPermissionToRole)
	g.Post("/assign-role", middleware.RequireAnyPermission("rbac.assign.role"), h.AssignRoleToUser)
	g.Post("/permissions/sync", middleware.RequireAnyPermission("rbac.permission.sync"), h.SyncPermissions)
}

func mapRBACError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, usecase.ErrRoleNotFound), errors.Is(err, usecase.ErrPermissionNotFound):
		return response.NotFound(c, err.Error())
	case errors.Is(err, usecase.ErrRoleExists), errors.Is(err, usecase.ErrPermissionExists):
		return response.Conflict(c, err.Error())
	case errors.Is(err, usecase.ErrInvalidInput):
		return response.BadRequest(c, err.Error(), nil)
	default:
		return response.Internal(c, err.Error())
	}
}
