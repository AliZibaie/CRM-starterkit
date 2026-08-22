package middleware

import (
	"crm-starterkit/internal/core/http/response"

	"github.com/gofiber/fiber/v2"
)

type PermissionChecker interface {
	HasPermission(userID string, permission string) (bool, error)
}

type PermissionFunc func(userID string, permission string) (bool, error)

func (f PermissionFunc) HasPermission(userID string, permission string) (bool, error) {
	return f(userID, permission)
}

func RBAC(checker PermissionChecker, permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := UserIDFromContext(c)
		if userID == "" {
			return response.Unauthorized(c, "authentication required")
		}

		for _, p := range PermissionsFromContext(c) {
			if p == permission || p == "*" {
				return c.Next()
			}
		}

		if checker != nil {
			ok, err := checker.HasPermission(userID, permission)
			if err != nil {
				return response.Internal(c, "permission check failed")
			}
			if ok {
				return c.Next()
			}
		}

		return response.Forbidden(c, "insufficient permissions")
	}
}

func RequireAnyPermission(permissions ...string) fiber.Handler {
	set := make(map[string]struct{}, len(permissions))
	for _, p := range permissions {
		set[p] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		userID := UserIDFromContext(c)
		if userID == "" {
			return response.Unauthorized(c, "authentication required")
		}
		for _, p := range PermissionsFromContext(c) {
			if _, ok := set[p]; ok || p == "*" {
				return c.Next()
			}
		}
		return response.Forbidden(c, "insufficient permissions")
	}
}

func RequireRoles(roles ...string) fiber.Handler {
	set := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		set[r] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		userID := UserIDFromContext(c)
		if userID == "" {
			return response.Unauthorized(c, "authentication required")
		}
		for _, r := range RolesFromContext(c) {
			if _, ok := set[r]; ok {
				return c.Next()
			}
		}
		return response.Forbidden(c, "insufficient role")
	}
}
