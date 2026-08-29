package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/payment/domain"

	"gorm.io/gorm"
)

type PaymentModel struct {
	ID          string    `gorm:"primaryKey;size:36"`
	OrderID     string    `gorm:"size:36;index;not null"`
	Gateway     string    `gorm:"size:50;not null"`
	Amount      int64     `gorm:"not null"`
	Currency    string    `gorm:"size:10;not null"`
	Status      string    `gorm:"size:50;index;not null"`
	Authority   string    `gorm:"size:255;uniqueIndex"`
	ProviderRef string    `gorm:"size:255"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (PaymentModel) TableName() string { return "payments" }

type OrderAmountRow struct {
	TotalAmount int64  `gorm:"column:total_amount"`
	Currency    string `gorm:"column:currency"`
	Status      string `gorm:"column:status"`
}

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error {
	return s.db.AutoMigrate(&PaymentModel{})
}

func (s *Store) Create(ctx context.Context, p *domain.Payment) error {
	return s.db.WithContext(ctx).Create(toModel(p)).Error
}

func (s *Store) Update(ctx context.Context, p *domain.Payment) error {
	return s.db.WithContext(ctx).Save(toModel(p)).Error
}

func (s *Store) FindByID(ctx context.Context, id string) (*domain.Payment, error) {
	var m PaymentModel
	if err := s.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (s *Store) FindByAuthority(ctx context.Context, authority string) (*domain.Payment, error) {
	var m PaymentModel
	if err := s.db.WithContext(ctx).First(&m, "authority = ?", authority).Error; err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (s *Store) FindByOrder(ctx context.Context, orderID string) (*domain.Payment, error) {
	var m PaymentModel
	if err := s.db.WithContext(ctx).Where("order_id = ?", orderID).Order("created_at desc").First(&m).Error; err != nil {
		return nil, err
	}
	return toDomain(&m), nil
}

func (s *Store) FindOrderAmount(ctx context.Context, orderID string) (int64, string, string, error) {
	var row OrderAmountRow
	err := s.db.WithContext(ctx).Table("orders").
		Select("total_amount, currency, status").
		Where("id = ?", orderID).
		Take(&row).Error
	if err != nil {
		return 0, "", "", err
	}
	return row.TotalAmount, row.Currency, row.Status, nil
}

func (s *Store) ListByStatus(ctx context.Context, status string, limit int) ([]domain.Payment, error) {
	if limit <= 0 {
		limit = 50
	}
	var rows []PaymentModel
	q := s.db.WithContext(ctx).Where("status = ?", status).Order("created_at desc").Limit(limit)
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Payment, 0, len(rows))
	for i := range rows {
		out = append(out, *toDomain(&rows[i]))
	}
	return out, nil
}

func toModel(p *domain.Payment) *PaymentModel {
	return &PaymentModel{
		ID: p.ID, OrderID: p.OrderID, Gateway: p.Gateway, Amount: p.Amount, Currency: p.Currency,
		Status: p.Status, Authority: p.Authority, ProviderRef: p.ProviderRef,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func toDomain(m *PaymentModel) *domain.Payment {
	return &domain.Payment{
		ID: m.ID, OrderID: m.OrderID, Gateway: m.Gateway, Amount: m.Amount, Currency: m.Currency,
		Status: m.Status, Authority: m.Authority, ProviderRef: m.ProviderRef,
		CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}
