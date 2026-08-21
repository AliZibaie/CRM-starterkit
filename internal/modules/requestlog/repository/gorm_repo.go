package repository

import (
	"context"
	"strconv"
	"time"

	"crm-starterkit/internal/modules/requestlog/domain"

	"gorm.io/gorm"
)

type Model struct {
	ID          string    `gorm:"primaryKey;size:36"`
	RequestID   string    `gorm:"size:36;index"`
	IP          string    `gorm:"size:64;index"`
	Method      string    `gorm:"size:16"`
	Path        string    `gorm:"size:500;index"`
	HeadersJSON string    `gorm:"type:text"`
	Payload     string    `gorm:"type:text"`
	DurationMs  int64     `gorm:"not null"`
	Status      int       `gorm:"index"`
	UserID      string    `gorm:"size:36;index"`
	UserAgent   string    `gorm:"size:500"`
	Country     string    `gorm:"size:100;index"`
	City        string    `gorm:"size:100"`
	Suspicious  bool      `gorm:"index"`
	CreatedAt   time.Time `gorm:"index;not null"`
}

func (Model) TableName() string { return "request_logs" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error { return s.db.AutoMigrate(&Model{}) }

func (s *Store) Create(ctx context.Context, e *domain.RequestEntry) error {
	return s.db.WithContext(ctx).Create(&Model{
		ID: e.ID, RequestID: e.RequestID, IP: e.IP, Method: e.Method, Path: e.Path,
		HeadersJSON: e.HeadersJSON, Payload: e.Payload, DurationMs: e.DurationMs, Status: e.Status,
		UserID: e.UserID, UserAgent: e.UserAgent, Country: e.Country, City: e.City,
		Suspicious: e.Suspicious, CreatedAt: e.CreatedAt,
	}).Error
}

func (s *Store) applyFilter(q *gorm.DB, f domain.ListFilter) *gorm.DB {
	if f.IP != "" {
		q = q.Where("ip = ?", f.IP)
	}
	if f.Path != "" {
		q = q.Where("path LIKE ?", "%"+f.Path+"%")
	}
	if f.UserID != "" {
		q = q.Where("user_id = ?", f.UserID)
	}
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	if f.From != nil {
		q = q.Where("created_at >= ?", *f.From)
	}
	if f.To != nil {
		q = q.Where("created_at <= ?", *f.To)
	}
	return q
}

func (s *Store) List(ctx context.Context, f domain.ListFilter) ([]domain.RequestEntry, error) {
	q := s.applyFilter(s.db.WithContext(ctx).Model(&Model{}), f)
	limit := f.Limit
	if limit <= 0 {
		limit = 100
	}
	var rows []Model
	if err := q.Order("created_at desc").Limit(limit).Offset(f.Offset).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.RequestEntry, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.RequestEntry{
			ID: r.ID, RequestID: r.RequestID, IP: r.IP, Method: r.Method, Path: r.Path,
			HeadersJSON: r.HeadersJSON, Payload: r.Payload, DurationMs: r.DurationMs, Status: r.Status,
			UserID: r.UserID, UserAgent: r.UserAgent, Country: r.Country, City: r.City,
			Suspicious: r.Suspicious, CreatedAt: r.CreatedAt,
		})
	}
	return out, nil
}

func (s *Store) Stats(ctx context.Context, f domain.ListFilter) (*domain.Stats, error) {
	q := s.applyFilter(s.db.WithContext(ctx).Model(&Model{}), f)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}
	var suspicious int64
	if err := s.applyFilter(s.db.WithContext(ctx).Model(&Model{}), f).Where("suspicious = ?", true).Count(&suspicious).Error; err != nil {
		return nil, err
	}
	type rowCount struct {
		Key   string
		Count int64
	}
	var statusRows []struct {
		Status int
		Count  int64
	}
	_ = s.applyFilter(s.db.WithContext(ctx).Model(&Model{}), f).Select("status, count(*) as count").Group("status").Scan(&statusRows)
	byStatus := map[string]int64{}
	for _, r := range statusRows {
		byStatus[strconv.Itoa(r.Status)] = r.Count
	}
	var countryRows []rowCount
	_ = s.applyFilter(s.db.WithContext(ctx).Model(&Model{}), f).Select("country as key, count(*) as count").Group("country").Scan(&countryRows)
	byCountry := map[string]int64{}
	for _, r := range countryRows {
		byCountry[r.Key] = r.Count
	}
	return &domain.Stats{Total: total, Suspicious: suspicious, ByStatus: byStatus, ByCountry: byCountry}, nil
}
