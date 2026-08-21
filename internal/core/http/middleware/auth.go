package middleware

import (
	"strings"

	"crm-starterkit/internal/core/http/response"
	"crm-starterkit/pkg/jwtutil"

	"github.com/gofiber/fiber/v2"
)

const (
	LocalsUserID      = "user_id"
	LocalsEmail       = "email"
	LocalsRoles       = "roles"
	LocalsPermissions = "permissions"
	LocalsClaims      = "claims"
)

func Auth(jwtManager *jwtutil.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return response.Unauthorized(c, "missing authorization header")
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			return response.Unauthorized(c, "invalid authorization header")
		}
		claims, err := jwtManager.Parse(parts[1])
		if err != nil {
			return response.Unauthorized(c, "invalid or expired token")
		}
		if claims.TokenType != "access" {
			return response.Unauthorized(c, "access token required")
		}
		c.Locals(LocalsUserID, claims.UserID)
		c.Locals(LocalsEmail, claims.Email)
		c.Locals(LocalsRoles, claims.Roles)
		c.Locals(LocalsPermissions, claims.Permissions)
		c.Locals(LocalsClaims, claims)
		return c.Next()
	}
}

func OptionalAuth(jwtManager *jwtutil.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Next()
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			return c.Next()
		}
		claims, err := jwtManager.Parse(parts[1])
		if err != nil || claims.TokenType != "access" {
			return c.Next()
		}
		c.Locals(LocalsUserID, claims.UserID)
		c.Locals(LocalsEmail, claims.Email)
		c.Locals(LocalsRoles, claims.Roles)
		c.Locals(LocalsPermissions, claims.Permissions)
		c.Locals(LocalsClaims, claims)
		return c.Next()
	}
}

func UserIDFromContext(c *fiber.Ctx) string {
	v, _ := c.Locals(LocalsUserID).(string)
	return v
}

func PermissionsFromContext(c *fiber.Ctx) []string {
	v, _ := c.Locals(LocalsPermissions).([]string)
	return v
}

func RolesFromContext(c *fiber.Ctx) []string {
	v, _ := c.Locals(LocalsRoles).([]string)
	return v
}
