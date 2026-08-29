package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"crm-starterkit/internal/modules/user/domain"
	"crm-starterkit/pkg/crypto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyUsed  = errors.New("email already used")
	ErrInvalidUserInput  = errors.New("invalid user input")
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, in domain.CreateUserInput) (*domain.User, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" || len(in.Password) < 6 {
		return nil, ErrInvalidUserInput
	}
	existing, err := s.repo.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, ErrEmailAlreadyUsed
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	hash, err := crypto.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	user := &domain.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: hash,
		FirstName:    strings.TrimSpace(in.FirstName),
		LastName:     strings.TrimSpace(in.LastName),
		IsActive:     in.IsActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	if len(in.RoleIDs) > 0 {
		if err := s.repo.ReplaceRoles(ctx, user.ID, in.RoleIDs); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (s *Service) Get(ctx context.Context, id string) (*domain.User, []string, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrUserNotFound
		}
		return nil, nil, err
	}
	roleIDs, err := s.repo.RoleIDsByUser(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return user, roleIDs, nil
}

func (s *Service) List(ctx context.Context, page, perPage int) ([]domain.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	offset := (page - 1) * perPage
	return s.repo.List(ctx, offset, perPage)
}

func (s *Service) Update(ctx context.Context, id string, in domain.UpdateUserInput) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if in.Email != nil {
		email := strings.TrimSpace(strings.ToLower(*in.Email))
		if email == "" {
			return nil, ErrInvalidUserInput
		}
		other, err := s.repo.FindByEmail(ctx, email)
		if err == nil && other != nil && other.ID != id {
			return nil, ErrEmailAlreadyUsed
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		user.Email = email
	}
	if in.Password != nil {
		if len(*in.Password) < 6 {
			return nil, ErrInvalidUserInput
		}
		hash, err := crypto.HashPassword(*in.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hash
	}
	if in.FirstName != nil {
		user.FirstName = strings.TrimSpace(*in.FirstName)
	}
	if in.LastName != nil {
		user.LastName = strings.TrimSpace(*in.LastName)
	}
	if in.IsActive != nil {
		user.IsActive = *in.IsActive
	}
	user.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) RoleIDsByUser(ctx context.Context, userID string) ([]string, error) {
	return s.repo.RoleIDsByUser(ctx, userID)
}

func (s *Service) AssignRoles(ctx context.Context, userID string, roleIDs []string) error {
	_, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return s.repo.ReplaceRoles(ctx, userID, roleIDs)
}

func (s *Service) EnsureDefaultAdmin(ctx context.Context, email, password string, roleID string) (*domain.User, error) {
	_, total, err := s.repo.List(ctx, 0, 1)
	if err != nil {
		return nil, err
	}
	if total > 0 {
		return nil, nil
	}
	return s.Create(ctx, domain.CreateUserInput{
		Email:     email,
		Password:  password,
		FirstName: "Admin",
		LastName:  "User",
		IsActive:  true,
		RoleIDs:   []string{roleID},
	})
}

func (s *Service) Count(ctx context.Context) (int64, error) {
	_, total, err := s.repo.List(ctx, 0, 1)
	return total, err
}
