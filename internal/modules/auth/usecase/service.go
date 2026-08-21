package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"crm-starterkit/internal/modules/auth/domain"
	"crm-starterkit/pkg/crypto"
	"crm-starterkit/pkg/jwtutil"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")
	ErrInvalidRefresh     = errors.New("invalid refresh token")
	ErrRefreshRevoked     = errors.New("refresh token revoked")
)

type Service struct {
	users    domain.UserReader
	profile  domain.AccessProfileReader
	tokens   domain.RefreshTokenRepository
	jwt      *jwtutil.Manager
}

func NewService(
	users domain.UserReader,
	profile domain.AccessProfileReader,
	tokens domain.RefreshTokenRepository,
	jwtManager *jwtutil.Manager,
) *Service {
	return &Service{users: users, profile: profile, tokens: tokens, jwt: jwtManager}
}

func (s *Service) Login(ctx context.Context, in domain.LoginInput) (*domain.TokenPair, *domain.AuthUser, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" || in.Password == "" {
		return nil, nil, ErrInvalidCredentials
	}
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrInvalidCredentials
		}
		return nil, nil, err
	}
	if !user.IsActive {
		return nil, nil, ErrUserInactive
	}
	if !crypto.CheckPassword(user.PasswordHash, in.Password) {
		return nil, nil, ErrInvalidCredentials
	}
	pair, err := s.issueTokens(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	return pair, user, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	claims, err := s.jwt.Parse(refreshToken)
	if err != nil || claims.TokenType != "refresh" {
		return nil, ErrInvalidRefresh
	}
	stored, err := s.tokens.FindByJTI(ctx, claims.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidRefresh
		}
		return nil, err
	}
	if stored.RevokedAt != nil {
		return nil, ErrRefreshRevoked
	}
	if time.Now().UTC().After(stored.ExpiresAt) {
		return nil, ErrInvalidRefresh
	}
	user, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrInvalidRefresh
	}
	if !user.IsActive {
		return nil, ErrUserInactive
	}
	_ = s.tokens.RevokeByJTI(ctx, claims.ID, time.Now().UTC())
	return s.issueTokens(ctx, user)
}

func (s *Service) Logout(ctx context.Context, refreshToken string, userID string) error {
	if refreshToken != "" {
		claims, err := s.jwt.Parse(refreshToken)
		if err == nil && claims.TokenType == "refresh" {
			return s.tokens.RevokeByJTI(ctx, claims.ID, time.Now().UTC())
		}
	}
	if userID != "" {
		return s.tokens.RevokeAllByUser(ctx, userID, time.Now().UTC())
	}
	return ErrInvalidRefresh
}

func (s *Service) issueTokens(ctx context.Context, user *domain.AuthUser) (*domain.TokenPair, error) {
	roles, err := s.profile.RoleNamesByUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	perms, err := s.profile.PermissionKeysByUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	access, accessExp, err := s.jwt.GenerateAccessToken(user.ID, user.Email, roles, perms)
	if err != nil {
		return nil, err
	}
	refresh, refreshExp, err := s.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}
	claims, err := s.jwt.Parse(refresh)
	if err != nil {
		return nil, err
	}
	token := &domain.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenJTI:  claims.ID,
		ExpiresAt: refreshExp.UTC(),
		CreatedAt: time.Now().UTC(),
	}
	if err := s.tokens.Create(ctx, token); err != nil {
		return nil, err
	}
	return &domain.TokenPair{
		AccessToken:      access,
		RefreshToken:     refresh,
		AccessExpiresAt:  accessExp.UTC(),
		RefreshExpiresAt: refreshExp.UTC(),
	}, nil
}
