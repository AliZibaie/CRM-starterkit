package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"crm-starterkit/internal/core/module"
	"crm-starterkit/internal/modules/rbac/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrRoleNotFound       = errors.New("role not found")
	ErrPermissionNotFound = errors.New("permission not found")
	ErrRoleExists         = errors.New("role already exists")
	ErrPermissionExists   = errors.New("permission already exists")
	ErrInvalidInput       = errors.New("invalid input")
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateRole(ctx context.Context, in domain.CreateRoleInput) (*domain.Role, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, ErrInvalidInput
	}
	existing, err := s.repo.FindRoleByName(ctx, name)
	if err == nil && existing != nil {
		return nil, ErrRoleExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	now := time.Now().UTC()
	role := &domain.Role{
		ID:          uuid.NewString(),
		Name:        name,
		Description: strings.TrimSpace(in.Description),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.CreateRole(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *Service) ListRoles(ctx context.Context) ([]domain.Role, error) {
	return s.repo.ListRoles(ctx)
}

func (s *Service) GetRole(ctx context.Context, id string) (*domain.Role, []string, error) {
	role, err := s.repo.FindRoleByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrRoleNotFound
		}
		return nil, nil, err
	}
	permIDs, err := s.repo.PermissionIDsByRole(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return role, permIDs, nil
}

func (s *Service) CreatePermission(ctx context.Context, in domain.CreatePermissionInput) (*domain.Permission, error) {
	key := strings.TrimSpace(in.Key)
	if key == "" || strings.TrimSpace(in.Name) == "" {
		return nil, ErrInvalidInput
	}
	existing, err := s.repo.FindPermissionByKey(ctx, key)
	if err == nil && existing != nil {
		return nil, ErrPermissionExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	now := time.Now().UTC()
	perm := &domain.Permission{
		ID:          uuid.NewString(),
		Key:         key,
		Name:        strings.TrimSpace(in.Name),
		Description: strings.TrimSpace(in.Description),
		Module:      strings.TrimSpace(in.Module),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.CreatePermission(ctx, perm); err != nil {
		return nil, err
	}
	return perm, nil
}

func (s *Service) EditPermission(ctx context.Context, id string, in domain.EditPermissionInput) (*domain.Permission, error) {
	perm, err := s.repo.FindPermissionByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPermissionNotFound
		}
		return nil, err
	}
	if in.Name != nil {
		perm.Name = strings.TrimSpace(*in.Name)
	}
	if in.Description != nil {
		perm.Description = strings.TrimSpace(*in.Description)
	}
	if in.Module != nil {
		perm.Module = strings.TrimSpace(*in.Module)
	}
	perm.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdatePermission(ctx, perm); err != nil {
		return nil, err
	}
	return perm, nil
}

func (s *Service) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	return s.repo.ListPermissions(ctx)
}

func (s *Service) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	if _, err := s.repo.FindRoleByID(ctx, roleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		return err
	}
	if _, err := s.repo.FindPermissionByID(ctx, permissionID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPermissionNotFound
		}
		return err
	}
	return s.repo.AssignPermissionToRole(ctx, roleID, permissionID)
}

func (s *Service) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	if userID == "" || roleID == "" {
		return ErrInvalidInput
	}
	if _, err := s.repo.FindRoleByID(ctx, roleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
		return err
	}
	return s.repo.AssignRoleToUser(ctx, userID, roleID)
}

func (s *Service) SyncPermissionsFromRegistry(ctx context.Context, perms []module.Permission) (int, error) {
	count := 0
	now := time.Now().UTC()
	for _, p := range perms {
		key := strings.TrimSpace(p.Key)
		if key == "" {
			continue
		}
		item := &domain.Permission{
			ID:          uuid.NewString(),
			Key:         key,
			Name:        strings.TrimSpace(p.Name),
			Description: strings.TrimSpace(p.Description),
			Module:      strings.TrimSpace(p.Module),
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if item.Name == "" {
			item.Name = key
		}
		if err := s.repo.UpsertPermissionByKey(ctx, item); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *Service) EnsureAdminRole(ctx context.Context) (*domain.Role, error) {
	role, err := s.repo.FindRoleByName(ctx, "admin")
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		role, err = s.CreateRole(ctx, domain.CreateRoleInput{
			Name:        "admin",
			Description: "Full access administrator",
		})
		if err != nil {
			return nil, err
		}
	}
	if err := s.GrantAllPermissionsToRole(ctx, role.ID); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *Service) GrantAllPermissionsToRole(ctx context.Context, roleID string) error {
	perms, err := s.repo.ListPermissions(ctx)
	if err != nil {
		return err
	}
	for _, p := range perms {
		if err := s.repo.AssignPermissionToRole(ctx, roleID, p.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) HasPermission(ctx context.Context, userID, permission string) (bool, error) {
	return s.repo.HasPermission(ctx, userID, permission)
}

func (s *Service) PermissionKeysByUser(ctx context.Context, userID string) ([]string, error) {
	return s.repo.PermissionKeysByUser(ctx, userID)
}

func (s *Service) RoleNamesByUser(ctx context.Context, userID string) ([]string, error) {
	return s.repo.RoleNamesByUser(ctx, userID)
}
