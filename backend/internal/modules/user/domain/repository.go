package domain

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, offset, limit int) ([]User, int64, error)
	Count(ctx context.Context) (int64, error)
	ReplaceRoles(ctx context.Context, userID string, roleIDs []string) error
	RoleIDsByUser(ctx context.Context, userID string) ([]string, error)
}
