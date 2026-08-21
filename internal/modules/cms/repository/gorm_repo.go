package repository

import (
	"context"
	"time"

	"crm-starterkit/internal/modules/cms/domain"

	"gorm.io/gorm"
)

type PageModel struct {
	ID, Slug, Title, Status string
	CreatedAt, UpdatedAt    time.Time
}

func (PageModel) TableName() string { return "cms_pages" }

type BlockModel struct {
	ID, PageID, Type, Content string
	SortOrder                 int
}

func (BlockModel) TableName() string { return "cms_blocks" }

type MenuModel struct {
	ID, Name, Location, ItemsJSON string
	CreatedAt, UpdatedAt          time.Time
}

func (MenuModel) TableName() string { return "cms_menus" }

type SEOModel struct {
	ID, EntityType, EntityID, Title, Description, Keywords string
}

func (SEOModel) TableName() string { return "cms_seo_meta" }

type FAQModel struct {
	ID, Question, Answer string
	SortOrder            int
	IsActive             bool
}

func (FAQModel) TableName() string { return "cms_faqs" }

type Store struct{ db *gorm.DB }

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

func (s *Store) AutoMigrate() error {
	return s.db.AutoMigrate(&PageModel{}, &BlockModel{}, &MenuModel{}, &SEOModel{}, &FAQModel{})
}

func (s *Store) CreatePage(ctx context.Context, p *domain.Page) error {
	return s.db.WithContext(ctx).Create(&PageModel{ID: p.ID, Slug: p.Slug, Title: p.Title, Status: p.Status, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}).Error
}

func (s *Store) UpdatePage(ctx context.Context, p *domain.Page) error {
	return s.db.WithContext(ctx).Save(&PageModel{ID: p.ID, Slug: p.Slug, Title: p.Title, Status: p.Status, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}).Error
}

func (s *Store) FindPageBySlug(ctx context.Context, slug string) (*domain.Page, error) {
	var m PageModel
	if err := s.db.WithContext(ctx).First(&m, "slug = ?", slug).Error; err != nil {
		return nil, err
	}
	return &domain.Page{ID: m.ID, Slug: m.Slug, Title: m.Title, Status: m.Status, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}, nil
}

func (s *Store) ListPages(ctx context.Context) ([]domain.Page, error) {
	var rows []PageModel
	if err := s.db.WithContext(ctx).Order("created_at desc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Page, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.Page{ID: r.ID, Slug: r.Slug, Title: r.Title, Status: r.Status, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt})
	}
	return out, nil
}

func (s *Store) SaveBlocks(ctx context.Context, pageID string, blocks []domain.Block) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("page_id = ?", pageID).Delete(&BlockModel{}).Error; err != nil {
			return err
		}
		for _, b := range blocks {
			if err := tx.Create(&BlockModel{ID: b.ID, PageID: pageID, Type: b.Type, Content: b.Content, SortOrder: b.SortOrder}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) ListBlocks(ctx context.Context, pageID string) ([]domain.Block, error) {
	var rows []BlockModel
	if err := s.db.WithContext(ctx).Where("page_id = ?", pageID).Order("sort_order asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Block, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.Block{ID: r.ID, PageID: r.PageID, Type: r.Type, Content: r.Content, SortOrder: r.SortOrder})
	}
	return out, nil
}

func (s *Store) UpsertMenu(ctx context.Context, m *domain.Menu) error {
	return s.db.WithContext(ctx).Save(&MenuModel{ID: m.ID, Name: m.Name, Location: m.Location, ItemsJSON: m.ItemsJSON, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}).Error
}

func (s *Store) ListMenus(ctx context.Context) ([]domain.Menu, error) {
	var rows []MenuModel
	if err := s.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Menu, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.Menu{ID: r.ID, Name: r.Name, Location: r.Location, ItemsJSON: r.ItemsJSON, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt})
	}
	return out, nil
}

func (s *Store) UpsertSEO(ctx context.Context, meta *domain.SEOMeta) error {
	return s.db.WithContext(ctx).Save(&SEOModel{ID: meta.ID, EntityType: meta.EntityType, EntityID: meta.EntityID, Title: meta.Title, Description: meta.Description, Keywords: meta.Keywords}).Error
}

func (s *Store) FindSEO(ctx context.Context, entityType, entityID string) (*domain.SEOMeta, error) {
	var m SEOModel
	if err := s.db.WithContext(ctx).First(&m, "entity_type = ? AND entity_id = ?", entityType, entityID).Error; err != nil {
		return nil, err
	}
	return &domain.SEOMeta{ID: m.ID, EntityType: m.EntityType, EntityID: m.EntityID, Title: m.Title, Description: m.Description, Keywords: m.Keywords}, nil
}

func (s *Store) UpsertFAQ(ctx context.Context, f *domain.FAQ) error {
	return s.db.WithContext(ctx).Save(&FAQModel{ID: f.ID, Question: f.Question, Answer: f.Answer, SortOrder: f.SortOrder, IsActive: f.IsActive}).Error
}

func (s *Store) ListFAQs(ctx context.Context) ([]domain.FAQ, error) {
	var rows []FAQModel
	if err := s.db.WithContext(ctx).Where("is_active = ?", true).Order("sort_order asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]domain.FAQ, 0, len(rows))
	for _, r := range rows {
		out = append(out, domain.FAQ{ID: r.ID, Question: r.Question, Answer: r.Answer, SortOrder: r.SortOrder, IsActive: r.IsActive})
	}
	return out, nil
}
