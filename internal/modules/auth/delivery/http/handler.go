package httpdelivery

import (
	"errors"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/internal/core/validator"
	"crm-starterkit/internal/modules/auth/domain"
	"crm-starterkit/internal/modules/auth/usecase"
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

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "email password"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	pair, user, err := h.service.Login(c.Context(), domain.LoginInput{
		Email: req.Email, Password: req.Password,
	})
	if err != nil {
		return mapAuthError(c, err)
	}
	return response.OK(c, fiber.Map{
		"access_token":       pair.AccessToken,
		"refresh_token":      pair.RefreshToken,
		"access_expires_at":  pair.AccessExpiresAt,
		"refresh_expires_at": pair.RefreshExpiresAt,
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

// @Summary Refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "refresh_token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/refresh [post]
func (h *Handler) Refresh(c *fiber.Ctx) error {
	var req refreshRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body", nil)
	}
	if err := h.validate.Struct(req); err != nil {
		return response.BadRequest(c, err.Error(), nil)
	}
	pair, err := h.service.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		return mapAuthError(c, err)
	}
	return response.OK(c, fiber.Map{
		"access_token":       pair.AccessToken,
		"refresh_token":      pair.RefreshToken,
		"access_expires_at":  pair.AccessExpiresAt,
		"refresh_expires_at": pair.RefreshExpiresAt,
	})
}

// @Summary Logout
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body map[string]interface{} false "refresh_token"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/logout [post]
func (h *Handler) Logout(c *fiber.Ctx) error {
	var req logoutRequest
	_ = c.BodyParser(&req)
	userID := middleware.UserIDFromContext(c)
	if err := h.service.Logout(c.Context(), req.RefreshToken, userID); err != nil {
		return mapAuthError(c, err)
	}
	return response.Message(c, fiber.StatusOK, "logged out", nil)
}

func RegisterRoutes(router fiber.Router, h *Handler, jwtManager *jwtutil.Manager) {
	g := router.Group("/auth")
	g.Post("/login", h.Login)
	g.Post("/refresh", h.Refresh)
	g.Post("/logout", middleware.Auth(jwtManager), h.Logout)
}

func mapAuthError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, usecase.ErrInvalidCredentials),
		errors.Is(err, usecase.ErrInvalidRefresh),
		errors.Is(err, usecase.ErrRefreshRevoked):
		return response.Unauthorized(c, err.Error())
	case errors.Is(err, usecase.ErrUserInactive):
		return response.Forbidden(c, err.Error())
	default:
		return response.Internal(c, err.Error())
	}
}
