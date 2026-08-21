package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/wishlist/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Model struct {
	ID        string    `gorm:"primaryKey;size:36"`
	UserID    string    `gorm:"size:36;uniqueIndex:idx_wish_user_item;not null"`
	ItemID    string    `gorm:"size:36;uniqueIndex:idx_wish_user_item;not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (Model) TableName() string { return "wishlist_items" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error { return s.db.AutoMigrate(&Model{}) }

func (s *Store) Add(ctx context.Context, item *domain.WishlistItem) error {
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&Model{
		ID: item.ID, UserID: item.UserID, ItemID: item.ItemID, CreatedAt: item.CreatedAt,
	}).Error
}

func (s *Store) Remove(ctx context.Context, userID, itemID string) error {
	return s.db.WithContext(ctx).Where("user_id = ? AND item_id = ?", userID, itemID).Delete(&Model{}).Error
}

func (s *Store) List(ctx context.Context, userID string) ([]domain.WishlistItem, error) {
	var rows []Model
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.WishlistItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.WishlistItem{ID: r.ID, UserID: r.UserID, ItemID: r.ItemID, CreatedAt: r.CreatedAt})
	}
	return out, nil
}
