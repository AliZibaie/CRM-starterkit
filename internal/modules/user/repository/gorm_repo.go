package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/user/domain"

	"gorm.io/gorm"
)

type UserModel struct {
	ID           string    `gorm:"primaryKey;size:36"`
	Email        string    `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	FirstName    string    `gorm:"size:100"`
	LastName     string    `gorm:"size:100"`
	IsActive     bool      `gorm:"not null;default:true"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

func (UserModel) TableName() string { return "users" }

type UserRoleModel struct {
	UserID string `gorm:"primaryKey;size:36"`
	RoleID string `gorm:"primaryKey;size:36;index"`
}

func (UserRoleModel) TableName() string { return "user_roles" }

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) AutoMigrate() error {
	return r.db.AutoMigrate(&UserModel{}, &UserRoleModel{})
}

func (r *GormRepository) Create(ctx context.Context, user *domain.User) error {
	m := toModel(user)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *GormRepository) Update(ctx context.Context, user *domain.User) error {
	m := toModel(user)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *GormRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", id).Delete(&UserRoleModel{}).Error; err != nil {
			return err
		}
		return tx.Delete(&UserModel{}, "id = ?", id).Error
	})
}

func (r *GormRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var m UserModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *GormRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m UserModel
	if err := r.db.WithContext(ctx).First(&m, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *GormRepository) List(ctx context.Context, offset, limit int) ([]domain.User, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&UserModel{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []UserModel
	if err := q.Order("created_at desc").Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.User, 0, len(rows))
	for i := range rows {
		out = append(out, *toDomain(&rows[i]))
	}
	return out, total, nil
}

func (r *GormRepository) ReplaceRoles(ctx context.Context, userID string, roleIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&UserRoleModel{}).Error; err != nil {
			return err
		}
		for _, roleID := range roleIDs {
			if err := tx.Create(&UserRoleModel{UserID: userID, RoleID: roleID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *GormRepository) Count(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&UserModel{}).Count(&total).Error
	return total, err
}

func (r *GormRepository) RoleIDsByUser(ctx context.Context, userID string) ([]string, error) {
	var rows []UserRoleModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.RoleID)
	}
	return ids, nil
}

func toModel(u *domain.User) *UserModel {
	return &UserModel{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		IsActive:     u.IsActive,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func toDomain(m *UserModel) *domain.User {
	return &domain.User{
		ID:           m.ID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		FirstName:    m.FirstName,
		LastName:     m.LastName,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
