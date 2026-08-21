package domain

import (
	"context"
	"time"
)

type UserReader interface {
	FindByEmail(ctx context.Context, email string) (*AuthUser, error)
	FindByID(ctx context.Context, id string) (*AuthUser, error)
}

type AccessProfileReader interface {
	RoleNamesByUser(ctx context.Context, userID string) ([]string, error)
	PermissionKeysByUser(ctx context.Context, userID string) ([]string, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	FindByJTI(ctx context.Context, jti string) (*RefreshToken, error)
	RevokeByJTI(ctx context.Context, jti string, revokedAt time.Time) error
	RevokeAllByUser(ctx context.Context, userID string, revokedAt time.Time) error
}
