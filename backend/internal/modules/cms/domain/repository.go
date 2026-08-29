package domain

import "context"

type Repository interface {
	CreatePage(ctx context.Context, p *Page) error
	UpdatePage(ctx context.Context, p *Page) error
	FindPageBySlug(ctx context.Context, slug string) (*Page, error)
	FindPageByID(ctx context.Context, id string) (*Page, error)
	ListPages(ctx context.Context) ([]Page, error)
	SaveBlocks(ctx context.Context, pageID string, blocks []Block) error
	ListBlocks(ctx context.Context, pageID string) ([]Block, error)
	UpsertMenu(ctx context.Context, m *Menu) error
	FindMenuByName(ctx context.Context, name string) (*Menu, error)
	ListMenus(ctx context.Context) ([]Menu, error)
	UpsertSEO(ctx context.Context, s *SEOMeta) error
	FindSEO(ctx context.Context, entityType, entityID string) (*SEOMeta, error)
	UpsertFAQ(ctx context.Context, f *FAQ) error
	ListFAQs(ctx context.Context) ([]FAQ, error)
}
