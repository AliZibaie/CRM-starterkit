package usecase

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"crm-starterkit/internal/modules/cms/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidInput = errors.New("invalid input")
	slugRe          = regexp.MustCompile(`[^a-z0-9]+`)
)

type Service struct{ repo domain.Repository }

func NewService(repo domain.Repository) *Service { return &Service{repo: repo} }

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return strings.Trim(slugRe.ReplaceAllString(s, "-"), "-")
}

func (s *Service) CreatePage(ctx context.Context, title, slug string) (*domain.Page, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrInvalidInput
	}
	if slug == "" {
		slug = slugify(title)
	}
	now := time.Now().UTC()
	p := &domain.Page{ID: uuid.NewString(), Title: title, Slug: slug, Status: "draft", CreatedAt: now, UpdatedAt: now}
	return p, s.repo.CreatePage(ctx, p)
}

func (s *Service) ListPages(ctx context.Context) ([]domain.Page, error) { return s.repo.ListPages(ctx) }

func (s *Service) GetPage(ctx context.Context, slug string) (*domain.Page, []domain.Block, error) {
	p, err := s.repo.FindPageBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrNotFound
		}
		return nil, nil, err
	}
	blocks, err := s.repo.ListBlocks(ctx, p.ID)
	return p, blocks, err
}

func (s *Service) SetBlocks(ctx context.Context, pageID string, blocks []domain.Block) error {
	for i := range blocks {
		if blocks[i].ID == "" {
			blocks[i].ID = uuid.NewString()
		}
		blocks[i].PageID = pageID
		blocks[i].SortOrder = i
	}
	return s.repo.SaveBlocks(ctx, pageID, blocks)
}

func (s *Service) GetPageByID(ctx context.Context, id string) (*domain.Page, []domain.Block, error) {
	p, err := s.repo.FindPageByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrNotFound
		}
		return nil, nil, err
	}
	blocks, err := s.repo.ListBlocks(ctx, p.ID)
	return p, blocks, err
}

func (s *Service) UpdatePage(ctx context.Context, id, title, slug, status string) (*domain.Page, error) {
	p, err := s.repo.FindPageByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if title = strings.TrimSpace(title); title != "" {
		p.Title = title
	}
	if slug = strings.TrimSpace(slug); slug != "" {
		p.Slug = slug
	}
	if status = strings.TrimSpace(status); status != "" {
		p.Status = status
	}
	p.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdatePage(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) UpsertMenu(ctx context.Context, id, name, location, itemsJSON string) (*domain.Menu, error) {
	now := time.Now().UTC()
	if existing, err := s.repo.FindMenuByName(ctx, name); err == nil && existing != nil {
		existing.ItemsJSON = itemsJSON
		existing.Location = location
		existing.UpdatedAt = now
		return existing, s.repo.UpsertMenu(ctx, existing)
	}
	menuID := id
	if menuID == "" {
		menuID = uuid.NewString()
	}
	m := &domain.Menu{ID: menuID, Name: name, Location: location, ItemsJSON: itemsJSON, CreatedAt: now, UpdatedAt: now}
	return m, s.repo.UpsertMenu(ctx, m)
}

func (s *Service) ListMenus(ctx context.Context) ([]domain.Menu, error) { return s.repo.ListMenus(ctx) }

func (s *Service) UpsertSEO(ctx context.Context, entityType, entityID, title, description, keywords string) (*domain.SEOMeta, error) {
	meta := &domain.SEOMeta{ID: uuid.NewString(), EntityType: entityType, EntityID: entityID, Title: title, Description: description, Keywords: keywords}
	return meta, s.repo.UpsertSEO(ctx, meta)
}

func (s *Service) GetSEO(ctx context.Context, entityType, entityID string) (*domain.SEOMeta, error) {
	meta, err := s.repo.FindSEO(ctx, entityType, entityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return meta, nil
}

func (s *Service) UpsertFAQ(ctx context.Context, question, answer string, sort int) (*domain.FAQ, error) {
	f := &domain.FAQ{ID: uuid.NewString(), Question: question, Answer: answer, SortOrder: sort, IsActive: true}
	return f, s.repo.UpsertFAQ(ctx, f)
}

func (s *Service) ListFAQs(ctx context.Context) ([]domain.FAQ, error) { return s.repo.ListFAQs(ctx) }
