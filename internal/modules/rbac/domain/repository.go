package domain

import "context"

type Repository interface {
	CreateRole(ctx context.Context, role *Role) error
	UpdateRole(ctx context.Context, role *Role) error
	FindRoleByID(ctx context.Context, id string) (*Role, error)
	FindRoleByName(ctx context.Context, name string) (*Role, error)
	ListRoles(ctx context.Context) ([]Role, error)
	DeleteRole(ctx context.Context, id string) error

	CreatePermission(ctx context.Context, permission *Permission) error
	UpdatePermission(ctx context.Context, permission *Permission) error
	FindPermissionByID(ctx context.Context, id string) (*Permission, error)
	FindPermissionByKey(ctx context.Context, key string) (*Permission, error)
	ListPermissions(ctx context.Context) ([]Permission, error)
	UpsertPermissionByKey(ctx context.Context, permission *Permission) error

	AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error
	AssignRoleToUser(ctx context.Context, userID, roleID string) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID string) error
	PermissionIDsByRole(ctx context.Context, roleID string) ([]string, error)
	RoleIDsByUser(ctx context.Context, userID string) ([]string, error)
	PermissionKeysByUser(ctx context.Context, userID string) ([]string, error)
	RoleNamesByUser(ctx context.Context, userID string) ([]string, error)
	HasPermission(ctx context.Context, userID, permissionKey string) (bool, error)
}
