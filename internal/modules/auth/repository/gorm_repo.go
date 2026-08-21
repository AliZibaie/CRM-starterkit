package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/auth/domain"

	"gorm.io/gorm"
)

type userRow struct {
	ID           string `gorm:"column:id"`
	Email        string `gorm:"column:email"`
	PasswordHash string `gorm:"column:password_hash"`
	IsActive     bool   `gorm:"column:is_active"`
}

func (userRow) TableName() string { return "users" }

type RefreshTokenModel struct {
	ID        string     `gorm:"primaryKey;size:36"`
	UserID    string     `gorm:"size:36;index;not null"`
	TokenJTI  string     `gorm:"uniqueIndex;size:36;not null"`
	ExpiresAt time.Time  `gorm:"not null"`
	RevokedAt *time.Time
	CreatedAt time.Time `gorm:"not null"`
}

func (RefreshTokenModel) TableName() string { return "refresh_tokens" }

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) AutoMigrate() error {
	return s.db.AutoMigrate(&RefreshTokenModel{})
}

func (s *Store) FindByEmail(ctx context.Context, email string) (*domain.AuthUser, error) {
	var row userRow
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&row).Error; err != nil {
		return nil, err
	}
	return &domain.AuthUser{
		ID: row.ID, Email: row.Email, PasswordHash: row.PasswordHash, IsActive: row.IsActive,
	}, nil
}

func (s *Store) FindByID(ctx context.Context, id string) (*domain.AuthUser, error) {
	var row userRow
	if err := s.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		return nil, err
	}
	return &domain.AuthUser{
		ID: row.ID, Email: row.Email, PasswordHash: row.PasswordHash, IsActive: row.IsActive,
	}, nil
}

func (s *Store) RoleNamesByUser(ctx context.Context, userID string) ([]string, error) {
	var names []string
	err := s.db.WithContext(ctx).Table("roles").
		Select("roles.name").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Pluck("roles.name", &names).Error
	return names, err
}

func (s *Store) PermissionKeysByUser(ctx context.Context, userID string) ([]string, error) {
	var keys []string
	err := s.db.WithContext(ctx).Table("permissions").
		Select("DISTINCT permissions.key").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("permissions.key", &keys).Error
	return keys, err
}

func (s *Store) Create(ctx context.Context, token *domain.RefreshToken) error {
	return s.db.WithContext(ctx).Create(&RefreshTokenModel{
		ID: token.ID, UserID: token.UserID, TokenJTI: token.TokenJTI,
		ExpiresAt: token.ExpiresAt, RevokedAt: token.RevokedAt, CreatedAt: token.CreatedAt,
	}).Error
}

func (s *Store) FindByJTI(ctx context.Context, jti string) (*domain.RefreshToken, error) {
	var m RefreshTokenModel
	if err := s.db.WithContext(ctx).First(&m, "token_jti = ?", jti).Error; err != nil {
		return nil, err
	}
	return &domain.RefreshToken{
		ID: m.ID, UserID: m.UserID, TokenJTI: m.TokenJTI,
		ExpiresAt: m.ExpiresAt, RevokedAt: m.RevokedAt, CreatedAt: m.CreatedAt,
	}, nil
}

func (s *Store) RevokeByJTI(ctx context.Context, jti string, revokedAt time.Time) error {
	return s.db.WithContext(ctx).Model(&RefreshTokenModel{}).
		Where("token_jti = ? AND revoked_at IS NULL", jti).
		Update("revoked_at", revokedAt).Error
}

func (s *Store) RevokeAllByUser(ctx context.Context, userID string, revokedAt time.Time) error {
	return s.db.WithContext(ctx).Model(&RefreshTokenModel{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", revokedAt).Error
}
