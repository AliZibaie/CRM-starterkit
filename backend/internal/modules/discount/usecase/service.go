package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"crm-starterkit/internal/modules/discount/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("discount not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrNotApplicable = errors.New("discount not applicable")
)

type Service struct{ repo domain.Repository }

func NewService(repo domain.Repository) *Service { return &Service{repo: repo} }

func (s *Service) Create(ctx context.Context, code, typ string, value, minAmount, maxUses int64, starts, ends *time.Time) (*domain.Discount, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" || value <= 0 {
		return nil, ErrInvalidInput
	}
	if typ != domain.TypePercent && typ != domain.TypeFixed {
		return nil, ErrInvalidInput
	}
	if typ == domain.TypePercent && value > 100 {
		return nil, ErrInvalidInput
	}
	now := time.Now().UTC()
	d := &domain.Discount{
		ID: uuid.NewString(), Code: code, Type: typ, Value: value, MinOrderAmount: minAmount,
		MaxUses: maxUses, StartsAt: starts, EndsAt: ends, IsActive: true, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Discount, error) { return s.repo.List(ctx) }

func (s *Service) Preview(ctx context.Context, code string, orderAmount int64) (*domain.ApplyResult, error) {
	d, err := s.repo.FindByCode(ctx, strings.ToUpper(strings.TrimSpace(code)))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	amount, err := s.calculate(d, orderAmount)
	if err != nil {
		return nil, err
	}
	return &domain.ApplyResult{Code: d.Code, DiscountAmount: amount, FinalAmount: orderAmount - amount}, nil
}

func (s *Service) ApplyAndConsume(ctx context.Context, code string, orderAmount int64) (*domain.ApplyResult, error) {
	d, err := s.repo.FindByCode(ctx, strings.ToUpper(strings.TrimSpace(code)))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	amount, err := s.calculate(d, orderAmount)
	if err != nil {
		return nil, err
	}
	if err := s.repo.IncrementUsage(ctx, d.ID); err != nil {
		return nil, err
	}
	return &domain.ApplyResult{Code: d.Code, DiscountAmount: amount, FinalAmount: orderAmount - amount}, nil
}

func (s *Service) calculate(d *domain.Discount, orderAmount int64) (int64, error) {
	now := time.Now().UTC()
	if !d.IsActive {
		return 0, ErrNotApplicable
	}
	if d.StartsAt != nil && now.Before(*d.StartsAt) {
		return 0, ErrNotApplicable
	}
	if d.EndsAt != nil && now.After(*d.EndsAt) {
		return 0, ErrNotApplicable
	}
	if d.MaxUses > 0 && d.UsedCount >= d.MaxUses {
		return 0, ErrNotApplicable
	}
	if orderAmount < d.MinOrderAmount {
		return 0, ErrNotApplicable
	}
	var amount int64
	if d.Type == domain.TypePercent {
		amount = orderAmount * d.Value / 100
	} else {
		amount = d.Value
	}
	if amount > orderAmount {
		amount = orderAmount
	}
	return amount, nil
}
