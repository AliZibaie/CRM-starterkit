package usecase

import (
	"context"
	"encoding/json"
	"time"

	"crm-starterkit/internal/core/http/middleware"
	"crm-starterkit/internal/modules/requestlog/domain"
	"crm-starterkit/internal/platform/geoip"
	"crm-starterkit/internal/platform/ratelimiter"

	"github.com/google/uuid"
)

type Service struct {
	repo    domain.Repository
	tracker *ratelimiter.Tracker
}

func NewService(repo domain.Repository, tracker *ratelimiter.Tracker) *Service {
	return &Service{repo: repo, tracker: tracker}
}

func (s *Service) RecordFromMiddleware(entry middleware.RequestLogEntry) {
	_ = s.Record(context.Background(), entry)
}

func (s *Service) Record(ctx context.Context, entry middleware.RequestLogEntry) error {
	headers, _ := json.Marshal(entry.Headers)
	geo := geoip.Lookup(entry.IP)
	suspicious := false
	if s.tracker != nil {
		_, suspicious = s.tracker.Hit(entry.IP)
		if entry.Status == 429 {
			s.tracker.MarkBlocked(entry.IP)
			suspicious = true
		}
	}
	return s.repo.Create(ctx, &domain.RequestEntry{
		ID: uuid.NewString(), RequestID: entry.RequestID, IP: entry.IP, Method: entry.Method, Path: entry.Path,
		HeadersJSON: string(headers), Payload: entry.Body, DurationMs: entry.Duration.Milliseconds(),
		Status: entry.Status, UserID: entry.UserID, UserAgent: entry.UserAgent,
		Country: geo.Country, City: geo.City, Suspicious: suspicious, CreatedAt: time.Now().UTC(),
	})
}

func (s *Service) List(ctx context.Context, f domain.ListFilter) ([]domain.RequestEntry, error) {
	return s.repo.List(ctx, f)
}

func (s *Service) Stats(ctx context.Context, f domain.ListFilter) (*domain.Stats, error) {
	return s.repo.Stats(ctx, f)
}

func (s *Service) IsSuspiciousIP(ip string) bool {
	if s.tracker == nil {
		return false
	}
	return s.tracker.IsSuspicious(ip)
}
