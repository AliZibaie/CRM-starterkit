package httpdelivery

import (
	"errors"
	"strconv"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/internal/modules/user/domain"
	"crm-starterkit/internal/modules/user/usecase"
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

type createUserRequest struct {
	Email     string   `json:"email" validate:"required,email"`
	Password  string   `json:"password" validate:"required,min=6"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	IsActive  *bool    `json:"is_active"`
	RoleIDs   []string `json:"role_ids"`
}

type updateUserRequest struct {
	Email     *string `json:"email" validate:"omitempty,email"`
	Password  *string `json:"password" validate:"omitempty,min=6"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	IsActive  *bool   `json:"is_active"`
}

type assignRolesRequest struct {
	RoleIDs []string `json:"role_ids" validate:"required"`
}

type userResponse struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	IsActive  bool     `json:"is_active"`
	RoleIDs   []string `json:"role_ids,omitempty"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// @Summary Create user
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "user payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /api/v1/users [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req createUserRequest
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
	user, err := h.service.Create(c.Context(), domain.CreateUserInput{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  active,
		RoleIDs:   req.RoleIDs,
	})
	if err != nil {
		return mapUserError(c, err)
	}
	return response.Created(c, toUserResponse(user, req.RoleIDs))
}

// @Summary List users
// @Tags user
// @Security BearerAuth
// @Produce json
// @Param page query int false "page"
// @Param per_page query int false "per page"
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/users [get]
func (h *Handler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))
	users, total, err := h.service.List(c.Context(), page, perPage)
	if err != nil {
		return response.Internal(c, err.Error())
	}
	items := make([]userResponse, 0, len(users))
	for i := range users {
		items = append(items, toUserResponse(&users[i], nil))
	}
	return response.WithMeta(c, items, fiber.Map{
		"page":     page,
		"per_page": perPage,
		"total":    total,
	})
}

// @Summary Get user
// @Tags user
// @Security BearerAuth
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/users/{id} [get]
func (h *Handler) Get(c *fiber.Ctx) error {
	user, roleIDs, err := h.service.Get(c.Context(), c.Params("id"))
	if err != nil {
		return mapUserError(c, err)
	}
	return response.OK(c, toUserResponse(user, roleIDs))
}

// @Summary Update user
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Param body body map[string]interface{} true "user payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/users/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	var req updateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	user, err := h.service.Update(c.Context(), c.Params("id"), domain.UpdateUserInput{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  req.IsActive,
	})
	if err != nil {
		return mapUserError(c, err)
	}
	return response.OK(c, toUserResponse(user, nil))
}

// @Summary Delete user
// @Tags user
// @Security BearerAuth
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/users/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.service.Delete(c.Context(), c.Params("id")); err != nil {
		return mapUserError(c, err)
	}
	return response.Message(c, fiber.StatusOK, "user deleted", nil)
}

// @Summary Assign roles to user
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Param body body map[string]interface{} true "role_ids"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/users/{id}/roles [put]
func (h *Handler) AssignRoles(c *fiber.Ctx) error {
	var req assignRolesRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	if err := h.service.AssignRoles(c.Context(), c.Params("id"), req.RoleIDs); err != nil {
		return mapUserError(c, err)
	}
	return response.Message(c, fiber.StatusOK, "roles assigned", fiber.Map{"role_ids": req.RoleIDs})
}

func RegisterRoutes(router fiber.Router, h *Handler, jwtManager *jwtutil.Manager) {
	g := router.Group("/users", middleware.Auth(jwtManager))
	g.Post("/", middleware.RequireAnyPermission("user.create"), h.Create)
	g.Get("/", middleware.RequireAnyPermission("user.read"), h.List)
	g.Get("/:id", middleware.RequireAnyPermission("user.read"), h.Get)
	g.Put("/:id", middleware.RequireAnyPermission("user.update"), h.Update)
	g.Delete("/:id", middleware.RequireAnyPermission("user.delete"), h.Delete)
	g.Put("/:id/roles", middleware.RequireAnyPermission("user.assign_role"), h.AssignRoles)
}

func toUserResponse(u *domain.User, roleIDs []string) userResponse {
	return userResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		IsActive:  u.IsActive,
		RoleIDs:   roleIDs,
		CreatedAt: u.CreatedAt.UTC().Format(timeRFC3339),
		UpdatedAt: u.UpdatedAt.UTC().Format(timeRFC3339),
	}
}

const timeRFC3339 = "2006-01-02T15:04:05Z07:00"

func mapUserError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, usecase.ErrUserNotFound):
		return response.NotFound(c, err.Error())
	case errors.Is(err, usecase.ErrEmailAlreadyUsed):
		return response.Conflict(c, err.Error())
	case errors.Is(err, usecase.ErrInvalidUserInput):
		return response.BadRequest(c, err.Error(), nil)
	default:
		return response.Internal(c, err.Error())
	}
}
