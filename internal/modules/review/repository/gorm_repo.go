package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/review/domain"

	"gorm.io/gorm"
)

type Model struct {
	ID        string    `gorm:"primaryKey;size:36"`
	ItemID    string    `gorm:"size:36;index;not null"`
	UserID    string    `gorm:"size:36;index;not null"`
	Rating    int       `gorm:"not null"`
	Comment   string    `gorm:"type:text"`
	IsVisible bool      `gorm:"not null;default:true"`
	CreatedAt time.Time `gorm:"not null"`
}

func (Model) TableName() string { return "product_reviews" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error { return s.db.AutoMigrate(&Model{}) }

func (s *Store) Create(ctx context.Context, r *domain.Review) error {
	return s.db.WithContext(ctx).Create(&Model{
		ID: r.ID, ItemID: r.ItemID, UserID: r.UserID, Rating: r.Rating, Comment: r.Comment, IsVisible: r.IsVisible, CreatedAt: r.CreatedAt,
	}).Error
}

func (s *Store) ListByItem(ctx context.Context, itemID string) ([]domain.Review, error) {
	var rows []Model
	if err := s.db.WithContext(ctx).Where("item_id = ? AND is_visible = ?", itemID, true).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Review, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.Review{ID: r.ID, ItemID: r.ItemID, UserID: r.UserID, Rating: r.Rating, Comment: r.Comment, IsVisible: r.IsVisible, CreatedAt: r.CreatedAt})
	}
	return out, nil
}

func (s *Store) AverageRating(ctx context.Context, itemID string) (float64, error) {
	var avg float64
	err := s.db.WithContext(ctx).Model(&Model{}).Select("COALESCE(AVG(rating),0)").Where("item_id = ? AND is_visible = ?", itemID, true).Scan(&avg).Error
	return avg, err
}

func (s *Store) Delete(ctx context.Context, id, userID string) error {
	return s.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&Model{}).Error
}
