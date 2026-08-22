package domain

import "time"

type AuthUser struct {
	ID           string
	Email        string
	PasswordHash string
	IsActive     bool
}

type RefreshToken struct {
	ID        string
	UserID    string
	TokenJTI  string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

type LoginInput struct {
	Email    string
	Password string
}
