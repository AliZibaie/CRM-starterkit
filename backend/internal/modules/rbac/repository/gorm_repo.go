package repository

import (
	"context"
	"errors"
	"time"

	"crm-starterkit/internal/modules/rbac/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleModel struct {
	ID          string    `gorm:"primaryKey;size:36"`
	Name        string    `gorm:"uniqueIndex;size:100;not null"`
	Description string    `gorm:"size:255"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (RoleModel) TableName() string { return "roles" }

type PermissionModel struct {
	ID          string    `gorm:"primaryKey;size:36"`
	Key         string    `gorm:"uniqueIndex;size:150;not null"`
	Name        string    `gorm:"size:150;not null"`
	Description string    `gorm:"size:255"`
	Module      string    `gorm:"size:100;index"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (PermissionModel) TableName() string { return "permissions" }

type RolePermissionModel struct {
	RoleID       string `gorm:"primaryKey;size:36"`
	PermissionID string `gorm:"primaryKey;size:36;index"`
}

func (RolePermissionModel) TableName() string { return "role_permissions" }

type UserRoleModel struct {
	UserID string `gorm:"primaryKey;size:36"`
	RoleID string `gorm:"primaryKey;size:36;index"`
}

func (UserRoleModel) TableName() string { return "user_roles" }

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) AutoMigrate() error {
	return s.db.AutoMigrate(&RoleModel{}, &PermissionModel{}, &RolePermissionModel{}, &UserRoleModel{})
}

func (s *Store) CreateRole(ctx context.Context, role *domain.Role) error {
	return s.db.WithContext(ctx).Create(&RoleModel{
		ID: role.ID, Name: role.Name, Description: role.Description,
		CreatedAt: role.CreatedAt, UpdatedAt: role.UpdatedAt,
	}).Error
}

func (s *Store) UpdateRole(ctx context.Context, role *domain.Role) error {
	return s.db.WithContext(ctx).Save(&RoleModel{
		ID: role.ID, Name: role.Name, Description: role.Description,
		CreatedAt: role.CreatedAt, UpdatedAt: role.UpdatedAt,
	}).Error
}

func (s *Store) FindRoleByID(ctx context.Context, id string) (*domain.Role, error) {
	var m RoleModel
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &domain.Role{ID: m.ID, Name: m.Name, Description: m.Description, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (s *Store) FindRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	var m RoleModel
	if err := s.db.WithContext(ctx).First(&m, "name = ?", name).Error; err != nil {
		return nil, err
	}
	return &domain.Role{ID: m.ID, Name: m.Name, Description: m.Description, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (s *Store) ListRoles(ctx context.Context) ([]domain.Role, error) {
	var rows []RoleModel
	if err := s.db.WithContext(ctx).Order("name asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Role, 0, len(rows))
	for _, m := range rows {
		out = append(out, domain.Role{ID: m.ID, Name: m.Name, Description: m.Description, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt})
	}
	return out, nil
}

func (s *Store) DeleteRole(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&RolePermissionModel{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id).Delete(&UserRoleModel{}).Error; err != nil {
			return err
		}
		return tx.Delete(&RoleModel{}, "id = ?", id).Error
	})
}

func (s *Store) CreatePermission(ctx context.Context, permission *domain.Permission) error {
	return s.db.WithContext(ctx).Create(&PermissionModel{
		ID: permission.ID, Key: permission.Key, Name: permission.Name,
		Description: permission.Description, Module: permission.Module,
		CreatedAt: permission.CreatedAt, UpdatedAt: permission.UpdatedAt,
	}).Error
}

func (s *Store) UpdatePermission(ctx context.Context, permission *domain.Permission) error {
	return s.db.WithContext(ctx).Save(&PermissionModel{
		ID: permission.ID, Key: permission.Key, Name: permission.Name,
		Description: permission.Description, Module: permission.Module,
		CreatedAt: permission.CreatedAt, UpdatedAt: permission.UpdatedAt,
	}).Error
}

func (s *Store) FindPermissionByID(ctx context.Context, id string) (*domain.Permission, error) {
	var m PermissionModel
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toPermission(&m), nil
}

func (s *Store) FindPermissionByKey(ctx context.Context, key string) (*domain.Permission, error) {
	var m PermissionModel
	if err := s.db.WithContext(ctx).First(&m, "key = ?", key).Error; err != nil {
		return nil, err
	}
	return toPermission(&m), nil
}

func (s *Store) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	var rows []PermissionModel
	if err := s.db.WithContext(ctx).Order("module asc, key asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Permission, 0, len(rows))
	for i := range rows {
		out = append(out, *toPermission(&rows[i]))
	}
	return out, nil
}

func (s *Store) UpsertPermissionByKey(ctx context.Context, permission *domain.Permission) error {
	existing, err := s.FindPermissionByKey(ctx, permission.Key)
	if err == nil && existing != nil {
		permission.ID = existing.ID
		permission.CreatedAt = existing.CreatedAt
		permission.UpdatedAt = time.Now().UTC()
		return s.UpdatePermission(ctx, permission)
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return s.CreatePermission(ctx, permission)
}

func (s *Store) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&RolePermissionModel{
		RoleID: roleID, PermissionID: permissionID,
	}).Error
}

func (s *Store) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	return s.db.WithContext(ctx).Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&RolePermissionModel{}).Error
}

func (s *Store) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&UserRoleModel{
		UserID: userID, RoleID: roleID,
	}).Error
}

func (s *Store) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	return s.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&UserRoleModel{}).Error
}

func (s *Store) PermissionIDsByRole(ctx context.Context, roleID string) ([]string, error) {
	var rows []RolePermissionModel
	if err := s.db.WithContext(ctx).Where("role_id = ?", roleID).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.PermissionID)
	}
	return ids, nil
}

func (s *Store) RoleIDsByUser(ctx context.Context, userID string) ([]string, error) {
	var rows []UserRoleModel
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.RoleID)
	}
	return ids, nil
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

func (s *Store) HasPermission(ctx context.Context, userID, permissionKey string) (bool, error) {
	var count int64
	err := s.db.WithContext(ctx).Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.key = ?", userID, permissionKey).
		Count(&count).Error
	return count > 0, err
}

func toPermission(m *PermissionModel) *domain.Permission {
	return &domain.Permission{
		ID: m.ID, Key: m.Key, Name: m.Name, Description: m.Description,
		Module: m.Module, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}
